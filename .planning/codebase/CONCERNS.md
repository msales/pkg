# Concerns & Risks

## Technical Debt
- Go version pinned at 1.14 — significantly outdated (no generics, no newer stdlib features, no module graph pruning)
- `go-redis` v6 is outdated — current major version is v9 (`github.com/redis/go-redis/v9`)
- `syncx.Mutex.TryLock()` relies on internal memory layout of `sync.Mutex` via `unsafe.Pointer` — fragile, could break on Go version upgrades. Note: Go 1.18+ added `sync.Mutex.TryLock()` natively, making this package unnecessary.
- `bytes.Buffer.String()` uses unsafe zero-copy conversion — could cause issues if buffer is returned to pool while string is still referenced
- `retry.Run()` uses recursion instead of iteration — could stack overflow on high retry counts (though unlikely with exponential backoff)
- Many packages were already migrated to external repos (per commit history) — remaining packages may be candidates for migration or deprecation

## Missing Coverage
- No benchmarks in any package — performance-sensitive code (`bytes.Buffer`, `cache`, `syncx.Mutex.TryLock`) would benefit from benchmarks
- `redisx.ClusterScanIterator` error aggregation across masters is not tested for partial failure scenarios
- `retry.Run()` recursion depth is not bounded beyond the policy's attempt count

## Security
- `cryptox` provides PKCS7 padding utilities — these are building blocks for encryption but do not perform encryption themselves. No cryptographic keys handled.
- `cache` package does not sanitize keys — callers must ensure keys do not contain unexpected characters
- No TLS/auth configuration exposed in Redis or Memcache client creation — relies on connection string and defaults

## Performance
- `bytes.Pool` uses `sync.Pool` for buffer reuse — efficient for high-throughput scenarios
- `bytes.Buffer.String()` zero-copy avoids allocation but requires careful pool return discipline
- `cache.redisCache.GetMulti` uses `MGet` for batch retrieval — efficient
- `syncx.Mutex.TryLock()` uses atomic CAS — non-blocking, suitable for contention-sensitive paths
- `retry.ExponentialPolicy` doubles sleep time on each attempt — standard exponential backoff without jitter (could cause thundering herd)
