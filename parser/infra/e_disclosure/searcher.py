import logging
import re
import time
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from infra.config import config

logger = logging.getLogger(__name__)


class CompanySearcher:
    """Поиск компаний на e-disclosure.ru."""

    def __init__(self, driver):
        self.driver = driver
        self.base_url = config.base_url

    def search(self, query: str) -> list[dict]:
        """Поиск компании по запросу (ИНН или название)."""
        search_url = f"{self.base_url}/poisk-po-kompaniyam"

        try:
            logger.info("=== Поиск компании через Selenium ===")
            logger.info(f"URL: {search_url}, Query: {query}")

            self.driver.get(search_url)
            time.sleep(config.timeout_between_requests)

            search_input = WebDriverWait(self.driver, config.timeout_element_wait).until(
                EC.presence_of_element_located((By.ID, "textfield"))
            )

            search_input.clear()
            search_input.send_keys(query)

            try:
                search_button = WebDriverWait(self.driver, config.timeout_page_load).until(
                    EC.element_to_be_clickable((By.ID, "sendButton"))
                )
                search_button.click()
            except Exception as e:
                logger.error(f"Ошибка при нажатии кнопки 'Искать': {e}", exc_info=True)

            time.sleep(config.timeout_between_requests)

            companies = []
            try:
                company_links = WebDriverWait(self.driver, 15).until(
                    EC.presence_of_all_elements_located(
                        (By.XPATH, "//a[contains(@href, '/portal/company.aspx?id=')]")
                    )
                )

                for link in company_links:
                    try:
                        href = link.get_attribute("href")
                        name = link.text.strip()
                        match = re.search(r"id=(\d+)", href)
                        if match and name:
                            company_id = match.group(1)
                            companies.append({
                                "name": name,
                                "id": company_id,
                                "url": href,
                                "element": link,
                            })
                    except Exception as e:
                        logger.debug(f"Ошибка обработки ссылки: {e}")
                        continue

            except Exception as e:
                logger.debug(f"Ошибка парсинга ссылок компаний: {e}")

            logger.info(f"=== Найдено компаний: {len(companies)} ===")
            for i, comp in enumerate(companies, 1):
                logger.info(f"{i}. {comp['name']} (ID: {comp['id']})")

            return companies

        except Exception as e:
            logger.error(f"Ошибка поиска: {e}", exc_info=True)
            return []
