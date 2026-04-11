"use client"

import { Fragment, useMemo } from "react"
import { buildAnnualSnapshots, NumericRawDataKey } from "@/lib/build-annual-snapshots"
import { formatShort } from "@/lib/utils"
import {
  RawData,
  MetricFieldConfig,
  BANK_FIELDS,
  PNL_FIELDS,
  BALANCE_SHEET_FIELDS,
  CASH_FLOW_FIELDS,
  MARKET_DATA_FIELDS,
  CALCULATED_FIELDS,
} from "@/types/raw-data"

interface RawDataTableProps {
  data: RawData[]
}

const ALL_NUMERIC_FIELDS = [
  ...BANK_FIELDS,
  ...PNL_FIELDS,
  ...BALANCE_SHEET_FIELDS,
  ...CASH_FLOW_FIELDS,
  ...MARKET_DATA_FIELDS,
  ...CALCULATED_FIELDS,
] as const

function detectCompanyType(data: RawData[]): string | null {
  for (const row of data) {
    if (row.companyType) return row.companyType
  }
  return null
}

interface Section {
  title: string
  fields: MetricFieldConfig[]
}

export const RawDataTable = ({ data }: RawDataTableProps) => {
  const companyType = useMemo(() => detectCompanyType(data), [data])

  const allKeys = useMemo(
    () => ALL_NUMERIC_FIELDS.map((f) => f.key) as NumericRawDataKey[],
    [],
  )

  const snapshots = useMemo(
    () => buildAnnualSnapshots(data, allKeys),
    [data, allKeys],
  )

  const periods = snapshots.map((s) => s.period)

  const sections: Section[] = useMemo(() => {
    const result: Section[] = []
    if (companyType === "bank") {
      result.push({ title: "Банковские показатели", fields: BANK_FIELDS })
    }
    result.push(
      { title: "Отчёт о прибылях и убытках", fields: PNL_FIELDS },
      { title: "Баланс", fields: BALANCE_SHEET_FIELDS },
      { title: "Денежный поток", fields: CASH_FLOW_FIELDS },
      { title: "Рыночные данные", fields: MARKET_DATA_FIELDS },
      { title: "Расчётные показатели", fields: CALCULATED_FIELDS },
    )
    return result
  }, [companyType])

  if (snapshots.length === 0) return null

  return (
    <div className="overflow-x-auto rounded-lg border border-border">
      <table className="w-full min-w-max text-sm">
        <thead>
          <tr className="border-b border-border bg-muted/50">
            <th className="sticky left-0 z-10 bg-muted/50 px-4 py-3 text-left font-medium text-muted-foreground min-w-[220px]">
              Метрика
            </th>
            {periods.map((period) => (
              <th
                key={period}
                className="px-4 py-3 text-right font-medium text-muted-foreground whitespace-nowrap min-w-[100px]"
              >
                {period}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {sections.map((section) => (
            <Fragment key={section.title}>
              <tr className="border-t border-border">
                <td
                  colSpan={1 + periods.length}
                  className="bg-muted/30 px-4 py-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground"
                >
                  {section.title}
                </td>
              </tr>
              {section.fields.map((field, rowIndex) => (
                <tr
                  key={field.key}
                  className={rowIndex % 2 === 1 ? "bg-muted/20" : undefined}
                >
                  <td className="sticky left-0 z-10 bg-background px-4 py-2.5 text-foreground">
                    {field.label}
                  </td>
                  {snapshots.map((snapshot) => {
                    const raw = snapshot[field.key as string]
                    const value = typeof raw === "number" ? raw : null
                    return (
                      <td
                        key={snapshot.period}
                        className="px-4 py-2.5 text-right font-mono text-foreground tabular-nums"
                      >
                        {formatShort(value)}
                      </td>
                    )
                  })}
                </tr>
              ))}
            </Fragment>
          ))}
        </tbody>
      </table>
    </div>
  )
}
