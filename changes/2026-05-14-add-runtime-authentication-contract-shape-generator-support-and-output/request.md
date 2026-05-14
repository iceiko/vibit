# Request

## Original Request

The maintainer asked the agent to continue advancing the next planned work item.

## Clarified Requirement

Advance `W-0089` by adding generator, inspection, and generated-output drift-check support for the runtime authentication contract family. Generate metadata-only Go contract shape files from `contracts/runtime/authentication/**` under `runtime/internal/generated/contracts/runtime/authentication/`.

## User-Visible Outcome

Maintainers and agents can inspect generated authentication contract shape output through `node tools/vibit inspect generated --json`, regenerate it through `node tools/vibit generate contract-shapes all`, and verify it through `node tools/vibit check generated --json`.

The generated files expose authentication contract metadata only. They do not implement authentication behavior.

## Non-Goals

- Do not add application authentication service implementations.
- Do not add token generation, verifier comparison, login execution, token validation, logout execution, refresh, or cleanup jobs.
- Do not add Protobuf authentication messages.
- Do not add WebSocket proof carriers or routes.
- Do not add authentication dependencies.
- Do not change `authentication.Repository`.
- Do not change migration schemas.
- Do not hand-edit generated files.

## Unknowns

- The exact application service interface boundary remains for `W-0090`.
- Runtime authentication behavior remains for later gated work.
- Whether runtime session contract shapes should later use the same family-aware layout remains a future decision.

## Acceptance Criteria

- [x] Runtime authentication contract shape generation is supported.
- [x] Runtime authentication generated shape inspection is supported.
- [x] Runtime authentication generated-output checks cover missing, stale, and drifted files.
- [x] Metadata-only Go authentication contract shape files are generated from semantic sources.
- [x] Runtime boundary checks allow only source-traced metadata output and still forbid runtime authentication behavior.
- [x] `W-0089` is completed and a next ready work item is created.
