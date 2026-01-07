import logging
import os
import time
import requests

from application.utils import extract_year_and_period, normalize_filename

from bs4 import BeautifulSoup

from domain.report_metadata import ReportMetadata

from infra.config import config
from infra.e_disclosure.metadata_parser import ReportMetadataParser
from infra.e_disclosure.unzipper import FileUnzipper

from urllib.parse import urljoin

logger = logging.getLogger(__name__)


class ReportDownloader:
    ALLOWED_DOCUMENT_TYPES = ["промежуточная", "годовая"]

    def __init__(self, driver, metadata_parser: ReportMetadataParser):
        self.driver = driver
        self.base_url = config.base_url
        self.metadata_parser = metadata_parser or ReportMetadataParser()

    def _is_allowed_document_type(self, document_type: str) -> bool:
        if not document_type:
            return False
        doc_type_lower = document_type.lower()
        return any(allowed in doc_type_lower for allowed in self.ALLOWED_DOCUMENT_TYPES)

    def download_reports(
        self, company: dict, download_dir: str = None
    ) -> list[dict]:
        download_dir = download_dir or config.download_dir
        os.makedirs(download_dir, exist_ok=True)

        if not self.driver:
            logger.error("WebDriver не инициализирован")
            return []

        company_id = company["id"]
        company_name = company["name"]
        company_element = company.get("element")

        try:
            self._navigate_to_company(company_element, company.get("url"))

            reports_url = f"{self.base_url}/portal/files.aspx?id={company_id}&type=4"
            self.driver.get(reports_url)
            time.sleep(config.timeout_page_load)

            soup = BeautifulSoup(self.driver.page_source, "html.parser")
            filtered_metadata: list[ReportMetadata] = list(
                filter(
                    lambda x: self._is_allowed_document_type(x.document_type), 
                    self.metadata_parser.parse_table(soup),
                )
            )

            reports = []
            for i, metadata in enumerate(filtered_metadata[:config.max_reports_per_company], 1):
                report = self._download_single_report(
                    metadata, company_name, download_dir
                )
                reports.append(report)

                if i < len(filtered_metadata[:config.max_reports_per_company]):
                    time.sleep(config.timeout_between_files)

            return reports

        except Exception as e:
            logger.error(f"Ошибка при получении отчетов: {e}", exc_info=True)
            return []

    def _navigate_to_company(self, element, url: str):
        if element:
            try:
                element.click()
            except Exception as e:
                logger.warning(f"Не удалось кликнуть по элементу: {e}")
                self.driver.get(url)
        else:
            self.driver.get(url)

        time.sleep(config.timeout_page_load)

    def _download_single_report(
        self, metadata: ReportMetadata, company_name: str, download_dir: str
    ) -> dict:
        file_url = metadata.file_url
        if not file_url.startswith("http"):
            file_url = urljoin(self.base_url, file_url)

        document_type = metadata.document_type
        year, period_months = extract_year_and_period(metadata.reporting_period)

        file_name = self._generate_filename(company_name, year, period_months)

        try:
            file_path = os.path.join(download_dir, file_name)

            time.sleep(config.timeout_between_requests)

            response = requests.get(file_url, timeout=config.timeout_download)
            response.raise_for_status()

            with open(file_path, "wb") as f:
                f.write(response.content)

            unzipped_path = FileUnzipper.unzip_and_rename(file_path)
            if not unzipped_path:
                return {
                    "name": file_name,
                    "year": year,
                    "period_months": period_months,
                    "document_type": document_type,
                    "url": file_url,
                    "status": "error",
                    "metadata": metadata,
                }

            file_size = len(response.content)

            return {
                "name": file_name,
                "year": year,
                "period_months": period_months,
                "document_type": document_type,
                "publication_date": metadata.publication_date,
                "base_date": metadata.base_date,
                "url": file_url,
                "status": "downloaded",
                "path": unzipped_path,
                "size": file_size,
                "metadata": metadata,
            }

        except Exception as e:
            return {
                "name": file_name,
                "year": year,
                "period_months": period_months,
                "document_type": document_type,
                "url": file_url,
                "status": "error",
                "metadata": metadata,
            }

    def _generate_filename(self, company_name: str, year: int, period_months: int) -> str:
        file_name = f"{company_name}_{year}_{period_months}M.zip"
        return normalize_filename(file_name)
