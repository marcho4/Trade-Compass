import logging
from typing import cast

from infra.e_disclosure import EDisclosureClient
from infra.s3_storage import S3ReportsStorage
from infra.db_repo import ReportsRepository
from infra.models import ReportORM
from companies import get_ticker_by_inn
from application.vectorization_service import VectorizationService
from domain.processing_result import ProcessingResult, ProcessingError, SingleCompanyProcessingResult
from domain.report import DownloadedReport

logger = logging.getLogger(__name__)


class ReportProcessor:
    def __init__(self, s3_client: S3ReportsStorage, repo: ReportsRepository, vectorization_service: VectorizationService):
        self.s3_client = s3_client
        self.repo = repo
        self.vectorization_service = vectorization_service

    def process_companies(self, companies_inn: list[str], skip_indexing: bool = False) -> ProcessingResult:
        results = ProcessingResult(
            processed=0,
            errors=[],
            saved=0
        )

        with EDisclosureClient() as client:
            for inn in companies_inn:
                try:
                    result = self.process_company(client, inn, skip_indexing=skip_indexing)
                    results.processed += 1
                    results.saved += result.saved
                except Exception as e:
                    logger.error(f"Ошибка обработки {inn}: {e}")
                    results.errors.append(ProcessingError(inn=inn, error=str(e)))

        return results

    def process_company(self, client: EDisclosureClient, inn: str, skip_indexing: bool = False) -> SingleCompanyProcessingResult:
        companies = client.search_company(inn)

        if not companies:
            logger.warning(f"Компания не найдена: {inn}")
            raise ValueError(f"Компания не найдена: {inn}")

        first_company = companies[0]
        ticker = get_ticker_by_inn(inn)
        logger.info(f"Выбрана компания: {first_company['name']} (ID: {first_company['id']}, тикер: {ticker})")

        logger.info("Получение отчетности эмитента для %s", ticker)
        reports = client.get_reports(first_company)
        downloaded = self._convert_to_downloaded_reports(reports)

        saved = self.upload_and_save_reports(downloaded, ticker, skip_indexing=skip_indexing)

        self._log_results(len(downloaded), saved)

        return SingleCompanyProcessingResult(ticker=ticker, saved=saved)

    def _convert_to_downloaded_reports(self, reports: list[dict]) -> list[DownloadedReport]:
        result = []
        for r in reports:
            if r.get("status") == "downloaded":
                result.append(DownloadedReport(
                    path=r["path"],
                    size=r.get("size", 0),
                    status=r["status"],
                    year=r.get("year"),
                    period_months=r.get("period_months"),
                    name=r.get("name"),
                ))
        return result

    def upload_and_save_reports(self, reports: list[DownloadedReport], ticker: str, skip_indexing: bool = False) -> int:
        saved = 0
        for report in reports:
            if self._process_single_report(report, ticker, skip_indexing=skip_indexing):
                saved += 1
        return saved

    def _process_single_report(self, report: DownloadedReport, ticker: str, skip_indexing: bool = False) -> bool:
        if not report.is_valid():
            logger.warning("Пропуск: отсутствуют year или period_months в данных отчета")
            return False

        logger.info(f"{report.year}, {report.period_months} месяцев - {report.size_mb:.2f} MB")
        logger.info(f"  Локальный путь: {report.path}")

        if self._report_exists_in_s3(ticker, report):
            return False

        s3_path = self._upload_to_s3(ticker, report)
        if not s3_path:
            return False

        report_orm = self._save_report_to_db(ticker, report.year, report.period, s3_path)
        if not report_orm:
            return False

        if not skip_indexing:
            self._vectorize_report(report_orm, report, ticker)
        return True

    def _report_exists_in_s3(self, ticker: str, report: DownloadedReport) -> bool:
        existing_s3_path = self.s3_client.get_s3_report_link(ticker, report.year, report.period)
        if existing_s3_path:
            logger.info(f"  Пропуск: файл уже существует в S3: {existing_s3_path}")
            return True
        return False

    def _upload_to_s3(self, ticker: str, report: DownloadedReport) -> str | None:
        s3_path = self.s3_client.upload_report(ticker, report.year, report.period, report.path)
        if s3_path is None:
            logger.error("Ошибка загрузки в S3")
            return None
        logger.info(f"  Загружено в S3: {s3_path}")
        return s3_path

    def _vectorize_report(self, report_orm: ReportORM, report: DownloadedReport, ticker: str) -> None:
        result = self.vectorization_service.vectorize_report(
            report_id=cast(int, report_orm.id),
            file_path=report.path,
            ticker=ticker,
            year=report.year,
            period=report.period,
        )
        if not result.get("success"):
            logger.warning(f"Ошибка векторизации отчёта {ticker} {report.year}: {result.get('error')}")

    def _save_report_to_db(self, ticker: str, year: int, period: str, s3_path: str) -> ReportORM | None:
        try:
            report = self.repo.create_report(ticker, year, period, s3_path)
            if report is None:
                logger.info(f"Отчёт уже существует в БД: {ticker} {year} {period}")
                return None
            logger.info(f"Отчёт сохранён в БД: {ticker} {year} {period}")
            return report
        except Exception as e:
            logger.error(f"Ошибка сохранения отчёта в БД: {e}")
            return None

    def _log_results(self, downloaded_count: int, saved: int) -> None:
        logger.info(f"Всего файлов обработано: {downloaded_count}")
        logger.info("=" * 60)
        logger.info(f"Сохранено в БД: {saved} отчётов")
        logger.info("=" * 60)
