"use client"

import { T } from "./tokens"

export interface TickerTab {
  id: string
  label: string
  count?: number | null
}

interface TickerTabStripProps {
  tabs: TickerTab[]
  active: string
  onChange: (id: string) => void
}

export function TickerTabStrip({ tabs, active, onChange }: TickerTabStripProps) {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "stretch",
        borderBottom: `1px solid ${T.border}`,
        marginTop: 2,
        flexWrap: "wrap",
      }}
    >
      {tabs.map((tab) => {
        const isActive = tab.id === active
        return (
          <button
            key={tab.id}
            type="button"
            onClick={() => onChange(tab.id)}
            style={{
              appearance: "none",
              background: "transparent",
              border: "none",
              padding: "14px 22px 12px",
              cursor: "pointer",
              fontFamily: T.mono,
              fontSize: 11,
              fontWeight: 700,
              letterSpacing: 1.2,
              textTransform: "uppercase",
              color: isActive ? T.text : T.textDim,
              borderBottom: isActive ? `2px solid ${T.accent}` : "2px solid transparent",
              marginBottom: -1,
              display: "flex",
              alignItems: "center",
              gap: 8,
            }}
          >
            <span>{tab.label}</span>
            {tab.count != null && (
              <span
                style={{
                  fontSize: 9.5,
                  color: isActive ? T.accent : T.textFaint,
                  fontWeight: 600,
                }}
              >
                [{String(tab.count).padStart(3, "0")}]
              </span>
            )}
          </button>
        )
      })}
    </div>
  )
}
