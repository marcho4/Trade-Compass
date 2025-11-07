// Типы для пользователя и подписки

export type SubscriptionLevel = "free" | "premium" | "pro"

export interface User {
  id: string
  name: string
  email: string
  createdAt: Date
  lastLoginAt: Date
}

export interface Subscription {
  id: string
  userId: string
  level: SubscriptionLevel
  startDate: Date
  endDate: Date | null
  isActive: boolean
  price?: number
  renewsAt?: Date
}

export interface UsageLimits {
  // AI запросы
  aiQueries: {
    used: number
    limit: number
    resetsAt: Date
  }
  // Анализы компаний
  companyAnalyses: {
    used: number
    limit: number
    resetsAt: Date
  }
  // Портфели
  portfolios: {
    used: number
    limit: number
  }
  // Алерты и уведомления
  alerts: {
    used: number
    limit: number
  }
}

export interface SubscriptionFeature {
  name: string
  included: boolean
  description?: string
}

export interface SubscriptionPlan {
  id: string
  level: SubscriptionLevel
  name: string
  price: number
  period: "month" | "year"
  features: SubscriptionFeature[]
  limits: {
    aiQueries: number
    companyAnalyses: number
    portfolios: number
    alerts: number
  }
  isPopular?: boolean
}

// Планы подписки
export const SUBSCRIPTION_PLANS: SubscriptionPlan[] = [
  {
    id: "free",
    level: "free",
    name: "Бесплатный",
    price: 0,
    period: "month",
    features: [
      { name: "Базовые метрики всех компаний", included: true },
      { name: "1 портфель", included: true },
      { name: "Скринер акций", included: true },
      { name: "10 AI запросов в месяц", included: true },
      { name: "1 полный анализ компании", included: true },
      { name: "AI анализ отчетов", included: false },
      { name: "Алерты и уведомления", included: false },
      { name: "Экспорт данных", included: false },
    ],
    limits: {
      aiQueries: 10,
      companyAnalyses: 1,
      portfolios: 1,
      alerts: 0,
    },
  },
  {
    id: "premium",
    level: "premium",
    name: "Premium",
    price: 490,
    period: "month",
    features: [
      { name: "Все функции бесплатного", included: true },
      { name: "5 портфелей", included: true },
      { name: "100 AI запросов в месяц", included: true },
      { name: "Безлимитный анализ компаний", included: true },
      { name: "AI анализ отчетов", included: true },
      { name: "10 алертов", included: true },
      { name: "Экспорт данных", included: true },
      { name: "Приоритетная поддержка", included: false },
    ],
    limits: {
      aiQueries: 100,
      companyAnalyses: -1, // -1 = безлимит
      portfolios: 5,
      alerts: 10,
    },
    isPopular: true,
  },
  {
    id: "pro",
    level: "pro",
    name: "Professional",
    price: 1990,
    period: "month",
    features: [
      { name: "Все функции Premium", included: true },
      { name: "Безлимитные портфели", included: true },
      { name: "Безлимитные AI запросы", included: true },
      { name: "Безлимитные алерты", included: true },
      { name: "API доступ", included: true },
      { name: "Приоритетная поддержка", included: true },
      { name: "Кастомные отчеты", included: true },
      { name: "Белая метка", included: true },
    ],
    limits: {
      aiQueries: -1,
      companyAnalyses: -1,
      portfolios: -1,
      alerts: -1,
    },
  },
]

