"use client"

// import { useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { ArrowLeft } from "lucide-react"
import { Button } from "@/components/ui/button"
// import {
//   PortfolioComposition,
//   RiskSelector,
//   PortfolioGoal,
//   PortfolioPerformanceChart,
//   type RiskLevel
// } from "@/components/portfolio"
// import {
//   generatePerformanceData,
//   getPortfolioById,
//   getPositionsByPortfolioId
// } from "@/lib/mock-data"

const PortfolioDetailPage = () => {
  const params = useParams()
  const router = useRouter()
  const portfolioId = params.id as string

  // TODO: Загружать данные портфеля из API
  // const portfolio = getPortfolioById(portfolioId)
  // const positions = getPositionsByPortfolioId(portfolioId)
  // const [selectedRisk, setSelectedRisk] = useState<RiskLevel>(portfolio?.risk || "moderate")
  // const performanceData = generatePerformanceData()

  return (
    <div className="space-y-8">
      {/* Навигация */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => router.push("/dashboard/portfolio")}
          aria-label="Вернуться к списку портфелей"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Назад к портфелям
        </Button>
      </div>

      {/* TODO: Вернуть содержимое когда будет API для портфелей */}
      <div className="rounded-lg border bg-card p-12 text-center">
        <h2 className="text-2xl font-semibold mb-2">Портфель {portfolioId}</h2>
        <p className="text-muted-foreground">
          Данные портфеля будут доступны после подключения API
        </p>
      </div>
    </div>
  )
}

export default PortfolioDetailPage

