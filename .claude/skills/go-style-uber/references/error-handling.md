# Error Handling, Defer, and Panic — Reference

## Table of Contents
1. Error wrapping and context
2. Sentinel errors vs custom types
3. Error matching (errors.Is / errors.As)
4. Defer patterns
5. Panic rules

---

## 1. Error Wrapping and Context

Always add context when propagating errors. The caller should understand the full chain.

```go
// Good — wraps with context using %w
func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
    u, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return User{}, fmt.Errorf("get user %s: %w", id, err)
    }
    return u, nil
}

// Bad — raw propagation loses context
func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
    return s.repo.FindByID(ctx, id)
}
```

Error string conventions:
- Start lowercase: `"connect to database: %w"`, not `"Connect to database: %w"`.
- No ending punctuation.
- Describe the operation that failed, not the error itself.
- Use `: ` as separator to build a readable chain: `"get user mark: query db: connection refused"`.

When NOT to wrap with `%w`:
- When the underlying error is an implementation detail that callers should not match against.
- Use `fmt.Errorf("...: %v", err)` instead to break the chain intentionally.

## 2. Sentinel Errors vs Custom Types

**Sentinel errors** — use when callers need to check for a specific condition:

```go
// Define at package level
var (
    ErrNotFound   = errors.New("not found")
    ErrConflict   = errors.New("conflict")
    ErrForbidden  = errors.New("forbidden")
)

// Usage
func (r *Repo) FindByID(ctx context.Context, id string) (User, error) {
    row := r.db.QueryRowContext(ctx, query, id)
    if err := row.Scan(&u); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return User{}, ErrNotFound
        }
        return User{}, fmt.Errorf("scan user: %w", err)
    }
    return u, nil
}
```

**Custom error types** — use when callers need structured information:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s — %s", e.Field, e.Message)
}

// Caller side
var ve *ValidationError
if errors.As(err, &ve) {
    // access ve.Field, ve.Message
}
```

Decision matrix:
- Caller just needs to know "what kind" → sentinel error.
- Caller needs extra data from the error → custom type.
- Nobody checks the error → don't export it, just wrap.

## 3. Error Matching

Always use `errors.Is` and `errors.As`, never `==` or type assertion:

```go
// Good
if errors.Is(err, ErrNotFound) { ... }

// Bad — breaks if error is wrapped
if err == ErrNotFound { ... }

// Good
var ve *ValidationError
if errors.As(err, &ve) { ... }

// Bad
if ve, ok := err.(*ValidationError); ok { ... }
```

## 4. Defer Patterns

Order: defers execute LIFO. Place cleanup immediately after resource acquisition.

```go
func (s *Service) Process(ctx context.Context) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer func() {
        if err != nil {
            _ = tx.Rollback() // best-effort rollback
        }
    }()

    // ... work with tx ...

    if err = tx.Commit(); err != nil {
        return fmt.Errorf("commit: %w", err)
    }
    return nil
}
```

Handle errors in deferred closers:

```go
// Good — capture close error
defer func() {
    if cerr := f.Close(); cerr != nil && err == nil {
        err = fmt.Errorf("close file: %w", cerr)
    }
}()

// Acceptable — log if close error is non-critical
defer func() {
    if cerr := resp.Body.Close(); cerr != nil {
        s.log.Warn("close response body", zap.Error(cerr))
    }
}()
```

Never defer inside a loop — resources won't release until function exits. Extract to a helper.

## 5. Panic Rules

**Never panic in library/package code.** Period.

Acceptable panic locations:
- `main()` or init-time setup when a truly fatal misconfiguration is detected.
- Test helpers via `t.Fatal`.

When you must panic, use a descriptive message:

```go
func MustParseConfig(path string) Config {
    cfg, err := ParseConfig(path)
    if err != nil {
        panic(fmt.Sprintf("parse config %s: %v", path, err))
    }
    return cfg
}
```

`Must*` prefix signals to callers that the function panics on error. Use only at init time.

## Handling Errors in Goroutines

Goroutines must not silently swallow errors:

```go
// Good — errgroup propagates errors
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error {
    return s.processA(ctx)
})
g.Go(func() error {
    return s.processB(ctx)
})
if err := g.Wait(); err != nil {
    return fmt.Errorf("parallel processing: %w", err)
}

// Bad — error disappears
go func() {
    s.processA(ctx) // error ignored
}()
```

## Multierr

When collecting multiple independent errors:

```go
import "go.uber.org/multierr"

func (s *Service) Shutdown() error {
    var errs error
    errs = multierr.Append(errs, s.server.Close())
    errs = multierr.Append(errs, s.db.Close())
    errs = multierr.Append(errs, s.cache.Close())
    return errs
}
```
