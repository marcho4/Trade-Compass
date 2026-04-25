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
import { cn } from "@/lib/utils"
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
    <Card className="overflow-hidden rounded-[2px] shadow-[0_1px_0_rgba(20,20,20,0.02)]">
      <CardHeader className="gap-1.5 border-b border-border/60 pb-4">
        <CardTitle className="text-base font-semibold tracking-tight">{title}</CardTitle>
        {description && (
          <CardDescription className="text-xs text-muted-foreground">{description}</CardDescription>
        )}
        <div className="flex flex-wrap gap-1.5 pt-3">
          {lines.map((line) => {
            const isActive = visibleKeys.has(line.key)
            return (
              <button
                key={line.key}
                type="button"
                onClick={() => toggleLine(line.key)}
                className={cn(
                  "inline-flex items-center gap-1.5 rounded-[2px] border px-2.5 py-1 text-xs font-medium transition-all",
                  "hover:bg-accent/60",
                  isActive
                    ? "border-border bg-card text-foreground shadow-xs"
                    : "border-transparent bg-muted/40 text-muted-foreground/70"
                )}
              >
                <span
                  aria-hidden
                  className={cn("h-1.5 w-1.5 rounded-full transition-opacity", !isActive && "opacity-40")}
                  style={{ backgroundColor: line.color }}
                />
                {line.label}
              </button>
            )
          })}
        </div>
      </CardHeader>
      <CardContent className="pt-6 pb-4">
        <ChartContainer config={chartConfig} className="h-[320px] w-full">
          <ComposedChart data={data} accessibilityLayer margin={{ top: 4, right: 8, bottom: 0, left: 0 }}>
            <defs>
              {lines.map((line) => (
                <linearGradient key={line.key} id={`fill-${line.key}`} x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor={`var(--color-${line.key})`} stopOpacity={0.18} />
                  <stop offset="95%" stopColor={`var(--color-${line.key})`} stopOpacity={0} />
                </linearGradient>
              ))}
            </defs>
            <CartesianGrid
              vertical={false}
              stroke="var(--border)"
              strokeOpacity={0.6}
              strokeDasharray="3 4"
            />
            <XAxis
              dataKey="period"
              tickLine={false}
              axisLine={false}
              tickMargin={10}
              className="text-[11px] font-mono"
              stroke="var(--muted-foreground)"
            />
            <YAxis
              tickFormatter={valueFormatter}
              tickLine={false}
              axisLine={false}
              tickMargin={6}
              className="text-[11px] font-mono"
              stroke="var(--muted-foreground)"
              width={72}
            />
            <ChartTooltip
              cursor={{ stroke: "var(--border)", strokeWidth: 1, strokeDasharray: "3 3" }}
              content={
                <ChartTooltipContent
                  className="rounded-lg border-border/60 bg-popover/95 backdrop-blur-sm shadow-md"
                  formatter={(value, name) => {
                    const config = chartConfig[name as string]
                    const formatted = typeof value === "number" ? valueFormatter(value) : value
                    return (
                      <div className="flex w-full items-center justify-between gap-4">
                        <span className="text-muted-foreground">{config?.label ?? name}</span>
                        <span className="font-mono font-medium tabular-nums text-foreground">{formatted}</span>
                      </div>
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
                  strokeWidth={1.75}
                  dot={false}
                  activeDot={{
                    r: 4,
                    strokeWidth: 2,
                    stroke: "var(--background)",
                    fill: `var(--color-${line.key})`,
                  }}
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
