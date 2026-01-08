import logging
from typing import cast

from infra.e_disclosure import EDisclosureClient
from infra.s3_storage import S3ReportsStorage
from infra.db_repo import ReportsRepository
from infra.models import ReportORM
from companies import get_ticker_by_inn
from application.vectorization_service import VectorizationService
from domain.processing_result import *

logger = logging.getLogger(__name__)


class ReportProcessor:
    def __init__(self, s3_client: S3ReportsStorage, repo: ReportsRepository, vectorization_service: VectorizationService):
        self.s3_client = s3_client
        self.repo = repo
        self.vectorization_service = vectorization_service

    def process_companies(self, companies_inn: list[str]) -> ProcessingResult:
        results = ProcessingResult(
            processed=0,
            errors=[],
            saved=0
        )

        with EDisclosureClient() as client:
            for inn in companies_inn:
                try:
                    result = self.process_company(client, inn)
                    results.processed += 1
                    results.saved += result.saved
                except Exception as e:
                    logger.error(f"Ошибка обработки {inn}: {e}")
                    results.errors.append(ProcessingError(inn=inn, error=str(e)))

        return results

    def process_company(self, client: EDisclosureClient, inn: str) -> SingleCompanyProcessingResult:
        companies = client.search_company(inn)

        if not companies:
            logger.warning(f"Компания не найдена: {inn}")
            raise ValueError(f"Компания не найдена: {inn}")

        first_company = companies[0]
        ticker = get_ticker_by_inn(inn)
        logger.info(f"Выбрана компания: {first_company['name']} (ID: {first_company['id']}, тикер: {ticker})")

        logger.info("Получение отчетности эмитента для %s", ticker)
        reports = client.get_reports(first_company)
        downloaded = [r for r in reports if r["status"] == "downloaded"]

        saved = self.upload_and_save_reports(downloaded, ticker)

        self.log_results(downloaded, saved)

        return SingleCompanyProcessingResult(ticker=ticker, saved=saved)

    def upload_and_save_reports(self, reports: list[dict], ticker: str) -> int:
        saved = 0

        for report in reports:
            size_mb = report.get("size", 0) / 1024 / 1024
            local_path = report["path"]

            year = report.get("year")
            period_months = report.get("period_months")

            if year is None or period_months is None:
                logger.warning(f"  Пропуск: отсутствуют year или period_months в данных отчета")
                continue

            logger.info(f"{year}, {period_months} месяцев - {size_mb:.2f} MB")
            logger.info(f"  Локальный путь: {local_path}")

            period = str(period_months)

            existing_s3_path = self.s3_client.get_s3_report_link(ticker, year, period)
            if existing_s3_path:
                logger.info(f"  Пропуск: файл уже существует в S3: {existing_s3_path}")
                continue

            s3_path = self.s3_client.upload_report(ticker, year, period, local_path)
            if s3_path is None:
                logger.error(f"  Ошибка загрузки в S3")
                continue
            logger.info(f"  Загружено в S3: {s3_path}")

            report_orm = self.save_report_to_db(ticker, year, period, s3_path)
            if report_orm:
                saved += 1

                self.vectorization_service.vectorize_report(
                    report_id=cast(int, report_orm.id),
                    file_path=local_path,
                    ticker=ticker,
                    year=year,
                    period=period,
                )
                
                    
        return saved

    def save_report_to_db(self, ticker: str, year: int, period: str, s3_path: str) -> ReportORM | None:
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

    def log_results(self, reports: list[dict], saved: int):
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
