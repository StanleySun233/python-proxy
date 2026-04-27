# Master Development Plan: Accounts Management Refactor

## Team

| Role | ID | Plan | Scope |
|------|----|------|-------|
| Frontend | FE-1 | [plan-fe-1.md](./plan-fe-1.md) | Sidebar nav, create page, list page, edit modal, delete API client |
| Backend | BE-1 | [plan-be-1.md](./plan-be-1.md) | DELETE account endpoint + admin protection |

## Frontend Tasks

### FE-1 — [plan](./plan-fe-1.md)

- [x] Update sidebar navigation with accounts sub-menus — [ref](../accounts.md#atomic-1)
- [x] Create accounts redirect page — [ref](../accounts.md#atomic-2)
- [x] Create dedicated "Create Account" page — [ref](../accounts.md#atomic-3)
- [x] Create Edit Account Modal component — [ref](../accounts.md#atomic-5)
- [x] Add deleteAccount API client function — [ref](../accounts.md#atomic-6)
- [x] Create "Account List" page with query, edit, and delete — [ref](../accounts.md#atomic-7)

## Backend Tasks

### BE-1 — [plan](./plan-be-1.md)

- [x] Add DELETE account endpoint with admin protection — [ref](../accounts.md#atomic-4)

## Integration

- [ ] Verify all sub-menu links navigate correctly
- [ ] Verify create → list navigation flow (create account, see it in list)
- [ ] Verify edit modal opens, pre-fills data, saves correctly
- [ ] Verify delete: non-admin account can be deleted, admin account rejected
- [ ] Verify backend DELETE returns proper error for admin account
