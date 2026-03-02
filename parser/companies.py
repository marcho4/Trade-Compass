COMPANIES = {
    "2309085638": "MGNT",   # Магнит
    "4823006703": "NLMK",   # НЛМК
    "7831000027": "BSPB",   # Банк Санкт-Петербург
    "7702077840": "MOEX",   # Московская биржа
    "7801268965": "SPBE",   # СПб Биржа
    "4401116480": "SVCB",   # Совкомбанк
    "7702070139": "VTBR",   # ВТБ
    "7712040126": "AFLT",   # Аэрофлот
    "7736050003": "GAZP",   # Газпром
    "7708004767": "LKOH",   # Лукойл
    "7707083893": "SBER",   # Сбербанк
    "9722079341": "X5",     # Икс 5
    "3906399157": "LENT",
    "7706107510": "ROSN",
    "7736216869": "PHOR",
    # "7536033929": "CHMF",
    "7706061801": "TRNFP",
    # "8604035373": "SNGS", 
    "1644003838": "TATN",
    "8401005730": "GMKN",
    # "7736033003": "SIBN",
    # "2463000007": "PLZL",
    "7731581426": "NVTK",

    # Tech Sector
    "7714716995": "NSVZ",
    "7736227885": "SOFL",
    "3900015862": "VKCO",
    "9715302870": "DIAS",
    "7726476459": "ASTR",
    "9722103241": "BAZA",
    "9718077239": "POSI",
    "9717163245": "DATA",
    "1683020450": "IVAT",
    "3900019850": "YDEX",
    "3900045916": "OZON",
}


TICKER_TO_INN = {ticker: inn for inn, ticker in COMPANIES.items()}


def get_ticker_by_inn(inn: str) -> str:
    return COMPANIES.get(inn, inn)


def get_inn_by_ticker(ticker: str) -> str | None:
    return TICKER_TO_INN.get(ticker)