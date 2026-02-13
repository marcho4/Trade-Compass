COMPANIES = {
    # "2309085638": "MGNT",   # Магнит
    # "4823006703": "NLMK",   # НЛМК
    # "7831000027": "BSPB",   # Банк Санкт-Петербург
    "7702077840": "MOEX",   # Московская биржа
    # "7801268965": "SPBE",   # СПб Биржа
    # "4401116480": "SVCB",   # Совкомбанк
    # "7702070139": "VTBR",   # ВТБ
    # "7712040126": "AFLT",   # Аэрофлот
    "7736050003": "GAZP",   # Газпром
    "7708004767": "LKOH",   # Лукойл
    "7707083893": "SBER",   # Сбербанк
    "9722079341": "X5", 
    "3900019850": "YDEX"
}


TICKER_TO_INN = {ticker: inn for inn, ticker in COMPANIES.items()}


def get_ticker_by_inn(inn: str) -> str:
    return COMPANIES.get(inn, inn)


def get_inn_by_ticker(ticker: str) -> str | None:
    return TICKER_TO_INN.get(ticker)