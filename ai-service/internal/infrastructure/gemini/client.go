package gemini

import (
	"ai-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	model  string
}

func NewClient(apiKey string) (*Client, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &Client{
		client: client,
		model:  "gemini-2.0-flash",
	}, nil
}

const extractionPrompt = `Ты финансовый аналитик. Извлеки из этого финансового отчета (PDF) следующие метрики.
Все значения должны быть в тысячах рублей (если в отчете млн — умножь на 1000, если в рублях — раздели на 1000).
Если метрика не найдена в отчете, ставь null.

Верни ТОЛЬКО валидный JSON объект (без markdown, без комментариев) со следующими полями:

{
  "revenue": <int64 или null>,
  "costOfRevenue": <int64 или null>,
  "grossProfit": <int64 или null>,
  "operatingExpenses": <int64 или null>,
  "ebit": <int64 или null>,
  "ebitda": <int64 или null>,
  "interestExpense": <int64 или null>,
  "taxExpense": <int64 или null>,
  "netProfit": <int64 или null>,
  "totalAssets": <int64 или null>,
  "currentAssets": <int64 или null>,
  "cashAndEquivalents": <int64 или null>,
  "inventories": <int64 или null>,
  "receivables": <int64 или null>,
  "totalLiabilities": <int64 или null>,
  "currentLiabilities": <int64 или null>,
  "debt": <int64 или null>,
  "longTermDebt": <int64 или null>,
  "shortTermDebt": <int64 или null>,
  "equity": <int64 или null>,
  "retainedEarnings": <int64 или null>,
  "operatingCashFlow": <int64 или null>,
  "investingCashFlow": <int64 или null>,
  "financingCashFlow": <int64 или null>,
  "capex": <int64 или null>,
  "freeCashFlow": <int64 или null>,
  "sharesOutstanding": <int64 или null>,
  "marketCap": <int64 или null>,
  "workingCapital": <int64 или null>,
  "capitalEmployed": <int64 или null>,
  "enterpriseValue": <int64 или null>,
  "netDebt": <int64 или null>
}

Подсказки:
- revenue = Выручка
- costOfRevenue = Себестоимость продаж
- grossProfit = Валовая прибыль (revenue - costOfRevenue)
- operatingExpenses = Операционные расходы (коммерческие + управленческие + прочие)
- ebit = Прибыль от продаж / операционная прибыль
- ebitda = EBITDA (ebit + амортизация)
- interestExpense = Проценты к уплате
- taxExpense = Налог на прибыль
- netProfit = Чистая прибыль
- totalAssets = Итого активы (баланс)
- currentAssets = Оборотные активы
- cashAndEquivalents = Денежные средства и денежные эквиваленты
- inventories = Запасы
- receivables = Дебиторская задолженность
- totalLiabilities = Итого обязательства
- currentLiabilities = Краткосрочные обязательства
- debt = Долгосрочные + краткосрочные заемные средства
- longTermDebt = Долгосрочные заемные средства
- shortTermDebt = Краткосрочные заемные средства
- equity = Собственный капитал (Итого капитал)
- retainedEarnings = Нераспределённая прибыль
- operatingCashFlow = Чистые денежные средства от текущих операций
- investingCashFlow = Чистые денежные средства от инвестиционных операций
- financingCashFlow = Чистые денежные средства от финансовых операций
- capex = Капитальные затраты (приобретение основных средств, обычно отрицательное число — верни как положительное)
- freeCashFlow = operatingCashFlow - capex
- sharesOutstanding = Количество акций в обращении (если есть в отчете)
- workingCapital = currentAssets - currentLiabilities
- capitalEmployed = totalAssets - currentLiabilities
- netDebt = debt - cashAndEquivalents
`

func (c *Client) ExtractFromPDF(ctx context.Context, pdfBytes []byte, ticker string, year int, period string) (*domain.RawData, error) {
	log.Printf("Extracting financial data from PDF for %s %d/%s", ticker, year, period)

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				genai.NewPartFromBytes(pdfBytes, "application/pdf"),
				genai.NewPartFromText(extractionPrompt),
			},
		},
	}

	result, err := c.client.Models.GenerateContent(ctx, c.model, contents, &genai.GenerateContentConfig{
		Temperature:     genai.Ptr(float32(0.1)),
		TopP:            genai.Ptr(float32(0.95)),
		MaxOutputTokens: 4096,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini API call failed: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini returned empty response")
	}

	text := result.Candidates[0].Content.Parts[0].Text
	text = cleanJSONResponse(text)

	var rawData domain.RawData
	if err := json.Unmarshal([]byte(text), &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response as JSON: %w\nresponse: %s", err, text)
	}

	rawData.Ticker = ticker
	rawData.Year = year

	mappedPeriod, ok := domain.MonthsToPeriod[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s", period)
	}
	rawData.Period = mappedPeriod
	rawData.Status = "draft"

	return &rawData, nil
}

func cleanJSONResponse(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	} else if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}
	return text
}
