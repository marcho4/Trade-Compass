"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ChartContainer } from "@/components/ui/chart"
import {
  Radar,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
} from "recharts"
import { curveCatmullRomClosed } from "d3-shape"
import { TrendingUp, TrendingDown } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import { formatLargeNumber } from "@/lib/utils"
import type { CompanyRating } from "./types"

interface CompanyCardProps {
  id: number
  ticker: string
  name: string
  sector: string
  price: number
  priceChange: number
  priceChangePercent: number
  priceLoading?: boolean
  rating: CompanyRating
  marketCap?: number
  onClick?: () => void
}

const RoundedRadar = (props: any) => {
  const { points, stroke, fill, fillOpacity } = props
  if (!points || points.length === 0) return null

  const pathD: string[] = []
  const context = {
    moveTo(x: number, y: number) { pathD.push(`M${x},${y}`) },
    lineTo(x: number, y: number) { pathD.push(`L${x},${y}`) },
    closePath() { pathD.push("Z") },
    bezierCurveTo(cp1x: number, cp1y: number, cp2x: number, cp2y: number, x: number, y: number) {
      pathD.push(`C${cp1x},${cp1y},${cp2x},${cp2y},${x},${y}`)
    },
  }

  const curve = curveCatmullRomClosed(context as any)
  curve.lineStart()
  for (const p of points) {
    curve.point(p.x, p.y)
  }
  curve.lineEnd()

  return <path d={pathD.join("")} stroke={stroke} fill={fill} fillOpacity={fillOpacity} />
}

function getChartFillColor(rating: number, el: HTMLElement | null): string {
  if (!el) return "oklch(0.6 0.15 145)"
  const n = Math.max(0, Math.min(100, rating))
  const style = getComputedStyle(el)
  const v = (name: string) => parseFloat(style.getPropertyValue(name))

  const lowL = v("--chart-rating-low-l"), lowC = v("--chart-rating-low-c"), lowH = v("--chart-rating-low-h")
  const midL = v("--chart-rating-mid-l"), midC = v("--chart-rating-mid-c"), midH = v("--chart-rating-mid-h")
  const highL = v("--chart-rating-high-l"), highC = v("--chart-rating-high-c"), highH = v("--chart-rating-high-h")

  const lerp = (a: number, b: number, t: number) => a + (b - a) * t

  let l: number, c: number, h: number
  if (n <= 50) {
    const t = n / 50
    l = lerp(lowL, midL, t); c = lerp(lowC, midC, t); h = lerp(lowH, midH, t)
  } else {
    const t = (n - 50) / 50
    l = lerp(midL, highL, t); c = lerp(midC, highC, t); h = lerp(midH, highH, t)
  }

  return `oklch(${l.toFixed(3)} ${c.toFixed(3)} ${h.toFixed(1)})`
}

