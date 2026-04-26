# 14 Node Onboarding Implementation

## Implementation Strategy

Build onboarding in layers and avoid mixing product model with transport execution.

Phase 1 should land the control-plane model first:

- onboarding access path model
- onboarding task model
- CRUD and task APIs
- node-agent runtime tolerance for missing control-plane address

Phase 2 can add execution workers:

- direct push worker
- relay push worker
- upstream pull forwarding worker

Current implementation status:

- upstream pull forwarding is available through a node-agent bound to control plane
- direct push task performs a minimal `/healthz` reachability probe
- relay-chain task executes hop-by-hop probe through ordered relay nodes

## Core Objects

### Node Access Path

Purpose:
Describe how the system can reach a target node for onboarding.

Fields:

- `id`
- `name`
- `mode`
- `targetNodeId`
- `entryNodeId`
- `relayNodeIds`
- `targetHost`
- `targetPort`
- `enabled`

Rules:

- `mode = direct` means control plane connects to target directly
- `mode = relay_chain` means control plane reaches target through ordered relay hops
- `mode = upstream_pull` means target node uses the configured entry node as its upstream onboarding entry

### Node Onboarding Task

Purpose:
Represent a concrete operator-triggered onboarding attempt.

Fields:

- `id`
- `mode`
- `pathId`
- `targetNodeId`
- `targetHost`
- `targetPort`
- `status`
- `statusMessage`
- `requestedByAccountId`
- `createdAt`
- `updatedAt`

Rules:

- task creation is synchronous
- execution can remain manual or stubbed in the first implementation
- status values should be stable because frontend state will depend on them

Recommended statuses:

- `planned`
- `dispatching`
- `pending`
- `connected`
- `failed`
- `cancelled`

## Control Plane Changes

### Domain Layer

Add:

- `NodeAccessPath`
- `CreateNodeAccessPathInput`
- `UpdateNodeAccessPathInput`
- `NodeOnboardingTask`
- `CreateNodeOnboardingTaskInput`

### Store Layer

Add persistence boundaries for:

- list/create/update/delete node access path
- list/create onboarding task

### Service Layer

Add validation rules:

- path name required
- mode required
- target host required for direct and relay modes
- target port must be positive when target host is set
- onboarding task must specify mode and either path or direct target

### HTTP API Layer

Add:

- `/api/v1/node-access-paths`
- `/api/v1/node-access-paths/{pathId}`
- `/api/v1/node-onboarding-tasks`

## SQLite Changes

### Tables

Add `node_access_paths`:

- basic path metadata
- JSON field for relay hop node ids

Add `node_onboarding_tasks`:

- task metadata
- status and message
- operator account id

### Migration Policy

For now, add these with startup `CREATE TABLE IF NOT EXISTS` guards so current databases keep working.

## Node-Agent Changes

The node-agent should continue supporting:

- full control-plane binding with enrollment and heartbeat
- local-only startup without control-plane URL
- forwarding selected control-plane onboarding APIs when bound to control plane
- exposing relay probe execution for control-plane-driven multi-hop reachability check

That allows onboarding mode selection later without blocking agent runtime.

## Deferred Work

- task executor and retry loop
- per-hop error detail
- panel wizard UX
