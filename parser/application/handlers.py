import logging
from fastapi import APIRouter, Depends, BackgroundTasks
from sqlalchemy.orm import Session

from application.parser import ReportProcessor
from companies import COMPANIES
from infra.database import get_db, get_db_session
from infra.db_repo import ReportsRepository
from infra.s3_storage import S3ReportsStorage

logger = logging.getLogger(__name__)

router = APIRouter()


def get_reports_repo(db: Session = Depends(get_db)) -> ReportsRepository:
    return ReportsRepository(db)


@router.get("/health")
def health():
    return {"status": "ok"}


@router.post("/start_parsing")
def start_parsing(background_tasks: BackgroundTasks):
    """Запустить парсинг в фоне."""
    background_tasks.add_task(run_parsing)
    return {"message": "Parsing started in background"}


def run_parsing():
    """Фоновая задача парсинга."""
    with get_db_session() as db:
        repo = ReportsRepository(db)
        s3_client = S3ReportsStorage()
        # processor = ReportProcessor(s3_client, repo)
        # results = processor.process_companies(COMPANIES)
        # logger.info(f"Parsing completed: {results}")


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
