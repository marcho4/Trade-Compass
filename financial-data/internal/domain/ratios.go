package domain

// Ratios содержит все финансовые коэффициенты для анализа компании
type Ratios struct {
	// Основные мультипликаторы оценки
	// P/E - Отношение цены к прибыли (Price to Earnings)
	// Формула: Цена акции / EPS
	PriceToEarnings *float64 `json:"priceToEarnings,omitempty"`

	// P/BV - Отношение цены к балансовой стоимости (Price to Book Value)
	// Формула: Рыночная цена / Балансовая стоимость на акцию
	PriceToBook *float64 `json:"priceToBook,omitempty"`

	// P/CF - Отношение цены к денежному потоку (Price to Cash Flow)
	// Формула: Рыночная цена / Денежный поток на акцию
	PriceToCashFlow *float64 `json:"priceToCashFlow,omitempty"`

	// EV/EBITDA - Стоимость предприятия к EBITDA
	// Формула: Enterprise Value / EBITDA
	// Показывает сколько лет нужно для окупаемости при текущей EBITDA
	EVToEBITDA *float64 `json:"evToEbitda,omitempty"`

	// EV/Sales - Стоимость предприятия к выручке
	// Формула: Enterprise Value / Выручка
	EVToSales *float64 `json:"evToSales,omitempty"`

	// EV/FCF - Стоимость предприятия к свободному денежному потоку
	// Формула: Enterprise Value / Free Cash Flow
	EVToFCF *float64 `json:"evToFcf,omitempty"`

	// PEG - Отношение P/E к темпу роста прибыли
	// Формула: P/E / Темп роста прибыли
	PEG *float64 `json:"peg,omitempty"`

	// Показатели прибыльности (Profitability Ratios)
	// ROE - Рентабельность собственного капитала (Return on Equity)
	// Формула: Чистая прибыль / Собственный капитал × 100%
	ROE *float64 `json:"roe,omitempty"`

	// ROA - Рентабельность активов (Return on Assets)
	// Формула: Чистая прибыль / Активы × 100%
	ROA *float64 `json:"roa,omitempty"`

	// ROIC - Рентабельность инвестированного капитала (Return on Invested Capital)
	// Формула: NOPAT / Инвестированный капитал × 100%
	ROIC *float64 `json:"roic,omitempty"`

	// Gross Profit Margin - Валовая маржа
	// Формула: Валовая прибыль / Выручка × 100%
	GrossProfitMargin *float64 `json:"grossProfitMargin,omitempty"`

	// Operating Profit Margin - Операционная маржа
	// Формула: Операционная прибыль / Выручка × 100%
	OperatingProfitMargin *float64 `json:"operatingProfitMargin,omitempty"`

	// Net Profit Margin - Чистая маржа
	// Формула: Чистая прибыль / Выручка × 100%
	NetProfitMargin *float64 `json:"netProfitMargin,omitempty"`

	// Показатели ликвидности (Liquidity Ratios)
	// Current Ratio - Коэффициент текущей ликвидности
	// Формула: Оборотные активы / Краткосрочные обязательства
	// Показывает способность покрыть краткосрочные обязательства
	CurrentRatio *float64 `json:"currentRatio,omitempty"`

	// Quick Ratio - Коэффициент быстрой ликвидности
	// Формула: (Денежные средства + Краткосрочные финвложения + Дебиторка) / Краткосрочные обязательства
	// Должен быть >= 1
	QuickRatio *float64 `json:"quickRatio,omitempty"`

	// Показатели долговой нагрузки (Leverage Ratios)
	// Net Debt/EBITDA - Долговая нагрузка относительно прибыли
	// Формула: (Долг - Денежные средства) / EBITDA
	NetDebtToEBITDA *float64 `json:"netDebtToEbitda,omitempty"`

	// Debt to Equity - Отношение долга к собственному капиталу
	// Формула: Долговые обязательства / Собственный капитал
	DebtToEquity *float64 `json:"debtToEquity,omitempty"`

	// Interest Coverage Ratio - Коэффициент покрытия процентов
	// Формула: EBIT / Процентные расходы
	// Показывает способность покрывать процентные платежи
	InterestCoverageRatio *float64 `json:"interestCoverageRatio,omitempty"`

	// Показатели качества прибыли и эффективности
	// Income Quality - Качество прибыли
	// Формула: Денежный поток от операционной деятельности / Чистая прибыль
	IncomeQuality *float64 `json:"incomeQuality,omitempty"`

	// Asset Turnover - Оборачиваемость активов
	// Формула: Выручка / Средняя стоимость активов
	AssetTurnover *float64 `json:"assetTurnover,omitempty"`

	// Inventory Turnover - Оборачиваемость запасов
	// Формула: Себестоимость / Средняя стоимость запасов
	InventoryTurnover *float64 `json:"inventoryTurnover,omitempty"`

	// Receivables Turnover - Оборачиваемость дебиторской задолженности
	// Формула: Выручка / Средняя дебиторская задолженность
	ReceivablesTurnover *float64 `json:"receivablesTurnover,omitempty"`

	// Показатели на акцию (Per Share Metrics)
	// EPS - Прибыль на акцию (Earnings Per Share)
	// Формула: (Чистая прибыль - Дивиденды по привилегированным акциям) / Количество обыкновенных акций
	EPS *float64 `json:"eps,omitempty"`

	// Book Value Per Share - Балансовая стоимость на акцию
	// Формула: Собственный капитал / Количество акций
	BookValuePerShare *float64 `json:"bookValuePerShare,omitempty"`

	// Cash Flow Per Share - Денежный поток на акцию
	// Формула: Операционный денежный поток / Количество акций
	CashFlowPerShare *float64 `json:"cashFlowPerShare,omitempty"`

	// Dividend Per Share - Дивиденд на акцию
	DividendPerShare *float64 `json:"dividendPerShare,omitempty"`

	// Dividend Yield - Дивидендная доходность
	// Формула: Дивиденд на акцию / Цена акции × 100%
	DividendYield *float64 `json:"dividendYield,omitempty"`

	// Payout Ratio - Коэффициент выплаты дивидендов
	// Формула: Дивиденды / Чистая прибыль × 100%
	PayoutRatio *float64 `json:"payoutRatio,omitempty"`

	// Абсолютные показатели
	// Enterprise Value - Стоимость предприятия
	// Формула: Рыночная капитализация + Долг + Пенсионные обязательства + Миноритарный интерес - Денежные средства
	EnterpriseValue *float64 `json:"enterpriseValue,omitempty"`

	// Market Cap - Рыночная капитализация
	// Формула: Цена акции × Количество акций в обращении
	MarketCap *float64 `json:"marketCap,omitempty"`

	// Free Cash Flow - Свободный денежный поток
	// Формула: Денежный поток от операционной деятельности - CAPEX
	FreeCashFlow *float64 `json:"freeCashFlow,omitempty"`

	// CAPEX - Капитальные затраты
	// Формула: Приобретение ОС + Приобретение НМА
	CAPEX *float64 `json:"capex,omitempty"`

	// EBITDA - Прибыль до вычета процентов, налогов, износа и амортизации
	// Формула: Операционная прибыль + Амортизация
	EBITDA *float64 `json:"ebitda,omitempty"`

	// Net Debt - Чистый долг
	// Формула: Долговые обязательства - Денежные средства
	NetDebt *float64 `json:"netDebt,omitempty"`

	// Working Capital - Оборотный капитал
	// Формула: Оборотные активы - Краткосрочные обязательства
	WorkingCapital *float64 `json:"workingCapital,omitempty"`

	// Показатели роста (Growth Metrics) - в процентах год к году
	// Revenue Growth - Рост выручки
	RevenueGrowth *float64 `json:"revenueGrowth,omitempty"`

	// Earnings Growth - Рост прибыли
	EarningsGrowth *float64 `json:"earningsGrowth,omitempty"`

	// EBITDA Growth - Рост EBITDA
	EBITDAGrowth *float64 `json:"ebitdaGrowth,omitempty"`

	// FCF Growth - Рост свободного денежного потока
	FCFGrowth *float64 `json:"fcfGrowth,omitempty"`
}
