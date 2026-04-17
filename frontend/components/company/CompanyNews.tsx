"use client"

import { useEffect, useState } from "react"
import { Newspaper } from "lucide-react"
import { aiApi, type NewsItem, type DependencyNewsItem, type NewsResponse } from "@/lib/api/ai-api"
import Link from "next/link"

interface CompanyNewsProps {
  ticker: string
}

const impactColor = (s: string) =>
  s === "positive" ? "text-positive" : s === "negative" ? "text-negative" : "text-muted-foreground"

const impactBg = (s: string) =>
  s === "positive"
    ? "bg-positive/10"
    : s === "negative"
      ? "bg-negative/10"
      : "bg-muted"

const impactLabel = (s: string) =>
  s === "positive" ? "POS" : s === "negative" ? "NEG" : "NEU"

const impactDotBg = (s: string) =>
  s === "positive" ? "bg-positive" : s === "negative" ? "bg-negative" : "bg-muted-foreground"

function ImpactPill({ impact }: { impact: string }) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 font-mono text-[10px] font-bold tracking-wide uppercase ${impactColor(impact)} ${impactBg(impact)} px-1.5 py-0.5 rounded-sm leading-tight`}
    >
      <span className={`w-1.5 h-1.5 rounded-full ${impactDotBg(impact)}`} />
      {impactLabel(impact)}
    </span>
  )
}

function SeverityDots({ level }: { level: string }) {
  const n = level === "high" ? 3 : level === "medium" ? 2 : 1
  const activeClass =
    level === "high"
      ? "bg-primary"
      : level === "medium"
        ? "bg-foreground/60"
        : "bg-muted-foreground/40"
  const label =
    level === "high"
      ? "Высокая значимость"
      : level === "medium"
        ? "Средняя значимость"
        : "Низкая значимость"
  return (
    <span className="relative inline-flex gap-0.5 items-center group cursor-help">
      {[0, 1, 2].map((i) => (
        <span
          key={i}
          className={`w-1 h-1 rounded-[1px] ${i < n ? activeClass : "bg-border"}`}
        />
      ))}
      <span
        role="tooltip"
        className="pointer-events-none absolute left-full top-1/2 -translate-y-1/2 ml-2 whitespace-nowrap rounded-sm border border-border bg-popover px-1.5 py-1 font-mono text-[10px] uppercase tracking-wide text-popover-foreground opacity-0 shadow-sm transition-opacity duration-150 group-hover:opacity-100 z-[9999]"
      >
        {label}
      </span>
    </span>
  )
}

function PanelHeader({
  title,
  count,
  accentClass,
}: {
  title: string
  count: number
  accentClass: string
}) {
  return (
    <div className="flex items-center justify-between px-3.5 py-2.5 border-b bg-muted/50 sticky top-0 z-[2]">
      <div className="flex items-center gap-2.5">
        <span className={`w-[3px] h-3 rounded-[1px] ${accentClass}`} />
        <span className="font-mono text-[11px] font-bold tracking-[1.2px] uppercase text-foreground">
          {title}
        </span>
        <span className="font-mono text-[10px] text-muted-foreground">
          [{String(count).padStart(3, "0")}]
        </span>
      </div>
    </div>
  )
}

function safeDate(d: string): Date | null {
  if (!d) return null
  const m = d.match(/^(\d{2})\.(\d{2})\.(\d{4})(?:[ T](.+))?$/)
  const input = m ? `${m[3]}-${m[2]}-${m[1]}${m[4] ? "T" + m[4] : ""}` : d
  const date = new Date(input)
  return isNaN(date.getTime()) ? null : date
}

function fmtRelTime(d: string) {
  const date = safeDate(d)
  if (!date) return d
  const now = new Date()
  const sameDay = date.toDateString() === now.toDateString()
  if (sameDay) {
    return date.toLocaleString("ru-RU", { hour: "2-digit", minute: "2-digit" })
  }
  return date.toLocaleString("ru-RU", { day: "2-digit", month: "short" })
}

function fmtShortDate(d: string) {
  const date = safeDate(d)
  if (!date) return d
  return date.toLocaleString("ru-RU", { day: "2-digit", month: "2-digit" })
}

function fmtHistDate(d: string) {
  const date = safeDate(d)
  if (!date) return d
  return date.toLocaleString("ru-RU", { day: "2-digit", month: "2-digit", year: "2-digit" })
}

function fmtEventDay(d: string) {
  const date = safeDate(d)
  if (!date) return { day: "—", month: "" }
  return {
    day: String(date.getDate()),
    month: date.toLocaleString("ru-RU", { month: "short" }),
  }
}

function dayOffset(d: string): number {
  const date = safeDate(d)
  if (!date) return 0
  return Math.ceil((date.getTime() - Date.now()) / (1000 * 60 * 60 * 24))
}

function NewsRow({ n }: { n: NewsItem }) {
  return (
    <div className="grid grid-cols-[56px_20px_1fr_48px] gap-3 items-baseline px-3.5 py-2.5 border-b border-border/60 text-xs leading-snug cursor-pointer hover:bg-muted/30 transition-colors">
      <span className="font-mono text-[10px] text-muted-foreground">
        {fmtRelTime(n.date)}
      </span>
      <SeverityDots level={n.severity} />
      <div>
        <div className="text-foreground mb-0.5">{n.news}</div>
        <div className="font-mono text-[10px] text-muted-foreground tracking-wide">
          <Link href={n.source}>{n.source.toUpperCase()}</Link>
        </div>
      </div>
      <div className="flex justify-end">
        <ImpactPill impact={n.impact_type} />
      </div>
    </div>
  )
}

function UpcomingCompanyRow({ e }: { e: NewsItem }) {
  const { day, month } = fmtEventDay(e.date)
  const offset = dayOffset(e.date)

  const typeLabel =
    e.source === "earnings" || e.news.toLowerCase().includes("отчёт")
      ? "ER"
      : e.source === "dividend" || e.news.toLowerCase().includes("дивиденд")
        ? "DV"
        : e.source === "conference" || e.news.toLowerCase().includes("investor")
          ? "CF"
          : "EV"

  return (
    <div className="grid grid-cols-[56px_28px_1fr_48px] gap-2.5 items-center px-3.5 py-3 border-b border-border/60 text-xs">
      <div className="text-center">
        <div className="font-mono text-base font-bold text-foreground leading-none">
          {day}
        </div>
        <div className="font-mono text-[9px] text-muted-foreground uppercase mt-0.5">
          {month}
        </div>
      </div>
      <span className="inline-flex items-center justify-center font-mono text-[10px] font-bold tracking-wide uppercase text-primary bg-primary/10 border border-primary/20 px-1 py-0.5 rounded-sm leading-tight">
        {typeLabel}
      </span>
      <div>
        <div className="text-foreground text-[12.5px] font-semibold">{e.news}</div>
        <div className="font-mono text-[10px] text-muted-foreground mt-0.5">
          {e.date} · <Link href={e.source}>{e.source}</Link>
        </div>
      </div>
      <div className="flex justify-end">
        <ImpactPill impact={e.impact_type} />
      </div>
    </div>
  )
}

function UpcomingDepsRow({ e }: { e: DependencyNewsItem }) {
  const { day, month } = fmtEventDay(e.date)
  const offset = dayOffset(e.date)

  return (
    <div className="grid grid-cols-[48px_1fr_48px] gap-2.5 items-center px-3.5 py-2.5 border-b border-border/60 text-xs">
      <div className="text-center">
        <div className="font-mono text-sm font-bold text-foreground leading-none">
          {day}
        </div>
        <div className="font-mono text-[8px] text-muted-foreground uppercase mt-0.5">
          {month}
        </div>
      </div>
      <div>
        <div className="text-foreground text-xs">
          <span className="font-mono text-blue-600 dark:text-blue-400 mr-1.5 font-bold">
            {e.dependency}
          </span>
          {e.news}
        </div>
        <div className="font-mono text-[9px] text-muted-foreground mt-0.5">
          <Link href={e.source}>{e.source.toUpperCase()}</Link>
        </div>
      </div>
      <div className="flex justify-end">
        <ImpactPill impact={e.impact_type} />
      </div>
    </div>
  )
}

function PastDepsRow({ e }: { e: DependencyNewsItem }) {
  return (
    <div className="grid grid-cols-[52px_1fr_72px] gap-3 items-center px-3.5 py-2.5 border-b border-border/60 text-xs">
      <div className="font-mono text-[10px] text-muted-foreground">
        {fmtHistDate(e.date)}
      </div>
      <div>
        <div className="text-foreground text-xs">
          <span className="font-mono text-blue-600 dark:text-blue-400 mr-1.5 font-bold">
            {e.dependency}
          </span>
          {e.news}
        </div>
        <div className="font-mono text-[9px] text-muted-foreground mt-0.5">
          <Link href={e.source}>{e.source.toUpperCase()}</Link>
        </div>
      </div>
      <div className="flex justify-end">
        <ImpactPill impact={e.impact_type} />
      </div>
    </div>
  )
}

function HistoricRow({ e }: { e: NewsItem }) {
  return (
    <div className="grid grid-cols-[76px_1fr_72px] gap-2.5 items-center px-3.5 py-2.5 border-b border-border/60 text-xs">
      <div className="font-mono text-[10px] text-muted-foreground">
        {fmtHistDate(e.date)}
      </div>
      <div>
        <div className="text-foreground text-xs">{e.news}</div>
        <div className="font-mono text-[9px] text-muted-foreground mt-0.5">
          <Link href={e.source}>{e.source.toUpperCase()}</Link>
        </div>
      </div>
      <div className="flex justify-end">
        <ImpactPill impact={e.impact_type} />
      </div>
    </div>
  )
}

function TermPanel({
  title,
  count,
  accentClass,
  children,
  className = "",
}: {
  title: string
  count: number
  accentClass: string
  children: React.ReactNode
  className?: string
}) {
  return (
    <div
      className={`min-w-0 bg-card border rounded-sm flex flex-col shadow-[0_1px_0_rgba(0,0,0,0.02)] ${className}`}
    >
      <PanelHeader title={title} count={count} accentClass={accentClass} />
      <div className="overflow-y-auto max-h-[360px]">{children}</div>
    </div>
  )
}

function LoadingSkeleton() {
  return (
    <div className="space-y-3.5">
      <div className="flex flex-col md:flex-row gap-3.5">
        <div className="flex-[1.6] bg-card border rounded-sm overflow-hidden">
          <div className="px-3.5 py-2.5 border-b bg-muted/50">
            <div className="h-3 w-32 rounded bg-muted animate-pulse" />
          </div>
          {[...Array(4)].map((_, i) => (
            <div key={i} className="px-3.5 py-2.5 border-b border-border/60">
              <div className="h-3 w-full rounded bg-muted animate-pulse mb-1.5" />
              <div className="h-2.5 w-24 rounded bg-muted animate-pulse" />
            </div>
          ))}
        </div>
        <div className="flex-1 bg-card border rounded-sm overflow-hidden">
          <div className="px-3.5 py-2.5 border-b bg-muted/50">
            <div className="h-3 w-40 rounded bg-muted animate-pulse" />
          </div>
          {[...Array(3)].map((_, i) => (
            <div key={i} className="px-3.5 py-3 border-b border-border/60">
              <div className="h-3.5 w-full rounded bg-muted animate-pulse mb-1" />
              <div className="h-2.5 w-32 rounded bg-muted animate-pulse" />
            </div>
          ))}
        </div>
      </div>
      <div className="flex flex-col md:flex-row gap-3.5">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="flex-1 bg-card border rounded-sm overflow-hidden">
            <div className="px-3.5 py-2.5 border-b bg-muted/50">
              <div className="h-3 w-36 rounded bg-muted animate-pulse" />
            </div>
            {[...Array(3)].map((_, j) => (
              <div key={j} className="px-3.5 py-2.5 border-b border-border/60">
                <div className="h-3 w-full rounded bg-muted animate-pulse mb-1" />
                <div className="h-2.5 w-20 rounded bg-muted animate-pulse" />
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  )
}

export const CompanyNews = ({ ticker }: CompanyNewsProps) => {
  const [data, setData] = useState<NewsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [timedOut, setTimedOut] = useState(false)

  useEffect(() => {
    const controller = new AbortController()
    let pollInterval: ReturnType<typeof setInterval> | null = null
    let pollCount = 0
    const MAX_POLLS = 12

    const fetchNews = async (): Promise<NewsResponse | null> => {
      const result = await aiApi.getNews(ticker, controller.signal)
      return result
    }

    const startPolling = () => {
      pollInterval = setInterval(async () => {
        pollCount++
        try {
          const result = await fetchNews()
          if (result) {
            if (pollInterval) clearInterval(pollInterval)
            setData(result)
            setLoading(false)
            return
          }
        } catch {
          // aborted or network error — stop polling
          if (pollInterval) clearInterval(pollInterval)
          setLoading(false)
          return
        }

        if (pollCount >= MAX_POLLS) {
          if (pollInterval) clearInterval(pollInterval)
          setTimedOut(true)
          setLoading(false)
        }
      }, 5000)
    }

    const init = async () => {
      try {
        setLoading(true)
        const result = await fetchNews()
        if (result) {
          setData(result)
          setLoading(false)
          return
        }

        // 404 — trigger and start polling
        try {
          await aiApi.triggerNews(ticker, controller.signal)
        } catch {
          // trigger failed — show empty state immediately
          setTimedOut(true)
          setLoading(false)
          return
        }

        startPolling()
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return
        setTimedOut(true)
        setLoading(false)
      }
    }

    init()

    return () => {
      controller.abort()
      if (pollInterval) clearInterval(pollInterval)
    }
  }, [ticker])

  if (loading) {
    return <LoadingSkeleton />
  }

  if (timedOut || !data) {
    return (
      <div className="bg-card border rounded-sm flex flex-col items-center justify-center py-12 gap-2">
        <Newspaper className="h-8 w-8 text-muted-foreground/50" />
        <p className="text-sm text-muted-foreground">Новостей пока нет</p>
      </div>
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
      <div className="bg-card border rounded-sm flex flex-col items-center justify-center py-12 gap-2">
        <Newspaper className="h-8 w-8 text-muted-foreground/50" />
        <p className="text-sm text-muted-foreground">Новостей пока нет</p>
      </div>
    )
  }

  return (
    <div className="space-y-3.5 font-sans text-foreground">
      {/* Row 1: Latest news (wide) + Upcoming company events (narrow) */}
      <div className="flex gap-3.5">
        {(data.latest_news?.length ?? 0) > 0 && (
          <TermPanel
            title="Последние новости"
            count={data.latest_news.length}
            accentClass="bg-primary"
            className="flex-[1.6]"
          >
            {data.latest_news.map((n, i) => (
              <NewsRow key={i} n={n} />
            ))}
          </TermPanel>
        )}
        {(data.upcoming_company_events?.length ?? 0) > 0 && (
          <TermPanel
            title="Предстоящие · компания"
            count={data.upcoming_company_events.length}
            accentClass="bg-primary"
            className="flex-1"
          >
            {data.upcoming_company_events.map((e, i) => (
              <UpcomingCompanyRow key={i} e={e} />
            ))}
          </TermPanel>
        )}
      </div>

      {/* Row 2: Upcoming deps | Past deps | Historic */}
      <div className="flex gap-3.5">
        {(data.upcoming_dependency_events?.length ?? 0) > 0 && (
          <TermPanel
            title="Предстоящие · зависимости"
            count={data.upcoming_dependency_events.length}
            accentClass="bg-blue-600 dark:bg-blue-400"
            className="flex-1"
          >
            {data.upcoming_dependency_events.map((e, i) => (
              <UpcomingDepsRow key={i} e={e} />
            ))}
          </TermPanel>
        )}
        {(data.past_dependency_events?.length ?? 0) > 0 && (
          <TermPanel
            title="Прошедшие · зависимости"
            count={data.past_dependency_events.length}
            accentClass="bg-muted-foreground"
            className="flex-1"
          >
            {data.past_dependency_events.map((e, i) => (
              <PastDepsRow key={i} e={e} />
            ))}
          </TermPanel>
        )}
        {(data.historical_events?.length ?? 0) > 0 && (
          <TermPanel
            title="Исторические события"
            count={data.historical_events.length}
            accentClass="bg-violet-600 dark:bg-violet-400"
            className="flex-1"
          >
            {data.historical_events.map((e, i) => (
              <HistoricRow key={i} e={e} />
            ))}
          </TermPanel>
        )}
      </div>
    </div>
  )
}
