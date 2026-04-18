"use client"

import { T } from "./tokens"

interface TickerTopBarProps {
  ticker: string
  marketOpen: boolean
  marketLabel: string
}

export function TickerTopBar({ ticker, marketOpen, marketLabel }: TickerTopBarProps) {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        padding: "16px 28px",
        borderBottom: `1px solid ${T.border}`,
        background: T.panel,
      }}
    >
      <div style={{ display: "flex", alignItems: "center", gap: 28 }}>
        <div
          style={{
            fontFamily: T.mono,
            fontSize: 13,
            fontWeight: 700,
            color: T.accent,
            letterSpacing: 2,
          }}
        >
          TC<span style={{ color: T.textDim }}>/</span>TICKER
        </div>
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 12,
            fontFamily: T.mono,
            fontSize: 11,
            color: T.textDim,
            letterSpacing: 0.5,
          }}
        >
          <span style={{ color: T.textFaint }}>&gt;</span>
          <span>СКРИНЕР</span>
          <span style={{ color: T.textFaint }}>/</span>
          <span style={{ color: T.text }}>ТИКЕР «{ticker}»</span>
        </div>
      </div>
      <div style={{ display: "flex", alignItems: "center", gap: 20 }}>
        <span
          style={{
            fontFamily: T.mono,
            fontSize: 11,
            color: T.textDim,
            letterSpacing: 0.5,
          }}
        >
          <span
            style={{
              display: "inline-block",
              width: 6,
              height: 6,
              borderRadius: "50%",
              background: marketOpen ? T.pos : T.textFaint,
              marginRight: 6,
              verticalAlign: "middle",
            }}
          />
          {marketLabel}
        </span>
      </div>
    </div>
  )
}
