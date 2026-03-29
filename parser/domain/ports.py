from abc import ABC, abstractmethod
from typing import Any, Dict

class VectorizationService(ABC):
    @abstractmethod
    def vectorize_report(self, report_id: int, file_path: str, ticker: str, year: int, period: str) -> Dict[str, Any]:
        ...
