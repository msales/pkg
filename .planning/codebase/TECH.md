# Technology Stack

## Language & Runtime
Go 1.14

## Dependencies
| Dependency | Purpose |
|---|---|
| `github.com/go-redis/redis` | Redis client (v6) for cache and scan iterator implementations |
| `github.com/bradfitz/gomemcache` | Memcache client for cache implementation |
| `github.com/alicebob/miniredis` | In-memory Redis server for testing |
| `github.com/stretchr/testify` | Test assertions |

## Build System
- **Makefile:** Includes shared build targets via `github.com/msales/make/golang`
- **Dockerfile:** Minimal — uses `msales/go-builder:1.14-base-1.0.0`, copies source only (no explicit build step)

## Testing
- `testify/assert` for assertions
- Table-driven tests with `t.Run` subtests
- Internal test files (`*_internal_test.go`) for unexported function testing
- `alicebob/miniredis` for Redis integration tests without external dependencies

## CI/CD
GitHub Actions workflow (`.github/workflows/test.yaml`):
- Triggers on every push
- Runs on `ubuntu-latest` (not self-hosted)
- Uses `msales/github-actions` private action (`go-test`)
- No staticcheck configured in CI
