import json
import logging

from confluent_kafka import Producer

from infra.config import config
from parser.usecase.interfaces import AnalyzeTaskGateway

logger = logging.getLogger(__name__)

MONTHS_TO_PERIOD = {
    "3": "Q1",
    "6": "Q2",
    "9": "Q3",
    "12": "YEAR",
}


class AiAnalyzeGateway(AnalyzeTaskGateway):
    def __init__(self):
        self._producer = Producer({
            "bootstrap.servers": config.kafka_bootstrap_servers,
            "transactional.id": "parser-txn-id-553",
        })
        self._producer.init_transactions(timeout=10)

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

        self._producer.produce(
            topic=config.kafka_ai_analyze_topic,
            value=json.dumps(message).encode("utf-8"),
            callback=self._delivery_callback,
        )

    def begin_transaction(self) -> None:
        self._producer.begin_transaction()

    def commit_transaction(self) -> None:
        self._producer.commit_transaction()

    def abort_transaction(self) -> None:
        self._producer.abort_transaction()

    def send_offsets_to_transaction(self, positions, group_metadata) -> None:
        self._producer.send_offsets_to_transaction(positions, group_metadata)

    @staticmethod
    def _delivery_callback(err, msg):
        if err:
            logger.error("Kafka delivery failed: %s", err)
        else:
            logger.debug("Message delivered to %s [%d]", msg.topic(), msg.partition())
