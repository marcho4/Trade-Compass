"use client"

import { Company } from "@/types"
import { Badge } from "@/components/ui/badge"
import { Card } from "@/components/ui/card"
import { ArrowUp, ArrowDown, TrendingUp, Users, Loader2 } from "lucide-react"
import { usePriceData } from "@/hooks/use-price-data"

interface CompanyHeaderProps {
  company: Company
}

export const CompanyHeader = ({ company }: CompanyHeaderProps) => {
  const { price, priceChange, priceChangePercent, loading: priceLoading } = usePriceData(company.ticker)
  const isPricePositive = priceChange >= 0

  const formatPrice = (value: number) => {
    return new Intl.NumberFormat("ru-RU", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(value)
  }

  return (
    <Card className="p-6">
      <div className="flex flex-col gap-6">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold">{company.ticker}</h1>
              {company.sector && (
                <Badge variant="outline">{company.sector}</Badge>
              )}
            </div>
          </div>

          <div className="text-right">
            {priceLoading ? (
              <div className="flex items-center gap-2 text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" />
              </div>
            ) : price > 0 ? (
              <>
                <div className="text-3xl font-bold mb-1">
                  ₽{formatPrice(price)}
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
                  {isPricePositive ? "+" : ""}
                  {priceChangePercent.toFixed(2)}%
                  <span className="text-muted-foreground ml-1">
                    ({isPricePositive ? "+" : ""}
                    {formatPrice(priceChange)} ₽)
                  </span>
                </div>
              </>
            ) : null}
          </div>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-2 gap-4 pt-4 border-t">
          {company.lotSize && (
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg">
                <TrendingUp className="w-5 h-5 text-primary" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Размер лота</p>
                <p className="font-medium">{company.lotSize}</p>
              </div>
            </div>
          )}

          {company.ceo && (
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg">
                <Users className="w-5 h-5 text-primary" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">CEO</p>
                <p className="font-medium text-sm">{company.ceo}</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </Card>
  )
}

