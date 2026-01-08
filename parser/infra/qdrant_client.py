import logging
from typing import List, Dict, Any, Optional

from qdrant_client import QdrantClient
from qdrant_client.models import (
    Distance,
    VectorParams,
    PointStruct,
    Filter,
    FieldCondition,
    MatchValue,
)

from infra.config import config

logger = logging.getLogger(__name__)


class QdrantVectorStore:
    def __init__(self, host: Optional[str] = None, port: Optional[int] = None):
        self.host = host or config.qdrant_host
        self.port = port or config.qdrant_port
        self.collection_name = config.qdrant_collection_name

        self.client = QdrantClient(host=self.host, port=self.port)
        logger.info(f"Qdrant client initialized: {self.host}:{self.port}")

    def ensure_collection_exists(self, vector_size: int = 3072) -> None:
        try:
            collections = self.client.get_collections().collections
            collection_names = [col.name for col in collections]

            if self.collection_name not in collection_names:
                self.client.create_collection(
                    collection_name=self.collection_name,
                    vectors_config=VectorParams(
                        size=vector_size,
                        distance=Distance.COSINE,
                    ),
                )
                logger.info(
                    f"Collection '{self.collection_name}' created with vector size {vector_size}"
                )
            else:
                logger.info(f"Collection '{self.collection_name}' already exists")
        except Exception as e:
            logger.error(f"Error ensuring collection exists: {e}")
            raise

    def upsert_vectors(
        self,
        vectors: List[List[float]],
        report_id: int,
        ticker: str,
        year: int,
        period: str,
        chunk_ids: List[int],
        texts: Optional[List[str]] = None,
    ) -> None:
        if len(vectors) != len(chunk_ids):
            raise ValueError("Number of vectors must match number of chunk_ids")

        if texts and len(texts) != len(vectors):
            raise ValueError("If provided, number of texts must match number of vectors")

        points = []
        for idx, (vector, chunk_id) in enumerate(zip(vectors, chunk_ids)):
            payload = {
                "report_id": report_id,
                "ticker": ticker,
                "year": year,
                "period": period,
                "chunk_id": chunk_id,
            }

            if texts:
                payload["text"] = texts[idx]

            point_id = report_id * 100000 + chunk_id

            points.append(
                PointStruct(
                    id=point_id,
                    vector=vector,
                    payload=payload,
                )
            )

        try:
            self.client.upsert(
                collection_name=self.collection_name,
                points=points,
            )
            logger.info(
                f"Upserted {len(points)} vectors for report_id={report_id}, ticker={ticker}"
            )
        except Exception as e:
            logger.error(f"Error upserting vectors: {e}")
            raise

    def delete_report_vectors(self, report_id: int) -> None:
        try:
            self.client.delete(
                collection_name=self.collection_name,
                points_selector=Filter(
                    must=[
                        FieldCondition(
                            key="report_id", match=MatchValue(value=report_id)
                        )
                    ]
                ),
            )
            logger.info(f"Deleted vectors for report_id={report_id}")
        except Exception as e:
            logger.error(f"Error deleting vectors: {e}")
            raise
