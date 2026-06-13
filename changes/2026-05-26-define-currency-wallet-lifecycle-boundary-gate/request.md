# Request

## Original Request

Continue repository work from the current `next_ready` item for up to 20 steps.

## Interpreted Work Item

Advance `W-0292 Define currency wallet lifecycle boundary gate`.

## User-Visible Outcome

The repository records a semantic-only currency wallet lifecycle boundary gate and opens the next bounded work item:

```text
W-0293 Define currency wallet persistence schema gate
```

## Scope

This is a boundary-gate change. It defines future currency wallet lifecycle vocabulary, Nakama/Hiro economy capability mapping, identity and permission posture, invariants, redaction expectations, future tests, and deferrals.

## Non-Goals

- No currency wallet behavior.
- No balance tables.
- No wallet transaction behavior.
- No reward integration.
- No inventory integration.
- No purchase behavior.
- No grant execution.
- No spend execution.
- No transfer behavior.
- No reservation or settlement behavior.
- No audit/event tables.
- No runtime behavior changes.
- No protocol messages or routes.
- No Protobuf source.
- No generated output.
- No repository interface changes.
- No PostgreSQL adapters.
- No migrations.
- No dependencies.
- No startup wiring.
- No authentication/session behavior changes.
- No hosted deployment.
- No SDK publication.
- No release artifacts.
- No distributed runtime behavior.
- No direct Nakama/Pitaya API compatibility.
