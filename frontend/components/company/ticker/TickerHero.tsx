"use client"

import { T } from "./tokens"
import { Chip, RangeBar, Stat } from "./primitives"
import { formatShortNumber } from "@/hooks/use-ticker-screen-data"

export interface TickerHeroData {
  symbol: string
  exchange: string
  currency: string
  isin: string
  name: string
  sector: string
  industry: string
  description: string
  price: number
  change: number
  changePct: number
  lastUpdate: string
  w52Low: number
  w52High: number
  open: number
  prevClose: number
  dayLow: number
  dayHigh: number
  volumeShares: number
  turnoverRub: number
  marketCap: string
  freeFloat: string
  pe: number | null
  pb: number | null
  divYield: string
  beta: number | null
  loading?: boolean
}

interface TickerHeroProps {
  data: TickerHeroData
}

function formatNumber(value: number, digits = 2): string {
  return new Intl.NumberFormat("ru-RU", {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  }).format(value)
}

export function TickerHero({ data }: TickerHeroProps) {
  const isPos = data.changePct > 0
  const isNeg = data.changePct < 0
  const changeColor = isPos ? T.pos : isNeg ? T.neg : T.textDim
  const changeArrow = isPos ? "↑" : isNeg ? "↓" : "→"

  const displayName = data.name || (data.loading ? "Загрузка..." : data.symbol)
  const displayDescription =
    data.description ||
    (data.loading
      ? "Загрузка описания компании..."
      : "Описание компании пока недоступно.")

  return (
    <div
      style={{
        background: T.panel,
        border: `1px solid ${T.border}`,
        borderRadius: 2,
        boxShadow: "0 1px 0 rgba(20,20,20,0.02)",
        display: "grid",
        gridTemplateColumns: "minmax(0, 1.2fr) minmax(0, 1fr)",
      }}
    >
      <div style={{ padding: "20px 24px", borderRight: `1px solid ${T.border}` }}>
        <div
          style={{
            display: "flex",
            alignItems: "baseline",
            gap: 14,
            marginBottom: 14,
          }}
        >
          <div
            style={{
              fontFamily: T.mono,
              fontSize: 48,
              fontWeight: 700,
              color: T.text,
              lineHeight: 1,
              letterSpacing: -1,
            }}
          >
            {data.symbol}
          </div>
          <div
            style={{
              fontFamily: T.mono,
              fontSize: 11,
              color: T.textDim,
              letterSpacing: 0.8,
            }}
          >
            {data.exchange} · {data.currency} · ISIN {data.isin}
          </div>
        </div>

        <div
          style={{
            fontFamily: T.sans,
            fontSize: 22,
            fontWeight: 600,
            color: T.text,
            letterSpacing: -0.4,
            marginBottom: 8,
            lineHeight: 1.15,
          }}
        >
          {displayName}
        </div>

        <div style={{ display: "flex", flexWrap: "wrap", gap: 6, marginBottom: 14 }}>
          {data.sector && (
            <Chip color={T.accent} bg={T.accentSoft} bd={T.accentLine}>
              {data.sector}
            </Chip>
          )}
          {data.industry && <Chip>{data.industry}</Chip>}
        </div>

        <div
          style={{
            fontFamily: T.sans,
            fontSize: 12.5,
            lineHeight: 1.55,
            color: T.text2,
            maxWidth: 620,
            textWrap: "pretty" as never,
          }}
        >
          {displayDescription}
        </div>
      </div>

      <div style={{ padding: "20px 24px", minWidth: 0 }}>
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "baseline",
            marginBottom: 6,
          }}
        >
          <span
            style={{
              fontFamily: T.mono,
              fontSize: 9.5,
              fontWeight: 600,
              letterSpacing: 1,
              textTransform: "uppercase",
              color: T.textDim,
            }}
          >
            Последняя цена
          </span>
          <span style={{ fontFamily: T.mono, fontSize: 10, color: T.textFaint }}>
            {data.lastUpdate}
          </span>
        </div>

        <div
          style={{
            display: "flex",
            alignItems: "baseline",
            gap: 14,
            marginBottom: 4,
            flexWrap: "wrap",
          }}
        >
          <div
            style={{
              fontFamily: T.mono,
              fontSize: 40,
              fontWeight: 700,
              color: T.text,
              letterSpacing: -0.6,
              lineHeight: 1,
            }}
          >
            {data.price > 0 ? `₽${formatNumber(data.price)}` : "—"}
          </div>
          <div
            style={{
              fontFamily: T.mono,
              fontSize: 13,
              fontWeight: 600,
              color: changeColor,
              display: "flex",
              alignItems: "baseline",
              gap: 6,
            }}
          >
            <span>
              {changeArrow} {data.changePct >= 0 ? "+" : ""}
              {data.changePct.toFixed(2)}%
            </span>
            <span style={{ color: T.textDim, fontWeight: 500 }}>
              ({data.change >= 0 ? "+" : ""}
              {data.change.toFixed(2)} ₽)
            </span>
          </div>
        </div>

        {data.w52Low > 0 && data.w52High > 0 && (
          <div style={{ marginTop: 14, marginBottom: 16 }}>
            <RangeBar
              low={data.w52Low}
              high={data.w52High}
              value={data.price > 0 ? data.price : data.w52Low}
              width={360}
            />
          </div>
        )}

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(4, minmax(0, 1fr))",
            gap: "14px 18px",
            paddingTop: 12,
            borderTop: `1px solid ${T.borderSoft}`,
          }}
        >
          <Stat label="Открытие" value={data.open > 0 ? formatNumber(data.open) : "—"} />
          <Stat
            label="Пред. закр."
            value={data.prevClose > 0 ? formatNumber(data.prevClose) : "—"}
          />
          <Stat label="Мин. дня" value={data.dayLow > 0 ? formatNumber(data.dayLow) : "—"} />
          <Stat label="Макс. дня" value={data.dayHigh > 0 ? formatNumber(data.dayHigh) : "—"} />
          <Stat
            label="Объём"
            value={data.volumeShares > 0 ? formatShortNumber(data.volumeShares) : "—"}
          />
          <Stat
            label="Оборот"
            value={data.turnoverRub > 0 ? formatShortNumber(data.turnoverRub, true) : "—"}
          />
          <Stat label="Кап." value={data.marketCap} />
          <Stat label="Free float" value={data.freeFloat} />
          <Stat label="P/E" value={data.pe != null ? data.pe.toFixed(1) : "—"} />
          <Stat label="P/B" value={data.pb != null ? data.pb.toFixed(1) : "—"} />
          <Stat label="Див. доходн." value={data.divYield} />
          <Stat label="Beta" value={data.beta != null ? data.beta.toFixed(2) : "—"} />
        </div>
      </div>
    </div>
  )
}
