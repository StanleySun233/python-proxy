# 3 API Design

## Principles

- control plane API is the source of truth
- admin web and Chrome extension consume control plane APIs only
- node agents consume compiled policy snapshots only
- protected control plane APIs use `Authorization: Bearer <access-token>`
- node agent sync APIs use `Authorization: Bearer <node-access-token>`
- success responses use `APIResponse[T]`
- error responses use the same envelope with non-zero `code`
- edge public certificate provider defaults to `lets_encrypt`

## Response Shape

### Success

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### Error

```json
{
  "code": 400,
  "message": "invalid_json"
}
```

## Auth APIs

### `POST /api/v1/auth/login`

Purpose:
Authenticate `account/password` and return `accessToken` plus `refreshToken`.

Auth:
No Bearer token required.

### `POST /api/v1/auth/refresh`

Purpose:
Refresh session tokens.

Auth:
Accepts `Authorization: Bearer <refresh-token>` or request body `refreshToken`.

### `POST /api/v1/auth/logout`

Purpose:
Invalidate current access session.

Auth:
Requires `Authorization: Bearer <access-token>`.

## Account APIs

### `GET /api/v1/accounts`

Purpose:
List accounts.

Auth:
Control plane Bearer token required.

### `POST /api/v1/accounts`

Purpose:
Create account with role.

Auth:
Control plane Bearer token required.

### `PATCH /api/v1/accounts/{accountId}`

Purpose:
Update password, role, or status.

Auth:
Control plane Bearer token required.

## Node APIs

### `GET /api/v1/nodes`

Purpose:
List nodes and current status.

Auth:
Control plane Bearer token required.

### `POST /api/v1/nodes`

Purpose:
Create node metadata manually.

Auth:
Control plane Bearer token required.

### `PATCH /api/v1/nodes/{nodeId}`

Purpose:
Update node metadata, parent, scope, enabled state, or status.

Auth:
Control plane Bearer token required.

### `DELETE /api/v1/nodes/{nodeId}`

Purpose:
Delete node.

Auth:
Control plane Bearer token required.

### `POST /api/v1/nodes/bootstrap-token`

Purpose:
Generate bootstrap token for node enrollment.

Auth:
Control plane Bearer token required.

### `POST /api/v1/nodes/enroll`

Purpose:
Create pending node enrollment request and return one-time enrollment secret.

Auth:
No Bearer token required.

### `POST /api/v1/nodes/approve/{nodeId}`

Purpose:
Approve pending node enrollment, activate trust material, and issue node access token.

Auth:
Control plane Bearer token required.

### `POST /api/v1/nodes/exchange`

Purpose:
Node exchanges approved enrollment secret for long-lived node access token and trust material.

Auth:
No Bearer token required.

## Topology APIs

### `GET /api/v1/node-links`

Purpose:
List node-to-node relationship metadata.

Auth:
Control plane Bearer token required.

### `POST /api/v1/node-links`

Purpose:
Create node-to-node relationship metadata.

Auth:
Control plane Bearer token required.

### `GET /api/v1/node-access-paths`

Purpose:
List saved onboarding access paths for direct, relay-chain, and upstream-pull node onboarding.

Auth:
Control plane Bearer token required.

### `POST /api/v1/node-access-paths`

Purpose:
Create a saved onboarding access path.

Auth:
Control plane Bearer token required.

### `PATCH /api/v1/node-access-paths/{pathId}`

Purpose:
Update onboarding access path metadata.

Auth:
Control plane Bearer token required.

### `DELETE /api/v1/node-access-paths/{pathId}`

Purpose:
Delete onboarding access path.

Auth:
Control plane Bearer token required.

### `GET /api/v1/node-onboarding-tasks`

Purpose:
List panel-created node onboarding tasks.

Auth:
Control plane Bearer token required.

### `POST /api/v1/node-onboarding-tasks`

Purpose:
Create a node onboarding task using direct or relay access intent.

Auth:
Control plane Bearer token required.

## Chain APIs

### `GET /api/v1/chains`

Purpose:
List chains.

Auth:
Control plane Bearer token required.

### `POST /api/v1/chains`

Purpose:
Create chain with ordered hops.

Auth:
Control plane Bearer token required.

### `PATCH /api/v1/chains/{chainId}`

Purpose:
Update chain hops and enabled state.

Auth:
Control plane Bearer token required.

### `DELETE /api/v1/chains/{chainId}`

Purpose:
Delete chain.

Auth:
Control plane Bearer token required.

## Route Policy APIs

### `GET /api/v1/route-rules`

Purpose:
List route rules.

Auth:
Control plane Bearer token required.

### `POST /api/v1/route-rules`

Purpose:
Create whitelist route rule.

Auth:
Control plane Bearer token required.

### `PATCH /api/v1/route-rules/{ruleId}`

Purpose:
Update route rule.

Auth:
Control plane Bearer token required.

### `DELETE /api/v1/route-rules/{ruleId}`

Purpose:
Delete route rule.

Auth:
Control plane Bearer token required.

### `GET /api/v1/policies/revisions`

Purpose:
List policy revisions.

Auth:
Control plane Bearer token required.

### `POST /api/v1/policies/publish`

Purpose:
Compile current graph and publish a new revision.

Auth:
Control plane Bearer token required.

## Runtime Status APIs

### `GET /api/v1/overview`

Purpose:
Summarize node, policy, and certificate state.

Auth:
Control plane Bearer token required.

### `GET /api/v1/nodes/health`

Purpose:
List node heartbeat and policy status.

Auth:
Control plane Bearer token required.

### `GET /api/v1/certificates`

Purpose:
List certificate status and expiry metadata.

Auth:
Control plane Bearer token required.

## Node Sync APIs

### `GET /api/v1/node-agent/policy`

Purpose:
Node fetches current assigned node-specific compiled policy.

Auth:
Node Bearer token required.

### `POST /api/v1/node-agent/heartbeat`

Purpose:
Node reports active revision and listener/certificate status.

Auth:
Node Bearer token required.

### `POST /api/v1/node-agent/cert/renew`

Purpose:
Node requests internal certificate renewal metadata update.

Auth:
Node Bearer token required.
