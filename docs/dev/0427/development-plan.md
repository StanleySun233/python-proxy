# Development Plan - 2026-04-27

## Overview

This development plan covers four major features based on product requirements:
1. Node ID Simplification (Backend + Frontend)
2. Bootstrap Approval Flow (Backend + Frontend)
3. Chains Management UI Enhancement (Backend + Frontend)
4. Routes Management UI Enhancement (Backend + Frontend)

## Technology Stack

- **Backend**: Go 1.23.0, GORM, MySQL
- **Frontend**: Next.js 14.2.5, React 18.3.1, TypeScript 5.5.4, @dnd-kit, @tanstack/react-query

## Task Breakdown

### Phase 1: Foundation - Node ID Simplification

**Priority**: HIGHEST - This is a foundational change that affects all other features

#### Task 1.1: Backend - Database Migration for Node ID
**Owner**: Backend Engineer  
**Estimated Time**: 4 hours  
**Dependencies**: None

**Description**:
- Create database migration script to convert `nodes.id` from TEXT to INTEGER AUTO_INCREMENT
- Update all foreign key references in related tables:
  - `nodes.parent_node_id`
  - `node_links.source_node_id`, `node_links.target_node_id`
  - `chain_hops.node_id`
  - `node_policy_assignments.node_id`
  - `node_health_snapshots.node_id`
  - `node_enrollment_approvals.parent_node_id` (for future use)

**Acceptance Criteria**:
- [ ] Migration script successfully converts existing node IDs
- [ ] All foreign key constraints remain valid
- [ ] No data loss during migration
- [ ] Rollback script is available

**Files to Modify**:
- `/apps/one-panel-api/migrations/` (new migration file)

**Commit Message**: `feat(db): migrate node ID from TEXT to INTEGER AUTO_INCREMENT`

---

#### Task 1.2: Backend - Update GORM Models and API Handlers
**Owner**: Backend Engineer  
**Estimated Time**: 3 hours  
**Dependencies**: Task 1.1

**Description**:
- Update GORM model `Node` to use `uint` for ID field
- Update all API handlers to accept and return numeric IDs:
  - `GET /api/v1/nodes`
  - `POST /api/v1/nodes`
  - `PATCH /api/v1/nodes/{nodeId}`
  - `DELETE /api/v1/nodes/{nodeId}`
  - `GET /api/v1/chains` (include numeric node IDs in hops)
  - `GET /api/v1/route-rules` (include numeric node IDs in chain details)

**Acceptance Criteria**:
- [ ] All node APIs accept and return numeric IDs
- [ ] Chain APIs correctly reference numeric node IDs
- [ ] Route APIs correctly reference numeric node IDs
- [ ] API responses match updated schema

**Files to Modify**:
- `/apps/one-panel-api/models/node.go`
- `/apps/one-panel-api/handlers/nodes.go`
- `/apps/one-panel-api/handlers/chains.go`
- `/apps/one-panel-api/handlers/routes.go`

**Commit Message**: `feat(api): update node APIs to use numeric IDs`

---

#### Task 1.3: Frontend - Update TypeScript Interfaces
**Owner**: Frontend Engineer  
**Estimated Time**: 1 hour  
**Dependencies**: Task 1.2

**Description**:
- Update TypeScript interfaces to use `number` for node ID
- Update API client types

**Acceptance Criteria**:
- [ ] All node-related interfaces use `number` for ID
- [ ] No TypeScript compilation errors

**Files to Modify**:
- `/apps/one-proxy-panel/types/node.ts`
- `/apps/one-proxy-panel/types/chain.ts`
- `/apps/one-proxy-panel/types/route.ts`

**Commit Message**: `feat(types): update node ID type to number`

---

#### Task 1.4: Frontend - Update UI Components
**Owner**: Frontend Engineer  
**Estimated Time**: 2 hours  
**Dependencies**: Task 1.3

**Description**:
- Update all UI components to display numeric node IDs
- Update node selection dropdowns to show format: `1 - edge-node-a`
- Update chain hop visualization to show format: `1 → 2 → 3`
- Update route rule display to show numeric node IDs

**Acceptance Criteria**:
- [ ] Node list displays numeric IDs
- [ ] Chain editor displays numeric IDs in hop sequence
- [ ] Route rule editor displays numeric IDs
- [ ] No broken references or display errors

