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
  PieChart,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  Minus,
  Loader2,
  Info,
} from "lucide-react"
import { aiApi, type BusinessResearch } from "@/lib/api/ai-api"

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

  return (
    <div className="space-y-6">
      {/* Профиль компании */}
      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <Building2 className="h-5 w-5" />
            Профиль компании
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Продукты и услуги */}
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

            {/* Рынки */}
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10">
                  <Globe className="h-4 w-4 text-primary" />
                </div>
                <h4 className="text-sm font-medium">Рынки</h4>
              </div>
              <div className="space-y-2">
                {(profile.markets || []).map((m) => (
                  <div key={m.market} className="flex items-start gap-2">
                    <span className="text-sm font-medium">{m.market}</span>
                    <span className="text-sm text-muted-foreground">— {m.role}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Клиенты */}
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10">
                  <Users className="h-4 w-4 text-primary" />
                </div>
                <h4 className="text-sm font-medium">Ключевые клиенты</h4>
              </div>
              <p className="text-sm text-muted-foreground">{profile.key_clients}</p>
            </div>

            {/* Бизнес-модель */}
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

      {/* Источники выручки */}
      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <PieChart className="h-5 w-5" />
            Источники выручки
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {revenue.map((rs) => (
              <div
                key={rs.segment}
                className="p-4 rounded-lg border bg-card"
              >
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium">{rs.segment}</span>
                    {rs.approximate && (
                      <Badge variant="outline" className="text-xs">~</Badge>
                    )}
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="flex items-center gap-1.5">
                      {trendIcon[rs.trend as keyof typeof trendIcon] || trendIcon.stable}
                      <span className="text-xs text-muted-foreground">
                        {trendLabel[rs.trend as keyof typeof trendLabel] || rs.trend}
                      </span>
                    </div>
                    <span className="text-sm font-semibold min-w-[3rem] text-right">
                      {rs.share_pct}%
                    </span>
                  </div>
                </div>
                <div className="w-full bg-muted rounded-full h-2 mb-2">
                  <div
                    className="bg-primary h-2 rounded-full transition-all"
                    style={{ width: `${Math.min(rs.share_pct, 100)}%` }}
                  />
                </div>
                <p className="text-xs text-muted-foreground">{rs.description}</p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

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
