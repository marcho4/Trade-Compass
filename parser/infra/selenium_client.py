import logging
import os
import requests
from bs4 import BeautifulSoup
import time
from urllib.parse import urljoin
import re
import logging
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from webdriver_manager.chrome import ChromeDriverManager

logger = logging.getLogger(__name__)

class EDisclosureParser:
    def __init__(self):
        self.base_url = "https://www.e-disclosure.ru"
        self.driver = None

        chrome_options = Options()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--disable-dev-shm-usage')
        chrome_options.add_argument('--disable-blink-features=AutomationControlled')
        chrome_options.add_argument('--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36')
        chrome_options.add_experimental_option("excludeSwitches", ["enable-automation"])
        chrome_options.add_experimental_option('useAutomationExtension', False)

        try:
            service = Service(ChromeDriverManager().install())
            self.driver = webdriver.Chrome(service=service, options=chrome_options)
            logger.info("Selenium WebDriver инициализирован")
        except Exception as e:
            logger.error(f"Не удалось инициализировать Selenium: {e}")
            return

    def __del__(self):
        if self.driver:
            self.driver.quit()

    def _parse_report_metadata(self, tr_element):
        try:
            cells = tr_element.find_all('td')
            if len(cells) < 6:
                return None

            metadata = {
                'row_number': cells[0].text.strip() if len(cells) > 0 else '',
                'document_type': cells[1].text.strip() if len(cells) > 1 else '',
                'reporting_period': cells[2].text.strip() if len(cells) > 2 else '',
                'base_date': cells[3].text.strip() if len(cells) > 3 else '',
                'publication_date': cells[4].text.strip() if len(cells) > 4 else '',
            }

            file_cell = cells[5] if len(cells) > 5 else None
            if file_cell:
                file_link = file_cell.find('a', class_='file-link')
                if file_link:
                    metadata['file_url'] = file_link.get('href', '')
                    metadata['file_id'] = file_link.get('data-fileid', '')
                    metadata['file_info'] = file_link.text.strip()

            return metadata

        except Exception as e:
            logger.debug(f"Ошибка парсинга метаданных строки: {e}")
            return None

    def search_company(self, query):
        if self.driver:
            return self._search_company_selenium(query)
    
    def _search_company_selenium(self, query):
        search_url = f"{self.base_url}/poisk-po-kompaniyam"

        try:
            logger.info("=== Поиск компании через Selenium ===")
            logger.info(f"URL: {search_url}, Query: {query}")

            self.driver.get(search_url)
            time.sleep(3)
            try:
                search_input = WebDriverWait(self.driver, 10).until(
                    EC.presence_of_element_located((By.ID, "textfield"))
                )

                search_input.clear()
                search_input.send_keys(query)

                try:
                    search_button = WebDriverWait(self.driver, 5).until(
                        EC.element_to_be_clickable((By.ID, "sendButton"))
                    )
                    search_button.click()
                except Exception as e:
                    logger.error(f"Ошибка при нажатии кнопки 'Искать': {e}", exc_info=True)
                time.sleep(2)

                companies = []
                try:
                    company_links = WebDriverWait(self.driver, 15).until(
                        EC.presence_of_all_elements_located((By.XPATH, "//a[contains(@href, '/portal/company.aspx?id=')]"))
                    )

                    for link in company_links:
                        try:
                            href = link.get_attribute('href')
                            name = link.text.strip()
                            match = re.search(r'id=(\d+)', href)
                            if match and name:
                                company_id = match.group(1)
                                companies.append({
                                    'name': name,
                                    'id': company_id,
                                    'url': href,
                                    'element': link
                                })
                        except Exception as e:
                            logger.debug(f"Ошибка обработки ссылки: {e}")
                            continue

                except Exception as e:
                    logger.debug(f"Ошибка парсинга ссылки: {e}")

                logger.info(f"=== Найдено компаний: {len(companies)} ===")
                for i, comp in enumerate(companies, 1):
                    logger.info(f"{i}. {comp['name']} (ID: {comp['id']})")

                return companies

            except Exception as e:
                logger.error(f"Ошибка при работе с элементами страницы: {e}", exc_info=True)
                return []

        except Exception as e:
            logger.error(f"Ошибка поиска: {e}", exc_info=True)
            return []
    
    def get_issuer_reports_by_click(self, company_data, download_dir='./parser/downloads'):
        """Получить отчетность эмитента кликнув на ссылку компании"""
        import os

        # Создаем директорию для загрузок
        os.makedirs(download_dir, exist_ok=True)

        # Используем текущий драйвер с уже открытой страницей поиска
        if not self.driver:
            logger.error("WebDriver не инициализирован")
            return []

        company_id = company_data['id']
        company_name = company_data['name']
        company_element = company_data.get('element')

        try:
            logger.info("=== Переход на страницу компании ===")
            logger.info(f"Компания: {company_name} (ID: {company_id})")

            # Кликаем на ссылку компании
            if company_element:
                try:
                    company_element.click()
                    logger.debug("Клик по ссылке компании выполнен")
                except Exception as e:
                    logger.warning(f"Не удалось кликнуть по элементу: {e}")
                    # Fallback - переходим по URL
                    self.driver.get(company_data['url'])
                    
            else:
                # Если элемента нет, переходим по URL
                self.driver.get(company_data['url'])

            time.sleep(5)  # Ждем загрузки страницы компании
            logger.info("Страница компании загружена")

            reports_url = f"{self.base_url}/portal/files.aspx?id={company_id}&type=5"
            logger.info("=== Переход к отчетности эмитента ===")
            logger.info(f"URL: {reports_url}")
            
            try:
                self.driver.get(reports_url)
                time.sleep(4)
                logger.info("Страница с отчетностью эмитента загружена")

                soup = BeautifulSoup(self.driver.page_source, 'html.parser')
                reports = []

                # Ищем таблицу с отчетами
                table = soup.find('table', class_='files-table')
                if not table:
                    logger.warning("Таблица с отчетами не найдена")
                    return reports

                # Ищем все строки таблицы (пропускаем заголовок)
                table_rows = table.find('tbody').find_all('tr') if table.find('tbody') else table.find_all('tr')[1:]
                logger.info(f"Найдено строк в таблице отчетов: {len(table_rows)}")

                # Обрабатываем каждую строку таблицы
                for i, tr in enumerate(table_rows, 1):
                    if i > 10:  # Ограничение на количество скачиваний
                        break

                    try:
                        # Парсим метаданные из строки таблицы
                        metadata = self._parse_report_metadata(tr)
                        if not metadata or not metadata.get('file_url'):
                            logger.debug(f"Строка {i}: метаданные не извлечены, пропускаем")
                            continue

                        # Формируем полный URL
                        file_url = metadata['file_url']
                        if not file_url.startswith('http'):
                            file_url = urljoin(self.base_url, file_url)

                        # Извлекаем данные из метаданных
                        period = metadata.get('reporting_period', 'unknown')
                        document_type = metadata.get('document_type', '')
                        file_id = metadata.get('file_id', '')
                        file_info = metadata.get('file_info', '')

                        # Формируем имя файла
                        file_name = f"{company_name}_{period}_{file_id}.zip"
                        file_name = re.sub(r'[^\w\s.-]', '_', file_name)
                        file_name = re.sub(r'_+', '_', file_name)

                        logger.info(f"\n{i}. Отчетный период: {period}")
                        logger.info(f"   Тип документа: {document_type}")
                        logger.info(f"   Дата публикации: {metadata.get('publication_date', '')}")
                        logger.info(f"   Файл: {file_info}")
                        logger.debug(f"   URL: {file_url}")

                        # Скачиваем файл через requests
                        try:
                            file_path = os.path.join(download_dir, file_name)
                            logger.info(f"   Скачивание в: {file_path}")

                            # Небольшая задержка перед запросом
                            time.sleep(2)

                            response = requests.get(file_url, timeout=60)
                            response.raise_for_status()

                            with open(file_path, 'wb') as f:
                                f.write(response.content)

                            file_size = len(response.content)
                            logger.info(f"   Файл скачан: {file_size} bytes")

                            reports.append({
                                'name': file_name,
                                'period': period,
                                'document_type': document_type,
                                'publication_date': metadata.get('publication_date', ''),
                                'base_date': metadata.get('base_date', ''),
                                'url': file_url,
                                'status': 'downloaded',
                                'path': file_path,
                                'size': file_size,
                                'metadata': metadata
                            })

                            # Задержка после успешного скачивания
                            if i < len(table_rows):
                                logger.debug("Пауза 3 секунды перед следующим файлом...")
                                time.sleep(3)

                        except Exception as e:
                            logger.error(f"   Ошибка при скачивании: {e}")
                            reports.append({
                                'name': file_name,
                                'period': period,
                                'document_type': document_type,
                                'url': file_url,
                                'status': 'error',
                                'metadata': metadata
                            })
                            time.sleep(1)

                    except Exception as e:
                        logger.error(f"   Ошибка обработки строки {i}: {e}")
                        time.sleep(1)
                
                # Проверяем загруженные файлы
                if os.path.exists(download_dir):
                    downloaded_files = [f for f in os.listdir(download_dir) if not f.startswith('.')]
                    logger.info("="*60)
                    logger.info(f"Всего файлов в {download_dir}: {len(downloaded_files)}")
                    logger.info("="*60)
                    for file in downloaded_files:
                        file_size = os.path.getsize(os.path.join(download_dir, file))
                        logger.info(f"  - {file} ({file_size:,} bytes)")

                return reports

            except Exception as e:
                logger.error(f"Ошибка при работе с меню: {e}", exc_info=True)
                return []

        except Exception as e:
            logger.error(f"Ошибка при получении отчетов: {e}", exc_info=True)
            return []

    def unzip_downloaded_files(self, path="./downloads/unzipped"):
        import zipfile
        from pathlib import Path

        os.makedirs(path, exist_ok=True)
        downloads_dir = Path("./downloads")

        # Найти все ZIP-файлы в папке downloads
        zip_files = list(downloads_dir.glob("*.zip"))

        if not zip_files:
            logger.info("ZIP-файлы не найдены в папке ./downloads")
            return

        logger.info(f"Найдено {len(zip_files)} ZIP-файлов для распаковки")

        # Распаковать каждый ZIP-файл
        for zip_file in zip_files:
            try:
                logger.info(f"Распаковка {zip_file.name}...")
                with zipfile.ZipFile(zip_file, 'r') as zip_ref:
                    zip_ref.extractall(path)
                logger.info(f"Файл {zip_file.name} успешно распакован")
            except Exception as e:
                logger.error(f"Ошибка при распаковке {zip_file.name}: {e}")

        logger.info(f"Распаковка завершена. Файлы находятся в {path}")
