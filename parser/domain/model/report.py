from dataclasses import dataclass
from typing import Optional


BYTES_IN_MB = 1024 * 1024


@dataclass
class DownloadedReport:
    path: str
    size: int
    status: str
    year: Optional[int] = None
    period_months: Optional[int] = None
    name: Optional[str] = None

    @property
    def size_mb(self) -> float:
        return self.size / BYTES_IN_MB

    @property
    def period(self) -> str:
        return str(self.period_months)

    def is_valid(self) -> bool:
        return self.year is not None and self.period_months is not None


@dataclass
class ReportEntity:
    ticker: str
    year: int
    period: str
    s3_path: str
    id: Optional[int] = None

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "ticker": self.ticker,
            "year": self.year,
            "period": self.period,
            "s3_path": self.s3_path,
        }
