"use client"

import { ReactNode } from "react"
import { cn } from "@/lib/utils"

interface LimitsData {
  plan: string
  period: string
  used: number
  total: number
  resetsAt: string
  daysToReset: number
  hoursToReset: number
  dailyUsed: number
  dailyTotal: number
  savedUsed: number
  savedTotal: number
  screenersUsed: number
  screenersTotal: number
  peakPerDay: number
  peakDate: string
  avgPerDay: number
  heatmap: number[]
}

interface RecentItem {
  sym: string
  name: string
  sector: string
  when: string
  views: number
}

interface PlanOption {
  id: string
  views: number | string
  daily: number | string
  saved: number | string
  screeners: number | string
  price: string
}

interface LimitsCardProps {
  limits: LimitsData
  recent: RecentItem[]
  plans: PlanOption[]
  onExtend?: () => void
  onSwitchPlan?: (planId: string) => void
  onFullJournal?: () => void
}

function SectionLabel({ children, right }: { children: ReactNode; right?: ReactNode }) {
  return (
    <div className="flex items-center justify-between border-b border-border bg-muted/40 px-3.5 py-2.5">
      <div className="flex items-center gap-2.5">
        <span className="h-3 w-[3px] bg-primary" />
        <span className="font-mono text-[11px] font-bold uppercase tracking-[0.1em] text-foreground">
          {children}
        </span>
      </div>
      {right}
    </div>
  )
}

function DashBtn({
  children,
  primary,
  onClick,
  className,
}: {
  children: ReactNode
  primary?: boolean
  onClick?: () => void
  className?: string
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "cursor-pointer rounded-[2px] border px-2.5 py-1.5 font-mono text-[10.5px] font-bold uppercase tracking-[0.08em] transition-colors",
        primary
          ? "border-primary bg-primary text-primary-foreground hover:bg-primary/90"
          : "border-border bg-card text-foreground hover:bg-muted/50",
        className,
      )}
    >
      {children}
    </button>
  )
}

function QuotaMeter({
  label,
  sub,
  used,
  total,
}: {
  label: string
  sub?: string
  used: number
  total: number
}) {
  const pct = Math.min(100, Math.round((used / total) * 100))
  const over80 = pct >= 80
  const filledCells = Math.round((used / total) * 40)

  return (
    <div className="border-b border-dashed border-border/50 py-3.5">
      <div className="mb-2 flex items-baseline justify-between">
        <div>
          <div className="font-mono text-[10px] font-bold uppercase tracking-[0.08em] text-muted-foreground">
            {label}
          </div>
          {sub && (
            <div className="mt-0.5 font-mono text-[10px] tracking-[0.02em] text-muted-foreground/70">
              {sub}
            </div>
          )}
        </div>
        <div className="font-mono text-[14px] font-bold tracking-[-0.01em] text-foreground">
          {used}
          <span className="font-medium text-muted-foreground/60"> / {total}</span>
          <span className="ml-2 text-[11px] text-muted-foreground">{pct}%</span>
        </div>
      </div>
      <div className="grid h-2.5 grid-cols-[repeat(40,minmax(0,1fr))] gap-[2px]">
        {Array.from({ length: 40 }).map((_, i) => (
          <span
            key={i}
            className={cn(
              i < filledCells ? (over80 ? "bg-[#b8751a]" : "bg-primary") : "bg-border/60",
            )}
          />
        ))}
      </div>
    </div>
  )
}

function HeatStrip({ heatmap }: { heatmap: number[] }) {
  const scale = [
    "bg-[#ece9df]",
    "bg-[#d7e2ef]",
    "bg-[#a9c5e5]",
    "bg-[#6b99c9]",
    "bg-primary",
  ]
  return (
    <div>
      <div className="mb-1.5 flex justify-between font-mono text-[9.5px] uppercase tracking-[0.05em] text-muted-foreground">
        <span>30 дней активности</span>
        <span className="text-muted-foreground/70">
          меньше{" "}
          {scale.map((c, i) => (
            <span key={i} className={cn("ml-0.5 inline-block h-2 w-2 align-middle", c)} />
          ))}{" "}
          больше
        </span>
      </div>
      <div className="grid h-[22px] grid-cols-[repeat(30,minmax(0,1fr))] gap-[2px]">
        {heatmap.map((v, i) => (
          <span
            key={i}
            title={`день -${29 - i}: ${v} просмотров`}
            className={cn(scale[Math.min(v, 4)])}
          />
        ))}
      </div>
    </div>
  )
}

