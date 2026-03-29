from dataclasses import dataclass

@dataclass
class ReportMetadata:
    row_number: str
    document_type: str
    reporting_period: str
    base_date: str
    publication_date: str
    file_url: str
    file_id: str
    file_info: str
