"use client"

import { CartesianGrid, Line, LineChart, XAxis, YAxis } from "recharts"
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { RawData } from "@/types/raw-data"

export interface MetricsLineConfig {
  key: keyof RawData
  label: string
  color: string
}

interface MetricsLineChartProps {
  title: string
  description?: string
  data: RawData[]
  lines: MetricsLineConfig[]
}

const PERIOD_LABELS: Record<string, string> = {
  Q1: "3 мес.",
  Q2: "6 мес.",
  Q3: "9 мес.",
  YEAR: "Год",
}

function formatCompactValue(value: number): string {
  const abs = Math.abs(value)
  if (abs >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(1)} млрд`
  if (abs >= 1_000_000) return `${(value / 1_000_000).toFixed(1)} млн`
  if (abs >= 1_000) return `${(value / 1_000).toFixed(1)} тыс`
  return value.toString()
}

function buildPeriodLabel(item: RawData): string {
  const periodLabel = PERIOD_LABELS[item.period] ?? item.period
  return `${periodLabel} ${item.year}`
}

export const MetricsLineChart = ({ title, description, data, lines }: MetricsLineChartProps) => {
  const sorted = [...data].sort((a, b) => {
    if (a.year !== b.year) return a.year - b.year
    const periodOrder = ["Q1", "Q2", "Q3", "YEAR"]
    return periodOrder.indexOf(a.period) - periodOrder.indexOf(b.period)
  })

  const chartConfig: ChartConfig = Object.fromEntries(
    lines.map((line) => [line.key, { label: line.label, color: line.color }])
  )

  const chartData = sorted.map((item) => ({
    label: buildPeriodLabel(item),
    ...Object.fromEntries(
      lines.map((line) => [line.key, (item[line.key] as number | null | undefined) ?? null])
    ),
  }))

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="h-[350px] w-full">
          <LineChart data={chartData} accessibilityLayer>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="label"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              className="text-xs"
            />
            <YAxis
              tickFormatter={formatCompactValue}
              tickLine={false}
              axisLine={false}
              className="text-xs"
              width={80}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  formatter={(value, name) => {
                    const config = chartConfig[name as string]
                    const formatted = typeof value === "number" ? formatCompactValue(value) : value
                    return (
                      <span>
                        {config?.label ?? name}: <strong>{formatted}</strong>
                      </span>
                    )
                  }}
                />
              }
            />
            <ChartLegend content={<ChartLegendContent />} />
            {lines.map((line) => (
              <Line
                key={String(line.key)}
                type="monotone"
                dataKey={String(line.key)}
                stroke={`var(--color-${String(line.key)})`}
                strokeWidth={2}
                dot={{ r: 4 }}
                activeDot={{ r: 6 }}
                connectNulls
              />
            ))}
          </LineChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