function Stat({ label, value, sub }: { label: string; value: ReactNode; sub?: ReactNode }) {
  return (
    <div>
      <div className="mb-1 font-mono text-[9.5px] font-semibold uppercase tracking-[0.08em] text-muted-foreground">
        {label}
      </div>
      <div className="font-mono text-[16px] font-semibold leading-[1.1] text-foreground">
        {value}
      </div>
      {sub && (
        <div className="mt-0.5 font-mono text-[10px] tracking-[0.02em] text-muted-foreground">
          {sub}
        </div>
      )}
    </div>
  )
}

export const LimitsCard = ({
  limits,
  recent,
  plans,
  onExtend,
  onSwitchPlan,
  onFullJournal,
}: LimitsCardProps) => {
  const mainCells = 60
  const filledMain = Math.round((limits.used / limits.total) * mainCells)

  return (
    <div className="flex flex-col rounded-[2px] border border-border bg-card">
      <SectionLabel
        right={
          <div className="flex items-center gap-2.5">
            <span className="font-mono text-[10px] tracking-[0.05em] text-muted-foreground">
              Сброс через{" "}
              <span className="font-bold text-foreground">
                {limits.daysToReset}д {limits.hoursToReset}ч
              </span>
            </span>
            <DashBtn primary onClick={onExtend}>
              Продлить {limits.plan}
            </DashBtn>
          </div>
        }
      >
        Лимиты просмотров · {limits.plan}
      </SectionLabel>

      {/* Hero meter */}
      <div className="grid grid-cols-[1.4fr_1fr] gap-7 px-[22px] pt-5 pb-1.5">
        <div>
          <div className="mb-1 font-mono text-[10px] font-bold uppercase tracking-[0.1em] text-muted-foreground">
            {limits.period} · просмотры компаний
          </div>
          <div className="flex items-baseline gap-2.5 font-mono tracking-[-0.03em]">
            <span className="text-[52px] font-bold leading-[1] text-foreground">
              {limits.used}
            </span>
            <span className="text-[22px] font-medium text-muted-foreground/60">
              / {limits.total}
            </span>
            <div className="ml-auto text-right font-mono text-[11px] tracking-[0.03em] text-muted-foreground">
              осталось{" "}
              <span className="text-[14px] font-bold text-foreground">
                {limits.total - limits.used}
              </span>{" "}
              просмотров
              <div className="mt-0.5 text-muted-foreground/70">
                сброс: {limits.resetsAt}
              </div>
            </div>
          </div>

          <div className="mt-4 grid h-3.5 grid-cols-[repeat(60,minmax(0,1fr))] gap-[2px]">
            {Array.from({ length: mainCells }).map((_, i) => {
              const filled = i < filledMain
              const over80 = i / mainCells >= 0.8
              return (
                <span
                  key={i}
                  className={cn(
                    filled ? (over80 ? "bg-[#b8751a]" : "bg-primary") : "bg-border/60",
                  )}
                />
              )
            })}
          </div>
          <div className="mt-1.5 flex justify-between font-mono text-[9.5px] tracking-[0.05em] text-muted-foreground/70">
            <span>0</span>
            <span>80% · Предупреждение</span>
            <span>{limits.total}</span>
          </div>
        </div>

        <div className="rounded-[2px] border border-border bg-muted/40 p-4">
          <HeatStrip heatmap={limits.heatmap} />
          <div className="mt-3.5 grid grid-cols-2 gap-3.5 border-t border-dashed border-border pt-3">
            <Stat label="Пик / день" value={limits.peakPerDay} sub={limits.peakDate} />
            <Stat
              label="Среднее / день"
              value={limits.avgPerDay.toFixed(1)}
              sub="за 30 дней"
            />
          </div>
        </div>
      </div>

      {/* Sub-quotas */}
      <div className="grid grid-cols-3 gap-x-7 px-[22px] pt-0.5 pb-3">
        <QuotaMeter
          label="Сегодня"
          sub="суточный лимит просмотров"
          used={limits.dailyUsed}
          total={limits.dailyTotal}
        />
        <QuotaMeter
          label="Сохранённые компании"
          sub="избранное в списках"
          used={limits.savedUsed}
          total={limits.savedTotal}
        />
        <QuotaMeter
          label="Свои скринеры"
          sub="сохранённые пресеты фильтров"
          used={limits.screenersUsed}
          total={limits.screenersTotal}
        />
      </div>

      {/* Recent */}
      <div className="mx-[22px] border-t border-border pt-3.5 pb-1.5">
        <div className="mb-2.5 flex items-baseline justify-between">
          <div className="font-mono text-[10px] font-bold uppercase tracking-[0.1em] text-muted-foreground">
            Недавно просмотренные · списано из квоты
          </div>
          <button
            type="button"
            onClick={onFullJournal}
            className="cursor-pointer font-mono text-[10px] font-bold uppercase tracking-[0.07em] text-primary hover:underline"
          >
            Полный журнал →
          </button>
        </div>
        <div className="grid grid-cols-[repeat(auto-fill,minmax(220px,1fr))] gap-1.5">
          {recent.map((r) => (
            <div
              key={r.sym}
              className="grid grid-cols-[44px_1fr_auto] items-center gap-2.5 rounded-[2px] border border-border bg-card p-2.5"
            >
              <span className="border-l-[3px] border-primary pl-1.5 font-mono text-[11px] font-bold tracking-[-0.015em] text-foreground">
                {r.sym}
              </span>
              <span className="min-w-0">
                <div className="truncate font-sans text-[12px] leading-[1.2] text-foreground">
                  {r.name}
                </div>
                <div className="mt-0.5 font-mono text-[9.5px] tracking-[0.03em] text-muted-foreground/70">
                  {r.sector} · {r.when}
                </div>
              </span>
              <span className="rounded-[2px] border border-border bg-muted/40 px-1.5 py-0.5 font-mono text-[10px] font-bold tracking-[0.05em] text-muted-foreground">
                ×{r.views}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Plan comparison */}
      <div className="mt-3.5 border-t border-border bg-muted/40 px-[22px] pt-3.5 pb-4">
        <div className="mb-2.5 font-mono text-[10px] font-bold uppercase tracking-[0.1em] text-muted-foreground">
          Нужно больше просмотров?
        </div>
        <div className="grid grid-cols-3 gap-2.5">
          {plans.map((p) => {
            const current = p.id === limits.plan
            return (
              <div
                key={p.id}
                className={cn(
                  "relative rounded-[2px] border p-3.5",
                  current
                    ? "border-primary bg-primary/10"
                    : "border-border bg-card",
                )}
              >
                <div className="flex items-baseline justify-between">
                  <span
                    className={cn(
                      "font-mono text-[12px] font-bold tracking-[0.08em]",
                      current ? "text-primary" : "text-foreground",
                    )}
                  >
                    {p.id}
                  </span>
                  <span className="font-mono text-[11px] font-bold text-foreground">
                    {p.price}
                    <span className="font-medium text-muted-foreground/60">/мес</span>
                  </span>
                </div>
                <div className="mt-2 space-y-0.5 font-mono text-[10px] leading-[1.6] tracking-[0.03em] text-muted-foreground">
                  <div>{p.views} просмотров/мес</div>
                  <div>
                    {p.daily} в день · {p.saved} сохр.
                  </div>
                  <div>{p.screeners} своих скринеров</div>
                </div>
                <div className="mt-2.5">
                  {current ? (
                    <DashBtn className="w-full">Текущий тариф</DashBtn>
                  ) : (
                    <DashBtn
                      primary={p.id === "PRO+"}
                      className="w-full"
                      onClick={() => onSwitchPlan?.(p.id)}
                    >
                      {p.id === "FREE" ? "Понизить" : "Перейти →"}
                    </DashBtn>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
