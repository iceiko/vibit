# Request

## Original Request

The maintainer asked the agent to continue the next planned work item.

## Clarified Requirement

Advance `W-0088` by deciding whether generated Go authentication contract shapes should be introduced before application service interfaces and runtime authentication behavior. Record the source, output path, immutable status, check requirements, and relationship to semantic authentication contracts.

## User-Visible Outcome

Maintainers and agents can see that generated authentication contract shapes are recommended before service interfaces, but no generated files or runtime authentication behavior are added by this change.

## Non-Goals

- Do not generate authentication contract shape files.
- Do not add or change application authentication service code.
- Do not add token generation, verifier comparison, login execution, token validation, logout execution, refresh, or cleanup jobs.
- Do not add Protobuf authentication messages.
- Do not add WebSocket proof carriers or routes.
- Do not add authentication dependencies.
- Do not change `authentication.Repository`.
- Do not change migration schemas.

## Unknowns

- The exact generator implementation details remain for the next bounded work item.
- Whether runtime session contract shapes should later use the same family-aware layout is not decided here.

## Acceptance Criteria

- [x] Timing decision is recorded.
- [x] Source and output path are recorded.
- [x] Generated file immutability is recorded.
- [x] Check requirements are recorded.
- [x] Actual generated files remain deferred.

