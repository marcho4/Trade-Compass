import logging
import os
import requests
from bs4 import BeautifulSoup
import pandas as pd
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

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

class EDisclosureParser:
    def __init__(self, use_selenium=True):
        self.base_url = "https://www.e-disclosure.ru"
        self.use_selenium = use_selenium
        self.driver = None
        
        if use_selenium:
            # Настройка Selenium с headless Chrome
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
        """Закрываем браузер при удалении объекта"""
        if self.driver:
            self.driver.quit()
    
    def search_company(self, query):
        """Поиск компании по названию/тикеру"""
        if self.driver:
            return self._search_company_selenium(query)
    
    def _search_company_selenium(self, query):
        """Поиск через Selenium (обходит JavaScript-защиту)"""
        search_url = f"{self.base_url}/poisk-po-kompaniyam"

        try:
            logger.info("=== Поиск компании через Selenium ===")
            logger.info(f"URL: {search_url}, Query: {query}")

            self.driver.get(search_url)

            # Ждем загрузки страницы
            logger.debug("Ожидание загрузки страницы...")
            time.sleep(3)

            # Находим поле ввода и вводим запрос
            try:
                search_input = WebDriverWait(self.driver, 10).until(
                    EC.presence_of_element_located((By.ID, "textfield"))
                )
                logger.debug("Поле поиска найдено")

                search_input.clear()
                search_input.send_keys(query)
                logger.debug(f"Текст '{query}' введен в поле поиска")

                # Ищем и кликаем кнопку "Искать"
                try:
                    search_button = WebDriverWait(self.driver, 5).until(
                        EC.element_to_be_clickable((By.ID, "sendButton"))
                    )
                    search_button.click()
                    logger.debug("Нажата кнопка 'Искать'")
                except Exception as e:
                    logger.error(f"Ошибка при нажатии кнопки 'Искать': {e}", exc_info=True)

                # Ждем появления результатов (AJAX запрос)
                logger.debug("Ожидание результатов поиска...")
                time.sleep(2)  # Увеличили время ожидания

                # Ищем ссылки на компании через Selenium (не через BeautifulSoup!)
                companies = []
                try:
                    logger.debug("Поиск ссылок на компании...")
                    company_links = WebDriverWait(self.driver, 15).until(
                        EC.presence_of_all_elements_located((By.XPATH, "//a[contains(@href, '/portal/company.aspx?id=')]"))
                    )

                    logger.info(f"Найдено ссылок на компании через Selenium: {len(company_links)}")

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
                                    'element': link  # Сохраняем сам элемент для клика
                                })
                        except Exception as e:
                            logger.debug(f"Ошибка обработки ссылки: {e}")
                            continue

                except Exception as e:
                    logger.warning(f"Не удалось найти ссылки на компании через Selenium: {e}")
                    # Fallback на BeautifulSoup
                    soup = BeautifulSoup(html, 'html.parser')
                    for link in soup.find_all('a', href=re.compile(r'/portal/company\.aspx\?id=')):
                        try:
                            company_id = re.search(r'id=(\d+)', link['href']).group(1)
                            name = link.text.strip()
                            if name:
                                companies.append({
                                    'name': name,
                                    'id': company_id,
                                    'url': urljoin(self.base_url, link['href'])
                                })
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
                time.sleep(4)  # Ждем загрузки страницы (увеличено время)
                logger.info("Страница с отчетностью эмитента загружена")

                soup = BeautifulSoup(self.driver.page_source, 'html.parser')

                # Ищем ссылки на файлы отчетов в таблице
                reports = []

                # Ищем ссылки с классом file-link и href содержащим FileLoad.ashx
                download_links = soup.find_all('a', class_='file-link', href=re.compile(r'FileLoad\.ashx'))
                logger.info(f"Найдено файлов отчетов (file-link): {len(download_links)}")

                # Скачиваем файлы
                for i, link in enumerate(download_links, 1):
                    if i > 10: break # временно скачиваем последние 10 файлов
                    try:
                        file_url = link.get('href', '')
                        if not file_url.startswith('http'):
                            file_url = urljoin(self.base_url, file_url)

                        # Получаем file_id из data-fileid или из URL
                        file_id = link.get('data-fileid', '')
                        if not file_id:
                            match = re.search(r'Fileid=(\d+)', file_url)
                            if match:
                                file_id = match.group(1)

                        # Получаем информацию о файле из текста ссылки
                        file_info = link.text.strip()

                        # Ищем отчетный период в строке таблицы (родительский tr)
                        period = "unknown"
                        try:
                            tr = link.find_parent('tr')
                            if tr:
                                # Ищем ячейку с отчетным периодом (обычно это одна из первых колонок)
                                cells = tr.find_all('td')
                                if len(cells) > 0:
                                    # Предполагаем, что период в первой или второй колонке
                                    for cell in cells[:3]:
                                        cell_text = cell.text.strip()
                                        # Проверяем, похоже ли на период (содержит год, квартал и т.д.)
                                        if re.search(r'(20\d{2}|квартал|год|полугодие)', cell_text, re.IGNORECASE):
                                            period = cell_text
                                            break
                        except Exception as e:
                            logger.debug(f"Ошибка при поиске периода: {e}")

                        # Формируем имя файла
                        file_name = f"{company_name}_{period}_{file_id}.zip"
                        file_name = re.sub(r'[^\w\s.-]', '_', file_name)
                        file_name = re.sub(r'_+', '_', file_name)

                        logger.info(f"\n{i}. Период: {period}")
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
                                'url': file_url,
                                'status': 'downloaded',
                                'path': file_path,
                                'size': file_size
                            })

                            # Задержка после успешного скачивания
                            if i < len(download_links):  # Если это не последний файл
                                logger.debug("Пауза 3 секунды перед следующим файлом...")
                                time.sleep(3)

                        except Exception as e:
                            logger.error(f"   Ошибка при скачивании: {e}")
                            reports.append({
                                'name': file_name,
                                'period': period,
                                'url': file_url,
                                'status': 'error'
                            })
                            # Задержка даже после ошибки
                            time.sleep(1)

                    except Exception as e:
                        logger.error(f"   Ошибка обработки ссылки: {e}")
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
    
    def get_issuer_reports_selenium(self, company_id, download_dir='/tmp/reports'):
        """Получить отчетность эмитента через Selenium и скачать файлы"""
        import os
        
        # Создаем директорию для загрузок
        os.makedirs(download_dir, exist_ok=True)
        
        # Настраиваем Chrome для автоматического скачивания
        chrome_options = Options()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--disable-dev-shm-usage')
        
        prefs = {
            "download.default_directory": download_dir,
            "download.prompt_for_download": False,
            "download.directory_upgrade": True,
            "safebrowsing.enabled": True
        }
        chrome_options.add_experimental_option("prefs", prefs)
        
        # Создаем новый драйвер с настройками загрузки
        try:
            from webdriver_manager.chrome import ChromeDriverManager
            service = Service(ChromeDriverManager().install())
            download_driver = webdriver.Chrome(service=service, options=chrome_options)
        except Exception as e:
            logger.error(f"Не удалось создать драйвер для загрузки: {e}")
            return []

        company_url = f"{self.base_url}/portal/company.aspx?id={company_id}"

        try:
            logger.info("=== Переход на страницу компании ===")
            logger.info(f"URL: {company_url}")

            download_driver.get(company_url)
            time.sleep(5)

            logger.info("Страница загружена")

            # Делаем скриншот страницы компании
            try:
                download_driver.save_screenshot('/tmp/company_page_screenshot.png')
                logger.debug("Скриншот страницы компании сохранен")
            except Exception as e:
                logger.warning(f"Не удалось сохранить скриншот: {e}")

            # Сохраняем HTML страницы компании
            with open('/tmp/company_page.html', 'w', encoding='utf-8') as f:
                f.write(download_driver.page_source)
            logger.debug("HTML страницы компании сохранен")

            # Ищем меню "Отчетность" и наводимся на него
            try:
                # Ищем элемент меню "Отчетность"
                reporting_menu = WebDriverWait(download_driver, 10).until(
                    EC.presence_of_element_located((By.LINK_TEXT, "Отчетность"))
                )
                logger.debug("Найдено меню 'Отчетность'")

                # Наводимся на меню (hover)
                from selenium.webdriver.common.action_chains import ActionChains
                actions = ActionChains(download_driver)
                actions.move_to_element(reporting_menu).perform()
                time.sleep(2)

                logger.debug("Навелись на меню 'Отчетность'")

                # Ищем и кликаем на "Отчетность эмитента"
                issuer_reports_link = WebDriverWait(download_driver, 10).until(
                    EC.element_to_be_clickable((By.LINK_TEXT, "Отчетность эмитента"))
                )
                logger.debug("Найдена ссылка 'Отчетность эмитента'")

                issuer_reports_link.click()
                time.sleep(3)

                logger.info("Перешли в раздел 'Отчетность эмитента'")
                
                # Получаем HTML страницы с отчетами
                html = download_driver.page_source
                soup = BeautifulSoup(html, 'html.parser')

                # Ищем все ссылки для скачивания файлов
                reports = []
                download_links = soup.find_all('a', href=re.compile(r'(\.zip|\.pdf|\.xls|\.xlsx|\.doc|\.docx)', re.IGNORECASE))

                logger.info(f"Найдено файлов для скачивания: {len(download_links)}")

                for i, link in enumerate(download_links[:10], 1):  # Ограничим первыми 10 файлами
                    try:
                        file_url = urljoin(self.base_url, link.get('href', ''))
                        file_name = link.text.strip() or f"report_{i}"

                        logger.info(f"{i}. {file_name}")
                        logger.debug(f"   URL: {file_url}")

                        # Кликаем на ссылку для скачивания
                        try:
                            download_link_element = download_driver.find_element(By.XPATH, f"//a[@href='{link.get('href')}']")
                            download_link_element.click()
                            logger.debug("Клик выполнен, файл загружается...")
                            time.sleep(2)  # Даем время на начало загрузки

                            reports.append({
                                'name': file_name,
                                'url': file_url,
                                'status': 'downloading'
                            })
                        except Exception as e:
                            logger.error(f"   Ошибка при клике: {e}")
                            reports.append({
                                'name': file_name,
                                'url': file_url,
                                'status': 'error'
                            })
                    except Exception as e:
                        logger.error(f"   Ошибка обработки ссылки: {e}")

                # Ждем завершения загрузок
                logger.info("Ожидание завершения загрузок (10 секунд)...")
                time.sleep(10)

                # Проверяем загруженные файлы
                downloaded_files = os.listdir(download_dir)
                logger.info(f"Загружено файлов в {download_dir}: {len(downloaded_files)}")
                for file in downloaded_files:
                    logger.info(f"  - {file}")

                return reports

            except Exception as e:
                logger.error(f"Ошибка при работе с меню: {e}", exc_info=True)
                return []

        except Exception as e:
            logger.error(f"Ошибка при получении отчетов: {e}", exc_info=True)
            return []
        finally:
            download_driver.quit()
    
    def get_company_reports(self, company_id):
        """Получить список отчетов компании"""
        company_url = f"{self.base_url}/portal/company.aspx?id={company_id}"
        
        try:
            response = self.session.get(company_url, timeout=10)
            response.raise_for_status()
            soup = BeautifulSoup(response.text, 'html.parser')
            
            reports = []
            # Ищем раздел с финансовой отчетностью
            for link in soup.find_all('a', href=re.compile(r'financial-results')):
                reports.append({
                    'title': link.text.strip(),
                    'url': urljoin(self.base_url, link['href'])
                })
            
            return reports
        except Exception as e:
            logger.error(f"Ошибка получения отчетов: {e}")
            return []

    def get_ifrs_reports(self, company_id):
        """Получить МСФО отчеты"""
        # Прямая ссылка на раздел МСФО
        ifrs_url = f"{self.base_url}/portal/files.aspx?id={company_id}&type=4"

        try:
            response = self.session.get(ifrs_url, timeout=10)
            response.raise_for_status()
            soup = BeautifulSoup(response.text, 'html.parser')

            reports = []
            # Парсим таблицу с отчетами
            table = soup.find('table', {'class': 'table'})
            if table:
                for row in table.find_all('tr')[1:]:  # Пропускаем заголовок
                    cols = row.find_all('td')
                    if len(cols) >= 3:
                        date = cols[0].text.strip()
                        period = cols[1].text.strip()
                        file_link = cols[2].find('a')

                        if file_link:
                            reports.append({
                                'date': date,
                                'period': period,
                                'file_name': file_link.text.strip(),
                                'file_url': urljoin(self.base_url, file_link['href'])
                            })

            return reports
        except Exception as e:
            logger.error(f"Ошибка получения МСФО: {e}")
            return []

    def download_report(self, url, filename):
        """Скачать отчет"""
        try:
            response = self.session.get(url, timeout=30)
            response.raise_for_status()

            with open(filename, 'wb') as f:
                f.write(response.content)

            logger.info(f"Скачан: {filename}")
            return True
        except Exception as e:
            logger.error(f"Ошибка скачивания {filename}: {e}")
            return False

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
