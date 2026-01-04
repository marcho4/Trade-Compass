import zipfile
from pathlib import Path

import pytest

from infra.e_disclosure.unzipper import FileUnzipper


class TestUnzipAndRename:
    @pytest.fixture
    def temp_zip(self, tmp_path):
        """Создаёт временный zip-архив с файлом внутри."""
        def _create_zip(zip_name: str, inner_file_name: str, content: bytes = b"test content"):
            zip_path = tmp_path / zip_name
            with zipfile.ZipFile(zip_path, "w") as zf:
                zf.writestr(inner_file_name, content)
            return zip_path
        return _create_zip

    def test_unzip_creates_file_with_zip_name(self, temp_zip):
        """После распаковки файл должен иметь имя архива с оригинальным расширением."""
        zip_path = temp_zip("Company_2024_12M_123.zip", "Отчет_компании.xlsx")

        result = FileUnzipper.unzip_and_rename(str(zip_path))

        assert result is not None
        result_path = Path(result)
        assert result_path.exists()
        assert result_path.name == "Company_2024_12M_123.xlsx"
        assert result_path.parent == zip_path.parent

    def test_unzip_preserves_extension(self, temp_zip):
        """Расширение файла должно сохраняться из оригинального файла в архиве."""
        zip_path = temp_zip("Report.zip", "data.csv")

        result = FileUnzipper.unzip_and_rename(str(zip_path))

        assert Path(result).suffix == ".csv"

    def test_unzip_creates_directory_with_zip_name(self, temp_zip, tmp_path):
        """Временная папка должна быть удалена после распаковки."""
        zip_path = temp_zip("MyArchive.zip", "inner.txt")

        result = FileUnzipper.unzip_and_rename(str(zip_path))

        result_path = Path(result)
        assert result_path.parent == tmp_path

        temp_dir = tmp_path / "MyArchive"
        assert not temp_dir.exists()

    def test_unzip_empty_archive_returns_none(self, tmp_path):
        """Пустой архив должен вернуть None."""
        zip_path = tmp_path / "empty.zip"
        with zipfile.ZipFile(zip_path, "w"):
            pass

        result = FileUnzipper.unzip_and_rename(str(zip_path))

        assert result is None

    def test_unzip_bad_zip_returns_none(self, tmp_path):
        """Некорректный zip должен вернуть None."""
        bad_zip = tmp_path / "bad.zip"
        bad_zip.write_text("not a zip file")

        result = FileUnzipper.unzip_and_rename(str(bad_zip))

        assert result is None
