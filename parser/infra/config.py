import os
from dataclasses import dataclass
from dotenv import load_dotenv

load_dotenv()


@dataclass
class ParserConfig:
    # E-Disclosure
    base_url: str = "https://www.e-disclosure.ru"
    user_agent: str = (
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "
        "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
    )

    # Таймауты (секунды)
    timeout_page_load: int = 5
    timeout_element_wait: int = 10
    timeout_between_requests: float = 2
    timeout_between_files: float = 3
    timeout_download: int = 60

    # Лимиты
    max_reports_per_company: int = 10

    # Пути
    download_dir: str = "./downloads"
    unzip_dir: str = "./downloads/unzipped"

    # Безопасность
    admin_api_key: str = ""

    @classmethod
    def from_env(cls) -> "ParserConfig":
        return cls(
            base_url=os.getenv("E_DISCLOSURE_URL", cls.base_url),
            max_reports_per_company=int(
                os.getenv("MAX_REPORTS", str(cls.max_reports_per_company))
            ),
            download_dir=os.getenv("DOWNLOADS_DIR", cls.download_dir),
            unzip_dir=os.getenv("UNZIP_DIR", cls.unzip_dir),
            admin_api_key=os.getenv("ADMIN_API_KEY", ""),
        )


config = ParserConfig.from_env()
