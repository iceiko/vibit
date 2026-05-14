# Request

Define the future application-owned authentication service interface boundary before adding runtime authentication behavior.

## Source

The maintainer said to continue advancing the work queue. The current next-ready work item is:

```text
W-0090 Define application authentication service interface boundary
```

## Scope

This change may:

- Define the future application-owned authentication service interface boundary under `runtime/internal/app`.
- Map generated authentication contract shapes to service-level request/result vocabulary.
- Record redaction expectations for credential and token material.
- Define error, permission, audit, and request identity handoff expectations.
- Declare how the future service boundary may use application unit of work and `authentication.Repository`.
- Update manifests, guides, checks, memory records, and work queue state.

This change must not:

- Add Go application authentication service code.
- Add token generation.
- Add verifier comparison.
- Add login execution.
- Add token validation.
- Add logout execution.
- Add cleanup jobs.
- Add Protobuf authentication messages.
- Add WebSocket proof carriers.
- Add authentication dependencies.
- Change `authentication.Repository`.
- Change migration schemas.
- Add production authentication behavior.
