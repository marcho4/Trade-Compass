import json
import logging
import threading

from confluent_kafka import Consumer, KafkaError

from application.reports_processor import ReportProcessor
from application.vectorization_service import VectorizationService
from infra.e_disclosure import EDisclosureClient
from infra.database import get_db_session
from infra.db_repo import ReportsRepository
from infra.s3_storage import S3ReportsStorage
from infra.config import config

logger = logging.getLogger(__name__)


class TickerParseConsumer:
    def __init__(self):
        self._thread: threading.Thread | None = None
        self._running = False

    def start(self):
        self._running = True
        self._thread = threading.Thread(target=self._consume_loop, daemon=True)
        self._thread.start()
        logger.info("Kafka consumer started for topic: %s", config.kafka_parse_ticker_topic)

    def stop(self):
        self._running = False
        if self._thread:
            self._thread.join(timeout=10)
            logger.info("Kafka consumer stopped")

    def _consume_loop(self):
        consumer = Consumer({
            "bootstrap.servers": config.kafka_bootstrap_servers,
            "group.id": config.kafka_consumer_group,
            "auto.offset.reset": "earliest",
        })
        consumer.subscribe([config.kafka_parse_ticker_topic])

        try:
            while self._running:
                msg = consumer.poll(timeout=1.0)
                if msg is None:
                    continue
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        continue
                    logger.error("Kafka consumer error: %s", msg.error())
                    continue

                try:
                    data = json.loads(msg.value().decode("utf-8"))
                    self._handle_message(data)
                except Exception as e:
                    logger.error("Failed to process Kafka message: %s", e)
        finally:
            consumer.close()

    def _handle_message(self, data: dict):
        ticker = data.get("ticker")
        name = data.get("name")
        if not ticker:
            logger.warning("Received message without ticker: %s", data)
            return

        if not name:
            logger.warning("Received message without name for ticker %s: %s", ticker, data)
            return

        logger.info("Processing parsing request for ticker=%s, name=%s", ticker, name)

        try:
            with get_db_session() as db:
                repo = ReportsRepository(db)
                s3_client = S3ReportsStorage()
                vectorization_service = VectorizationService()
                processor = ReportProcessor(s3_client, repo, vectorization_service)
                with EDisclosureClient() as client:
                    result = processor.process_company_by_query(
                        client,
                        query=name,
                        ticker=ticker,
                        skip_indexing=True,
                    )
                logger.info("Parsing result for %s: %s", ticker, result)
        except Exception as e:
            logger.error("Failed to process ticker %s: %s", ticker, e)
