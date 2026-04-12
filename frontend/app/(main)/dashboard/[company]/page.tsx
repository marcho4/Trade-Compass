"use client"

import { use, useEffect, useState } from "react"
import { notFound } from "next/navigation"
import {
  CompanyHeader,
  CompanyReports,
  CompanyAnalyses,
  CompanyMetrics,
  CompanyAbout,
  CompanyNews,
} from "@/components/company"
// import { getMockCompanyAnalysis } from "@/lib/mock-data"
import { financialDataApi, Sector } from "@/lib/api"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
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

  // TODO: Заменить на реальные данные из API
  // const mockAnalysis = getMockCompanyAnalysis(decodedTicker)
  // const latestMetrics = mockAnalysis?.latestMetrics
  // const latestIndicators = mockAnalysis?.latestIndicators
  // const historicalMetrics = mockAnalysis?.historicalMetrics || []
  // const historicalIndicators = mockAnalysis?.historicalIndicators || []
  // const industryAverages = mockAnalysis?.industryAverages

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
      <Tabs defaultValue="about" className="w-full">
        <TabsList className="grid w-full lg:w-auto" style={{ gridTemplateColumns: 'repeat(5, minmax(0, 1fr))' }}>
          <TabsTrigger value="about">О компании</TabsTrigger>
          <TabsTrigger value="metrics">Показатели</TabsTrigger>
          <TabsTrigger value="ai-analysis">AI Анализ</TabsTrigger>
          <TabsTrigger value="news">Новости</TabsTrigger>
          <TabsTrigger value="reports">Отчёты</TabsTrigger>
        </TabsList>

        {/* О компании */}
        <TabsContent value="about" className="mt-6">
          <CompanyAbout ticker={decodedTicker} />
        </TabsContent>

        
        {/* Показатели */}
        <TabsContent value="metrics" className="mt-6">
          <CompanyMetrics ticker={decodedTicker} />
        </TabsContent>

        {/* AI Анализ */}
        <TabsContent value="ai-analysis" className="mt-6">
          <CompanyAnalyses ticker={decodedTicker} />
        </TabsContent>

        {/* Новости */}
        <TabsContent value="news" className="mt-6">
          <CompanyNews ticker={decodedTicker} />
        </TabsContent>


        {/* Отчёты */}
        <TabsContent value="reports" className="mt-6">
          <CompanyReports ticker={decodedTicker} />
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default CompanyDashboardPage
