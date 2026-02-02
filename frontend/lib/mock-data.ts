import { Position, RiskLevel, Company, FinancialMetrics, FinancialIndicators, CompanyAnalysis, User, Subscription, UsageLimits } from "@/types"

// Демо данные для позиций портфеля
export const mockPositions: Position[] = [
  {
    id: "1",
    portfolioId: "1",
    companyId: "1",
    companyTicker: "SBER",
    companyName: "ПАО Сбербанк",
    quantity: 100,
    avgPrice: 250.5,
    currentPrice: 285.3,
    lastBuyDate: new Date(2024, 0, 15),
    sector: "Финансы",
  },
  {
    id: "2",
    portfolioId: "1",
    companyId: "2",
    companyTicker: "GAZP",
    companyName: "ПАО Газпром",
    quantity: 50,
    avgPrice: 180.2,
    currentPrice: 175.8,
    lastBuyDate: new Date(2024, 1, 20),
    sector: "Энергетика",
  },
  {
    id: "3",
    portfolioId: "1",
    companyId: "3",
    companyTicker: "LKOH",
    companyName: "ПАО НК Лукойл",
    quantity: 30,
    avgPrice: 6500.0,
    currentPrice: 7250.5,
    lastBuyDate: new Date(2024, 2, 10),
    sector: "Энергетика",
  },
  {
    id: "4",
    portfolioId: "1",
    companyId: "4",
    companyTicker: "YNDX",
    companyName: "Яндекс",
    quantity: 20,
    avgPrice: 3200.0,
    currentPrice: 3450.0,
    lastBuyDate: new Date(2024, 3, 5),
    sector: "Технологии",
  },
  {
    id: "5",
    portfolioId: "2",
    companyId: "1",
    companyTicker: "SBER",
    companyName: "ПАО Сбербанк",
    quantity: 150,
    avgPrice: 240.0,
    currentPrice: 285.3,
    lastBuyDate: new Date(2024, 2, 1),
    sector: "Финансы",
  },
  {
    id: "6",
    portfolioId: "2",
    companyId: "5",
    companyTicker: "TATN",
    companyName: "ПАО Татнефть",
    quantity: 40,
    avgPrice: 580.0,
    currentPrice: 620.5,
    lastBuyDate: new Date(2024, 3, 15),
    sector: "Энергетика",
  },
]

// Демо данные для портфелей
export const mockPortfolios: Record<
  string,
  {
    name: string
    description: string
    risk: RiskLevel
    currentValue: number
    goalValue: number
    goalDescription: string
  }
> = {
  "1": {
    name: "Основной портфель",
    description: "Сбалансированный портфель для долгосрочного роста капитала",
    risk: "moderate",
    currentValue: 1250000,
    goalValue: 2000000,
    goalDescription: "Достичь 2 млн рублей к концу 2025 года",
  },
  "2": {
    name: "Дивидендный портфель",
    description: "Портфель с акцентом на получение регулярного дохода",
    risk: "conservative",
    currentValue: 850000,
    goalValue: 1500000,
    goalDescription: "Накопить 1.5 млн рублей для пассивного дохода",
  },
}

// Генерация данных для графика производительности
export const generatePerformanceData = () => {
  const data = []
  const startDate = new Date()
  startDate.setMonth(startDate.getMonth() - 12)

  for (let i = 0; i < 365; i += 7) {
    const date = new Date(startDate)
    date.setDate(date.getDate() + i)

    // Генерируем случайные данные с трендом роста для портфеля
    const dayIndex = i / 7
    const portfolioBase = dayIndex * 0.6 + Math.sin(dayIndex * 0.1) * 3 + Math.random() * 2
    const moexBase = dayIndex * 0.4 + Math.sin(dayIndex * 0.15) * 2 + Math.random() * 1.5
    const sp500Base = dayIndex * 0.7 + Math.sin(dayIndex * 0.12) * 2.5 + Math.random() * 2
    const nasdaqBase = dayIndex * 0.8 + Math.sin(dayIndex * 0.08) * 4 + Math.random() * 2.5

    data.push({
      date: date.toISOString().split("T")[0],
      portfolio: Math.max(0, portfolioBase),
      moex: Math.max(0, moexBase),
      sp500: Math.max(0, sp500Base),
      nasdaq: Math.max(0, nasdaqBase),
    })
  }

  return data
}

