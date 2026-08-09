# Resilient Path, Topology, and Remote Defaults Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver ordered chain candidate groups with automatic connection-time replacement, operational topology views, and explicit tenant SSH/RDP access-path defaults.

**Architecture:** Replace flat chain hops with ordered candidate groups throughout storage, policy, and node runtime. Make the chain authoritative for access-path node membership, derive client entrypoints and topology, persist remote defaults per tenant and protocol, and render one shared topology projection in the dashboard and topology administration page.

**Tech Stack:** Go 1.23.12, MySQL 8, Bun, Next.js 16, React 19, TypeScript 6, React Flow 11, Docker BuildKit.

## Global Constraints

- Implement only the new contracts; do not add compatibility branches for flat `hops` or alphabetic remote fallback.
- Automatic replacement applies to connection establishment and reconnection, not transparent migration of established byte streams.
- Candidate order is deterministic priority order; do not add weights or percentages.
- Use English for durable code, UI text, and documents, with existing `zh` and `en` locale files for translated UI copy.
- Do not add source comments, explanatory blocks, or unfinished-work markers.
- Do not install dependencies.
- Do not stop, restart, kill, or replace an existing service.
- Docker cleanup is limited to `docker builder prune -af`; do not prune containers, images, networks, or volumes.

---

### Task 1: Persist ordered chain candidate groups

**Files:**
- Modify: `apps/panel/api/schema/final.sql`
- Modify: `apps/panel/api/internal/store/bun_models.go`
- Modify: `apps/panel/api/internal/features/proxy/domain/chain.go`
- Modify: `apps/panel/api/internal/store/proxy_repository.go`
- Modify: `apps/panel/api/internal/store/mysql_chain.go`
- Modify: `apps/panel/api/internal/store/schema_test.go`
- Create: `apps/panel/api/internal/store/mysql_chain_test.go`

**Interfaces:**
- Produces: `proxy.ChainHopGroup{Candidates []string}`, `Chain.HopGroups []ChainHopGroup`, and matching create, update, validate, and preview inputs.
- Produces: ordered persistence keyed by `(chain_id, hop_index, candidate_index)`.

- [ ] **Step 1: Write failing schema and repository tests**

```go
func TestFinalSchemaStoresMultipleCandidatesPerHop(t *testing.T) {
    raw, err := os.ReadFile("../../schema/final.sql")
    if err != nil { t.Fatal(err) }
    for _, required := range []string{"candidate_index INT NOT NULL", "PRIMARY KEY (chain_id, hop_index, candidate_index)"} {
        if !strings.Contains(string(raw), required) { t.Fatalf("missing %q", required) }
    }
}

func TestListChainHopGroupsPreservesPositionAndPriority(t *testing.T) {
    got := groupChainHopModels([]ChainHopModel{{HopIndex: 0, CandidateIndex: 0, NodeID: "a"}, {HopIndex: 0, CandidateIndex: 1, NodeID: "e"}, {HopIndex: 1, CandidateIndex: 0, NodeID: "b"}})
    want := []proxy.ChainHopGroup{{Candidates: []string{"a", "e"}}, {Candidates: []string{"b"}}}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got=%+v want=%+v", got, want) }
}
```

- [ ] **Step 2: Run the focused tests and verify failure**

Run: `cd apps/panel/api && go test ./internal/store -run 'TestFinalSchemaStoresMultipleCandidatesPerHop|TestListChainHopGroupsPreservesPositionAndPriority'`

Expected: FAIL because `candidate_index`, `ChainHopGroup`, and grouped repository loading do not exist.

- [ ] **Step 3: Replace the storage and domain contract**

```go
type ChainHopGroup struct {
    Candidates []string `json:"candidates"`
}

type Chain struct {
    ID               string          `json:"id"`
    Name             string          `json:"name"`
    DestinationScope string          `json:"destinationScope"`
    Enabled          bool            `json:"enabled"`
    HopGroups        []ChainHopGroup `json:"hopGroups"`
}
```

