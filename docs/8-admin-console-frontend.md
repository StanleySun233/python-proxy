# 8 Admin Console Frontend

## Visual Direction

Use the same color family as the Chrome extension:

- light theme uses warm paper, moss green, sand, and clay accents
- dark theme uses navy, cyan-mint, and electric blue accents
- typography stays in the same pragmatic family: `IBM Plex Sans` and `IBM Plex Mono`

## Main Views

### Dashboard

- top summary for healthy nodes, degraded nodes, active chains, and latest policy revision
- recent publish events
- certificate expiry panel

### Nodes

- node grid with status, scope, parent, and trust state
- node detail drawer with listeners, links, heartbeats, and latest assigned policy

### Chains

- visual hop editor
- chain validation summary
- destination scope visibility

### Route Rules

- rule table sorted by priority
- create/edit panel with match type, match value, action, chain, and destination scope

### Health And Certificates

- node heartbeat timeline
- certificate expiry and renewal status

## Information Architecture

- `Overview`
- `Nodes`
- `Chains`
- `Route Rules`
- `Policies`
- `Accounts`
- `Certificates`
- `Health`

## Interaction Rules

- destructive actions require confirmation
- publish action shows compile summary before confirm
- node and chain pages must expose where traffic will terminate
- route rule editing must surface overlapping scope conflicts

## Responsive Behavior

- desktop first for operators
- tablet keeps left navigation and collapses secondary panes
- mobile only needs read-only health and emergency actions in V1
