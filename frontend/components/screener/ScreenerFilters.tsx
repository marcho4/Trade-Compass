"use client"

import { useState } from "react"
import type { FilterValues } from "./types"

const SORT_OPTIONS = [
  { k: "rating", label: "Рейтинг" },
  { k: "name", label: "Тикер" },
  { k: "cap", label: "Капитализация" },
]

interface ScreenerFiltersProps {
  filters: FilterValues
  onFilterChange: (key: keyof FilterValues, value: string) => void
  onReset: () => void
  sectors: Array<{ id: number; name: string }>
  sort: string
  sortDir: "asc" | "desc"
  onSortCycle: () => void
  found: number
  total: number
}

function RatingDiamond({ filled }: { filled: boolean }) {
  return (
    <span
      className={[
        "inline-block w-[11px] h-[11px] rotate-45 flex-shrink-0",
        filled
          ? "bg-primary border border-primary"
          : "bg-transparent border border-border",
      ].join(" ")}
    />
  )
}

function RatingPicker({
  value,
  onChange,
}: {
  value: number
  onChange: (v: number) => void
}) {
  const [hover, setHover] = useState(0)
  const active = hover || value

  return (
    <div className="flex items-center gap-1.5" onMouseLeave={() => setHover(0)}>
      <button
        onClick={() => onChange(0)}
        className={[
          "font-mono text-[10px] font-bold tracking-[0.8px] uppercase px-2 py-1 rounded-[2px] leading-none cursor-pointer border transition-colors",
          value === 0
            ? "bg-foreground text-background border-foreground"
            : "bg-transparent text-muted-foreground border-border",
        ].join(" ")}
      >
        ЛЮБ
      </button>
      <div className="flex items-center gap-[3px]">
        {[1, 2, 3, 4, 5].map((i) => (
          <button
            key={i}
            onMouseEnter={() => setHover(i)}
            onClick={() => onChange(i === value ? 0 : i)}
            className="p-[2px] cursor-pointer bg-transparent border-none leading-none"
          >
            <RatingDiamond filled={i <= active} />
          </button>
        ))}
      </div>
      {value > 0 && (
        <span className="font-mono text-[10px] text-muted-foreground tracking-[0.8px]">
          {value}+
        </span>
      )}
    </div>
  )
}

function SortButton({
  sort,
  dir,
  onCycle,
}: {
  sort: string
  dir: "asc" | "desc"
  onCycle: () => void
}) {
  const label = SORT_OPTIONS.find((s) => s.k === sort)?.label ?? "—"
  return (
    <button
      onClick={onCycle}
      className="flex items-center gap-2 px-2.5 h-8 bg-card border border-border rounded-[2px] font-mono text-[11px] font-semibold tracking-[0.4px] cursor-pointer whitespace-nowrap flex-shrink-0 hover:bg-muted/50 transition-colors"
    >
      <span className="text-[9px] text-muted-foreground tracking-[1px] uppercase font-bold">
        сорт
      </span>
      <span className="text-foreground">{label}</span>
      <span className="text-primary font-bold">{dir === "asc" ? "↑" : "↓"}</span>
    </button>
  )
}

function ActiveChips({
  filters,
  sectors,
  onClear,
  onClearAll,
}: {
  filters: FilterValues
  sectors: Array<{ id: number; name: string }>
  onClear: (k: keyof FilterValues) => void
  onClearAll: () => void
}) {
  const chips: { k: keyof FilterValues; label: string }[] = []

  if (filters.search)
    chips.push({ k: "search", label: `ПОИСК · ${filters.search}` })
  if (filters.sector) {
    const s = sectors.find((x) => x.id.toString() === filters.sector)
    if (s)
      chips.push({ k: "sector", label: `СЕКТОР · ${s.name.toUpperCase()}` })
  }
  if (filters.ratingMin)
    chips.push({ k: "ratingMin", label: `РЕЙТИНГ · ${filters.ratingMin}+` })

  if (chips.length === 0) return null

  return (
    <div className="flex items-center gap-1.5 flex-wrap">
      {chips.map((c) => (
        <span
          key={c.k}
          className="inline-flex items-center gap-1.5 font-mono text-[10px] font-bold tracking-[0.8px] text-foreground bg-primary/10 border border-primary px-2 py-[3px] rounded-[2px]"
        >
          <span className="w-[5px] h-[5px] bg-primary flex-shrink-0" />
          {c.label}
          <button
            onClick={() => onClear(c.k)}
            className="border-none bg-transparent cursor-pointer font-mono text-[12px] text-muted-foreground font-bold leading-none ml-0.5 px-[2px]"
          >
            ×
          </button>
        </span>
      ))}
      <button
        onClick={onClearAll}
        className="bg-transparent border border-dashed border-border px-2 py-[3px] rounded-[2px] cursor-pointer font-mono text-[10px] font-bold text-muted-foreground tracking-[1px] uppercase hover:border-muted-foreground transition-colors"
      >
        сброс
      </button>
    </div>
  )
}

