# Request

## Original Request

继续推进十步，除非需要遇到我确认的问题，否则不要停下来。

## Clarified Requirement

Advance `W-0086` by defining the runtime authentication implementation boundary after the authentication PostgreSQL adapter milestone, without implementing runtime authentication behavior.

## User-Visible Outcome

Maintainers and agents can see how future runtime authentication work must be split before code appears:

- Application service ownership.
- Repository usage through `authentication.Repository`.
- PostgreSQL adapter persistence-only role.
- Token material lifecycle placeholder.
- Verifier comparison boundary.
- Request identity handoff.
- Error, permission, and audit mapping.
- Separate gates for token generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, generated authentication shapes, Protobuf messages, WebSocket proof carriers, and dependency adoption.

## Non-Goals

- Do not implement login.
- Do not generate tokens.
- Do not validate tokens.
- Do not compare credential or token verifiers.
- Do not execute logout.
- Do not add refresh behavior.
- Do not add cleanup jobs.
- Do not add Protobuf authentication messages.
- Do not add WebSocket proof carriers.
- Do not add generated authentication shapes.
- Do not add authentication dependencies.
- Do not change `authentication.Repository`.
- Do not change ratified migration schemas.

## Unknowns

- Exact verifier algorithms remain future work.
- Exact generated authentication shape timing remains future work.
- Exact Protobuf authentication payload shape remains future work.
- Exact WebSocket proof-carrier posture remains deferred.

## Acceptance Criteria

- [x] Add a runtime authentication implementation boundary standard and Simplified Chinese translation.
- [x] Add an ADR for the boundary.
- [x] Update manifests and guides with the new boundary.
- [x] Mark `W-0086` completed and create a conservative next work item.
- [x] Preserve all runtime authentication behavior deferrals.
