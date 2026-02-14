import uvicorn

from infra.logging_config import get_log_config

if __name__ == "__main__":
    uvicorn.run(
        "main:api",
        host="0.0.0.0",
        port=8081,
        log_config=get_log_config(),
    )
