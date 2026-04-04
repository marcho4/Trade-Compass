# AI Service Architecture Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Устранить архитектурные нарушения в ai-service: мёртвый код, нарушения SRP в репозиториях, конкретные типы вместо интерфейсов в адаптерах, дублирование параметра, лишний запрос в БД.

**Architecture:** Clean Architecture — usecase определяет интерфейсы, адаптеры и gateway зависят только от них. Каждый репозиторий отвечает за один домен.

**Tech Stack:** Go, pgx/v5, segmentio/kafka-go

---

## Затронутые файлы

| Действие | Файл |
|---|---|
| Удалить | `internal/repository/postgres/client.go` |
| Изменить | `internal/repository/postgres/analysis_repository.go` |
| Изменить | `internal/repository/postgres/report_results_repository.go` |
| Изменить | `internal/repository/postgres/task_repository.go` |
| Изменить | `internal/usecase/ai_service.go` |
| Изменить | `internal/usecase/business_research.go` |
| Изменить | `internal/gateway/gemini/ai_service.go` |
| Изменить | `internal/adapters/kafka/dispatcher.go` |
| Изменить | `internal/adapters/http/analysis.go` |
| Изменить | `internal/app/app.go` |

---

### Task 1: Удалить мёртвый `client.go`

**Files:**
- Delete: `internal/repository/postgres/client.go`

- [ ] **Step 1: Убедиться, что `DBRepo` нигде не используется**

```bash
cd /home/marcho/bull-run/ai-service && grep -r "DBRepo\|NewDBRepo" --include="*.go" .
```
Ожидаемый результат: только `client.go` упоминает эти символы.

- [ ] **Step 2: Удалить файл**

```bash
rm internal/repository/postgres/client.go
```

- [ ] **Step 3: Проверить компиляцию**

```bash
go build ./...
```

- [ ] **Step 4: Коммит**

```bash
git add -A && git commit -m "refactor: remove dead DBRepo from postgres package"
```

---

### Task 2: Исправить SRP в `analysis_repository.go`

Убрать методы чужих доменов из `AnalysisRepository` и добавить их в правильные репозитории.

**Files:**
- Modify: `internal/repository/postgres/analysis_repository.go`
- Modify: `internal/repository/postgres/report_results_repository.go`

- [ ] **Step 1: Добавить `GetReportResults` и `GetLatestReportResults` в `ReportResultsRepository`**

В файле `internal/repository/postgres/report_results_repository.go` добавить после `SaveReportResults`:

```go
func (r *ReportResultsRepository) GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error) {
	db := Executor(ctx, r.db)

	sql := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1 AND year = $2 AND period = $3
	`

	var res entity.ReportResults
	err := db.QueryRow(ctx, sql, ticker, year, period).Scan(
		&res.Health, &res.Growth, &res.Moat, &res.Dividends, &res.Value, &res.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: report results for %s year=%d period=%d", domain.ErrNotFound, ticker, year, period)
		}
		return nil, fmt.Errorf("get report results: %w", err)
	}

	return &res, nil
}

func (r *ReportResultsRepository) GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error) {
	db := Executor(ctx, r.db)

	sql := `
		SELECT health, growth, moat, dividends, value, total
		FROM report_results
		WHERE ticker = $1
		ORDER BY year DESC, period DESC
		LIMIT 1
	`

	var res entity.ReportResults
	err := db.QueryRow(ctx, sql, ticker).Scan(
		&res.Health, &res.Growth, &res.Moat, &res.Dividends, &res.Value, &res.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: latest report results for %s", domain.ErrNotFound, ticker)
		}
		return nil, fmt.Errorf("get latest report results: %w", err)
	}

	return &res, nil
}
```

Добавить нужные импорты в начало файла:
```go
import (
	"context"
	"errors"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
```

- [ ] **Step 2: Удалить из `analysis_repository.go` методы `GetReportResults`, `GetLatestReportResults`, `GetBusinessResearch`**

Финальный вид `analysis_repository.go` — только методы анализа:

```go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalysisRepository struct {
	db *pgxpool.Pool
}

func NewAnalysisRepository(db *pgxpool.Pool) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error) {
	db := Executor(ctx, r.db)

	sql := `SELECT analysis FROM analysis_reports WHERE ticker = $1 AND year = $2 AND period = $3`

	var analysis string
	err := db.QueryRow(ctx, sql, ticker, year, period).Scan(&analysis)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%w: analysis for %s year=%d period=%d", domain.ErrNotFound, ticker, year, period)
		}
		return "", fmt.Errorf("get analysis: %w", err)
	}

	return analysis, nil
}

func (r *AnalysisRepository) GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error) {
	db := Executor(ctx, r.db)

	sql := `SELECT year, period FROM analysis_reports WHERE ticker = $1 ORDER BY year DESC, period DESC`

	rows, err := db.Query(ctx, sql, ticker)
	if err != nil {
		return nil, fmt.Errorf("get available periods: %w", err)
	}
	defer rows.Close()

	var periods []entity.AvailablePeriod
	for rows.Next() {
		var p entity.AvailablePeriod
		if err := rows.Scan(&p.Year, &p.Period); err != nil {
			return nil, fmt.Errorf("scan period: %w", err)
		}
		periods = append(periods, p)
	}

	return periods, nil
}