export function ScreenerFilters({
  filters,
  onFilterChange,
  onReset,
  sectors,
  sort,
  sortDir,
  onSortCycle,
  found,
  total,
}: ScreenerFiltersProps) {
  const [advOpen, setAdvOpen] = useState(false)

  const activeCount =
    (filters.search ? 1 : 0) +
    (filters.sector ? 1 : 0) +
    (filters.ratingMin ? 1 : 0)

  const allSectors = [{ id: 0, name: "Все" }, ...sectors]

  return (
    <div className="bg-card border border-border rounded-[2px] overflow-hidden shadow-[0_1px_0_rgba(20,20,20,0.02)]">
      {/* Title strip */}
      <div className="flex items-center justify-between px-4 py-2.5 border-b border-border bg-muted/30 gap-4">
        <div className="flex items-center gap-2.5">
          <span className="w-[3px] h-[14px] bg-primary inline-block flex-shrink-0" />
          <span className="font-mono text-[12px] font-bold tracking-[1.2px] uppercase text-foreground">
            Скринер акций
          </span>
          {total > 0 && (
            <span className="font-mono text-[10px] text-muted-foreground/60 tracking-[0.5px]">
              · {total} бумаг MOEX
            </span>
          )}
        </div>
        <span className="font-mono text-[10px] text-muted-foreground tracking-[0.6px] uppercase whitespace-nowrap">
          найдено{" "}
          <span className="text-foreground font-bold text-[12px]">
            {String(found).padStart(3, "0")}
          </span>
        </span>
      </div>

      {/* Search + filters toggle + sort */}
      <div className="px-4 py-2.5 flex items-center gap-2 flex-wrap">
        {/* Search */}
        <div className="flex items-center gap-2 border border-border bg-card px-2.5 rounded-[2px] h-8 min-w-[260px] flex-1 max-w-[400px]">
          <span className="font-mono text-[11px] font-bold text-primary">
            &gt;_
          </span>
          <input
            value={filters.search}
            onChange={(e) => onFilterChange("search", e.target.value)}
            placeholder="Тикер или название…"
            className="flex-1 border-none outline-none bg-transparent font-mono text-[11.5px] font-medium text-foreground tracking-[0.2px] placeholder:text-muted-foreground/60"
          />
          {filters.search && (
            <button
              onClick={() => onFilterChange("search", "")}
              className="border-none bg-transparent cursor-pointer font-mono text-[13px] text-muted-foreground font-bold leading-none"
            >
              ×
            </button>
          )}
        </div>

        {/* Filters toggle */}
        <button
          onClick={() => setAdvOpen((v) => !v)}
          className={[
            "flex items-center gap-2 px-3 h-8 rounded-[2px] font-mono text-[11px] font-bold tracking-[0.8px] uppercase cursor-pointer border transition-colors flex-shrink-0",
            advOpen
              ? "bg-foreground text-background border-foreground"
              : "bg-card text-foreground border-border hover:bg-muted/50",
          ].join(" ")}
        >
          <span className="flex flex-col gap-[2.5px] items-end">
            <span className="block w-[11px] h-[1px] bg-current" />
            <span className="block w-[7px] h-[1px] bg-current" />
            <span className="block w-[4px] h-[1px] bg-current" />
          </span>
          Фильтры
          {activeCount > 0 && (
            <span
              className={[
                "font-mono text-[10px] text-white px-[5px] py-[1px] font-bold rounded-[2px]",
                advOpen ? "bg-white/20" : "bg-primary",
              ].join(" ")}
            >
              {activeCount}
            </span>
          )}
          <span className="text-[9px] opacity-70">{advOpen ? "▲" : "▼"}</span>
        </button>

        <div className="flex-1" />
        <SortButton sort={sort} dir={sortDir} onCycle={onSortCycle} />
      </div>

      {/* Expandable drawer */}
      {advOpen && (
        <div className="px-4 pb-4 pt-3.5 border-t border-dashed border-border/60 bg-muted/20 grid grid-cols-[1.6fr_1fr] gap-6">
          {/* Sector chips */}
          <div>
            <p className="font-mono text-[10px] font-bold text-muted-foreground tracking-[1px] uppercase mb-2">
              Сектор
            </p>
            <div className="flex flex-wrap gap-1">
              {allSectors.map((s) => {
                const val = s.id === 0 ? "" : s.id.toString()
                const on = filters.sector === val
                return (
                  <button
                    key={s.id}
                    onClick={() => onFilterChange("sector", val)}
                    className={[
                      "font-mono text-[10.5px] font-bold tracking-[0.6px] uppercase px-[9px] py-1 rounded-[2px] cursor-pointer border transition-colors",
                      on
                        ? "bg-primary text-white border-primary"
                        : "bg-card text-foreground border-border hover:bg-muted/50",
                    ].join(" ")}
                  >
                    {s.name}
                  </button>
                )
              })}
            </div>
          </div>

          {/* Rating */}
          <div>
            <p className="font-mono text-[10px] font-bold text-muted-foreground tracking-[1px] uppercase mb-2">
              Минимальный рейтинг
            </p>
            <RatingPicker
              value={parseInt(filters.ratingMin) || 0}
              onChange={(v) =>
                onFilterChange("ratingMin", v > 0 ? String(v) : "")
              }
            />
          </div>
        </div>
      )}

      {/* Active chips */}
      {activeCount > 0 && (
        <div className="px-4 pb-2.5 pt-2 border-t border-border">
          <ActiveChips
            filters={filters}
            sectors={sectors}
            onClear={(k) => onFilterChange(k, "")}
            onClearAll={onReset}
          />
        </div>
      )}
    </div>
  )
}
