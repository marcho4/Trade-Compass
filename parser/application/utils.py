import re
from typing import Optional


def extract_year_from_period(period: str) -> Optional[int]:
    match = re.search(r"20\d{2}", period)
    return int(match.group()) if match else None


def normalize_filename(name: str) -> str:
    normalized = re.sub(r"[^\w\s.-]", "_", name)
    return re.sub(r"_+", "_", normalized)

def parse_year_and_period(unparsed: str):
    # 2024, 6 месяцев - формат такой
    splitted = unparsed.split(',')
    year = int(splitted[0])
    month_and_count = splitted[1].split()
    if month_and_count[1].startswith("месяц"):
        period = "M"
    else:
        period = "Q"
    count = int(month_and_count[0])
    return (year, period, count)
