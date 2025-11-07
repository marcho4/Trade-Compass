"use client"

import { useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { ArrowLeft } from "lucide-react"
import { Button } from "@/components/ui/button"
import { 
  PortfolioComposition, 
  RiskSelector, 
  PortfolioGoal, 
  PortfolioPerformanceChart,
  type RiskLevel 
} from "@/components/portfolio"
import { 
  generatePerformanceData, 
  getPortfolioById, 
  getPositionsByPortfolioId 
} from "@/lib/mock-data"

const PortfolioDetailPage = () => {
  const params = useParams()
  const router = useRouter()
  const portfolioId = params.id as string

  const portfolio = getPortfolioById(portfolioId)
  const positions = getPositionsByPortfolioId(portfolioId)
  const [selectedRisk, setSelectedRisk] = useState<RiskLevel>(portfolio?.risk || "moderate")
  const performanceData = generatePerformanceData()

  if (!portfolio) {
    return (
      <div className="space-y-8">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => router.push("/dashboard/portfolio")}
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Назад к портфелям
          </Button>
        </div>
        <div className="rounded-lg border bg-card p-12 text-center">
          <h2 className="text-2xl font-semibold mb-2">Портфель не найден</h2>
          <p className="text-muted-foreground">
            Портфель с таким ID не существует или был удален
          </p>
        </div>
      </div>
    )
  }

  const handleRiskChange = (risk: RiskLevel) => {
    setSelectedRisk(risk)
    console.log("Изменение риска:", risk)
    // Здесь будет логика обновления риска через API
  }

  const handleUpdateGoal = (goalValue: number, description: string) => {
    console.log("Обновление цели:", goalValue, description)
    // Здесь будет логика обновления цели через API
  }

  const totalPortfolioValue = positions.reduce(
    (sum, position) => sum + position.currentPrice * position.quantity,
    0
  )

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

      {/* Заголовок портфеля */}
      <div>
        <h1 className="text-3xl font-bold text-foreground">{portfolio.name}</h1>
        <p className="mt-2 text-muted-foreground">{portfolio.description}</p>
      </div>

      {/* Состав портфеля */}
      <PortfolioComposition positions={positions} totalValue={totalPortfolioValue} />

      {/* График производительности */}
      <PortfolioPerformanceChart
        data={performanceData}
        portfolioName={portfolio.name}
      />

      {/* Настройки портфеля */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Уровень риска */}
        <RiskSelector value={selectedRisk} onChange={handleRiskChange} />

        {/* Цель портфеля */}
        <PortfolioGoal
          currentValue={portfolio.currentValue}
          goalValue={portfolio.goalValue}
          goalDescription={portfolio.goalDescription}
          onUpdateGoal={handleUpdateGoal}
        />
      </div>

      {/* Дополнительная информация */}
      <div className="rounded-lg border bg-card p-6">
        <h3 className="text-lg font-semibold mb-4">О портфельной стратегии</h3>
        <div className="space-y-3 text-sm text-muted-foreground">
          <p>
            <strong className="text-foreground">Рейтинг платформы</strong> — это
            оценка качества портфеля на основе анализа диверсификации, волатильности,
            фундаментальных показателей компаний и соответствия выбранной стратегии риска.
          </p>
          <p>
            <strong className="text-foreground">Ребалансировка</strong> —
            регулярно проверяйте рекомендации платформы по оптимизации портфеля
            для достижения лучших результатов.
          </p>
          <p>
            <strong className="text-foreground">Сравнение с индексами</strong> —
            отслеживайте как ваш портфель показывает себя относительно рыночных
            индексов, чтобы оценить эффективность стратегии.
          </p>
        </div>
      </div>
    </div>
  )
}

export default PortfolioDetailPage

