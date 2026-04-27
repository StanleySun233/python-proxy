# 15 Node Onboarding API

## Principles

- onboarding APIs belong to control plane
- onboarding path and onboarding task are separate resources
- path defines reachability
- task represents one operator-triggered action

## Node Access Path APIs

### `GET /api/v1/node-access-paths`

Purpose:
List saved onboarding paths.

Auth:
Control plane Bearer token required.

### `POST /api/v1/node-access-paths`

Purpose:
Create a saved onboarding path.

Auth:
Control plane Bearer token required.

Request:

```json
{
  "name": "panel-to-node2-via-node1",
  "mode": "relay_chain",
  "targetNodeId": "node2",
  "entryNodeId": "node1",
  "relayNodeIds": ["node1"],
  "targetHost": "10.20.0.12",
  "targetPort": 2888
}
```

### `PATCH /api/v1/node-access-paths/{pathId}`

Purpose:
Update a saved onboarding path.

Auth:
Control plane Bearer token required.

### `DELETE /api/v1/node-access-paths/{pathId}`

Purpose:
Delete a saved onboarding path.

Auth:
Control plane Bearer token required.

## Node Onboarding Task APIs

### `GET /api/v1/node-onboarding-tasks`

Purpose:
List onboarding tasks.

Auth:
Control plane Bearer token required.

### `POST /api/v1/node-onboarding-tasks`

Purpose:
Create a new onboarding task from panel.

Auth:
Control plane Bearer token required.

Request for direct push:

```json
{
  "mode": "direct",
  "targetNodeId": "node2",
  "targetHost": "node2.internal.example.com",
  "targetPort": 2888
}
```

Request for relay push:

```json
{
  "mode": "relay_chain",
  "pathId": "path-relay-node2",
  "targetNodeId": "node2"
}
```

Request for upstream pull:

```json
{
  "mode": "upstream_pull",
  "pathId": "path-node2-via-node1",
  "targetNodeId": "node2"
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "task-123",
    "mode": "relay_chain",
    "pathId": "path-relay-node2",
    "targetNodeId": "node2",
    "targetHost": "",
    "targetPort": 0,
    "status": "planned",
    "statusMessage": "task created",
    "createdAt": "2026-04-25T14:00:00Z",
    "updatedAt": "2026-04-25T14:00:00Z"
  }
}
```

## Validation Rules

- `mode` must be one of `direct`, `relay_chain`, `upstream_pull`
- `pathId` is required for `relay_chain` and `upstream_pull`
- `targetHost` and `targetPort` are required for `direct`
- `targetNodeId` is strongly recommended for all modes

## Execution Contract

Current behavior:

- `direct` task performs a synchronous `GET /healthz` probe and returns `connected` or `failed`
- `upstream_pull` task validates path existence and returns `pending`
- `relay_chain` task validates path existence, executes ordered relay probe, and returns `connected` or `failed`

Future worker execution can update:

- `status`
- `statusMessage`
- final node lifecycle transitions
