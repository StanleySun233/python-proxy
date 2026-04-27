# FE-2 Development Plan

**Master plan**: [plan.md](./plan.md)

## Tasks

- [x] Add health history API client function — [ref](../health.md#atomic-6)
- [ ] Create Heartbeat Monitor page — [ref](../health.md#atomic-8)
- [x] Create Certificate Monitor page — [ref](../health.md#atomic-9)

## Review Log

### 2026-04-27 — Atomic-6 (getNodeHealthHistory API client)
- Added `getNodeHealthHistory` export function in `control-plane-api.ts`
- Supports optional `window` query param for time range

### 2026-04-27 — Atomic-9 (Certificate Monitor page)
- Created `/health/certificates/page.tsx` with status filter, text search, expiry range filter
- Added expiry visualization with colored dots and remaining days count
- Added CSS for `.expiry-indicator`, `.expiry-dot`, `.expiry-days` in globals.css
- States: loading, error (with retry), empty, filtered-empty
<!-- PM fills this on approval -->