// Получить позиции для конкретного портфеля
export const getPositionsByPortfolioId = (portfolioId: string): Position[] => {
  return mockPositions.filter((position) => position.portfolioId === portfolioId)
}

// Получить информацию о портфеле
export const getPortfolioById = (portfolioId: string) => {
  return mockPortfolios[portfolioId] || null
}

// Mock данные для компаний
export const mockCompanies: Record<string, Company> = {
  SBER: {
    id: 1,
    ticker: "SBER",
    sectorId: 1,
    sector: "Финансы",
    lotSize: 10,
    ceo: "Герман Греф",
    currentPrice: 285.30,
    priceChange24h: 2.4,
  },
  GAZP: {
    id: 2,
    ticker: "GAZP",
    sectorId: 2,
    sector: "Энергетика",
    lotSize: 10,
    ceo: "Алексей Миллер",
    currentPrice: 175.80,
    priceChange24h: -1.2,
  },
  LKOH: {
    id: 3,
    ticker: "LKOH",
    sectorId: 2,
    sector: "Энергетика",
    lotSize: 1,
    ceo: "Вагит Алекперов",
    currentPrice: 7250.50,
    priceChange24h: 1.8,
  },
  YNDX: {
    id: 4,
    ticker: "YNDX",
    sectorId: 3,
    sector: "Технологии",
    lotSize: 1,
    ceo: "Аркадий Волож",
    currentPrice: 3450.00,
    priceChange24h: 3.2,
  },
}

// Генерация исторических метрик
const generateHistoricalMetrics = (baseMetrics: Partial<FinancialMetrics>, quarters: number = 8): FinancialMetrics[] => {
  const metrics: FinancialMetrics[] = []
  const currentYear = new Date().getFullYear()
  
  for (let i = quarters - 1; i >= 0; i--) {
    const year = currentYear - Math.floor(i / 4)
    const quarter = (4 - (i % 4)) as 1 | 2 | 3 | 4
    const period = `Q${quarter}` as "Q1" | "Q2" | "Q3" | "Q4"
    
    // Генерируем с небольшим трендом роста
    const growthFactor = 1 + (quarters - i) * 0.03 + Math.random() * 0.1
    
    metrics.push({
      reportId: `report-${i}`,
      reportYear: year,
      reportPeriod: period,
      revenue: Math.floor((baseMetrics.revenue || 1000000) * growthFactor),
      costOfRevenue: Math.floor((baseMetrics.costOfRevenue || 600000) * growthFactor),
      grossProfit: Math.floor((baseMetrics.grossProfit || 400000) * growthFactor),
      operatingExpenses: Math.floor((baseMetrics.operatingExpenses || 200000) * growthFactor),
      ebit: Math.floor((baseMetrics.ebit || 200000) * growthFactor),
      ebitda: Math.floor((baseMetrics.ebitda || 250000) * growthFactor),
      interestExpense: Math.floor((baseMetrics.interestExpense || 10000) * growthFactor),
      taxExpense: Math.floor((baseMetrics.taxExpense || 40000) * growthFactor),
      netProfit: Math.floor((baseMetrics.netProfit || 150000) * growthFactor),
      totalAssets: Math.floor((baseMetrics.totalAssets || 5000000) * growthFactor),
      currentAssets: Math.floor((baseMetrics.currentAssets || 2000000) * growthFactor),
      cashAndEquivalents: Math.floor((baseMetrics.cashAndEquivalents || 500000) * growthFactor),
      inventories: Math.floor((baseMetrics.inventories || 300000) * growthFactor),
      receivables: Math.floor((baseMetrics.receivables || 400000) * growthFactor),
      totalLiabilities: Math.floor((baseMetrics.totalLiabilities || 3000000) * growthFactor),
      currentLiabilities: Math.floor((baseMetrics.currentLiabilities || 1000000) * growthFactor),
      debt: Math.floor((baseMetrics.debt || 1500000) * growthFactor),
      longTermDebt: Math.floor((baseMetrics.longTermDebt || 1000000) * growthFactor),
      shortTermDebt: Math.floor((baseMetrics.shortTermDebt || 500000) * growthFactor),
      equity: Math.floor((baseMetrics.equity || 2000000) * growthFactor),
      retainedEarnings: Math.floor((baseMetrics.retainedEarnings || 1000000) * growthFactor),
      operatingCashFlow: Math.floor((baseMetrics.operatingCashFlow || 300000) * growthFactor),
      investingCashFlow: Math.floor((baseMetrics.investingCashFlow || -100000) * growthFactor),
      financingCashFlow: Math.floor((baseMetrics.financingCashFlow || -50000) * growthFactor),
      capex: Math.floor((baseMetrics.capex || 80000) * growthFactor),
      freeCashFlow: Math.floor((baseMetrics.freeCashFlow || 220000) * growthFactor),
      sharesOutstanding: baseMetrics.sharesOutstanding || 1000000000,
      marketCap: Math.floor((baseMetrics.marketCap || 6000000) * growthFactor),
      workingCapital: Math.floor((baseMetrics.workingCapital || 1000000) * growthFactor),
      capitalEmployed: Math.floor((baseMetrics.capitalEmployed || 4000000) * growthFactor),
      enterpriseValue: Math.floor((baseMetrics.enterpriseValue || 7000000) * growthFactor),
      netDebt: Math.floor((baseMetrics.netDebt || 1000000) * growthFactor),
    })
  }
  
  return metrics
}

