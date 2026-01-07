import re
from typing import Tuple

from application.exceptions import PeriodParseError


QUARTER_TO_MONTHS = {
    1: 3,
    2: 6,
    3: 9,
    4: 12,
}


def extract_year_and_period(period_str: str) -> Tuple[int, int]:
    if not period_str or not isinstance(period_str, str):
        raise PeriodParseError(f"Некорректная строка периода: {period_str}")

    try:
        parts = period_str.split(',')

        if len(parts) == 1:
            return int(parts[0].strip()), 12

        if len(parts) != 2:
            raise PeriodParseError(f"Неверный формат периода: {period_str}")

        year = int(parts[0].strip())

        period_part = parts[1].strip().lower()
        tokens = period_part.split()

        if len(tokens) < 2:
            raise PeriodParseError(f"Неверный формат периода: {period_str}")

        count = int(tokens[0])
        period_type = tokens[1]

        if period_type.startswith("месяц"):
            period_months = count
        elif period_type.startswith("квартал"):
            if count not in QUARTER_TO_MONTHS:
                raise PeriodParseError(f"Некорректный номер квартала: {count}")
            period_months = QUARTER_TO_MONTHS[count]
        else:
            raise PeriodParseError(f"Неизвестный тип периода: {period_type}")

        return (year, period_months)

    except ValueError as e:
        raise PeriodParseError(f"Ошибка парсинга периода '{period_str}': {e}")
    except Exception as e:
        if isinstance(e, PeriodParseError):
            raise
        raise PeriodParseError(f"Неожиданная ошибка парсинга '{period_str}': {e}")


def normalize_filename(name: str) -> str:
    normalized = re.sub(r"[^\w\s.-]", "_", name)
    return re.sub(r"_+", "_", normalized)
