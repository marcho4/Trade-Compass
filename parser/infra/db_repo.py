import logging
from typing import List, Optional
from sqlalchemy.orm import Session
from sqlalchemy.exc import IntegrityError
from domain.report import Report

logger = logging.getLogger(__name__)


class ReportsRepository:
    def __init__(self, session: Session):
        self.session = session

    def create_report(self, ticker: str, year: int, period: str, s3_path: str) -> Optional[Report]:
        report = Report(
            ticker=ticker,
            year=year,
            period=period,
            s3_path=s3_path
        )
        try:
            self.session.add(report)
            self.session.commit()
            self.session.refresh(report)
            return report
        except IntegrityError:
            self.session.rollback()
            logger.warning(f"Отчёт уже существует: {ticker} {year} {period}")
            return None

    def get_report_by_id(self, report_id: int) -> Optional[Report]:
        return self.session.query(Report).filter(Report.id == report_id).first()

    def get_reports_by_ticker(self, ticker: str) -> List[Report]:
        return self.session.query(Report).filter(Report.ticker == ticker).all()

    def get_reports_by_year(self, year: int) -> List[Report]:
        return self.session.query(Report).filter(Report.year == year).all()

    def get_report_by_params(self, ticker: str, year: int, period: str) -> Optional[Report]:
        return self.session.query(Report).filter(
            Report.ticker == ticker,
            Report.year == year,
            Report.period == period
        ).first()

    def get_all_reports(self, skip: int = 0, limit: int = 100) -> List[Report]:
        return self.session.query(Report).offset(skip).limit(limit).all()

    def update_report_s3_path(self, report_id: int, new_s3_path: str) -> Optional[Report]:
        report = self.get_report_by_id(report_id)
        if report:
            report.s3_path = new_s3_path
            self.session.commit()
            self.session.refresh(report)
        return report

    def delete_report(self, report_id: int) -> bool:
        report = self.get_report_by_id(report_id)
        if report:
            self.session.delete(report)
            self.session.commit()
            return True
        return False

    def report_exists(self, ticker: str, year: int, period: str) -> bool:
        return self.get_report_by_params(ticker, year, period) is not None
