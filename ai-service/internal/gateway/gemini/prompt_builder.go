package gemini

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	docs "ai-service/internal/docs"
	"ai-service/internal/domain/entity"
)

type analysisContext struct {
	Ticker              string
	Year                int
	Period              entity.ReportPeriod
	RawDataHistory      []entity.RawData
	Candles             []entity.Candle
	CBRate              *entity.CBRate
	MarketCap           float64
	News                *entity.NewsResponse
	BusinessResearch    entity.BusinessResearchResult
	RisksAndGrowth      entity.RiskAndGrowthResponse
	Scenarios           []entity.Scenario
	DCFScenariosResults []entity.DCFResult
}

func buildNewsAgentPrompt(ticker string) string {
	return docs.NewsCollectorAgent() + "\n\n" + "## Тикер для анализа\n" + ticker
}

func buildExtractPrompt(reportText string) string {
	return docs.ExtractPrompt() + "\n\n" + reportText
}

func buildAnalysisPrompt(ctx analysisContext) string {
	var b strings.Builder

	writeRole(&b, ctx)
	writeAnalysisMethodology(&b)
	writeMacroContext(&b)
	writeFinancialHistory(&b, ctx.RawDataHistory)
	writeMarketData(&b, ctx.CBRate, ctx.MarketCap)
	writePriceHistory(&b, ctx.Candles)
	writeNews(&b, ctx.News)

	b.WriteString(ctx.RisksAndGrowth.String())
	b.WriteString(ctx.BusinessResearch.String())
	writeScenarios(&b, ctx.Scenarios, ctx.DCFScenariosResults)

	slog.Info("Full Prompt", slog.String("Prompt", b.String()))
	return b.String()
}

func writeRole(b *strings.Builder, ctx analysisContext) {
	fmt.Fprintf(b,
		`Ты — старший инвестиционный аналитик с 15-летним опытом работы на российском фондовом рынке (MOEX).
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

func writeFinancialHistory(b *strings.Builder, history []entity.RawData) {
	b.WriteString("<financial_data>\n")

	if len(history) == 0 {
		b.WriteString("Исторические финансовые данные не предоставлены.\n")
		b.WriteString("</financial_data>\n\n")
		return
	}

	for i, rd := range history {
		if i > 0 {
			b.WriteString("---\n\n")
		}

		units := "тыс. руб."
		switch rd.ReportUnits {
		case "millions":
			units = "млн руб."
		case "units":
			units = "руб."
		}

		fmt.Fprintf(b, "## Период: %d / %s (%s)\n\n", rd.Year, rd.Period, units)
		data, _ := json.Marshal(rd)
		b.WriteString(string(data))
		b.WriteString("\n")
	}

	b.WriteString("</financial_data>\n\n")
}

func writeMarketData(b *strings.Builder, rate *entity.CBRate, marketCap float64) {
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

func writePriceHistory(b *strings.Builder, candles []entity.Candle) {
	b.WriteString("<price_history>\n")

	if len(candles) == 0 {
		b.WriteString("История цен не предоставлена.\n")
		b.WriteString("</price_history>\n\n")
		return
	}

	b.WriteString("Для анализа используй цены только отсюда.")
	b.WriteString("При расчете оценки акции бери цену за дату генерации отчета, или за последний самый близкий к дате отчета день из таблицы")
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

func writeNews(b *strings.Builder, news *entity.NewsResponse) {
	b.WriteString("<news>\n")
	data, _ := json.Marshal(news)
	b.WriteString(string(data))
	b.WriteString("\n</news>\n\n")
}

func writeScenarios(b *strings.Builder, scenarios []entity.Scenario, dcf []entity.DCFResult) {
	m := make(map[string]entity.Scenario)
	for _, i := range scenarios {
		m[i.Name] = i
	}

	for _, x := range dcf {
		scenario, ok := m[x.ID]
		if !ok {
			slog.Error("Scenario not found for dcf", "dcf", dcf)
		}

		b.WriteString("---- Сценарий ------\n")
		b.WriteString(scenario.String() + "\n")
		b.WriteString(x.String() + "\n")
		b.WriteString("----------\n")
	}
}
