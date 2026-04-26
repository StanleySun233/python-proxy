# 17 Frontend Information Architecture

## Goal

Design the admin console around operator workflows instead of mirroring backend tables one-to-one.

## Top-Level Navigation

- `Overview`
- `Nodes`
- `Onboarding`
- `Chains`
- `Routes`
- `Health`
- `Accounts`
- `Certificates`

## Primary Operator Flows

### Overview

Purpose:

- summarize node state
- summarize onboarding task state
- summarize policy and certificate posture
- expose the current topology shape immediately

Required surfaces:

- metric cards
- relay and onboarding topology preview
- pending task queue
- quick status bands

### Nodes

Purpose:

- inspect node identity and health
- inspect parent-child relationships
- inspect public host and relay positioning

Required surfaces:

- table view
- topology view
- detail drawer or split panel

### Onboarding

Purpose:

- manage `direct`, `relay_chain`, and `upstream_pull`
- create and inspect `node_access_paths`
- create and inspect `node_onboarding_tasks`

Required surfaces:

- mode switcher
- access-path list
- onboarding task list
- path preview canvas
- create-task form

### Chains

Purpose:

- edit relay chains and ordered hops
- validate destination scope and hop ownership

Required surfaces:

- ordered chain list
- drag-and-drop hop editor
- compile summary panel

### Routes

Purpose:

- manage whitelist forwarding rules
- map matches to destination scope and execution path

Required surfaces:

- rule table
- create and edit form
- overlap and conflict hints

### Health

Purpose:

- centralize node heartbeat, task connectivity, and certificate renewal state

Required surfaces:

- heartbeat panels
- health timeline or charts
- renewal warnings

### Accounts

Purpose:

- manage operator identities and role boundaries

### Certificates

Purpose:

- track public and internal certificate state
- expose renewal posture without introducing audit-heavy complexity

## Layout Rules

- persistent left navigation on desktop
- content-first stacked layout on smaller viewports
- graph and table views should coexist for topology-heavy areas
- raw IDs, ports, and hop paths should use mono typography

## State Rules

- server truth should come from query state, not duplicated local stores
- page filters, selected tabs, and detail context should be URL-friendly
- task and health views should support polling from the start

## V1 Frontend Delivery Order

1. app shell, theme, locale, query provider
2. overview page
3. onboarding page
4. nodes page
5. chains and routes pages
6. health, certificates, accounts
