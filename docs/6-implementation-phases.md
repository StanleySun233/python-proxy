# 6 Implementation Phases

## Phase 1

- freeze V1 product scope
- finalize core objects and route model
- finalize MySQL 8.0 schema and migration flow
- finalize control plane API contract

## Phase 2

- scaffold Go control-plane API
- implement account bootstrap and login
- implement nodes, chains, and route rules CRUD
- implement policy compilation and publish flow

## Phase 3

- scaffold Go node agent
- implement node enrollment
- implement node heartbeat and policy sync
- implement HTTP, HTTPS, and WS forwarding through configured chains

## Phase 4

- scaffold Next.js admin console
- implement auth screens
- implement nodes, chains, and route rules management pages
- implement node health and certificate status views

## Phase 5

- connect Chrome extension to auth and profile APIs
- add node selection UX
- add policy state display

## Phase 6

- certificate automation
- health checks
- failure recovery
- hardening
