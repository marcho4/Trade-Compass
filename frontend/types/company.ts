// Типы для анализа компаний

export interface Company {
  id: number
  ticker: string
  sectorId: number
  sector?: string
  lotSize?: number
  ceo?: string
  currentPrice?: number
  priceChange24h?: number
}

export interface FinancialMetrics {
  reportId: string
  reportYear: number
  reportPeriod: "Q1" | "Q2" | "Q3" | "Q4" | "ANNUAL"
  
  // P&L (Отчёт о прибылях и убытках)
  revenue: number // Выручка
  costOfRevenue: number // Себестоимость
  grossProfit: number // Валовая прибыль
  operatingExpenses: number // Операционные расходы
  ebit: number // Прибыль до вычета процентов и налогов
  ebitda: number // EBITDA
  interestExpense: number // Проценты к уплате
  taxExpense: number // Налоги
  netProfit: number // Чистая прибыль
  
  // Balance Sheet (Баланс)
  totalAssets: number // Всего активов
  currentAssets: number // Оборотные активы
  cashAndEquivalents: number // Денежные средства и эквиваленты
  inventories: number // Запасы
  receivables: number // Дебиторская задолженность
  
  totalLiabilities: number // Всего обязательств
  currentLiabilities: number // Краткосрочные обязательства
  debt: number // Долг (краткосрочный + долгосрочный)
  longTermDebt: number // Долгосрочный долг
  shortTermDebt: number // Краткосрочный долг
  
  equity: number // Собственный капитал
  retainedEarnings: number // Нераспределённая прибыль
  
  // Cash Flow Statement (Отчёт о движении денежных средств)
  operatingCashFlow: number // Операционный денежный поток
  investingCashFlow: number // Инвестиционный денежный поток
  financingCashFlow: number // Финансовый денежный поток
  capex: number // Капитальные затраты
  freeCashFlow: number // Свободный денежный поток (OCF - CapEx)
  
  // Market Data (для мультипликаторов)
  sharesOutstanding: number // Количество акций в обращении
  marketCap: number // Рыночная капитализация на дату отчёта
  
  // Дополнительные расчётные поля
  workingCapital: number // Оборотный капитал (current_assets - current_liabilities)
  capitalEmployed: number // Задействованный капитал (total_assets - current_liabilities)
  enterpriseValue: number // EV = market_cap + debt - cash
  netDebt: number // Чистый долг (debt - cash)
}

export interface FinancialIndicators {
  reportId: string
  
  // Мультипликаторы
  pe: number | null // Price to Earnings
  pb: number | null // Price to Book
  ps: number | null // Price to Sales
  peg: number | null // Price/Earnings to Growth
  
  // Рентабельность
  roe: number | null // Return on Equity
  roa: number | null // Return on Assets
  roce: number | null // Return on Capital Employed
  roic: number | null // Return on Invested Capital
  
  // Маржинальность
  grossProfitMargin: number | null // Валовая маржа
  operatingProfitMargin: number | null // Операционная маржа
  netProfitMargin: number | null // Чистая маржа
  
  // Enterprise Value мультипликаторы
  evEbitda: number | null
  evSales: number | null
  evFcf: number | null
  
  // Ликвидность
  currentRatio: number | null // Коэффициент текущей ликвидности
  quickRatio: number | null // Коэффициент быстрой ликвидности
  
  // Долговая нагрузка
  debtToEquity: number | null
  debtToEbitda: number | null
  netDebtToEbitda: number | null
  interestCoverage: number | null // EBIT / Interest Expense
  
  // Денежные потоки
  fcfYield: number | null // Free Cash Flow Yield
  incomeQuality: number | null // Operating CF / Net Income
  
  // Рост
  revenueGrowthYoy: number | null
  profitGrowthYoy: number | null
  epsGrowthYoy: number | null
  
  // Прочее
  eps: number | null // Earnings per Share
  dividendYield: number | null
  payoutRatio: number | null
}

export interface CompanyAnalysis {
  company: Company
  latestMetrics: FinancialMetrics
  latestIndicators: FinancialIndicators
  historicalMetrics: FinancialMetrics[]
  historicalIndicators: FinancialIndicators[]
  industryAverages?: Partial<FinancialIndicators>
}

export interface MetricDescription {
  name: string
  formula?: string
  description: string
  interpretation: string
  goodValue?: string
}

// Описания метрик для подсказок
export const METRIC_DESCRIPTIONS: Record<string, MetricDescription> = {
  pe: {
    name: "P/E (Price to Earnings)",
    formula: "Цена акции / Прибыль на акцию",
    description: "Показывает, сколько инвесторы готовы платить за 1 рубль прибыли компании",
    interpretation: "Чем ниже показатель, тем дешевле акция относительно прибыли. Однако высокий P/E может указывать на ожидания роста.",
    goodValue: "15-25 для стабильных компаний, может быть выше для растущих"
  },
  roe: {
    name: "ROE (Return on Equity)",
    formula: "Чистая прибыль / Собственный капитал × 100%",
    description: "Показывает эффективность использования собственного капитала",
    interpretation: "Чем выше, тем лучше компания использует вложения акционеров",
    goodValue: "> 15% считается хорошим показателем"
  },
  evEbitda: {
    name: "EV/EBITDA",
    formula: "Enterprise Value / EBITDA",
    description: "Показывает, за сколько лет окупится бизнес при текущем уровне EBITDA",
    interpretation: "Чем ниже, тем дешевле компания. Позволяет сравнивать компании с разной долговой нагрузкой",
    goodValue: "< 10 считается недооцененной, > 15 переоцененной"
  },
  currentRatio: {
    name: "Current Ratio (Коэффициент текущей ликвидности)",
    formula: "Оборотные активы / Краткосрочные обязательства",
    description: "Способность компании погасить краткосрочные обязательства",
    interpretation: "Показывает, сможет ли компания покрыть свои краткосрочные долги",
    goodValue: "> 1.5 считается безопасным"
  },
  debtToEquity: {
    name: "Debt/Equity (Долг к капиталу)",
    formula: "Общий долг / Собственный капитал",
    description: "Показывает финансовый леверидж компании",
    interpretation: "Чем выше, тем больше долговая нагрузка и финансовый риск",
    goodValue: "< 1 считается консервативным, < 2 умеренным"
  },
  fcfYield: {
    name: "FCF Yield (Доходность по свободному денежному потоку)",
    formula: "Свободный денежный поток / Рыночная капитализация × 100%",
    description: "Показывает, сколько свободных денег генерирует компания относительно своей стоимости",
    interpretation: "Чем выше, тем лучше. Показывает реальную способность генерировать деньги",
    goodValue: "> 5% считается хорошим"
  },
  netProfitMargin: {
    name: "Net Profit Margin (Чистая маржа)",
    formula: "Чистая прибыль / Выручка × 100%",
    description: "Показывает, какая доля выручки остается в виде чистой прибыли",
    interpretation: "Чем выше, тем эффективнее компания",
    goodValue: "> 10% считается хорошим, зависит от отрасли"
  }
}

