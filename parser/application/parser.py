import logging
import re
from typing import Optional
from infra.s3_storage import S3ReportsStorage
from infra.selenium_client import EDisclosureParser
from infra.db_repo import ReportsRepository

logger = logging.getLogger(__name__)


class Parser:
    def __init__(self, repo: ReportsRepository) -> None:
        self.s3_client = S3ReportsStorage()
        self.parser = EDisclosureParser()
        self.repo = repo

    def _extract_year_from_period(self, period: str) -> Optional[int]:
        """
        Извлечь год из строки периода.
        
        Примеры входных данных:
        - "3 квартал 2024 года" → 2024
        - "2024 год" → 2024
        - "1 полугодие 2023" → 2023
        - "2023" → 2023
        """
        # Ищем 4-значное число, начинающееся с 20 (2000-2099)
        match = re.search(r'20\d{2}', period)
        if match:
            return int(match.group())
        
        logger.warning(f"Не удалось извлечь год из периода: {period}")
        return None

    def _save_report_to_db_with_year(self, ticker: str, year: int, period: str, s3_path: str) -> bool:
        """Сохранить отчёт в базу данных"""
        try:
            report = self.repo.create_report(ticker, year, period, s3_path)
            if report is None:
                # Отчёт уже существует (IntegrityError был обработан в репозитории)
                logger.info(f"Отчёт уже существует в БД: {ticker} {year} {period}")
                return False
            logger.info(f"Отчёт сохранён в БД: {ticker} {year} {period}")
            return True
        except Exception as e:
            logger.error(f"Ошибка сохранения отчёта в БД: {e}")
            return False

    def run(self, companies_inn: list[str]):
        for company_inn in companies_inn:
            logger.info(f"Поиск компании... {company_inn}")
            companies = self.parser.search_company(company_inn)

            if not companies:
                logger.warning("Компания не найдена")
                continue

            logger.info("=" * 60)
            logger.info(f"Найдено компаний: {len(companies)}")
            for i, company in enumerate(companies[:5], 1):
                logger.info(f"{i}. {company['name']} (ID: {company['id']})")
            logger.info("=" * 60)

            first_company = companies[0]
            ticker = first_company.get('ticker', company_inn)  # Используем ИНН как fallback
            logger.info(f"Выбрана компания: {first_company['name']} (ID: {first_company['id']})")

            logger.info("=" * 60)
            logger.info("Получение отчетности эмитента...")
            logger.info("=" * 60)

            reports = self.parser.get_issuer_reports_by_click(first_company, download_dir='./downloads')

            downloaded = [r for r in reports if r['status'] == 'downloaded']
            errors = [r for r in reports if r['status'] == 'error']

            logger.info("=" * 60)
            logger.info("ИТОГОВАЯ СТАТИСТИКА")
            logger.info("=" * 60)
            logger.info(f"Всего файлов обработано: {len(reports)}")
            logger.info(f"Успешно скачано: {len(downloaded)}")
            logger.info(f"Ошибок: {len(errors)}")
            logger.info("=" * 60)

            # Загружаем скачанные файлы в S3 и сохраняем в БД
            saved_count = 0
            if downloaded:
                logger.info("ОБРАБОТКА СКАЧАННЫХ ФАЙЛОВ:")
                for report in downloaded:
                    size_mb = report.get('size', 0) / 1024 / 1024
                    local_path = report['path']
                    period = report['period']

                    logger.info(f"{period} - {size_mb:.2f} MB")
                    logger.info(f"  Локальный путь: {local_path}")

                    # Загружаем в S3
                    try:
                        year = self._extract_year_from_period(period)
                        if year is None:
                            logger.warning(f"  Пропуск: не удалось определить год для {period}")
                            continue
                        s3_path = self.s3_client.upload_report(ticker, year, period, local_path)
                        if s3_path is None:
                            logger.error(f"  Ошибка загрузки в S3")
                            continue
                        logger.info(f"  Загружено в S3: {s3_path}")

                        # Сохраняем в БД
                        if self._save_report_to_db_with_year(ticker, year, period, s3_path):
                            saved_count += 1
                    except Exception as e:
                        logger.error(f"  Ошибка загрузки в S3: {e}")

            if errors:
                logger.warning("ОШИБКИ:")
                for report in errors:
                    logger.warning(f"  {report.get('name', 'Unknown')}")

            logger.info("=" * 60)
            logger.info(f"Сохранено в БД: {saved_count} отчётов")
            logger.info("=" * 60)

            self.parser.unzip_downloaded_files()

    