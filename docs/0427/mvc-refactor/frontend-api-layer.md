# Frontend API Layer Refactor

## Overview

Split `lib/control-plane-api.ts` (540 lines) and `lib/control-plane-types.ts` (313 lines) by domain, following a DAO/DTO pattern. Each domain gets its own API module and types module.

## Current State

- `lib/control-plane-api.ts` — 540 lines, 40+ exported functions, single `request()` helper
- `lib/control-plane-types.ts` — 313 lines, 30+ exported types, enum helper types

Both files mix concerns across auth, nodes, chains, routes, health, groups, onboarding, certificates, policies, setup, and enums.

## Atomic Requirements

### Atomic 1: Split types by domain
Create `lib/types/` directory. Extract types into domain files:
- `lib/types/auth.ts` — Account, LoginResult, Session
- `lib/types/nodes.ts` — Node, NodeLink, NodeTransport, NodeHealth, NodeHealthHistory, BootstrapToken, UnconsumedBootstrapToken, ConnectedNodeResult, NodeAccessPath, NodeOnboardingTask
- `lib/types/chains.ts` — Chain, ChainProbeHop, ChainProbeResult, ChainValidationResult, CompileChainHop, CompiledChainConfig, ChainPreviewResult
- `lib/types/routes.ts` — RouteRule, MatchValueValidation, ChainValidation, ScopeValidation, RouteRuleValidationResult
- `lib/types/health.ts` — NodeHealth, NodeHealthHistory (re-export from nodes or dedicated)
- `lib/types/groups.ts` — Group, GroupDetail
- `lib/types/common.ts` — APIResponse<T>, FieldEnumEntry, FieldEnumMap, InitRequest, InitResult, SetupStatus, TestConnectionRequest, TestConnectionResult, GenerateKeyResult, PolicyRevision, Certificate
- `lib/types/index.ts` — barrel re-exports all types

Each type file must contain only types related to that domain. No runtime code.

**Dependencies**: None (pure extraction, no logic changes)

### Atomic 2: Split API functions by domain
Create `lib/api/` directory. Extract API functions into domain modules:
- `lib/api/client.ts` — ControlPlaneAPIError, Session, request(), notifyUnauthorized(), constants (CONTROL_PLANE_PROXY_BASE, SESSION_STORAGE_KEY, AUTH_INVALID_EVENT)
- `lib/api/auth.ts` — login(), logout()
- `lib/api/nodes.ts` — getNodes(), createNode(), updateNode(), deleteNode(), connectNode(), approveNode(), rejectNode(), getPendingNodes(), getNodeLinks(), createNodeLink(), getNodeTransports(), createBootstrapToken(), getUnconsumedBootstrapTokens(), getNodeAccessPaths(), createNodeAccessPath(), updateNodeAccessPath(), deleteNodeAccessPath(), getNodeOnboardingTasks(), createNodeOnboardingTask(), updateNodeOnboardingTaskStatus()
- `lib/api/chains.ts` — getChains(), createChain(), probeChain(), validateChain(), previewChain()
- `lib/api/routes.ts` — getRouteRules(), createRouteRule(), validateRouteRule()
- `lib/api/health.ts` — getNodeHealth(), getNodeHealthHistory()
- `lib/api/groups.ts` — listGroups(), createGroup(), getGroup(), updateGroup(), deleteGroup(), setGroupAccounts(), setGroupScopes()
- `lib/api/accounts.ts` — getAccounts(), createAccount(), updateAccount(), deleteAccount()
- `lib/api/policies.ts` — getPolicyRevisions(), publishPolicy()
- `lib/api/certificates.ts` — getCertificates()
- `lib/api/setup.ts` — getSetupStatus(), testSetupConnection(), generateSetupKey(), submitSetupInit()
- `lib/api/enums.ts` — fetchEnums()
- `lib/api/index.ts` — barrel re-exports

Each API module imports `request` from `client.ts`. No circular dependencies.

**Dependencies**: Atomic 1 (types must exist first)

### Atomic 3: Update all imports across the codebase
After Atomic 1 and Atomic 2 are complete, update every file that imports from `@/lib/control-plane-api` or `@/lib/control-plane-types` to import from the new locations. Keep old files as re-export barrels during transition, then remove them.

**Dependencies**: Atomic 1, Atomic 2

## Acceptance Criteria

- No single API/types file exceeds 150 lines
- Each domain module is self-contained
- All existing imports updated
- TypeScript compilation passes
- No runtime regressions

## User Stories

- As a developer, I can find API functions for a specific domain (e.g., nodes) by looking at a single small file
- As a developer, I can find types for a specific entity without scrolling through 30+ unrelated type definitions
- As a developer, I can add a new API endpoint by editing only the relevant domain module
