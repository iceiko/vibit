# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0009` by adding a repeatable request-loop test fixture for WebSocket Protobuf command/query tests. The fixture should reduce duplicated test-only inventory dispatcher setup without widening runtime dependencies or changing production behavior.

## User-Visible Outcome

Agents extending request-loop coverage can reuse a package-local fixture for:

- In-memory inventory dispatcher setup.
- Protobuf envelope construction and frame payload marshaling.
- Protobuf response envelope unmarshaling.
- Envelope-to-application request routing helpers.

## Non-Goals

- Do not add a new public runtime API.
- Do not add a new test framework dependency.
- Do not import WebSocket packages outside the transport adapter package.
- Do not add authentication, session validation, PostgreSQL persistence, or MinIO.
- Do not change the Protobuf envelope shape or generated Protobuf files.

## Unknowns

- None for this bounded step.

## Acceptance Criteria

- [x] Add a reusable request-loop test fixture under the Protobuf protocol adapter package.
- [x] Replace duplicated test-only inventory repository, permission, event id, and clock helpers.
- [x] Keep the fixture test-only and package-local.
- [x] Preserve transport, protocol, application, domain, and generated output boundaries.
- [x] Run Go and repository verification.
