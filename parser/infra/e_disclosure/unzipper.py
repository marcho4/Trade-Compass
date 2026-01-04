import logging
import os
import zipfile
from pathlib import Path
from infra.config import config

logger = logging.getLogger(__name__)


class FileUnzipper:
    def unzip_all(self, source_dir: str = None, target_dir: str = None):
        source_dir = source_dir or config.download_dir
        target_dir = target_dir or config.unzip_dir

        os.makedirs(target_dir, exist_ok=True)
        downloads_path = Path(source_dir)

        zip_files = list(downloads_path.glob("*.zip"))

        if not zip_files:
            return []

        extracted = []
        for zip_file in zip_files:
            try:
                with zipfile.ZipFile(zip_file, "r") as zip_ref:
                    zip_ref.extractall(target_dir)
                extracted.append(zip_file.name)
            except Exception as e:
                logger.error(f"Ошибка при распаковке {zip_file.name}: {e}")

        return extracted
