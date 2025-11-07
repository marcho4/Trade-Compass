"use client"

import { Card } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { FinancialMetrics } from "@/types"

interface FinancialStatementsProps {
  metrics: FinancialMetrics
}

export const FinancialStatements = ({ metrics }: FinancialStatementsProps) => {
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("ru-RU", {
      style: "currency",
      currency: "RUB",
      maximumFractionDigits: 0,
    }).format(value)
  }

  const StatementRow = ({
    label,
    value,
    bold = false,
    indent = 0,
  }: {
    label: string
    value: number
    bold?: boolean
    indent?: number
  }) => (
    <div
      className={`flex justify-between py-2 border-b border-border/50 ${
        bold ? "font-semibold" : ""
      }`}
      style={{ paddingLeft: `${indent * 16}px` }}
    >
      <span className={!bold ? "text-muted-foreground" : ""}>{label}</span>
      <span>{formatCurrency(value)}</span>
    </div>
  )

  return (
    <Card className="p-6">
      <h3 className="text-xl font-semibold mb-6">Финансовая отчетность</h3>
      <Tabs defaultValue="pnl" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="pnl">Отчёт о прибылях и убытках</TabsTrigger>
          <TabsTrigger value="balance">Баланс</TabsTrigger>
          <TabsTrigger value="cashflow">Движение денежных средств</TabsTrigger>
        </TabsList>

        <TabsContent value="pnl" className="space-y-2 mt-6">
          <StatementRow label="Выручка" value={metrics.revenue} bold />
          <StatementRow
            label="Себестоимость"
            value={-metrics.costOfRevenue}
            indent={1}
          />
          <StatementRow label="Валовая прибыль" value={metrics.grossProfit} bold />
          <StatementRow
            label="Операционные расходы"
            value={-metrics.operatingExpenses}
            indent={1}
          />
          <StatementRow
            label="Операционная прибыль (EBIT)"
            value={metrics.ebit}
            bold
          />
          <StatementRow label="EBITDA" value={metrics.ebitda} bold />
          <StatementRow
            label="Процентные расходы"
            value={-metrics.interestExpense}
            indent={1}
          />
          <StatementRow label="Налог на прибыль" value={-metrics.taxExpense} indent={1} />
          <StatementRow label="Чистая прибыль" value={metrics.netProfit} bold />
          <div className="pt-4 mt-4 border-t">
            <div className="text-sm text-muted-foreground">
              Период: {metrics.reportPeriod} {metrics.reportYear}
            </div>
          </div>
        </TabsContent>

        <TabsContent value="balance" className="space-y-2 mt-6">
          <div className="mb-4">
            <h4 className="font-semibold text-lg mb-2">Активы</h4>
            <StatementRow label="Оборотные активы" value={metrics.currentAssets} />
            <StatementRow
              label="Денежные средства"
              value={metrics.cashAndEquivalents}
              indent={1}
            />
            <StatementRow label="Дебиторская задолженность" value={metrics.receivables} indent={1} />
            <StatementRow label="Запасы" value={metrics.inventories} indent={1} />
            <StatementRow label="Всего активов" value={metrics.totalAssets} bold />
          </div>

          <div className="mb-4 pt-4 border-t">
            <h4 className="font-semibold text-lg mb-2">Обязательства</h4>
            <StatementRow
              label="Краткосрочные обязательства"
              value={metrics.currentLiabilities}
            />
            <StatementRow label="Краткосрочный долг" value={metrics.shortTermDebt} indent={1} />
            <StatementRow label="Долгосрочный долг" value={metrics.longTermDebt} />
            <StatementRow label="Всего обязательств" value={metrics.totalLiabilities} bold />
          </div>

          <div className="pt-4 border-t">
            <h4 className="font-semibold text-lg mb-2">Капитал</h4>
            <StatementRow label="Собственный капитал" value={metrics.equity} bold />
            <StatementRow
              label="Нераспределённая прибыль"
              value={metrics.retainedEarnings}
              indent={1}
            />
          </div>

          <div className="pt-4 mt-4 border-t">
            <div className="text-sm text-muted-foreground">
              На дату: {metrics.reportPeriod} {metrics.reportYear}
            </div>
          </div>
        </TabsContent>

        <TabsContent value="cashflow" className="space-y-2 mt-6">
          <StatementRow
            label="Денежный поток от операционной деятельности"
            value={metrics.operatingCashFlow}
            bold
          />
          <StatementRow
            label="Денежный поток от инвестиционной деятельности"
            value={metrics.investingCashFlow}
          />
          <StatementRow label="Капитальные затраты (CAPEX)" value={-metrics.capex} indent={1} />
          <StatementRow
            label="Денежный поток от финансовой деятельности"
            value={metrics.financingCashFlow}
          />
          <StatementRow label="Свободный денежный поток (FCF)" value={metrics.freeCashFlow} bold />
          <div className="pt-4 mt-4 border-t">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Оборотный капитал</p>
                <p className="text-lg font-semibold">{formatCurrency(metrics.workingCapital)}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Чистый долг</p>
                <p className="text-lg font-semibold">{formatCurrency(metrics.netDebt)}</p>
              </div>
            </div>
          </div>
          <div className="pt-4 border-t">
            <div className="text-sm text-muted-foreground">
              Период: {metrics.reportPeriod} {metrics.reportYear}
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </Card>
  )
}

