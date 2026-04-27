# Bootstrap Approval Flow

## Background

Current node bootstrap flow creates a token via `/nodes/bootstrap-token`, but provides no visibility into pending nodes that have not yet connected. Operators need a dedicated approval interface to review, approve, or reject pending node enrollment requests before nodes become active in the topology.

## Goal

Add a `/nodes/approvals` page that displays pending node enrollment requests after bootstrap token creation, allowing operators to review and approve nodes before they join the network.

## Functional Requirements

### FR1. Pending Node State Tracking

- After bootstrap token is created, track pending enrollment state
- Store pending node metadata:
  - Bootstrap token ID
  - Target node name (if specified)
  - Target node type (edge/relay)
  - Creation timestamp
  - Expiry timestamp
  - Current status: `pending`, `approved`, `rejected`, `expired`

### FR2. Approval Queue Display

- `/nodes/approvals` page shows all pending enrollment requests
- Display columns:
  - Token ID (first 8 characters)
  - Target node name
  - Node type
  - Created at
  - Expires at
  - Status
  - Actions (Approve/Reject)

### FR3. Approval Actions

- Operator can approve pending enrollment
- Operator can reject pending enrollment
- Approved enrollment allows node to complete exchange and become active
- Rejected enrollment prevents node from joining network

### FR4. Notification and Alerts

- Badge count on `/nodes/approvals` navigation item shows pending count
- Alert banner on dashboard if pending enrollments exist
- Auto-refresh approval queue every 30 seconds

## User Interaction Flow

### Flow 1: Bootstrap Token Creation

1. Operator navigates to `/nodes/bootstrap`
2. Fills in node metadata:
   - Node name: `relay-node-c`
   - Node type: `relay`
   - Parent node: `1 - edge-node-a`
3. Clicks "Generate Bootstrap Token"
4. System creates bootstrap token and pending enrollment record
5. System displays token and shows notification: "Pending enrollment created. Review in Approvals."
6. Badge appears on `/nodes/approvals` navigation item: `(1)`

### Flow 2: Reviewing Pending Enrollments

1. Operator clicks `/nodes/approvals` in navigation
2. Page displays pending enrollment table:
   ```
   Token ID    | Node Name      | Type  | Created At          | Expires At          | Status  | Actions
   a1b2c3d4    | relay-node-c   | relay | 2026-04-27 10:00:00 | 2026-04-27 11:00:00 | pending | [Approve] [Reject]
   ```
3. Operator reviews node details
4. Operator clicks "Approve"
5. System updates enrollment status to `approved`
6. System displays success message: "Enrollment approved. Node can now connect."
7. Badge count decrements

### Flow 3: Node Connection After Approval

1. Node agent runs with bootstrap token
2. Node calls `/api/v1/nodes/enroll` with token
3. Backend checks enrollment status
4. If status is `approved`, proceed with enrollment
5. If status is `pending`, return error: "Enrollment pending approval"
6. If status is `rejected`, return error: "Enrollment rejected"
7. If status is `expired`, return error: "Bootstrap token expired"

## Backend API Requirements

### New API Endpoints

#### `GET /api/v1/nodes/approvals`

Purpose: List pending enrollment requests

Auth: Control plane Bearer token required

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "uuid-1",
      "token_id": "a1b2c3d4-e5f6-7890",
      "target_node_name": "relay-node-c",
      "target_node_type": "relay",
      "parent_node_id": 1,
      "status": "pending",
      "created_at": "2026-04-27T10:00:00Z",
      "expires_at": "2026-04-27T11:00:00Z"
    }
  ]
}
```

#### `POST /api/v1/nodes/approvals/{approvalId}/approve`

Purpose: Approve pending enrollment

Auth: Control plane Bearer token required

Request:
```json
{
  "operator_note": "Approved for production relay"
}
```

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "uuid-1",
    "status": "approved",
    "approved_at": "2026-04-27T10:05:00Z",
    "approved_by": "admin"
  }
}
```

#### `POST /api/v1/nodes/approvals/{approvalId}/reject`

Purpose: Reject pending enrollment

Auth: Control plane Bearer token required