// Генерация индикаторов на основе метрик
const generateIndicatorsFromMetrics = (metrics: FinancialMetrics): FinancialIndicators => {
  return {
    reportId: metrics.reportId,
    pe: metrics.netProfit > 0 ? parseFloat((metrics.marketCap / metrics.netProfit).toFixed(2)) : null,
    pb: metrics.equity > 0 ? parseFloat((metrics.marketCap / metrics.equity).toFixed(2)) : null,
    ps: metrics.revenue > 0 ? parseFloat((metrics.marketCap / metrics.revenue).toFixed(2)) : null,
    peg: null,
    roe: metrics.equity > 0 ? parseFloat(((metrics.netProfit / metrics.equity) * 100).toFixed(2)) : null,
    roa: metrics.totalAssets > 0 ? parseFloat(((metrics.netProfit / metrics.totalAssets) * 100).toFixed(2)) : null,
    roce: null,
    roic: null,
    grossProfitMargin: metrics.revenue > 0 ? parseFloat(((metrics.grossProfit / metrics.revenue) * 100).toFixed(2)) : null,
    operatingProfitMargin: metrics.revenue > 0 ? parseFloat(((metrics.ebit / metrics.revenue) * 100).toFixed(2)) : null,
    netProfitMargin: metrics.revenue > 0 ? parseFloat(((metrics.netProfit / metrics.revenue) * 100).toFixed(2)) : null,
    evEbitda: metrics.ebitda > 0 ? parseFloat((metrics.enterpriseValue / metrics.ebitda).toFixed(2)) : null,
    evSales: metrics.revenue > 0 ? parseFloat((metrics.enterpriseValue / metrics.revenue).toFixed(2)) : null,
    evFcf: metrics.freeCashFlow > 0 ? parseFloat((metrics.enterpriseValue / metrics.freeCashFlow).toFixed(2)) : null,
    currentRatio: metrics.currentLiabilities > 0 ? parseFloat((metrics.currentAssets / metrics.currentLiabilities).toFixed(2)) : null,
    quickRatio: metrics.currentLiabilities > 0 ? parseFloat(((metrics.currentAssets - metrics.inventories) / metrics.currentLiabilities).toFixed(2)) : null,
    debtToEquity: metrics.equity > 0 ? parseFloat((metrics.debt / metrics.equity).toFixed(2)) : null,
    debtToEbitda: metrics.ebitda > 0 ? parseFloat((metrics.debt / metrics.ebitda).toFixed(2)) : null,
    netDebtToEbitda: metrics.ebitda > 0 ? parseFloat((metrics.netDebt / metrics.ebitda).toFixed(2)) : null,
    interestCoverage: metrics.interestExpense > 0 ? parseFloat((metrics.ebit / metrics.interestExpense).toFixed(2)) : null,
    fcfYield: metrics.marketCap > 0 ? parseFloat(((metrics.freeCashFlow / metrics.marketCap) * 100).toFixed(2)) : null,
    incomeQuality: metrics.netProfit > 0 ? parseFloat((metrics.operatingCashFlow / metrics.netProfit).toFixed(2)) : null,
    revenueGrowthYoy: null,
    profitGrowthYoy: null,
    epsGrowthYoy: null,
    eps: metrics.sharesOutstanding > 0 ? parseFloat((metrics.netProfit / metrics.sharesOutstanding).toFixed(2)) : null,
    dividendYield: 5.2,
    payoutRatio: 0.45,
  }
}

