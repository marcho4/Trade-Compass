import json
import logging
import threading

from confluent_kafka import Consumer, KafkaError, Producer

from application.reports_processor import ReportProcessor
from application.vectorization_service import VectorizationService
from infra.e_disclosure import EDisclosureClient
from infra.database import get_db_session
from infra.db_repo import ReportsRepository
from infra.s3_storage import S3ReportsStorage
from infra.config import config
from companies import get_inn_by_ticker

logger = logging.getLogger(__name__)

MONTHS_TO_PERIOD = {
    "3": "Q1",
    "6": "Q2",
    "9": "Q3",
    "12": "YEAR",
}

class TickerParseConsumer:
    def __init__(self):
        self._thread: threading.Thread | None = None
        self._running = False

        self._producer = Producer({
            "bootstrap.servers": config.kafka_bootstrap_servers,
            "transactional.id": "parser-txn-id-553",
        })
        self._producer.init_transactions(timeout=10) 

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
            self._producer.begin_transaction()

            with get_db_session() as db:
                repo = ReportsRepository(db)
                s3_client = S3ReportsStorage()
                vectorization_service = VectorizationService()
                processor = ReportProcessor(s3_client, repo, vectorization_service)
                with EDisclosureClient() as client:
                    inn = get_inn_by_ticker(ticker)
                    query = inn if inn else name

                    result = processor.process_company_by_query(
                        client,
                        query=query,
                        ticker=ticker,
                        skip_indexing=True,
                    )

                for meta in result.reports_metadata:
                    self.send_task(ticker, meta.year, meta.period, meta.s3_path, task_id, "extract")
                    self.send_task(ticker, meta.year, meta.period, meta.s3_path, task_id, "expect")  

                self._producer.send_offsets_to_transaction(
                    self._consumer.position(self._consumer.assignment()),                                                                                      
                    self._consumer.consumer_group_metadata(),
                )                                                                                                                                  

                self._producer.commit_transaction() 

                logger.info("Parsing result for %s: %s", ticker, result)
        except Exception as e:
            self._producer.abort_transaction()
            logger.error("Failed to process ticker %s: %s", ticker, e)

    @staticmethod
    def _delivery_callback(err, msg):
        if err:
            logger.error("Kafka delivery failed: %s", err)
        else:
            logger.debug("Message delivered to %s [%d]", msg.topic(), msg.partition())


    def send_task(self, ticker: str, year: int, period: str, report_url: str, task_id: str, task_type: str) -> None:
        ai_period = MONTHS_TO_PERIOD.get(period)
        if not ai_period:
            logger.warning("unknown period '%s' for ticker %s, skipping analyze task", period, ticker)
            return

        message = {
            "ticker": ticker,
            "year": year,
            "period": ai_period,
            "report_url": report_url,
            "type": task_type,
            "id": task_id,
        }

        try:
            self._producer.produce(
                topic=config.kafka_ai_analyze_topic,
                value=json.dumps(message).encode("utf-8"),
                callback=self._delivery_callback,
            )

        except Exception as e:
            logger.error("Failed to send analyze task for %s: %s", ticker, e)
