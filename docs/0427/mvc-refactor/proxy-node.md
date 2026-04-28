# Proxy Node Refactor (one-proxy-node)

## Overview

one-proxy-node is already reasonably well-structured with files under 400 lines. The refactor focuses on ensuring consistent layering and package organization.

## Current State

| File | Lines | Notes |
|------|-------|-------|
| `tunnel/registry.go` | 388 | Largest file |
| `proxy/server.go` | 349 | Proxy server |
| `tunnel/controller.go` | 327 | Tunnel controller |
| `cmd/main.go` | 214 | Entry point |
| `runtime/manager.go` | 199 | Runtime manager |
| — 18 other files under 170 lines each | — | Already reasonable |

## Atomic Requirements

### Atomic 16: Review and extract tunnel/registry.go if needed
If `tunnel/registry.go` (388 lines) contains multiple concerns (connection pooling, registration, heartbeat), split into:
- `tunnel/registry.go` — core registry logic
- `tunnel/pool.go` — connection pooling

**Dependencies**: None

### Atomic 17: Ensure consistent package naming and import paths
Review all package names match directory names. Verify no circular imports. Clean up any unused imports.

**Dependencies**: None

## Acceptance Criteria

- All files under 350 lines
- Package naming is consistent
- `go build ./...` passes

## User Stories

- As a developer, the proxy node codebase is navigable without opening any single very large file
