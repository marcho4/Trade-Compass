import Link from "next/link"

// SVG-only color refs (CSS classes don't apply to SVG presentation attrs)
const SVG_BORDER_SOFT = "#edeae0"
const SVG_TEXT_DIM = "var(--color-muted-foreground)"
const SVG_ACCENT = "var(--color-primary)"

// ── Data ─────────────────────────────────────────────────────────────────────
const TICKERS = [
  {
    symbol: "MOEX", name: "ПАО Московская Биржа", sector: "Финансы",
    price: 172.20, changePct: 0.06,
    marketCap: "₽391.7B", pe: 8.4, divYield: "6.1%", rating: 5,
    scores: { health: 0.92, growth: 0.34, ratio: 0.88, dividends: 0.74, valuation: 0.81 },
    spark: [168.1, 168.4, 167.9, 169.2, 170.1, 169.8, 170.6, 171.0, 170.4, 171.2, 171.8, 172.1, 171.6, 172.0, 172.20],
  },
  {
    symbol: "GAZP", name: "ПАО Газпром", sector: "Энергетика",
    price: 138.42, changePct: -1.31,
    marketCap: "₽3.27T", pe: 4.1, divYield: "0.0%", rating: 2,
    scores: { health: 0.41, growth: 0.22, ratio: 0.55, dividends: 0.10, valuation: 0.62 },
    spark: [144.0, 143.4, 142.9, 142.1, 141.8, 141.0, 140.6, 140.4, 140.0, 139.5, 139.2, 139.0, 138.6, 138.5, 138.42],
  },
  {
    symbol: "SBER", name: "ПАО Сбербанк", sector: "Финансы",
    price: 312.85, changePct: 1.36,
    marketCap: "₽6.75T", pe: 5.2, divYield: "11.8%", rating: 4,
    scores: { health: 0.85, growth: 0.62, ratio: 0.91, dividends: 0.94, valuation: 0.76 },
    spark: [305.4, 306.1, 306.0, 307.3, 308.0, 308.9, 309.6, 310.2, 310.8, 311.5, 311.9, 312.3, 312.6, 312.7, 312.85],
  },
]

const TICKER_TAPE_ITEMS = [
  { sym: "MOEX", price: "172.20", chg: "+0.06%", pos: true },
  { sym: "SBER", price: "312.85", chg: "+1.36%", pos: true },
  { sym: "GAZP", price: "138.42", chg: "−1.31%", pos: false },
  { sym: "LKOH", price: "7345.50", chg: "+0.42%", pos: true },
  { sym: "YNDX", price: "3412.00", chg: "−0.18%", pos: false },
  { sym: "ROSN", price: "562.40", chg: "+0.77%", pos: true },
  { sym: "NVTK", price: "1087.20", chg: "+2.14%", pos: true },
  { sym: "MTSS", price: "264.80", chg: "−0.44%", pos: false },
]

const SCORE_AXES = [
  { key: "health",    label: "Здоровье"  },
  { key: "growth",    label: "Рост"      },
  { key: "ratio",     label: "Ров"       },
  { key: "dividends", label: "Дивиденды" },
  { key: "valuation", label: "Оценка"    },
]