Request:
```json
{
  "operator_note": "Duplicate request"
}
```

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "uuid-1",
    "status": "rejected",
    "rejected_at": "2026-04-27T10:05:00Z",
    "rejected_by": "admin"
  }
}
```

### Database Schema Changes

Add new table `node_enrollment_approvals`:

```sql
CREATE TABLE node_enrollment_approvals (
  id TEXT PRIMARY KEY,
  bootstrap_token_id TEXT NOT NULL,
  target_node_name TEXT,
  target_node_type TEXT NOT NULL,
  parent_node_id INTEGER,
  status TEXT NOT NULL,
  operator_note TEXT,
  approved_by TEXT,
  approved_at TEXT,
  rejected_by TEXT,
  rejected_at TEXT,
  created_at TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  FOREIGN KEY (bootstrap_token_id) REFERENCES bootstrap_tokens(id),
  FOREIGN KEY (parent_node_id) REFERENCES nodes(id)
);
```

### Modified API Behavior

#### `POST /api/v1/nodes/bootstrap-token`

After creating bootstrap token, also create enrollment approval record:

```go
// Create bootstrap token
token := createBootstrapToken(...)

// Create enrollment approval record
approval := NodeEnrollmentApproval{
  ID: uuid.New(),
  BootstrapTokenID: token.ID,
  TargetNodeName: req.NodeName,
  TargetNodeType: req.NodeType,
  ParentNodeID: req.ParentNodeID,
  Status: "pending",
  CreatedAt: time.Now(),
  ExpiresAt: token.ExpiresAt,
}
db.Create(&approval)
```

#### `POST /api/v1/nodes/enroll`

Check enrollment approval status before proceeding:

```go
// Verify bootstrap token
token := verifyBootstrapToken(req.Token)

// Check approval status
approval := getApprovalByTokenID(token.ID)
if approval.Status != "approved" {
  return error("enrollment_not_approved")
}

// Proceed with enrollment
...
```

## Frontend Requirements

### New Page: `/nodes/approvals`

Component structure:
```
/nodes/approvals
├── ApprovalQueueTable
│   ├── ApprovalRow
│   │   ├── TokenIDCell
│   │   ├── NodeNameCell
│   │   ├── NodeTypeCell
│   │   ├── TimestampCell
│   │   ├── StatusBadge
│   │   └── ActionButtons
│   └── EmptyState
└── ApprovalDetailDrawer
```

### Navigation Badge

Update navigation component to show pending count:

```tsx
<NavItem href="/nodes/approvals">
  Approvals
  {pendingCount > 0 && <Badge>{pendingCount}</Badge>}
</NavItem>
```

### Dashboard Alert

Add alert banner on dashboard:

```tsx
{pendingApprovals.length > 0 && (
  <Alert variant="info">
    {pendingApprovals.length} pending node enrollment(s) require approval.
    <Link href="/nodes/approvals">Review now</Link>
  </Alert>
)}
```

### Auto-Refresh

Implement polling for approval queue:

```tsx
useEffect(() => {
  const interval = setInterval(() => {
    fetchApprovals();
  }, 30000); // 30 seconds
  return () => clearInterval(interval);
}, []);
```

## Acceptance Criteria

### AC1. Bootstrap Token Creation

- [ ] Creating bootstrap token creates pending enrollment record
- [ ] Pending enrollment record contains correct metadata
- [ ] Badge appears on `/nodes/approvals` navigation item

### AC2. Approval Queue Display

- [ ] `/nodes/approvals` page displays all pending enrollments
- [ ] Table shows correct columns and data
- [ ] Auto-refresh updates queue every 30 seconds
- [ ] Empty state displays when no pending enrollments

### AC3. Approval Actions

- [ ] Clicking "Approve" updates enrollment status to `approved`
- [ ] Clicking "Reject" updates enrollment status to `rejected`
- [ ] Success message displays after action
- [ ] Badge count updates after action

### AC4. Node Enrollment Flow

- [ ] Node with approved enrollment can complete enrollment
- [ ] Node with pending enrollment receives error
- [ ] Node with rejected enrollment receives error
- [ ] Node with expired token receives error

### AC5. Dashboard Integration

- [ ] Dashboard shows alert when pending enrollments exist
- [ ] Alert links to `/nodes/approvals` page
- [ ] Alert disappears when no pending enrollments

## Security Considerations

### Authorization

- Only operators with `admin` role can approve/reject enrollments
- Approval actions are logged with operator identity
- Rejected enrollments cannot be re-approved without creating new token

### Token Expiry

- Expired tokens automatically mark enrollment as `expired`
- Expired enrollments cannot be approved
- Background job cleans up expired enrollments after 7 days

### Audit Trail

- All approval/rejection actions are logged
- Logs include:
  - Operator identity
  - Action timestamp
  - Enrollment ID
  - Operator note

## Future Enhancements

- Auto-approval rules based on node metadata
- Approval workflow with multiple reviewers
- Slack/email notifications for pending enrollments
- Bulk approval/rejection actions
- Enrollment request history and audit log viewer
