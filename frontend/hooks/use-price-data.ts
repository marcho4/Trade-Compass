"use client"

import { useEffect, useState } from "react"
import { financialDataApi } from "@/lib/api"

export interface PriceData {
  price: number
  previousPrice: number
  priceChange: number
  priceChangePercent: number
  loading: boolean
}

const INITIAL_PRICE_DATA: PriceData = {
  price: 0,
  previousPrice: 0,
  priceChange: 0,
  priceChangePercent: 0,
  loading: true,
}

export const usePriceData = (ticker: string): PriceData => {
  const [data, setData] = useState<PriceData>(INITIAL_PRICE_DATA)

  useEffect(() => {
    const controller = new AbortController()

    const fetchPrice = async () => {
      try {
        const candles = await financialDataApi.getPriceCandles(
          ticker,
          2,
          24,
          controller.signal,
        )

        if (candles.length === 0) {
          setData((prev) => ({ ...prev, loading: false }))
          return
        }

        const currentCandle = candles[candles.length - 1]
        const price = currentCandle.close

        let previousPrice = price
        if (candles.length >= 2) {
          previousPrice = candles[candles.length - 2].close
        }

        const priceChange = price - previousPrice
        const priceChangePercent =
          previousPrice !== 0 ? (priceChange / previousPrice) * 100 : 0

        setData({
          price,
          previousPrice,
          priceChange,
          priceChangePercent,
          loading: false,
        })
      } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") return
        console.error(`Failed to fetch price data for ${ticker}:`, error)
        setData((prev) => ({ ...prev, loading: false }))
      }
    }

    fetchPrice()
    return () => controller.abort()
  }, [ticker])

  return data
}
