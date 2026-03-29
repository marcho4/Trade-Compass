import json
import logging
import threading

from confluent_kafka import Consumer, KafkaError

from companies import get_inn_by_ticker
from infra.config import config
from infra.database import get_db_session
from gateway.e_disclosure import EDisclosureClient
from repository.qdrant.vectorizator import QdrantVectorizationService
from usecase.interfaces import AnalyzeTaskGateway
from repository.postgres.report import PostgresReportsRepository
from gateway.s3.storage import S3ReportsStorage
from usecase.reports import ReportProcessor

logger = logging.getLogger(__name__)


class TickerParseConsumer:
    def __init__(self, gateway: AnalyzeTaskGateway):
        self._gateway = gateway
        self._thread: threading.Thread | None = None
        self._running = False

        self._s3_client = S3ReportsStorage()
        self._vector_store = QdrantVectorizationService()

        self._consumer = Consumer({
            "bootstrap.servers": config.kafka_bootstrap_servers,
            "group.id": config.kafka_consumer_group,
            "auto.offset.reset": "earliest",
            "enable.auto.commit": False,
        })
        self._consumer.subscribe([config.kafka_parse_ticker_topic])

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
        try:
            while self._running:
                msg = self._consumer.poll(timeout=1.0)
                if msg is None:
                    continue

                error = msg.error()
                if error is not None:
                    if error.code() == KafkaError._PARTITION_EOF:
                        continue
                    logger.error("Kafka consumer error: %s", error)
                    continue

                raw = msg.value()
                if raw is None:
                    continue

                try:
                    data = json.loads(raw.decode("utf-8"))
                    self._handle_message(data)
                except Exception as e:
                    logger.error("Failed to process Kafka message: %s", e)
        finally:
            self._consumer.close()

    def _handle_message(self, data: dict):
        ticker = data.get("ticker")
        name = data.get("name")
        task_id = data.get("id")

        if not task_id:
            logger.warning("Received message without id: %s", data)
            return

        if not ticker:
            logger.warning("Received message without ticker: %s", data)
            return

        if not name:
            logger.warning("Received message without name for ticker %s: %s", ticker, data)
            return

        logger.info("Processing parsing request for ticker=%s, name=%s", ticker, name)

        try:
            self._gateway.begin_transaction()

            with get_db_session() as db:
                repo = PostgresReportsRepository(db)
                with EDisclosureClient() as client:
                    inn = get_inn_by_ticker(ticker)
                    query = inn if inn else name

                    processor = ReportProcessor(client, self._s3_client, repo, self._vector_store)
                    result = processor.process_company_by_query(
                        query=query,
                        ticker=ticker,
                        skip_indexing=True,
                    )

            for meta in result.reports_metadata:
                self._gateway.send_task(ticker, meta.year, meta.period, meta.s3_path, task_id, "extract")
                self._gateway.send_task(ticker, meta.year, meta.period, meta.s3_path, task_id, "expect-raw-data")

            self._gateway.send_offsets_to_transaction(
                self._consumer.position(self._consumer.assignment()),
                self._consumer.consumer_group_metadata(),
            )
            self._gateway.commit_transaction()

            logger.info("Parsing result for %s: %s", ticker, result)
        except Exception as e:
            self._gateway.abort_transaction()
            logger.error("Failed to process ticker %s: %s", ticker, e)
