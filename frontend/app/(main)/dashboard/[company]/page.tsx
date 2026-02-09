"use client"

import { use, useEffect, useState } from "react"
import { notFound } from "next/navigation"
import { 
  CompanyHeader,
  CompanyReports,
  KeyMetricsGrid, 
  FinancialChart, 
  FinancialStatements 
} from "@/components/company"
import { getMockCompanyAnalysis } from "@/lib/mock-data"
import { financialDataApi, Sector } from "@/lib/api"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card } from "@/components/ui/card"
import { Company as CompanyType } from "@/types"

type PageProps = {
  params: Promise<{
    company: string
  }>
}

const CompanyDashboardPage = ({ params }: PageProps) => {
  const { company: ticker } = use(params)
  const decodedTicker = decodeURIComponent(ticker).toUpperCase()
  
  const [company, setCompany] = useState<CompanyType | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  useEffect(() => {
    const loadCompany = async () => {
      try {
        const [companyData, sectorsData] = await Promise.all([
          financialDataApi.getCompanyByTicker(decodedTicker),
          financialDataApi.getSectors(),
        ])
        
        const sector = sectorsData.find((s: Sector) => s.id === companyData.sectorId)
        
        setCompany({
          id: companyData.id,
          ticker: companyData.ticker,
          sectorId: companyData.sectorId,
          sector: sector?.name,
          lotSize: companyData.lotSize,
          ceo: companyData.ceo,
        })
      } catch (err) {
        console.error("Failed to load company:", err)
        setError("Компания не найдена")
      } finally {
        setLoading(false)
      }
    }

    loadCompany()
  }, [decodedTicker])

  // Получаем моковые финансовые данные (пока нет API)
  const mockAnalysis = getMockCompanyAnalysis(decodedTicker)
  const latestMetrics = mockAnalysis?.latestMetrics
  const latestIndicators = mockAnalysis?.latestIndicators
  const historicalMetrics = mockAnalysis?.historicalMetrics || []
  const historicalIndicators = mockAnalysis?.historicalIndicators || []
  const industryAverages = mockAnalysis?.industryAverages
  const hasFinancials = !!(latestMetrics && latestIndicators)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">Загрузка...</p>
      </div>
    )
  }

  if (error || !company) {
    notFound()
  }

  return (
    <div className="space-y-8">
      {/* Шапка с основной информацией о компании */}
      <CompanyHeader company={company} />

      {/* Tabs для переключения между разными видами данных */}
      <Tabs defaultValue={hasFinancials ? "overview" : "reports"} className="w-full">
        <TabsList className="grid w-full lg:w-auto" style={{ gridTemplateColumns: `repeat(${hasFinancials ? 4 : 1}, minmax(0, 1fr))` }}>
          {hasFinancials && <TabsTrigger value="overview">Обзор</TabsTrigger>}
          {hasFinancials && <TabsTrigger value="financials">Финансы</TabsTrigger>}
          <TabsTrigger value="reports">Отчёты</TabsTrigger>
          {hasFinancials && <TabsTrigger value="charts">Графики</TabsTrigger>}
        </TabsList>

        {hasFinancials && latestMetrics && latestIndicators && (
          <>
            {/* Обзор - Ключевые метрики */}
            <TabsContent value="overview" className="space-y-8 mt-6">
              <KeyMetricsGrid 
                indicators={latestIndicators} 
                industryAverages={industryAverages} 
              />

              {/* Краткая финансовая информация */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <Card className="p-4">
                  <p className="text-sm text-muted-foreground mb-1">Выручка</p>
                  <p className="text-2xl font-bold">
                    {new Intl.NumberFormat("ru-RU", {
                      notation: "compact",
                      compactDisplay: "short",
                    }).format(latestMetrics.revenue)} ₽
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {latestMetrics.reportPeriod} {latestMetrics.reportYear}
                  </p>
                </Card>

                <Card className="p-4">
                  <p className="text-sm text-muted-foreground mb-1">Чистая прибыль</p>
                  <p className="text-2xl font-bold">
                    {new Intl.NumberFormat("ru-RU", {
                      notation: "compact",
                      compactDisplay: "short",
                    }).format(latestMetrics.netProfit)} ₽
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {latestMetrics.reportPeriod} {latestMetrics.reportYear}
                  </p>
                </Card>

                <Card className="p-4">
                  <p className="text-sm text-muted-foreground mb-1">Рыночная капитализация</p>
                  <p className="text-2xl font-bold">
                    {new Intl.NumberFormat("ru-RU", {
                      notation: "compact",
                      compactDisplay: "short",
                    }).format(latestMetrics.marketCap)} ₽
                  </p>
                </Card>

                <Card className="p-4">
                  <p className="text-sm text-muted-foreground mb-1">Свободный денежный поток</p>
                  <p className="text-2xl font-bold">
                    {new Intl.NumberFormat("ru-RU", {
                      notation: "compact",
                      compactDisplay: "short",
                    }).format(latestMetrics.freeCashFlow)} ₽
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {latestMetrics.reportPeriod} {latestMetrics.reportYear}
                  </p>
                </Card>
              </div>
            </TabsContent>

            {/* Финансовая отчетность */}
            <TabsContent value="financials" className="space-y-6 mt-6">
              <FinancialStatements metrics={latestMetrics} />
            </TabsContent>

            {/* Графики */}
            <TabsContent value="charts" className="space-y-6 mt-6">
              <FinancialChart
                title="Выручка и прибыль"
                description="Динамика выручки, операционной и чистой прибыли по кварталам"
                data={historicalMetrics}
                dataKeys={[
                  { key: "revenue", label: "Выручка", color: "hsl(var(--chart-1))" },
                  { key: "ebit", label: "EBIT", color: "hsl(var(--chart-2))" },
                  { key: "netProfit", label: "Чистая прибыль", color: "hsl(var(--chart-3))" },
                ]}
                chartType="bar"
              />

              <FinancialChart
                title="Денежные потоки"
                description="Операционный, инвестиционный и свободный денежный поток"
                data={historicalMetrics}
                dataKeys={[
                  { key: "operatingCashFlow", label: "Операционный CF", color: "hsl(var(--chart-1))" },
                  { key: "freeCashFlow", label: "Свободный CF", color: "hsl(var(--chart-3))" },
                ]}
                chartType="area"
              />

              <FinancialChart
                title="Маржинальность"
                description="Динамика маржи по кварталам"
                data={historicalIndicators.map((ind, idx) => {
                  const baseMetrics = historicalMetrics[idx]
                  return {
                    ...baseMetrics,
                    revenue: ind.grossProfitMargin || 0,
                    ebit: ind.operatingProfitMargin || 0,
                    netProfit: ind.netProfitMargin || 0,
                  }
                })}
                dataKeys={[
                  { key: "revenue", label: "Валовая маржа", color: "hsl(var(--chart-1))" },
                  { key: "ebit", label: "Операционная маржа", color: "hsl(var(--chart-2))" },
                  { key: "netProfit", label: "Чистая маржа", color: "hsl(var(--chart-3))" },
                ]}
                chartType="line"
                formatValue={(value) => `${value.toFixed(1)}%`}
              />

              <FinancialChart
                title="Рентабельность"
                description="ROE и ROA по кварталам"
                data={historicalIndicators.map((ind, idx) => {
                  const baseMetrics = historicalMetrics[idx]
                  return {
                    ...baseMetrics,
                    grossProfit: ind.roe || 0,
                    operatingExpenses: ind.roa || 0,
                  }
                })}
                dataKeys={[
                  { key: "grossProfit", label: "ROE", color: "hsl(var(--chart-1))" },
                  { key: "operatingExpenses", label: "ROA", color: "hsl(var(--chart-2))" },
                ]}
                chartType="line"
                formatValue={(value) => `${value.toFixed(1)}%`}
              />

              <FinancialChart
                title="Долговая нагрузка"
                description="Общий долг и чистый долг компании"
                data={historicalMetrics}
                dataKeys={[
                  { key: "debt", label: "Долг", color: "hsl(var(--chart-1))" },
                  { key: "netDebt", label: "Чистый долг", color: "hsl(var(--chart-2))" },
                ]}
                chartType="bar"
              />
            </TabsContent>
          </>
        )}

        {/* Отчёты — доступны всегда */}
        <TabsContent value="reports" className="mt-6">
          <CompanyReports ticker={decodedTicker} />
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default CompanyDashboardPage
