"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { PortfolioCard } from "@/components/portfolio/PortfolioCard"
import { CreatePortfolioDialog } from "@/components/portfolio/CreatePortfolioDialog"

const PortfolioPage = () => {
  const router = useRouter()
  
  const [portfolios] = useState([
    {
      id: "1",
      name: "Основной портфель",
      value: 1250000,
      createdAt: new Date(2024, 0, 15),
      profitPercent: 18.5,
      profitAmount: 195000,
      rating: 8,
    },
    {
      id: "2",
      name: "Дивидендный портфель",
      value: 850000,
      createdAt: new Date(2024, 3, 10),
      profitPercent: 12.3,
      profitAmount: 93000,
      rating: 7,
    },
  ])

  const handleCreatePortfolio = (name: string, initialAmount: number) => {
    console.log("Создание портфеля:", name, initialAmount)
    // Здесь будет логика создания портфеля через API
  }

  const handlePortfolioClick = (portfolioId: string) => {
    router.push(`/dashboard/portfolio/${portfolioId}`)
  }

  return (
    <div className="space-y-8">
      {/* Заголовок и кнопка создания */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground">
            Мои портфели
          </h1>
          <p className="mt-2 text-muted-foreground">
            Управляйте своими инвестиционными портфелями и отслеживайте их динамику
          </p>
        </div>
        <CreatePortfolioDialog onCreatePortfolio={handleCreatePortfolio} />
      </div>

      {/* Список портфелей */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {portfolios.map((portfolio) => (
          <PortfolioCard
            key={portfolio.id}
            name={portfolio.name}
            value={portfolio.value}
            createdAt={portfolio.createdAt}
            profitPercent={portfolio.profitPercent}
            profitAmount={portfolio.profitAmount}
            rating={portfolio.rating}
            onClick={() => handlePortfolioClick(portfolio.id)}
          />
        ))}
      </div>

      {/* Информационная подсказка */}
      <div className="rounded-lg border bg-card p-6">
        <h3 className="text-lg font-semibold mb-4">Выберите портфель</h3>
        <p className="text-sm text-muted-foreground">
          Кликните на карточку портфеля, чтобы увидеть его детальный состав, графики производительности, 
          настройки риска и цели инвестирования.
        </p>
      </div>
    </div>
  )
}

export default PortfolioPage
