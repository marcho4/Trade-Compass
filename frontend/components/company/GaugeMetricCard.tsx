"use client"

import { useId } from "react"
import { Card } from "@/components/ui/card"

interface GaugeMetricCardProps {
  label: string
  min: number
  max: number
  companyValue: number | null
  sectorValue: number | null
  format?: "number" | "percent" | "ratio"
  higherIsBetter?: boolean
}

const CX = 100
const CY = 138
const R_OUTER = 88
const R_INNER = 68
const R_MID = (R_OUTER + R_INNER) / 2
const ARC_WIDTH = R_OUTER - R_INNER

function p2c(r: number, angleDeg: number): { x: number; y: number } {
  const rad = (angleDeg * Math.PI) / 180
  return { x: CX + r * Math.cos(rad), y: CY - r * Math.sin(rad) }
}

function valToAngle(value: number, min: number, max: number): number {
  const clamped = Math.max(min, Math.min(max, value))
  return 180 - ((clamped - min) / (max - min)) * 180
}

function fmt(val: number | null, format: string): string {
  if (val === null) return "—"
  switch (format) {
    case "percent":
      return `${val.toFixed(1)}%`
    case "ratio":
      return val.toFixed(2)
    default:
      return val.toLocaleString("ru-RU", { maximumFractionDigits: 2 })
  }
}

function companyNeedlePoints(angle: number): string {
  const rad = (angle * Math.PI) / 180
  const px = Math.sin(rad)
  const py = Math.cos(rad)
  const tip = p2c(R_INNER - 1, angle)
  const near = p2c(R_INNER - 13, angle)
  const base = p2c(34, angle)
  return [
    `${tip.x},${tip.y}`,
    `${near.x + 1.5 * px},${near.y + 1.5 * py}`,
    `${base.x + 4 * px},${base.y + 4 * py}`,
    `${base.x - 4 * px},${base.y - 4 * py}`,
    `${near.x - 1.5 * px},${near.y - 1.5 * py}`,
  ].join(" ")
}

function sectorArrowPoints(angle: number): string {
  const rad = (angle * Math.PI) / 180
  const px = Math.sin(rad)
  const py = Math.cos(rad)
  const tip = p2c(R_OUTER + 3, angle)
  const base = p2c(R_OUTER + 15, angle)
  const hw = 6
  return [
    `${tip.x},${tip.y}`,
    `${base.x + hw * px},${base.y + hw * py}`,
    `${base.x - hw * px},${base.y - hw * py}`,
  ].join(" ")
}

export function GaugeMetricCard({
  label,
  min,
  max,
  companyValue,
  sectorValue,
  format = "ratio",
  higherIsBetter = true,
}: GaugeMetricCardProps) {
  const uid = useId().replace(/:/g, "")
  const gradId = `gauge-grad-${uid}`

  const arcStart = p2c(R_MID, 180)
  const arcEnd = p2c(R_MID, 0)
  const arcPath = `M ${arcStart.x} ${arcStart.y} A ${R_MID} ${R_MID} 0 0 0 ${arcEnd.x} ${arcEnd.y}`

  const companyAngle = companyValue !== null ? valToAngle(companyValue, min, max) : null
  const sectorAngle = sectorValue !== null ? valToAngle(sectorValue, min, max) : null

  const minPt = p2c(R_MID, 180)
  const maxPt = p2c(R_MID, 0)

  const colorStart = higherIsBetter ? "#ef4444" : "#22c55e"
  const colorMid = "#eab308"
  const colorEnd = higherIsBetter ? "#22c55e" : "#ef4444"

  return (
    <Card className="p-3">
      <svg viewBox="0 0 200 200" className="w-full">
        <defs>
          <linearGradient
            id={gradId}
            x1={CX - R_MID}
            y1="0"
            x2={CX + R_MID}
            y2="0"
            gradientUnits="userSpaceOnUse"
          >
            <stop offset="0%" stopColor={colorStart} />
            <stop offset="50%" stopColor={colorMid} />
            <stop offset="100%" stopColor={colorEnd} />
          </linearGradient>
        </defs>

        {/* Title */}
        <text
          x="100"
          y="18"
          textAnchor="middle"
          fontSize="13"
          fontWeight="600"
          fill="currentColor"
        >
          {label}
        </text>

        {/* Arc track (background) */}
        <path
          d={arcPath}
          fill="none"
          stroke="currentColor"
          strokeOpacity={0.1}
          strokeWidth={ARC_WIDTH}
        />

        {/* Colored gradient arc */}
        <path
          d={arcPath}
          fill="none"
          stroke={`url(#${gradId})`}
          strokeWidth={ARC_WIDTH}
          strokeOpacity={0.9}
        />

        {/* Min label */}
        <text
          x={minPt.x}
          y={CY + 18}
          textAnchor="middle"
          fontSize="9"
          fill="currentColor"
          fillOpacity={0.45}
        >
          {fmt(min, format)}
        </text>

        {/* Max label */}
        <text
          x={maxPt.x}
          y={CY + 18}
          textAnchor="middle"
          fontSize="9"
          fill="currentColor"
          fillOpacity={0.45}
        >
          {fmt(max, format)}
        </text>

        {/* Sector arrow — from outside, pointing inward */}
        {sectorAngle !== null && (
          <polygon points={sectorArrowPoints(sectorAngle)} fill="#3b82f6" />
        )}

        {/* Company needle — from center, pointing to arc */}
        {companyAngle !== null && (
          <polygon points={companyNeedlePoints(companyAngle)} fill="currentColor" />
        )}

        {/* Pivot */}
        <circle cx={CX} cy={CY} r="5" fill="currentColor" />
        <circle cx={CX} cy={CY} r="2.5" fill="var(--card)" />

        {/* Legend */}
        <g transform="translate(0, 170)">
          {/* Company */}
          <rect x="18" y="0" width="9" height="9" rx="2" fill="currentColor" opacity={0.8} />
          <text x="31" y="8" fontSize="10" fill="currentColor" fillOpacity={0.55}>
            Компания:
          </text>
          <text x="84" y="8" fontSize="10" fontWeight="600" fill="currentColor">
            {fmt(companyValue, format)}
          </text>

          {/* Sector */}
          <rect x="18" y="16" width="9" height="9" rx="2" fill="#3b82f6" />
          <text x="31" y="24" fontSize="10" fill="currentColor" fillOpacity={0.55}>
            Сектор:
          </text>
          <text x="78" y="24" fontSize="10" fontWeight="600" fill="currentColor">
            {fmt(sectorValue, format)}
          </text>
        </g>
      </svg>
    </Card>
  )
}
