import logging
import os
import re
import time
import requests
from urllib.parse import urljoin
from bs4 import BeautifulSoup
from infra.config import config
from infra.e_disclosure.metadata_parser import ReportMetadataParser

logger = logging.getLogger(__name__)


class ReportDownloader:
    """Загрузчик отчетов с e-disclosure.ru."""

    def __init__(self, driver, metadata_parser: ReportMetadataParser = None):
        self.driver = driver
        self.base_url = config.base_url
        self.metadata_parser = metadata_parser or ReportMetadataParser()

    def download_reports(
        self, company: dict, download_dir: str = None
    ) -> list[dict]:
        """Скачать отчеты компании."""
        download_dir = download_dir or config.download_dir
        os.makedirs(download_dir, exist_ok=True)

        if not self.driver:
            logger.error("WebDriver не инициализирован")
            return []

        company_id = company["id"]
        company_name = company["name"]
        company_element = company.get("element")

        try:
            logger.info("=== Переход на страницу компании ===")
            logger.info(f"Компания: {company_name} (ID: {company_id})")

            self._navigate_to_company(company_element, company.get("url"))

            reports_url = f"{self.base_url}/portal/files.aspx?id={company_id}&type=5"
            logger.info("=== Переход к отчетности эмитента ===")
            logger.info(f"URL: {reports_url}")

            self.driver.get(reports_url)
            time.sleep(config.timeout_page_load)
            logger.info("Страница с отчетностью эмитента загружена")

            soup = BeautifulSoup(self.driver.page_source, "html.parser")
            metadata_list = self.metadata_parser.parse_table(soup)

            reports = []
            for i, metadata in enumerate(metadata_list[:config.max_reports_per_company], 1):
                report = self._download_single_report(
                    metadata, company_name, download_dir, i
                )
                reports.append(report)

                if i < len(metadata_list):
                    time.sleep(config.timeout_between_files)

            self._log_downloaded_files(download_dir)
            return reports

        except Exception as e:
            logger.error(f"Ошибка при получении отчетов: {e}", exc_info=True)
            return []

    def _navigate_to_company(self, element, url: str):
        """Переход на страницу компании."""
        if element:
            try:
                element.click()
                logger.debug("Клик по ссылке компании выполнен")
            except Exception as e:
                logger.warning(f"Не удалось кликнуть по элементу: {e}")
                self.driver.get(url)
        else:
            self.driver.get(url)

        time.sleep(config.timeout_page_load)
        logger.info("Страница компании загружена")

    def _download_single_report(
        self, metadata: dict, company_name: str, download_dir: str, index: int
    ) -> dict:
        """Скачать один отчет."""
        file_url = metadata["file_url"]
        if not file_url.startswith("http"):
            file_url = urljoin(self.base_url, file_url)

        period = metadata.get("reporting_period", "unknown")
        document_type = metadata.get("document_type", "")
        file_id = metadata.get("file_id", "")
        file_info = metadata.get("file_info", "")

        file_name = self._generate_filename(company_name, period, file_id)

        logger.info(f"\n{index}. Отчетный период: {period}")
        logger.info(f"   Тип документа: {document_type}")
        logger.info(f"   Дата публикации: {metadata.get('publication_date', '')}")
        logger.info(f"   Файл: {file_info}")
        logger.debug(f"   URL: {file_url}")

        try:
            file_path = os.path.join(download_dir, file_name)
            logger.info(f"   Скачивание в: {file_path}")

            time.sleep(config.timeout_between_requests)

            response = requests.get(file_url, timeout=config.timeout_download)
            response.raise_for_status()

            with open(file_path, "wb") as f:
                f.write(response.content)

            file_size = len(response.content)
            logger.info(f"   Файл скачан: {file_size} bytes")

            return {
                "name": file_name,
                "period": period,
                "document_type": document_type,
                "publication_date": metadata.get("publication_date", ""),
                "base_date": metadata.get("base_date", ""),
                "url": file_url,
                "status": "downloaded",
                "path": file_path,
                "size": file_size,
                "metadata": metadata,
            }

        except Exception as e:
            logger.error(f"   Ошибка при скачивании: {e}")
            return {
                "name": file_name,
                "period": period,
                "document_type": document_type,
                "url": file_url,
                "status": "error",
                "metadata": metadata,
            }

    def _generate_filename(self, company_name: str, period: str, file_id: str) -> str:
        """Генерация имени файла."""
        file_name = f"{company_name}_{period}_{file_id}.zip"
        file_name = re.sub(r"[^\w\s.-]", "_", file_name)
        return re.sub(r"_+", "_", file_name)

    def _log_downloaded_files(self, download_dir: str):
        """Логирование скачанных файлов."""
        if os.path.exists(download_dir):
            downloaded_files = [
                f for f in os.listdir(download_dir) if not f.startswith(".")
            ]
            logger.info("=" * 60)
            logger.info(f"Всего файлов в {download_dir}: {len(downloaded_files)}")
            logger.info("=" * 60)
            for file in downloaded_files:
                file_size = os.path.getsize(os.path.join(download_dir, file))
                logger.info(f"  - {file} ({file_size:,} bytes)")
