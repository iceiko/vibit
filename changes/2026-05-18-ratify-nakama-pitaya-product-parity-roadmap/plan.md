# Plan

## Files To Create

- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- `decisions/ADR-0078-nakama-pitaya-product-parity-roadmap.md`
- `conversations/2026-05-18-nakama-pitaya-product-parity-roadmap.md`
- `changes/2026-05-18-ratify-nakama-pitaya-product-parity-roadmap/*`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/reference.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`
- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-implementation/*`

## Generated Artifacts

None.

## Handwritten Logic

Only repository check logic in `tools/vibit`; no runtime Go behavior.

## Tests

No Go tests. Use repository checks and JavaScript syntax check.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change ratify-nakama-pitaya-product-parity-roadmap --json`
- `node tools/vibit check all --json`
- `git diff --check`
- `node tools/vibit inspect next --json`

## Rollback Or Migration Notes

This is a roadmap-only change. Reversal would remove ADR-0078, the roadmap standard, and the corresponding work/check markers.
