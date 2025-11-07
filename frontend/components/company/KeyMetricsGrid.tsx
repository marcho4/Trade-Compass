"use client"

import { FinancialIndicators } from "@/types"
import { METRIC_DESCRIPTIONS } from "@/types"
import { MetricCard } from "./MetricCard"

interface KeyMetricsGridProps {
  indicators: FinancialIndicators
  industryAverages?: Partial<FinancialIndicators>
}

export const KeyMetricsGrid = ({
  indicators,
  industryAverages,
}: KeyMetricsGridProps) => {
  return (
    <div className="space-y-8">
      {/* Мультипликаторы */}
      <div>
        <h3 className="text-xl font-semibold mb-4">Мультипликаторы оценки</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <MetricCard
            label="P/E"
            value={indicators.pe}
            format="ratio"
            description={METRIC_DESCRIPTIONS.pe}
            comparisonValue={industryAverages?.pe}
          />
          <MetricCard
            label="P/B"
            value={indicators.pb}
            format="ratio"
            comparisonValue={industryAverages?.pb}
          />
          <MetricCard
            label="EV/EBITDA"
            value={indicators.evEbitda}
            format="ratio"
            description={METRIC_DESCRIPTIONS.evEbitda}
            comparisonValue={industryAverages?.evEbitda}
          />
          <MetricCard
            label="P/S"
            value={indicators.ps}
            format="ratio"
            comparisonValue={industryAverages?.ps}
          />
        </div>
      </div>

      {/* Рентабельность */}
      <div>
        <h3 className="text-xl font-semibold mb-4">Рентабельность</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <MetricCard
            label="ROE"
            value={indicators.roe}
            format="percent"
            description={METRIC_DESCRIPTIONS.roe}
            comparisonValue={industryAverages?.roe}
          />
          <MetricCard
            label="ROA"
            value={indicators.roa}
            format="percent"
            comparisonValue={industryAverages?.roa}
          />
          <MetricCard
            label="Чистая маржа"
            value={indicators.netProfitMargin}
            format="percent"
            description={METRIC_DESCRIPTIONS.netProfitMargin}
            comparisonValue={industryAverages?.netProfitMargin}
          />
          <MetricCard
            label="Валовая маржа"
            value={indicators.grossProfitMargin}
            format="percent"
            comparisonValue={industryAverages?.grossProfitMargin}
          />
        </div>
      </div>

      {/* Ликвидность и долг */}
      <div>
        <h3 className="text-xl font-semibold mb-4">Ликвидность и долговая нагрузка</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <MetricCard
            label="Current Ratio"
            value={indicators.currentRatio}
            format="ratio"
            description={METRIC_DESCRIPTIONS.currentRatio}
            comparisonValue={industryAverages?.currentRatio}
          />
          <MetricCard
            label="Quick Ratio"
            value={indicators.quickRatio}
            format="ratio"
            comparisonValue={industryAverages?.quickRatio}
          />
          <MetricCard
            label="Debt/Equity"
            value={indicators.debtToEquity}
            format="ratio"
            description={METRIC_DESCRIPTIONS.debtToEquity}
            comparisonValue={industryAverages?.debtToEquity}
          />
          <MetricCard
            label="Debt/EBITDA"
            value={indicators.debtToEbitda}
            format="ratio"
            comparisonValue={industryAverages?.debtToEbitda}
          />
        </div>
      </div>

      {/* Денежные потоки */}
      <div>
        <h3 className="text-xl font-semibold mb-4">Денежные потоки и эффективность</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <MetricCard
            label="FCF Yield"
            value={indicators.fcfYield}
            format="percent"
            description={METRIC_DESCRIPTIONS.fcfYield}
            comparisonValue={industryAverages?.fcfYield}
          />
          <MetricCard
            label="Income Quality"
            value={indicators.incomeQuality}
            format="ratio"
            comparisonValue={industryAverages?.incomeQuality}
          />
          <MetricCard
            label="Interest Coverage"
            value={indicators.interestCoverage}
            format="ratio"
            comparisonValue={industryAverages?.interestCoverage}
          />
          <MetricCard
            label="Дивидендная доходность"
            value={indicators.dividendYield}
            format="percent"
            comparisonValue={industryAverages?.dividendYield}
          />
        </div>
      </div>
    </div>
  )
}

