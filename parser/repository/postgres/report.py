import logging

from typing import List, Optional
from sqlalchemy.orm import Session
from sqlalchemy.exc import IntegrityError
from models import ReportORM

from domain.model.report import ReportEntity
from usecase.interfaces import ReportsRepository

logger = logging.getLogger(__name__)


class PostgresReportsRepository(ReportsRepository):
    def __init__(self, session: Session):
        self.session = session

    def create_report(self, ticker: str, year: int, period: str, s3_path: str) -> Optional[ReportEntity]:
        report = ReportORM(
            ticker=ticker,
            year=year,
            period=period,
            s3_path=s3_path
        )
        try:
            self.session.add(report)
            self.session.commit()
            self.session.refresh(report)
            return self._map_orm_to_domain(report)
        except IntegrityError:
            self.session.rollback()
            logger.warning(f"Отчёт уже существует: {ticker} {year} {period}")
            return None

    def get_report_by_id(self, report_id: int) -> Optional[ReportEntity]:
        orm = self.session.query(ReportORM).filter(ReportORM.id == report_id).first()
        return self._map_orm_to_domain(orm) if orm else None

    def get_reports_by_ticker(self, ticker: str) -> List[ReportEntity]:
        rows = self.session.query(ReportORM).filter(ReportORM.ticker == ticker).all()
        return [self._map_orm_to_domain(r) for r in rows]

    def get_reports_by_year(self, year: int) -> List[ReportEntity]:
        rows = self.session.query(ReportORM).filter(ReportORM.year == year).all()
        return [self._map_orm_to_domain(r) for r in rows]

    def get_report_by_params(self, ticker: str, year: int, period: str) -> Optional[ReportEntity]:
        orm = self.session.query(ReportORM).filter(
            ReportORM.ticker == ticker,
            ReportORM.year == year,
            ReportORM.period == period
        ).first()
        return self._map_orm_to_domain(orm) if orm else None

    def get_all_reports(self, skip: int = 0, limit: int = 100) -> List[ReportEntity]:
        rows = self.session.query(ReportORM).offset(skip).limit(limit).all()
        return [self._map_orm_to_domain(r) for r in rows]

    def update_report_s3_path(self, report_id: int, new_s3_path: str) -> Optional[ReportEntity]:
        orm = self.session.query(ReportORM).filter(ReportORM.id == report_id).first()
        if orm:
            orm.s3_path = new_s3_path  # type: ignore[assignment]
            self.session.commit()
            self.session.refresh(orm)
            return self._map_orm_to_domain(orm)
        return None

    def delete_report(self, report_id: int) -> bool:
        orm = self.session.query(ReportORM).filter(ReportORM.id == report_id).first()
        if orm:
            self.session.delete(orm)
            self.session.commit()
            return True
        return False

    def report_exists(self, ticker: str, year: int, period: str) -> bool:
        return self.get_report_by_params(ticker, year, period) is not None

    def get_latest_report(self, ticker: str) -> Optional[ReportEntity]:
        orm = (
            self.session.query(ReportORM)
            .filter(ReportORM.ticker == ticker)
            .order_by(ReportORM.year.desc(), ReportORM.period.desc())
            .first()
        )
        return self._map_orm_to_domain(orm) if orm else None

    def _map_orm_to_domain(self, orm: ReportORM) -> ReportEntity:
        return ReportEntity(
            ticker=str(orm.ticker),
            year=int(orm.year),  # type: ignore[arg-type]
            period=str(orm.period),
            s3_path=str(orm.s3_path),
            id=int(orm.id) if orm.id is not None else None,  # type: ignore[arg-type]
        )
