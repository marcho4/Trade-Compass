from dataclasses import dataclass
from typing import List


@dataclass
class ProcessingError:
    error: str
    inn: str

@dataclass
class ProcessingResult:
    processed: int
    saved: int
    errors: List[ProcessingError]

@dataclass
class SingleCompanyProcessingResult:
    saved: int
    ticker: str

@dataclass
class VectorizationError:
    success: bool
    error: str
    report_id: int