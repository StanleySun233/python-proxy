# Node ID Simplification

## Background

Current node ID format uses `node-{timestamp}` pattern, which creates unnecessarily long and hard-to-read identifiers. For a multi-node proxy system where operators frequently reference nodes in chains, routes, and topology views, simpler numeric IDs improve usability and reduce cognitive load.

## Goal

Replace timestamp-based node IDs with auto-incrementing numeric IDs (1, 2, 3, 4...) to improve readability and operator experience across the admin console and API responses.

## Functional Requirements

### FR1. Auto-Incrementing Node ID

- Node ID must be a positive integer starting from 1
- Each new node receives the next available ID in sequence
- ID assignment must be atomic and collision-free
- Deleted node IDs should not be reused in V1

### FR2. Database Schema Change

- Change `nodes.id` column from `TEXT` to `INTEGER PRIMARY KEY AUTO_INCREMENT`
- Update all foreign key references in related tables:
  - `nodes.parent_node_id`
  - `node_links.source_node_id` and `target_node_id`
  - `chain_hops.node_id`
  - `node_policy_assignments.node_id`
  - `node_health_snapshots.node_id`

### FR3. API Response Format

- All node-related API responses must return numeric ID
- Example: `{"id": 1, "name": "edge-node-a", ...}`
- Maintain backward compatibility during migration by supporting both formats temporarily

### FR4. Frontend Display

- Display node ID as plain number in all UI components
- Node selection dropdowns show: `1 - edge-node-a`
- Chain hop display shows: `1 → 2 → 3`
- Route rule display shows: `Chain: 1 → 2 (scope: k8s-prod)`

## User Interaction Flow

### Operator Creating a Chain

1. Navigate to `/chains`
2. Click "Create Chain"
3. Select nodes from dropdown showing: `1 - edge-node-a`, `2 - relay-node-b`
4. Drag to reorder hops
5. Save chain
6. Chain list displays: `Chain: prod-path (1 → 2 → 3)`

### Operator Viewing Node List

1. Navigate to `/nodes`
2. Node table displays:
   - ID column: `1`, `2`, `3`
   - Name column: `edge-node-a`, `relay-node-b`
   - Status column: `active`, `pending`

## Backend API Requirements

### API Changes

#### `GET /api/v1/nodes`

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "edge-node-a",
      "mode": "edge",
      "status": "active"
    }
  ]
}
```

#### `POST /api/v1/nodes`

Request:
```json
{
  "name": "new-relay-node",
  "mode": "relay",
  "parent_node_id": 1
}
```

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 4,
    "name": "new-relay-node"
  }
}
```

#### `PATCH /api/v1/nodes/{nodeId}`

- Path parameter `nodeId` accepts integer
- Example: `PATCH /api/v1/nodes/1`

#### `DELETE /api/v1/nodes/{nodeId}`

- Path parameter `nodeId` accepts integer
- Example: `DELETE /api/v1/nodes/1`

### Database Migration

Backend must provide migration script:

```sql
-- Step 1: Create new table with integer ID
CREATE TABLE nodes_new (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  name TEXT NOT NULL,
  mode TEXT NOT NULL,
  public_host TEXT,
  public_port INTEGER,
  scope_key TEXT NOT NULL,
  parent_node_id INTEGER,
  enabled INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY (parent_node_id) REFERENCES nodes_new(id)
);

-- Step 2: Migrate data with ID mapping
-- Step 3: Update foreign key references in related tables
-- Step 4: Drop old table and rename new table
```

## Frontend Requirements

### Component Updates

1. Node selection components must handle numeric ID
2. Chain hop display must format as `1 → 2 → 3`
3. Route rule display must show numeric node IDs
4. Node detail drawer must display numeric ID prominently

### API Client Updates

- Update TypeScript interfaces:
  ```typescript
  interface Node {
    id: number;  // changed from string
    name: string;
    mode: 'edge' | 'relay';
    status: string;
  }
  ```

## Acceptance Criteria

### AC1. Database Migration

- [ ] Migration script successfully converts existing node IDs
- [ ] All foreign key constraints remain valid
- [ ] No data loss during migration

### AC2. API Compatibility

- [ ] All node APIs accept and return numeric IDs
- [ ] Chain APIs correctly reference numeric node IDs
- [ ] Route APIs correctly reference numeric node IDs

### AC3. Frontend Display

- [ ] Node list displays numeric IDs
- [ ] Chain editor displays numeric IDs in hop sequence
- [ ] Route rule editor displays numeric IDs
- [ ] No broken references or display errors

### AC4. End-to-End Flow

- [ ] Create new node receives auto-incremented ID
- [ ] Create chain with numeric node IDs succeeds
- [ ] Create route rule with numeric node IDs succeeds
- [ ] Policy compilation uses numeric node IDs correctly

## Migration Strategy

### Phase 1: Backend Migration

1. Create database migration script
2. Update GORM models to use `uint` for ID fields
3. Update all API handlers to accept numeric IDs
4. Test migration on staging database

### Phase 2: Frontend Migration

1. Update TypeScript interfaces
2. Update API client to handle numeric IDs
3. Update all UI components
4. Test all node-related flows

### Phase 3: Deployment

1. Schedule maintenance window
2. Run database migration
3. Deploy backend changes
4. Deploy frontend changes
5. Verify all flows work correctly

## Risks and Mitigations

### Risk 1: Data Loss During Migration

Mitigation: Create full database backup before migration, test migration script on staging environment

### Risk 2: Frontend-Backend Incompatibility

Mitigation: Deploy backend first with dual format support, then deploy frontend, then remove legacy format support

### Risk 3: Existing Node References Break

Mitigation: Comprehensive testing of all node reference paths (chains, routes, policies) before production deployment
