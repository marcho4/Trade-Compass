import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatShort(value: number | null | undefined): string {
  if (value === null || value === undefined) return "—"

  const abs = Math.abs(value)
  const sign = value < 0 ? "-" : ""

  if (abs >= 1_000_000_000_000) {
    return sign + (abs / 1_000_000_000_000).toFixed(1).replace(/\.0$/, "") + "T"
  }
  if (abs >= 1_000_000_000) {
    return sign + (abs / 1_000_000_000).toFixed(1).replace(/\.0$/, "") + "B"
  }
  if (abs >= 1_000_000) {
    return sign + (abs / 1_000_000).toFixed(1).replace(/\.0$/, "") + "M"
  }
  if (abs >= 1_000) {
    return sign + (abs / 1_000).toFixed(1).replace(/\.0$/, "") + "K"
  }
  return value.toString()
}

export function formatLargeNumber(value: number): string {
  if (!value) return "0"
  
  if (value >= 1_000_000_000_000) {
    return (value / 1_000_000_000_000).toFixed(1).replace(/\.0$/, "") + " трлн"
  }
  if (value >= 1_000_000_000) {
    return (value / 1_000_000_000).toFixed(1).replace(/\.0$/, "") + " млрд"
  }
  if (value >= 1_000_000) {
    return (value / 1_000_000).toFixed(1).replace(/\.0$/, "") + " млн"
  }
  if (value >= 1_000) {
    return (value / 1_000).toFixed(1).replace(/\.0$/, "") + " тыс"
  }
  
  return value.toString()
}
