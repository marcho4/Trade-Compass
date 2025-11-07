"use client"

import { Company } from "@/types"
import { Badge } from "@/components/ui/badge"
import { Card } from "@/components/ui/card"
import { ArrowUp, ArrowDown, TrendingUp, Building2, Users } from "lucide-react"

interface CompanyHeaderProps {
  company: Company
}

export const CompanyHeader = ({ company }: CompanyHeaderProps) => {
  const isPricePositive = (company.priceChange24h || 0) >= 0

  const formatNumber = (num: number) => {
    return new Intl.NumberFormat("ru-RU").format(num)
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("ru-RU", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(price)
  }

  return (
    <Card className="p-6">
      <div className="flex flex-col gap-6">
        {/* Основная информация */}
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold">{company.name}</h1>
              <Badge variant="secondary" className="text-sm">
                {company.ticker}
              </Badge>
              <Badge variant="outline">{company.sector}</Badge>
            </div>
            <p className="text-muted-foreground">ИНН: {company.inn}</p>
          </div>

          {/* Цена */}
          {company.currentPrice && (
            <div className="text-right">
              <div className="text-3xl font-bold mb-1">
                ₽{formatPrice(company.currentPrice)}
              </div>
              <div
                className={`flex items-center justify-end gap-1 text-sm font-medium ${
                  isPricePositive ? "text-green-600" : "text-red-600"
                }`}
              >
                {isPricePositive ? (
                  <ArrowUp className="w-4 h-4" />
                ) : (
                  <ArrowDown className="w-4 h-4" />
                )}
                {Math.abs(company.priceChange24h || 0).toFixed(2)}%
                <span className="text-muted-foreground ml-1">(24ч)</span>
              </div>
            </div>
          )}
        </div>

        {/* Дополнительная информация */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Building2 className="w-5 h-5 text-primary" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Владелец</p>
              <p className="font-medium text-sm">{company.owner}</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Users className="w-5 h-5 text-primary" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Сотрудники</p>
              <p className="font-medium">{formatNumber(company.employees)}</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <TrendingUp className="w-5 h-5 text-primary" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Размер лота</p>
              <p className="font-medium">{company.lotSize}</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Users className="w-5 h-5 text-primary" />
            </div>
            <div>
              <p className="text-sm text-muted-foreground">CEO</p>
              <p className="font-medium text-sm">{company.ceo}</p>
            </div>
          </div>
        </div>
      </div>
    </Card>
  )
}

