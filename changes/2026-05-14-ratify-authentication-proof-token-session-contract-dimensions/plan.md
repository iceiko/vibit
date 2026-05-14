# Plan

1. Add `docs/authentication-proof-token-session-contract-dimensions.md`.
2. Add `docs/authentication-proof-token-session-contract-dimensions.zh-CN.md`.
3. Update the existing authentication/token/session validation standard to point to the dimensions standard.
4. Update runtime session semantic contract sources with dimension references and failure-class metadata.
5. Update `.arch/contracts.yaml`, `.arch/runtime.yaml`, `.arch/conventions.yaml`, `.arch/reference.yaml`, `AGENTS.md`, `AGENTS.zh-CN.md`, `runtime/AGENTS.md`, and `runtime/AGENTS.zh-CN.md`.
6. Mark `W-0059` completed and `W-0060` next-ready in `.arch/work-items.yaml`.
7. Run verification.
8. Record verification results.

## Generated Artifacts

None.

## Runtime Code

No runtime Go code should change in this step.

## Rollback Notes

Because this is a design and semantic contract clarification, rollback would remove the new standard and restore the previous manifest/work-item references. No data migration rollback is needed.
