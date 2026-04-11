import pytest
from bs4 import BeautifulSoup

from gateway.e_disclosure.metadata_parser import ReportMetadataParser


ROW_HTML = """
<tr>
  <td>1</td>
  <td>Годовой отчет</td>
  <td>2024</td>
  <td>31.12.2024</td>
  <td>15.03.2025</td>
  <td>
    <a class="file-link" href="/files/report_123.zip" data-fileid="456">
      report_123.zip
    </a>
  </td>
</tr>
"""

TABLE_HTML = f"""
<table class="files-table">
  <tbody>
    {ROW_HTML}
  </tbody>
</table>
"""


@pytest.fixture
def parser():
    return ReportMetadataParser()


def test_parse_row_returns_metadata(parser):
    row = BeautifulSoup(ROW_HTML, "html.parser").find("tr")

    result = parser.parse_row(row)

    assert result is not None
    assert result.document_type == "Годовой отчет"
    assert result.reporting_period == "2024"
    assert result.file_url == "/files/report_123.zip"
    assert result.file_id == "456"


def test_parse_table_returns_list(parser):
    soup = BeautifulSoup(TABLE_HTML, "html.parser")

    results = parser.parse_table(soup)

    assert len(results) == 1
    assert results[0].file_url == "/files/report_123.zip"
