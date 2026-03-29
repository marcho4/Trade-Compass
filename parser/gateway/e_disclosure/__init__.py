from .driver import SeleniumDriver
from .searcher import CompanySearcher
from .downloader import ReportDownloader
from .metadata_parser import ReportMetadataParser
from unzipper import FileUnzipper
from usecase.interfaces import ReportsParser


class EDisclosureClient(ReportsParser):
    def __init__(self):
        self._driver_manager = SeleniumDriver()
        self._driver = None
        self._searcher = None
        self._downloader = None
        self._unzipper = FileUnzipper()

    def __enter__(self):
        self._driver = self._driver_manager.__enter__()
        metadata_parser = ReportMetadataParser()
        self._searcher = CompanySearcher(self._driver)
        self._downloader = ReportDownloader(self._driver, metadata_parser)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self._driver_manager.__exit__(exc_type, exc_val, exc_tb)

    def search_company(self, query: str) -> list[dict]:
        if self._searcher is None:
            return []
    
        return self._searcher.search(query)

    def download_reports(self, company: dict, download_dir: str = "") -> list[dict]:
        if self._downloader is None:
            return []

        return self._downloader.download_reports(company, download_dir)

    def unzip_files(self, source_dir: str = "", target_dir: str = "") -> list[str]:
        return self._unzipper.unzip_all(source_dir, target_dir)


__all__ = [
    "SeleniumDriver",
    "CompanySearcher",
    "ReportDownloader",
    "ReportMetadataParser",
    "FileUnzipper",
    "EDisclosureClient",
]
