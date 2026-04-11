package usecase

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
	Ticker           string
	Year             int
	Period           entity.ReportPeriod
	RawDataHistory   []entity.RawData
	Candles          []entity.Candle
	CBRate           *entity.CBRate
	News             *entity.NewsResponse
	BusinessResearch *entity.BusinessResearchResult
	RisksAndGrowth   *entity.RiskAndGrowthResponse
	Scenarios        []entity.Scenario
	DCFResult        *entity.DCFResult
}

func buildNewsAgentPrompt(ticker string, dependencies []entity.CompanyDependency) string {
	prompt := docs.NewsCollectorAgent() + "\n\n## Текущая дата\n" + time.Now().Format("02.01.2006") + "\n\n## Тикер для анализа\n" + ticker

	if len(dependencies) > 0 {
		prompt += "\n\n## Зависимости компании из бизнес-анализа\n"
		for _, d := range dependencies {
			prompt += fmt.Sprintf("- %s [тип: %s, критичность: %s]: %s\n", d.Factor, d.Type, d.Severity, d.Description)
		}
	}

	return prompt
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
	writePrecomputedDCF(&b, ctx.Scenarios, ctx.DCFResult)
	writeMarketData(&b, ctx.CBRate)
	writePriceHistory(&b, ctx.Candles)
	writeNews(&b, ctx.News)

	if ctx.RisksAndGrowth != nil {
		b.WriteString(ctx.RisksAndGrowth.String())
	}
	if ctx.BusinessResearch != nil {
		b.WriteString(ctx.BusinessResearch.String())
	}

	slog.Debug("Full Prompt", slog.String("Prompt", b.String()))
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
		case "billions":
			units = "млрд руб."
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

func writeMarketData(b *strings.Builder, rate *entity.CBRate) {
	b.WriteString("<market_data>\n")

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

func writePrecomputedDCF(b *strings.Builder, scenarios []entity.Scenario, dcf *entity.DCFResult) {
	b.WriteString("<precomputed_dcf>\n")

	if len(scenarios) == 0 || dcf == nil {
		b.WriteString("Заранее рассчитанный DCF отсутствует. Явно укажи в отчёте невозможность определить Fair Value и не пытайся рассчитать DCF самостоятельно.\n")
		b.WriteString("</precomputed_dcf>\n\n")
		return
	}

	b.WriteString("ЭТО AUTHORITATIVE-ИСТОЧНИК СПРАВЕДЛИВОЙ СТОИМОСТИ. DCF, WACC, FCFF, терминал и цены за акцию уже рассчитаны внешней моделью. Пересчитывать, корректировать или \"уточнять\" эти числа ЗАПРЕЩЕНО. Вероятности сценариев тоже фиксированы — не меняй их. Твоя задача — интерпретировать сценарии и их допущения, а не переоценивать.\n\n")

	fmt.Fprintf(b, "Взвешенная цена за акцию (WeightedPrice): %.2f руб.\n", dcf.WeightedPrice)
	fmt.Fprintf(b, "Взвешенный Enterprise Value (WeightedEV): %.0f руб.\n", dcf.WeightedEV)
	fmt.Fprintf(b, "Количество сценариев: %d\n\n", len(dcf.Scenarios))

	scenariosByID := make(map[string]entity.Scenario, len(scenarios))
	for _, s := range scenarios {
		scenariosByID[s.ID] = s
	}

	for _, sr := range dcf.Scenarios {
		scenario, ok := scenariosByID[sr.ScenarioID]
		if !ok {
			slog.Error("Scenario not found for dcf result", slog.String("scenario_id", sr.ScenarioID))
			continue
		}

		b.WriteString("---- Сценарий ----\n")
		b.WriteString(scenario.String())
		b.WriteString(sr.String())
		b.WriteString("------------------\n\n")
	}

	b.WriteString("</precomputed_dcf>\n\n")
}
