# 10 Node Agent Data Plane

## Responsibilities

- maintain trusted connection with control plane
- fetch and cache latest policy revision
- execute HTTP, HTTPS, and WS forwarding
- execute node-to-node CONNECT relays
- expose heartbeat and local health
- manage local inbound certificate state

## Runtime Objects

- `local node identity`
- `trusted peers`
- `active policy revision`
- `listener set`
- `route matcher`
- `chain executor`

## Execution Rules

- unmatched traffic follows explicit default policy only
- chain loops must never execute
- final hop owns destination resolution for the configured scope
- last valid policy remains active if sync fails

## Local Persistence

- current assigned policy revision
- last known good compiled snapshot
- node cert material reference
- transient request metrics

## Go Layout Target

- `apps/proxy-node/cmd/proxy-node`
- `apps/proxy-node/internal/agentconfig`
- `apps/proxy-node/internal/policystore`
- `apps/proxy-node/internal/proxy`
- `apps/proxy-node/internal/controlrelay`
- `apps/proxy-node/internal/heartbeat`
- `apps/proxy-node/internal/controlplane`
