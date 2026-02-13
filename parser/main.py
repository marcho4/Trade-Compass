from contextlib import asynccontextmanager
from fastapi import FastAPI

from application.handlers import router
from infra.database import init_db
from infra.kafka_consumer import TickerParseConsumer


@asynccontextmanager
async def lifespan(app: FastAPI):
    init_db()
    consumer = TickerParseConsumer()
    consumer.start()
    yield
    consumer.stop()


api = FastAPI(
    title="Reports Storage API",
    lifespan=lifespan
)

api.include_router(router)
