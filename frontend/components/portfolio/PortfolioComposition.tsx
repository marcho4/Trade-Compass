"use client"

import { useRouter } from "next/navigation"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { TrendingUp, TrendingDown, Minus } from "lucide-react"
import { Position } from "@/types/portfolio"

interface PortfolioCompositionProps {
  positions: Position[]
  totalValue: number
}

export const PortfolioComposition = ({
  positions,
  totalValue,
}: PortfolioCompositionProps) => {
  const router = useRouter()

  const handlePositionClick = (ticker: string) => {
    router.push(`/dashboard/${ticker}`)
  }

  if (positions.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Состав портфеля</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-12 text-center">
            <p className="text-muted-foreground">
              Портфель пуст. Добавьте позиции для начала инвестирования.
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  const calculatePositionValue = (position: Position) => {
    return position.currentPrice * position.quantity
  }

  const calculatePositionProfit = (position: Position) => {
    const currentValue = calculatePositionValue(position)
    const investedValue = position.avgPrice * position.quantity
    return currentValue - investedValue
  }

  const calculatePositionProfitPercent = (position: Position) => {
    const profit = calculatePositionProfit(position)
    const investedValue = position.avgPrice * position.quantity
    return investedValue > 0 ? (profit / investedValue) * 100 : 0
  }

  const calculatePositionWeight = (position: Position) => {
    const positionValue = calculatePositionValue(position)
    return totalValue > 0 ? (positionValue / totalValue) * 100 : 0
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("ru-RU", {
      style: "currency",
      currency: "RUB",
      maximumFractionDigits: 0,
    }).format(value)
  }

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat("ru-RU", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    }).format(date)
  }

  const renderProfitIndicator = (profitPercent: number) => {
    if (profitPercent > 0) {
      return (
        <div className="flex items-center gap-1 text-green-600">
          <TrendingUp className="h-4 w-4" />
          <span className="font-semibold">+{profitPercent.toFixed(2)}%</span>
        </div>
      )
    }
    
    if (profitPercent < 0) {
      return (
        <div className="flex items-center gap-1 text-red-600">
          <TrendingDown className="h-4 w-4" />
          <span className="font-semibold">{profitPercent.toFixed(2)}%</span>
        </div>
      )
    }
    
    return (
      <div className="flex items-center gap-1 text-muted-foreground">
        <Minus className="h-4 w-4" />
        <span className="font-semibold">0.00%</span>
      </div>
    )
  }

  // Сортируем позиции по весу в портфеле (от большего к меньшему)
  const sortedPositions = [...positions].sort(
    (a, b) => calculatePositionWeight(b) - calculatePositionWeight(a)
  )

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <CardTitle>Состав портфеля</CardTitle>
          <div className="text-sm text-muted-foreground">
            Всего позиций: {positions.length}
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {sortedPositions.map((position) => {
            const positionValue = calculatePositionValue(position)
            const profit = calculatePositionProfit(position)
            const profitPercent = calculatePositionProfitPercent(position)
            const weight = calculatePositionWeight(position)

            return (
              <div
                key={position.id}
                className="rounded-lg border bg-card p-4 hover:shadow-md transition-all cursor-pointer hover:scale-[1.01]"
                onClick={() => handlePositionClick(position.companyTicker)}
                onKeyDown={(e) => {
                  if (e.key === "Enter" || e.key === " ") {
                    handlePositionClick(position.companyTicker)
                  }
                }}
                role="button"
                tabIndex={0}
                aria-label={`Открыть анализ компании ${position.companyName}`}
              >
                <div className="space-y-3">
                  {/* Заголовок позиции */}
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <h4 className="font-semibold text-lg">
                          {position.companyTicker}
                        </h4>
                        {position.sector && (
                          <Badge variant="secondary" className="text-xs">
                            {position.sector}
                          </Badge>
                        )}
                      </div>
                      <p className="text-sm text-muted-foreground mt-1">
                        {position.companyName}
                      </p>
                    </div>
                    {renderProfitIndicator(profitPercent)}
                  </div>

                  {/* Основная информация */}
                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
                    <div>
                      <p className="text-xs text-muted-foreground">Количество</p>
                      <p className="font-medium">
                        {position.quantity} {position.quantity === 1 ? "лот" : "лотов"}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Средняя цена</p>
                      <p className="font-medium">{formatCurrency(position.avgPrice)}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Текущая цена</p>
                      <p className="font-medium">{formatCurrency(position.currentPrice)}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Последняя покупка</p>
                      <p className="font-medium text-sm">{formatDate(position.lastBuyDate)}</p>
                    </div>
                  </div>

                  {/* Доп. информация */}
                  <div className="flex flex-wrap items-center justify-between gap-4 pt-2 border-t">
                    <div className="flex items-center gap-6">
                      <div>
                        <p className="text-xs text-muted-foreground">Стоимость позиции</p>
                        <p className="font-semibold">{formatCurrency(positionValue)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Прибыль/убыток</p>
                        <p
                          className={`font-semibold ${
                            profit >= 0 ? "text-green-600" : "text-red-600"
                          }`}
                        >
                          {profit >= 0 ? "+" : ""}
                          {formatCurrency(profit)}
                        </p>
                      </div>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Доля в портфеле</p>
                      <div className="flex items-center gap-2">
                        <div className="w-24 h-2 bg-secondary rounded-full overflow-hidden">
                          <div
                            className="h-full bg-primary rounded-full"
                            style={{ width: `${Math.min(weight, 100)}%` }}
                          />
                        </div>
                        <span className="font-semibold text-sm">{weight.toFixed(1)}%</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}

