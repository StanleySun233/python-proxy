# 12 Control Plane API Payloads

## `POST /api/v1/auth/login`

Request:

```json
{
  "account": "admin",
  "password": "<bootstrap-generated-password>"
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "account": {
      "id": "acct-admin",
      "account": "admin",
      "role": "super_admin",
      "status": "active",
      "mustRotatePassword": true
    },
    "accessToken": "access-token",
    "refreshToken": "refresh-token",
    "expiresAt": "2026-04-25T14:00:00Z",
    "mustRotatePassword": true
  }
}
```

## `POST /api/v1/auth/refresh`

Request:

```json
{
  "refreshToken": "refresh-token"
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "account": {
      "id": "acct-admin",
      "account": "admin",
      "role": "super_admin",
      "status": "active",
      "mustRotatePassword": true
    },
    "accessToken": "next-access-token",
    "refreshToken": "next-refresh-token",
    "expiresAt": "2026-04-25T16:00:00Z",
    "mustRotatePassword": true
  }
}
```

## `POST /api/v1/auth/logout`

Auth:

```text
Authorization: Bearer access-token
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "status": "logged_out"
  }
}
```

## `GET /api/v1/overview`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "nodes": {
      "healthy": 3,
      "degraded": 1
    },
    "policies": {
      "activeRevision": "rev-0007",
      "publishedAt": "2026-04-25T00:00:00Z"
    },
    "certificates": {
      "renewSoon": 1
    }
  }
}
```

## `GET /api/v1/accounts`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "acct-admin",
      "account": "admin",
      "role": "super_admin",
      "status": "active",
      "mustRotatePassword": true
    }
  ]
}
```

## `POST /api/v1/accounts`

Request:

```json
{
  "account": "ops-user",
  "password": "change-me",
  "role": "operator"
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "acct-123",
    "account": "ops-user",
    "role": "operator",
    "status": "active",
    "mustRotatePassword": false
  }
}
```

## `PATCH /api/v1/accounts/{accountId}`

Request:

```json
{
  "password": "rotate-now",
  "role": "admin",
  "status": "active"
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "acct-admin",
    "account": "admin",
    "role": "admin",
    "status": "active",
    "mustRotatePassword": false
  }
}
```

## `GET /api/v1/nodes`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "edge-a",
      "name": "edge-a",
      "mode": "edge",
      "scopeKey": "public-edge",
      "parentNodeId": "",
      "enabled": true,
      "status": "healthy",
      "publicHost": "edge-a.example.com",
      "publicPort": 443
    }
  ]
}
```

## `POST /api/v1/nodes`

Request:

```json
{
  "name": "relay-d",
  "mode": "relay",
  "scopeKey": "d-office",
  "parentNodeId": "edge-a",
  "publicHost": "",
  "publicPort": 0
}
```

## `PATCH /api/v1/nodes/{nodeId}`

Request:

```json
{
  "name": "relay-d",
  "mode": "relay",
  "scopeKey": "d-office",
  "parentNodeId": "edge-a",
  "publicHost": "",
  "publicPort": 0,
  "enabled": true,
  "status": "healthy"
}
```

## `POST /api/v1/nodes/bootstrap-token`

Request:

```json
{
  "targetType": "node",
  "targetId": ""
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "bootstrap-123",
    "token": "bootstrap-token",
    "targetType": "node",
    "targetId": "",
    "expiresAt": "2026-04-25T12:15:00Z"
  }
}
```

## `POST /api/v1/nodes/enroll`

Request:

```json
{
  "token": "bootstrap-token",
  "name": "relay-e",
  "mode": "relay",
  "scopeKey": "e-lan",
  "parentNodeId": "edge-a",
  "publicHost": "",
  "publicPort": 0
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "node": {
      "id": "node-123",
      "name": "relay-e",
      "mode": "relay",
      "scopeKey": "e-lan",
      "parentNodeId": "edge-a",
      "enabled": true,
      "status": "pending"
    },
    "enrollmentSecret": "enrollment-secret",
    "approvalState": "pending"
  }
}
```

## `POST /api/v1/nodes/approve/{nodeId}`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "node": {
      "id": "node-123",
      "name": "relay-e",
      "mode": "relay",
      "scopeKey": "e-lan",
      "parentNodeId": "edge-a",
      "enabled": true,
      "status": "degraded"
    },
    "accessToken": "node-access-token",
    "trustMaterial": "shared-secret",
    "expiresAt": "2026-05-25T12:00:00Z"
  }
}
```

## `POST /api/v1/nodes/exchange`

Request:

