import logging

from companies import COMPANIES
from selenium_client import EDisclosureParser

# Настройка логгирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

if __name__ == "__main__":
    parser = EDisclosureParser()

    for company in COMPANIES:
        logger.info(f"Поиск компании... {company}")
        companies = parser.search_company(company)

        if companies:
            logger.info("="*60)
            logger.info(f"Найдено компаний: {len(companies)}")
            for i, company in enumerate(companies[:5], 1):
                logger.info(f"{i}. {company['name']} (ID: {company['id']})")
            logger.info("="*60)

            first_company = companies[0]
            logger.info(f"Выбрана компания: {first_company['name']} (ID: {first_company['id']})")

            logger.info("="*60)
            logger.info("Получение отчетности эмитента...")
            logger.info("="*60)

            reports = parser.get_issuer_reports_by_click(first_company, download_dir='./downloads')

            logger.info("="*60)
            logger.info("ИТОГОВАЯ СТАТИСТИКА")
            logger.info("="*60)
            logger.info(f"Всего файлов обработано: {len(reports)}")

            downloaded = [r for r in reports if r['status'] == 'downloaded']
            errors = [r for r in reports if r['status'] == 'error']

            logger.info(f"Успешно скачано: {len(downloaded)}")
            logger.info(f"Ошибок: {len(errors)}")
            logger.info("="*60)

            if downloaded:
                logger.info("СКАЧАННЫЕ ФАЙЛЫ:")
                for report in downloaded:
                    size_mb = report.get('size', 0) / 1024 / 1024
                    logger.info(f"{report['period']} - {size_mb:.2f} MB")
                    logger.info(f"  {report['path']}")

            if errors:
                logger.warning("ОШИБКИ:")
                for report in errors:
                    logger.warning(f"{report['name']}")

            logger.info("="*60)

        else:
            logger.warning("Компания не найдена")

    parser.unzip_downloaded_files()
# import pdfplumber


# def print_table(table):
#     """
#     Аккуратно форматирует и выводит таблицу по столбцам.

#     Args:
#         table: Список списков, представляющий таблицу
#     """
#     if not table or not any(table):
#         print("Пустая таблица")
#         return

#     # Преобразуем None в пустые строки и все значения в строки
#     formatted_table = []
#     for row in table:
#         if row is None:
#             continue
#         formatted_row = [str(cell) if cell is not None else '' for cell in row]
#         formatted_table.append(formatted_row)

#     if not formatted_table:
#         print("Пустая таблица")
#         return

#     # Находим количество столбцов
#     max_cols = max(len(row) for row in formatted_table)

#     # Дополняем короткие строки пустыми ячейками
#     for row in formatted_table:
#         while len(row) < max_cols:
#             row.append('')

#     # Вычисляем максимальную ширину для каждого столбца
#     col_widths = [0] * max_cols
#     for row in formatted_table:
#         for i, cell in enumerate(row):
#             col_widths[i] = max(col_widths[i], len(cell))

#     # Выводим таблицу с разделителями
#     separator = '+' + '+'.join('-' * (width + 2) for width in col_widths) + '+'

#     print(separator)
#     for row in formatted_table:
#         cells = []
#         for i, cell in enumerate(row):
#             # Выравниваем числа по правому краю, текст по левому
#             if cell and cell.replace(' ', '').replace('-', '').replace(',', '').replace('.', '').replace('(', '').replace(')', '').isdigit():
#                 cells.append(f" {cell.rjust(col_widths[i])} ")
#             else:
#                 cells.append(f" {cell.ljust(col_widths[i])} ")
#         print('|' + '|'.join(cells) + '|')
#         print(separator)


# with pdfplumber.open("parser/downloads/T-Technologies_IFRS REVIEW RPRT_6m2025_rus.pdf") as pdf:
    # for page in pdf.pages:
    #     tables = page.extract_tables()
    #     print(f"\n{'='*80}")
    #     print(f"Количество таблиц на странице: {len(tables)}")
    #     print(f"{'='*80}\n")

    #     for idx, table in enumerate(tables, 1):
    #         print(f"\nТаблица {idx}:")
    #         print_table(table)

        # for table in tables:
        #     for row in table:
        #         if "Выручка" in str(row[0]) or "Revenue" in str(row[0]):
        #             revenue = row[1]  # значение за период