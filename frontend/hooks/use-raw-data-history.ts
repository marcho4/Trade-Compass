"use client"

import { useEffect, useState } from "react"
import { financialDataApi } from "@/lib/api"
import { RawData } from "@/types/raw-data"

interface RawDataHistoryState {
  data: RawData[]
  loading: boolean
  error: string | null
}

export const useRawDataHistory = (ticker: string): RawDataHistoryState => {
  const [state, setState] = useState<RawDataHistoryState>({
    data: [],
    loading: true,
    error: null,
  })

  useEffect(() => {
    let cancelled = false

    const fetch = async () => {
      setState((prev) => ({ ...prev, loading: true, error: null }))
      try {
        const history = await financialDataApi.getRawDataHistory(ticker)
        if (!cancelled) {
          setState({ data: history, loading: false, error: null })
        }
      } catch (err) {
        if (!cancelled) {
          setState({
            data: [],
            loading: false,
            error: err instanceof Error ? err.message : "Ошибка загрузки данных",
          })
        }
      }
    }

    fetch()
    return () => {
      cancelled = true
    }
  }, [ticker])

  return state
}
