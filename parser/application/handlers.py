import logging
import os
import tempfile
from fastapi import APIRouter, Depends, BackgroundTasks, UploadFile, File, Form, HTTPException, Header
from sqlalchemy.orm import Session

from application.reports_processor import ReportProcessor
from companies import COMPANIES
from infra.database import get_db, get_db_session
from infra.db_repo import ReportsRepository
from infra.s3_storage import S3ReportsStorage
from infra.config import config

logger = logging.getLogger(__name__)

router = APIRouter()


def get_reports_repo(db: Session = Depends(get_db)) -> ReportsRepository:
    return ReportsRepository(db)


def verify_admin_key(x_api_key: str = Header(..., alias="X-API-Key")):
    if not config.admin_api_key:
        raise HTTPException(
            status_code=500,
            detail="Admin API key not configured"
        )
    if x_api_key != config.admin_api_key:
        raise HTTPException(
            status_code=403,
            detail="Invalid API key"
        )
    return True


@router.get("/health")
def health():
    return {"status": "ok"}


@router.post("/start_parsing")
def start_parsing(
    background_tasks: BackgroundTasks,
    _: bool = Depends(verify_admin_key)
):
    background_tasks.add_task(run_parsing)
    return {"message": "Parsing started in background"}


def run_parsing():
    with get_db_session() as db:
        repo = ReportsRepository(db)
        s3_client = S3ReportsStorage()
        processor = ReportProcessor(s3_client, repo)
        results = processor.process_companies(list(COMPANIES.keys()))
        logger.info(f"Parsing completed: {results}")


@router.get("/reports")
def get_all_reports(
    skip: int = 0,
    limit: int = 100,
    repo: ReportsRepository = Depends(get_reports_repo),
):
    reports = repo.get_all_reports(skip=skip, limit=limit)
    return {"reports": [r.to_dict() for r in reports], "total": len(reports)}


@router.get("/reports/{ticker}")
def get_reports_by_ticker(
    ticker: str,
    repo: ReportsRepository = Depends(get_reports_repo),
):
    reports = repo.get_reports_by_ticker(ticker)
    return {"ticker": ticker, "reports": [r.to_dict() for r in reports], "total": len(reports)}


@router.post("/reports")
def create_report(
    ticker: str,
    year: int,
    period: str,
    s3_path: str,
    repo: ReportsRepository = Depends(get_reports_repo),
):
    report = repo.create_report(ticker, year, period, s3_path)
    if report is None:
        return {"error": "Report already exists", "ticker": ticker, "year": year, "period": period}
    return {"message": "Report created successfully", "report_id": report.id}


@router.post("/reports/upload")
async def upload_report(
    file: UploadFile = File(...),
    ticker: str = Form(...),
    year: int = Form(...),
    period: int = Form(..., description="Период в месяцах: 3, 6, 9 или 12"),
    repo: ReportsRepository = Depends(get_reports_repo),
):
    if period not in [3, 6, 9, 12]:
        raise HTTPException(
            status_code=400,
            detail=f"Некорректный период: {period}. Допустимые значения: 3, 6, 9, 12"
        )

    period_str = str(period)
    s3_client = S3ReportsStorage()

    existing_s3_path = s3_client.get_s3_report_link(ticker, year, period_str)
    if existing_s3_path:
        raise HTTPException(
            status_code=409,
            detail=f"Отчет уже существует: {existing_s3_path}"
        )

    file_extension = os.path.splitext(file.filename)[1] if file.filename else ".zip"

    with tempfile.NamedTemporaryFile(delete=False, suffix=file_extension) as tmp_file:
        content = await file.read()
        tmp_file.write(content)
        tmp_path = tmp_file.name

    try:
        s3_path = s3_client.upload_report(ticker, year, period_str, tmp_path)
        if s3_path is None:
            raise HTTPException(status_code=500, detail="Ошибка загрузки в S3")

        report = repo.create_report(ticker, year, period_str, s3_path)
        if report is None:
            raise HTTPException(
                status_code=409,
                detail=f"Отчет уже существует в БД: {ticker} {year} {period_str}"
            )

        logger.info(f"Отчет загружен вручную: {ticker} {year} {period_str} -> {s3_path}")

        return {
            "message": "Report uploaded successfully",
            "report_id": report.id,
            "s3_path": s3_path,
            "ticker": ticker,
            "year": year,
            "period": period_str
        }
    finally:
        if os.path.exists(tmp_path):
            os.unlink(tmp_path)
