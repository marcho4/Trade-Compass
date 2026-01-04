from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from application.parser import Parser
from companies import COMPANIES
from domain.report import report_to_dict
from infra.database import get_db, get_db_session
from infra.db_repo import ReportsRepository

router = APIRouter()


def get_reports_repo(db: Session = Depends(get_db)) -> ReportsRepository:
    return ReportsRepository(db)


@router.get("/health")
def health():
    return {"status": "ok"}  


@router.post("/start_parsing")
def start_parsing():
    with get_db_session() as db:
        repo = ReportsRepository(db)
        parser = Parser(repo)
        parser.run(COMPANIES)
    return {"message": "ok"}


@router.get("/reports")
def get_all_reports(
    skip: int = 0,
    limit: int = 100,
    repo: ReportsRepository = Depends(get_reports_repo)
):
    reports = repo.get_all_reports(skip=skip, limit=limit)
    return {"reports": [report_to_dict(r) for r in reports], "total": len(reports)}


@router.get("/reports/{ticker}")
def get_reports_by_ticker(
    ticker: str,
    repo: ReportsRepository = Depends(get_reports_repo)
):
    reports = repo.get_reports_by_ticker(ticker)
    return {"ticker": ticker, "reports": [report_to_dict(r) for r in reports], "total": len(reports)}


@router.post("/reports")
def create_report(
    ticker: str,
    year: int,
    period: str,
    s3_path: str,
    repo: ReportsRepository = Depends(get_reports_repo)
):
    report = repo.create_report(ticker, year, period, s3_path)
    if report is None:
        return {"error": "Report already exists", "ticker": ticker, "year": year, "period": period}
    return {"message": "Report created successfully", "report_id": report.id}
