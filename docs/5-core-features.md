# 5 Core Features

## Admin Bootstrap

- node first start creates `admin`
- random password is shown once
- password is stored only as hash
- first login forces password rotation

## Node Management

- register node
- approve node enrollment
- assign parent node
- define node scope
- enable or disable node

## Route And Chain Management

- create whitelist route rules by domain, suffix, IP, CIDR, or protocol
- map a rule to a chain
- map a rule to a `destination_scope`
- validate chain loops before publish

## Policy Distribution

- publish a compiled policy revision
- nodes pull or receive latest policy revision
- nodes keep the last working revision locally

## Traffic Handling

- support `http`
- support `https`
- support `ws`
- support node-to-node CONNECT tunnels

## Certificates

- auto-renew public certificates for edge nodes
- manage private certificates for node trust
- expose expiry and renewal status in admin UI

## Runtime Status

- node health and heartbeat status
- current policy revision visibility
- certificate expiry visibility

## Chrome Extension

- login to selected node
- choose active profile
- show current active node
- show whether routing is direct or chained
