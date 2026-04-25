"use client"

import { T } from "./tokens"

interface TickerFooterProps {
  ticker: string
  buildTag: string
}

export function TickerFooter({ ticker, buildTag }: TickerFooterProps) {
  return (
    <div
      style={{
        marginTop: 24,
        paddingTop: 14,
        borderTop: `1px solid ${T.border}`,
        display: "flex",
        justifyContent: "space-between",
        fontFamily: T.mono,
        fontSize: 10,
        color: T.textFaint,
        letterSpacing: 0.5,
        textTransform: "uppercase",
        gap: 12,
        flexWrap: "wrap",
      }}
    >
      <span>
        TRADE COMPASS · TICKER {ticker} · DATA DELAYED 15MIN · MOEX / INTERFAX
      </span>
      <span>не индивидуальная инвестиционная рекомендация</span>
    </div>
  )
}
