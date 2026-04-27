# Code Quality

## Test Coverage
- `bytes/` — `buffer_test.go` present
- `cache/` — comprehensive: `cache_test.go`, `cache_internal_test.go`, `item_internal_test.go`, `memcache_test.go`, `memcache_internal_test.go`, `redis_test.go`, `redis_internal_test.go`
- `cryptox/` — `padding_test.go` present
- `redisx/` — `client_test.go` present
- `retry/` — `retry_test.go` present
- `syncx/` — `mutex_test.go` present
- `utils/` — `strings_test.go` present
- Every package has test coverage — good discipline

## Code Patterns
- **Functional options** — `WithPoolSize()`, `WithTimeout()` etc. for cache configuration
- **Null object pattern** — `nullCache` and `nullDecoder` for safe default behavior
- **Context injection** — `cache.WithCache(ctx, c)` / `cache.FromContext(ctx)` for DI via context
- **Sentinel errors** — `ErrCacheMiss`, `ErrNotStored` as package-level error variables
- **Internal test files** — `*_internal_test.go` for testing unexported functions without exposing them
- **sync.Pool reuse** — `bytes.Pool` for allocation-efficient buffer management
- **Unsafe pointer tricks** — `bytes.Buffer.String()` uses zero-copy `unsafe.Pointer` cast; `syncx.Mutex.TryLock()` uses atomic CAS on internal mutex state

## Error Handling
- Sentinel errors for expected cache misses and conditional write failures
- Package-level unexported errors for internal validation (`errInvalidBlockSize`, etc.)
- `retry.Stop()` sentinel pattern to distinguish "abort retry" from "retry again"
- Consistent error message prefixing with package name

## Documentation
- Godoc comments on all exported types, functions, methods, and variables
- README.md and LICENCE present at root
- Error variables documented with their meaning
- Cache interface methods documented with semantic differences (Add vs Replace vs Set)
