# Architecture

## Overview
A shared utility library providing independent, low-level building blocks used across the organization's Go services. Each sub-package is self-contained with no cross-dependencies between packages. The library was trimmed down from a larger set of packages (per commit history: "Removed all packages that have been moved to external repos").

## Package Structure

| Package | Responsibility |
|---|---|
| `bytes` | `Buffer` — pooled, allocation-efficient byte buffer with typed append methods (int, uint, float, bool, time, string). Uses `sync.Pool` for reuse. Zero-copy `String()` via `unsafe.Pointer`. |
| `cache` | `Cache` interface with `Get`/`GetMulti`/`Set`/`Add`/`Replace`/`Delete`/`Inc`/`Dec` operations. Implementations for Redis (`go-redis`) and Memcache (`gomemcache`). `Item` type with typed getters (`Bool`, `Int64`, `Uint64`, `Float64`, `String`, `Bytes`). Context-based cache injection. `nullCache` no-op implementation. |
| `cryptox` | PKCS7 padding (`PKCS7Pad`) and unpadding (`PKCS7Unpad`) for block cipher operations |
| `redisx` | `ScanIterator` interface abstracting Redis SCAN across single-node and cluster clients. `ClusterScanIterator` iterates over all master nodes. |
| `retry` | `Policy` interface with `ExponentialPolicy` implementation. `Run()` executes a function with retries. `Stop()` sentinel wraps an error to abort retries. |
| `syncx` | `Mutex` wrapping `sync.Mutex` with an additional `TryLock()` method using atomic CAS on internal state. |
| `utils` | `SplitMap()` — splits a string slice into a key-value map using a separator. |

## Key Interfaces

- **`cache.Cache`** — unified cache operations (Get, Set, Add, Replace, Delete, Inc, Dec)
- **`retry.Policy`** — retry policy returning next sleep duration and whether to continue
- **`redisx.ScanIterator`** — Val/Next/Err iterator for Redis SCAN results
- **`redisx.ClusterClient`** — abstraction for ForEachMaster on cluster clients

## Data Flow
Each package is standalone — there is no cross-package data flow. Consumers import individual packages as needed.

## External Dependencies
- **Redis** — used by `cache` (via `go-redis` UniversalClient) and `redisx` (scan iteration)
- **Memcache** — used by `cache` (via `gomemcache`)
