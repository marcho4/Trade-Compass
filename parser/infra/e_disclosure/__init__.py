from infra.e_disclosure.driver import SeleniumDriver
from infra.e_disclosure.searcher import CompanySearcher
from infra.e_disclosure.downloader import ReportDownloader
from infra.e_disclosure.metadata_parser import ReportMetadataParser
from infra.e_disclosure.unzipper import FileUnzipper


class EDisclosureClient:
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
        return self._searcher.search(query)

    def get_reports(self, company: dict, download_dir: str = None) -> list[dict]:
        return self._downloader.download_reports(company, download_dir)

    def unzip_files(self, source_dir: str = None, target_dir: str = None) -> list[str]:
        return self._unzipper.unzip_all(source_dir, target_dir)


__all__ = [
    "SeleniumDriver",
    "CompanySearcher",
    "ReportDownloader",
    "ReportMetadataParser",
    "FileUnzipper",
    "EDisclosureClient",
]
