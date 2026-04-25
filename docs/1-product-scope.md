# 1 Product Scope

## Goal

Build a multi-node proxy product where a user connects once from a Chrome extension and reaches internal targets through explicit chains such as `user -> a -> b -> c`.

## Primary Roles

- `user`: uses the Chrome extension to select a node and authenticate
- `edge node`: public entry node such as `a`
- `relay node`: private or downstream node such as `b c d`
- `control plane`: central Go backend for accounts, nodes, chains, and policy distribution
- `admin`: operator using the Next.js console

## V1 Scope

- Chrome extension stays plain JavaScript
- control plane backend uses Go and SQLite
- admin console uses Next.js
- node routing is whitelist-based
- supported traffic types are `http`, `https`, and `ws`
- public-facing certificates support automatic renewal
- node bootstrapping creates a temporary `admin` account with a one-time random password

## Key Product Objects

- `account`
- `role`
- `node`
- `node link`
- `chain`
- `route rule`
- `policy revision`
- `certificate`

## Explicit Non-Goals For V1

- automatic path finding
- distributed multi-primary database
- dynamic mesh routing
- fully adaptive per-request routing
- enterprise SSO
