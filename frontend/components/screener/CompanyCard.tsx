"use client"

import { cn } from "@/lib/utils"
import type { CompanyRating } from "./types"

function fmtPrice(p: number): string {
  return p.toLocaleString("ru-RU", { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function fmtCap(n: number): string {
  if (n >= 1e12) return `₽${(n / 1e12).toFixed(1)}T`
  if (n >= 1e9) return `₽${(n / 1e9).toFixed(1)}B`
  if (n >= 1e6) return `₽${(n / 1e6).toFixed(0)}M`
  return `₽${n.toFixed(0)}`
}

const RADAR_LABELS = ["Здоровье", "Рост", "Ров", "Дивиденды", "Оценка"]

function RadarViz({ scores, size = 220 }: { scores: number[]; size?: number }) {
  const c = size / 2
  const r = c - 38
  const n = scores.length
  const rings = [0.25, 0.5, 0.75, 1.0]

  const pt = (i: number, mag: number): [number, number] => {
    const a = -Math.PI / 2 + (i * 2 * Math.PI) / n
    return [c + Math.cos(a) * r * mag, c + Math.sin(a) * r * mag]
  }

  const poly = (vals: number[]): string =>
    vals.map((v, i) => {
      const [x, y] = pt(i, v)
      return `${i === 0 ? "M" : "L"}${x.toFixed(1)},${y.toFixed(1)}`
    }).join(" ") + " Z"

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} className="block overflow-visible">
      {rings.map((m, i) => (
        <path
          key={i}
          d={poly(Array(n).fill(m))}
          fill="none"
          className="stroke-border"
          strokeWidth={1}
          strokeDasharray={i < rings.length - 1 ? "2 3" : undefined}
        />
      ))}

      {scores.map((_, i) => {
        const [x, y] = pt(i, 1)
        return (
          <line
            key={i}
            x1={c} y1={c}
            x2={x.toFixed(1)} y2={y.toFixed(1)}
            className="stroke-border"
            strokeWidth={1}
            strokeDasharray="2 3"
          />
        )
      })}

      <path d={poly(scores)} className="fill-primary/20 stroke-primary" strokeWidth={1.4} />

      {scores.map((v, i) => {
        const [x, y] = pt(i, v)
        return <rect key={i} x={x - 2} y={y - 2} width={4} height={4} className="fill-primary" />
      })}

      {RADAR_LABELS.map((label, i) => {
        const [x, y] = pt(i, 1.22)
        const a = -Math.PI / 2 + (i * 2 * Math.PI) / n
        const cos = Math.cos(a)
        const anchor = Math.abs(cos) < 0.2 ? "middle" : cos > 0 ? "start" : "end"
        const dy = Math.sin(a) > 0.3 ? 9 : Math.sin(a) < -0.3 ? -2 : 4
        return (
          <text
            key={i}
            x={x.toFixed(1)}
            y={y.toFixed(1)}
            textAnchor={anchor}
            dy={dy}
            className="font-mono text-[9px] font-bold fill-muted-foreground tracking-[1px]"
          >
            {label}
          </text>
        )
      })}
    </svg>
  )
}

function RatingDots({ value, max = 5 }: { value: number; max?: number }) {
  return (
    <div className="flex items-center gap-[3px]">
      {Array.from({ length: max }).map((_, i) => (
        <span
          key={i}
          className={cn(
            "w-[10px] h-[10px] rounded-[1px] border",
            i < value ? "bg-primary border-primary" : "border-border",
          )}
        />
      ))}
    </div>
  )
}

interface CompanyCardProps {
  id: number
  ticker: string
  name: string
  sector: string
  price: number
  priceChange: number
  priceChangePercent: number
  priceLoading?: boolean
  rating: CompanyRating
  marketCap?: number
  onClick?: () => void
}

export const CompanyCard = ({
  ticker,
  name,
  sector,
  price,
  priceChange,
  priceChangePercent,
  priceLoading,
  rating,
  marketCap,
  onClick,
}: CompanyCardProps) => {
  const isPos = priceChangePercent > 0
  const isNeg = priceChangePercent < 0
  const arrow = isPos ? "↑" : isNeg ? "↓" : "→"
  const changeClass = isPos ? "text-positive" : isNeg ? "text-negative" : "text-muted-foreground"

  const scores = [
    rating.health / 6,
    rating.growth / 6,
    rating.moat / 6,
    rating.dividends / 6,
    rating.value / 6,
  ]

  const totalDots = Math.min(5, Math.max(1, Math.round(rating.total)))

  const footer: [string, string][] = [
    ["Кап.", marketCap ? fmtCap(marketCap) : "—"],
    ["Рейтинг", `${totalDots}/5`],
    ["Здоровье", `${Math.round((rating.health / 6) * 100)}`],
    ["Рост", `${Math.round((rating.growth / 6) * 100)}`],
  ]

  return (
    <div
      className="flex flex-col bg-card border border-border rounded-[2px] shadow-[0_1px_0_rgba(20,20,20,0.02)] cursor-pointer hover:shadow-md transition-shadow duration-150"
      onClick={onClick}
      onKeyDown={(e) => { if (e.key === "Enter" || e.key === " ") onClick?.() }}
      role="button"
      tabIndex={0}
      aria-label={`Открыть анализ компании ${name}`}
    >
      {/* Header strip */}
      <div className="flex items-center justify-between px-[14px] py-[10px] bg-muted/40 border-b border-border">
        <div className="flex items-center gap-[10px]">
          <span className="w-[3px] h-3 bg-primary shrink-0" />
          <span className="font-mono text-[11px] font-bold text-foreground tracking-[1.2px]">
            {ticker}
            <span className="text-muted-foreground/50">.MOEX</span>
          </span>
          <span className="inline-flex items-center font-mono text-[10px] font-semibold tracking-[0.8px] uppercase text-foreground bg-card border border-border px-[7px] py-[3px] rounded-[2px] leading-[1.4]">
            {sector}
          </span>
        </div>
      </div>

      {/* Identity + Rating */}
      <div className="px-4 pt-[14px] pb-[10px] border-b border-dashed border-border">
        <div className="flex justify-between items-start gap-3">
          <div>
            <div className="font-mono text-[32px] font-bold text-foreground leading-none tracking-[-1px]">
              {ticker}
            </div>
            <div className="text-[13px] text-foreground/70 mt-[6px] leading-[1.3] max-w-[200px]">
              {name}
            </div>
          </div>
          <div className="text-right shrink-0">
            <div className="font-mono text-[9px] font-bold text-muted-foreground tracking-[1px] uppercase mb-[5px]">
              Рейтинг
            </div>
            <div className="font-mono text-[22px] font-bold text-foreground leading-none">
              {totalDots}
              <span className="text-muted-foreground/50 text-[12px]">/5</span>
            </div>
            <div className="mt-[6px]">
              <RatingDots value={totalDots} />
            </div>
          </div>
        </div>
      </div>

      {/* Radar */}
      <div className="flex justify-center px-2 py-[6px] border-b border-dashed border-border">
        <RadarViz scores={scores} size={220} />
      </div>

      {/* Price row */}
      <div className="flex-1 px-4 pt-3 pb-[10px]">
        <div className="font-mono text-[9px] font-bold text-muted-foreground tracking-[1px] uppercase mb-1">
          Цена
        </div>
        {priceLoading ? (
          <div className="h-6 w-24 bg-muted animate-pulse" />
        ) : (
          <>
            <div className="font-mono text-[26px] font-bold text-foreground leading-none tracking-[-0.4px]">
              {fmtPrice(price)}{" "}
              <span className="text-muted-foreground text-[14px] font-semibold">₽</span>
            </div>
            <div className={cn("font-mono text-[11px] font-semibold mt-[6px]", changeClass)}>
              {arrow} {priceChangePercent >= 0 ? "+" : ""}
              {priceChangePercent.toFixed(2)}%{" "}
              <span className="text-muted-foreground font-normal">
                ({priceChange >= 0 ? "+" : ""}
                {priceChange.toFixed(2)} ₽)
              </span>
            </div>
          </>
        )}
      </div>

      {/* Footer KV */}
      <div className="grid grid-cols-4 border-t border-border bg-muted/30">
        {footer.map(([k, v], i) => (
          <div
            key={k}
            className={cn(
              "px-3 py-[10px]",
              i < footer.length - 1 && "border-r border-border",
            )}
          >
            <div className="font-mono text-[9px] font-bold text-muted-foreground tracking-[1px] uppercase">
              {k}
            </div>
            <div className="font-mono text-[12px] font-semibold text-foreground mt-[3px]">
              {v}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
