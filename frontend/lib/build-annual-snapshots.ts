import { RawData } from "@/types/raw-data"

export interface AnnualSnapshot {
  period: string
  [field: string]: string | number | null
}

export type NumericRawDataKey = Extract<keyof RawData, "revenue" | "costOfRevenue" | "grossProfit" | "operatingExpenses" | "ebit" | "ebitda" | "interestExpense" | "taxExpense" | "netProfit" | "totalAssets" | "currentAssets" | "cashAndEquivalents" | "inventories" | "receivables" | "totalLiabilities" | "currentLiabilities" | "debt" | "longTermDebt" | "shortTermDebt" | "equity" | "retainedEarnings" | "operatingCashFlow" | "investingCashFlow" | "financingCashFlow" | "capex" | "freeCashFlow" | "sharesOutstanding" | "marketCap" | "workingCapital" | "capitalEmployed" | "enterpriseValue" | "netDebt" | "netInterestIncome" | "commissionIncome" | "commissionExpense" | "netCommissionIncome" | "creditLossProvision">

export function buildAnnualSnapshots(
  rawData: RawData[],
  fields: NumericRawDataKey[] = ["revenue", "netProfit"],
): AnnualSnapshot[] {
  const idx: Record<number, Record<string, RawData>> = {}
  for (const row of rawData) {
    if (!idx[row.year]) idx[row.year] = {}
    idx[row.year][row.period] = row
  }

  const years = Object.keys(idx).map(Number).sort((a, b) => a - b)
  const result: (AnnualSnapshot & { _sortKey: number })[] = []

  for (const y of years) {
    const periods = idx[y]

    if (periods.YEAR) {
      result.push(makeEntry(`${y}`, y * 100 + 12, periods.YEAR, fields))
    }

    if (periods.Q2) {
      const prev = idx[y - 1]
      if (prev?.YEAR && prev?.Q2) {
        const ttm: AnnualSnapshot & { _sortKey: number } = {
          period: `${y} H1 TTM`,
          _sortKey: y * 100 + 6,
        }
        for (const f of fields) {
          const h2Prev = (val(prev.YEAR, f)) - (val(prev.Q2, f))
          const h1Curr = val(periods.Q2, f)
          ttm[f] = h2Prev + h1Curr
        }
        result.push(ttm)
      }
    }

    if (periods.Q3) {
      const prev = idx[y - 1]
      if (prev?.YEAR && prev?.Q3) {
        const ttm: AnnualSnapshot & { _sortKey: number } = {
          period: `${y} 9М TTM`,
          _sortKey: y * 100 + 9,
        }
        for (const f of fields) {
          const q4Prev = (val(prev.YEAR, f)) - (val(prev.Q3, f))
          const nineMCurr = val(periods.Q3, f)
          ttm[f] = q4Prev + nineMCurr
        }
        result.push(ttm)
      }
    }
  }

  result.sort((a, b) => a._sortKey - b._sortKey)

  return result.map(({ _sortKey, ...rest }) => rest)
}

const UNITS_MAP: Record<string, number> = {
  units: 1,
  thousands: 1000,
  millions: 1000000,
}

function multiplier(row: RawData): number {
  return (row.reportUnits && UNITS_MAP[row.reportUnits]) || 1
}

function val(row: RawData, field: NumericRawDataKey): number {
  const raw = (row[field] as number | null | undefined) ?? 0
  return raw * multiplier(row)
}

function makeEntry(
  period: string,
  sortKey: number,
  row: RawData,
  fields: NumericRawDataKey[],
): AnnualSnapshot & { _sortKey: number } {
  const entry: AnnualSnapshot & { _sortKey: number } = { period, _sortKey: sortKey }
  for (const f of fields) entry[f] = val(row, f)
  return entry
}
