# Plan

1. Add focused application presence tests for online self-presence and offline after close/invalidation.
2. Add authenticated local alpha Protobuf flow proof for online, close/offline, and invalidation/offline using existing presence route and connection lifecycle primitives.
3. Add `ADR-0130` to record the hardening decision.
4. Add conversation memory preserving maintainer intent and the selected scope.
5. Update `.arch/work-items.yaml`:
   - mark `M-150/W-0222` completed;
   - open `M-151/W-0223` as next-ready.
6. Update runtime/reference/contracts/conventions/modules manifests and relevant guides so agents see `W-0223` as the next continuation step.
7. Register `runtime.presence_status_local_proof_hardening` in `rules/check-rules.json`.
8. Update `tools/vibit` expected-next helpers and runtime check logic.
9. Run verification commands and record the results.

## Files To Create

- `changes/2026-05-24-harden-presence-status-local-proof/request.md`
- `changes/2026-05-24-harden-presence-status-local-proof/spec.yaml`
- `changes/2026-05-24-harden-presence-status-local-proof/impact.md`
- `changes/2026-05-24-harden-presence-status-local-proof/plan.md`
- `changes/2026-05-24-harden-presence-status-local-proof/checklist.md`
- `changes/2026-05-24-harden-presence-status-local-proof/verification.md`
- `decisions/ADR-0130-presence-status-local-proof-hardening.md`
- `conversations/2026-05-24-presence-status-local-proof-hardening.md`

## Files To Edit

- `runtime/internal/app/presence/presence_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Generated Output

No generated output.

## Handwritten Logic

No production runtime handwritten behavior is changed. The only Go changes are tests.

`tools/vibit` receives check and next-ready helper updates only.

## Tests And Verification

Run:

```bash
go test ./internal/app/presence ./internal/platform/protocol/protobuf
go test ./...
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.presence_status_local_proof_hardening
node tools/vibit check change harden-presence-status-local-proof --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Notes

Rollback removes test additions, the hardening records, ADR-0130, rule registration, and W-0223 next-ready update. No data, protocol, generated output, dependency, or migration rollback is required.
