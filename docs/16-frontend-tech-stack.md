# 16 Frontend Tech Stack

## Goal

Establish a frontend stack that matches the control-plane product shape instead of optimizing for a temporary landing page.

## Core Choices

### Framework

- `Next.js 14`
- App Router
- React Server Components for route-level shells and metadata

Reason:

- matches the current repository direction
- supports route-level composition, locale segments, and future server-side data shaping

### Internationalization

- `next-intl`

Reason:

- clean App Router integration
- route-aware locale handling
- server and client translation support without forcing a custom i18n layer

### Theme System

- `next-themes`
- CSS variables as the primary token layer

Reason:

- theme switching remains simple
- `灵动` light theme and matching dark theme can share the same token structure
- avoids locking the project into a component library theme model

### Server State

- `@tanstack/react-query`

Reason:

- the control plane is server-state-heavy
- nodes, onboarding tasks, health, policies, and certificates all need cache, refetch, and mutation support

### Forms And Validation

- `react-hook-form`
- `zod`

Reason:

- onboarding, node editing, path editing, rule editing, and account management all require structured forms
- payload validation should match API constraints early in the frontend

### Visual Topology And Path Editing

- `reactflow`
- `dnd-kit`

Reason:

- topology and relay-chain visibility are first-class product features
- hop ordering and path editing need reliable drag-and-drop support

### Base UI Utilities

- `lucide-react`
- `sonner`
- `recharts`

Reason:

- iconography, operator feedback, and health or certificate charts are baseline console needs

## Style Direction

The frontend must inherit the theme family defined in [style.css](/home/sijin/workspace/proxy/docs/style.css:1).

### Light Theme

- warm paper background
- moss green primary
- camel and clay accents
- deep brown text

### Dark Theme

- dark olive-charcoal base
- preserved moss green highlight
- copper-clay accents
- warm parchment text

The dark theme should remain in the same visual family instead of switching to generic neon blue dashboard colors.

## Runtime Libraries To Add

- `next-intl`
- `next-themes`
- `@tanstack/react-query`
- `react-hook-form`
- `zod`
- `reactflow`
- `@dnd-kit/core`
- `@dnd-kit/sortable`
- `@dnd-kit/utilities`
- `lucide-react`
- `sonner`
- `recharts`

## Deferred Additions

- generated OpenAPI client from `apps/one-panel-api/openapi.yaml`
- richer chart package only if `recharts` becomes limiting
- state store such as `zustand` only if local UI state complexity truly grows beyond route and query state
