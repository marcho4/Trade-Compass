from contextlib import asynccontextmanager
from fastapi import FastAPI

from parser.adapter.http.handlers import router
from adapter.kafka.consumer import TickerParseConsumer
from gateway.kafka.producer import AiAnalyzeGateway
from infra.database import init_db


@asynccontextmanager
async def lifespan(app: FastAPI):
    init_db()
    gateway = AiAnalyzeGateway()
    consumer = TickerParseConsumer(gateway)
    consumer.start()
    yield
    consumer.stop()


api = FastAPI(
    title="Reports Storage API",
    lifespan=lifespan
)

api.include_router(router)
