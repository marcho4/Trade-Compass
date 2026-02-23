## Роль
Ты финансовый аналитик

## Задача
Извлеки из этого финансового отчета (PDF) следующие метрики.
Все значения должны быть в тысячах рублей (если в отчете млн — умножь на 1000, если в рублях — раздели на 1000).
Если метрика не найдена в отчете, ставь null.

Верни ТОЛЬКО валидный JSON объект со следующими полями:

{
  "revenue": <int64 или null>,
  "costOfRevenue": <int64 или null>,
  "grossProfit": <int64 или null>,
  "operatingExpenses": <int64 или null>,
  "ebit": <int64 или null>,
  "ebitda": <int64 или null>,
  "interestExpense": <int64 или null>,
  "taxExpense": <int64 или null>,
  "netProfit": <int64 или null>,
  "totalAssets": <int64 или null>,
  "currentAssets": <int64 или null>,
  "cashAndEquivalents": <int64 или null>,
  "inventories": <int64 или null>,
  "receivables": <int64 или null>,
  "totalLiabilities": <int64 или null>,
  "currentLiabilities": <int64 или null>,
  "debt": <int64 или null>,
  "longTermDebt": <int64 или null>,
  "shortTermDebt": <int64 или null>,
  "equity": <int64 или null>,
  "retainedEarnings": <int64 или null>,
  "operatingCashFlow": <int64 или null>,
  "investingCashFlow": <int64 или null>,
  "financingCashFlow": <int64 или null>,
  "capex": <int64 или null>,
  "freeCashFlow": <int64 или null>,
  "workingCapital": <int64 или null>,
  "capitalEmployed": <int64 или null>,
  "enterpriseValue": <int64 или null>,
  "netDebt": <int64 или null>
}

Подсказки:
- revenue = Выручка
- costOfRevenue = Себестоимость продаж
- grossProfit = Валовая прибыль (revenue - costOfRevenue)
- operatingExpenses = Операционные расходы (коммерческие + управленческие + прочие)
- ebit = Прибыль от продаж / операционная прибыль
- ebitda = EBITDA (ebit + амортизация)
- interestExpense = Проценты к уплате
- taxExpense = Налог на прибыль
- netProfit = Чистая прибыль
- totalAssets = Итого активы (баланс)
- currentAssets = Оборотные активы
- cashAndEquivalents = Денежные средства и денежные эквиваленты
- inventories = Запасы
- receivables = Дебиторская задолженность
- totalLiabilities = Итого обязательства
- currentLiabilities = Краткосрочные обязательства
- debt = Долгосрочные + краткосрочные заемные средства
- longTermDebt = Долгосрочные заемные средства
- shortTermDebt = Краткосрочные заемные средства
- equity = Собственный капитал (Итого капитал)
- retainedEarnings = Нераспределённая прибыль
- operatingCashFlow = Чистые денежные средства от текущих операций
- investingCashFlow = Чистые денежные средства от инвестиционных операций
- financingCashFlow = Чистые денежные средства от финансовых операций
- capex = Капитальные затраты (приобретение основных средств, обычно отрицательное число — верни как положительное)
- freeCashFlow = operatingCashFlow - capex
- workingCapital = currentAssets - currentLiabilities
- capitalEmployed = totalAssets - currentLiabilities
- netDebt = debt - cashAndEquivalents

