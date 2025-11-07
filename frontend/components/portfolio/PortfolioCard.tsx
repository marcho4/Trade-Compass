"use client"

import { TrendingUp, TrendingDown } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

interface PortfolioCardProps {
  name: string
  value: number
  createdAt: Date
  profitPercent: number
  profitAmount: number
  rating: number
  onClick?: () => void
}

export const PortfolioCard = ({
  name,
  value,
  createdAt,
  profitPercent,
  profitAmount,
  rating,
  onClick,
}: PortfolioCardProps) => {
  const isProfit = profitPercent >= 0
  const formattedDate = new Intl.DateTimeFormat("ru-RU", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  }).format(createdAt)

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
      aria-label={`Открыть портфель ${name}`}
    >
      <Card className="h-full">
        <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
          <div className="space-y-1">
            <CardTitle className="text-xl font-bold">{name}</CardTitle>
            <p className="text-sm text-muted-foreground">Создан {formattedDate}</p>
          </div>
          <div className="flex items-center gap-2">
            <Badge
              variant="secondary"
              className="h-8 w-8 flex items-center justify-center rounded-full text-sm font-semibold"
            >
              {rating}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div>
              <p className="text-sm text-muted-foreground">Стоимость портфеля</p>
              <p className="text-2xl font-bold">
                {new Intl.NumberFormat("ru-RU", {
                  style: "currency",
                  currency: "RUB",
                  maximumFractionDigits: 0,
                }).format(value)}
              </p>
            </div>
            <div className="flex items-center gap-2">
              {isProfit ? (
                <TrendingUp className="h-4 w-4 text-green-600" />
              ) : (
                <TrendingDown className="h-4 w-4 text-red-600" />
              )}
              <div className="flex items-center gap-2">
                <span
                  className={`text-lg font-semibold ${
                    isProfit ? "text-green-600" : "text-red-600"
                  }`}
                >
                  {isProfit ? "+" : ""}
                  {profitPercent.toFixed(2)}%
                </span>
                <span
                  className={`text-sm ${
                    isProfit ? "text-green-600" : "text-red-600"
                  }`}
                >
                  ({isProfit ? "+" : ""}
                  {new Intl.NumberFormat("ru-RU", {
                    style: "currency",
                    currency: "RUB",
                    maximumFractionDigits: 0,
                  }).format(profitAmount)})
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

