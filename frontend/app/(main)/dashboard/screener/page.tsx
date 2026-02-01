"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { CompanyCard, ScreenerFilters } from "@/components/screener"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { financialDataApi, Sector } from "@/lib/api-client"

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

const MOCK_COMPANIES = [
  {
    id: 1,
    ticker: "GAZP",
    name: "ПАО Газпром",
    sector: "Нефть и газ",
    price: 155.32,
    priceChange: 2.45,
    priceChangePercent: 1.6,
    rating: {
      profitability: 5,
      growth: 10,
      valuation: 8,
      financial_health: 12,
      efficiency: 7,
    },
    marketCap: 3680000000000,
    pe: 3.2,
    dividendYield: 12.5,
  },
  {
    id: 2,
    ticker: "SBER",
    name: "ПАО Сбербанк",
    sector: "Финансы",
    price: 289.5,
    priceChange: -3.2,
    priceChangePercent: -1.09,
    rating: {
      profitability: 25,
      growth: 30,
      valuation: 28,
      financial_health: 32,
      efficiency: 27,
    },
    marketCap: 6450000000000,
    pe: 5.8,
    dividendYield: 8.2,
  },
  {
    id: 3,
    ticker: "LKOH",
    name: "ПАО НК ЛУКОЙЛ",
    sector: "Нефть и газ",
    price: 7245.0,
    priceChange: 125.5,
    priceChangePercent: 1.76,
    rating: {
      profitability: 45,
      growth: 50,
      valuation: 48,
      financial_health: 52,
      efficiency: 47,
    },
    marketCap: 5100000000000,
    pe: 4.5,
    dividendYield: 10.3,
  },
  {
    id: 4,
    ticker: "YNDX",
    name: "Яндекс",
    sector: "Технологии",
    price: 4250.0,
    priceChange: 85.0,
    priceChangePercent: 2.04,
    rating: {
      profitability: 65,
      growth: 68,
      valuation: 62,
      financial_health: 70,
      efficiency: 67,
    },
    marketCap: 1350000000000,
    pe: 18.5,
    dividendYield: 0,
  },
  {
    id: 5,
    ticker: "ROSN",
    name: "ПАО НК Роснефть",
    sector: "Нефть и газ",
    price: 625.8,
    priceChange: -8.2,
    priceChangePercent: -1.29,
    rating: {
      profitability: 80,
      growth: 75,
      valuation: 82,
      financial_health: 78,
      efficiency: 77,
    },
    marketCap: 6620000000000,
    pe: 3.8,
    dividendYield: 11.2,
  },
  {
    id: 6,
    ticker: "GMKN",
    name: "ПАО ГМК Норильский никель",
    sector: "Металлургия",
    price: 16780.0,
    priceChange: 220.0,
    priceChangePercent: 1.33,
    rating: {
      profitability: 92,
      growth: 88,
      valuation: 95,
      financial_health: 90,
      efficiency: 93,
    },
    marketCap: 2660000000000,
    pe: 6.2,
    dividendYield: 15.8,
  },
  {
    id: 7,
    ticker: "MTSS",
    name: "ПАО МТС",
    sector: "Телекоммуникации",
    price: 310.5,
    priceChange: 4.3,
    priceChangePercent: 1.4,
    rating: {
      profitability: 98,
      growth: 95,
      valuation: 100,
      financial_health: 97,
      efficiency: 96,
    },
    marketCap: 620000000000,
    pe: 7.5,
    dividendYield: 9.5,
  },
  {
    id: 8,
    ticker: "NVTK",
    name: "ПАО НОВАТЭК",
    sector: "Нефть и газ",
    price: 1285.0,
    priceChange: 15.0,
    priceChangePercent: 1.18,
    rating: {
      profitability: 35,
      growth: 40,
      valuation: 38,
      financial_health: 42,
      efficiency: 37,
    },
    marketCap: 3840000000000,
    pe: 8.9,
    dividendYield: 7.8,
  },
  {
    id: 9,
    ticker: "TATN",
    name: "ПАО Татнефть",
    sector: "Нефть и газ",
    price: 685.2,
    priceChange: -5.8,
    priceChangePercent: -0.84,
    rating: {
      profitability: 55,
      growth: 58,
      valuation: 52,
      financial_health: 60,
      efficiency: 57,
    },
    marketCap: 1450000000000,
    pe: 4.2,
    dividendYield: 13.2,
  },
]

const ITEMS_PER_PAGE = 9

export default function ScreenerPage() {
  const router = useRouter()
  const [currentPage, setCurrentPage] = useState(1)
  const [sectors, setSectors] = useState<Sector[]>([])
  const [sectorsLoading, setSectorsLoading] = useState(true)
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
    const loadSectors = async () => {
      try {
        const data = await financialDataApi.getSectors()
        setSectors(data)
      } catch (error) {
        console.error("Failed to load sectors:", error)
      } finally {
        setSectorsLoading(false)
      }
    }
    
    loadSectors()
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

  // Фильтрация компаний (в реальности будет на бэкенде)
  const filteredCompanies = MOCK_COMPANIES.filter((company) => {
    // Поиск по названию или тикеру
    if (
      filters.search &&
      !company.name.toLowerCase().includes(filters.search.toLowerCase()) &&
      !company.ticker.toLowerCase().includes(filters.search.toLowerCase())
    ) {
      return false
    }

    // Фильтр по сектору
    if (filters.sector) {
      const sector = sectors.find((s) => s.id === parseInt(filters.sector))
      if (sector && company.sector !== sector.name) {
        return false
      }
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

