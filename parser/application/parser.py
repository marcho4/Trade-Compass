import logging
from infra.e_disclosure import EDisclosureClient
from infra.s3_storage import S3ReportsStorage
from infra.db_repo import ReportsRepository
from application.utils import extract_year_and_period
from application.exceptions import PeriodParseError
from companies import get_ticker_by_inn

logger = logging.getLogger(__name__)


class ReportProcessor:
    """Процессор для обработки отчетов компаний."""

    def __init__(self, 
                 s3_client: S3ReportsStorage, 
                 repo: ReportsRepository
                 ):
        self.s3_client = s3_client
        self.repo = repo

    def process_companies(self, companies_inn: list[str]) -> dict:
        """Обработать список компаний по ИНН."""
        results = {"processed": 0, "saved": 0, "errors": []}

        with EDisclosureClient() as client:
            for inn in companies_inn:
                try:
                    result = self._process_company(client, inn)
                    results["processed"] += 1
                    results["saved"] += result["saved"]
                except Exception as e:
                    logger.error(f"Ошибка обработки {inn}: {e}")
                    results["errors"].append({"inn": inn, "error": str(e)})

            client.unzip_files()

        return results

    def _process_company(self, client: EDisclosureClient, inn: str) -> dict:
        """Обработать одну компанию."""
        logger.info(f"Поиск компании... {inn}")
        companies = client.search_company(inn)

        if not companies:
            logger.warning(f"Компания не найдена: {inn}")
            raise ValueError(f"Компания не найдена: {inn}")

        logger.info(f"Найдено компаний: {len(companies)}")
        for i, company in enumerate(companies[:5], 1):
            logger.info(f"{i}. {company['name']} (ID: {company['id']})")

        first_company = companies[0]
        ticker = get_ticker_by_inn(inn)
        logger.info(f"Выбрана компания: {first_company['name']} (ID: {first_company['id']}, тикер: {ticker})")

        logger.info("Получение отчетности эмитента...")
        reports = client.get_reports(first_company)

        saved = self._upload_and_save_reports(reports, ticker)

        self._log_results(reports, saved)

        return {"company": first_company["name"], "saved": None}

    def _upload_and_save_reports(self, reports: list[dict], ticker: str) -> int:
        """Загрузить в S3 и сохранить в БД."""
        saved = 0
        downloaded = [r for r in reports if r["status"] == "downloaded"]

        logger.info("ОБРАБОТКА СКАЧАННЫХ ФАЙЛОВ:")
        for report in downloaded:
            size_mb = report.get("size", 0) / 1024 / 1024
            local_path = report["path"]
            period_raw = report["period"]

            logger.info(f"{period_raw} - {size_mb:.2f} MB")
            logger.info(f"  Локальный путь: {local_path}")

            try:
                year, period_months = extract_year_and_period(period_raw)
            except PeriodParseError as e:
                logger.warning(f"  Пропуск: {e}")
                continue

            period = str(period_months)

            existing_s3_path = self.s3_client.get_s3_report_link(ticker, year, period)
            if existing_s3_path:
                logger.info(f"  Пропуск: файл уже существует в S3: {existing_s3_path}")
                continue

            try:
                s3_path = self.s3_client.upload_report(ticker, year, period, local_path)
                if s3_path is None:
                    logger.error(f"  Ошибка загрузки в S3")
                    continue
                logger.info(f"  Загружено в S3: {s3_path}")

                if self._save_report_to_db(ticker, year, period, s3_path):
                    saved += 1
            except Exception as e:
                logger.error(f"  Ошибка загрузки в S3: {e}")

        return saved

    def _save_report_to_db(self, ticker: str, year: int, period: str, s3_path: str) -> bool:
        """Сохранить отчёт в базу данных."""
        try:
            report = self.repo.create_report(ticker, year, period, s3_path)
            if report is None:
                logger.info(f"Отчёт уже существует в БД: {ticker} {year} {period}")
                return False
            logger.info(f"Отчёт сохранён в БД: {ticker} {year} {period}")
            return True
        except Exception as e:
            logger.error(f"Ошибка сохранения отчёта в БД: {e}")
            return False

    def _log_results(self, reports: list[dict], saved: int):
        """Логирование результатов обработки."""
        downloaded = [r for r in reports if r["status"] == "downloaded"]
        errors = [r for r in reports if r["status"] == "error"]

        logger.info(f"Всего файлов обработано: {len(reports)}")
        logger.info(f"Успешно скачано: {len(downloaded)}")
        logger.info(f"Ошибок: {len(errors)}")

        if errors:
            logger.warning("ОШИБКИ:")
            for report in errors:
                logger.warning(f"  {report.get('name', 'Unknown')}")

        logger.info("=" * 60)
        logger.info(f"Сохранено в БД: {saved} отчётов")
        logger.info("=" * 60)
