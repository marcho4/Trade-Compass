import logging
from bs4 import BeautifulSoup

from domain.report_metadata import ReportMetadata

logger = logging.getLogger(__name__)


class ReportMetadataParser:
    def parse_row(self, tr_element) -> ReportMetadata | None:
        try:
            cells = tr_element.find_all("td")
            if len(cells) < 6:
                return None

            metadata = ReportMetadata(
                row_number=cells[0].text.strip() if len(cells) > 0 else "",
                document_type=cells[1].text.strip() if len(cells) > 1 else "",
                reporting_period=cells[2].text.strip() if len(cells) > 2 else "",
                base_date=cells[3].text.strip() if len(cells) > 3 else "",
                publication_date=cells[4].text.strip() if len(cells) > 4 else "",
                file_url="",
                file_id="",
                file_info="",
            )

            file_cell = cells[5] if len(cells) > 5 else None
            if file_cell:
                file_link = file_cell.find("a", class_="file-link")
                if file_link:
                    metadata.file_url = file_link.get("href", "")
                    metadata.file_id = file_link.get("data-fileid", "")
                    metadata.file_info = file_link.text.strip()

            return metadata

        except Exception as e:
            logger.debug(f"Ошибка парсинга метаданных строки: {e}")
            return None

    def parse_table(self, soup: BeautifulSoup) -> list[ReportMetadata]:
        table = soup.find("table", class_="files-table")
        if not table:
            return []

        tbody = table.find("tbody")
        if tbody:
            rows = tbody.find_all("tr")
        else:
            rows = table.find_all("tr")[1:]

        results: list[ReportMetadata] = []
        for tr in rows:
            try:
                metadata = self.parse_row(tr)
                if metadata and metadata.file_url:
                    results.append(metadata)
            except Exception as e:
                logger.error(f"Error while parsing table {e}")

        return results
