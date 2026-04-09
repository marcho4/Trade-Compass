from abc import ABC, abstractmethod
from typing import Any, Dict, List, Optional

from domain.model.report import ReportEntity

class VectorizationService(ABC):
    @abstractmethod
    def vectorize_report(self, report_id: int, file_path: str, ticker: str, year: int, period: str) -> Dict[str, Any]:
        ...

class ColdStorage(ABC):
    @abstractmethod
    def get_report_link(self, ticker: str, year: int, period: str, extension: str = ".pdf") -> Optional[str]:
        ...

    @abstractmethod
    def upload_report(self, ticker: str, year: int, period: str, file_path: str) -> Optional[str]:
        ...

class ReportsRepository(ABC):
    @abstractmethod
    def create_report(self, ticker: str, year: int, period: str, s3_path: str) -> Optional[ReportEntity]:
        ...

    @abstractmethod
    def get_report_by_id(self, report_id: int) -> Optional[ReportEntity]:
        ...

    @abstractmethod
    def get_reports_by_ticker(self, ticker: str) -> List[ReportEntity]:
        ...

    @abstractmethod
    def get_reports_by_year(self, year: int) -> List[ReportEntity]:
        ...
    
    @abstractmethod
    def get_report_by_params(self, ticker: str, year: int, period: str) -> Optional[ReportEntity]:
        ...

    @abstractmethod
    def get_all_reports(self, skip: int = 0, limit: int = 100) -> List[ReportEntity]:
        ...
    
    @abstractmethod
    def update_report_s3_path(self, report_id: int, new_s3_path: str) -> Optional[ReportEntity]:
        ...

    @abstractmethod
    def delete_report(self, report_id: int) -> bool:
        ...

    @abstractmethod
    def report_exists(self, ticker: str, year: int, period: str) -> bool:
        ...

    @abstractmethod
    def get_latest_report(self, ticker: str) -> Optional[ReportEntity]:
        ...

class AnalyzeTaskGateway(ABC):
    @abstractmethod
    def send_task(self, ticker: str, year: int, period: str, report_url: str, task_id: str, task_type: str) -> None: ...

class ReportsParser(ABC):
    @abstractmethod
    def search_company(self, query: str) -> List[dict]:
        ...

    @abstractmethod
    def download_reports(self, company: dict, download_dir: str = "") -> list[dict]:
        ...

    @abstractmethod
    def unzip_files(self, source_dir: str = "", target_dir: str = "") -> list[str]:
        ...