Change `ChainHopModel` to include `CandidateIndex int`, add `groupChainHopModels(models []ChainHopModel) []proxy.ChainHopGroup`, and replace flat hop readers and writers with grouped equivalents ordered by `hop_index, candidate_index`.

- [ ] **Step 4: Run store tests and formatting**

Run: `cd apps/panel/api && gofmt -w internal/features/proxy/domain/chain.go internal/store/bun_models.go internal/store/proxy_repository.go internal/store/mysql_chain.go internal/store/mysql_chain_test.go internal/store/schema_test.go && go test ./internal/store`

Expected: PASS.

- [ ] **Step 5: Commit the candidate-group persistence**

```bash
git add apps/panel/api/schema/final.sql apps/panel/api/internal/features/proxy/domain/chain.go apps/panel/api/internal/store
git commit -m "feat: persist chain candidate groups"
```

### Task 2: Validate, probe, compile, and publish candidate graphs

**Files:**
- Modify: `apps/panel/api/internal/features/proxy/service/chain.go`
- Modify: `apps/panel/api/internal/features/proxy/service/route.go`
- Create: `apps/panel/api/internal/features/proxy/service/chain_test.go`
- Modify: `apps/panel/api/internal/policy/compiler.go`
- Modify: `apps/panel/api/internal/policy/compiler_test.go`
- Modify: `apps/panel/api/internal/service/policy.go`
- Modify: `apps/panel/api/internal/store/mysql_policy.go`
- Modify: `apps/panel/api/internal/store/mysql_policy_test.go`
- Modify: `apps/panel/api/internal/store/mysql_node.go`
- Modify: `apps/panel/api/internal/store/mysql_node_test.go`
- Modify: `apps/panel/api/internal/store/mysql_delete_impact.go`

**Interfaces:**
- Consumes: `Chain.HopGroups` from Task 1.
- Produces: graph validation, deterministic probe result with `ResolvedHops`, `BlockingGroupIndex`, and candidate states.
- Produces: policy snapshots containing only `hopGroups`.

- [ ] **Step 1: Write failing invariant and compiler tests**

```go
func TestValidateChainAcceptsConnectedCandidateGroups(t *testing.T) {
    input := proxy.ValidateChainInput{DestinationScope: "scope-c", HopGroups: []proxy.ChainHopGroup{{Candidates: []string{"a", "e"}}, {Candidates: []string{"b", "f"}}, {Candidates: []string{"c"}}}}
    result, err := service.ValidateChain(tenantCtx, input)
    if err != nil || !result.Valid { t.Fatalf("result=%+v err=%v", result, err) }
}

func TestCompileRejectsIsolatedCandidate(t *testing.T) {
    _, err := Compile(nodes, linksWithoutF, []proxy.Chain{{ID: "chain", DestinationScope: "scope-c", Enabled: true, HopGroups: groups}}, nil, nil)
    if err == nil || !strings.Contains(err.Error(), "isolated_candidate") { t.Fatalf("err=%v", err) }
}
```

Add table cases for an empty group, duplicate node, non-edge entry candidate, wrong exit scope, missing inbound edge, missing outbound edge, disabled candidate, deterministic priority, and no viable candidate.

- [ ] **Step 2: Run focused tests and verify failure**

Run: `cd apps/panel/api && go test ./internal/features/proxy/service ./internal/policy -run 'Candidate|Group|Isolated|Duplicate|Entry|Exit'`

Expected: FAIL because service and compiler still traverse flat hops.

- [ ] **Step 3: Implement graph helpers and replace flat traversal**

```go
func flattenHopGroups(groups []proxy.ChainHopGroup) []string
func validateHopGroups(groups []proxy.ChainHopGroup, nodes []domain.Node, links []domain.NodeLink, destinationScope string) []string
func primaryHopIDs(groups []proxy.ChainHopGroup) []string
func resolveProbePath(groups []proxy.ChainHopGroup, nodes []domain.Node, links []domain.NodeLink, transports []domain.NodeTransport) proxy.ChainProbeResult
```

