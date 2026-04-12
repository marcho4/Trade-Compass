"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  Building2,
  ShoppingBag,
  Globe,
  Users,
  Briefcase,
  PieChart as PieChartIcon,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  Minus,
  Loader2,
  Info,
} from "lucide-react"
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } from "recharts"
import { aiApi, type BusinessResearch } from "@/lib/api/ai-api"

const CHART_COLORS = [
  "#ef4444", "#f97316", "#eab308", "#22c55e", "#14b8a6",
  "#06b6d4", "#3b82f6", "#6366f1", "#a855f7", "#ec4899",
  "#f43f5e", "#84cc16", "#10b981", "#0ea5e9", "#8b5cf6",
  "#d946ef",
]

function getColorByHash(str: string, index: number, set: Set<number>): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }

  const finalIdx = Math.abs(hash + index) % CHART_COLORS.length
  if (set.has(finalIdx)) {
    return getColorByHash(str, index + 1, set)
  }

  set.add(finalIdx)
  return CHART_COLORS[Math.abs(hash + index) % CHART_COLORS.length]
}

interface CompanyAboutProps {
  ticker: string
}

const severityConfig = {
  critical: { label: "Критический", className: "bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/20" },
  high: { label: "Высокий", className: "bg-orange-500/10 text-orange-600 dark:text-orange-400 border-orange-500/20" },
  moderate: { label: "Умеренный", className: "bg-yellow-500/10 text-yellow-600 dark:text-yellow-400 border-yellow-500/20" },
} as const

const trendIcon = {
  growing: <TrendingUp className="h-4 w-4 text-green-500" />,
  stable: <Minus className="h-4 w-4 text-muted-foreground" />,
  declining: <TrendingDown className="h-4 w-4 text-red-500" />,
}

const trendLabel = {
  growing: "Растёт",
  stable: "Стабильно",
  declining: "Снижается",
}

const depTypeLabel: Record<string, string> = {
  commodity: "Сырьё",
  currency: "Валюта",
  regulation: "Регулирование",
  macro: "Макро",
  technology: "Технологии",
  geopolitics: "Геополитика",
  infrastructure: "Инфраструктура",
  demand: "Спрос",
}

export const CompanyAbout = ({ ticker }: CompanyAboutProps) => {
  const [data, setData] = useState<BusinessResearch | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()

    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const result = await aiApi.getBusinessResearch(ticker, controller.signal)
        setData(result)
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return
        setError("Не удалось загрузить информацию о компании")
      } finally {
        setLoading(false)
      }
    }

    fetchData()
    return () => controller.abort()
  }, [ticker])

  if (loading) {
    return (
      <Card>
        <CardContent className="flex items-center justify-center py-12">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          <span className="ml-2 text-sm text-muted-foreground">Загрузка...</span>
        </CardContent>
      </Card>
    )
  }

  if (error || !data) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12 gap-2">
          <Info className="h-8 w-8 text-muted-foreground/50" />
          <p className="text-sm text-muted-foreground">
            {error || "Информация о компании пока не доступна"}
          </p>
        </CardContent>
      </Card>
    )
  }

  const profile = data.profile
  const revenue = data.revenue_sources || []
  const dependencies = data.dependencies || []
  const colors = new Set<number>()

  const chartData = revenue.map((rs, i) => ({
    name: rs.segment,
    value: rs.share_pct,
    color: getColorByHash(rs.segment, i, colors),
    trend: rs.trend,
    description: rs.description,
    approximate: rs.approximate,
  }))

  return (
    <div className="space-y-6">
      {/* Профиль компании + Структура выручки */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Профиль компании — 2/3 на десктопе */}
        <Card className="lg:col-span-2">
          <CardHeader className="pb-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Building2 className="h-5 w-5" />
              Профиль компании
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 rounded-md bg-primary/10">
                    <ShoppingBag className="h-4 w-4 text-primary" />
                  </div>
                  <h4 className="text-sm font-medium">Продукты и услуги</h4>
                </div>
                <div className="flex flex-wrap gap-2">
                  {(profile.products_and_services || []).map((product) => (
                    <Badge key={product} variant="secondary">{product}</Badge>
                  ))}
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 rounded-md bg-primary/10">
                    <Globe className="h-4 w-4 text-primary" />
                  </div>
                  <h4 className="text-sm font-medium">Рынки</h4>
                </div>
                <div className="space-y-2">
                  {(profile.markets || []).map((m) => (
                    <div key={m.market} className="grid grid-cols-[8rem_1fr] gap-2">
                      <span className="text-sm font-medium">{m.market}</span>
                      <span className="text-sm text-muted-foreground">— {m.role}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 rounded-md bg-primary/10">
                    <Users className="h-4 w-4 text-primary" />
                  </div>
                  <h4 className="text-sm font-medium">Ключевые клиенты</h4>
                </div>
                <p className="text-sm text-muted-foreground">{profile.key_clients}</p>
              </div>

              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <div className="p-1.5 rounded-md bg-primary/10">
                    <Briefcase className="h-4 w-4 text-primary" />
                  </div>
                  <h4 className="text-sm font-medium">Бизнес-модель</h4>
                </div>
                <p className="text-sm text-muted-foreground">{profile.business_model}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Структура выручки — 1/3 на десктопе */}
        <Card className="lg:col-span-1">
          <CardHeader className="pb-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <PieChartIcon className="h-5 w-5" />
              Структура выручки
            </CardTitle>
          </CardHeader>
          <CardContent>
            {revenue.length > 0 ? (
              <div className="flex flex-col items-center">
                <ResponsiveContainer width="100%" height={200}>
                  <PieChart>
                    <Pie
                      data={chartData}
                      cx="50%"
                      cy="50%"
                      innerRadius={55}
                      outerRadius={85}
                      paddingAngle={2}
                      dataKey="value"
                      stroke="none"
                    >
                      {chartData.map((entry, index) => (
                        <Cell key={index} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip
                      content={({ active, payload }) => {
                        if (!active || !payload?.length) return null
                        const d = payload[0].payload
                        return (
                          <div className="rounded-lg border bg-popover px-3 py-2 text-sm shadow-md">
                            <p className="font-medium">{d.name}</p>
                            <p className="text-muted-foreground">{d.value}%</p>
                          </div>
                        )
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>
                <div className="w-full mt-4 space-y-2">
                  {chartData.map((entry) => (
                    <div key={entry.name} className="flex items-center justify-between text-sm">
                      <div className="flex items-center gap-2 min-w-0">
                        <span
                          className="h-3 w-3 rounded-full shrink-0"
                          style={{ backgroundColor: entry.color }}
                        />
                        <span>{entry.name}</span>
                      </div>
                      <div className="flex items-center gap-2 shrink-0 ml-2">
                        {trendIcon[entry.trend as keyof typeof trendIcon] || trendIcon.stable}
                        <span className="font-medium">{entry.value}%</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground text-center py-8">
                Нет данных о выручке
              </p>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Ключевые зависимости */}
      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <AlertTriangle className="h-5 w-5" />
            Ключевые зависимости
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {dependencies.map((dep) => {
              const severity = severityConfig[dep.severity as keyof typeof severityConfig] || severityConfig.moderate
              return (
                <div
                  key={dep.factor}
                  className="p-4 rounded-lg border bg-card space-y-2"
                >
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">{dep.factor}</span>
                    <Badge variant="outline" className={severity.className}>
                      {severity.label}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary" className="text-xs">
                      {depTypeLabel[dep.type] || dep.type}
                    </Badge>
                  </div>
                  <p className="text-xs text-muted-foreground">{dep.description}</p>
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
