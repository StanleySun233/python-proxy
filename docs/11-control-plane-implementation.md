# 11 Control Plane Implementation

## Current Implementation Target

Build the control plane in layers:

- `config`: process configuration
- `domain`: core objects returned by APIs
- `store`: data access boundary
- `service`: use-case composition
- `httpapi`: transport layer

## V1 Execution Order

### Step 1

- implement read-only overview APIs with stable response shapes
- expose nodes, chains, route rules, and node health

### Step 2

- add account bootstrap and login
- add password hashing and token issuing

### Step 3

- add create/update APIs for nodes, chains, and route rules
- add policy publish flow

### Step 4

- replace seed store with MySQL-backed store
- keep seed store as emergency fallback only during early integration

### Step 5

- expose full control-plane CRUD and node-agent sync routes
- stabilize payload shapes for frontend integration

## Backend Contract Rules

- handlers should not know about SQL
- service layer should return domain objects
- store layer should be replaceable from seed-memory to MySQL
- API response shapes should stabilize before frontend wiring
