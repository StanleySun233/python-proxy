# 2 System Architecture

## Top Level

- `apps/chrome-extension`: browser entry point
- `apps/one-panel-api`: Go control plane
- `apps/one-proxy-panel`: Next.js admin console
- `prototypes/proxy-node-demo`: Python reference only

## Control Plane And Data Plane

- control plane owns accounts, topology, route policy, enrollment, certificate lifecycle, and node health state
- data plane runs on nodes and executes the compiled proxy policy
- nodes must keep the last valid policy locally so traffic still flows during control plane outages

## Multi-Node Routing Model

Do not model routing as "node `a` can reach node `b`, so use it." Model it as explicit objects:

- `node`
- `node_link`
- `chain`
- `route_rule`
- `policy_revision`

## Private Network Ownership

The hard problem is overlapping internal networks.

Use `destination_scope` in route policy. A rule should resolve not only to a chain but also to the node or network context that owns the final target interpretation.

Example:

- request target `10.0.0.12`
- chain `a -> b -> c`
- `destination_scope = c-k8s`

That tells the final hop which network view is authoritative.

## Chain Execution

Recommended V1 model:

- `user -> edge node`: browser proxy
- `edge node -> relay node`: node-to-node CONNECT tunnel
- final node opens the target connection in its local network context

This keeps the first version simpler than building an overlay network.

## Certificates

Treat these as separate systems:

- public inbound certs for user-facing HTTPS and WSS
- private node-to-node trust certs for relay and control traffic

## High-Risk Areas

- overlapping RFC1918 ranges
- chain loops
- stale route policy
- compromised downstream node
- extension UX confusion about which path is active
