import logging
import time
from typing import List, Optional

from google import genai
from google.genai.types import ContentEmbedding, EmbedContentConfig

from infra.config import config

logger = logging.getLogger(__name__)


class GeminiEmbeddingService:
    def __init__(
        self,
        api_key: Optional[str] = None,
        model: Optional[str] = None,
        batch_size: Optional[int] = None,
    ):
        self.api_key = api_key or config.gemini_api_key
        self.model = model or config.embedding_model
        self.batch_size = batch_size or config.embedding_batch_size

        if not self.api_key:
            raise ValueError(
                "Gemini API key not provided. Set GEMINI_API_KEY in environment."
            )

        self.client = genai.Client(
            api_key=self.api_key,
        )

    def generate_embeddings(
        self, texts: List[str], retry_delay: float = 1.0, max_retries: int = 3
    ) -> List[ContentEmbedding]:
        if not texts:
            logger.warning("Empty texts list provided")
            return []

        all_embeddings: List[ContentEmbedding] = []
        total_batches = (len(texts) + self.batch_size - 1) // self.batch_size

        logger.info(
            f"Generating embeddings for {len(texts)} texts in {total_batches} batch(es)"
        )

        for i in range(0, len(texts), self.batch_size):
            batch = texts[i : i + self.batch_size]
            batch_num = i // self.batch_size + 1

            logger.debug(
                f"Processing batch {batch_num}/{total_batches} ({len(batch)} texts)"
            )

            batch_embeddings = self._generate_batch_with_retry(
                batch, retry_delay, max_retries
            )

            all_embeddings.extend(batch_embeddings)

            if i + self.batch_size < len(texts):
                time.sleep(0.5)

        logger.info(f"Successfully generated {len(all_embeddings)} embeddings")
        return all_embeddings

    def _generate_batch_with_retry(
        self,
        texts: List[str],
        retry_delay: float,
        max_retries: int,
    ) -> List[ContentEmbedding]:
        for attempt in range(max_retries):
            try:
                result = self.client.models.embed_content(
                    model=self.model,
                    contents=[text for text in texts],
                    config=EmbedContentConfig(
                        output_dimensionality=3072,
                        task_type='retrieval_document'
                    )
                )
                if result.embeddings is None:
                    raise ValueError("No embeddings returned")

                return list(result.embeddings)

            except Exception as e:
                logger.warning(
                    f"Attempt {attempt + 1}/{max_retries} failed: {str(e)[:100]}"
                )

                if attempt < max_retries - 1:
                    sleep_time = retry_delay * (2**attempt)
                    logger.info(f"Retrying in {sleep_time:.1f} seconds...")
                    time.sleep(sleep_time)
                else:
                    logger.error(f"Failed to generate embeddings after {max_retries} attempts")
                    raise

        raise RuntimeError("Unexpected exit from retry loop")

    # Берем 3072 чтобы не нормализовать вывод
    def get_embedding_dimension(self) -> int:
        return 3072
