import logging
import os
from typing import List, Optional
from pathlib import Path

from pypdf import PdfReader
from langchain_text_splitters import RecursiveCharacterTextSplitter

from infra.config import config

logger = logging.getLogger(__name__)


class PDFTextExtractor:
    def extract_text_from_pdf(self, file_path: str) -> str:
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"PDF file not found: {file_path}")

        if not file_path.lower().endswith(".pdf"):
            raise ValueError(f"File is not a PDF: {file_path}")

        try:
            logger.info(f"Extracting text from PDF: {file_path}")

            reader = PdfReader(file_path)
            total_pages = len(reader.pages)

            logger.debug(f"PDF has {total_pages} pages")

            all_text = []
            for page_num, page in enumerate(reader.pages, start=1):
                text = page.extract_text()

                if text:
                    all_text.append(text)
                    logger.debug(
                        f"Extracted {len(text)} chars from page {page_num}/{total_pages}"
                    )
                else:
                    logger.warning(f"No text found on page {page_num}/{total_pages}")

            full_text = "\n\n".join(all_text)
            logger.info(
                f"Successfully extracted {len(full_text)} characters from {total_pages} pages"
            )

            return full_text

        except Exception as e:
            logger.error(f"Error extracting text from PDF {file_path}: {e}")
            raise

    def extract_text_by_pages(self, file_path: str) -> List[str]:
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"PDF file not found: {file_path}")

        try:
            logger.info(f"Extracting text by pages from PDF: {file_path}")

            reader = PdfReader(file_path)
            pages_text = []

            for _, page in enumerate(reader.pages, start=1):
                text = page.extract_text()
                pages_text.append(text if text else "")

            logger.info(f"Extracted text from {len(pages_text)} pages")
            return pages_text

        except Exception as e:
            logger.error(f"Error extracting text by pages from {file_path}: {e}")
            raise


class TextChunker:
    def __init__(self, chunk_size: Optional[int] = None, chunk_overlap: Optional[int] = None):
        self.chunk_size = chunk_size or config.chunk_size
        self.chunk_overlap = chunk_overlap or config.chunk_overlap

        self.splitter = RecursiveCharacterTextSplitter(
            chunk_size=self.chunk_size,
            chunk_overlap=self.chunk_overlap,
            length_function=len,
            separators=["\n\n", "\n", ". ", " ", ""],
        )

        logger.info(
            f"TextChunker initialized: chunk_size={self.chunk_size}, "
            f"overlap={self.chunk_overlap}"
        )

    def chunk_text(self, text: str) -> List[str]:
        if not text or not text.strip():
            logger.warning("Empty text provided for chunking")
            return []

        try:
            return self.splitter.split_text(text)

        except Exception as e:
            logger.error(f"Error chunking text: {e}")
            raise