function fmtPrice(p: number) {
  return p.toLocaleString("ru-RU", { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// ── Sparkline ─────────────────────────────────────────────────────────────────
function Sparkline({ data, w = 90, h = 32, color }: { data: number[]; w?: number; h?: number; color: string }) {
  const min = Math.min(...data)
  const max = Math.max(...data)
  const range = max - min || 1
  const step = w / (data.length - 1)
  const pts = data.map((v, i) =>
    `${(i * step).toFixed(1)},${(h - ((v - min) / range) * (h - 4) - 2).toFixed(1)}`
  )
  const [lx, ly] = pts[pts.length - 1].split(",")
  return (
    <svg width={w} height={h} viewBox={`0 0 ${w} ${h}`} style={{ display: "block" }}>
      <polyline points={pts.join(" ")} fill="none" stroke={color} strokeWidth={1.4} strokeLinejoin="round" strokeLinecap="round" />
      <circle cx={lx} cy={ly} r={2} fill={color} />
    </svg>
  )
}

// ── Radar ─────────────────────────────────────────────────────────────────────
type Scores = Record<string, number>

function Radar({ scores, size = 160 }: { scores: Scores; size?: number }) {
  const c = size / 2
  const r = c - 36
  const n = SCORE_AXES.length

  const axisPt = (i: number, mag: number): [number, number] => {
    const a = -Math.PI / 2 + (i * 2 * Math.PI) / n
    return [c + Math.cos(a) * r * mag, c + Math.sin(a) * r * mag]
  }

  const ringPath = (mag: number) =>
    SCORE_AXES.map((_, i) => {
      const [x, y] = axisPt(i, mag)
      return `${i === 0 ? "M" : "L"}${x.toFixed(1)},${y.toFixed(1)}`
    }).join(" ") + " Z"

  const dataPath =
    SCORE_AXES.map((ax, i) => {
      const [x, y] = axisPt(i, scores[ax.key])
      return `${i === 0 ? "M" : "L"}${x.toFixed(1)},${y.toFixed(1)}`
    }).join(" ") + " Z"

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} style={{ display: "block", overflow: "visible" }}>
      {[0.25, 0.5, 0.75, 1.0].map((m, i) => (
        <path key={i} d={ringPath(m)} fill="none" stroke={SVG_BORDER_SOFT} strokeWidth={1} strokeDasharray={i === 3 ? "0" : "2 3"} />
      ))}
      {SCORE_AXES.map((_, i) => {
        const [x, y] = axisPt(i, 1)
        return <line key={i} x1={c} y1={c} x2={x} y2={y} stroke={SVG_BORDER_SOFT} strokeWidth={1} strokeDasharray="2 3" />
      })}
      <path d={dataPath} fill={SVG_ACCENT} fillOpacity={0.18} stroke={SVG_ACCENT} strokeWidth={1.4} />
      {SCORE_AXES.map((ax, i) => {
        const [x, y] = axisPt(i, scores[ax.key])
        return <rect key={i} x={x - 2} y={y - 2} width={4} height={4} fill={SVG_ACCENT} />
      })}
      {SCORE_AXES.map((ax, i) => {
        const [x, y] = axisPt(i, 1.18)
        const a = -Math.PI / 2 + (i * 2 * Math.PI) / n
        const cos = Math.cos(a)
        const anchor = Math.abs(cos) < 0.2 ? "middle" : cos > 0 ? "start" : "end"
        const dy = Math.sin(a) > 0.3 ? 9 : Math.sin(a) < -0.3 ? -2 : 4
        return (
          <text
            key={i} x={x} y={y} textAnchor={anchor} dy={dy}
            style={{ fontFamily: "ui-monospace, monospace", fontSize: 9, fontWeight: 700, fill: SVG_TEXT_DIM, letterSpacing: 1, textTransform: "uppercase" }}
          >
            {ax.label}
          </text>
        )
      })}
    </svg>
  )
}

// ── RatingDots ────────────────────────────────────────────────────────────────
function RatingDots({ value }: { value: number }) {
  return (
    <div className="flex gap-[2px]">
      {Array.from({ length: 5 }).map((_, i) => (
        <span
          key={i}
          className={`inline-block h-1.5 w-1.5 border ${i < value ? "border-primary bg-primary" : "border-border bg-transparent"}`}
        />
      ))}
    </div>
  )
}

// ── HeroCard ──────────────────────────────────────────────────────────────────
type TickerData = (typeof TICKERS)[0]

function HeroCard({ ticker: t, variant = "full", size = 260 }: {
  ticker: TickerData
  variant?: "full" | "compact"
  size?: number
}) {
  const isPos = t.changePct > 0
  const isNeg = t.changePct < 0
  const sparkColor = isPos ? "var(--color-positive)" : isNeg ? "var(--color-negative)" : "var(--color-muted-foreground)"
  const cArrow = isPos ? "↑" : isNeg ? "↓" : "→"
  const cClass = isPos ? "text-positive" : isNeg ? "text-negative" : "text-muted-foreground"

  return (
    <div
      className="flex flex-col rounded-[2px] border border-border bg-card shadow-[0_1px_0_rgba(20,20,20,0.02),0_18px_40px_-18px_rgba(20,30,50,0.22),0_8px_20px_-12px_rgba(20,30,50,0.14)]"
      style={{ width: size }}
    >
      {/* Header strip */}
      <div className="flex items-center justify-between border-b border-border bg-[#fafaf7] px-3 py-2">
        <div className="flex items-center gap-2">
          <span className="inline-block h-[10px] w-[3px] bg-primary" />
          <span className="font-mono text-[10px] font-bold tracking-[1.2px] text-foreground">
            {t.symbol}<span className="text-muted-foreground/60">.MOEX</span>
          </span>
        </div>
        <span className="font-mono text-[9px] tracking-[0.5px] text-muted-foreground/60">{t.sector}</span>
      </div>

      {/* Identity */}
      <div className="flex items-start justify-between gap-2.5 border-b border-dashed border-[#edeae0] px-3 pb-2 pt-2.5">
        <div>
          <div className="font-mono text-2xl font-bold leading-none tracking-[-0.8px] text-foreground">{t.symbol}</div>
          <div className="mt-1 max-w-[150px] font-sans text-[11px] leading-tight text-[#363b45]">{t.name}</div>
        </div>
        <div className="text-right">
          <div className="mb-[3px] font-mono text-[8px] font-bold uppercase tracking-[1px] text-muted-foreground">Рейт.</div>
          <div className="font-mono text-base font-bold leading-none text-foreground">
            {t.rating}<span className="text-[10px] text-muted-foreground/60">/5</span>
          </div>
          <div className="mt-1 flex justify-end">
            <RatingDots value={t.rating} />
          </div>
        </div>
      </div>

      {/* Radar — full variant only */}
      {variant !== "compact" && (
        <div className="flex justify-center border-b border-dashed border-[#edeae0] px-1.5 pb-1.5 pt-1">
          <Radar scores={t.scores} size={160} />
        </div>
      )}

      {/* Price */}
      <div className="grid grid-cols-[1fr_auto] items-end gap-2.5 px-3 py-2">
        <div>
          <div className="mb-[3px] font-mono text-[8px] font-bold uppercase tracking-[1px] text-muted-foreground">Цена</div>
          <div className="font-mono text-[18px] font-bold leading-none tracking-[-0.3px] text-foreground">
            {fmtPrice(t.price)}{" "}
            <span className="text-[11px] font-semibold text-muted-foreground">₽</span>
          </div>
          <div className={`mt-1 font-mono text-[10px] font-semibold ${cClass}`}>
            {cArrow} {t.changePct >= 0 ? "+" : ""}{t.changePct.toFixed(2)}%
          </div>
        </div>
        <Sparkline data={t.spark} color={sparkColor} w={90} h={32} />
      </div>

      {/* Footer KV */}
      <div className="grid grid-cols-3 border-t border-border bg-[#fafaf7]">
        {(
          [
            ["Кап", t.marketCap],
            ["P/E", t.pe.toFixed(1)],
            ["Див", t.divYield],
          ] as [string, string][]
        ).map(([k, v], i) => (
          <div key={k} className={`px-2.5 py-[7px] ${i < 2 ? "border-r border-border" : ""}`}>
            <div className="font-mono text-[8px] font-bold uppercase tracking-[1px] text-muted-foreground">{k}</div>
            <div className="mt-0.5 font-mono text-[11px] font-semibold text-foreground">{v}</div>
          </div>
        ))}
      </div>
    </div>
  )
}

// ── Single composition: 1 main card + 1 ghost behind ─────────────────────────
function CardSingle() {
  return (
    <div className="relative h-[560px] w-[620px] shrink-0">
      {/* Ghost — GAZP, compact, dimmed */}
      <div className="absolute left-[60px] top-[50px] z-[1] rotate-[-4deg] opacity-50 [filter:grayscale(0.4)]">
        <HeroCard ticker={TICKERS[1]} variant="compact" size={280} />
      </div>
      {/* Main — MOEX, full */}
      <div className="absolute left-[160px] top-[90px] z-[2] rotate-[3deg]">
        <HeroCard ticker={TICKERS[0]} size={320} />
      </div>
    </div>
  )
}

// ── TickerTape ────────────────────────────────────────────────────────────────
function TickerTape() {
  const row = [...TICKER_TAPE_ITEMS, ...TICKER_TAPE_ITEMS]
  return (
    <div className="overflow-hidden border-y border-border bg-card py-2.5">
      <div className="animate-ticker-scroll flex w-max gap-11 whitespace-nowrap">
        {row.map((it, i) => (
          <span key={i} className="inline-flex items-baseline gap-2.5 font-mono text-[12px] tracking-[0.5px]">
            <span className="inline-block h-[10px] w-[3px] align-middle bg-primary" />
            <span className="font-bold text-foreground">{it.sym}</span>
            <span className="text-[#363b45]">{it.price}</span>
            <span className={`font-semibold ${it.pos ? "text-positive" : "text-negative"}`}>{it.chg}</span>
          </span>
        ))}
      </div>
    </div>
  )
}

// ── HeroSection ───────────────────────────────────────────────────────────────
export const HeroSection = () => (
  <section className="relative overflow-hidden bg-background">
    {/* Grid background */}
    <div className="pointer-events-none absolute inset-0 [background-image:linear-gradient(to_right,#edeae0_1px,transparent_1px),linear-gradient(to_bottom,#edeae0_1px,transparent_1px)] [background-size:60px_60px] [mask-image:linear-gradient(to_bottom,rgba(0,0,0,0.6),rgba(0,0,0,0.1))]" />

    <div className="relative z-[2] mx-auto max-w-[1440px] px-11 py-[72px]">
      <div
        className="grid min-h-[560px] items-center gap-10"
        style={{ gridTemplateColumns: "minmax(480px, 1fr) 620px" }}
      >
        {/* Left: headline block */}
        <div>
          {/* Eyebrow */}
          <div className="inline-flex items-center gap-2.5 font-mono text-[11px] font-bold uppercase tracking-[1.8px] text-primary">
            <span className="inline-block h-px w-6 bg-primary" />
            <span>ПЛАТФОРМА · АНАЛИЗ · РОССИЙСКИЕ АКЦИИ</span>
          </div>

          {/* Terminal headline */}
          <h1 className="mt-[22px] font-mono text-[clamp(38px,4.6vw,64px)] font-bold uppercase leading-[0.98] tracking-[-1.5px] text-foreground">
            <span className="font-medium text-muted-foreground/50">&gt;&nbsp;</span>
            Укажем
            <br />
            направление
            <br />
            <span className="text-primary">для инвестиций_</span>
          </h1>

          {/* Sub */}
          <p className="mb-8 mt-6 max-w-[420px] font-sans text-base leading-[1.5] text-[#363b45]">
            Анализируйте рынок в 10 раз быстрее. Без финансового образования.
          </p>

          {/* Outline CTA */}
          <Link
            href="/dashboard/screener"
            className="inline-flex items-center gap-2.5 rounded-[2px] border-[1.5px] border-foreground bg-transparent px-[22px] py-[14px] font-mono text-[13px] font-bold uppercase tracking-[1.2px] text-foreground no-underline transition-colors hover:bg-foreground hover:text-background"
          >
            <span>ПОПРОБОВАТЬ БЕСПЛАТНО</span>
            <span className="text-[15px]">→</span>
          </Link>

          {/* Badges */}
          <div className="mt-5 flex items-center gap-3.5 font-mono text-[11px] tracking-[0.5px] text-muted-foreground">
            <span className="inline-flex items-center gap-1.5">
              <span className="inline-block h-1.5 w-1.5 rounded-full bg-positive" />
              Без карты
            </span>
            <span className="text-muted-foreground/40">·</span>
            <span>Данные по всем крупным компаниям Мосбиржи</span>
          </div>
        </div>

        {/* Right: cards */}
        <div className="flex items-center justify-center">
          <CardSingle />
        </div>
      </div>
    </div>

    <TickerTape />
  </section>
)
