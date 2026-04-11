import pytest
from unittest.mock import MagicMock

from domain.model.report import DownloadedReport, ReportEntity
from domain.model.processing_result import ReportMetadata
from usecase.reports import ReportProcessor


@pytest.fixture
def mock_client():
    client = MagicMock()
    client.search_company.return_value = [{"id": "123", "name": "Магнит"}]
    client.download_reports.return_value = [
        {
            "status": "downloaded",
            "path": "/tmp/report.pdf",
            "year": 2024,
            "period_months": 12,
            "name": "Годовой отчет",
            "size": 2 * 1024 * 1024,
        }
    ]
    return client


@pytest.fixture
def mock_s3():
    s3 = MagicMock()
    s3.get_report_link.return_value = None
    s3.upload_report.return_value = "reports/MGNT/2024/12/report.pdf"
    return s3


@pytest.fixture
def mock_repo():
    repo = MagicMock()
    repo.create_report.return_value = ReportEntity(
        id=1,
        ticker="MGNT",
        year=2024,
        period="12",
        s3_path="reports/MGNT/2024/12/report.pdf",
    )
    return repo


@pytest.fixture
def mock_vector_store():
    vs = MagicMock()
    vs.vectorize_report.return_value = {"success": True}
    return vs


@pytest.fixture
def processor(mock_client, mock_s3, mock_repo, mock_vector_store):
    return ReportProcessor(mock_client, mock_s3, mock_repo, mock_vector_store)


def test_process_company_returns_result(processor, mock_client):
    result = processor.process_company_by_query(query="7707503127", ticker="MGNT")

    assert result.ticker == "MGNT"
    assert result.saved == 1
    assert len(result.reports_metadata) == 1


def test_process_company_metadata_fields(processor):
    result = processor.process_company_by_query(query="7707503127", ticker="MGNT")

    meta = result.reports_metadata[0]
    assert meta.year == 2024
    assert meta.period == "12"
    assert meta.s3_path == "reports/MGNT/2024/12/report.pdf"


def test_process_company_uploads_to_s3(processor, mock_s3):
    processor.process_company_by_query(query="7707503127", ticker="MGNT")

    mock_s3.upload_report.assert_called_once_with("MGNT", 2024, "12", "/tmp/report.pdf")


def test_process_company_saves_to_db(processor, mock_repo):
    processor.process_company_by_query(query="7707503127", ticker="MGNT")

    mock_repo.create_report.assert_called_once_with(
        "MGNT", 2024, "12", "reports/MGNT/2024/12/report.pdf"
    )


def test_process_company_skips_vectorization_when_skip_indexing(processor, mock_vector_store):
    processor.process_company_by_query(query="7707503127", ticker="MGNT", skip_indexing=True)

    mock_vector_store.vectorize_report.assert_not_called()


def test_process_company_not_found_raises(processor, mock_client):
    mock_client.search_company.return_value = []

    with pytest.raises(ValueError):
        processor.process_company_by_query(query="unknown", ticker="UNKNOWN")
