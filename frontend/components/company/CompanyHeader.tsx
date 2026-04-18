"use client"

import { useEffect, useState } from "react"
import { Company } from "@/types"
import { Badge } from "@/components/ui/badge"
import { Card } from "@/components/ui/card"
import { ArrowUp, ArrowDown } from "lucide-react"
import { usePriceData } from "@/hooks/use-price-data"
import { aiApi } from "@/lib/api/ai-api"

interface CompanyHeaderProps {
  company: Company
}

export const CompanyHeader = ({ company }: CompanyHeaderProps) => {
  const { price, priceChange, priceChangePercent, loading: priceLoading } = usePriceData(company.ticker)
  const isPricePositive = priceChange >= 0

  const [companyName, setCompanyName] = useState<string | null>(null)
  const [description, setDescription] = useState<string | null>(null)

  useEffect(() => {
    const controller = new AbortController()
    aiApi.getBusinessResearch(company.ticker, controller.signal)
      .then((data) => {
        if (data) {
          setCompanyName(data.profile.company_name)
          setDescription(data.profile.description)
        } else {
          console.warn("data not found (business research)")
        }
      })
      .catch(() => {})
    return () => controller.abort()
  }, [company.ticker])

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
            {companyName ? (
              <p className="text-lg font-medium text-foreground/80 mb-1">{companyName}</p>
            ) : (
              <div className="h-5 w-48 rounded bg-muted animate-pulse mb-1" />
            )}
            {description ? (
              <p className="text-sm text-muted-foreground leading-relaxed max-w-2xl">{description}</p>
            ) : (
              <div className="space-y-1.5 max-w-2xl">
                <div className="h-3.5 w-full rounded bg-muted animate-pulse" />
                <div className="h-3.5 w-3/4 rounded bg-muted animate-pulse" />
              </div>
            )}
          </div>

          <div className="text-right">
            {priceLoading ? (
              <div className="space-y-2">
                <div className="h-8 w-32 rounded bg-muted animate-pulse ml-auto" />
                <div className="h-4 w-48 rounded bg-muted animate-pulse ml-auto" />
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
                  Изменение за день: 
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
      </div>
    </Card>
  )
}

