import logging.config
from contextlib import asynccontextmanager

from fastapi import FastAPI

from adapter.http.handlers import router
from adapter.kafka.consumer import TickerParseConsumer
from gateway.kafka.producer import AiAnalyzeGateway
from infra.database import init_db
from infra.logging_config import get_log_config


logging.config.dictConfig(get_log_config())

logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    init_db()
    logger.info("Database initialized")

    gateway = AiAnalyzeGateway()
    consumer = TickerParseConsumer(gateway)
    consumer.start()

    logger.info("Parser service started")
    yield

    consumer.stop()
    logger.info("Parser service stopped")


api = FastAPI(
    title="Reports Storage API",
    lifespan=lifespan,
)

api.include_router(router)
