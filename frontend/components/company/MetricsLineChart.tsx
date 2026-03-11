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
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import { AnnualSnapshot } from "@/lib/build-annual-snapshots"

export interface MetricsLineConfig {
  key: string
  label: string
  color: string
}

interface MetricsLineChartProps {
  title: string
  description?: string
  data: AnnualSnapshot[]
  lines: MetricsLineConfig[]
}

function formatCompactValue(value: number): string {
  const abs = Math.abs(value)
  if (abs >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(1)} млрд`
  if (abs >= 1_000_000) return `${(value / 1_000_000).toFixed(1)} млн`
  if (abs >= 1_000) return `${(value / 1_000).toFixed(1)} тыс`
  return value.toString()
}

export const MetricsLineChart = ({ title, description, data, lines }: MetricsLineChartProps) => {
  const [visibleKeys, setVisibleKeys] = useState<Set<string>>(
    () => new Set(lines.map((l) => l.key))
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
        <div className="flex flex-wrap gap-4 pt-2">
          {lines.map((line) => (
            <div key={line.key} className="flex items-center gap-2">
              <Checkbox
                id={`toggle-${line.key}`}
                checked={visibleKeys.has(line.key)}
                onCheckedChange={() => toggleLine(line.key)}
                style={{ borderColor: line.color, color: line.color, backgroundColor: "white" } as React.CSSProperties}
              />
              <Label htmlFor={`toggle-${line.key}`} className="text-sm cursor-pointer">
                {line.label}
              </Label>
            </div>
          ))}
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
