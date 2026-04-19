"use client"

import { ReactNode, CSSProperties } from "react"
import { T } from "./tokens"

interface TPanelHeaderProps {
  title: string
  count?: number | null
  accent?: string
  right?: ReactNode
}

export function TPanelHeader({ title, count, accent, right }: TPanelHeaderProps) {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        padding: "10px 14px",
        borderBottom: `1px solid ${T.border}`,
        background: T.panelAlt,
      }}
    >
      <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
        <span style={{ width: 3, height: 12, background: accent || T.accent }} />
        <span
          style={{
            fontFamily: T.mono,
            fontSize: 11,
            fontWeight: 700,
            letterSpacing: 1.2,
            textTransform: "uppercase",
            color: T.text,
          }}
        >
          {title}
        </span>
        {count != null && (
          <span style={{ fontFamily: T.mono, fontSize: 10, color: T.textFaint }}>
            [{String(count).padStart(3, "0")}]
          </span>
        )}
      </div>
      {right}
    </div>
  )
}

interface TPanelProps extends TPanelHeaderProps {
  children: ReactNode
  maxH?: number
  style?: CSSProperties
}

export function TPanel({ title, count, accent, children, right, maxH, style }: TPanelProps) {
  return (
    <div
      style={{
        background: T.panel,
        border: `1px solid ${T.border}`,
        borderRadius: 2,
        overflow: "hidden",
        display: "flex",
        flexDirection: "column",
        boxShadow: "0 1px 0 rgba(20,20,20,0.02)",
        ...style,
      }}
    >
      <TPanelHeader title={title} count={count} accent={accent} right={right} />
      <div style={{ overflowY: maxH ? "auto" : "visible", maxHeight: maxH }}>{children}</div>
    </div>
  )
}

interface KVProps {
  k: string
  v: ReactNode
  width?: number
}

export function KV({ k, v, width = 120 }: KVProps) {
  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: `${width}px 1fr`,
        gap: 14,
        alignItems: "baseline",
        padding: "8px 0",
        borderBottom: `1px dashed ${T.borderSoft}`,
        fontSize: 12.5,
      }}
    >
      <span
        style={{
          fontFamily: T.mono,
          fontSize: 10,
          fontWeight: 600,
          textTransform: "uppercase",
          letterSpacing: 0.8,
          color: T.textDim,
        }}
      >
        {k}
      </span>
      <span style={{ color: T.text, lineHeight: 1.45 }}>{v}</span>
    </div>
  )
}

interface StatProps {
  label: string
  value: ReactNode
  sub?: ReactNode
  subColor?: string
  mono?: boolean
  align?: "left" | "center" | "right"
}

export function Stat({ label, value, sub, subColor, mono = true, align = "left" }: StatProps) {
  return (
    <div style={{ textAlign: align }}>
      <div
        style={{
          fontFamily: T.mono,
          fontSize: 9.5,
          fontWeight: 600,
          textTransform: "uppercase",
          letterSpacing: 1,
          color: T.textDim,
          marginBottom: 4,
        }}
      >
        {label}
      </div>
      <div
        style={{
          fontFamily: mono ? T.mono : T.sans,
          fontSize: 16,
          fontWeight: 600,
          color: T.text,
          lineHeight: 1.1,
          letterSpacing: mono ? 0 : -0.2,
        }}
      >
        {value}
      </div>
      {sub && (
        <div
          style={{
            fontFamily: T.mono,
            fontSize: 10,
            color: subColor || T.textDim,
            marginTop: 3,
            letterSpacing: 0.2,
          }}
        >
          {sub}
        </div>
      )}
    </div>
  )
}

interface ChipProps {
  children: ReactNode
  color?: string
  bg?: string
  bd?: string
}

export function Chip({ children, color = T.text, bg = T.panelAlt, bd = T.border }: ChipProps) {
  return (
    <span
      style={{
        display: "inline-flex",
        alignItems: "center",
        fontFamily: T.mono,
        fontSize: 10,
        fontWeight: 600,
        letterSpacing: 0.8,
        textTransform: "uppercase",
        color,
        background: bg,
        border: `1px solid ${bd}`,
        padding: "3px 7px",
        borderRadius: 2,
        lineHeight: 1.4,
      }}
    >
      {children}
    </span>
  )
}

interface DonutSlice {
  label: string
  pct: number
  color: string
}

interface DonutProps {
  data: DonutSlice[]
  size?: number
  thickness?: number
  centerLabel?: string
  centerValue?: string
}

export function Donut({ data, size = 176, thickness = 28, centerLabel, centerValue }: DonutProps) {
  const r = (size - thickness) / 2
  const c = size / 2
  const circ = 2 * Math.PI * r

  const segments = data.reduce<
    { len: number; offset: number; color: string; key: number }[]
  >((acc, d, i) => {
    const len = (d.pct / 100) * circ
    const offset = acc.reduce((s, seg) => s + seg.len, 0)
    acc.push({ len, offset, color: d.color, key: i })
    return acc
  }, [])

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      <circle cx={c} cy={c} r={r} fill="none" stroke={T.borderSoft} strokeWidth={thickness} />
      {segments.map((seg) => (
        <circle
          key={seg.key}
          cx={c}
          cy={c}
          r={r}
          fill="none"
          stroke={seg.color}
          strokeWidth={thickness}
          strokeDasharray={`${seg.len} ${circ - seg.len}`}
          strokeDashoffset={-seg.offset}
          transform={`rotate(-90 ${c} ${c})`}
        />
      ))}
      {centerValue && (
        <>
          <text
            x={c}
            y={c - 4}
            textAnchor="middle"
            style={{ fontFamily: T.mono, fontSize: 9, fill: T.textDim, letterSpacing: 1 }}
          >
            {centerLabel}
          </text>
          <text
            x={c}
            y={c + 14}
            textAnchor="middle"
            style={{ fontFamily: T.mono, fontSize: 18, fontWeight: 700, fill: T.text }}
          >
            {centerValue}
          </text>
        </>
      )}
    </svg>
  )
}

interface RangeBarProps {
  low: number
  high: number
  value: number
  width?: number
}

export function RangeBar({ low, high, value, width = 220 }: RangeBarProps) {
  const pct = high === low ? 50 : ((value - low) / (high - low)) * 100
  const clamped = Math.max(0, Math.min(100, pct))

  return (
    <div style={{ width }}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          fontFamily: T.mono,
          fontSize: 9,
          color: T.textDim,
          letterSpacing: 0.5,
          marginBottom: 5,
        }}
      >
        <span>L {low.toFixed(2)}</span>
        <span style={{ color: T.textFaint }}>52W RANGE</span>
        <span>H {high.toFixed(2)}</span>
      </div>
      <div style={{ position: "relative", height: 6, background: T.borderSoft, borderRadius: 1 }}>
        <div
          style={{
            position: "absolute",
            left: 0,
            top: 0,
            height: "100%",
            width: `${clamped}%`,
            background: `linear-gradient(90deg, ${T.border}, ${T.accent})`,
          }}
        />
        <div
          style={{
            position: "absolute",
            left: `calc(${clamped}% - 1px)`,
            top: -3,
            width: 2,
            height: 12,
            background: T.text,
          }}
        />
      </div>
    </div>
  )
}
