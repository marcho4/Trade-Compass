"use client"

import { useState, useEffect, useMemo } from "react"
import { useRouter } from "next/navigation"
import { CompanyCard, ScreenerFilters } from "@/components/screener"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { financialDataApi, Sector, Company } from "@/lib/api"
import { aiApi } from "@/lib/api/ai-api"
import type { FilterValues, CompanyRating } from "@/components/screener/types"

interface CompanyWithRating extends Company {
  name: string
  sectorName: string
  rating: CompanyRating
}

interface PriceInfo {
  price: number
  priceChange: number
  priceChangePercent: number
}

const ITEMS_PER_PAGE = 9

const DEFAULT_RATING: CompanyRating = {
  health: 0,
  growth: 0,
  moat: 0,
  dividends: 0,
  value: 0,
  total: 0,
}

export default function ScreenerPage() {
  const router = useRouter()
  const [currentPage, setCurrentPage] = useState(1)
  const [sectors, setSectors] = useState<Sector[]>([])
  const [companies, setCompanies] = useState<CompanyWithRating[]>([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState<FilterValues>({
    search: "",
    sector: "",
    ratingMin: "",
  })

  const [priceMap, setPriceMap] = useState<Map<string, PriceInfo>>(new Map())
  const [marketCapMap, setMarketCapMap] = useState<Map<string, number>>(new Map())

  useEffect(() => {
    const controller = new AbortController()

    const loadData = async () => {
      try {
        const [sectorsData, companiesData] = await Promise.all([
          financialDataApi.getSectors(),
          financialDataApi.getCompanies(),
        ])

        setSectors(sectorsData)

        const sectorMap = new Map(sectorsData.map((s) => [s.id, s.name]))

        const reportResults = await Promise.allSettled(
          companiesData.map((c) => aiApi.getReportResults(c.ticker, controller.signal))
        )
        const reportMap = new Map<string, CompanyRating>()
        companiesData.forEach((company, i) => {
          const result = reportResults[i]
          if (result.status === "fulfilled" && result.value) {
            reportMap.set(company.ticker, result.value)
          }
        })

        const enriched: CompanyWithRating[] = companiesData.map((company) => ({
          ...company,
          name: company.name || company.ticker,
          sectorName: sectorMap.get(company.sectorId) || "Неизвестно",
          rating: reportMap.get(company.ticker) || DEFAULT_RATING,
        }))

        setCompanies(enriched)
      } catch (error) {
        console.error("Failed to load data:", error)
      } finally {
        setLoading(false)
      }
    }

    loadData()
    return () => controller.abort()
  }, [])

  const handleFilterChange = (key: keyof FilterValues, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }))
    setCurrentPage(1)
  }

  const handleResetFilters = () => {
    setFilters({
      search: "",
      sector: "",
      ratingMin: "",
    })
    setCurrentPage(1)
  }

  const filteredCompanies = useMemo(() => {
    return companies.filter((company) => {
      if (filters.search) {
        const q = filters.search.toLowerCase()
        if (
          !company.ticker.toLowerCase().includes(q) &&
          !company.name.toLowerCase().includes(q)
        ) {
          return false
        }
      }

      if (filters.sector && company.sectorId !== parseInt(filters.sector)) {
        return false
      }

      if (filters.ratingMin) {
        if (company.rating.total < parseFloat(filters.ratingMin)) {
          return false
        }
      }

      return true
    })
  }, [companies, filters])

  const totalPages = Math.ceil(filteredCompanies.length / ITEMS_PER_PAGE)
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE
  const paginatedCompanies = filteredCompanies.slice(
    startIndex,
    startIndex + ITEMS_PER_PAGE
  )

  const visibleTickers = useMemo(
    () => paginatedCompanies.map((c) => c.ticker),
    [paginatedCompanies]
  )

  useEffect(() => {
    if (visibleTickers.length === 0) return
    const controller = new AbortController()

    const fetchVisibleData = async () => {
      const [priceResults, mcapResults] = await Promise.all([
        Promise.allSettled(
          visibleTickers.map((ticker) =>
            financialDataApi.getPriceCandles(ticker, 2, 24, controller.signal)
          )
        ),
        Promise.allSettled(
          visibleTickers.map((ticker) =>
            financialDataApi.getMarketCap(ticker, controller.signal)
          )
        ),
      ])

      if (controller.signal.aborted) return

      const newPriceMap = new Map(priceMap)
      visibleTickers.forEach((ticker, i) => {
        const result = priceResults[i]
        if (result.status === "fulfilled" && result.value.length > 0) {
          const candles = result.value
          const current = candles[candles.length - 1].close
          const previous = candles.length >= 2 ? candles[candles.length - 2].close : current
          const change = current - previous
          const changePercent = previous !== 0 ? (change / previous) * 100 : 0
          newPriceMap.set(ticker, { price: current, priceChange: change, priceChangePercent: changePercent })
        }
      })
      setPriceMap(newPriceMap)

      const newMcapMap = new Map(marketCapMap)
      visibleTickers.forEach((ticker, i) => {
        const result = mcapResults[i]
        if (result.status === "fulfilled" && result.value) {
          newMcapMap.set(ticker, result.value)
        }
      })
      setMarketCapMap(newMcapMap)
    }

    fetchVisibleData()
    return () => controller.abort()
  }, [visibleTickers.join(",")])

  const handlePreviousPage = () => {
    setCurrentPage((prev) => Math.max(1, prev - 1))
  }

  const handleNextPage = () => {
    setCurrentPage((prev) => Math.min(totalPages, prev + 1))
  }

  const handleCompanyClick = (ticker: string) => {
    router.push(`/dashboard/${ticker}`)
  }

  const paginationPages = useMemo(() => {
    if (totalPages <= 1) return []
    const pages: (number | "ellipsis")[] = []
    const addPage = (p: number) => { if (!pages.includes(p)) pages.push(p) }

    addPage(1)
    if (currentPage - 1 > 2) pages.push("ellipsis")
    for (let p = Math.max(2, currentPage - 1); p <= Math.min(totalPages - 1, currentPage + 1); p++) {
      addPage(p)
    }
    if (currentPage + 1 < totalPages - 1) pages.push("ellipsis")
    if (totalPages > 1) addPage(totalPages)

    return pages
  }, [currentPage, totalPages])

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

      <div className="mb-6">
        <ScreenerFilters
          filters={filters}
          onFilterChange={handleFilterChange}
          onReset={handleResetFilters}
          sectors={sectors}
        />
      </div>

      <div className="space-y-6">
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

        {paginatedCompanies.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
            {paginatedCompanies.map((company) => {
              const priceInfo = priceMap.get(company.ticker)
              const mcap = marketCapMap.get(company.ticker)
              return (
                <CompanyCard
                  key={company.id}
                  id={company.id}
                  ticker={company.ticker}
                  name={company.name}
                  sector={company.sectorName}
                  price={priceInfo?.price ?? 0}
                  priceChange={priceInfo?.priceChange ?? 0}
                  priceChangePercent={priceInfo?.priceChangePercent ?? 0}
                  priceLoading={!priceInfo}
                  rating={company.rating}
                  marketCap={mcap}
                  onClick={() => handleCompanyClick(company.ticker)}
                />
              )
            })}
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

        {paginationPages.length > 0 && paginatedCompanies.length > 0 && (
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
              {paginationPages.map((item, idx) =>
                item === "ellipsis" ? (
                  <span key={`ellipsis-${idx}`} className="px-2 text-muted-foreground" aria-hidden="true">
                    ...
                  </span>
                ) : (
                  <Button
                    key={item}
                    variant={item === currentPage ? "default" : "outline"}
                    size="icon"
                    onClick={() => setCurrentPage(item)}
                    aria-label={`Страница ${item}`}
                    className="w-10"
                  >
                    {item}
                  </Button>
                )
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
