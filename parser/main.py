from contextlib import asynccontextmanager
from fastapi import FastAPI

from application.handlers import router
from infra.database import init_db


@asynccontextmanager
async def lifespan(app: FastAPI):
    init_db()
    yield


api = FastAPI(
    title="Reports Storage API",
    lifespan=lifespan
)

api.include_router(router)
