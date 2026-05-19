# ADR-0066: Runtime Session Validation Implementation

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-runtime-session-validation-gate/`
- `changes/2026-05-17-implement-runtime-session-validation/`

Related conversations:

- `conversations/2026-05-17-runtime-session-validation-implementation.md`

Related artifacts:

- `runtime/internal/app/runtime_session_validator.go`
- `runtime/internal/app/runtime_session_validator_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0142` defined the runtime session validation gate after durable `runtime_sessions`, the storage-neutral session repository interface, and the PostgreSQL adapter were in place.

The work queue reached `M-071/W-0143`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key reference baselines.

## Decision

Select:

```text
implement_runtime_session_validation
```

Implement the bounded persistent session validator in:

```text
runtime/internal/app/runtime_session_validator.go
runtime/internal/app/runtime_session_validator_test.go
```

The validator is application-owned and depends only on:

```text
runtime/internal/app/session.Repository
```

It uses:

```text
FindActiveSessionByID
```

The validator requires an already validated player identity before repository lookup, validates an active and unexpired persisted session row, requires actor and player identity match, and returns `RequestIdentity.SessionValidated = true` only after successful durable validation.

This ADR does not create sessions at login or `BindConnection`, update `last_seen_at`, change route policy, change WebSocket handshake authentication, add transport credential carriers, add Protobuf session messages, change the existing Protobuf envelope, add generated output, add logout/revocation active-connection invalidation, add reconnect/epoch behavior, add cleanup jobs, add dependencies, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Wire persisted session validation directly into protected route policy.
- Create durable runtime sessions during login or first-message binding in the same slice.
- Update `last_seen_at` on every validation.
- Treat a client-supplied `session_id` as proof without an already validated actor identity.
- Put validation into the PostgreSQL adapter, WebSocket transport, or Protobuf adapter.
- Add logout/revocation active-connection invalidation together with validation.
- Copy Nakama or Pitaya session APIs directly.

## Rationale

Nakama shows that session validity is a lifecycle concern with expiration, refresh, logout, and active socket implications. It also shows that ended or expired session state does not automatically mean every surrounding play or sign-in layer ends at the same time.

Pitaya shows that user session context, acceptors, handler services, agents, pipelines, serializers, and route handlers are separate responsibilities. The useful vibit adaptation is an application-owned validator that can build request identity without moving durable validation into transport or protocol code.

The bounded implementation gives later route-policy and session-creation work a concrete validator to compose while keeping those decisions separate.

## Agent Reasoning Summary

The agent selected implementation because the repository already had the prerequisites: `runtime_sessions`, `session.Repository`, the PostgreSQL adapter, and a ratified validation gate. The implementation is deliberately narrow: lookup-only, fake-repository tested, application-owned, and not wired into startup or routing.

This keeps the design aligned with Nakama's mature session lifecycle pressure and Pitaya's transport/session/handler separation without adopting their public API shapes.

## Decision Weights

```yaml
decision_weights:
  validated_identity_safety: high
  application_layer_ownership: high
  nakama_pitaya_alignment: high
  route_policy_deferral: high
  session_creation_deferral: high
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `PersistentSessionValidator` implements `SessionValidator`.
- Validation uses `session.Repository.FindActiveSessionByID`.
- Successful validation can produce `RequestIdentity.SessionValidated = true`.
- Public invalid-session failures collapse to a stable redacted reason.
- Focused tests cover success, missing or malformed session id, metadata-only identity, unvalidated identity, lookup/repository failures, record mismatches, nil repository, and no mutation.
- `runtime.runtime_session_validation_implementation` becomes the repository check rule for this slice.
- The work queue blocks again after implementation at `M-073/W-0145`.

## Reversal Conditions

Revisit this decision if a future ADR chooses a different durable session store, changes the session repository contract, requires handshake-level authentication to own session proof, requires every validation to update `last_seen_at`, or adopts direct Nakama/Pitaya compatibility.

## Follow-Up

- Define session creation composition before login or `BindConnection` creates durable runtime sessions.
- Define bound/session identity route policy before protected routes can rely on session-validated identity.
- Define logout/revocation active-connection behavior before revocation closes or invalidates WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