export const CompanyCard = ({
  ticker,
  name,
  sector,
  price,
  priceChange,
  priceChangePercent,
  priceLoading,
  rating,
  marketCap,
  onClick,
}: CompanyCardProps) => {
  const cardRef = useRef<HTMLDivElement>(null)
  const [chartFillColor, setChartFillColor] = useState("oklch(0.6 0.15 145)")

  const isPositiveChange = priceChange >= 0

  const toPercent = (v: number) => Math.round((v / 6) * 100)

  const radarData = [
    { metric: "Здоровье", value: toPercent(rating.health), fullMark: 100 },
    { metric: "Рост", value: toPercent(rating.growth), fullMark: 100 },
    { metric: "Ров", value: toPercent(rating.moat), fullMark: 100 },
    { metric: "Дивиденды", value: toPercent(rating.dividends), fullMark: 100 },
    { metric: "Оценка", value: toPercent(rating.value), fullMark: 100 },
  ]

  const averageRating = Math.round(
    (toPercent(rating.health) +
      toPercent(rating.growth) +
      toPercent(rating.moat) +
      toPercent(rating.dividends) +
      toPercent(rating.value)) /
      5
  )

  const totalRating = rating.total

  const getRatingColor = (total: number) => {
    if (total >= 5) return "bg-rating-5 text-rating-5-foreground"
    if (total === 4) return "bg-rating-4 text-rating-4-foreground"
    if (total === 3) return "bg-rating-3 text-rating-3-foreground"
    if (total === 2) return "bg-rating-2 text-rating-2-foreground"
    return "bg-rating-1 text-rating-1-foreground"
  }

  useEffect(() => {
    setChartFillColor(getChartFillColor(averageRating, cardRef.current))
  }, [averageRating])

  return (
    <div
      ref={cardRef}
      className="cursor-pointer transition-[transform,box-shadow] hover:shadow-lg hover:scale-[1.02]"
      onClick={onClick}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          onClick?.()
        }
      }}
      role="button"
      tabIndex={0}
      aria-label={`Открыть анализ компании ${name}`}
    >
      <Card className="h-full">
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="space-y-1 flex-1">
              <div className="flex items-center gap-2">
                <CardTitle className="text-2xl font-bold">{ticker}</CardTitle>
                <Badge variant="outline" className="text-xs">
                  {sector}
                </Badge>
              </div>
              <p className="text-sm text-muted-foreground line-clamp-1">
                {name}
              </p>
            </div>
            <div
              className={`ml-2 h-10 w-10 flex items-center justify-center rounded-full text-base font-bold ${getRatingColor(totalRating)}`}
            >
              {totalRating}
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col items-center justify-center">
              <ChartContainer
                config={{
                  value: {
                    label: "Рейтинг",
                    color: "hsl(var(--chart-1))",
                  },
                }}
                className="h-[200px] w-full"
              >
                <RadarChart data={radarData}>
                  <PolarGrid strokeDasharray="3 3" />
                  <PolarAngleAxis
                    dataKey="metric"
                    tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
                  />
                  <PolarRadiusAxis
                    angle={90}
                    domain={[0, 100]}
                    tick={false}
                  />
                  <Radar
                    name="Рейтинг"
                    dataKey="value"
                    stroke={chartFillColor}
                    fill={chartFillColor}
                    fillOpacity={0.6}
                    shape={<RoundedRadar />}
                  />
                </RadarChart>
              </ChartContainer>
            </div>

            <div className="flex flex-col space-y-3">
              <div>
                <p className="text-xs text-muted-foreground">Цена</p>
                {priceLoading ? (
                  <div className="h-6 w-24 rounded bg-muted animate-pulse mt-1" />
                ) : (
                  <>
                    <div className="flex items-baseline gap-2">
                      <p className="text-lg font-bold">
                        {new Intl.NumberFormat("ru-RU", {
                          style: "currency",
                          currency: "RUB",
                          maximumFractionDigits: 2,
                        }).format(price)}
                      </p>
                    </div>
                    <div className="flex items-center gap-1 mt-1">
                      {isPositiveChange ? (
                        <TrendingUp className="h-3 w-3 text-positive" />
                      ) : (
                        <TrendingDown className="h-3 w-3 text-negative" />
                      )}
                      <span
                        className={`text-xs font-semibold ${
                          isPositiveChange ? "text-positive" : "text-negative"
                        }`}
                      >
                        {isPositiveChange ? "+" : ""}
                        {priceChangePercent.toFixed(2)}%
                      </span>
                      <span
                        className={`text-xs ${
                          isPositiveChange ? "text-positive" : "text-negative"
                        }`}
                      >
                        ({isPositiveChange ? "+" : ""}
                        {new Intl.NumberFormat("ru-RU", {
                          style: "currency",
                          currency: "RUB",
                          maximumFractionDigits: 2,
                        }).format(priceChange)})
                      </span>
                    </div>
                  </>
                )}
              </div>

              <div className="space-y-2 text-xs">
                {marketCap ? (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Капитализация:</span>
                    <span className="font-medium">
                      {formatLargeNumber(marketCap)} ₽
                    </span>
                  </div>
                ) : (
                  <div className="h-4 w-32 rounded bg-muted animate-pulse" />
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
