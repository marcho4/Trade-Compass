"use client"

import { useRawDataHistory } from "@/hooks/use-raw-data-history"
import { MetricsLineChart, MetricsLineConfig } from "./MetricsLineChart"

interface CompanyMetricsProps {
  ticker: string
}

const REVENUE_NET_PROFIT_LINES: MetricsLineConfig[] = [
  { key: "revenue", label: "Выручка", color: "hsl(221, 83%, 53%)" },
  { key: "netProfit", label: "Чистая прибыль", color: "hsl(142, 71%, 45%)" },
]

export const CompanyMetrics = ({ ticker }: CompanyMetricsProps) => {
  const { data, loading, error } = useRawDataHistory(ticker)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Загрузка показателей...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-destructive">{error}</p>
      </div>
    )
  }

  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Нет данных по показателям</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <MetricsLineChart
        title="Выручка и чистая прибыль"
        description="Динамика ключевых финансовых показателей по периодам"
        data={data}
        lines={REVENUE_NET_PROFIT_LINES}
      />
    </div>
  )
}
