"use client"

import { useState } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { X, Filter, ChevronDown, ChevronUp } from "lucide-react"

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

interface ScreenerFiltersProps {
  filters: FilterValues
  onFilterChange: (key: keyof FilterValues, value: string) => void
  onReset: () => void
  sectors: Array<{ id: number; name: string }>
}

export const ScreenerFilters = ({
  filters,
  onFilterChange,
  onReset,
  sectors,
}: ScreenerFiltersProps) => {
  const [showAdvanced, setShowAdvanced] = useState(false)

  const activeFiltersCount = Object.entries(filters).filter(
    ([key, value]) => value !== "" && key !== "search"
  ).length

  const handleToggleAdvanced = () => {
    setShowAdvanced((prev) => !prev)
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <div className="space-y-4">
          {/* Первый ряд - основные фильтры */}
          <div className="flex flex-wrap items-end gap-4">
            {/* Поиск */}
            <div className="flex-1 min-w-[250px] space-y-2">
              <Label htmlFor="search" className="text-sm font-medium">
                Поиск
              </Label>
              <Input
                id="search"
                type="text"
                placeholder="Название или тикер"
                value={filters.search}
                onChange={(e) => onFilterChange("search", e.target.value)}
                aria-label="Поиск по названию или тикеру компании"
                className="h-10"
              />
            </div>

            {/* Сектор */}
            <div className="w-[200px] space-y-2">
              <Label htmlFor="sector" className="text-sm font-medium">
                Сектор
              </Label>
              <Select
                value={filters.sector || "all"}
                onValueChange={(value) =>
                  onFilterChange("sector", value === "all" ? "" : value)
                }
              >
                <SelectTrigger id="sector" aria-label="Выбрать сектор" className="h-10">
                  <SelectValue placeholder="Все секторы" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Все секторы</SelectItem>
                  {sectors.map((sector) => (
                    <SelectItem key={sector.id} value={sector.id.toString()}>
                      {sector.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Рейтинг */}
            <div className="w-[200px] space-y-2">
              <Label htmlFor="ratingMin" className="text-sm font-medium">
                Рейтинг
              </Label>
              <Select
                value={filters.ratingMin || "any"}
                onValueChange={(value) =>
                  onFilterChange("ratingMin", value === "any" ? "" : value)
                }
              >
                <SelectTrigger id="ratingMin" aria-label="Минимальный рейтинг" className="h-10">
                  <SelectValue placeholder="Любой" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="any">Любой</SelectItem>
                  <SelectItem value="90">90+</SelectItem>
                  <SelectItem value="80">80+</SelectItem>
                  <SelectItem value="70">70+</SelectItem>
                  <SelectItem value="60">60+</SelectItem>
                  <SelectItem value="50">50+</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Кнопки действий */}
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="default"
                onClick={handleToggleAdvanced}
                className="h-10 gap-2"
                aria-label={showAdvanced ? "Скрыть дополнительные фильтры" : "Показать дополнительные фильтры"}
              >
                <Filter className="h-4 w-4" />
                Еще
                {activeFiltersCount > 0 && (
                  <Badge variant="secondary" className="ml-1 h-5 min-w-5 px-1.5">
                    {activeFiltersCount}
                  </Badge>
                )}
                {showAdvanced ? (
                  <ChevronUp className="h-4 w-4" />
                ) : (
                  <ChevronDown className="h-4 w-4" />
                )}
              </Button>

              {activeFiltersCount > 0 && (
                <Button
                  variant="ghost"
                  size="default"
                  onClick={onReset}
                  className="h-10 gap-2"
                  aria-label="Сбросить фильтры"
                >
                  <X className="h-4 w-4" />
                  Сбросить
                </Button>
              )}
            </div>
          </div>

          {/* Расширенные фильтры */}
          {showAdvanced && (
            <div className="pt-4 border-t space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {/* P/E */}
                <div className="space-y-2">
                  <Label className="text-sm font-medium">P/E</Label>
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      placeholder="Мин"
                      value={filters.peMin}
                      onChange={(e) => onFilterChange("peMin", e.target.value)}
                      aria-label="Минимальное значение P/E"
                      className="h-10"
                    />
                    <Input
                      type="number"
                      placeholder="Макс"
                      value={filters.peMax}
                      onChange={(e) => onFilterChange("peMax", e.target.value)}
                      aria-label="Максимальное значение P/E"
                      className="h-10"
                    />
                  </div>
                </div>

                {/* Капитализация */}
                <div className="space-y-2">
                  <Label className="text-sm font-medium">Капитализация (млрд ₽)</Label>
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      placeholder="Мин"
                      value={filters.marketCapMin}
                      onChange={(e) => onFilterChange("marketCapMin", e.target.value)}
                      aria-label="Минимальная капитализация"
                      className="h-10"
                    />
                    <Input
                      type="number"
                      placeholder="Макс"
                      value={filters.marketCapMax}
                      onChange={(e) => onFilterChange("marketCapMax", e.target.value)}
                      aria-label="Максимальная капитализация"
                      className="h-10"
                    />
                  </div>
                </div>

                {/* Дивидендная доходность */}
                <div className="space-y-2">
                  <Label htmlFor="dividendYieldMin" className="text-sm font-medium">
                    Дивиденды (мин %)
                  </Label>
                  <Input
                    id="dividendYieldMin"
                    type="number"
                    placeholder="Например: 5"
                    value={filters.dividendYieldMin}
                    onChange={(e) => onFilterChange("dividendYieldMin", e.target.value)}
                    aria-label="Минимальная дивидендная доходность"
                    className="h-10"
                  />
                </div>

                {/* ROE */}
                <div className="space-y-2">
                  <Label htmlFor="roeMin" className="text-sm font-medium">
                    ROE (мин %)
                  </Label>
                  <Input
                    id="roeMin"
                    type="number"
                    placeholder="Например: 15"
                    value={filters.roeMin}
                    onChange={(e) => onFilterChange("roeMin", e.target.value)}
                    aria-label="Минимальный ROE"
                    className="h-10"
                  />
                </div>

                {/* Долг к капиталу */}
                <div className="space-y-2">
                  <Label htmlFor="debtToEquityMax" className="text-sm font-medium">
                    Долг/Капитал (макс)
                  </Label>
                  <Input
                    id="debtToEquityMax"
                    type="number"
                    placeholder="Например: 1.5"
                    value={filters.debtToEquityMax}
                    onChange={(e) => onFilterChange("debtToEquityMax", e.target.value)}
                    aria-label="Максимальный Долг/Капитал"
                    className="h-10"
                  />
                </div>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

