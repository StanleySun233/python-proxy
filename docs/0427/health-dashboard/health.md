# Health Dashboard — Multi-Submenu Refactor with ECharts & Historical Trends

## Overview

Refactor the single-page health dashboard (`/health`) into a multi-submenu section (following the Nodes pattern). Replace plain metric cards with rich ECharts visualizations, split heartbeat and certificate views into dedicated query pages, and add backend historical data APIs for time-series trend charts.

## Current State

- Single page at `/health` with: 4 summary metric cards, heartbeat registry table, certificate registry table
- Sidebar: one sub-menu `健康面板` → `/health`
- Backend: `GET /api/v1/nodes/health` (latest snapshot per node), `GET /api/v1/certificates`, `GET /api/v1/nodes`
- No historical data is stored — `node_health_snapshots` table uses `ON DUPLICATE KEY UPDATE` (upsert)
- `recharts` is installed but unused; zero chart components exist in the project

## Target Architecture

```
/health
├── /health/overview      — ECharts dashboard with summary charts + trend lines
├── /health/heartbeat     — Node heartbeat query with filtering (moved from current page)
└── /health/certificates  — Certificate query with filtering (moved from current page)
```

## Atomic Requirements

### Atomic-1: Install ECharts dependencies

**Scope**: Frontend — `package.json`

- Add `echarts` (latest v5)
- Add `echarts-for-react` (latest v3)

Run: `npm install echarts echarts-for-react` in `apps/one-proxy-panel/`.

### Atomic-2: Update sidebar navigation for health section

**Scope**: Frontend — `console-shell.tsx`, `zh.json`, `en.json`

Change the `health` nav section from a single sub-item to three sub-items:
- `健康总览` (`shell.healthOverview`) → `/health/overview`
- `心跳监控` (`shell.healthHeartbeat`) → `/health/heartbeat`
- `证书监控` (`shell.healthCertificates`) → `/health/certificates`

Parent section href: `/health/overview` (was `/health`).

### Atomic-3: Create health redirect page

**Scope**: Frontend — `/health/page.tsx`

Replace the current health page with a simple redirect to `/health/overview`.

### Atomic-4: Backend — Add node_health_history table

**Scope**: Backend — schema migration + store layer

Create new migration SQL file `schema/004_node_health_history.sql`:
```sql
CREATE TABLE IF NOT EXISTS node_health_history (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  node_id VARCHAR(191) NOT NULL,
  heartbeat_at VARCHAR(64) NOT NULL,
  policy_revision_id VARCHAR(191) DEFAULT '',
  listener_status_json LONGTEXT NOT NULL,
  cert_status_json LONGTEXT NOT NULL,
  created_at VARCHAR(64) NOT NULL,
  INDEX idx_history_node_time (node_id, heartbeat_at),
  INDEX idx_history_time (heartbeat_at)
);
```

Modify `UpsertNodeHeartbeat` in `mysql.go` to INSERT into `node_health_history` (in the same transaction) after upserting the snapshot.

Add to `Store` interface:
- `ListNodeHealthHistory(nodeID string, window time.Duration) ([]domain.NodeHealth, error)`

Implement in `mysql.go`: `SELECT ... FROM node_health_history WHERE node_id = ? AND heartbeat_at >= ? ORDER BY heartbeat_at`.

Implement in `seed.go`: return empty slice.

### Atomic-5: Backend — Add health history API endpoint

**Scope**: Backend — router + handler + service

- Register route: `GET /api/v1/nodes/health/history` → `r.handleNodeHealthHistory` (with account auth)
- Handler in `resources.go`: parse `nodeId` (required) and `window` (optional, default `24h`) query params, call service
- Service in `controlplane.go`: `NodeHealthHistory(nodeID string, window time.Duration) ([]domain.NodeHealth, error)`
- Returns time-series `NodeHealth[]` suitable for ECharts line chart

### Atomic-6: Frontend — Add health history API client

**Scope**: Frontend — `control-plane-api.ts`

