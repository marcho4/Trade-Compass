"use client"

import { useMemo } from "react"
import { useRawDataHistory } from "@/hooks/use-raw-data-history"
import { MetricsLineChart, MetricsLineConfig } from "./MetricsLineChart"
import { buildAnnualSnapshots } from "@/lib/build-annual-snapshots"

interface CompanyMetricsProps {
  ticker: string
}

const REVENUE_NET_PROFIT_LINES: MetricsLineConfig[] = [
  { key: "revenue", label: "Выручка", color: "hsl(221, 83%, 53%)" },
  { key: "netProfit", label: "Чистая прибыль", color: "hsl(142, 71%, 45%)" },
]

export const CompanyMetrics = ({ ticker }: CompanyMetricsProps) => {
  const { data, loading, error } = useRawDataHistory(ticker)

  const snapshots = useMemo(
    () => (data.length > 0 ? buildAnnualSnapshots(data, ["revenue", "netProfit"]) : []),
    [data],
  )

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

  if (snapshots.length === 0) {
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
        description="Годовые данные и TTM (trailing twelve months)"
        data={snapshots}
        lines={REVENUE_NET_PROFIT_LINES}
      />
    </div>
  )
}
