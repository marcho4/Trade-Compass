import os
import re
import boto3
from typing import Optional
from botocore.exceptions import ClientError


class S3ReportsStorage:
    '''
    Класс для работы с S3 хранением отчетов. Интерфейс - upload_report, get_s3_report_link
    '''
    def __init__(self):
        key_id = os.environ.get("YANDEX_CLOUD_S3_ACCESS_KEY_ID")
        key_secret = os.environ.get("YANDEX_CLOUD_S3_SECRET_ACCESS_KEY")
        bucket_name = os.environ.get("BUCKET_NAME")

        self.validate_env(key_id, "YANDEX_CLOUD_S3_ACCESS_KEY_ID")
        self.validate_env(key_secret, "YANDEX_CLOUD_S3_SECRET_ACCESS_KEY")
        self.validate_env(bucket_name, "BUCKET_NAME")

        self.bucket_name: str = bucket_name  # type: ignore

        self.client = boto3.client(
            's3',
            endpoint_url='https://storage.yandexcloud.net',
            aws_access_key_id=key_id,
            aws_secret_access_key=key_secret,
            region_name='ru-central1'
        )

    def upload_report(self, ticker: str, year: int, period: str, file_path: str) -> Optional[str]:
        try:
            ticker_normalized = self.__normalize_string(ticker)
            period_normalized = self.__normalize_string(period)

            file_extension = os.path.splitext(file_path)[1]

            object_key = f"reports/{ticker_normalized}/{year}/{period_normalized}/report{file_extension}"

            with open(file_path, 'rb') as file:
                file_content = file.read()

            self.__upload_file(self.bucket_name, object_key, file_content)
            url = f"https://storage.yandexcloud.net/{self.bucket_name}/{object_key}"
            return url

        except FileNotFoundError:
            print(f"✗ Ошибка: файл не найден - {file_path}")
            return None
        except Exception as e:
            print(f"✗ Ошибка при загрузке файла: {e}")
            return None

    def get_s3_report_link(self, ticker: str, year: int, period: str, extension: str = ".zip") -> Optional[str]:
        try:
            ticker_normalized = self.__normalize_string(ticker)
            period_normalized = self.__normalize_string(period)

            object_key = f"reports/{ticker_normalized}/{year}/{period_normalized}/report{extension}"

            try:
                self.client.head_object(Bucket=self.bucket_name, Key=object_key)
                url = f"https://storage.yandexcloud.net/{self.bucket_name}/{object_key}"
                return url
            except ClientError as e:
                if e.response['Error']['Code'] == '404':
                    print(f"✗ Файл не найден в S3: {object_key}")
                    return None
                raise

        except Exception as e:
            print(f"✗ Ошибка при получении ссылки: {e}")
            return None

    def generate_presigned_url(self, ticker: str, year: int, period: str,
                              extension: str = ".zip", expiration: int = 3600) -> Optional[str]:
        try:
            ticker_normalized = self.__normalize_string(ticker)
            period_normalized = self.__normalize_string(period)
            object_key = f"reports/{ticker_normalized}/{year}/{period_normalized}/report{extension}"

            url = self.client.generate_presigned_url(
                'get_object',
                Params={'Bucket': self.bucket_name, 'Key': object_key},
                ExpiresIn=expiration
            )
            return url

        except Exception as e:
            print(f"✗ Ошибка при генерации подписанной ссылки: {e}")
            return None

    def __upload_file(self, bucket_name: str, object_key: str, file: bytes) -> dict:
        response = self.client.put_object(
            Bucket=bucket_name,
            Key=object_key,
            Body=file
        )
        return response

    @staticmethod
    def __normalize_string(value: str) -> str:
        normalized = re.sub(r'[^\w\s-]', '', value)
        normalized = re.sub(r'[-\s]+', '-', normalized)
        return normalized.strip('-').lower()

    @staticmethod
    def validate_env(variable, name: str):
        if not variable:
            print(f"Missing ENV variable: {name}")
            exit(1)