**Files to Modify**:
- `/apps/one-proxy-panel/components/nodes/NodeList.tsx`
- `/apps/one-proxy-panel/components/chains/ChainEditor.tsx`
- `/apps/one-proxy-panel/components/routes/RouteRuleEditor.tsx`

**Commit Message**: `feat(ui): display numeric node IDs across all components`

---

### Phase 2: Bootstrap Approval Flow

#### Task 2.1: Backend - Database Schema for Enrollment Approvals
**Owner**: Backend Engineer  
**Estimated Time**: 2 hours  
**Dependencies**: Task 1.2

**Description**:
- Create `node_enrollment_approvals` table
- Add migration script

**Acceptance Criteria**:
- [ ] Table created with correct schema
- [ ] Foreign key constraints are valid
- [ ] Indexes are created for performance

**Files to Modify**:
- `/apps/one-panel-api/migrations/` (new migration file)
- `/apps/one-panel-api/models/enrollment_approval.go` (new file)

**Commit Message**: `feat(db): add node_enrollment_approvals table`

---

#### Task 2.2: Backend - Implement Approval APIs
**Owner**: Backend Engineer  
**Estimated Time**: 4 hours  
**Dependencies**: Task 2.1

**Description**:
- Implement `GET /api/v1/nodes/approvals`
- Implement `POST /api/v1/nodes/approvals/{approvalId}/approve`
- Implement `POST /api/v1/nodes/approvals/{approvalId}/reject`
- Modify `POST /api/v1/nodes/bootstrap-token` to create approval record
- Modify `POST /api/v1/nodes/enroll` to check approval status

**Acceptance Criteria**:
- [ ] All approval APIs work correctly
- [ ] Bootstrap token creation creates approval record
- [ ] Node enrollment checks approval status
- [ ] Proper error messages for pending/rejected/expired enrollments

**Files to Modify**:
- `/apps/one-panel-api/handlers/approvals.go` (new file)
- `/apps/one-panel-api/handlers/nodes.go`
- `/apps/one-panel-api/routes/routes.go`

**Commit Message**: `feat(api): implement bootstrap approval flow APIs`

---

#### Task 2.3: Frontend - Implement Approvals Page
**Owner**: Frontend Engineer  
**Estimated Time**: 5 hours  
**Dependencies**: Task 2.2

**Description**:
- Create `/nodes/approvals` page
- Implement approval queue table
- Implement approve/reject actions
- Add navigation badge for pending count
- Add dashboard alert for pending enrollments
- Implement auto-refresh (30 seconds)

**Acceptance Criteria**:
- [ ] Approvals page displays all pending enrollments
- [ ] Approve/reject actions work correctly
- [ ] Badge count updates after action
- [ ] Auto-refresh updates queue every 30 seconds
- [ ] Dashboard alert displays when pending enrollments exist

**Files to Modify**:
- `/apps/one-proxy-panel/app/nodes/approvals/page.tsx` (new file)
- `/apps/one-proxy-panel/components/approvals/ApprovalQueueTable.tsx` (new file)
- `/apps/one-proxy-panel/components/layout/Navigation.tsx`
- `/apps/one-proxy-panel/app/dashboard/page.tsx`

**Commit Message**: `feat(ui): implement bootstrap approval flow UI`

---

### Phase 3: Chains Management UI Enhancement

#### Task 3.1: Backend - Implement Chains Validation and Preview APIs
**Owner**: Backend Engineer  
**Estimated Time**: 4 hours  
**Dependencies**: Task 1.2

**Description**:
- Implement `POST /api/v1/chains/validate`
- Implement `POST /api/v1/chains/preview`
- Implement `GET /api/v1/nodes/scopes`
- Add validation logic for hop connectivity, scope ownership, first hop edge node
- Enhance `GET /api/v1/chains` with query parameters and validation status

**Acceptance Criteria**:
- [ ] Validation API returns correct validation results
- [ ] Preview API returns compiled chain configuration
- [ ] Scopes API returns all available scopes with ownership info
- [ ] Chains list API supports filtering and includes validation status

**Files to Modify**:
- `/apps/one-panel-api/handlers/chains.go`
- `/apps/one-panel-api/handlers/nodes.go`
- `/apps/one-panel-api/services/chain_validator.go` (new file)

