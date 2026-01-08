from infra.pdf_processor import PDFTextExtractor, TextChunker
import pytest

@pytest.fixture()
def local_path():
    return "tests/demo_pdf_mgnt.pdf"

def test_read_pdf(local_path):
    extractor = PDFTextExtractor()
    text = extractor.extract_text_from_pdf(local_path)
    print(text[:100])
    assert len(text) != 0

def test_chunking_pdf():
    chunker = TextChunker()
    chunks = chunker.chunk_text("some very\n\nvery big text")
    assert len(chunks) != 0
