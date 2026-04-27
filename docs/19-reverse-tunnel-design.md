# 19 Reverse Tunnel Design

## Architecture Summary

The control plane remains out of the data path.

Runtime traffic flows through an overlay formed by parent-child tunnels between proxy nodes.

### Control Plane Responsibilities

- store node topology
- store route rules and chains
- publish compiled policy
- observe transport state
- trigger chain probes

### Node Responsibilities

- maintain parent tunnel when required
- accept child tunnels when acting as upstream relay
- execute multiplexed streams
- forward probe and proxy sessions to the next hop
- reach final target on the last hop

## Core Concepts

### Node

Logical identity and operator-managed topology record.

Existing fields remain:

- `id`
- `name`
- `scopeKey`
- `parentNodeId`
- `enabled`

### Node Transport

Represents a concrete runtime transport binding for a node.

Recommended fields:

- `id`
- `nodeId`
- `transportType`
- `direction`
- `address`
- `status`
- `parentNodeId`
- `connectedAt`
- `lastHeartbeatAt`
- `latencyMs`
- `detailsJson`

`transportType` examples:

- `public_http`
- `public_https`
- `reverse_ws_parent`
- `child_ws_listener`

`direction` examples:

- `inbound`
- `outbound`

### Chain Probe

A chain probe is a control-plane initiated diagnostic that resolves each hop using the best available transport.

Result shape:

- `chainId`
- `status`
- `message`
- `resolvedHops`
- `blockingNodeId`
- `blockingReason`
- `targetHost`
- `targetPort`

### Tunnel Session

Persistent parent-child tunnel used for probe and future proxy streams.

Session responsibilities:

- register child identity
- authenticate parent-child relationship
- send keepalive
- open multiplexed stream
- forward stream to child or next hop

## Execution Model

### Case 1. Public To Private

Topology:

- `edge(public)`
- `private-child(egress only)`

Flow:

1. child establishes `wss` tunnel to parent
2. parent registers active child transport
3. route on parent resolves next hop through active child tunnel
4. parent opens stream on tunnel
5. child reaches local target

### Case 2. A To B To C

Topology:

- `a` public ingress
- `b` can reach `a` and `c`
- `c` can reach `b` only

Flow:

1. `b -> a` parent tunnel is active
2. `c -> b` parent tunnel is active
3. incoming stream lands on `a`
4. `a` forwards session to `b`
5. `b` forwards same session to `c`
6. `c` accesses target
7. response frames return on the reverse path

## Transport Resolution Rules

When runtime needs to reach a hop node:

1. prefer active reverse child tunnel if current node already has one for that child
2. else use public endpoint if available
3. else fail with an explicit blocking reason

Control-plane probe follows the same resolver so frontend behavior stays stable while transport implementations evolve.

## Database Additions

### `node_transports`

Stores current or recently observed transport bindings.

Suggested columns:

- `id`
- `node_id`
- `transport_type`
- `direction`
- `address`
- `status`
- `parent_node_id`
- `connected_at`
- `last_heartbeat_at`
- `latency_ms`
- `details_json`
- `created_at`
- `updated_at`

### `chain_probe_results`

Stores latest observable probe result for each chain.

Suggested columns:

- `chain_id`
- `status`
- `message`
- `resolved_hops_json`
- `blocking_node_id`
- `blocking_reason`
- `target_host`
- `target_port`
- `probed_at`

## API Additions

### Control Plane

- `GET /api/v1/node-transports`
- `POST /api/v1/chains/{chainId}/probe`
- `GET /api/v1/chains/{chainId}/probe`

### Node Tunnel

Parent-facing node tunnel entrypoint:

- `GET /api/v1/node-tunnel/connect`

Tunnel message families:

- `register`
- `heartbeat`
- `open_stream`
- `stream_data`
- `close_stream`
- `probe_request`
- `probe_result`

## Delivery Phases

### Phase 1

- add requirements and design docs
- add transport domain models and schema
- add transport list API
- add chain probe API that reports blockers using current public-endpoint-only execution

### Phase 2

- add parent-child reverse ws tunnel server in proxy node
- add child tunnel client in proxy node
- persist active transport state to control plane
- update chain probe to use active child tunnels

### Phase 3

- execute proxy streams over tunnel sessions
- support child-assisted bootstrap through parent
- expose richer path and chain health in panel

## Current Implementation Decision

Start with transport-aware diagnostics first.

That gives operators immediate visibility into why a chain such as `node-hk -> node-astar-91` cannot currently probe, while keeping the API contract compatible with the later reverse tunnel implementation.
