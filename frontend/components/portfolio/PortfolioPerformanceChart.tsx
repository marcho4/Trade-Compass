"use client"

import { useState } from "react"
import { TrendingUp } from "lucide-react"
import { CartesianGrid, Line, LineChart, XAxis, YAxis } from "recharts"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"

type TimeRange = "1M" | "3M" | "6M" | "1Y" | "ALL"

interface PerformanceDataPoint {
  date: string
  portfolio: number
  moex?: number
  sp500?: number
  nasdaq?: number
}

interface PortfolioPerformanceChartProps {
  data: PerformanceDataPoint[]
  portfolioName: string
}

const chartConfig = {
  portfolio: {
    label: "Портфель",
    color: "#8b5cf6",
  },
  moex: {
    label: "MOEX",
    color: "#06b6d4",
  },
  sp500: {
    label: "S&P 500",
    color: "#10b981",
  },
  nasdaq: {
    label: "NASDAQ",
    color: "#f59e0b",
  },
} satisfies ChartConfig

export const PortfolioPerformanceChart = ({
  data,
  portfolioName,
}: PortfolioPerformanceChartProps) => {
  const [timeRange, setTimeRange] = useState<TimeRange>("6M")
  const [showIndices, setShowIndices] = useState({
    moex: true,
    sp500: false,
    nasdaq: false,
  })

  // Фильтруем данные по выбранному периоду
  const getFilteredData = () => {
    const now = new Date()
    const startDate = new Date()

    switch (timeRange) {
      case "1M":
        startDate.setMonth(now.getMonth() - 1)
        break
      case "3M":
        startDate.setMonth(now.getMonth() - 3)
        break
      case "6M":
        startDate.setMonth(now.getMonth() - 6)
        break
      case "1Y":
        startDate.setFullYear(now.getFullYear() - 1)
        break
      case "ALL":
        return data
    }

    return data.filter((point) => new Date(point.date) >= startDate)
  }

  const filteredData = getFilteredData()

  // Вычисляем процентное изменение
  const calculateChange = (dataKey: keyof PerformanceDataPoint) => {
    if (filteredData.length < 2) return 0
    const firstValue = filteredData[0][dataKey] as number
    const lastValue = filteredData[filteredData.length - 1][dataKey] as number
    return ((lastValue - firstValue) / firstValue) * 100
  }

  const portfolioChange = calculateChange("portfolio")

  // Отладка: проверяем данные
  if (filteredData.length > 0) {
    console.log("Chart data sample:", filteredData[0], "Total points:", filteredData.length)
  }

  const toggleIndex = (index: keyof typeof showIndices) => {
    setShowIndices((prev) => ({
      ...prev,
      [index]: !prev[index],
    }))
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          <div>
            <CardTitle>Динамика портфеля</CardTitle>
            <CardDescription>
              Сравнение доходности {portfolioName} с рыночными индексами
            </CardDescription>
          </div>
          <Tabs value={timeRange} onValueChange={(val) => setTimeRange(val as TimeRange)}>
            <TabsList>
              <TabsTrigger value="1M">1М</TabsTrigger>
              <TabsTrigger value="3M">3М</TabsTrigger>
              <TabsTrigger value="6M">6М</TabsTrigger>
              <TabsTrigger value="1Y">1Г</TabsTrigger>
              <TabsTrigger value="ALL">Все</TabsTrigger>
            </TabsList>
          </Tabs>
        </div>
      </CardHeader>
      <CardContent>
        <div className="mb-4 flex flex-wrap gap-3">
          <button
            onClick={() => toggleIndex("moex")}
            className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
              showIndices.moex
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            }`}
          >
            MOEX
          </button>
          <button
            onClick={() => toggleIndex("sp500")}
            className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
              showIndices.sp500
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            }`}
          >
            S&P 500
          </button>
          <button
            onClick={() => toggleIndex("nasdaq")}
            className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
              showIndices.nasdaq
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            }`}
          >
            NASDAQ
          </button>
        </div>

        <ChartContainer config={chartConfig} className="aspect-auto h-[400px] w-full">
          <LineChart
            data={filteredData}
            margin={{
              left: 12,
              right: 12,
              top: 12,
              bottom: 12,
            }}
          >
            <CartesianGrid vertical={false} strokeDasharray="3 3" />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => {
                const date = new Date(value)
                return date.toLocaleDateString("ru-RU", {
                  month: "short",
                  day: "numeric",
                })
              }}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => `${value.toFixed(0)}%`}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  labelFormatter={(value) => {
                    const date = new Date(value)
                    return date.toLocaleDateString("ru-RU", {
                      day: "numeric",
                      month: "long",
                      year: "numeric",
                    })
                  }}
                  formatter={(value) => `${(value as number).toFixed(2)}%`}
                />
              }
            />
            <Line
              dataKey="portfolio"
              type="monotone"
              stroke="#8b5cf6"
              strokeWidth={3}
              dot={false}
              isAnimationActive={false}
            />
            {showIndices.moex && (
              <Line
                dataKey="moex"
                type="monotone"
                stroke="#06b6d4"
                strokeWidth={2}
                dot={false}
                strokeDasharray="5 5"
                isAnimationActive={false}
              />
            )}
            {showIndices.sp500 && (
              <Line
                dataKey="sp500"
                type="monotone"
                stroke="#10b981"
                strokeWidth={2}
                dot={false}
                strokeDasharray="5 5"
                isAnimationActive={false}
              />
            )}
            {showIndices.nasdaq && (
              <Line
                dataKey="nasdaq"
                type="monotone"
                stroke="#f59e0b"
                strokeWidth={2}
                dot={false}
                strokeDasharray="5 5"
                isAnimationActive={false}
              />
            )}
          </LineChart>
        </ChartContainer>
      </CardContent>
      <CardFooter>
        <div className="flex w-full items-start gap-2 text-sm">
          <div className="grid gap-2">
            <div className="flex items-center gap-2 leading-none font-medium">
              {portfolioChange >= 0 ? (
                <>
                  Доходность +{portfolioChange.toFixed(2)}% за период{" "}
                  <TrendingUp className="h-4 w-4 text-green-600" />
                </>
              ) : (
                <>
                  Доходность {portfolioChange.toFixed(2)}% за период{" "}
                  <TrendingUp className="h-4 w-4 text-red-600 rotate-180" />
                </>
              )}
            </div>
            <div className="text-muted-foreground flex items-center gap-2 leading-none">
              Сплошная линия — ваш портфель, пунктир — индексы
            </div>
          </div>
        </div>
      </CardFooter>
    </Card>
  )
}

