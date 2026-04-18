"use client"

import { useState } from "react"
import { Area, CartesianGrid, ComposedChart, XAxis, YAxis } from "recharts"
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { AnnualSnapshot } from "@/lib/build-annual-snapshots"

export interface MetricsLineConfig {
  key: string
  label: string
  color: string
  hiddenByDefault?: boolean
}

interface MetricsLineChartProps {
  title: string
  description?: string
  data: AnnualSnapshot[]
  lines: MetricsLineConfig[]
  valueFormatter?: (value: number) => string
}

function formatCompactValue(value: number): string {
  const abs = Math.abs(value)
  if (abs >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(1)} млрд`
  if (abs >= 1_000_000) return `${(value / 1_000_000).toFixed(1)} млн`
  if (abs >= 1_000) return `${(value / 1_000).toFixed(1)} тыс`
  return value.toString()
}

export const MetricsLineChart = ({ title, description, data, lines, valueFormatter = formatCompactValue }: MetricsLineChartProps) => {
  const [visibleKeys, setVisibleKeys] = useState<Set<string>>(
    () => new Set(lines.filter((l) => !l.hiddenByDefault).map((l) => l.key))
  )

  const toggleLine = (key: string) => {
    setVisibleKeys((prev) => {
      const next = new Set(prev)
      if (next.has(key)) next.delete(key)
      else next.add(key)
      return next
    })
  }

  const chartConfig: ChartConfig = Object.fromEntries(
    lines.map((line) => [line.key, { label: line.label, color: line.color }])
  )

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
        <div className="flex flex-wrap gap-2 pt-2">
          {lines.map((line) => {
            const isActive = visibleKeys.has(line.key)
            return (
              <Button
                key={line.key}
                variant="outline"
                size="sm"
                onClick={() => toggleLine(line.key)}
                className="h-7 rounded-full text-xs gap-1.5 transition-colors"
                style={isActive ? {
                  borderColor: line.color,
                  backgroundColor: `color-mix(in oklch, ${line.color} 15%, transparent)`,
                  color: line.color,
                } : {
                  opacity: 0.5,
                }}
              >
                <span
                  className="h-2 w-2 rounded-full shrink-0"
                  style={{ backgroundColor: line.color }}
                />
                {line.label}
              </Button>
            )
          })}
        </div>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="h-[350px] w-full">
          <ComposedChart data={data} accessibilityLayer>
            <defs>
              {lines.map((line) => (
                <linearGradient key={line.key} id={`fill-${line.key}`} x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor={`var(--color-${line.key})`} stopOpacity={0.25} />
                  <stop offset="100%" stopColor={`var(--color-${line.key})`} stopOpacity={0.02} />
                </linearGradient>
              ))}
            </defs>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="period"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              className="text-xs"
            />
            <YAxis
              tickFormatter={valueFormatter}
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
                    const formatted = typeof value === "number" ? valueFormatter(value) : value
                    return (
                      <span>
                        {config?.label ?? name}: <strong>{formatted}</strong>
                      </span>
                    )
                  }}
                />
              }
            />
            {lines.map((line) =>
              visibleKeys.has(line.key) ? (
                <Area
                  key={line.key}
                  type="monotone"
                  dataKey={line.key}
                  stroke={`var(--color-${line.key})`}
                  fill={`url(#fill-${line.key})`}
                  strokeWidth={2}
                  dot={{ r: 4, fill: `var(--color-${line.key})` }}
                  activeDot={{ r: 6 }}
                  connectNulls
                />
              ) : null,
            )}
          </ComposedChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
