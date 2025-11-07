"use client"

import { Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts"
import { FinancialMetrics } from "@/types"

interface FinancialChartProps {
  title: string
  description?: string
  data: FinancialMetrics[]
  dataKeys: {
    key: keyof FinancialMetrics
    label: string
    color: string
  }[]
  chartType?: "line" | "bar" | "area"
  formatValue?: (value: number) => string
}

export const FinancialChart = ({
  title,
  description,
  data,
  dataKeys,
  chartType = "line",
  formatValue,
}: FinancialChartProps) => {
  const formatXAxis = (value: string) => {
    const item = data.find((d) => d.reportId === value)
    if (!item) return value
    return `${item.reportPeriod} ${item.reportYear}`
  }

  const defaultFormatValue = (value: number) => {
    return new Intl.NumberFormat("ru-RU", {
      notation: "compact",
      compactDisplay: "short",
      maximumFractionDigits: 1,
    }).format(value)
  }

  const valueFormatter = formatValue || defaultFormatValue

  const chartData = data.map((item) => ({
    name: item.reportId,
    ...dataKeys.reduce(
      (acc, dk) => ({
        ...acc,
        [dk.label]: item[dk.key],
      }),
      {}
    ),
  }))

  const renderChart = () => {
    const commonProps = {
      data: chartData,
      margin: { top: 10, right: 30, left: 0, bottom: 0 },
    }

    const axes = (
      <>
        <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
        <XAxis
          dataKey="name"
          tickFormatter={formatXAxis}
          className="text-xs"
          stroke="hsl(var(--muted-foreground))"
        />
        <YAxis
          tickFormatter={valueFormatter}
          className="text-xs"
          stroke="hsl(var(--muted-foreground))"
        />
        <Tooltip
          formatter={valueFormatter}
          contentStyle={{
            backgroundColor: "hsl(var(--background))",
            border: "1px solid hsl(var(--border))",
            borderRadius: "8px",
          }}
        />
        <Legend />
      </>
    )

    switch (chartType) {
      case "bar":
        return (
          <BarChart {...commonProps}>
            {axes}
            {dataKeys.map((dk) => (
              <Bar key={dk.label} dataKey={dk.label} fill={dk.color} radius={[4, 4, 0, 0]} />
            ))}
          </BarChart>
        )
      case "area":
        return (
          <AreaChart {...commonProps}>
            {axes}
            {dataKeys.map((dk) => (
              <Area
                key={dk.label}
                type="monotone"
                dataKey={dk.label}
                stroke={dk.color}
                fill={dk.color}
                fillOpacity={0.2}
              />
            ))}
          </AreaChart>
        )
      default:
        return (
          <LineChart {...commonProps}>
            {axes}
            {dataKeys.map((dk) => (
              <Line
                key={dk.label}
                type="monotone"
                dataKey={dk.label}
                stroke={dk.color}
                strokeWidth={2}
                dot={{ fill: dk.color, r: 4 }}
              />
            ))}
          </LineChart>
        )
    }
  }

  return (
    <Card className="p-6">
      <div className="mb-6">
        <div className="flex items-center gap-3 mb-2">
          <h3 className="text-lg font-semibold">{title}</h3>
          <Badge variant="outline">По кварталам</Badge>
        </div>
        {description && <p className="text-sm text-muted-foreground">{description}</p>}
      </div>
      <ResponsiveContainer width="100%" height={300}>
        {renderChart()}
      </ResponsiveContainer>
    </Card>
  )
}

