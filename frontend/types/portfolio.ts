export interface Position {
  id: string
  portfolioId: string
  companyId: string
  companyTicker: string
  companyName: string
  quantity: number // количество лотов
  avgPrice: number // средняя цена покупки
  currentPrice: number // текущая цена
  lastBuyDate: Date
  sector?: string
  createdAt?: Date
  updatedAt?: Date
}

export interface Portfolio {
  id: string
  name: string
  userId: string
  description?: string
  value: number // текущая стоимость портфеля
  createdAt: Date
  updatedAt?: Date
  profitPercent?: number // процент прибыли/убытка
  profitAmount?: number // сумма прибыли/убытка
  rating?: number // рейтинг от платформы (1-10)
  positions?: Position[]
}

export interface PortfolioGoalData {
  currentValue: number
  goalValue: number
  goalDescription: string
}

export type RiskLevel = "conservative" | "moderate" | "aggressive"

export interface PortfolioSettings {
  riskLevel: RiskLevel
  goal?: PortfolioGoalData
}