Use the helpers in create, update, validate, preview, route validation, policy compilation, per-node policy filtering, node delete impact, and policy publication.

- [ ] **Step 4: Run panel API tests**

Run: `cd apps/panel/api && gofmt -w internal/features/proxy internal/policy internal/service/policy.go internal/store && go test ./...`

Expected: PASS.

- [ ] **Step 5: Commit control-plane graph semantics**

```bash
git add apps/panel/api
git commit -m "feat: validate and publish resilient chains"
```

### Task 3: Select and replace node next-hop candidates

**Files:**
- Modify: `apps/node/api/internal/domain/types.go`
- Modify: `apps/node/api/internal/proxy/server.go`
- Modify: `apps/node/api/internal/proxy/chain.go`
- Modify: `apps/node/api/internal/proxy/forward_http.go`
- Modify: `apps/node/api/internal/proxy/connect_tunnel.go`
- Modify: `apps/node/api/internal/proxy/websocket_upgrade.go`
- Modify: `apps/node/api/internal/proxy/direct_stream.go`
- Modify: `apps/node/api/internal/proxy/server_test.go`
- Modify: `apps/node/api/internal/proxy/direct_stream_test.go`
- Modify: `apps/node/api/internal/policystore/store_test.go`
- Modify: `apps/node/api/internal/tcpaccess/server.go`
- Modify: `apps/node/api/internal/tcpaccess/server_test.go`

**Interfaces:**
- Consumes: policy `Chain.HopGroups` from Task 2.
- Produces: `resolveChainCandidates(snapshot, chainID) []chainHop` in priority order.
- Produces: connection-time fallback across viable candidates and a fail-closed `chain_candidates_unavailable` result.

- [ ] **Step 1: Write failing resolver and failover tests**

```go
func TestResolveChainCandidatesUsesStandbyWhenPrimaryTunnelIsDown(t *testing.T) {
    server := testServerAtNode("a", registryWithChild("f"))
    got := server.resolveChainCandidates(candidateSnapshot(), "chain")
    if len(got) != 1 || got[0].node.ID != "f" { t.Fatalf("got=%+v", got) }
}

func TestConnectTriesNextCandidateAfterOpenFailure(t *testing.T) {
    opener := &candidateOpener{errors: map[string]error{"b": errors.New("down")}}
    conn, selected, err := openCandidateStream(context.Background(), opener, []chainHop{hop("b"), hop("f")}, "target", 443)
    if err != nil || conn == nil || selected.node.ID != "f" { t.Fatalf("selected=%+v err=%v", selected, err) }
}
```

Cover HTTP reusable bodies, CONNECT, WebSocket upgrade, TCP access, deterministic priority, and all candidates failing before any success response is written.

- [ ] **Step 2: Run node tests and verify failure**

Run: `cd apps/node/api && go test ./internal/proxy ./internal/tcpaccess ./internal/policystore -run 'Candidate|Standby|NextCandidate'`

Expected: FAIL because the node policy and runtime expose one next hop.

- [ ] **Step 3: Implement candidate selection and connection fallback**

```go
type chainHop struct {
    node                domain.Node
    remainingHopGroups  []domain.ChainHopGroup
}

func (s *Server) resolveChainCandidates(snapshot policystore.Snapshot, chainID string) []chainHop
func (s *Server) candidateUsable(snapshot policystore.Snapshot, currentNodeID string, candidate domain.Node) bool
```

Select only candidates joined by a configured link and backed by a child tunnel, direct peer, or public endpoint. Refactor transport opening so failures can advance to the next candidate before a response or tunnel acknowledgement is committed.

- [ ] **Step 4: Run all node tests and formatting**

Run: `cd apps/node/api && gofmt -w internal/domain internal/proxy internal/tcpaccess && go test ./...`

