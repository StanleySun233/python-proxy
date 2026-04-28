# Go Backend MVC Refactor (one-panel-api)

## Overview

Refactor the Go backend from monolithic files into a clean layered architecture:
- **Domain layer** (`internal/domain/`) — entity types, input/output DTOs, enum definitions (one file per entity)
- **Repository layer** (`internal/store/`) — data access, SQL queries (one file per entity)
- **Service layer** (`internal/service/`) — business logic (one file per domain)
- **Controller layer** (`internal/httpapi/`) — HTTP handlers (one file per entity)

## Current State

| File | Lines | Contains |
|------|-------|----------|
| `store/mysql.go` | 2473 | 65+ methods: init, accounts, nodes, chains, routes, health, groups, sessions, enums, policy, maintenance |
| `service/controlplane.go` | 1532 | 55+ methods: all business logic across all domains |
| `httpapi/resources.go` | 891 | 35+ HTTP handler functions for all routes |
| `domain/types.go` | 588 | All domain types, input/output structs, enums |
| `store/store.go` | 83 | Store interface with 60+ method signatures |
| `httpapi/router.go` | 213 | Route registration |

## Atomic Requirements

### Atomic 10: Split domain/types.go by entity
Create entity-specific type files under `internal/domain/`:
- `domain/account.go` — Account, CreateAccountInput, UpdateAccountInput, LoginResult, ExtensionBootstrap
- `domain/node.go` — Node, CreateNodeInput, UpdateNodeInput, ConnectNodeInput, ConnectedNodeResult, NodeScope, NodeAccessPath + inputs, NodeOnboardingTask + inputs, NodeLink + inputs, NodeTransport + inputs, BootstrapToken + inputs, EnrollNodeInput/Result, ApproveNodeEnrollmentResult, ExchangeNodeEnrollmentInput
- `domain/chain.go` — Chain, ChainWithDetails, CreateChainInput, UpdateChainInput, ValidateChainInput, ChainValidationResult, CompileChainHop, CompiledChainConfig, PreviewChainInput, ChainPreviewResult, ChainProbeResult, SaveChainProbeResultInput
- `domain/route.go` — RouteRule, RouteRuleWithDetails, CreateRouteRuleInput, UpdateRouteRuleInput, MatchType, RouteRuleValidationResult
- `domain/health.go` — NodeHealth, NodeHealthHistory, NodeHeartbeatInput, NodeCertRenewInput, NodeCertRenewResult
- `domain/group.go` — Group, GroupDetail, CreateGroupInput, UpdateGroupInput, SetGroupAccountsInput, SetGroupScopesInput
- `domain/policy.go` — PolicyRevision, NodeAgentPolicy
- `domain/enums.go` — FieldEnum, FieldEnumMap (already exists, keep or split if >200 lines)
- `domain/common.go` — Overview, Certificate, shared utility types

Keep `domain/types.go` as barrel re-export during transition, then remove.

**Dependencies**: None

### Atomic 11: Reorganize store interface into entity interfaces
Keep `store/store.go` as the main interface but compose it from entity-specific interfaces:

```go
type AccountStore interface { ... }
type NodeStore interface { ... }
type ChainStore interface { ... }
type RouteStore interface { ... }
type HealthStore interface { ... }
type GroupStore interface { ... }
type PolicyStore interface { ... }
type SessionStore interface { ... }
type MaintenanceStore interface { ... }

type Store interface {
    AccountStore
    NodeStore
    ChainStore
    // ...
}
```

Or alternatively: keep single Store interface but structure `mysql.go` methods into separate implementation files.

**Dependencies**: Atomic 10

### Atomic 12: Split store/mysql.go into entity repository files
Create entity-specific MySQL implementation files:
- `store/mysql_account.go` — Account CRUD, authentication, sessions
- `store/mysql_node.go` — Node CRUD, transports, links, access paths, onboarding tasks, bootstrap tokens, enrollment, approval
- `store/mysql_chain.go` — Chain CRUD, probe results
- `store/mysql_route.go` — Route rule CRUD
- `store/mysql_health.go` — Node health, heartbeat, certificate renewal
- `store/mysql_group.go` — Group CRUD, account/scope memberships
- `store/mysql_policy.go` — Policy revisions, publishing, node agent policy
- `store/mysql_enums.go` — Field enum queries
- `store/mysql_maintenance.go` — Cleanup jobs, status refresh
- `store/mysql.go` — MySQLStore struct, NewMySQLStore(), init(), helper utils

Each repository file should contain only methods related to that entity.

**Dependencies**: Atomic 11 (interface structure drives file organization)

### Atomic 13: Split service/controlplane.go by domain
Create domain-specific service files:
- `service/account.go` — Login, logout, accounts CRUD, session management
- `service/node.go` — Node CRUD, connect, approve, enroll, exchange, transports, links, access paths, onboarding tasks
- `service/chain.go` — Chain CRUD, probe, validate, preview
- `service/route.go` — Route rule CRUD, validation, match type suggestions
- `service/health.go` — Health queries, heartbeat processing
- `service/group.go` — Group CRUD, account/scope management
- `service/policy.go` — Policy revisions, publishing
- `service/enums.go` — Enum loading and querying
- `service/overview.go` — Overview aggregation, extension bootstrap
- `service/controlplane.go` — ControlPlane struct, NewControlPlane(), shared helpers

**Dependencies**: Atomic 10 (types), Atomic 12 (store methods referenced by name)

### Atomic 14: Split httpapi/resources.go into entity handler files
Create entity-specific handler files:
- `httpapi/handler_account.go` — handleAccounts, handleAccountByID
- `httpapi/handler_node.go` — handleNodes, handleNodeByID, handleNodeConnect, handleNodeApprove, handleNodeBootstrapToken, handleUnconsumedBootstrapTokens, handleNodeEnroll, handleNodeExchange, handlePendingNodes, handleNodeReject
- `httpapi/handler_nodelink.go` — handleNodeLinks, handleNodeTransports
- `httpapi/handler_accesspath.go` — handleNodeAccessPaths, handleNodeAccessPathByID
- `httpapi/handler_onboarding.go` — handleNodeOnboardingTasks, handleNodeOnboardingTaskByID
- `httpapi/handler_chain.go` — handleChains, handleChainByID, handleChainProbe, handleChainValidate, handleChainPreview
- `httpapi/handler_route.go` — handleRouteRules, handleRouteRuleByID, handleRouteRuleValidate, handleRouteRuleSuggestions
- `httpapi/handler_health.go` — handleNodeHealth, handleNodeHealthHistory
- `httpapi/handler_group.go` — handleGroups, handleGroupByID, handleGroupAccounts, handleGroupScopes
- `httpapi/handler_other.go` — handleOverview, handleExtensionBootstrap, handleCertificates, handleEnums, handlePolicyRevisions, handlePolicyPublish, handleNodeScopes

**Dependencies**: Atomic 13 (service methods must match)

### Atomic 15: Update router.go
After splitting, update `router.go` to reference new handler method locations. File structure stays the same; method references are implicit since all handlers are methods on `*Router`.

**Dependencies**: Atomic 14

## Acceptance Criteria

- No Go file exceeds 500 lines (except unavoidable cases like complex business logic)
- Each entity has its own domain types file
- Each entity has its own store implementation file
- Each entity has its own HTTP handler file
- Service methods grouped by domain
- `go build ./...` passes
- Existing tests pass (if any)

## User Stories

- As a developer, I can add a new entity CRUD by creating type + store + service + handler files
- As a developer, I can find all SQL queries for a specific entity in one file
- As a developer, I can understand HTTP routing without scrolling through 35 handler functions