**Commit Message**: `feat(api): implement chain validation and preview APIs`

---

#### Task 3.2: Frontend - Implement Visual Chain Editor
**Owner**: Frontend Engineer  
**Estimated Time**: 6 hours  
**Dependencies**: Task 3.1

**Description**:
- Enhance `/chains` page with visual chain editor
- Implement drag-and-drop hop reordering using @dnd-kit
- Implement chain preview card
- Implement scope selector with recommendations
- Implement real-time validation
- Implement compilation preview modal

**Acceptance Criteria**:
- [ ] Chain editor opens when clicking "Create" or "Edit"
- [ ] Hop cards can be dragged and reordered
- [ ] Validation runs automatically on input change
- [ ] Scope selector shows recommended scopes
- [ ] Preview modal displays compiled configuration
- [ ] All validation errors and warnings display correctly

**Files to Modify**:
- `/apps/one-proxy-panel/app/chains/page.tsx`
- `/apps/one-proxy-panel/components/chains/ChainEditor.tsx`
- `/apps/one-proxy-panel/components/chains/HopList.tsx` (new file)
- `/apps/one-proxy-panel/components/chains/ScopeSelector.tsx` (new file)
- `/apps/one-proxy-panel/components/chains/CompilationPreviewModal.tsx` (new file)

**Commit Message**: `feat(ui): implement visual chain editor with drag-and-drop`

---

### Phase 4: Routes Management UI Enhancement

#### Task 4.1: Backend - Implement Routes Validation and Suggestions APIs
**Owner**: Backend Engineer  
**Estimated Time**: 4 hours  
**Dependencies**: Task 3.1

**Description**:
- Implement `POST /api/v1/route-rules/validate`
- Implement `GET /api/v1/route-rules/match-types`
- Implement `GET /api/v1/route-rules/suggestions`
- Add match value validation logic (domain, CIDR, port, regex, etc.)
- Enhance `GET /api/v1/route-rules` with query parameters and validation status

**Acceptance Criteria**:
- [ ] Validation API validates match value format correctly
- [ ] Match types API returns all available types with metadata
- [ ] Suggestions API returns relevant suggestions based on match type
- [ ] Route rules list API supports filtering and includes validation status

**Files to Modify**:
- `/apps/one-panel-api/handlers/routes.go`
- `/apps/one-panel-api/services/route_validator.go` (new file)
- `/apps/one-panel-api/services/match_value_validator.go` (new file)

**Commit Message**: `feat(api): implement route validation and suggestions APIs`

---

#### Task 4.2: Frontend - Implement Smart Route Rule Editor
**Owner**: Frontend Engineer  
**Estimated Time**: 6 hours  
**Dependencies**: Task 4.1

**Description**:
- Enhance `/routes` page with smart route rule editor
- Implement match-type-specific input components
- Implement chain selector with preview
- Implement scope selector with recommendations
- Implement real-time validation
- Implement regex tester modal

**Acceptance Criteria**:
- [ ] Route rule editor adapts input based on match type
- [ ] Match value validation works for all types
- [ ] Chain preview card displays when chain is selected
- [ ] Scope selector shows recommended scopes
- [ ] Validation errors display in real-time
- [ ] Regex tester modal works correctly

**Files to Modify**:
- `/apps/one-proxy-panel/app/routes/page.tsx`
- `/apps/one-proxy-panel/components/routes/RouteRuleEditor.tsx`
- `/apps/one-proxy-panel/components/routes/MatchValueInput.tsx` (new file)
- `/apps/one-proxy-panel/components/routes/ChainSelector.tsx` (new file)
- `/apps/one-proxy-panel/components/routes/RegexTesterModal.tsx` (new file)

**Commit Message**: `feat(ui): implement smart route rule editor with validation`

---

## Development Sequence and Dependencies

