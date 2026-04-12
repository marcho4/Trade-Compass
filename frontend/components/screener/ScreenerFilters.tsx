"use client"

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
import { Button } from "@/components/ui/button"
import { X } from "lucide-react"
import type { FilterValues } from "./types"

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
  const activeFiltersCount = Object.entries(filters).filter(
    ([key, value]) => value !== "" && key !== "search"
  ).length

  return (
    <Card>
      <CardContent className="pt-1">
        <div className="flex flex-wrap items-end gap-4">
          <div className="flex-1 min-w-[180px] space-y-2">
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

          <div className="min-w-[140px] w-[200px] flex-1 sm:flex-none space-y-2">
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

          <div className="min-w-[140px] w-[200px] flex-1 sm:flex-none space-y-2">
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
                <SelectItem value="1">1+</SelectItem>
                <SelectItem value="2">2+</SelectItem>
                <SelectItem value="3">3+</SelectItem>
                <SelectItem value="4">4+</SelectItem>
                <SelectItem value="5">5</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center gap-2">
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
      </CardContent>
    </Card>
  )
}