Expected: PASS.

- [ ] **Step 5: Commit node failover**

```bash
git add apps/node/api
git commit -m "feat: replace unavailable chain hops"
```

### Task 4: Derive access-path topology and entrypoints from chains

**Files:**
- Modify: `apps/panel/api/schema/final.sql`
- Modify: `apps/panel/api/internal/domain/node.go`
- Modify: `apps/panel/api/internal/store/bun_models.go`
- Modify: `apps/panel/api/internal/store/proxy_repository.go`
- Modify: `apps/panel/api/internal/store/mysql_node_access_path.go`
- Modify: `apps/panel/api/internal/store/mysql_node.go`
- Modify: `apps/panel/api/internal/store/mysql_delete_impact.go`
- Modify: `apps/panel/api/internal/features/proxy/service/access_path.go`
- Modify: `apps/panel/api/internal/httpapi/proxy_token.go`
- Modify: `apps/panel/api/internal/service/policy.go`
- Modify: `apps/panel/api/internal/service/remote.go`
- Create: `apps/panel/api/internal/features/proxy/service/access_path_test.go`
- Modify: `apps/panel/api/internal/service/policy_test.go`

**Interfaces:**
- Consumes: authoritative chain candidate groups.
- Produces: `NodeAccessPath.RemoteProtocol`, `Entrypoints []AccessEntrypoint`, and `TopologyGroups []ChainHopGroup`.
- Removes: writable `TargetNodeID`, `EntryNodeID`, and `RelayNodeIDs` fields from access-path inputs and persistence.

- [ ] **Step 1: Write failing access-path derivation tests**

```go
func TestAccessPathDerivesEntrypointsAndTopologyFromChain(t *testing.T) {
    got := service.AccessPaths(tenantCtx)[0]
    groups := []proxy.ChainHopGroup{{Candidates: []string{"a", "e"}}, {Candidates: []string{"b"}}}
    if want := []string{"a", "e"}; !reflect.DeepEqual(entrypointNodeIDs(got.Entrypoints), want) { t.Fatalf("got=%+v want=%+v", got.Entrypoints, want) }
    if !reflect.DeepEqual(got.TopologyGroups, groups) { t.Fatalf("got=%+v want=%+v", got.TopologyGroups, groups) }
}

func TestRemoteProtocolRequiresTCPAccess(t *testing.T) {
    _, err := service.CreateAccessPath(tenantCtx, domain.CreateNodeAccessPathInput{Mode: "forward", RemoteProtocol: "ssh"})
    if errorCode(err) != "invalid_remote_protocol_path" { t.Fatalf("err=%v", err) }
}
```

- [ ] **Step 2: Run focused tests and verify failure**

Run: `cd apps/panel/api && go test ./internal/features/proxy/service ./internal/service -run 'AccessPathDerives|RemoteProtocol'`

Expected: FAIL because access paths own a separate node sequence and have no remote protocol.

- [ ] **Step 3: Implement authoritative chain derivation**

```go
type AccessEntrypoint struct {
    NodeID string `json:"nodeId"`
    Host   string `json:"host"`
    Port   int    `json:"port"`
    Status string `json:"status"`
}
```

Remove access-path node columns from `final.sql` and store models. Derive membership, entrypoints, primary target, bootstrap topology, proxy-token authorization, remote TCP frames, and node delete impact from the referenced chain.

- [ ] **Step 4: Run all panel API tests**

Run: `cd apps/panel/api && gofmt -w internal && go test ./...`

Expected: PASS.

- [ ] **Step 5: Commit authoritative access paths**

```bash
git add apps/panel/api
git commit -m "feat: derive access path topology from chains"
```

### Task 5: Persist and expose tenant remote defaults

