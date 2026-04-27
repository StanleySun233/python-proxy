# Master Plan: Frontend Enhancement - Remaining Tasks

## Team
| Role | ID | Plan | Scope |
|------|----|------|-------|
| Backend | BE-1 | [plan-be-1.md](./plan-be-1.md) | Route validation & suggestions APIs |
| Frontend | FE-1 | [plan-fe-1.md](./plan-fe-1.md) | Missing UI components & real-time validation |

## BE-1 — [plan](./plan-be-1.md)
- [x] Implement `POST /api/v1/route-rules/validate` — [ref](../routes-management-ui.md#post-apiv1route-rulesvalidate)
- [x] Implement `GET /api/v1/route-rules/suggestions` — [ref](../routes-management-ui.md#get-apiv1route-rulessuggestions)

## FE-1 — [plan](./plan-fe-1.md)
- [x] Add CompilationPreviewModal to chains editor — [ref](../chains-management-ui.md#fr5-compilation-preview)
- [x] Add real-time API validation to chain editor — [ref](../chains-management-ui.md#fr6-chain-validation-rules)
- [x] Add RegexTesterModal to routes editor — [ref](../routes-management-ui.md#fr3-match-type-and-value-intelligence)
- [x] Add real-time API validation to route editor — [ref](../routes-management-ui.md#fr7-rule-validation)

## Integration
- [x] End-to-end verify: route validate API + frontend real-time validation — BE-1 + FE-1
