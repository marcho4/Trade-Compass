from sqlalchemy import Column, Integer, String, CheckConstraint, UniqueConstraint
from infra.database import Base


class ReportORM(Base):
    __tablename__ = "reports"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    ticker = Column(String(50), nullable=False, index=True)
    year = Column(Integer, nullable=False)
    period = Column(String(50), nullable=False)
    s3_path = Column(String(500), nullable=False)

    __table_args__ = (
        CheckConstraint("year >= 2000 AND year <= 2100", name="check_year_range"),
        UniqueConstraint("ticker", "year", "period", name="uq_report_ticker_year_period"),
    )

    def __repr__(self):
        return f"<ReportORM(id={self.id}, ticker={self.ticker}, year={self.year}, period={self.period})>"

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "ticker": self.ticker,
            "year": self.year,
            "period": self.period,
            "s3_path": self.s3_path,
        }
