# Accounts Management — Menu & Sub-Menu Refactor

## Overview

Refactor the current single-page accounts management (`/accounts`) to follow the same menu/sub-menu pattern used by the Nodes section. Split the single page into dedicated **Create** and **Query** sub-pages, add **Edit** capability via modal dialog, and add **Delete** with admin account protection.

## Current State

- Single page at `/accounts` with a two-column layout: create form (left) + accounts list (right)
- Backend supports: `GET /accounts`, `POST /accounts`, `PATCH /accounts/:id`
- Backend does NOT support DELETE for accounts
- API client has: `getAccounts()`, `createAccount()`, `updateAccount()`

## Atomic Requirements

### Atomic-1: Update sidebar navigation for accounts section

**Scope**: Frontend — `console-shell.tsx`, `zh.json`, `en.json`

Change the accounts nav section from a single sub-item to two sub-items:
- `创建账号` (`shell.accountCreate`) → `/accounts/create`
- `账号查询` (`shell.accountList`) → `/accounts/list`

The parent section href changes from `/accounts` to `/accounts/create` (following the nodes pattern where `/nodes` href is `/nodes/connect`).

### Atomic-2: Create accounts redirect page

**Scope**: Frontend — `/accounts/page.tsx`

Replace the current full accounts page with a simple redirect to `/accounts/create` (matching the nodes pattern where `/nodes` redirects to `/nodes/connect`).

### Atomic-3: Create dedicated "Create Account" page

**Scope**: Frontend — `/accounts/create/page.tsx`

Move the create account form from the existing combined page into a dedicated page at `/accounts/create`. The form remains functionally identical: account name, password, role fields with the same validation rules. On success, show a toast and reset the form.

### Atomic-4: Add backend DELETE account endpoint

**Scope**: Backend — `store.go`, `mysql.go`, `seed.go`, `controlplane.go`, `resources.go`

Add DELETE support to the `/api/v1/accounts/{id}` endpoint:
- `store.go`: Add `DeleteAccount(accountID string) error` to Store interface
- `mysql.go`: Implement `DeleteAccount` — DELETE FROM accounts WHERE id = ?
- `seed.go`: Implement `DeleteAccount` — no-op (seed store has no real persistence)
- `controlplane.go`: Add `DeleteAccount(accountID string) error` service method with admin protection:
  - Fetch account by ID, if `account == "admin"` → return error "cannot_delete_admin"
  - Otherwise delegate to store
- `resources.go`: Add `case http.MethodDelete` in `handleAccountByID`, with password rotation check

### Atomic-5: Create Edit Account Modal component

**Scope**: Frontend — `/accounts/_components/edit-account-dialog.tsx`

Create a reusable modal dialog component for editing an account. Fields:
- Password (optional, with hint "leave blank to keep current password")
- Role (input, pre-filled with current value)
- Status (select: active / disabled, pre-filled with current value)
- Cancel and Save buttons
- Form validation: password must be 8+ characters IF provided (optional field)

### Atomic-6: Add deleteAccount API client function

**Scope**: Frontend — `control-plane-api.ts`

Add `deleteAccount(accessToken, accountID)` function calling `DELETE /accounts/{id}`.

### Atomic-7: Create "Account List" page with query, edit, and delete

**Scope**: Frontend — `/accounts/list/page.tsx`

Move the accounts list from the existing combined page into a dedicated page at `/accounts/list`. Add:
- A table/list displaying all accounts (account name, role, status, ID)
- An "Edit" button on each account row → opens Edit Account Modal (Atomic-5)
- A "Delete" button on each account row (hidden for admin account)
- Delete confirmation dialog before sending delete request
- Delete uses `deleteAccount()` mutation → toast + list refresh on success
- Edit uses `updateAccount()` via modal → toast + list refresh on success

## User Stories

1. As an admin, I want to navigate between creating accounts and viewing accounts via sidebar sub-menus.
2. As an admin, I want a focused creation form without the distraction of the account list.
3. As an admin, I want to view all accounts and quickly edit an account's role, status, or password via a modal.
4. As an admin, I want to delete non-admin accounts with a confirmation step.
5. As an admin, I expect the system to prevent deletion of the built-in admin account.

## Acceptance Criteria

- Sidebar shows `创建账号` and `账号查询` as sub-menu items under `账号`
- `/accounts` redirects to `/accounts/create`
- `/accounts/create` shows the account creation form, works end-to-end
- `/accounts/list` shows all accounts in a list with Edit and Delete buttons
- Clicking "Edit" opens a modal with pre-filled role and status
- Changing password is optional in the edit modal
- Clicking "Delete" shows confirmation, then deletes the account
- The admin account (`account == "admin"`) has no Delete button in the list
- Backend `DELETE /api/v1/accounts/{id}` rejects deletion of admin account with error "cannot_delete_admin"
- Backend `DELETE /api/v1/accounts/{id}` enforces password rotation check
- All mutations show toast notifications on success/error
- Empty, loading, and error states are handled for all pages

## Dependencies

- Atomic-2 depends on Atomic-1 (redirect target path must exist)
- Atomic-7 depends on Atomic-5 and Atomic-6 (list page uses the edit modal and delete API)
- Atomic-4, Atomic-6 are backend/frontend counterparts — implement in that order
- Atomic-3 and Atomic-7 are independent of each other
