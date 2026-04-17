"use client"

import { useMemo } from "react"
import { useRawDataHistory } from "@/hooks/use-raw-data-history"
import { MetricsLineChart, MetricsLineConfig } from "./MetricsLineChart"
import { RawDataTable } from "./RawDataTable"
import { buildAnnualSnapshots } from "@/lib/build-annual-snapshots"

interface CompanyMetricsProps {
  ticker: string
}

const REVENUE_NET_PROFIT_LINES: MetricsLineConfig[] = [
  { key: "revenue", label: "Выручка", color: "hsl(221, 83%, 53%)" },
  { key: "netProfit", label: "Чистая прибыль", color: "hsl(142, 71%, 45%)" },
  { key: "operatingCashFlow", label: "Операционный Денежный поток", color: "hsl(32, 95%, 55%)", hiddenByDefault: true },
  { key: "freeCashFlow", label: "Свободный денежный поток", color: "hsl(280, 70%, 55%)", hiddenByDefault: true },
]

const DEBT_EQUITY_LINES: MetricsLineConfig[] = [
  { key: "debt", label: "Долг", color: "hsl(0, 72%, 51%)" },
  { key: "equity", label: "Собственный капитал", color: "hsl(221, 83%, 53%)" },
]

const MARGIN_LINES: MetricsLineConfig[] = [
  { key: "grossMargin", label: "Валовая маржа", color: "hsl(221, 83%, 53%)" },
  { key: "operatingMargin", label: "Операционная маржа", color: "hsl(32, 95%, 55%)" },
  { key: "netMargin", label: "Чистая маржа", color: "hsl(142, 71%, 45%)" },
]

const BALANCE_LINES: MetricsLineConfig[] = [
  { key: "totalAssets", label: "Всего активов", color: "hsl(221, 83%, 53%)" },
  { key: "totalLiabilities", label: "Всего обязательств", color: "hsl(0, 72%, 51%)" },
]

export const CompanyMetrics = ({ ticker }: CompanyMetricsProps) => {
  const { data, loading, error } = useRawDataHistory(ticker)

  const snapshots = useMemo(
    () => (data.length > 0 ? buildAnnualSnapshots(data, ["revenue", "netProfit", "operatingCashFlow", "freeCashFlow"]) : []),
    [data],
  )

  const debtEquitySnapshots = useMemo(
    () => (data.length > 0 ? buildAnnualSnapshots(data, ["debt", "equity"]) : []),
    [data],
  )

  const balanceSnapshots = useMemo(
    () => (data.length > 0 ? buildAnnualSnapshots(data, ["totalAssets", "totalLiabilities", "equity"]) : []),
    [data],
  )

  const marginSnapshots = useMemo(() => {
    if (data.length === 0) return []
    const raw = buildAnnualSnapshots(data, ["revenue", "grossProfit", "ebit", "netProfit"])
    return raw.map((s) => {
      const revenue = s.revenue as number
      if (!revenue) return { period: s.period, grossMargin: null, operatingMargin: null, netMargin: null }
      return {
        period: s.period,
        grossMargin: +((s.grossProfit as number) / revenue * 100).toFixed(1),
        operatingMargin: +((s.ebit as number) / revenue * 100).toFixed(1),
        netMargin: +((s.netProfit as number) / revenue * 100).toFixed(1),
      }
    })
  }, [data])

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
      <MetricsLineChart
        title="Маржинальность"
        description="Годовые данные и TTM (trailing twelve months)"
        data={marginSnapshots}
        lines={MARGIN_LINES}
        valueFormatter={(v) => `${v}%`}
      />
      <MetricsLineChart
        title="Долг и собственный капитал"
        description="Годовые данные и TTM (trailing twelve months)"
        data={debtEquitySnapshots}
        lines={DEBT_EQUITY_LINES}
      />
      <MetricsLineChart
        title="Баланс компании"
        description="Годовые данные и TTM (trailing twelve months)"
        data={balanceSnapshots}
        lines={BALANCE_LINES}
      />
      <RawDataTable data={data} />
    </div>
  )
}
