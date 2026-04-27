@~/CLAUDE.md

## Project Description

**Language:** Go 1.14  
**Module:** `github.com/msales/pkg/v5`  
**Purpose:** Shared utility library providing common low-level building blocks — pooled byte buffers, cache abstraction (Redis/Memcache), PKCS7 padding, Redis scan iterators (cluster-aware), retry with exponential backoff, try-lock mutex, and string utilities.  
**Key dependencies:** `github.com/go-redis/redis` (Redis client), `github.com/bradfitz/gomemcache` (Memcache client), `github.com/alicebob/miniredis` (in-memory Redis for testing), `github.com/stretchr/testify` (testing).

## Directory Structure

```
bytes/          Pooled byte buffer with typed append methods (int, uint, float, bool, time, string)
cache/          Cache interface with Redis and Memcache implementations, Item type with typed getters, context helpers
cryptox/        PKCS7 padding/unpadding utilities
redisx/         Redis scan iterator abstraction supporting both single-node and cluster clients
retry/          Retry execution with pluggable policy (exponential backoff) and stop sentinel
syncx/          Mutex with TryLock capability (atomic compare-and-swap)
utils/          String utilities (SplitMap)
```

## Build, Test & Lint

```bash
# Build (library only — no binaries)
go build ./...

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Vet
go vet ./...

# Static analysis
staticcheck ./...
```

CI runs on GitHub Actions (`.github/workflows/test.yaml`) using `ubuntu-latest` runners with `go-test` action.

Makefile includes shared build targets via `github.com/msales/make/golang`.

The Dockerfile is minimal — for CI testing only.

## Code Style

### Imports

Two groups separated by blank lines:
1. Standard library
2. Third-party packages

No internal packages — each sub-package is independent.

### Naming Conventions

- **Packages:** short, lowercase, suffixed with `x` for extended stdlib packages (e.g., `syncx`, `redisx`, `cryptox`)
- **Interfaces:** noun-based (e.g., `Cache`, `Policy`, `ScanIterator`)
- **Functional options:** `WithXxx` functions returning option func types (e.g., `WithPoolSize`, `WithTimeout`)
- **Null objects:** prefixed with `null` (e.g., `nullCache`, `nullDecoder`)
- **Sentinel errors:** package-level `var` with `Err` prefix (e.g., `ErrCacheMiss`, `ErrNotStored`)

### Error Handling

- Sentinel errors for expected conditions (`ErrCacheMiss`, `ErrNotStored`)
- Package-level error variables for internal errors (e.g., `errInvalidBlockSize`)
- Errors prefixed with package name (e.g., `"cryptox: ..."`, `"cache: ..."`)
- Retry pattern with `Stop()` sentinel to abort retries early

### Testing

- **Framework:** `testify/assert` for assertions
- **Pattern:** Table-driven tests with `t.Run` subtests
- **Internal tests:** `_internal_test.go` files for testing unexported functions
- **Redis tests:** Uses `alicebob/miniredis` for in-memory Redis testing

## Branching & Git

- **Never** work directly on `main` or `master`. Always create a new branch off the latest `main`/`master`.
- Before creating a branch, **always pull** the latest changes from the remote (`git pull origin main`).
- If the user provides a task identifier in the format `TRK-XXXX`, use it as the branch name (e.g., `TRK-1234`). Otherwise, **ask the user** for a task ID or branch name before proceeding.
- **Never push** to `main` or `master` directly.
- Commit changes after each meaningful iteration. Use your judgment to decide when progress should be saved — prefer smaller, atomic commits over large monolithic ones.
- Write clear, concise commit messages. Use conventional commit format when appropriate (e.g., `feat:`, `fix:`, `refactor:`, `test:`).
