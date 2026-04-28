# Master Plan: MVC Layered Architecture Refactor

## Team
| Role | ID | Plan | Scope |
|------|----|------|-------|
| Frontend | FE-1 | [plan-fe-1.md](./plan-fe-1.md) | Types split, imports update |
| Frontend | FE-2 | [plan-fe-2.md](./plan-fe-2.md) | API split, onboarding refactor |
| Frontend | FE-3 | [plan-fe-3.md](./plan-fe-3.md) | Routes, health overview, health heartbeat refactor |
| Backend | BE-1 | [plan-be-1.md](./plan-be-1.md) | Domain types, store interface, store split |
| Backend | BE-2 | [plan-be-2.md](./plan-be-2.md) | Service layer split |
| Backend | BE-3 | [plan-be-3.md](./plan-be-3.md) | HTTP handlers split, router, proxy-node review |

## FE-1 — [plan](./plan-fe-1.md)
- [x] Split types by domain (056f093) — [ref](./frontend-api-layer.md#atomic-1)
- [x] Update all imports (980498c) — [ref](./frontend-api-layer.md#atomic-3)

## FE-2 — [plan](./plan-fe-2.md)
- [x] Split API functions by domain (c9b1c1c) — [ref](./frontend-api-layer.md#atomic-2)
- [x] Extract onboarding hooks (2961b5a) — [ref](./frontend-pages.md#atomic-4)
- [x] Extract onboarding sub-components (00ae7bb) — [ref](./frontend-pages.md#atomic-5)

## FE-3 — [plan](./plan-fe-3.md)
- [x] Split node-pages.tsx into individual files (f37eb69) — [ref](./frontend-pages.md#atomic-6)
- [x] Extract routes hooks and validation (c06fdfb) — [ref](./frontend-pages.md#atomic-7)
- [x] Extract health overview hooks and sub-components (4fc8af3) — [ref](./frontend-pages.md#atomic-8)
- [x] Extract health heartbeat hooks and sub-components (63f489b) — [ref](./frontend-pages.md#atomic-9)

## BE-1 — [plan](./plan-be-1.md)
- [x] Split domain/types.go by entity (8b13a72) — [ref](./backend-go.md#atomic-10)
- [x] Reorganize store interface into entity interfaces (4d7ae29) — [ref](./backend-go.md#atomic-11)
- [x] Split store/mysql.go into entity repository files (8f68595) — [ref](./backend-go.md#atomic-12)

## BE-2 — [plan](./plan-be-2.md)
- [x] Split service/controlplane.go by domain (e1a97ed) — [ref](./backend-go.md#atomic-13)

## BE-3 — [plan](./plan-be-3.md)
- [x] Split httpapi/resources.go into entity handler files (9ccf512) — [ref](./backend-go.md#atomic-14)
- [x] Update router.go and verify go build — [ref](./backend-go.md#atomic-15)
- [x] Review one-proxy-node tunnel/registry.go (e6a90b7) — [ref](./proxy-node.md#atomic-16)
- [x] Verify one-proxy-node package consistency — [ref](./proxy-node.md#atomic-17)

## Dependency Graph

```
Phase 1 (parallel):
  FE-1: Atomic 1 (types split)
  FE-3: Atomic 6 (node-pages split — independent)

Phase 2 (after Atomic 1):
  FE-2: Atomic 2 (API split)
  BE-1: Atomic 10 (domain types split)

Phase 3 (after Atomic 2):
  FE-1: Atomic 3 (update imports)
  BE-2: Atomic 13 (service split — depends on BE-1 types)

Phase 4 (after Atomic 3):
  FE-2: Atomic 4,5 (onboarding refactor)
  FE-3: Atomic 7,8,9 (routes + health refactor)

Phase 5 (after BE-1 Atomic 10,11):
  BE-1: Atomic 12 (store split)

Phase 6 (after BE-2 Atomic 13):
  BE-3: Atomic 14,15 (handlers + router)

Phase 7 (after BE-3):
  BE-3: Atomic 16,17 (proxy-node review)
```

## Integration
- FE-1 + FE-2: Import paths must be consistent between types and API modules
- BE-1 + BE-2: Domain types must match service method signatures
- BE-2 + BE-3: Service methods must match handler calls
- All: `tsc --noEmit` and `go build ./...` must pass at each phase