**Files:**
- Modify: `apps/panel/api/schema/final.sql`
- Modify: `apps/panel/api/internal/domain/remote.go`
- Modify: `apps/panel/api/internal/store/store.go`
- Modify: `apps/panel/api/internal/store/bun_models.go`
- Create: `apps/panel/api/internal/store/mysql_remote_default.go`
- Create: `apps/panel/api/internal/store/mysql_remote_default_test.go`
- Modify: `apps/panel/api/internal/service/remote.go`
- Create: `apps/panel/api/internal/service/remote_test.go`
- Modify: `apps/panel/api/internal/httpapi/router.go`
- Modify: `apps/panel/api/internal/httpapi/handler_remote.go`

**Interfaces:**
- Produces: `RemoteAccessDefault{TenantID, Protocol, AccessPathID, UpdatedBy, UpdatedAt}`.
- Produces: `GET /api/remote/defaults`, `PUT /api/remote/defaults/{protocol}`.
- Changes: `RemoteSessionInput.AccessPathID` becomes optional and resolves an exact tenant default.

- [ ] **Step 1: Write failing store, service, and handler tests**

```go
func TestSetRemoteDefaultRejectsProtocolMismatch(t *testing.T) {
    _, err := control.SetRemoteAccessDefault(admin, tenantCtx, "ssh", "rdp-path")
    if errorCode(err) != "remote_default_protocol_mismatch" { t.Fatalf("err=%v", err) }
}

func TestCreateRemoteSessionUsesTenantDefault(t *testing.T) {
    session, err := control.CreateRemoteSession(account, tenantCtx, domain.RemoteSessionInput{Protocol: "ssh", Username: "ops", Password: "secret"})
    if err != nil || session.AccessPathID != "ssh-default" { t.Fatalf("session=%+v err=%v", session, err) }
}
```

Cover tenant-admin authorization, member read access, missing default, disabled default, explicit override, and protocol mismatch.

- [ ] **Step 2: Run focused tests and verify failure**

Run: `cd apps/panel/api && go test ./internal/store ./internal/service ./internal/httpapi -run 'RemoteDefault|TenantDefault'`

Expected: FAIL because the table, store methods, service methods, and routes do not exist.

- [ ] **Step 3: Implement strict default resolution**

```sql
CREATE TABLE tenant_remote_access_defaults (
  tenant_id VARCHAR(191) NOT NULL,
  remote_protocol VARCHAR(16) NOT NULL,
  access_path_id VARCHAR(191) NOT NULL,
  updated_by VARCHAR(191) NOT NULL,
  updated_at VARCHAR(64) NOT NULL,
  PRIMARY KEY (tenant_id, remote_protocol)
);
```

Do not select the first access path when no default exists. Return `remote_default_not_configured` or `remote_default_unavailable` with the exact condition.

- [ ] **Step 4: Run panel API tests**

Run: `cd apps/panel/api && gofmt -w internal && go test ./...`

Expected: PASS.

- [ ] **Step 5: Commit remote defaults**

```bash
git add apps/panel/api
git commit -m "feat: configure tenant remote defaults"
```

### Task 6: Replace the chain and access-path editors

**Files:**
- Modify: `apps/panel/web/lib/types/proxy.ts`
- Modify: `apps/panel/web/lib/types/nodes.ts`
- Modify: `apps/panel/web/lib/api/proxy.ts`
- Modify: `apps/panel/web/app/[locale]/(console)/proxy/studio/page.tsx`
- Modify: `apps/panel/web/app/[locale]/(console)/proxy/studio/_components/chain-editor.tsx`
- Modify: `apps/panel/web/app/[locale]/(console)/proxy/studio/_components/access-path-panel.tsx`
- Modify: `apps/panel/web/app/styles/proxy.css`
- Modify: `apps/panel/web/messages/en/proxyChains.json`
- Modify: `apps/panel/web/messages/zh/proxyChains.json`
- Modify: `apps/panel/web/messages/en/accessPaths.json`
- Modify: `apps/panel/web/messages/zh/accessPaths.json`

