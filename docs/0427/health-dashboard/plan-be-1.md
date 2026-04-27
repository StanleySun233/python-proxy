# BE-1 Development Plan

**Master plan**: [plan.md](./plan.md)

## Tasks

- [x] Add node_health_history table + store implementation — [ref](../health.md#atomic-4)
- [ ] Add health history API endpoint — [ref](../health.md#atomic-5)
- [ ] Add health history cleanup job — [ref](../health.md#atomic-10)

## Review Log
- 2026-04-27: BE-1 Task 1 (Atomic-4) completed. Created 004 migration, modified UpsertNodeHeartbeat, added ListNodeHealthHistory to store interface + mysql + seed.
