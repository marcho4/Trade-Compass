import logging
from bs4 import BeautifulSoup

logger = logging.getLogger(__name__)


class ReportMetadataParser:
    """Парсер метаданных отчетов из HTML таблицы."""

    def parse_row(self, tr_element) -> dict | None:
        """Парсинг одной строки таблицы отчетов."""
        try:
            cells = tr_element.find_all("td")
            if len(cells) < 6:
                return None

            metadata = {
                "row_number": cells[0].text.strip() if len(cells) > 0 else "",
                "document_type": cells[1].text.strip() if len(cells) > 1 else "",
                "reporting_period": cells[2].text.strip() if len(cells) > 2 else "",
                "base_date": cells[3].text.strip() if len(cells) > 3 else "",
                "publication_date": cells[4].text.strip() if len(cells) > 4 else "",
            }

            file_cell = cells[5] if len(cells) > 5 else None
            if file_cell:
                file_link = file_cell.find("a", class_="file-link")
                if file_link:
                    metadata["file_url"] = file_link.get("href", "")
                    metadata["file_id"] = file_link.get("data-fileid", "")
                    metadata["file_info"] = file_link.text.strip()

            return metadata

        except Exception as e:
            logger.debug(f"Ошибка парсинга метаданных строки: {e}")
            return None

    def parse_table(self, soup: BeautifulSoup) -> list[dict]:
        """Парсинг всей таблицы отчетов."""
        table = soup.find("table", class_="files-table")
        if not table:
            logger.warning("Таблица с отчетами не найдена")
            return []

        tbody = table.find("tbody")
        if tbody:
            rows = tbody.find_all("tr")
        else:
            rows = table.find_all("tr")[1:]  # Пропускаем заголовок

        logger.info(f"Найдено строк в таблице отчетов: {len(rows)}")

        results = []
        for tr in rows:
            metadata = self.parse_row(tr)
            if metadata and metadata.get("file_url"):
                results.append(metadata)

        return results
