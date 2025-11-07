"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ChartContainer } from "@/components/ui/chart"
import {
  Radar,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
} from "recharts"
import { TrendingUp, TrendingDown } from "lucide-react"

interface CompanyRating {
  profitability: number // Рентабельность (ROE, ROA)
  growth: number // Рост (revenue_growth, profit_growth)
  valuation: number // Оценка (P/E, P/B, P/S)
  financial_health: number // Финансовое здоровье (debt_to_equity, current_ratio)
  efficiency: number // Эффективность (ROCE, margin)
}

interface CompanyCardProps {
  id: number
  ticker: string
  name: string
  sector: string
  price: number
  priceChange: number
  priceChangePercent: number
  rating: CompanyRating
  marketCap?: number
  pe?: number
  dividendYield?: number
  onClick?: () => void
}

export const CompanyCard = ({
  ticker,
  name,
  sector,
  price,
  priceChange,
  priceChangePercent,
  rating,
  marketCap,
  pe,
  dividendYield,
  onClick,
}: CompanyCardProps) => {
  const isPositiveChange = priceChange >= 0

  // Подготовка данных для радар чарта
  const radarData = [
    {
      metric: "Прибыль",
      value: rating.profitability,
      fullMark: 100,
    },
    {
      metric: "Рост",
      value: rating.growth,
      fullMark: 100,
    },
    {
      metric: "Оценка",
      value: rating.valuation,
      fullMark: 100,
    },
    {
      metric: "Здоровье",
      value: rating.financial_health,
      fullMark: 100,
    },
    {
      metric: "Эффект.",
      value: rating.efficiency,
      fullMark: 100,
    },
  ]

  // Средний рейтинг
  const averageRating = Math.round(
    (rating.profitability +
      rating.growth +
      rating.valuation +
      rating.financial_health +
      rating.efficiency) /
      5
  )

  const getRatingBadgeVariant = (rating: number) => {
    if (rating >= 80) return "default"
    if (rating >= 60) return "secondary"
    return "destructive"
  }

  // Функция для вычисления цвета заливки графика на основе рейтинга
  // 0-33: красный → оранжевый, 33-67: оранжевый → желтый, 67-100: желтый → зеленый
  const getChartFillColor = (rating: number) => {
    const normalizedRating = Math.max(0, Math.min(100, rating))
    
    let r: number, g: number, b: number
    
    if (normalizedRating <= 33) {
      // 0-33: Красный (#ef4444) → Оранжевый (#f97316)
      const t = normalizedRating / 33
      const redR = 239, redG = 68, redB = 68
      const orangeR = 249, orangeG = 115, orangeB = 22
      
      r = Math.round(redR + (orangeR - redR) * t)
      g = Math.round(redG + (orangeG - redG) * t)
      b = Math.round(redB + (orangeB - redB) * t)
    } else if (normalizedRating <= 67) {
      // 33-67: Оранжевый (#f97316) → Желтый (#eab308)
      const t = (normalizedRating - 33) / 34
      const orangeR = 249, orangeG = 115, orangeB = 22
      const yellowR = 234, yellowG = 179, yellowB = 8
      
      r = Math.round(orangeR + (yellowR - orangeR) * t)
      g = Math.round(orangeG + (yellowG - orangeG) * t)
      b = Math.round(orangeB + (yellowB - orangeB) * t)
    } else {
      // 67-100: Желтый (#eab308) → Зеленый (#22c55e)
      const t = (normalizedRating - 67) / 33
      const yellowR = 234, yellowG = 179, yellowB = 8
      const greenR = 34, greenG = 197, greenB = 94
      
      r = Math.round(yellowR + (greenR - yellowR) * t)
      g = Math.round(yellowG + (greenG - yellowG) * t)
      b = Math.round(yellowB + (greenB - yellowB) * t)
    }
    
    return `rgb(${r}, ${g}, ${b})`
  }

  const chartFillColor = getChartFillColor(averageRating)

  return (
    <div
      className="cursor-pointer transition-all hover:shadow-lg hover:scale-[1.02]"
      onClick={onClick}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          onClick?.()
        }
      }}
      role="button"
      tabIndex={0}
      aria-label={`Открыть анализ компании ${name}`}
    >
      <Card className="h-full">
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="space-y-1 flex-1">
              <div className="flex items-center gap-2">
                <CardTitle className="text-2xl font-bold">{ticker}</CardTitle>
                <Badge variant="outline" className="text-xs">
                  {sector}
                </Badge>
              </div>
              <p className="text-sm text-muted-foreground line-clamp-1">
                {name}
              </p>
            </div>
            <Badge
              variant={getRatingBadgeVariant(averageRating)}
              className="ml-2 h-10 w-10 flex items-center justify-center rounded-full text-base font-bold"
            >
              {averageRating}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4">
            {/* Радар чарт сверху */}
            <div className="flex flex-col items-center justify-center">
              <ChartContainer
                config={{
                  value: {
                    label: "Рейтинг",
                    color: "hsl(var(--chart-1))",
                  },
                }}
                className="h-[200px] w-full"
              >
                <RadarChart data={radarData}>
                  <PolarGrid strokeDasharray="3 3" />
                  <PolarAngleAxis
                    dataKey="metric"
                    tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
                  />
                  <PolarRadiusAxis
                    angle={90}
                    domain={[0, 100]}
                    tick={{ fontSize: 10 }}
                  />
                  <Radar
                    name="Рейтинг"
                    dataKey="value"
                    stroke={chartFillColor}
                    fill={chartFillColor}
                    fillOpacity={0.6}
                  />
                </RadarChart>
              </ChartContainer>
            </div>

            {/* Ключевые метрики снизу */}
            <div className="flex flex-col space-y-3">
              {/* Цена */}
              <div>
                <p className="text-xs text-muted-foreground">Цена</p>
                <div className="flex items-baseline gap-2">
                  <p className="text-lg font-bold">
                    {new Intl.NumberFormat("ru-RU", {
                      style: "currency",
                      currency: "RUB",
                      maximumFractionDigits: 2,
                    }).format(price)}
                  </p>
                </div>
                <div className="flex items-center gap-1 mt-1">
                  {isPositiveChange ? (
                    <TrendingUp className="h-3 w-3 text-green-600" />
                  ) : (
                    <TrendingDown className="h-3 w-3 text-red-600" />
                  )}
                  <span
                    className={`text-xs font-semibold ${
                      isPositiveChange ? "text-green-600" : "text-red-600"
                    }`}
                  >
                    {isPositiveChange ? "+" : ""}
                    {priceChangePercent.toFixed(2)}%
                  </span>
                  <span
                    className={`text-xs ${
                      isPositiveChange ? "text-green-600" : "text-red-600"
                    }`}
                  >
                    ({isPositiveChange ? "+" : ""}
                    {new Intl.NumberFormat("ru-RU", {
                      style: "currency",
                      currency: "RUB",
                      maximumFractionDigits: 2,
                    }).format(priceChange)})
                  </span>
                </div>
              </div>

              {/* Дополнительные метрики */}
              <div className="space-y-2 text-xs">
                {marketCap && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Капитализация:</span>
                    <span className="font-medium">
                      {new Intl.NumberFormat("ru-RU", {
                        notation: "compact",
                        compactDisplay: "short",
                        maximumFractionDigits: 1,
                      }).format(marketCap)} ₽
                    </span>
                  </div>
                )}
                {pe && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">P/E:</span>
                    <span className="font-medium">{pe.toFixed(2)}</span>
                  </div>
                )}
                {dividendYield && dividendYield > 0 && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Див. доход:</span>
                    <span className="font-medium text-green-600">
                      {dividendYield.toFixed(2)}%
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

