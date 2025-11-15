"use client"

import { Card } from "@/components/ui/card"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import { Badge } from "@/components/ui/badge"
import { Info, TrendingUp, TrendingDown } from "lucide-react"
import { MetricDescription } from "@/types"

interface MetricCardProps {
  label: string
  value: number | null
  format?: "number" | "percent" | "currency" | "ratio"
  description?: MetricDescription
  comparisonValue?: number | null
  comparisonLabel?: string
  trend?: "up" | "down" | "neutral"
}

export const MetricCard = ({
  label,
  value,
  format = "number",
  description,
  comparisonValue,
  comparisonLabel = "Индустрия",
  trend,
}: MetricCardProps) => {
  const formatValue = (val: number | null): string => {
    if (val === null) return "—"

    switch (format) {
      case "percent":
        return `${val.toFixed(2)}%`
      case "currency":
        return new Intl.NumberFormat("ru-RU", {
          style: "currency",
          currency: "RUB",
          maximumFractionDigits: 0,
        }).format(val)
      case "ratio":
        return val.toFixed(2)
      default:
        return new Intl.NumberFormat("ru-RU", {
          maximumFractionDigits: 2,
        }).format(val)
    }
  }

  const getComparisonColor = () => {
    if (value === null || comparisonValue === null || comparisonValue === undefined) return "text-muted-foreground"
    
    const diff = value - (comparisonValue ? comparisonValue : 0)
    if (Math.abs(diff) < 0.01) return "text-muted-foreground"
    
    // Для некоторых метрик меньше - лучше (например, P/E, долг)
    const lowerIsBetter = ["P/E", "EV/EBITDA", "Debt/Equity"].some(m => label.includes(m))
    
    if (lowerIsBetter) {
      return diff < 0 ? "text-green-600" : "text-red-600"
    }
    return diff > 0 ? "text-green-600" : "text-red-600"
  }

  const getTrendIcon = () => {
    if (!trend || trend === "neutral") return null
    if (trend === "up") {
      return <TrendingUp className="w-4 h-4 text-green-600" />
    }
    return <TrendingDown className="w-4 h-4 text-red-600" />
  }

  return (
    <Card className="p-4 hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-2">
        <div className="flex items-center gap-2">
          <h3 className="text-sm font-medium text-muted-foreground">{label}</h3>
          {description && (
            <TooltipProvider>
              <Tooltip delayDuration={200}>
                <TooltipTrigger>
                  <Info className="w-4 h-4 text-muted-foreground hover:text-foreground transition-colors" />
                </TooltipTrigger>
                <TooltipContent side="right" className="max-w-xs">
                  <div className="space-y-2">
                    <p className="font-semibold">{description.name}</p>
                    {description.formula && (
                      <p className="text-xs">
                        <span className="font-medium">Формула:</span> {description.formula}
                      </p>
                    )}
                    <p className="text-xs">{description.description}</p>
                    <p className="text-xs">
                      <span className="font-medium">Интерпретация:</span>{" "}
                      {description.interpretation}
                    </p>
                    {description.goodValue && (
                      <Badge variant="secondary" className="text-xs">
                        {description.goodValue}
                      </Badge>
                    )}
                  </div>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          )}
        </div>
        {getTrendIcon()}
      </div>

      <div className="flex items-end justify-between">
        <div className="text-2xl font-bold">{formatValue(value)}</div>
      </div>

      {comparisonValue !== undefined && (
        <div className="mt-2 pt-2 border-t">
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground">{comparisonLabel}:</span>
            <span className={getComparisonColor()}>{formatValue(comparisonValue)}</span>
          </div>
        </div>
      )}
    </Card>
  )
}

