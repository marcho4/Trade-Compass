"use client"

import { useEffect, useState } from "react"
import { financialDataApi } from "@/lib/api"
import type { Candle } from "@/lib/api/financial-data-api"

export interface TickerScreenData {
  price: number
  prevClose: number
  open: number
  dayLow: number
  dayHigh: number
  change: number
  changePct: number
  volumeShares: number
  turnoverRub: number
  w52Low: number
  w52High: number
  marketCapRub: number | null
  loading: boolean
}

const INITIAL: TickerScreenData = {
  price: 0,
  prevClose: 0,
  open: 0,
  dayLow: 0,
  dayHigh: 0,
  change: 0,
  changePct: 0,
  volumeShares: 0,
  turnoverRub: 0,
  w52Low: 0,
  w52High: 0,
  marketCapRub: null,
  loading: true,
}

function aggregateLast(candles: Candle[]): {
  open: number
  close: number
  high: number
  low: number
  volume: number
  turnover: number
  prevClose: number
} | null {
  if (candles.length === 0) return null

  const lastDate = candles[candles.length - 1].begin.slice(0, 10)
  const today = candles.filter((c) => c.begin.startsWith(lastDate))
  if (today.length === 0) return null

  const open = today[0].open
  const close = today[today.length - 1].close
  const high = Math.max(...today.map((c) => c.high))
  const low = Math.min(...today.map((c) => c.low))
  const volume = today.reduce((s, c) => s + c.volume, 0)
  const turnover = today.reduce((s, c) => s + c.value, 0)

  const beforeToday = candles.filter((c) => !c.begin.startsWith(lastDate))
  const prevClose = beforeToday.length > 0 ? beforeToday[beforeToday.length - 1].close : open

  return { open, close, high, low, volume, turnover, prevClose }
}

export function useTickerScreenData(ticker: string): TickerScreenData {
  const [data, setData] = useState<TickerScreenData>(INITIAL)

  useEffect(() => {
    const controller = new AbortController()

    const run = async () => {
      try {
        const [intraday, yearly, marketCap] = await Promise.allSettled([
          financialDataApi.getPriceCandles(ticker, 2, 1, controller.signal),
          financialDataApi.getPriceCandles(ticker, 365, 24, controller.signal),
          financialDataApi.getMarketCap(ticker, controller.signal),
        ])

        let today: ReturnType<typeof aggregateLast> = null
        if (intraday.status === "fulfilled") {
          today = aggregateLast(intraday.value)
        }

        let w52Low = 0
        let w52High = 0
        if (yearly.status === "fulfilled" && yearly.value.length > 0) {
          w52Low = Math.min(...yearly.value.map((c) => c.low))
          w52High = Math.max(...yearly.value.map((c) => c.high))
        }

        let price = 0
        let prevClose = 0
        let open = 0
        let dayLow = 0
        let dayHigh = 0
        let volumeShares = 0
        let turnoverRub = 0

        if (today) {
          price = today.close
          prevClose = today.prevClose
          open = today.open
          dayLow = today.low
          dayHigh = today.high
          volumeShares = today.volume
          turnoverRub = today.turnover
        } else if (yearly.status === "fulfilled" && yearly.value.length > 0) {
          const last = yearly.value[yearly.value.length - 1]
          price = last.close
          prevClose = last.open
          open = last.open
          dayLow = last.low
          dayHigh = last.high
          volumeShares = last.volume
          turnoverRub = last.value
        }

        const change = price - prevClose
        const changePct = prevClose > 0 ? (change / prevClose) * 100 : 0

        setData({
          price,
          prevClose,
          open,
          dayLow,
          dayHigh,
          change,
          changePct,
          volumeShares,
          turnoverRub,
          w52Low,
          w52High,
          marketCapRub: marketCap.status === "fulfilled" ? marketCap.value : null,
          loading: false,
        })
      } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") return
        console.error(`Failed to fetch ticker screen data for ${ticker}:`, error)
        setData((prev) => ({ ...prev, loading: false }))
      }
    }

    run()
    return () => controller.abort()
  }, [ticker])

  return data
}

export function formatShortNumber(value: number, currency = false): string {
  if (!value) return currency ? "—" : "0"
  const abs = Math.abs(value)
  const sign = value < 0 ? "-" : ""
  const prefix = currency ? "₽" : ""

  if (abs >= 1e12) return `${sign}${prefix}${(abs / 1e12).toFixed(2)}T`
  if (abs >= 1e9) return `${sign}${prefix}${(abs / 1e9).toFixed(2)}B`
  if (abs >= 1e6) return `${sign}${prefix}${(abs / 1e6).toFixed(2)}M`
  if (abs >= 1e3) return `${sign}${prefix}${(abs / 1e3).toFixed(1)}K`
  return `${sign}${prefix}${abs.toFixed(0)}`
}
