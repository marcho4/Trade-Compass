# Performance & Concurrency — Reference

## Table of Contents
1. Goroutine lifecycle
2. Channel patterns
3. Mutex vs channels
4. sync package essentials
5. Context propagation
6. Performance hot-path rules
7. Memory allocation tips

---

## 1. Goroutine Lifecycle

**Never fire-and-forget goroutines.** Every goroutine must have:
- A way to signal it to stop.
- A way to wait for it to finish.

```go
// Good — controllable goroutine
func (s *Service) Start(ctx context.Context) {
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                s.poll()
            }
        }
    }()
}

func (s *Service) Stop() {
    s.cancel()
    s.wg.Wait()
}
```

Use `errgroup.Group` for structured concurrency:

```go
g, ctx := errgroup.WithContext(ctx)

g.Go(func() error {
    return s.fetchPrices(ctx)
})
g.Go(func() error {
    return s.fetchNews(ctx)
})

if err := g.Wait(); err != nil {
    return fmt.Errorf("parallel fetch: %w", err)
}
```

Use `goleak.VerifyTestMain` to catch leaks in tests:

```go
func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

## 2. Channel Patterns

**Channels for coordination, mutexes for state protection.**

Signal channel (stop/done):

```go
done := make(chan struct{})
go func() {
    defer close(done)
    // ... work ...
}()
<-done // wait for completion
```

Fan-out / fan-in:

```go
func processAll(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(10) // limit concurrency

    for _, item := range items {
        g.Go(func() error {
            return process(ctx, item)
        })
    }
    return g.Wait()
}
```

Buffered channels as semaphores:

```go
sem := make(chan struct{}, maxConcurrent)
for _, item := range items {
    sem <- struct{}{} // acquire
    go func() {
        defer func() { <-sem }() // release
        process(item)
    }()
}
```

Rules:
- Prefer `errgroup` over manual channel coordination in most cases.
- Close channels from the sender side only. Never close from receiver.
- Nil channels block forever — useful for disabling select cases dynamically.

## 3. Mutex vs Channels

Use **mutex** when:
- Protecting a shared data structure (map, slice, counter).
- Critical section is short (just read/write a field).
- No inter-goroutine coordination needed.

Use **channels** when:
- Passing ownership of data between goroutines.
- Signaling events (done, stop, ready).
- Coordinating producer/consumer pipelines.

```go
// Mutex — protecting shared state
type Cache struct {
    mu    sync.RWMutex
    items map[string]Item
}

func (c *Cache) Get(key string) (Item, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, ok := c.items[key]
    return item, ok
}

func (c *Cache) Set(key string, item Item) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = item
}
```

Use `sync.RWMutex` when reads significantly outnumber writes.

## 4. sync Package Essentials

**sync.Once** — lazy initialization (thread-safe):

```go
type Client struct {
    initOnce sync.Once
    conn     *grpc.ClientConn
}

func (c *Client) getConn() *grpc.ClientConn {
    c.initOnce.Do(func() {
        c.conn = dial() // expensive, happens once
    })
    return c.conn
}
```

**sync.Pool** — reuse allocations on hot paths:

```go
var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

func process(data []byte) {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)
    }()
    buf.Write(data)
    // ...
}
```

**sync.Map** — only for two specific cases:
1. Entry is written once and read many times.
2. Multiple goroutines read/write disjoint key sets.
Otherwise, use a regular map with `sync.RWMutex`.

**go.uber.org/atomic** — type-safe atomic operations:

```go
type Server struct {
    running atomic.Bool
    conns   atomic.Int64
}

func (s *Server) Start() {
    if s.running.Swap(true) {
        return // already running
    }
    // ...
}

func (s *Server) Connections() int64 {
    return s.conns.Load()
}
```

Always prefer `uber-go/atomic` over `sync/atomic` — it prevents mixing up value types.

## 5. Context Propagation

- Every function that does I/O or may block takes `context.Context` as first parameter.
- Never store context in a struct. Pass it through the call chain.
- Use `context.WithTimeout` / `context.WithCancel` for deadline propagation.
- Check `ctx.Err()` before expensive operations.

```go
func (s *Service) FetchData(ctx context.Context, id string) (Data, error) {
    // Check context before expensive call
    if err := ctx.Err(); err != nil {
        return Data{}, fmt.Errorf("fetch data: %w", err)
    }

    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return s.client.Get(ctx, id)
}
```

Context values: use only for request-scoped data that crosses API boundaries (trace ID, auth). Never for optional parameters or dependencies.

```go
// Define unexported key type to prevent collisions
type ctxKey struct{}
var traceIDKey = ctxKey{}

func WithTraceID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, traceIDKey, id)
}

func TraceID(ctx context.Context) string {
    id, _ := ctx.Value(traceIDKey).(string)
    return id
}
```

## 6. Performance Hot-Path Rules

These rules apply **only on hot paths** (tight loops, high-QPS handlers). Don't prematurely optimize cold code.

**String conversion**: `strconv` over `fmt`:

```go
// Hot path — use strconv
s := strconv.Itoa(42)
n, _ := strconv.Atoi("42")

// Cold path — fmt is fine for readability
s := fmt.Sprintf("user_%d", id)
```

**Pre-allocate slices and maps** when size is known:

```go
users := make([]User, 0, len(ids))  // known capacity
index := make(map[string]int, len(items))
```

**Avoid repeated string→[]byte conversions.** Convert once, reuse:

```go
// Bad — converts on every call
for _, s := range strs {
    hash.Write([]byte(s))
}

// Good — reuse buffer
buf := make([]byte, 0, 256)
for _, s := range strs {
    buf = append(buf[:0], s...)
    hash.Write(buf)
}
```

**strings.Builder** for concatenation:

```go
var b strings.Builder
b.Grow(estimatedSize) // pre-grow if size is known
for _, s := range parts {
    b.WriteString(s)
}
result := b.String()
```

**Avoid defer in hot loops** — overhead per call. Inline the cleanup or extract to a function.

## 7. Memory Allocation Tips

- Use `go test -benchmem` to measure allocations.
- Pointer vs value: use values for small structs (<= ~64 bytes). Fewer heap allocations.
- Return slices of values `[]T`, not pointers `[]*T`, when items are small.
- Reuse buffers with `sync.Pool` on hot paths.
- Group declarations to reduce fragmentation:
  ```go
  // Good — single allocation for related fields
  type request struct {
      header [8]byte  // fixed-size, on stack
      body   []byte
  }
  ```
- For JSON serialization hot paths, consider `json.NewEncoder` writing to a reusable buffer instead of `json.Marshal` (which allocates every time).
