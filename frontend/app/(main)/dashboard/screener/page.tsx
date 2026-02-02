"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { CompanyCard, ScreenerFilters } from "@/components/screener"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { financialDataApi, Sector, Company } from "@/lib/api-client"

interface FilterValues {
  search: string
  sector: string
  peMin: string
  peMax: string
  dividendYieldMin: string
  marketCapMin: string
  marketCapMax: string
  ratingMin: string
  debtToEquityMax: string
  roeMin: string
}

interface CompanyWithDetails extends Company {
  name: string
  sector: string
  price: number
  priceChange: number
  priceChangePercent: number
  rating: {
    profitability: number
    growth: number
    valuation: number
    financial_health: number
    efficiency: number
  }
  marketCap?: number
  pe?: number
  dividendYield?: number
}

const ITEMS_PER_PAGE = 9

export default function ScreenerPage() {
  const router = useRouter()
  const [currentPage, setCurrentPage] = useState(1)
  const [sectors, setSectors] = useState<Sector[]>([])
  const [companies, setCompanies] = useState<CompanyWithDetails[]>([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState<FilterValues>({
    search: "",
    sector: "",
    peMin: "",
    peMax: "",
    dividendYieldMin: "",
    marketCapMin: "",
    marketCapMax: "",
    ratingMin: "",
    debtToEquityMax: "",
    roeMin: "",
  })

  useEffect(() => {
    const loadData = async () => {
      try {
        const [sectorsData, companiesData] = await Promise.all([
          financialDataApi.getSectors(),
          financialDataApi.getCompanies(),
        ])
        
        setSectors(sectorsData)
        
        const companiesWithDetails: CompanyWithDetails[] = companiesData.map((company) => {
          const sector = sectorsData.find((s) => s.id === company.sectorId)
          return {
            ...company,
            name: company.ticker,
            sector: sector?.name || "Неизвестно",
            price: Math.random() * 1000 + 100,
            priceChange: (Math.random() - 0.5) * 20,
            priceChangePercent: (Math.random() - 0.5) * 5,
            rating: {
              profitability: Math.floor(Math.random() * 100),
              growth: Math.floor(Math.random() * 100),
              valuation: Math.floor(Math.random() * 100),
              financial_health: Math.floor(Math.random() * 100),
              efficiency: Math.floor(Math.random() * 100),
            },
            marketCap: Math.random() * 5000000000000 + 100000000000,
            pe: Math.random() * 20 + 2,
            dividendYield: Math.random() * 15,
          }
        })
        
        setCompanies(companiesWithDetails)
      } catch (error) {
        console.error("Failed to load data:", error)
      } finally {
        setLoading(false)
      }
    }
    
    loadData()
  }, [])

  // Функция для обработки изменения фильтра
  const handleFilterChange = (key: keyof FilterValues, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }))
    setCurrentPage(1) // Сбрасываем на первую страницу при изменении фильтров
  }

  // Функция сброса фильтров
  const handleResetFilters = () => {
    setFilters({
      search: "",
      sector: "",
      peMin: "",
      peMax: "",
      dividendYieldMin: "",
      marketCapMin: "",
      marketCapMax: "",
      ratingMin: "",
      debtToEquityMax: "",
      roeMin: "",
    })
    setCurrentPage(1)
  }

  // Фильтрация компаний
  const filteredCompanies = companies.filter((company) => {
    // Поиск по названию или тикеру
    if (
      filters.search &&
      !company.ticker.toLowerCase().includes(filters.search.toLowerCase())
    ) {
      return false
    }

    // Фильтр по сектору
    if (filters.sector && company.sectorId !== parseInt(filters.sector)) {
      return false
    }

    // Фильтр по P/E
    if (filters.peMin && company.pe && company.pe < parseFloat(filters.peMin)) {
      return false
    }
    if (filters.peMax && company.pe && company.pe > parseFloat(filters.peMax)) {
      return false
    }

    // Фильтр по дивидендной доходности
    if (
      filters.dividendYieldMin &&
      company.dividendYield &&
      company.dividendYield < parseFloat(filters.dividendYieldMin)
    ) {
      return false
    }

    // Фильтр по капитализации (конвертируем в млрд)
    const marketCapBillion = company.marketCap ? company.marketCap / 1000000000 : 0
    if (
      filters.marketCapMin &&
      marketCapBillion < parseFloat(filters.marketCapMin)
    ) {
      return false
    }
    if (
      filters.marketCapMax &&
      marketCapBillion > parseFloat(filters.marketCapMax)
    ) {
      return false
    }

    // Фильтр по рейтингу
    if (filters.ratingMin) {
      const averageRating =
        (company.rating.profitability +
          company.rating.growth +
          company.rating.valuation +
          company.rating.financial_health +
          company.rating.efficiency) /
        5
      if (averageRating < parseFloat(filters.ratingMin)) {
        return false
      }
    }

    return true
  })

  // Пагинация
  const totalPages = Math.ceil(filteredCompanies.length / ITEMS_PER_PAGE)
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE
  const paginatedCompanies = filteredCompanies.slice(
    startIndex,
    startIndex + ITEMS_PER_PAGE
  )

  const handlePreviousPage = () => {
    setCurrentPage((prev) => Math.max(1, prev - 1))
  }

  const handleNextPage = () => {
    setCurrentPage((prev) => Math.min(totalPages, prev + 1))
  }

  const handleCompanyClick = (ticker: string) => {
    router.push(`/dashboard/${ticker}`)
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-center h-64">
          <p className="text-muted-foreground">Загрузка компаний...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">Скринер акций</h1>
        <p className="text-muted-foreground">
          Найдите компании, которые соответствуют вашим инвестиционным критериям
        </p>
      </div>

      {/* Фильтры сверху */}
      <div className="mb-6">
        <ScreenerFilters
          filters={filters}
          onFilterChange={handleFilterChange}
          onReset={handleResetFilters}
          sectors={sectors}
        />
      </div>

      {/* Список компаний */}
      <div className="space-y-6">
        {/* Информация о результатах */}
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Найдено компаний: <span className="font-semibold text-foreground">{filteredCompanies.length}</span>
          </p>
          {totalPages > 1 && (
            <p className="text-sm text-muted-foreground">
              Страница {currentPage} из {totalPages}
            </p>
          )}
        </div>

        {/* Сетка с карточками компаний */}
        {paginatedCompanies.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
            {paginatedCompanies.map((company) => (
              <CompanyCard
                key={company.id}
                {...company}
                onClick={() => handleCompanyClick(company.ticker)}
              />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <p className="text-xl font-semibold text-muted-foreground mb-2">
              Компании не найдены
            </p>
            <p className="text-sm text-muted-foreground mb-4">
              Попробуйте изменить параметры фильтров
            </p>
            <Button onClick={handleResetFilters} variant="outline">
              Сбросить фильтры
            </Button>
          </div>
        )}

        {/* Пагинация */}
        {totalPages > 1 && paginatedCompanies.length > 0 && (
          <div className="flex items-center justify-center gap-2 mt-8">
            <Button
              variant="outline"
              size="icon"
              onClick={handlePreviousPage}
              disabled={currentPage === 1}
              aria-label="Предыдущая страница"
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>

            <div className="flex items-center gap-1">
              {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                (page) => {
                  // Показываем только несколько страниц около текущей
                  if (
                    page === 1 ||
                    page === totalPages ||
                    (page >= currentPage - 1 && page <= currentPage + 1)
                  ) {
                    return (
                      <Button
                        key={page}
                        variant={page === currentPage ? "default" : "outline"}
                        size="icon"
                        onClick={() => setCurrentPage(page)}
                        aria-label={`Страница ${page}`}
                        className="w-10"
                      >
                        {page}
                      </Button>
                    )
                  } else if (
                    page === currentPage - 2 ||
                    page === currentPage + 2
                  ) {
                    return (
                      <span key={page} className="px-2 text-muted-foreground">
                        ...
                      </span>
                    )
                  }
                  return null
                }
              )}
            </div>

            <Button
              variant="outline"
              size="icon"
              onClick={handleNextPage}
              disabled={currentPage === totalPages}
              aria-label="Следующая страница"
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}

