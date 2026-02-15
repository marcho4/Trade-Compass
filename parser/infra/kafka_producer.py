import json
import logging

from confluent_kafka import Producer

from infra.config import config

logger = logging.getLogger(__name__)

MONTHS_TO_PERIOD = {
    "3": "Q1",
    "6": "Q2",
    "9": "Q3",
    "12": "YEAR",
}


class AnalyzeTaskProducer:
    def __init__(self):
        self._producer = Producer({
            "bootstrap.servers": config.kafka_bootstrap_servers,
        })

    def send_analyze_task(self, ticker: str, year: int, period: str, report_url: str) -> None:
        ai_period = MONTHS_TO_PERIOD.get(period)
        if not ai_period:
            logger.warning("Unknown period '%s' for ticker %s, skipping analyze task", period, ticker)
            return

        message = {
            "ticker": ticker,
            "year": year,
            "period": ai_period,
            "report_url": report_url,
            "type": "analyze",
        }

        try:
            self._producer.produce(
                topic=config.kafka_ai_analyze_topic,
                value=json.dumps(message).encode("utf-8"),
                callback=self._delivery_callback,
            )
            self._producer.flush(timeout=5)
            logger.info("Sent analyze task for %s %d %s", ticker, year, ai_period)
        except Exception as e:
            logger.error("Failed to send analyze task for %s: %s", ticker, e)

    @staticmethod
    def _delivery_callback(err, msg):
        if err:
            logger.error("Kafka delivery failed: %s", err)
        else:
            logger.debug("Message delivered to %s [%d]", msg.topic(), msg.partition())