```json
{
  "nodeId": "node-123",
  "enrollmentSecret": "enrollment-secret"
}
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "node": {
      "id": "node-123",
      "name": "relay-e",
      "mode": "relay",
      "scopeKey": "e-lan",
      "parentNodeId": "edge-a",
      "enabled": true,
      "status": "degraded"
    },
    "accessToken": "node-access-token",
    "trustMaterial": "shared-secret",
    "expiresAt": "2026-05-25T12:00:00Z"
  }
}
```

## `GET /api/v1/node-links`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "link-edge-a-relay-b",
      "sourceNodeId": "edge-a",
      "targetNodeId": "relay-b",
      "linkType": "parent_child",
      "trustState": "trusted"
    }
  ]
}
```

## `POST /api/v1/node-links`

Request:

```json
{
  "sourceNodeId": "edge-a",
  "targetNodeId": "relay-e",
  "linkType": "parent_child",
  "trustState": "trusted"
}
```

## `GET /api/v1/chains`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "chain-corp-k8s",
      "name": "corp-k8s",
      "destinationScope": "c-k8s",
      "enabled": true,
      "hops": ["edge-a", "relay-b", "relay-c"]
    }
  ]
}
```

## `POST /api/v1/chains`

Request:

```json
{
  "name": "office-tools",
  "destinationScope": "d-office",
  "hops": ["edge-a", "relay-d"]
}
```

## `PATCH /api/v1/chains/{chainId}`

Request:

```json
{
  "name": "corp-k8s",
  "destinationScope": "c-k8s",
  "hops": ["edge-a", "relay-b", "relay-c"],
  "enabled": true
}
```

## `GET /api/v1/route-rules`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "rule-corp-domain",
      "priority": 100,
      "matchType": "domain_suffix",
      "matchValue": ".corp.internal",
      "actionType": "chain",
      "chainId": "chain-corp-k8s",
      "destinationScope": "c-k8s",
      "enabled": true
    }
  ]
}
```

## `POST /api/v1/route-rules`

Request:

```json
{
  "priority": 300,
  "matchType": "domain",
  "matchValue": "grafana.office.local",
  "actionType": "chain",
  "chainId": "chain-office-tools",
  "destinationScope": "d-office"
}
```

## `PATCH /api/v1/route-rules/{ruleId}`

Request:

```json
{
  "priority": 100,
  "matchType": "domain_suffix",
  "matchValue": ".corp.internal",
  "actionType": "chain",
  "chainId": "chain-corp-k8s",
  "destinationScope": "c-k8s",
  "enabled": true
}
```

## `GET /api/v1/policies/revisions`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "policy-rev-0007",
      "version": "rev-0007",
      "status": "published",
      "createdAt": "2026-04-25T00:00:00Z",
      "assignedNodes": 4
    }
  ]
}
```

## `POST /api/v1/policies/publish`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "policy-1745577600",
    "version": "rev-1745577600",
    "status": "published",
    "createdAt": "2026-04-25T12:00:00Z",
    "assignedNodes": 4
  }
}
```

## `GET /api/v1/certificates`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "cert-edge-a-public",
      "ownerType": "node",
      "ownerId": "edge-a",
      "certType": "public",
      "provider": "lets_encrypt",
      "status": "renew-soon",
      "notBefore": "2026-04-01T00:00:00Z",
      "notAfter": "2026-05-13T00:00:00Z"
    }
  ]
}
```

## `GET /api/v1/nodes/health`

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "nodeId": "edge-a",
      "heartbeatAt": "2026-04-25T12:00:00Z",
      "policyRevisionId": "rev-0007",
      "listenerStatus": {
        "http": "up",
        "https": "up"
      },
      "certStatus": {
        "public": "renew-soon",
        "internal": "healthy"
      }
    }
  ]
}
```

## `GET /api/v1/node-agent/policy`

Auth:

```text
Authorization: Bearer node-access-token
```

Response:

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "nodeId": "edge-a",
    "policyRevisionId": "rev-0007",
    "payloadJson": "{\"nodes\":[{\"id\":\"edge-a\"}],\"links\":[],\"chains\":[],\"routeRules\":[]}"
  }
}
```

## `POST /api/v1/node-agent/heartbeat`

Auth:

```text
Authorization: Bearer node-access-token
```

Request:

```json
{
  "policyRevisionId": "rev-0007",
  "listenerStatus": {
    "http": "healthy"
  },
  "certStatus": {
    "public": "healthy"
  }
}
```

## `POST /api/v1/node-agent/cert/renew`

Auth:

```text
Authorization: Bearer node-access-token
```

Request:

```json
{
  "certType": "internal"
}
```
