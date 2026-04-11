import pytest
import requests

ADMIN_API_KEY = "test-admin-key-integration"


@pytest.fixture(scope="session")
def seeded_reports(parser_url):
    reports = [
        ("SBER", 2023, "12", "reports/SBER/2023/12/test.pdf"),
        ("SBER", 2022, "6", "reports/SBER/2022/6/test.pdf"),
        ("GAZP", 2023, "12", "reports/GAZP/2023/12/test.pdf"),
    ]
    for ticker, year, period, s3_path in reports:
        requests.post(
            f"{parser_url}/reports",
            params={"ticker": ticker, "year": year, "period": period, "s3_path": s3_path},
        )
    return reports


@pytest.mark.integration
def test_health_returns_ok(parser_url):
    response = requests.get(f"{parser_url}/health")

    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


@pytest.mark.integration
def test_get_reports_returns_list(parser_url, seeded_reports):
    response = requests.get(f"{parser_url}/reports")

    assert response.status_code == 200
    body = response.json()
    assert isinstance(body["reports"], list)
    assert isinstance(body["total"], int)
    assert body["total"] >= 3


@pytest.mark.integration
def test_get_reports_pagination(parser_url, seeded_reports):
    response = requests.get(f"{parser_url}/reports", params={"skip": 0, "limit": 1})

    assert response.status_code == 200
    body = response.json()
    assert len(body["reports"]) <= 1


@pytest.mark.integration
def test_get_reports_by_ticker_filters_by_ticker(parser_url, seeded_reports):
    response = requests.get(f"{parser_url}/reports/SBER")

    assert response.status_code == 200
    body = response.json()
    assert body["total"] == 2
    for report in body["reports"]:
        assert report["ticker"] == "SBER"


@pytest.mark.integration
def test_get_reports_by_unknown_ticker_returns_empty(parser_url):
    response = requests.get(f"{parser_url}/reports/TICKER_THAT_DOES_NOT_EXIST_XYZ")

    assert response.status_code == 200
    body = response.json()
    assert body["total"] == 0
    assert body["reports"] == []


@pytest.mark.integration
@pytest.mark.parametrize(
    "ticker,year,period,s3_path",
    [
        ("LKOH", 2023, "12", "reports/LKOH/2023/12/test.pdf"),
        ("NVTK", 2022, "3", "reports/NVTK/2022/3/test.pdf"),
    ],
)
def test_create_report_returns_report_id(parser_url, ticker, year, period, s3_path):
    response = requests.post(
        f"{parser_url}/reports",
        params={"ticker": ticker, "year": year, "period": period, "s3_path": s3_path},
    )

    assert response.status_code == 200
    body = response.json()
    assert isinstance(body["report_id"], int)


@pytest.mark.integration
def test_create_duplicate_report_returns_error(parser_url, seeded_reports):
    response = requests.post(
        f"{parser_url}/reports",
        params={"ticker": "SBER", "year": 2023, "period": "12", "s3_path": "reports/SBER/2023/12/test.pdf"},
    )

    assert response.status_code == 200
    body = response.json()
    assert body["error"] == "Report already exists"


@pytest.mark.integration
def test_start_parsing_without_api_key_returns_422(parser_url):
    response = requests.post(f"{parser_url}/start_parsing")

    assert response.status_code == 422


@pytest.mark.integration
def test_start_parsing_with_wrong_key_returns_403(parser_url):
    response = requests.post(f"{parser_url}/start_parsing", headers={"X-API-Key": "wrong-key"})

    assert response.status_code == 403


@pytest.mark.integration
def test_start_parsing_with_correct_key_returns_200(parser_url):
    response = requests.post(
        f"{parser_url}/start_parsing",
        params={"skip_indexing": "true"},
        headers={"X-API-Key": ADMIN_API_KEY},
    )

    assert response.status_code == 200
    body = response.json()
    assert body["skip_indexing"] is True


@pytest.mark.integration
def test_upload_report_with_invalid_period_returns_400(parser_url):
    response = requests.post(
        f"{parser_url}/reports/upload",
        data={"ticker": "SBER", "year": "2023", "period": "7"},
        files={"file": ("test.pdf", b"%PDF-1.4 test", "application/pdf")},
    )

    assert response.status_code == 400
