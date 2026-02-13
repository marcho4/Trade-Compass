import os
from dataclasses import dataclass
from dotenv import load_dotenv

load_dotenv()


@dataclass
class ParserConfig:
    # E-Disclosure
    base_url: str = "https://www.e-disclosure.ru"
    user_agent: str = (
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "
        "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
    )

    # Таймауты (секунды)
    timeout_page_load: int = 5
    timeout_element_wait: int = 10
    timeout_between_requests: float = 2
    timeout_between_files: float = 3
    timeout_download: int = 60

    # Лимиты
    max_reports_per_company: int = 10

    # Пути
    download_dir: str = "./downloads"
    unzip_dir: str = "./downloads/unzipped"

    # Безопасность
    admin_api_key: str = ""

    # Qdrant
    qdrant_host: str = "qdrant"
    qdrant_port: int = 6333
    qdrant_collection_name: str = "reports_embeddings"

    # Gemini Embeddings
    gemini_api_key: str = ""
    embedding_model: str = "gemini-embedding-001"
    embedding_batch_size: int = 100

    # Kafka
    kafka_bootstrap_servers: str = "kafka:9092"
    kafka_parse_ticker_topic: str = "parser.parse_ticker"
    kafka_consumer_group: str = "parser-group"

    # PDF Processing
    chunk_size: int = 1000
    chunk_overlap: int = 200

    @classmethod
    def from_env(cls) -> "ParserConfig":
        return cls(
            base_url=os.getenv("E_DISCLOSURE_URL", cls.base_url),
            max_reports_per_company=int(
                os.getenv("MAX_REPORTS", str(cls.max_reports_per_company))
            ),
            download_dir=os.getenv("DOWNLOADS_DIR", cls.download_dir),
            unzip_dir=os.getenv("UNZIP_DIR", cls.unzip_dir),
            admin_api_key=os.getenv("ADMIN_API_KEY", ""),
            qdrant_host=os.getenv("QDRANT_HOST", cls.qdrant_host),
            qdrant_port=int(os.getenv("QDRANT_PORT", str(cls.qdrant_port))),
            qdrant_collection_name=os.getenv(
                "QDRANT_COLLECTION_NAME", cls.qdrant_collection_name
            ),
            gemini_api_key=os.getenv("GEMINI_API_KEY", ""),
            embedding_model=os.getenv("EMBEDDING_MODEL", cls.embedding_model),
            embedding_batch_size=int(
                os.getenv("EMBEDDING_BATCH_SIZE", str(cls.embedding_batch_size))
            ),
            kafka_bootstrap_servers=os.getenv("KAFKA_BOOTSTRAP_SERVERS", cls.kafka_bootstrap_servers),
            kafka_parse_ticker_topic=os.getenv("KAFKA_PARSE_TICKER_TOPIC", cls.kafka_parse_ticker_topic),
            kafka_consumer_group=os.getenv("KAFKA_CONSUMER_GROUP", cls.kafka_consumer_group),
            chunk_size=int(os.getenv("CHUNK_SIZE", str(cls.chunk_size))),
            chunk_overlap=int(os.getenv("CHUNK_OVERLAP", str(cls.chunk_overlap))),
        )


config = ParserConfig.from_env()
