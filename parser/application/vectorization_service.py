import logging
from typing import Optional, Dict, Any

from infra.pdf_processor import PDFTextExtractor, TextChunker
from infra.gemini_embeddings import GeminiEmbeddingService
from infra.qdrant_client import QdrantVectorStore

logger = logging.getLogger(__name__)


class VectorizationService:
    def __init__(
        self,
        pdf_extractor: Optional[PDFTextExtractor] = None,
        text_chunker: Optional[TextChunker] = None,
        embedding_service: Optional[GeminiEmbeddingService] = None,
        vector_store: Optional[QdrantVectorStore] = None,
    ):
        self.pdf_extractor = pdf_extractor or PDFTextExtractor()
        self.text_chunker = text_chunker or TextChunker()
        self.embedding_service = embedding_service or GeminiEmbeddingService()
        self.vector_store = vector_store or QdrantVectorStore()

        embedding_dim = self.embedding_service.get_embedding_dimension()
        self.vector_store.ensure_collection_exists(vector_size=embedding_dim)

    def vectorize_report(
        self, report_id: int, file_path: str, ticker: str, year: int, period: str
    ) -> Dict[str, Any]:

        try:
            logger.info("Step 1/4: Extracting text from PDF...")
            text = self.pdf_extractor.extract_text_from_pdf(file_path)

            if not text or len(text.strip()) < 100:
                logger.warning(
                    f"Extracted text is too short ({len(text)} chars). "
                    f"PDF might be empty or unreadable."
                )
                return {
                    "success": False,
                    "error": "Extracted text is too short or empty"
                }

            logger.info("Step 2/4: Splitting text into chunks...")
            chunks = self.text_chunker.chunk_text(text)

            if not chunks:
                logger.warning("No chunks created from text")
                return {
                    "success": False,
                    "error": "No chunks created",
                    "text_length": len(text),
                }

            logger.info(f"Step 3/4: Generating embeddings for {len(chunks)} chunks...")
            embeddings = self.embedding_service.generate_embeddings(chunks)

            if len(embeddings) != len(chunks):
                raise ValueError(
                    f"Mismatch: {len(chunks)} chunks but {len(embeddings)} embeddings"
                )

            logger.info("Step 4/4: Saving vectors to Qdrant...")

            try:
                self.vector_store.delete_report_vectors(report_id)
            except Exception as e:
                logger.warning(f"Could not delete old vectors: {e}")

            chunk_ids = list(range(len(chunks)))
            self.vector_store.upsert_vectors(
                vectors=embeddings,
                report_id=report_id,
                ticker=ticker,
                year=year,
                period=period,
                chunk_ids=chunk_ids,
                texts=chunks,
            )

            result = {
                "success": True,
                "report_id": report_id,
                "ticker": ticker,
                "year": year,
                "period": period,
                "file_path": file_path,
                "text_length": len(text),
                "chunks_count": len(chunks),
                "embeddings_count": len(embeddings),
                "vector_dimension": len(embeddings[0]) if embeddings else 0,
            }

            logger.info(
                f"Vectorization completed successfully: "
                f"{len(chunks)} chunks, {len(embeddings)} embeddings"
            )

            return result

        except FileNotFoundError as e:
            logger.error(f"File not found: {e}")
            return {
                "success": False,
                "error": f"File not found: {str(e)}",
                "report_id": report_id,
            }

        except Exception as e:
            logger.error(f"Error during vectorization: {e}", exc_info=True)
            return {
                "success": False,
                "error": str(e),
                "report_id": report_id,
            }
