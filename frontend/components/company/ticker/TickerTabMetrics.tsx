"use client"

import { useMemo } from "react"
import { useRawDataHistory } from "@/hooks/use-raw-data-history"
import { T } from "./tokens"
import { TPanel } from "./primitives"
import { RawData } from "@/types/raw-data"
import { formatShortNumber } from "@/hooks/use-ticker-screen-data"

interface TickerTabMetricsProps {
  ticker: string
}

interface MetricTile {
  label: string
  value: string
  delta: string | null
  positive: boolean
}

const UNITS_MAP: Record<string, number> = {
  units: 1,
  thousands: 1000,
  millions: 1000000,
}

function multiplier(row: RawData): number {
  return (row.reportUnits && UNITS_MAP[row.reportUnits]) || 1
}

function num(row: RawData, field: keyof RawData): number | null {
  const raw = row[field] as number | null | undefined
  if (raw == null) return null
  return raw * multiplier(row)
}

function formatDelta(curr: number, prev: number, percentPoints = false): { label: string; positive: boolean } {
  if (prev === 0) return { label: "—", positive: true }
  if (percentPoints) {
    const delta = curr - prev
    const sign = delta > 0 ? "+" : delta < 0 ? "" : ""
    return {
      label: `${sign}${delta.toFixed(1)} п.п.`,
      positive: delta >= 0,
    }
  }
  const pct = ((curr - prev) / Math.abs(prev)) * 100
  const sign = pct > 0 ? "+" : pct < 0 ? "" : ""
  return {
    label: `${sign}${pct.toFixed(1)}%`,
    positive: pct >= 0,
  }
}

function buildTiles(data: RawData[]): { tiles: MetricTile[]; latestYear: number | null } {
  const yearEntries = data
    .filter((r) => r.period === "YEAR" && r.status === "confirmed")
    .sort((a, b) => a.year - b.year)

  if (yearEntries.length === 0) {
    return { tiles: [], latestYear: null }
  }

  const latest = yearEntries[yearEntries.length - 1]
  const prev = yearEntries.length > 1 ? yearEntries[yearEntries.length - 2] : null

  const revenueCurr = num(latest, "revenue")
  const revenuePrev = prev ? num(prev, "revenue") : null
  const netCurr = num(latest, "netProfit")
  const netPrev = prev ? num(prev, "netProfit") : null
  const ebitdaCurr = num(latest, "ebitda")
  const ebitdaPrev = prev ? num(prev, "ebitda") : null
  const ocfCurr = num(latest, "operatingCashFlow")
  const ocfPrev = prev ? num(prev, "operatingCashFlow") : null
  const fcfCurr = num(latest, "freeCashFlow")
  const fcfPrev = prev ? num(prev, "freeCashFlow") : null
  const equityCurr = num(latest, "equity")
  const equityPrev = prev ? num(prev, "equity") : null
  const assetsCurr = num(latest, "totalAssets")
  const assetsPrev = prev ? num(prev, "totalAssets") : null
  const debtCurr = num(latest, "debt")
  const debtPrev = prev ? num(prev, "debt") : null

  const roeCurr = netCurr != null && equityCurr ? (netCurr / equityCurr) * 100 : null
  const roePrev = netPrev != null && equityPrev ? (netPrev / equityPrev) * 100 : null
  const roaCurr = netCurr != null && assetsCurr ? (netCurr / assetsCurr) * 100 : null
  const roaPrev = netPrev != null && assetsPrev ? (netPrev / assetsPrev) * 100 : null
  const deCurr = debtCurr != null && equityCurr ? debtCurr / equityCurr : null
  const dePrev = debtPrev != null && equityPrev ? debtPrev / equityPrev : null

  const tiles: MetricTile[] = []

  const push = (
    label: string,
    curr: number | null,
    prevVal: number | null,
    formatter: (v: number) => string,
    percentPoints = false,
  ) => {
    if (curr == null) return
    const d = prevVal != null ? formatDelta(curr, prevVal, percentPoints) : null
    tiles.push({
      label,
      value: formatter(curr),
      delta: d?.label || null,
      positive: d?.positive ?? true,
    })
  }

  push("Выручка", revenueCurr, revenuePrev, (v) => formatShortNumber(v, true))
  push("Чистая прибыль", netCurr, netPrev, (v) => formatShortNumber(v, true))
  push("EBITDA", ebitdaCurr, ebitdaPrev, (v) => formatShortNumber(v, true))
  push("Оп. ден. поток", ocfCurr, ocfPrev, (v) => formatShortNumber(v, true))
  push("FCF", fcfCurr, fcfPrev, (v) => formatShortNumber(v, true))
  push("ROE", roeCurr, roePrev, (v) => `${v.toFixed(1)}%`, true)
  push("ROA", roaCurr, roaPrev, (v) => `${v.toFixed(1)}%`, true)
  push("Долг / Капитал", deCurr, dePrev, (v) => v.toFixed(2))

  return { tiles, latestYear: latest.year }
}