```
Phase 1: Node ID Simplification (Foundation)
├── Task 1.1: DB Migration (4h)
├── Task 1.2: Backend API Update (3h) [depends on 1.1]
├── Task 1.3: Frontend Types Update (1h) [depends on 1.2]
└── Task 1.4: Frontend UI Update (2h) [depends on 1.3]
Total: 10 hours

Phase 2: Bootstrap Approval Flow
├── Task 2.1: DB Schema (2h) [depends on 1.2]
├── Task 2.2: Backend APIs (4h) [depends on 2.1]
└── Task 2.3: Frontend UI (5h) [depends on 2.2]
Total: 11 hours

Phase 3: Chains Management UI
├── Task 3.1: Backend APIs (4h) [depends on 1.2]
└── Task 3.2: Frontend UI (6h) [depends on 3.1]
Total: 10 hours

Phase 4: Routes Management UI
├── Task 4.1: Backend APIs (4h) [depends on 3.1]
└── Task 4.2: Frontend UI (6h) [depends on 4.1]
Total: 10 hours

Grand Total: 41 hours (~5 working days)
```

## Commit Strategy

Each task completion triggers one commit with the following format:

```
<type>(<scope>): <subject>

<body>

Task: <task-id>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring
- `test`: Adding tests
- `docs`: Documentation changes

Scopes:
- `db`: Database changes
- `api`: Backend API changes
- `ui`: Frontend UI changes
- `types`: TypeScript type changes

## Testing Strategy

### Backend Testing
- Unit tests for validation logic
- Integration tests for API endpoints
- Database migration tests

### Frontend Testing
- Component rendering tests
- User interaction tests
- API integration tests

### End-to-End Testing
After each phase completion:
1. Test complete user flow
2. Verify data consistency
3. Check error handling
4. Validate UI/UX

## Risk Mitigation

### Risk 1: Node ID Migration Data Loss
**Mitigation**: 
- Create full database backup before migration
- Test migration on staging environment
- Prepare rollback script

### Risk 2: Frontend-Backend API Incompatibility
**Mitigation**:
- Deploy backend first with dual format support
- Test API compatibility before frontend deployment
- Use feature flags for gradual rollout

### Risk 3: Complex Drag-and-Drop Implementation
**Mitigation**:
- Use proven library (@dnd-kit)
- Implement progressive enhancement
- Provide keyboard navigation fallback

### Risk 4: Real-Time Validation Performance
**Mitigation**:
- Implement debouncing (500ms)
- Cache validation results
- Use optimistic UI updates

## Deployment Plan

### Pre-Deployment Checklist
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Database backup created
- [ ] Rollback plan documented
- [ ] Staging environment tested

### Deployment Sequence
1. **Phase 1 Deployment** (Node ID Simplification)
   - Schedule maintenance window
   - Run database migration
   - Deploy backend changes
   - Deploy frontend changes
   - Verify all flows work correctly

2. **Phase 2 Deployment** (Bootstrap Approval Flow)
   - Run database migration
   - Deploy backend changes
   - Deploy frontend changes
   - Test approval flow end-to-end

3. **Phase 3 Deployment** (Chains Management UI)
   - Deploy backend changes
   - Deploy frontend changes
   - Test chain editor functionality

4. **Phase 4 Deployment** (Routes Management UI)
   - Deploy backend changes
   - Deploy frontend changes
   - Test route rule editor functionality

### Post-Deployment Verification
- [ ] All APIs respond correctly
- [ ] UI displays correctly
- [ ] No console errors
- [ ] Database integrity maintained
- [ ] Performance metrics acceptable

## Success Metrics

### Functional Metrics
- All acceptance criteria met
- Zero critical bugs in production
- All user flows work end-to-end

### Performance Metrics
- API response time < 200ms (p95)
- Frontend page load time < 2s
- Validation response time < 500ms

### User Experience Metrics
- Reduced operator errors in chain/route configuration
- Improved operator efficiency (time to create chain/route)
- Positive operator feedback on new UI features

## Timeline Estimate

Assuming 1 backend engineer and 1 frontend engineer working in parallel:

- **Week 1 (Days 1-2)**: Phase 1 - Node ID Simplification
- **Week 1 (Days 3-4)**: Phase 2 - Bootstrap Approval Flow
- **Week 2 (Days 1-2)**: Phase 3 - Chains Management UI
- **Week 2 (Days 3-4)**: Phase 4 - Routes Management UI
- **Week 2 (Day 5)**: Integration testing and bug fixes

**Total Duration**: 9-10 working days (2 weeks)

## Notes

- All features are independent after Phase 1 completion
- Phase 2, 3, 4 can be developed in parallel if resources allow
- Each phase should be deployed and tested before moving to next
- Maintain backward compatibility during migration periods
- Document all API changes in API documentation
- Update user documentation after each phase deployment
