# Resilient Path, Topology, and Remote Defaults Design

## Objective

Replace linear proxy chains with ordered candidate groups, expose the resulting traffic and dependency graph in the panel, and make SSH and RDP access-path defaults explicit and configurable.

## Current-State Findings

- A chain stores one node per `hop_index`, and the API exposes a flat `hops` array.
- A node resolves exactly one next hop. A disconnected transport cannot select a peer at the same logical position.
- An access path duplicates a single entry node, relay sequence, and target node even though it already references a chain.
- The dashboard topology receives only nodes, limits the graph to six nodes, and draws only `parentNodeId` edges.
- WebSSH and WebRDP select a `pathId` URL parameter or the first alphabetically sorted enabled TCP path.
- The access-path editor is labeled "Network", and TCP paths do not identify whether they are intended for SSH or RDP.

## Domain Model

### Chain candidate groups

A chain contains ordered hop groups. Each group contains ordered node candidates.

```json
{
  "id": "chain-1",
  "name": "production-egress",
  "destinationScope": "scope-1",
  "hopGroups": [
    {"candidates": ["node-a", "node-e"]},
    {"candidates": ["node-b", "node-f"]},
    {"candidates": ["node-c"]}
  ],
  "enabled": true
}
```

Candidate order is priority order. The model has no weights, load percentages, or generic graph-search settings.

The `chain_hops` table stores `candidate_index` in addition to `hop_index`. Its primary key is `(chain_id, hop_index, candidate_index)`. The old flat `hops` contract is replaced rather than retained.

### Validation invariants

- A chain contains at least one non-empty hop group.
- A node appears at most once in a chain.
- Every first-group candidate is an edge node.
- Every last-group candidate owns the chain destination scope.
- Every candidate outside the last group has at least one configured link to a candidate in the next group.
- Every candidate outside the first group has at least one configured link from a candidate in the previous group.
- Disabled, pending, unknown, duplicate, or structurally isolated candidates make the chain invalid.

### Access paths

An access path references the chain as the authoritative node topology. It no longer stores a separately editable entry node, relay list, or target node.

The client-facing access-path contract includes derived `entrypoints` and `topologyGroups`:

```json
{
  "id": "path-1",
  "chainId": "chain-1",
  "remoteProtocol": "ssh",
  "entrypoints": [
    {"nodeId": "node-a", "host": "a.example.com", "port": 2990, "status": "healthy"},
    {"nodeId": "node-e", "host": "e.example.com", "port": 2990, "status": "healthy"}
  ],
  "topologyGroups": [
    {"candidates": ["node-a", "node-e"]},
    {"candidates": ["node-b", "node-f"]},
    {"candidates": ["node-c"]}
  ]
}
```

`remoteProtocol` is optional for non-remote TCP access paths and is exactly `ssh` or `rdp` for paths shown in a remote console.

### Remote defaults

`tenant_remote_access_defaults` stores one default per `(tenant_id, remote_protocol)` and references an enabled tenant-visible access path with the same remote protocol.

Tenant administrators and super administrators can change the default. All tenant members can read it. A `pathId` URL parameter remains a one-connection override and does not mutate the default.

## Runtime Resolution

### Entry selection

OneProxy-aware clients receive all enabled entrypoints in candidate order. They select the first healthy endpoint and retry the remaining endpoints when connection establishment fails.

The panel remote-session service resolves the default path when no explicit path is supplied, selects an entrypoint using current node health, and records the selected concrete entry node.

### Next-hop selection

At each node, the policy snapshot identifies the current hop group. The node considers next-group candidates in configured order and filters them by the configured NodeLink and currently usable transport:

1. connected child tunnel;
2. available direct peer;
3. valid public endpoint.

The first usable candidate is selected. If the selected transport disappears before the stream is established, the node attempts the remaining candidates before returning an error. Selection is fail-closed when no candidate is usable.

### Session boundary

Automatic replacement applies during connection establishment and reconnection. An established TCP, SSH, RDP, CONNECT, or WebSocket byte stream is not migrated transparently after it breaks.

### Observability

Each connection records the resolved concrete node sequence and the replacement reason. Probe output returns candidate status per group, the selected path, the blocking group when no path exists, and the observation time.

## Panel Information Architecture

### Dashboard topology

The dashboard fetches nodes, configured links, transports, chains, access paths, and current health. It renders a selected access path or chain from left to right:

- client and access path;
- entry candidate group;
- relay candidate groups;
- exit candidate group;
- destination scope and target.

Solid directional edges identify the current preferred path. Dashed edges identify valid standby dependencies. Colors distinguish ingress, relay, and egress semantics. Nodes show role, health, and primary or standby state. A chain or access-path selector replaces the six-node truncation.

### Topology administration

The existing topology table remains the NodeLink editor. A graph above the table uses the same graph projection as the dashboard, so operational visualization and configured dependencies cannot disagree.

### Remote pages

The sidebar label "Network" becomes "Access Paths". Each remote page shows:

- current tenant default;
- selected path and whether it came from the default or URL override;
- path health and resolved entrypoint;
- an administrator action to set the default;
- a visible link to manage access paths.

SSH lists only `remoteProtocol=ssh` paths. RDP lists only `remoteProtocol=rdp` paths.

## API Boundaries

- Chain create, update, validate, preview, list, probe, policy, and bootstrap contracts use `hopGroups`.
- Access-path create and update derive topology from `chainId` and accept optional `remoteProtocol`.
- `GET /api/remote/defaults` returns the tenant SSH and RDP mappings.
- `PUT /api/remote/defaults/{protocol}` updates one tenant default.
- Remote session creation accepts an optional `accessPathId`; omission resolves the tenant default.
- Topology projection is returned by a panel API and reused by dashboard and topology administration views.

## Failure Behavior

- Invalid candidate graphs cannot be saved or published.
- A path with no healthy entrypoint remains visible with a blocking reason and cannot start a remote session.
- A missing remote default never falls back to alphabetical order. The UI requires explicit selection and explains how to configure the default.
- A configured default that becomes disabled or unauthorized is reported as unavailable and is not silently replaced.
- A node with no usable next candidate returns a concrete chain-unavailable error and records the blocking group.

## Verification

- Store tests prove candidate ordering and remote-default persistence.
- Service and compiler tests prove every graph invariant and fail-closed behavior.
- Node tests prove primary selection, standby selection after transport loss, deterministic priority, and complete failure.
- Remote tests prove protocol filtering, explicit override, default resolution, authorization, and unavailable-default behavior.
- Frontend type checking and production compilation prove the new contracts are consumed consistently.
- Docker scenario tests cover `A-B-C`, `[A|E]-B-C`, `A-[B|F]-C`, no-candidate failure, SSH default selection, and RDP isolation.
- A build-cache-only `docker builder prune -af` is followed by `--no-cache` image builds. Containers, images, networks, and volumes are not pruned.
- Final delivery requires a clean diff audit, committed changes, a successful push to the configured GitHub remote, and remote branch verification.