func (r *AnalysisRepository) SaveAnalysis(ctx context.Context, result, ticker string, year, period int) error {
	db := Executor(ctx, r.db)

	sql := `
		INSERT INTO analysis_reports (ticker, year, period, analysis)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (ticker, year, period)
		DO UPDATE SET analysis = EXCLUDED.analysis
	`

	_, err := db.Exec(ctx, sql, ticker, year, period, result)
	if err != nil {
		return fmt.Errorf("save analysis: %w", err)
	}

	return nil
}
```

- [ ] **Step 3: Обновить `app.go` — использовать `reportResultsRepo` вместо `analysisRepo` для `GetReportResultsUsecase`**

В `internal/app/app.go` строка:
```go
reportResultsUC := usecase.NewGetReportResultsUsecase(analysisRepo)
```
заменить на:
```go
reportResultsUC := usecase.NewGetReportResultsUsecase(reportResultsRepo)
```

- [ ] **Step 4: Проверить компиляцию**

```bash
go build ./...
```

- [ ] **Step 5: Коммит**

```bash
git add -A && git commit -m "refactor: fix SRP in analysis_repository, move report results reads to ReportResultsRepository"
```

---

### Task 3: Исправить `task_repository.go` — убрать лишний `ensureTaskExists`

**Files:**
- Modify: `internal/repository/postgres/task_repository.go`

- [ ] **Step 1: Переписать `DecrementPending`**

Заменить метод `DecrementPending` и удалить `ensureTaskExists`:

```go
func (r *TasksRepository) DecrementPending(ctx context.Context, taskID string, taskType string) (int, error) {
	db := Executor(ctx, r.db)

	sql := `UPDATE tasks SET pending_count = pending_count - 1 WHERE task_id = $1 AND task_type = $2 RETURNING pending_count`

	var pendingCount int
	err := db.QueryRow(ctx, sql, taskID, taskType).Scan(&pendingCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("task not found")
		}
		return 0, fmt.Errorf("update tasks: %w", err)
	}

	return pendingCount, nil
}
```

Удалить метод `ensureTaskExists` целиком.

Финальный вид файла (без `ensureTaskExists`):

```go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TasksRepository struct {
	db *pgxpool.Pool
}

func NewTasksRepository(db *pgxpool.Pool) *TasksRepository {
	return &TasksRepository{
		db: db,
	}
}

func (r *TasksRepository) IncrementPending(ctx context.Context, taskID string, taskType string, count int) error {
	db := Executor(ctx, r.db)

	sql := `
        INSERT INTO tasks (task_id, task_type, pending_count)
        VALUES ($1, $2, $3)
        ON CONFLICT (task_id, task_type)
        DO UPDATE SET pending_count = tasks.pending_count + $3
    `
	_, err := db.Exec(ctx, sql, taskID, taskType, count)
	if err != nil {
		return fmt.Errorf("increment: %w", err)
	}

	return nil
}

func (r *TasksRepository) DecrementPending(ctx context.Context, taskID string, taskType string) (int, error) {
	db := Executor(ctx, r.db)

	sql := `UPDATE tasks SET pending_count = pending_count - 1 WHERE task_id = $1 AND task_type = $2 RETURNING pending_count`

	var pendingCount int
	err := db.QueryRow(ctx, sql, taskID, taskType).Scan(&pendingCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("task not found")
		}
		return 0, fmt.Errorf("update tasks: %w", err)
	}

	return pendingCount, nil
}

func (r *TasksRepository) CheckIfTaskIsReady(ctx context.Context, taskID string, expectedTasks int) (bool, error) {
	db := Executor(ctx, r.db)

	sql := `select count(*) from tasks where task_id = $1 and pending_count = 0`

	var count int
	err := db.QueryRow(ctx, sql, taskID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("select pending count: %w", err)
	}

	return count == expectedTasks, nil
}