**Interfaces:**
- Consumes: `hopGroups`, derived access-path topology, and `remoteProtocol`.
- Produces: an ordered group editor with add candidate, add group, reorder, and remove actions.

- [ ] **Step 1: Change TypeScript contracts first and verify failure**

```ts
export type ChainHopGroup = {candidates: string[]};
export type Chain = ResourcePermissionMetadata & {
  id: string;
  name: string;
  destinationScope: string;
  enabled: boolean;
  hopGroups: ChainHopGroup[];
};
```

Run: `cd apps/panel/web && npm exec tsc -- --noEmit`

Expected: FAIL at every remaining flat `hops` consumer and obsolete access-path field.

- [ ] **Step 2: Implement the candidate-group and access-path forms**

Render each logical position as one bordered group with primary and standby candidates. Preserve configured order via the existing `@dnd-kit` dependency. Remove entry, relay, and target-node inputs from access paths and add `Remote protocol: None / SSH / RDP` for TCP access paths.

```tsx
<HopGroupCard
  candidates={group.candidates}
  index={index}
  nodes={nodes}
  onAddCandidate={nodeID => addCandidate(index, nodeID)}
  onRemoveCandidate={candidateIndex => removeCandidate(index, candidateIndex)}
/>
```

- [ ] **Step 3: Run frontend type checking**

Run: `cd apps/panel/web && npm exec tsc -- --noEmit`

Expected: PASS.

- [ ] **Step 4: Commit the editors**

```bash
git add apps/panel/web
git commit -m "feat: edit resilient access paths"
```

### Task 7: Build the operational topology projection and views

**Files:**
- Create: `apps/panel/api/internal/features/proxy/domain/topology.go`
- Create: `apps/panel/api/internal/features/proxy/service/topology.go`
- Create: `apps/panel/api/internal/features/proxy/service/topology_test.go`
- Modify: `apps/panel/api/internal/features/proxy/httpapi/router.go`
- Create: `apps/panel/api/internal/features/proxy/httpapi/handler_topology.go`
- Modify: `apps/panel/web/lib/api/proxy.ts`
- Modify: `apps/panel/web/lib/types/proxy.ts`
- Replace: `apps/panel/web/components/topology-preview.tsx`
- Modify: `apps/panel/web/app/[locale]/(console)/overview/dashboard/page.tsx`
- Modify: `apps/panel/web/app/[locale]/(console)/proxy/topology/_components/topology-page-content.tsx`
- Modify: `apps/panel/web/app/styles/features.css`
- Modify: `apps/panel/web/app/styles/nodes.css`
- Modify: `apps/panel/web/messages/en/overview.json`
- Modify: `apps/panel/web/messages/zh/overview.json`
- Modify: `apps/panel/web/messages/en/nodesConsole.json`
- Modify: `apps/panel/web/messages/zh/nodesConsole.json`

**Interfaces:**
- Produces: `GET /api/proxy/topology` with access paths, groups, candidates, configured dependency edges, health, selected primary edges, target scope, and blocking reason.
- Consumes: the same projection in dashboard and topology administration views.

- [ ] **Step 1: Write a failing topology projection test**

```go
func TestTopologyProjectionMarksPrimaryAndStandbyEdges(t *testing.T) {
    projection := service.Topology(tenantCtx)
    assertEdge(t, projection, "a", "b", "primary", true)
    assertEdge(t, projection, "a", "f", "standby", true)
    assertEdge(t, projection, "e", "b", "standby", false)
}
```

Run: `cd apps/panel/api && go test ./internal/features/proxy/service -run TestTopologyProjection`

Expected: FAIL because the projection does not exist.

- [ ] **Step 2: Implement and verify the backend projection**

Run: `cd apps/panel/api && gofmt -w internal/features/proxy && go test ./internal/features/proxy/...`

Expected: PASS.

- [ ] **Step 3: Replace the six-node graph with the projection view**

