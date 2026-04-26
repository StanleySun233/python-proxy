# 7 Control Plane Backend

## Responsibilities

- account and session management
- node inventory and enrollment
- chain and route rule management
- policy compilation and publish
- certificate issuance metadata and renewal scheduling
- node health and policy status query APIs

## Go Service Layout

- `cmd/control-plane`: process entrypoint
- `internal/config`: environment and runtime config
- `internal/httpapi`: routers, handlers, middleware
- `internal/service`: business logic
- `internal/store`: GORM-managed MySQL access
- `internal/policy`: compile route rules and chains into node snapshots
- `internal/auth`: password hashing and token flow
- `internal/nodeenroll`: bootstrap and trust handshake
- `internal/cert`: cert lifecycle orchestration

## Runtime Modules

- `api server`: admin-facing HTTP API
- `scheduler`: renewal jobs, stale session cleanup, node health evaluation
- `compiler`: policy revision generation
- `issuer`: bootstrap token and internal certificate issuance

## Configuration Keys

- `HTTP_ADDR`
- `MYSQL_DSN`
- `JWT_SIGNING_KEY`
- `BOOTSTRAP_TOKEN_TTL`
- `NODE_CERT_TTL`
- `PUBLIC_CERT_RENEW_WINDOW`

## Core Workflows

### Admin Bootstrap

- initialize database
- create roles
- create `admin` account
- generate one-time random password
- force password rotation on first login

### Policy Publish

- load enabled rules, chains, nodes, links
- validate chain references and loops
- validate `destination_scope`
- compile snapshot per node
- insert `policy_revision`
- assign node revisions

### Node Enrollment

- admin creates bootstrap token
- node presents token
- backend validates token and creates node record
- backend binds node trust material
- node starts heartbeat and policy sync

## Failure Rules

- invalid policy must never replace the active policy
- disabled nodes must be excluded from chain compilation
- downstream node revocation must invalidate future assignments
