# Frontend Large Page Component Refactor

## Overview

Decompose large page components (>300 lines) into hooks, sub-components, and utility modules. Target files: onboarding (1056 lines), node-pages (928 lines), routes (492 lines), health/overview (406 lines), health/heartbeat (391 lines).

## Current State

### onboarding/page.tsx (1056 lines)
Single file containing:
- 3 form type definitions (AccessPathFormValues, OnboardingTaskFormValues, PathEditorState, TaskEditorState)
- 12 useState declarations
- 8 useMutation hooks (create/update/delete for paths and tasks)
- 3 useQuery hooks
- 4 useMemo computed values
- Inline form rendering (~300 lines of JSX)
- Inline table rendering with edit/delete actions (~200 lines of JSX)
- Inline filter controls

### node-pages.tsx (928 lines)
Multiple page content functions in one file:
- NodeConnectPageContent
- NodeManualPageContent
- NodeBootstrapPageContent
- NodeApprovePageContent
- NodeTopologyPageContent
- NodeRegistryPageContent (likely large)

### routes/page.tsx (492 lines)
- Inline validation logic (validateMatchValue)
- Form state, mutations, debounced validation
- Rule registry table with inline editing

### health/overview/page.tsx (406 lines)
- Charts, health rows computation, detail panel

### health/heartbeat/page.tsx (391 lines)
- Time window selection, charts, node detail

## Atomic Requirements

### Atomic 4: Extract onboarding hooks
Create `app/[locale]/(console)/onboarding/_hooks/use-onboarding.ts`:
- Extract all useState, useQuery, useMutation, useMemo from onboarding/page.tsx
- Export a single `useOnboarding()` hook returning all state + actions

Create `_hooks/use-onboarding-paths.ts` for path-specific logic if the combined hook exceeds 200 lines.

**Dependencies**: Atomic 1, Atomic 2 (new import paths)

### Atomic 5: Extract onboarding sub-components
Create files under `_components/`:
- `onboarding-path-form.tsx` — create path form
- `onboarding-task-form.tsx` — create task form
- `onboarding-path-table.tsx` — path registry table with edit/delete
- `onboarding-task-table.tsx` — task registry table with status editing
- `onboarding-metrics.tsx` — metrics grid cards
- `onboarding-path-editor.tsx` — inline edit form for paths
- `onboarding-task-editor.tsx` — inline edit form for tasks

Keep `page.tsx` as thin composition (<100 lines).

**Dependencies**: Atomic 4

### Atomic 6: Split node-pages.tsx into individual files
Each exported component gets its own file:
- `_components/connect-page-content.tsx` (currently ~35 lines)
- `_components/manual-page-content.tsx` (currently ~30 lines)
- `_components/bootstrap-page-content.tsx` (currently ~30 lines)
- `_components/approvals-page-content.tsx`
- `_components/topology-page-content.tsx`
- `_components/registry-page-content.tsx`

Remove `node-pages.tsx` after extraction. Update page imports.

**Dependencies**: None (pure component extraction, imports stay within same directory)

### Atomic 7: Extract routes hooks and validation
- `routes/_lib/validation.ts` — validateMatchValue() function
- `routes/_hooks/use-routes.ts` — query, mutation, debounce logic
- `routes/_components/route-rule-form.tsx` — create rule form
- `routes/_components/route-rule-table.tsx` — rule registry with inline editing

Keep `page.tsx` <100 lines.

**Dependencies**: Atomic 1, Atomic 2

### Atomic 8: Extract health overview hooks and sub-components
- `health/overview/_hooks/use-health-overview.ts`
- `health/overview/_components/health-metrics.tsx`
- `health/overview/_components/health-table.tsx`
- `health/overview/_components/health-detail-panel.tsx`

**Dependencies**: Atomic 1, Atomic 2

### Atomic 9: Extract health heartbeat hooks and sub-components
- `health/heartbeat/_hooks/use-heartbeat.ts`
- `health/heartbeat/_components/heartbeat-chart.tsx`

**Dependencies**: Atomic 1, Atomic 2

## Acceptance Criteria

- No page component file exceeds 150 lines
- Each extracted hook returns a clearly typed interface
- Each sub-component has a single responsibility
- All existing functionality preserved
- TypeScript compilation passes

## User Stories

- As a developer, I can understand a page's data flow by reading a single hook file
- As a developer, I can modify a form component without touching table logic
- As a developer, I can add a new tab/node operation by creating one small file
