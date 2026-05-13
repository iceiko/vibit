# Request

## Original Request

Continue development.

## Clarified Requirement

Add a generated-output standard and repository check before generated Go Protobuf output is committed or runtime protocol code is added.

## User-Visible Outcome

Maintainers and agents can see where generated output belongs, what source traces it must carry, and how to verify that generated Go Protobuf files were produced from real `.proto` sources instead of hand-written under the generated tree.

## Non-Goals

- Do not run `buf generate`.
- Do not commit generated Go Protobuf output.
- Do not add WebSocket runtime handlers.
- Do not implement inventory runtime business logic.

## Unknowns

- The local environment may still lack Buf, Protobuf compiler tooling, and Go.
- Generated contract output rules will need deeper checks when the contract generator exists.

## Acceptance Criteria

- [ ] English and Simplified Chinese generated-output standards exist.
- [ ] An ADR records the generated output decision.
- [ ] Architecture manifests and agent guides reference the standard.
- [ ] `node tools/vibit check generated` verifies `runtime/internal/generated/proto/`.
- [ ] Repository verification records what was and was not run.