```tsx
<ReactFlow
  nodes={projectionNodes(selectedTopology)}
  edges={projectionEdges(selectedTopology)}
  fitView
  nodesDraggable={false}
  proOptions={{hideAttribution: true}}
/>
```

Use left-to-right group lanes, distinct client/path/scope nodes, solid primary edges, dashed standby edges, directional markers, health rings, a chain or path selector, and a blocking-state banner.

- [ ] **Step 4: Run frontend type checking**

Run: `cd apps/panel/web && npm exec tsc -- --noEmit`

Expected: PASS.

- [ ] **Step 5: Commit topology projection and UI**

```bash
git add apps/panel/api apps/panel/web
git commit -m "feat: visualize access path dependencies"
```

### Task 8: Expose remote defaults and protocol-specific paths in the panel

**Files:**
- Modify: `apps/panel/web/lib/types/remote.ts`
- Modify: `apps/panel/web/lib/api/remote.ts`
- Modify: `apps/panel/web/lib/api/index.ts`
- Modify: `apps/panel/web/app/[locale]/(console)/remote/_components/remote-helpers.ts`
- Modify: `apps/panel/web/app/[locale]/(console)/remote/_components/remote-page.tsx`
- Modify: `apps/panel/web/components/console-shell.tsx`
- Modify: `apps/panel/web/app/styles/remote.css`
- Modify: `apps/panel/web/messages/en/remote.json`
- Modify: `apps/panel/web/messages/zh/remote.json`
- Modify: `apps/panel/web/messages/en/shell.json`
- Modify: `apps/panel/web/messages/zh/shell.json`

**Interfaces:**
- Consumes: remote-default endpoints and protocol-specific access paths.
- Produces: visible default source, administrator update action, strict SSH/RDP filtering, unavailable-default state, and access-path management link.

- [ ] **Step 1: Replace remote types and verify compile failure**

```ts
export type RemoteAccessDefault = {
  tenantId: string;
  protocol: RemoteProtocol;
  accessPathId: string;
  updatedBy: string;
  updatedAt: string;
};

export type RemoteSessionPayload = {
  accessPathId?: string;
  protocol: RemoteProtocol;
  username: string;
  password?: string;
  privateKey?: string;
  passphrase?: string;
  width: number;
  height: number;
  dpi: number;
};
```

Run: `cd apps/panel/web && npm exec tsc -- --noEmit`

Expected: FAIL until APIs and page state use the new contract.

- [ ] **Step 2: Implement explicit default behavior**

Do not set `tcpPaths[0]`. Select a valid URL override first, otherwise the exact configured default. Show `Default not configured` with a management link when absent. Allow only tenant admins and super administrators to persist a new default.

- [ ] **Step 3: Run frontend type checking**

Run: `cd apps/panel/web && npm exec tsc -- --noEmit`

Expected: PASS.

- [ ] **Step 4: Commit remote configuration UX**

```bash
git add apps/panel/web
git commit -m "feat: expose remote access defaults"
```

### Task 9: Update client consumers and full local verification

**Files:**
- Modify: `apps/cli/src/control-plane.ts`
- Modify: `apps/cli/src/storage.ts`
- Modify: `apps/cli/src/init.ts`
- Modify: `apps/cli/src/daemon/router.ts`
- Modify: `apps/cli/test/cli.test.mjs`
- Modify: `apps/extension/vscode/src/extension.ts`
- Modify: `apps/extension/chrome/background/one-proxy-worker.js`
- Modify: `apps/extension/chrome/tools/background-source/state.js`
- Modify: `apps/extension/chrome/tools/background-source/status-bubble.js`
- Modify: `apps/extension/chrome/tools/background-source/proxy-auth.js`
- Modify: `apps/extension/chrome/tools/build_background_bundle.mjs`
- Modify: `apps/extension/cli/internal/proxycommand/direct.go`
- Modify: `apps/panel/api/openapi.yaml`
- Modify: `scripts/test-v210-docker-scenario.sh`