// Mock данные анализа компаний
export const getMockCompanyAnalysis = (ticker: string): CompanyAnalysis | null => {
  const company = mockCompanies[ticker]
  if (!company) return null
  
  // Базовые метрики в зависимости от компании
  const baseMetrics: Record<string, Partial<FinancialMetrics>> = {
    SBER: {
      revenue: 4_500_000_000,
      costOfRevenue: 2_000_000_000,
      grossProfit: 2_500_000_000,
      operatingExpenses: 1_200_000_000,
      ebit: 1_300_000_000,
      ebitda: 1_400_000_000,
      netProfit: 1_000_000_000,
      totalAssets: 35_000_000_000,
      equity: 4_500_000_000,
      debt: 5_000_000_000,
      cashAndEquivalents: 3_000_000_000,
      freeCashFlow: 800_000_000,
      marketCap: 6_000_000_000,
      sharesOutstanding: 21_000_000_000,
    },
    GAZP: {
      revenue: 8_000_000_000,
      costOfRevenue: 4_500_000_000,
      grossProfit: 3_500_000_000,
      operatingExpenses: 1_800_000_000,
      ebit: 1_700_000_000,
      ebitda: 2_000_000_000,
      netProfit: 1_200_000_000,
      totalAssets: 20_000_000_000,
      equity: 10_000_000_000,
      debt: 3_000_000_000,
      cashAndEquivalents: 1_500_000_000,
      freeCashFlow: 1_000_000_000,
      marketCap: 4_000_000_000,
      sharesOutstanding: 23_000_000_000,
    },
    LKOH: {
      revenue: 9_000_000_000,
      costOfRevenue: 5_000_000_000,
      grossProfit: 4_000_000_000,
      operatingExpenses: 2_000_000_000,
      ebit: 2_000_000_000,
      ebitda: 2_300_000_000,
      netProfit: 1_500_000_000,
      totalAssets: 12_000_000_000,
      equity: 7_000_000_000,
      debt: 2_000_000_000,
      cashAndEquivalents: 2_000_000_000,
      freeCashFlow: 1_200_000_000,
      marketCap: 5_500_000_000,
      sharesOutstanding: 760_000_000,
    },
    YNDX: {
      revenue: 500_000_000,
      costOfRevenue: 200_000_000,
      grossProfit: 300_000_000,
      operatingExpenses: 200_000_000,
      ebit: 100_000_000,
      ebitda: 120_000_000,
      netProfit: 80_000_000,
      totalAssets: 800_000_000,
      equity: 500_000_000,
      debt: 100_000_000,
      cashAndEquivalents: 200_000_000,
      freeCashFlow: 70_000_000,
      marketCap: 1_200_000_000,
      sharesOutstanding: 350_000_000,
    },
  }
  
  const metrics = baseMetrics[ticker] || baseMetrics.SBER
  const historicalMetrics = generateHistoricalMetrics(metrics, 8)
  const latestMetrics = historicalMetrics[historicalMetrics.length - 1]
  const historicalIndicators = historicalMetrics.map(generateIndicatorsFromMetrics)
  
  return {
    company,
    latestMetrics,
    latestIndicators: historicalIndicators[historicalIndicators.length - 1],
    historicalMetrics,
    historicalIndicators,
    industryAverages: {
      pe: 8.5,
      pb: 0.9,
      roe: 18.5,
      debtToEquity: 0.8,
      currentRatio: 1.8,
      netProfitMargin: 15.2,
    },
  }
}

// Mock данные пользователя
export const mockUser: User = {
  id: "user-123",
  name: "Иван Иванов",
  email: "ivan@example.com",
  createdAt: new Date(2024, 0, 15),
  lastLoginAt: new Date(),
}

// Mock данные подписки
export const mockSubscription: Subscription = {
  id: "sub-123",
  userId: "user-123",
  level: "premium",
  startDate: new Date(2024, 5, 1),
  endDate: new Date(2025, 5, 1),
  isActive: true,
  price: 490,
  renewsAt: new Date(2025, 5, 1),
}

// Mock данные лимитов
export const mockUsageLimits: UsageLimits = {
  aiQueries: {
    used: 45,
    limit: 100,
    resetsAt: new Date(2025, 0, 1),
  },
  companyAnalyses: {
    used: 12,
    limit: -1, // безлимит
    resetsAt: new Date(2025, 0, 1),
  },
  portfolios: {
    used: 3,
    limit: 5,
  },
  alerts: {
    used: 4,
    limit: 10,
  },
}

