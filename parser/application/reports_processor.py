import logging
from typing import cast

from infra.e_disclosure import EDisclosureClient
from infra.s3_storage import S3ReportsStorage
from infra.db_repo import ReportsRepository
from infra.models import ReportORM
from companies import get_ticker_by_inn
from domain.processing_result import ProcessingResult, ProcessingError, ReportMetadata, SingleCompanyProcessingResult
from domain.report import DownloadedReport
from parser.domain.ports import VectorizationService

logger = logging.getLogger(__name__)


class ReportProcessor:
    def __init__(
        self,
        s3_client: S3ReportsStorage,
        repo: ReportsRepository,
        vectorization_service: VectorizationService,
    ):
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
                    result = self.process_company_by_inn(client, inn, skip_indexing=skip_indexing)
                    results.processed += 1
                    results.saved += result.saved
                except Exception as e:
                    logger.error(f"error while processing {inn}: {e}")
                    results.errors.append(ProcessingError(inn=inn, error=str(e)))

        return results

    def process_company_by_inn(self, client: EDisclosureClient, inn: str, skip_indexing: bool = False) -> SingleCompanyProcessingResult:
        ticker = get_ticker_by_inn(inn)
        return self.process_company_by_query(
            client,
            query=inn,
            ticker=ticker,
            skip_indexing=skip_indexing,
        )

    def process_company_by_query(
        self,
        client: EDisclosureClient,
        query: str,
        ticker: str,
        skip_indexing: bool = False,
    ) -> SingleCompanyProcessingResult:
        companies = client.search_company(query)

        if not companies:
            logger.warning(f"company not found: {query}")
            raise ValueError(f"company not found: {query}")

        first_company = companies[0]
        logger.info(f"selected company: {first_company['name']} (ID: {first_company['id']}, {ticker})")

        logger.info("downloading reports for %s...", ticker)
        reports = client.download_reports(first_company)
        downloaded = self._convert_to_downloaded_reports(reports)

        saved = self.process_downloaded_reports(downloaded, ticker, skip_indexing=skip_indexing)

        self._log_results(len(downloaded), len(saved))

        return SingleCompanyProcessingResult(ticker=ticker, saved=len(saved), reports_metadata=saved)

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

    def process_downloaded_reports(self, reports: list[DownloadedReport], ticker: str, skip_indexing: bool = False) -> list[ReportMetadata]:
        saved : list[ReportMetadata] = []
        for report in reports:
            meta = self._process_downloaded_report(report, ticker, skip_indexing=skip_indexing)
            if meta is not None:
                saved.append(meta)
        return saved

    def _process_downloaded_report(self, report: DownloadedReport, ticker: str, skip_indexing: bool = False) -> ReportMetadata | None:
        if not report.is_valid():
            logger.warning(f"skipping: report {report} is not valid")
            return None

        logger.info(f"{report.year}, {report.period_months} months - {report.size_mb:.2f} MB")
        logger.info(f"local path: {report.path}")

        s3_path = self._ensure_s3_uploaded(ticker, report)
        if s3_path is None:
            logger.error("error while saving to s3")
            return None

        report_orm = self._save_report_to_db(ticker, report.year, report.period, s3_path)
        if not report_orm:
            return None

        if not skip_indexing:
            self._vectorize_report(report_orm, report, ticker)

        return ReportMetadata(s3_path=s3_path, year=report.year, period=report.period)

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
        logger.info(f"total file processed: {downloaded_count}")
        logger.info(f"{saved} reports saved")

    def _ensure_s3_uploaded(self, ticker, report) -> str | None:
        existing_s3_path = self.s3_client.get_s3_report_link(ticker, report.year, report.period)
        if existing_s3_path:
            logger.info(f"skipping: report already exists in S3: {existing_s3_path}")
            return existing_s3_path

        return self._upload_to_s3(ticker, report)