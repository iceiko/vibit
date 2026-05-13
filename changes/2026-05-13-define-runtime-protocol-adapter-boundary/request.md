# Request

## Original Request

Continue development.

## Clarified Requirement

Define the runtime protocol adapter boundary before Go runtime implementation starts, because the local toolchain cannot currently run Go, Buf, or Protobuf generation.

## User-Visible Outcome

Maintainers and agents can read one standard that explains how WebSocket transport, Protobuf protocol adaptation, application dispatch, generated code, and domain modules hand work to each other.

## Non-Goals

- Do not add Go runtime source files.
- Do not run `buf generate`.
- Do not commit generated Go Protobuf output.
- Do not implement WebSocket handlers.
- Do not implement inventory runtime business logic.

## Unknowns

- The local environment currently lacks Go, Buf, and `protoc`.
- The first concrete Go type names may change when implementation begins, but their responsibilities should preserve this boundary.

## Acceptance Criteria

- [ ] English and Simplified Chinese runtime protocol adapter standards exist.
- [ ] ADR-0018 records the boundary decision.
- [ ] Architecture manifests and agent guides reference the boundary.
- [ ] `node tools/vibit check runtime` verifies the boundary references.
- [ ] Verification records unavailable toolchain-dependent checks.
