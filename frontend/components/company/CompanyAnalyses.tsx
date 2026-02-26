"use client"

import { useEffect, useState, useCallback } from "react"
import ReactMarkdown from "react-markdown"
import remarkGfm from "remark-gfm"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { BrainCircuit, CalendarDays, ChevronDown, ChevronUp } from "lucide-react"
import { Button } from "@/components/ui/button"
import { aiApi } from "@/lib/api"
import type { AvailablePeriod } from "@/lib/api"

interface CompanyAnalysesProps {
  ticker: string
}

const periodLabel = (period: number): string => {
  switch (period) {
    case 3: return "1 квартал"
    case 6: return "Полугодие"
    case 9: return "9 месяцев"
    case 12: return "Годовой"
    default: return `${period} мес.`
  }
}

const periodBadgeVariant = (period: number) => {
  return period === 12 ? "default" as const : "secondary" as const
}

const periodKey = (p: AvailablePeriod) => `${p.year}-${p.period}`

export const CompanyAnalyses = ({ ticker }: CompanyAnalysesProps) => {
  const [periods, setPeriods] = useState<AvailablePeriod[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [expandedKey, setExpandedKey] = useState<string | null>(null)
  const [analysisCache, setAnalysisCache] = useState<Record<string, string>>({})
  const [analysisLoading, setAnalysisLoading] = useState<string | null>(null)
  const [analysisError, setAnalysisError] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()

    const fetchPeriods = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await aiApi.getAvailablePeriods(ticker, controller.signal)
        setPeriods(data)
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return
        console.error(`Failed to fetch periods for ${ticker}:`, err)
        setError("Не удалось загрузить анализы")
      } finally {
        setLoading(false)
      }
    }

    fetchPeriods()
    return () => controller.abort()
  }, [ticker])

  const toggleExpand = useCallback(async (p: AvailablePeriod) => {
    const key = periodKey(p)

    if (expandedKey === key) {
      setExpandedKey(null)
      return
    }

    setExpandedKey(key)
    setAnalysisError(null)

    if (analysisCache[key]) return

    try {
      setAnalysisLoading(key)
      const text = await aiApi.getAnalysis(ticker, p.year, p.period)
      setAnalysisCache(prev => ({ ...prev, [key]: text }))
    } catch (err) {
      console.error(`Failed to fetch analysis for ${key}:`, err)
      setAnalysisError("Не удалось загрузить анализ")
    } finally {
      setAnalysisLoading(null)
    }
  }, [ticker, expandedKey, analysisCache])

  const groupedByYear = periods.reduce<Record<number, AvailablePeriod[]>>((acc, p) => {
    if (!acc[p.year]) acc[p.year] = []
    acc[p.year].push(p)
    return acc
  }, {})

  const sortedYears = Object.keys(groupedByYear)
    .map(Number)
    .sort((a, b) => b - a)

  if (loading) {
    return (
      <Card>
        <CardHeader className="pb-4">
          <Skeleton className="h-6 w-32" />
        </CardHeader>
        <CardContent className="space-y-4">
          {[1, 2, 3].map(i => (
            <div key={i} className="space-y-2">
              <Skeleton className="h-4 w-20" />
              <Skeleton className="h-12 w-full rounded-lg" />
              <Skeleton className="h-12 w-full rounded-lg" />
            </div>
          ))}
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Card>
        <CardContent className="py-12 text-center">
          <p className="text-sm text-muted-foreground">{error}</p>
        </CardContent>
      </Card>
    )
  }

  if (periods.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12 gap-2">
          <BrainCircuit className="h-8 w-8 text-muted-foreground/50" />
          <p className="text-sm text-muted-foreground">
            AI анализы для {ticker} пока не готовы
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-lg">
            <BrainCircuit className="h-5 w-5" />
            AI Анализ
          </CardTitle>
          <span className="text-sm text-muted-foreground">
            {periods.length} {formatAnalysisCount(periods.length)}
          </span>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {sortedYears.map((year) => {
          const yearPeriods = groupedByYear[year].sort(
            (a, b) => b.period - a.period
          )

          return (
            <div key={year}>
              <div className="flex items-center gap-2 mb-3">
                <CalendarDays className="h-4 w-4 text-muted-foreground" />
                <h3 className="text-sm font-semibold text-muted-foreground">{year}</h3>
              </div>
              <div className="space-y-2">
                {yearPeriods.map((p) => {
                  const key = periodKey(p)
                  const isExpanded = expandedKey === key
                  const cachedText = analysisCache[key]
                  const isLoading = analysisLoading === key

                  return (
                    <div
                      key={key}
                      className="rounded-lg border bg-card transition-colors"
                    >
                      <div
                        className="flex items-center justify-between p-3 cursor-pointer hover:bg-accent/50 rounded-lg"
                        onClick={() => toggleExpand(p)}
                      >
                        <div className="flex items-center gap-3">
                          <div className="p-2 rounded-md bg-primary/10">
                            <BrainCircuit className="h-4 w-4 text-primary" />
                          </div>
                          <div>
                            <p className="text-sm font-medium">
                              {ticker} — {periodLabel(p.period)} {p.year}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge variant={periodBadgeVariant(p.period)}>
                            {periodLabel(p.period)}
                          </Badge>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8"
                          >
                            {isExpanded
                              ? <ChevronUp className="h-4 w-4" />
                              : <ChevronDown className="h-4 w-4" />}
                          </Button>
                        </div>
                      </div>
                      {isExpanded && (
                        <div className="px-4 pb-4 pt-2 border-t">
                          {isLoading ? (
                            <div className="space-y-3">
                              <Skeleton className="h-4 w-full" />
                              <Skeleton className="h-4 w-[90%]" />
                              <Skeleton className="h-4 w-[95%]" />
                              <Skeleton className="h-4 w-[85%]" />
                              <Skeleton className="h-4 w-[70%]" />
                            </div>
                          ) : analysisError && !cachedText ? (
                            <p className="text-sm text-muted-foreground">{analysisError}</p>
                          ) : (
                            <div className="prose prose-sm dark:prose-invert max-w-none
                              prose-headings:mt-6 prose-headings:mb-3
                              prose-h1:text-xl prose-h1:font-bold prose-h1:border-b prose-h1:pb-2
                              prose-h2:text-lg prose-h2:font-semibold prose-h2:border-b prose-h2:border-border prose-h2:pb-1.5
                              prose-h3:text-base prose-h3:font-medium
                              prose-p:my-2 prose-p:leading-relaxed
                              prose-table:my-4 prose-table:text-sm
                              prose-th:bg-muted prose-th:px-3 prose-th:py-2 prose-th:text-left prose-th:font-medium
                              prose-td:px-3 prose-td:py-2 prose-td:border-border
                              prose-hr:my-6 prose-hr:border-border
                              prose-ul:my-2 prose-ol:my-2
                              prose-li:my-0.5
                              prose-blockquote:text-xs prose-blockquote:text-muted-foreground prose-blockquote:border-border
                              prose-strong:text-foreground
                            ">
                              <ReactMarkdown remarkPlugins={[remarkGfm]}>{cachedText ?? ""}</ReactMarkdown>
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  )
                })}
              </div>
            </div>
          )
        })}
      </CardContent>
    </Card>
  )
}

const formatAnalysisCount = (count: number): string => {
  const lastTwo = count % 100
  const lastOne = count % 10

  if (lastTwo >= 11 && lastTwo <= 19) return "анализов"
  if (lastOne === 1) return "анализ"
  if (lastOne >= 2 && lastOne <= 4) return "анализа"
  return "анализов"
}
