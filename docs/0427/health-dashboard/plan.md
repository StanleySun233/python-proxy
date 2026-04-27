# Master Development Plan: Health Dashboard Refactor

## Team

| Role | ID | Plan | Scope |
|------|----|------|-------|
| Frontend | FE-1 | [plan-fe-1.md](./plan-fe-1.md) | ECharts install, sidebar, redirect, overview, heartbeat, certificates, history API client |
| Backend | BE-1 | [plan-be-1.md](./plan-be-1.md) | History table, history API, cleanup job |

## Frontend Tasks

### FE-1 — [plan](./plan-fe-1.md)

- [ ] Install ECharts dependencies — [ref](../health.md#atomic-1)
- [ ] Update sidebar navigation with health sub-menus — [ref](../health.md#atomic-2)
- [ ] Create health redirect page — [ref](../health.md#atomic-3)
- [ ] Add health history API client function — [ref](../health.md#atomic-6)
- [ ] Create Health Overview page with ECharts dashboard — [ref](../health.md#atomic-7)
- [ ] Create Heartbeat Monitor page — [ref](../health.md#atomic-8)
- [ ] Create Certificate Monitor page — [ref](../health.md#atomic-9)

## Backend Tasks

### BE-1 — [plan](./plan-be-1.md)

- [ ] Add node_health_history table + store implementation — [ref](../health.md#atomic-4)
- [ ] Add health history API endpoint — [ref](../health.md#atomic-5)
- [ ] Add health history cleanup job — [ref](../health.md#atomic-10)

## Integration

- [ ] Verify sidebar sub-menu links navigate correctly
- [ ] Verify ECharts charts render in both light and dark themes
- [ ] Verify overview page shows data from history API
- [ ] Verify heartbeat page inline detail with mini trend chart
- [ ] Verify certificate page search + expiry range filter
- [ ] Verify history cleanup job runs in maintenance cycle