Add:
```typescript
export function getNodeHealthHistory(accessToken: string, nodeId: string, window?: string) {
  const params = new URLSearchParams({nodeId});
  if (window) params.set('window', window);
  return request<NodeHealth[]>(`/nodes/health/history?${params.toString()}`, {accessToken});
}
```

### Atomic-7: Create Health Overview page with ECharts dashboard

**Scope**: Frontend — `/health/overview/page.tsx`

The main dashboard page featuring ECharts visualizations:

1. **Summary metric row** (top): 4 cards — healthy, stale, unreported, cert pressure (reused from current page, restyled)
2. **Node Health Distribution pie chart** (left): Shows healthy / stale / degraded / unreported proportions
3. **Health Trend line chart** (center, full-width): Shows healthy/degraded counts over time window (using history API). 
   - Has a node selector dropdown to filter by specific node
   - Default: aggregate all nodes
4. **Certificate Status bar chart** (bottom): Shows count per certificate status (healthy, renew-soon, rotate, failed, expired)
5. **Certificate Expiry Timeline** (bottom-right): Horizontal bar chart showing days until expiry per certificate

Use `echarts-for-react` `<ReactECharts>` component. Each chart is a separate `<ReactECharts option={...} />`.

Theme: match the existing dark theme. Use the `echarts` dark theme or manually configure colors to match CSS variables.

### Atomic-8: Create Heartbeat Monitor page

**Scope**: Frontend — `/health/heartbeat/page.tsx`

Extract the heartbeat registry section from the old health page. Enhancements:
- Same search + derived state filter (carried over from old page)
- Click on a node row → expand inline detail panel showing:
  - Listener status badges (individual key:value with colored badges)
  - Certificate status badges (individual key:value with colored badges)
  - Mini line chart of that node's heartbeat history (last 24h) using ECharts
- Keep the same data table columns

### Atomic-9: Create Certificate Monitor page

**Scope**: Frontend — `/health/certificates/page.tsx`

Extract the certificate registry section from the old health page. Enhancements:
- Same status filter (carried over from old page)
- Add text search (search by owner name, cert type, provider, ID)
- Add a "Valid to" date range filter (show certs expiring within N days)
- Show a small ECharts gauge or progress bar visualizing time until expiry per certificate

### Atomic-10: Backend — Add health history cleanup job

**Scope**: Backend — maintenance

Add to `RunMaintenance()`:
- `CleanupNodeHealthHistory(retention time.Duration)` — DELETE rows older than N days (default: 7 days)
- Add to `Store` interface: `CleanupNodeHealthHistory(retention time.Duration) (int64, error)`
- Implement in `mysql.go`: `DELETE FROM node_health_history WHERE heartbeat_at < ?`

## User Stories

1. As an operator, I want a visual dashboard with charts to quickly assess overall system health.
2. As an operator, I want to drill into a specific node's heartbeat history trend over time.
3. As an operator, I want to query and filter nodes by health status to troubleshoot issues.
4. As an operator, I want to query certificates by expiry window to proactively renew.
5. As an operator, I want the certificate lifecycle visualized so I can plan renewals.

## Acceptance Criteria

- Sidebar shows 3 sub-menu items: 健康总览, 心跳监控, 证书监控
- `/health` redirects to `/health/overview`
- Overview page has: 4 metric cards + at least 3 ECharts visualizations
- ECharts charts render correctly in both light and dark themes
- Heartbeat page supports search + status filter + inline node detail expansion
- Certificate page supports search + status filter + expiry range filter
- Backend stores heartbeat history and exposes `/api/v1/nodes/health/history`
- History data older than 7 days is automatically cleaned up
- All pages handle loading, empty, and error states

## Dependencies

- Atomic-3 depends on Atomic-2 (redirect target must exist)
- Atomic-6 depends on Atomic-4 + Atomic-5 (API must exist)
- Atomic-7 depends on Atomic-1 + Atomic-6 (charts lib + history API)
- Atomic-8 depends on Atomic-6 (history API for mini trend chart)
- Atomic-9 is independent
- Atomic-10 depends on Atomic-4 (history table must exist)
