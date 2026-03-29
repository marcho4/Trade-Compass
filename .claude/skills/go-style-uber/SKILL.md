---
name: go-style-uber
description: Enforce Uber Go Style Guide conventions when writing, reviewing, or refactoring Go code. Use this skill whenever the user asks to write Go code, create Go files, review Go code, refactor Go, build Go microservices, or mentions golang, .go files, Go handlers, Go tests. Also trigger when the user says "write it clean", "production-quality Go", "idiomatic Go", or asks about Go best practices.
---

# Uber Go Style Guide — Agent Instructions

Apply every rule below when writing or reviewing Go code. These rules are distilled from the Uber Go Style Guide and adapted for agent use. When in doubt, prefer clarity over cleverness.

Before writing any non-trivial Go code, load the relevant reference file:
- Error handling, defer, panic → read `references/error-handling.md`
- Structs, interfaces, enums, functional options, naming → read `references/patterns.md`
- Performance, concurrency, goroutines, channels → read `references/performance.md`

---

## Core Principles

1. **Clarity over cleverness** — code is read far more than it is written.
2. **Consistency** — follow existing codebase patterns; if none exist, follow this guide.
3. **Minimize surface area** — export only what must be exported.
4. **Fail explicitly** — always handle errors, never discard them silently.
5. **Let tools do formatting** — never manually align code; rely on `gofmt`/`goimports`.

---

## Package Layout

- Package names: all lowercase, singular, no underscores (`strconv`, not `str_conv`).
- Never use generic names: `common`, `util`, `shared`, `helpers`, `lib`, `base`, `misc`.
- One package = one purpose. If a package needs the word "and" to describe it, split it.
- `internal/` for code that must not leak outside the module.

## Import Grouping

Always group imports in this order, separated by blank lines:

```go
import (
    // 1. Standard library
    "context"
    "fmt"

    // 2. External / third-party
    "github.com/go-chi/chi/v5"
    "go.uber.org/zap"

    // 3. Internal / same module
    "github.com/yourorg/yourapp/internal/domain"
)
```

Use import aliases only when the package name conflicts or differs from the last path element.

## Naming

- `MixedCaps` / `mixedCaps` — never underscores (except `TestFoo_SubCase` in tests).
- Receivers: 1–2 letter abbreviation of the type (`s` for `Service`, `r` for `Repo`). Consistent across all methods.
- Interfaces: name by behavior, not by the struct (`Reader`, `Storer`), suffix `-er` where natural.
- Avoid stuttering: `http.Client`, not `http.HTTPClient`; `user.Service`, not `user.UserService`.
- Unexported top-level vars/consts: prefix with `_` only if needed to avoid conflicts.
- Use `ctx` for `context.Context`, `err` for errors, `ok` for boolean map/type-assert results.

## Variable Declarations

```go
// Zero-value initialization — use var
var s string
var m map[string]int  // nil map is fine if only read

// Non-zero initialization — use short declaration
s := "hello"
m := make(map[string]int, 16) // allocate if you'll write to it

// Group related declarations
var (
    mu       sync.Mutex
    registry map[string]Handler
)
```

Never use `new(T)` when `&T{}` is clearer. Prefer `make` for slices/maps with known capacity.

## Functions

- Max ~60 lines; if longer, extract helpers.
- Accept interfaces, return structs.
- Parameters: `context.Context` always first, `error` always last in returns.
- Avoid naked returns — always name what you return explicitly in the `return` statement.
- Prefer returning `(T, error)` over `(*T, error)` when the zero value of T is useful.

## Error Handling (summary — see references/error-handling.md for full rules)

- Always check errors. If truly unreachable, add a comment explaining why.
- Wrap errors with `fmt.Errorf("operation: %w", err)` — include context about what failed.
- Error strings: lowercase, no punctuation at end. E.g. `"connect to db: %w"`.
- Use sentinel errors (`var ErrNotFound = errors.New(...)`) for errors callers must handle.
- Use custom error types only when callers need structured info (`type ValidationError struct{}`).
- Never panic in library code. Panic only for truly unrecoverable programmer errors in main.

## Struct Initialization

Always use field names — never positional:

```go
// Good
user := User{
    Name:  "Mark",
    Email: "mark@example.com",
}

// Bad — breaks if fields are reordered
user := User{"Mark", "mark@example.com"}
```

Prefer `&T{...}` over `new(T)` followed by field assignment.

## Testing

- Use table-driven tests with `t.Run` subtests.
- Test function names: `TestFunctionName_Scenario`.
- Use `testify/assert` or `testify/require` — `require` for fatal preconditions, `assert` for checks.
- Use `t.Helper()` in test helper functions.
- Use `t.Cleanup()` instead of `defer` when possible.
- Test packages: use `_test` package suffix for black-box tests, same package for white-box.
- Use `go.uber.org/goleak` in `TestMain` to detect goroutine leaks.

## Concurrency (summary — see references/performance.md for full rules)

- Never fire-and-forget goroutines. Always provide a shutdown mechanism.
- Use `errgroup.Group` for coordinated goroutine lifecycle.
- Use `sync.Once` for lazy init, not double-checked locking.
- Use `go.uber.org/atomic` for atomic operations — it provides type safety.
- Channels for coordination; mutexes for state protection.
- Keep critical sections short.

## Logging & Observability

- Use structured logging (`zap`, `slog`). Never `log.Println` or `fmt.Printf` for production logs.
- Log at the right level: `Error` for failures requiring attention, `Warn` for degraded but functional, `Info` for significant events, `Debug` for development.
- Include context fields: `zap.String("user_id", uid)`, not string interpolation.

## Dependencies

- Prefer `uber-go` ecosystem where applicable: `slog` (logging), `fx` (DI), `goleak`, `atomic`, `multierr`.
- Pin dependencies with exact versions in `go.mod`.
- Minimize third-party imports — stdlib is preferred when it covers the use case.
