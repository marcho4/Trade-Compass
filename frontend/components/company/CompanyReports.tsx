"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { FileText, Download, CalendarDays, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { parserApi } from "@/lib/api-client"
import type { Report } from "@/types"

interface CompanyReportsProps {
  ticker: string
}

const periodLabel = (period: string): string => {
  switch (period) {
    case "3": return "1 квартал"
    case "6": return "Полугодие"
    case "9": return "9 месяцев"
    case "12": return "Годовой"
    default: return `${period} мес.`
  }
}

const periodBadgeVariant = (period: string) => {
  return period === "12" ? "default" as const : "secondary" as const
}

export const CompanyReports = ({ ticker }: CompanyReportsProps) => {
  const [reports, setReports] = useState<Report[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()

    const fetchReports = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await parserApi.getReportsByTicker(ticker, controller.signal)
        setReports(data)
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return
        console.error(`Failed to fetch reports for ${ticker}:`, err)
        setError("Не удалось загрузить отчёты")
      } finally {
        setLoading(false)
      }
    }

    fetchReports()
    return () => controller.abort()
  }, [ticker])

  const groupedByYear = reports.reduce<Record<number, Report[]>>((acc, report) => {
    if (!acc[report.year]) acc[report.year] = []
    acc[report.year].push(report)
    return acc
  }, {})

  const sortedYears = Object.keys(groupedByYear)
    .map(Number)
    .sort((a, b) => b - a)

  if (loading) {
    return (
      <Card>
        <CardContent className="flex items-center justify-center py-12">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          <span className="ml-2 text-sm text-muted-foreground">Загрузка отчётов...</span>
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

  if (reports.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12 gap-2">
          <FileText className="h-8 w-8 text-muted-foreground/50" />
          <p className="text-sm text-muted-foreground">
            Отчёты для {ticker} пока не загружены
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
            <FileText className="h-5 w-5" />
            Отчётность
          </CardTitle>
          <span className="text-sm text-muted-foreground">
            {reports.length} {formatReportCount(reports.length)}
          </span>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {sortedYears.map((year) => {
          const yearReports = groupedByYear[year].sort(
            (a, b) => Number(b.period) - Number(a.period)
          )

          return (
            <div key={year}>
              <div className="flex items-center gap-2 mb-3">
                <CalendarDays className="h-4 w-4 text-muted-foreground" />
                <h3 className="text-sm font-semibold text-muted-foreground">{year}</h3>
              </div>
              <div className="space-y-2">
                {yearReports.map((report) => (
                  <div
                    key={report.id}
                    className="flex items-center justify-between p-3 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center gap-3">
                      <div className="p-2 rounded-md bg-primary/10">
                        <FileText className="h-4 w-4 text-primary" />
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
                        asChild
                      >
                        <a
                          href={report.s3_path}
                          target="_blank"
                          rel="noopener noreferrer"
                          aria-label={`Скачать отчёт ${ticker} за ${periodLabel(report.period)} ${report.year}`}
                        >
                          <Download className="h-4 w-4" />
                        </a>
                      </Button>
                    </div>
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

const formatReportCount = (count: number): string => {
  const lastTwo = count % 100
  const lastOne = count % 10

  if (lastTwo >= 11 && lastTwo <= 19) return "отчётов"
  if (lastOne === 1) return "отчёт"
  if (lastOne >= 2 && lastOne <= 4) return "отчёта"
  return "отчётов"
}