**Interfaces:**
- Consumes: access-path `entrypoints`, `topologyGroups`, and `remoteProtocol`.
- Produces: candidate-order client entry selection and a Docker scenario that proves both replacement examples and remote defaults.

- [ ] **Step 1: Add failing CLI and scenario assertions**

```js
assert.deepEqual(selectEntrypoint(pathWithFailedPrimary), {nodeId: 'e', host: 'e.test', port: 2990});
assert.equal(selectEntrypoint(pathWithNoHealthyCandidates), null);
```

Add scenario assertions for `[A|E]-B-C`, `A-[B|F]-C`, all candidates down, SSH default, and RDP path isolation.

- [ ] **Step 2: Run checks and verify failure**

Run: `cd apps/cli && npm run check && node --test test/*.mjs`

Run: `cd apps/extension/vscode && npm exec tsc -- --noEmit`

Run: `cd apps/extension/cli && go test ./...`

Expected: at least one contract assertion or type check fails before consumer updates.

- [ ] **Step 3: Update all consumers and regenerate the Chrome bundle**

Run: `cd apps/extension/chrome && node tools/build_background_bundle.mjs`

Expected: generated background bundle matches the updated source modules.

- [ ] **Step 4: Run the complete local verification matrix**

Run: `cd apps/panel/api && go test ./...`

Run: `cd apps/node/api && go test ./...`

Run: `cd apps/extension/cli && go test ./...`

Run: `cd apps/cli && npm run check && node --test test/*.mjs`

Run: `cd apps/panel/web && npm exec tsc -- --noEmit`

Run: `cd apps/extension/vscode && npm exec tsc -- --noEmit`

Expected: all commands PASS.

- [ ] **Step 5: Commit client and scenario contracts**

```bash
git add apps/cli apps/extension apps/panel/api/openapi.yaml scripts/test-v210-docker-scenario.sh
git commit -m "feat: consume resilient access paths"
```

### Task 10: Pure Docker build, scenario test, audit, and GitHub delivery

**Files:**
- Verify: `docker/one-proxy-panel.Dockerfile`
- Verify: `docker/one-proxy-node.Dockerfile`
- Verify: `.github/workflows/*`
- Verify: all committed files from Tasks 1-9

**Interfaces:**
- Produces: clean-cache image evidence, scenario evidence, clean repository state, and remote GitHub commit evidence.

- [ ] **Step 1: Audit the implementation against the design**

Run: `git diff origin/main...HEAD --check`

Run: `rg -n '"hops"|\.Hops|tcpPaths\[0\]' apps --glob '!**/node_modules/**'`

Expected: no obsolete flat chain contract or alphabetic remote fallback remains; unrelated uses are inspected and justified.

- [ ] **Step 2: Clear BuildKit cache only**

Run: `docker builder prune -af`

Expected: command succeeds and reports reclaimed or zero bytes without deleting containers, images, networks, or volumes.

- [ ] **Step 3: Build panel and node images without cache**

Run: `docker build --no-cache -f docker/one-proxy-panel.Dockerfile -t one-proxy-panel:resilient-path-test .`

Run: `docker build --no-cache -f docker/one-proxy-node.Dockerfile -t one-proxy-node:resilient-path-test .`

Expected: both images build successfully from the current commit.

- [ ] **Step 4: Run the isolated Docker scenario**

Run: `bash scripts/test-v210-docker-scenario.sh`

Expected: the script proves primary routing, entry replacement, relay replacement, fail-closed exhaustion, SSH default selection, RDP isolation, and reports PASS.

- [ ] **Step 5: Perform final repository and commit audit**

Run: `git status --short`

Run: `git log --oneline --decorate -12`

Expected: working tree is clean and all planned commits are present.

- [ ] **Step 6: Push and verify GitHub state**

Run: `git push origin main`

Run: `git ls-remote --heads origin main`

Expected: the remote `main` hash equals local `HEAD`.
