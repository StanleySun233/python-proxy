# Master Plan: Frontend Enhancement - Remaining Tasks

## Team
| Role | ID | Plan | Scope |
|------|----|------|-------|
| Backend | BE-1 | [plan-be-1.md](./plan-be-1.md) | Route validation & suggestions APIs |
| Frontend | FE-1 | [plan-fe-1.md](./plan-fe-1.md) | Missing UI components & real-time validation |

## BE-1 — [plan](./plan-be-1.md)
- [ ] Implement `POST /api/v1/route-rules/validate` — [ref](../routes-management-ui.md#post-apiv1route-rulesvalidate)
- [ ] Implement `GET /api/v1/route-rules/suggestions` — [ref](../routes-management-ui.md#get-apiv1route-rulessuggestions)

## FE-1 — [plan](./plan-fe-1.md)
- [ ] Add CompilationPreviewModal to chains editor — [ref](../chains-management-ui.md#fr5-compilation-preview)
- [ ] Add real-time API validation to chain editor — [ref](../chains-management-ui.md#fr6-chain-validation-rules)
- [ ] Add RegexTesterModal to routes editor — [ref](../routes-management-ui.md#fr3-match-type-and-value-intelligence)
- [ ] Add real-time API validation to route editor — [ref](../routes-management-ui.md#fr7-rule-validation)

## Integration
- [ ] End-to-end verify: route validate API + frontend real-time validation — BE-1 + FE-1
