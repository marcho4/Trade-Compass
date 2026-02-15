"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { BrainCircuit, CalendarDays, Loader2, ChevronDown, ChevronUp } from "lucide-react"
import { Button } from "@/components/ui/button"
import { aiApi } from "@/lib/api"
import type { AnalysisReport } from "@/lib/api"

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

export const CompanyAnalyses = ({ ticker }: CompanyAnalysesProps) => {
  const [analyses, setAnalyses] = useState<AnalysisReport[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [expandedId, setExpandedId] = useState<number | null>(null)

  useEffect(() => {
    const controller = new AbortController()

    const fetchAnalyses = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await aiApi.getAnalysesByTicker(ticker, controller.signal)
        setAnalyses(data)
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return
        console.error(`Failed to fetch analyses for ${ticker}:`, err)
        setError("Не удалось загрузить анализы")
      } finally {
        setLoading(false)
      }
    }

    fetchAnalyses()
    return () => controller.abort()
  }, [ticker])

  const groupedByYear = analyses.reduce<Record<number, AnalysisReport[]>>((acc, report) => {
    if (!acc[report.year]) acc[report.year] = []
    acc[report.year].push(report)
    return acc
  }, {})

  const sortedYears = Object.keys(groupedByYear)
    .map(Number)
    .sort((a, b) => b - a)

  const toggleExpand = (id: number) => {
    setExpandedId(prev => prev === id ? null : id)
  }

  if (loading) {
    return (
      <Card>
        <CardContent className="flex items-center justify-center py-12">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          <span className="ml-2 text-sm text-muted-foreground">Загрузка анализов...</span>
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

  if (analyses.length === 0) {
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
            {analyses.length} {formatAnalysisCount(analyses.length)}
          </span>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {sortedYears.map((year) => {
          const yearAnalyses = groupedByYear[year].sort(
            (a, b) => b.period - a.period
          )

          return (
            <div key={year}>
              <div className="flex items-center gap-2 mb-3">
                <CalendarDays className="h-4 w-4 text-muted-foreground" />
                <h3 className="text-sm font-semibold text-muted-foreground">{year}</h3>
              </div>
              <div className="space-y-2">
                {yearAnalyses.map((report) => (
                  <div
                    key={report.id}
                    className="rounded-lg border bg-card transition-colors"
                  >
                    <div
                      className="flex items-center justify-between p-3 cursor-pointer hover:bg-accent/50 rounded-lg"
                      onClick={() => toggleExpand(report.id)}
                    >
                      <div className="flex items-center gap-3">
                        <div className="p-2 rounded-md bg-primary/10">
                          <BrainCircuit className="h-4 w-4 text-primary" />
                        </div>
                        <div>
                          <p className="text-sm font-medium">
                            {ticker} — {periodLabel(report.period)} {report.year}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant={periodBadgeVariant(report.period)}>
                          {periodLabel(report.period)}
                        </Badge>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                        >
                          {expandedId === report.id
                            ? <ChevronUp className="h-4 w-4" />
                            : <ChevronDown className="h-4 w-4" />}
                        </Button>
                      </div>
                    </div>
                    {expandedId === report.id && (
                      <div className="px-4 pb-4 pt-2 border-t">
                        <div className="prose prose-sm dark:prose-invert max-w-none whitespace-pre-wrap text-sm leading-relaxed">
                          {report.analysis}
                        </div>
                      </div>
                    )}
                  </div>
                ))}
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
