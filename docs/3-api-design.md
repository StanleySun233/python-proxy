# 3 API Design

## Principles

- control plane API is the source of truth
- node agents consume compiled policy snapshots rather than raw UI data
- admin web and Chrome extension should not talk directly to SQLite

## Auth APIs

### `POST /api/v1/auth/login`

Purpose:
Authenticate `account/password` and return tokens.

### `POST /api/v1/auth/refresh`

Purpose:
Refresh access token.

### `POST /api/v1/auth/logout`

Purpose:
Invalidate current session.

## Account APIs

### `GET /api/v1/accounts`

Purpose:
List accounts.

### `POST /api/v1/accounts`

Purpose:
Create account with role.

### `PATCH /api/v1/accounts/{accountId}`

Purpose:
Update password, role, or status.

## Node APIs

### `GET /api/v1/nodes`

Purpose:
List nodes and health status.

### `POST /api/v1/nodes/bootstrap-token`

Purpose:
Generate node enrollment token.

### `POST /api/v1/nodes/enroll`

Purpose:
Complete first-time node enrollment.

### `PATCH /api/v1/nodes/{nodeId}`

Purpose:
Update node metadata, parent, scope, or enabled state.

## Chain APIs

### `GET /api/v1/chains`

Purpose:
List chains.

### `POST /api/v1/chains`

Purpose:
Create chain with ordered hops.

### `PATCH /api/v1/chains/{chainId}`

Purpose:
Update chain hops and validation status.

## Route Policy APIs

### `GET /api/v1/route-rules`

Purpose:
List route rules.

### `POST /api/v1/route-rules`

Purpose:
Create whitelist route rule.

### `POST /api/v1/policies/publish`

Purpose:
Compile rules and chains into a new policy revision.

### `GET /api/v1/policies/revisions`

Purpose:
List published revisions.

## Node Sync APIs

### `GET /api/v1/node-agent/policy`

Purpose:
Node fetches latest compiled policy.

### `POST /api/v1/node-agent/heartbeat`

Purpose:
Node reports health, version, and active revision.

### `POST /api/v1/node-agent/cert/renew`

Purpose:
Node requests internal certificate renewal.

## Health APIs

### `GET /api/v1/nodes/health`

Purpose:
List node heartbeat and current policy status.