export function TickerTabMetrics({ ticker }: TickerTabMetricsProps) {
  const { data, loading, error } = useRawDataHistory(ticker)

  const { tiles, latestYear } = useMemo(() => buildTiles(data), [data])

  if (loading) {
    return (
      <div
        style={{
          padding: 40,
          textAlign: "center",
          fontFamily: T.mono,
          color: T.textDim,
          fontSize: 11,
          letterSpacing: 1,
          textTransform: "uppercase",
          border: `1px solid ${T.border}`,
          background: T.panel,
        }}
      >
        Загрузка показателей...
      </div>
    )
  }

  if (error) {
    return (
      <div
        style={{
          padding: 40,
          textAlign: "center",
          fontFamily: T.sans,
          color: T.neg,
          fontSize: 13,
          border: `1px solid ${T.border}`,
          background: T.panel,
        }}
      >
        {error}
      </div>
    )
  }

  if (tiles.length === 0) {
    return (
      <TPanel title="Ключевые показатели" accent={T.accent}>
        <div
          style={{
            padding: 40,
            textAlign: "center",
            fontFamily: T.sans,
            color: T.textDim,
            fontSize: 13,
          }}
        >
          Нет данных для расчёта показателей
        </div>
      </TPanel>
    )
  }

  return (
    <TPanel
      title="Ключевые показатели"
      accent={T.accent}
      right={
        <span
          style={{
            fontFamily: T.mono,
            fontSize: 10,
            color: T.textDim,
            letterSpacing: 0.5,
          }}
        >
          FY {latestYear ?? "—"} · Δ К FY {latestYear ? latestYear - 1 : "—"}
        </span>
      }
    >
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(4, minmax(0, 1fr))",
        }}
      >
        {tiles.map((m, i) => (
          <div
            key={m.label}
            style={{
              padding: "20px 22px",
              borderRight: i % 4 < 3 ? `1px solid ${T.borderSoft}` : "none",
              borderBottom: i < tiles.length - 4 ? `1px solid ${T.borderSoft}` : "none",
              minWidth: 0,
            }}
          >
            <div
              style={{
                fontFamily: T.mono,
                fontSize: 9.5,
                color: T.textDim,
                letterSpacing: 1,
                textTransform: "uppercase",
                marginBottom: 8,
              }}
            >
              {m.label}
            </div>
            <div
              style={{
                fontFamily: T.mono,
                fontSize: 22,
                fontWeight: 700,
                color: T.text,
                letterSpacing: -0.3,
                marginBottom: 6,
                lineHeight: 1,
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
              }}
              title={m.value}
            >
              {m.value}
            </div>
            {m.delta && (
              <div
                style={{
                  fontFamily: T.mono,
                  fontSize: 11,
                  fontWeight: 600,
                  color: m.positive ? T.pos : T.neg,
                  letterSpacing: 0.2,
                  display: "inline-flex",
                  alignItems: "center",
                  gap: 4,
                }}
              >
                <span>{m.positive ? "↑" : "↓"}</span>
                <span>{m.delta}</span>
              </div>
            )}
          </div>
        ))}
      </div>
    </TPanel>
  )
}
