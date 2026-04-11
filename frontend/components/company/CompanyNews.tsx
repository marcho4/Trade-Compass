"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  Newspaper,
  TrendingUp,
  TrendingDown,
  Minus,
  Clock,
  History,
  CalendarClock,
  Link2,
  Info,
  Loader2,
} from "lucide-react"
import { aiApi, type NewsItem, type DependencyNewsItem, type NewsResponse } from "@/lib/api/ai-api"

interface CompanyNewsProps {
  ticker: string
}

const impactConfig = {
  positive: {
    label: "Позитивно",
    className: "bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20",
    icon: <TrendingUp className="h-3.5 w-3.5" />,
  },
  negative: {
    label: "Негативно",
    className: "bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/20",
    icon: <TrendingDown className="h-3.5 w-3.5" />,
  },
  neutral: {
    label: "Нейтрально",
    className: "bg-muted text-muted-foreground border-border",
    icon: <Minus className="h-3.5 w-3.5" />,
  },
} as const

const severityConfig = {
  high: { label: "Высокая", dotClass: "bg-red-500" },
  medium: { label: "Средняя", dotClass: "bg-yellow-500" },
  low: { label: "Низкая", dotClass: "bg-muted-foreground/40" },
} as const

function NewsCard({ item }: { item: NewsItem }) {
  const impact = impactConfig[item.impact_type] ?? impactConfig.neutral
  const severity = severityConfig[item.severity] ?? severityConfig.low

  return (
    <div className="p-4 rounded-lg border bg-card space-y-2.5">
      <div className="flex items-start justify-between gap-3">
        <p className="text-sm leading-relaxed flex-1">{item.news}</p>
        <Badge variant="outline" className={`shrink-0 flex items-center gap-1 ${impact.className}`}>
          {impact.icon}
          {impact.label}
        </Badge>
      </div>
      <div className="flex items-center gap-3 text-xs text-muted-foreground">
        <span className="flex items-center gap-1.5">
          <span className={`inline-block h-1.5 w-1.5 rounded-full ${severity.dotClass}`} />
          {severity.label} значимость
        </span>
        {item.date && (
          <span className="flex items-center gap-1">
            <Clock className="h-3 w-3" />
            {item.date}
          </span>
        )}
        {item.source && (
          <span className="flex items-center gap-1">
            <Link2 className="h-3 w-3" />
            {item.source}
          </span>
        )}
      </div>
    </div>
  )
}

function DependencyNewsCard({ item }: { item: DependencyNewsItem }) {
  const impact = impactConfig[item.impact_type] ?? impactConfig.neutral
  const severity = severityConfig[item.severity] ?? severityConfig.low

  return (
    <div className="p-4 rounded-lg border bg-card space-y-2.5">
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 space-y-1">
          <span className="text-xs font-medium text-primary">{item.dependency}</span>
          <p className="text-sm leading-relaxed">{item.news}</p>
        </div>
        <Badge variant="outline" className={`shrink-0 flex items-center gap-1 ${impact.className}`}>
          {impact.icon}
          {impact.label}
        </Badge>
      </div>
      <div className="flex items-center gap-3 text-xs text-muted-foreground">
        <span className="flex items-center gap-1.5">
          <span className={`inline-block h-1.5 w-1.5 rounded-full ${severity.dotClass}`} />
          {severity.label} значимость
        </span>
        {item.date && (
          <span className="flex items-center gap-1">
            <Clock className="h-3 w-3" />
            {item.date}
          </span>
        )}
        {item.source && (
          <span className="flex items-center gap-1">
            <Link2 className="h-3 w-3" />
            {item.source}
          </span>
        )}
      </div>
    </div>
  )
}

function NewsSection({
  title,
  icon,
  items,
  renderItem,
}: {
  title: string
  icon: React.ReactNode
  items: (NewsItem | DependencyNewsItem)[]
  renderItem: (item: NewsItem | DependencyNewsItem, i: number) => React.ReactNode
}) {
  if (items.length === 0) return null

  return (
    <Card>
      <CardHeader className="pb-4">
        <CardTitle className="flex items-center gap-2 text-base">
          {icon}
          {title}
          <span className="ml-auto text-xs font-normal text-muted-foreground">{items.length}</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {items.map((item, i) => renderItem(item, i))}
      </CardContent>
    </Card>
  )
}

export const CompanyNews = ({ ticker }: CompanyNewsProps) => {
  const [data, setData] = useState<NewsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()

    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const result = await aiApi.getNews(ticker, controller.signal)
        setData(result)
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return
        setError("Не удалось загрузить новости")
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
            {error || "Новости по компании пока не доступны"}
          </p>
        </CardContent>
      </Card>
    )
  }

  const allEmpty =
    !data.latest_news?.length &&
    !data.upcoming_company_events?.length &&
    !data.upcoming_dependency_events?.length &&
    !data.past_dependency_events?.length &&
    !data.historical_events?.length

  if (allEmpty) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12 gap-2">
          <Newspaper className="h-8 w-8 text-muted-foreground/50" />
          <p className="text-sm text-muted-foreground">Новостей пока нет</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      <NewsSection
        title="Последние новости"
        icon={<Newspaper className="h-4 w-4" />}
        items={data.latest_news ?? []}
        renderItem={(item, i) => <NewsCard key={i} item={item as NewsItem} />}
      />

      <NewsSection
        title="Предстоящие события компании"
        icon={<CalendarClock className="h-4 w-4" />}
        items={data.upcoming_company_events ?? []}
        renderItem={(item, i) => <NewsCard key={i} item={item as NewsItem} />}
      />

      <NewsSection
        title="Предстоящие события зависимостей"
        icon={<CalendarClock className="h-4 w-4 text-muted-foreground" />}
        items={data.upcoming_dependency_events ?? []}
        renderItem={(item, i) => <DependencyNewsCard key={i} item={item as DependencyNewsItem} />}
      />

      <NewsSection
        title="Прошедшие события зависимостей"
        icon={<History className="h-4 w-4 text-muted-foreground" />}
        items={data.past_dependency_events ?? []}
        renderItem={(item, i) => <DependencyNewsCard key={i} item={item as DependencyNewsItem} />}
      />

      <NewsSection
        title="Исторические события"
        icon={<History className="h-4 w-4" />}
        items={data.historical_events ?? []}
        renderItem={(item, i) => <NewsCard key={i} item={item as NewsItem} />}
      />
    </div>
  )
}
