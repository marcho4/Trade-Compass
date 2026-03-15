package application

import (
	"ai-service/internal/domain/entity"
	docs "ai-service/internal/infrastructure/docs"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type AnalysisContext struct {
	Ticker         string
	Year           int
	Period         entity.ReportPeriod
	RawDataHistory []entity.RawData
	Candles        []entity.Candle
	CBRate         *entity.CBRate
	MarketCap      float64
	News           *entity.NewsResponse
}

func BuildNewsAgentPrompt(ticker string) string {
	return docs.NewsCollectorAgent() + "\n\n" + "## Тикер для анализа\n" + ticker
}

func BuildExtractPrompt(reportText string) string {
	return docs.ExtractPrompt() + "\n\n" + reportText
}

func BuildAnalysisPrompt(ctx AnalysisContext) string {
	var b strings.Builder

	writeRole(&b, ctx)
	writeAnalysisMethodology(&b)
	writeMacroContext(&b)
	writeFinancialHistory(&b, ctx.RawDataHistory)
	writeMarketData(&b, ctx.CBRate, ctx.MarketCap)
	writePriceHistory(&b, ctx.Candles)
	writeNews(&b, ctx.News)
	slog.Info("Full Prompt", slog.String("Prompt", b.String()))
	return b.String()
}

func writeRole(b *strings.Builder, ctx AnalysisContext) {
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
		json, _ := json.Marshal(rd)
		b.WriteString(string(json))
		b.WriteString("\n")
	}

	b.WriteString("</financial_data>\n\n")
}

func writeMetric(b *strings.Builder, name string, val *int64) {
	if val != nil {
		fmt.Fprintf(b, "%s: %d\n", name, *val)
	}
}

func writeFloatMetric(b *strings.Builder, name string, val *float64) {
	if val != nil {
		fmt.Fprintf(b, "%s: %.2f\n", name, *val)
	}
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
	json, _ := json.Marshal(news)
	b.WriteString(string(json))
	b.WriteString("\n</news>\n\n")
}
