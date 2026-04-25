# 9 Chrome Extension Client

## Responsibilities

- authenticate user against the selected entry node or control plane
- store active session and selected profile
- configure browser proxy mode
- display current node, profile, and routing state
- expose quick actions for current site routing overrides

## V1 Screens

- popup: quick status, active node, active profile, routing mode
- options: node list, profile list, login state, local overrides

## Data Model

- `selectedNodeId`
- `selectedProfileId`
- `accessToken`
- `refreshToken`
- `routingMode`
- `themeMode`
- `siteOverrides`

## API Dependencies

- login
- refresh token
- list profiles available to user
- fetch current policy metadata
- optional fetch node health summary

## UX Rules

- user should always know which node is active
- user should always know whether traffic is direct, chained, or partially bypassed
- token expiry should fail visibly and require relogin
- site overrides must clearly indicate local-only behavior versus centrally managed policy
