# 18 Reverse Tunnel Requirements

## Goal

Support multi-hop proxy routing when some nodes have no public inbound access and some nodes cannot reach the control plane directly.

## Problem Statement

Current relay execution assumes every hop can be addressed through `public_host:public_port`.

That fails in these common topologies:

- `edge -> relay-private`
- `a(public) -> b(private but egress allowed) -> c(private with no egress but reachable from b)`
- target node can reach its parent node but cannot be reached directly from the panel or from earlier relay hops

The product needs a transport model where child nodes establish long-lived reverse tunnels to parent nodes and runtime routing reuses those tunnels instead of assuming direct inbound connectivity.

## Functional Scope

### F1. Parent-Child Reverse Tunnel

Any node without stable public inbound reachability must be able to initiate and maintain a persistent tunnel to its configured parent node.

The tunnel must support:

- authenticated registration
- periodic keepalive
- multiplexed runtime streams
- parent-driven probe requests

### F2. Multi-Hop Overlay Routing

Logical relay chains must be executable even when intermediate or final hops do not have public endpoints.

The execution model must allow:

- `a -> b` through a reverse tunnel from `b` to `a`
- `a -> b -> c` through nested child tunnels
- final target access from the last reachable node in the chain

### F3. Transport-Aware Reachability

Node reachability must become a first-class runtime object instead of being inferred only from `publicHost` and `publicPort`.

The system must track:

- public endpoint transports
- reverse child tunnel transports
- parent tunnel state
- last tunnel heartbeat
- observed latency

### F4. Chain Probe From Control Plane

Operators must be able to trigger a probe for a configured chain from the panel.

The probe must:

- start from the control plane
- resolve the first reachable transport for the first hop
- traverse the full chain using public endpoints or active tunnels
- report success, failure, and the first blocking hop

### F5. Tunnel-Assisted Child Bootstrap

Some nodes cannot reach the control plane directly.

The product must support child bootstrap through a parent node:

- parent node receives a child tunnel join
- parent forwards enrollment or attachment metadata upward
- control plane can still create and manage the child node record

### F6. Stable Operator Model

Operators should continue managing:

- nodes
- parent-child topology
- chains
- route rules

The new transport layer must not require operators to reason about low-level socket details on every route rule.

## Example Topologies

### Scenario A. Public Edge To Private Relay

- `node-hk` is public
- `node-astar-91` is private and can reach `node-hk`
- traffic enters `node-hk`
- `node-astar-91` maintains a reverse tunnel to `node-hk`
- `node-hk` forwards chain streams through that tunnel

### Scenario B. Public A, Private B, Dark C

- `a` has public ingress
- `b` can access `a` and `c`
- `c` cannot access internet or control plane
- `b` establishes reverse tunnel to `a`
- `c` establishes reverse tunnel or LAN tunnel to `b`
- runtime chain `a -> b -> c` is executed through nested tunnel forwarding

## Non-Goals For First Iteration

- no automatic optimal path search
- no full mesh gossip protocol
- no dynamic peer-to-peer NAT traversal
- no traffic recording or billing pipeline
- no multi-parent tunnel failover

## Acceptance Criteria

- control plane has explicit runtime objects for node transports
- control plane can report why a chain is not currently probeable
- node-agent can distinguish parent tunnel state from generic node heartbeat
- public endpoint is no longer the only supported transport for relay execution
- the architecture cleanly supports `a -> b -> c` where `c` has no direct control-plane connectivity
