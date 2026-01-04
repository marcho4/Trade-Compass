import logging
import os
import shutil
import zipfile
from pathlib import Path
from infra.config import config

logger = logging.getLogger(__name__)


class FileUnzipper:
    @staticmethod
    def unzip_all(source_dir: str = None, target_dir: str = None):
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

    @staticmethod
    def unzip_and_rename(zip_path: str) -> str | None:
        zip_path = Path(zip_path)
        temp_dir = zip_path.parent / zip_path.stem

        try:
            with zipfile.ZipFile(zip_path, "r") as zf:
                files = [f for f in zf.namelist() if not f.endswith("/")]
                if not files:
                    logger.warning(f"Архив пустой: {zip_path}")
                    return None

                original_name = files[0]
                extension = Path(original_name).suffix

                temp_dir.mkdir(exist_ok=True)
                zf.extract(original_name, temp_dir)

                new_name = zip_path.stem + extension
                temp_file = temp_dir / original_name
                final_path = zip_path.parent / new_name

                temp_file.rename(final_path)

                if temp_dir.exists():
                    shutil.rmtree(temp_dir)

                return str(final_path)
                
        except zipfile.BadZipFile:
            logger.error(f"Некорректный zip-файл: {zip_path}")
            return None
        except Exception as e:
            logger.error(f"Ошибка при распаковке {zip_path}: {e}")
            return None