func (r *TasksRepository) DeleteTask(ctx context.Context, taskID string) error {
	db := Executor(ctx, r.db)

	sql := `delete from tasks where task_id = $1`
	_, err := db.Exec(ctx, sql, taskID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
```

- [ ] **Step 2: Проверить компиляцию**

```bash
go build ./...
```

- [ ] **Step 3: Коммит**

```bash
git add -A && git commit -m "refactor: simplify DecrementPending, remove redundant ensureTaskExists"
```

---

### Task 4: Убрать `companyName` из `ResearchBusiness`

**Files:**
- Modify: `internal/usecase/ai_service.go`
- Modify: `internal/usecase/business_research.go`
- Modify: `internal/gateway/gemini/ai_service.go`

- [ ] **Step 1: Обновить интерфейс в `internal/usecase/ai_service.go`**

```go
type AIService interface {
	AnalyzeReport(ctx context.Context, ticker, reportURL string, year int, period entity.ReportPeriod) (string, error)
	ExtractRawData(ctx context.Context, ticker, reportURL string, year int, period entity.ReportPeriod) (*entity.RawData, error)
	ExtractResultFromReport(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.ReportResults, error)
	CollectNews(ctx context.Context, ticker string) (*entity.NewsResponse, error)
	ResearchBusiness(ctx context.Context, ticker string) (*entity.BusinessResearchResponse, error)
	ExtractRiskAndGrowth(ctx context.Context, ticker string, news *entity.NewsResponse, business *entity.BusinessResearchResult) (*entity.RiskAndGrowthResponse, error)
}
```

- [ ] **Step 2: Обновить вызов в `internal/usecase/business_research.go:43`**

```go
result, err := u.ai.ResearchBusiness(ctx, task.Ticker)
```

- [ ] **Step 3: Обновить реализацию в `internal/gateway/gemini/ai_service.go`**

Сигнатуру метода изменить на:
```go
func (s *AIService) ResearchBusiness(ctx context.Context, ticker string) (*entity.BusinessResearchResponse, error) {
	prompt := docs.BusinessResearcherPrompt() +
		"\n\n## Компания для анализа\nТикер: " + ticker +
		"\n\nВАЖНО: В поле ticker ответа используй СТРОГО \"" + ticker + "\". Не заменяй тикер на альтернативный."
    // остальное без изменений
```

- [ ] **Step 4: Проверить компиляцию**

```bash
go build ./...
```

- [ ] **Step 5: Коммит**

```bash
git add -A && git commit -m "refactor: remove unused companyName param from ResearchBusiness"
```

---

### Task 5: Интерфейсы в `dispatcher.go`

**Files:**
- Modify: `internal/adapters/kafka/dispatcher.go`

- [ ] **Step 1: Добавить интерфейс `TaskCounter` в `internal/usecase/task_counter.go`**

В конец файла `internal/usecase/task_counter.go` добавить:

```go
type TaskCounter interface {
	Increment(ctx context.Context, task entity.Task) error
	Decrement(ctx context.Context, task entity.Task) error
}
```

- [ ] **Step 2: Переписать `dispatcher.go`**

```go
package kafka

import (
	"context"
	"fmt"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
	"ai-service/internal/usecase"
)

type TaskDispatcher struct {
	handlers map[entity.TaskType]usecase.TaskExecutor
	counter  usecase.TaskCounter
}

func NewTaskDispatcher(
	analyzeReport usecase.TaskExecutor,
	extractRawData usecase.TaskExecutor,
	extractResult usecase.TaskExecutor,
	businessResearch usecase.TaskExecutor,
	newsResearch usecase.TaskExecutor,
	riskAndGrowth usecase.TaskExecutor,
	counter usecase.TaskCounter,
) *TaskDispatcher {
	handlers := map[entity.TaskType]usecase.TaskExecutor{
		entity.Analyze:          analyzeReport,
		entity.Extract:          extractRawData,
		entity.ExtractResult:    extractResult,
		entity.BusinessResearch: businessResearch,
		entity.NewsResearch:     newsResearch,
		entity.RiskAndGrowth:    riskAndGrowth,
	}

	return &TaskDispatcher{
		handlers: handlers,
		counter:  counter,
	}
}

func (d *TaskDispatcher) Dispatch(ctx context.Context, task entity.Task) error {
	switch task.Type {
	case entity.RawDataExpect, entity.RiskAndGrowthExpect:
		return d.counter.Increment(ctx, task)
	case entity.RawDataSuccess, entity.RiskAndGrowthSuccess:
		return d.counter.Decrement(ctx, task)
	}

	handler, ok := d.handlers[task.Type]
	if !ok {
		return fmt.Errorf("%w: %s", domain.ErrUnknownTaskType, task.Type)
	}

	return handler.Execute(ctx, task)
}
```

- [ ] **Step 3: Проверить компиляцию**

```bash
go build ./...
```

- [ ] **Step 4: Коммит**

```bash
git add -A && git commit -m "refactor: use interfaces in TaskDispatcher instead of concrete usecase types"
```

---

### Task 6: Интерфейсы в HTTP-хендлере

**Files:**
- Modify: `internal/adapters/http/analysis.go`

- [ ] **Step 1: Переписать `analysis.go` с интерфейсами**

```go
package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
)

type AnalysisReader interface {
	GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error)
	GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error)
}

type ReportResultsReader interface {
	GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error)
	GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error)
}

type BusinessResearchReader interface {
	GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error)
}

type analysisHandler struct {
	analysis         AnalysisReader
	reportResults    ReportResultsReader
	businessResearch BusinessResearchReader
}

func NewAnalysisHandler(
	analysis AnalysisReader,
	reportResults ReportResultsReader,
	businessResearch BusinessResearchReader,
) *analysisHandler {
	return &analysisHandler{
		analysis:         analysis,
		reportResults:    reportResults,
		businessResearch: businessResearch,
	}
}
```

Тела методов-хендлеров остаются без изменений.

- [ ] **Step 2: Проверить компиляцию**

```bash
go build ./...
```

- [ ] **Step 3: Коммит**

```bash
git add -A && git commit -m "refactor: use interfaces in HTTP handler instead of concrete usecase types"
```
