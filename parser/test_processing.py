import logging
from infra.database import get_db_session, init_db
from infra.db_repo import ReportsRepository
from infra.s3_storage import S3ReportsStorage
from parser.application.reports_processor import ReportProcessor

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


def main():
    test_company = ["2309085638"]

    logger.info("Инициализация БД...")
    init_db()

    with get_db_session() as db:
        repo = ReportsRepository(db)
        s3_client = S3ReportsStorage()
        logger.info("Запуск тестового парсинга...")
        processor = ReportProcessor(
            s3_client, repo
        )

    results = processor.process_companies(test_company)
        

    logger.info(f"Результаты: {results}")


if __name__ == "__main__":
    main()

