from dataclasses import dataclass
from typing import Optional


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
