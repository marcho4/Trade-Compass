# Patterns — Structs, Interfaces, Enums, Options, Naming

## Table of Contents
1. Interface design
2. Struct embedding
3. Functional options
4. Enum patterns
5. Dependency injection
6. Constructor patterns
7. Method receivers
8. Nil handling

---

## 1. Interface Design

**Accept interfaces, return structs.** Define interfaces at the consumer, not the producer.

```go
// Good — interface defined where it's consumed
// In service/user.go:
type UserRepository interface {
    FindByID(ctx context.Context, id string) (User, error)
    Save(ctx context.Context, u User) error
}

type Service struct {
    repo UserRepository
}

// Bad — interface defined next to implementation
// In repo/user.go:
type UserRepository interface { ... }
type PostgresUserRepo struct { ... }
```

Keep interfaces small. One or two methods is ideal. Compose larger behaviors:

```go
type Reader interface { Read(ctx context.Context, id string) (T, error) }
type Writer interface { Write(ctx context.Context, v T) error }
type ReadWriter interface {
    Reader
    Writer
}
```

**Never export an interface just for mocking.** Use `_test.go` for test interfaces or use `go.uber.org/mock`.

## 2. Struct Embedding

Embed for behavior delegation, not for "inheritance":

```go
// Good — embedding adds the Lock/Unlock methods
type SafeMap struct {
    mu sync.Mutex
    m  map[string]string
}

// Bad — embedding leaks the full mutex API to callers
type SafeMap struct {
    sync.Mutex  // exported! callers can Lock/Unlock directly
    m map[string]string
}
```

Rules:
- Embed unexported types or interfaces, not exported structs.
- Never embed just to avoid writing a one-line delegation method.
- Embedded types go at the top of the struct, before regular fields.
- In outer types, use the embedded type's methods; don't shadow them unless intentional.

## 3. Functional Options

Use for constructors with many optional parameters:

```go
type Server struct {
    addr    string
    timeout time.Duration
    logger  *zap.Logger
    tls     *tls.Config
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *zap.Logger) Option {
    return func(s *Server) {
        s.logger = l
    }
}

func WithTLS(cfg *tls.Config) Option {
    return func(s *Server) {
        s.tls = cfg
    }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:    addr,
        timeout: 30 * time.Second,  // sensible default
        logger:  zap.NewNop(),       // sensible default
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer(":8080",
    WithTimeout(10*time.Second),
    WithLogger(logger),
)
```

When to use functional options vs config struct:
- Few options, most have good defaults → functional options.
- Many required fields → config struct with a constructor that validates.
- Need validation on option application → functional options that return error.

## 4. Enum Patterns

Go doesn't have enums. Use typed constants with `iota`:

```go
type Status int

const (
    StatusPending  Status = iota // 0
    StatusActive                  // 1
    StatusInactive                // 2
)

// Always implement Stringer for debuggability
func (s Status) String() string {
    switch s {
    case StatusPending:
        return "pending"
    case StatusActive:
        return "active"
    case StatusInactive:
        return "inactive"
    default:
        return fmt.Sprintf("Status(%d)", s)
    }
}
```

Rules:
- Start `iota` at 0 only if the zero value is meaningful. Otherwise start at 1:
  ```go
  type Role int
  const (
      _ Role = iota // skip 0 so zero value is "unknown"
      RoleAdmin
      RoleUser
  )
  ```
- Always add a `String()` method.
- Never depend on the numeric value of iota outside the package.
- For string-based enums (JSON APIs), use string constants instead:
  ```go
  type Currency string
  const (
      CurrencyRUB Currency = "RUB"
      CurrencyUSD Currency = "USD"
      CurrencyEUR Currency = "EUR"
  )
  ```

## 5. Dependency Injection

Prefer constructor injection. Pass dependencies explicitly:

```go
type Service struct {
    repo   UserRepository
    cache  Cache
    logger *zap.Logger
}

func NewService(repo UserRepository, cache Cache, logger *zap.Logger) *Service {
    return &Service{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}
```

For complex applications, use `go.uber.org/fx`:

```go
func main() {
    fx.New(
        fx.Provide(
            NewPostgresRepo,
            NewRedisCache,
            NewService,
            NewHTTPHandler,
        ),
        fx.Invoke(StartServer),
    ).Run()
}
```

Never use package-level global state for dependencies. No `init()` for connecting to databases.

## 6. Constructor Patterns

Naming: `New<Type>` returns `*Type`. If multiple constructors, `New<Type>With<Variant>` or use options.

```go
// Single constructor
func NewClient(baseURL string) *Client { ... }

// Constructor that can fail
func NewClient(baseURL string) (*Client, error) { ... }

// Constructor with options
func NewClient(baseURL string, opts ...Option) *Client { ... }
```

Validate in constructors, not in methods:

```go
func NewService(repo Repository, logger *zap.Logger) (*Service, error) {
    if repo == nil {
        return nil, errors.New("nil repository")
    }
    if logger == nil {
        logger = zap.NewNop()
    }
    return &Service{repo: repo, logger: logger}, nil
}
```

## 7. Method Receivers

- Use pointer receivers if the method mutates the receiver or if the struct is large.
- Use value receivers for small immutable types (coordinates, money).
- **Be consistent**: if any method has a pointer receiver, all methods on that type should too.

```go
// Pointer — mutates state
func (s *Service) Start() error { ... }

// Value — small, immutable
type Point struct{ X, Y float64 }
func (p Point) Distance(q Point) float64 { ... }
```

Receiver names: 1–2 letters, abbreviation of type. Never `this` or `self`.

## 8. Nil Handling

Prefer nil slices over empty slices — they JSON-marshal to `null` vs `[]`, so choose intentionally:

```go
// Returns nil if no results — marshals as null
func (r *Repo) FindAll(ctx context.Context) ([]User, error) {
    // ...
    return users, nil  // users is nil if no rows
}

// If API contract requires empty array, explicitly allocate:
func (r *Repo) FindAll(ctx context.Context) ([]User, error) {
    users := make([]User, 0)
    // ...
    return users, nil  // marshals as []
}
```

Nil maps are safe to read but panic on write. Always `make()` before writing:

```go
var m map[string]int
_ = m["key"]       // OK, returns zero value
m["key"] = 1       // PANIC

m = make(map[string]int)
m["key"] = 1       // OK
```

Guard against nil receivers in methods called on potentially nil pointers:

```go
func (s *Service) IsReady() bool {
    if s == nil {
        return false
    }
    return s.started.Load()
}
```
