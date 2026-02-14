package application

import (
	"ai-service/internal/domain"
	docs "ai-service/internal/infrastructure/docs"
	"fmt"
	"strings"
	"time"
)

type AnalysisContext struct {
	Ticker    string
	Year      int
	Period    domain.ReportPeriod
	RawData   *domain.RawData
	Candles   []domain.Candle
	CBRate    *domain.CBRate
	MarketCap float64
}

func BuildAnalysisPrompt(ctx AnalysisContext) string {
	var b strings.Builder

	writeRole(&b, ctx)
	writeAnalysisMethodology(&b)
	writeMacroContext(&b)
	// writeFinancialReport(&b, ctx.RawData)
	writeMarketData(&b, ctx.CBRate, ctx.MarketCap)
	writePriceHistory(&b, ctx.Candles)

	return b.String()
}

func writeRole(b *strings.Builder, ctx AnalysisContext) {
	fmt.Fprintf(b, `Ты — старший инвестиционный аналитик с 15-летним опытом работы на российском фондовом рынке (MOEX).
Твоя задача — провести комплексный фундаментальный анализ компании %s на основе предоставленных данных
и сформировать структурированный аналитический отчёт по методологии Morningstar, адаптированной для российского рынка.

КРИТИЧЕСКИ ВАЖНЫЕ ПРАВИЛА:
- Используй ТОЛЬКО предоставленные данные. Не выдумывай и не додумывай цифры.
- Если данных недостаточно для расчёта — явно укажи это и объясни, какие данные необходимы.
- Все расчёты должны быть прозрачными — показывай формулы и промежуточные вычисления.
- Каждый вывод должен быть подкреплён конкретными цифрами из предоставленных данных.
- Дата анализа: %s

`, ctx.Ticker, time.Now().Format("02.01.2006"))
}

func writeFinancialReport(b *strings.Builder, rd *domain.RawData) {
	if rd == nil {
		b.WriteString("<financial_report>\nДанные финансовой отчётности не предоставлены.\n</financial_report>\n\n")
		return
	}

	b.WriteString("<financial_report>\n")
	fmt.Fprintf(b, "Тикер: %s\n", rd.Ticker)
	fmt.Fprintf(b, "Период: %d / %s\n\n", rd.Year, rd.Period)

	b.WriteString("=== ОТЧЁТ О ПРИБЫЛЯХ И УБЫТКАХ (тыс. руб.) ===\n")
	writeMetric(b, "Выручка", rd.Revenue)
	writeMetric(b, "Себестоимость продаж", rd.CostOfRevenue)
	writeMetric(b, "Валовая прибыль", rd.GrossProfit)
	writeMetric(b, "Операционные расходы", rd.OperatingExpenses)
	writeMetric(b, "EBIT (операционная прибыль)", rd.EBIT)
	writeMetric(b, "EBITDA", rd.EBITDA)
	writeMetric(b, "Проценты к уплате", rd.InterestExpense)
	writeMetric(b, "Налог на прибыль", rd.TaxExpense)
	writeMetric(b, "Чистая прибыль", rd.NetProfit)

	b.WriteString("\n=== БАЛАНС (тыс. руб.) ===\n")
	writeMetric(b, "Итого активы", rd.TotalAssets)
	writeMetric(b, "Оборотные активы", rd.CurrentAssets)
	writeMetric(b, "Денежные средства", rd.CashAndEquivalents)
	writeMetric(b, "Запасы", rd.Inventories)
	writeMetric(b, "Дебиторская задолженность", rd.Receivables)
	writeMetric(b, "Итого обязательства", rd.TotalLiabilities)
	writeMetric(b, "Краткосрочные обязательства", rd.CurrentLiabilities)
	writeMetric(b, "Общий долг", rd.Debt)
	writeMetric(b, "Долгосрочный долг", rd.LongTermDebt)
	writeMetric(b, "Краткосрочный долг", rd.ShortTermDebt)
	writeMetric(b, "Собственный капитал", rd.Equity)
	writeMetric(b, "Нераспределённая прибыль", rd.RetainedEarnings)

	b.WriteString("\n=== ДЕНЕЖНЫЕ ПОТОКИ (тыс. руб.) ===\n")
	writeMetric(b, "Операционный денежный поток (OCF)", rd.OperatingCashFlow)
	writeMetric(b, "Инвестиционный денежный поток", rd.InvestingCashFlow)
	writeMetric(b, "Финансовый денежный поток", rd.FinancingCashFlow)
	writeMetric(b, "Капитальные затраты (CAPEX)", rd.CAPEX)
	writeMetric(b, "Свободный денежный поток (FCF)", rd.FreeCashFlow)

	b.WriteString("\n=== ПРОИЗВОДНЫЕ ПОКАЗАТЕЛИ (тыс. руб.) ===\n")
	writeMetric(b, "Оборотный капитал", rd.WorkingCapital)
	writeMetric(b, "Задействованный капитал", rd.CapitalEmployed)
	writeMetric(b, "Чистый долг", rd.NetDebt)
	writeMetric(b, "Стоимость предприятия (EV)", rd.EnterpriseValue)

	b.WriteString("</financial_report>\n\n")
}

func writeMetric(b *strings.Builder, name string, val *int64) {
	if val != nil {
		fmt.Fprintf(b, "%s: %d\n", name, *val)
	}
}

func writeMarketData(b *strings.Builder, rate *domain.CBRate, marketCap float64) {
	b.WriteString("<market_data>\n")

	if marketCap > 0 {
		fmt.Fprintf(b, "Рыночная капитализация: %.0f руб.\n", marketCap)
	} else {
		b.WriteString("Рыночная капитализация: нет данных\n")
	}

	if rate != nil {
		fmt.Fprintf(b, "Ключевая ставка ЦБ РФ: %.2f%% (дата: %s)\n", rate.Rate, rate.Date.Format("02.01.2006"))
	} else {
		b.WriteString("Ключевая ставка ЦБ РФ: нет данных\n")
	}

	b.WriteString("</market_data>\n\n")
}

func writePriceHistory(b *strings.Builder, candles []domain.Candle) {
	b.WriteString("<price_history>\n")

	if len(candles) == 0 {
		b.WriteString("История цен не предоставлена.\n")
		b.WriteString("</price_history>\n\n")
		return
	}

	b.WriteString("История цен за последние 12 месяцев (дневные свечи):\n")
	b.WriteString("Дата       | Открытие | Закрытие | Макс   | Мин    | Объём\n")
	b.WriteString("-----------|----------|----------|--------|--------|----------\n")

	for _, c := range candles {
		date := c.Begin
		if len(date) > 10 {
			date = date[:10]
		}
		fmt.Fprintf(b, "%s | %8.2f | %8.2f | %6.2f | %6.2f | %.0f\n",
			date, c.Open, c.Close, c.High, c.Low, c.Volume)
	}

	b.WriteString("</price_history>\n\n")
}

func writeAnalysisMethodology(b *strings.Builder) {
	b.WriteString("<analysis_methodology>\n")
	b.WriteString("Используй следующую методологию как справочник для проведения анализа.\n")
	b.WriteString("Все определения, шкалы, пороговые значения и формулы бери СТРОГО отсюда.\n\n")
	b.WriteString(docs.AnalysisFramework())
	b.WriteString("\n</analysis_methodology>\n\n")
}

func writeMacroContext(b *strings.Builder) {
	b.WriteString("<macro_context>\n")
	b.WriteString(docs.RussianHistory())
	b.WriteString("\n</macro_context>\n\n")
}
