# ADR-0083: Reconnect Connection Epoch Functional Slice

Status: Accepted
Date: 2026-05-20
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-20-define-reconnect-connection-epoch-functional-slice/`

Related conversations:

- `conversations/2026-05-20-reconnect-connection-epoch-functional-slice.md`

Related artifacts:

- `runtime/internal/app/connection/registry.go`
- `runtime/internal/app/connection/registry_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `rules/check-rules.json`

## Context

`ADR-0082` changed the workflow so non-security functional work can proceed as a single bounded slice with its boundary embedded in the change spec. `M-103/W-0175` is the first application of that strategy.

Before this slice, vibit had:

- first-message connection binding,
- a single-process active connection registry,
- a close policy that can invalidate active registry records,
- a concrete WebSocket transport close handoff by server-observed connection id and epoch.

The missing smallest reconnect/epoch primitive was an application-level rule for what happens when the same server-owned connection id is observed with a newer epoch. Without that rule, stale close intents are protected at transport, but the application registry still lacks an explicit lifecycle state for older active epochs.

## Decision

Select:

```text
define_reconnect_connection_epoch_functional_slice
```

as a Tier 2 functional slice and implement the smallest checkable behavior directly.

The application-owned connection registry now:

- adds `StateSuperseded`,
- adds `ErrorCodeConnectionEpochStale`,
- adds `SupersededAt` and `SupersededByEpoch` to lifecycle records,
- accepts a newer epoch for the same server-owned connection id and marks earlier active records superseded,
- rejects stale or repeated epochs after a newer epoch has been observed,
- keeps superseded records inspectable through `FindConnectionByID`,
- excludes superseded records from active target lists.

The implementation is in:

```text
runtime/internal/app/connection/registry.go
runtime/internal/app/connection/registry_test.go
```

## Boundaries

This ADR keeps these boundaries:

- Connection epoch progression is application registry behavior under `runtime/internal/app/connection`.
- WebSocket transport still owns only accepted socket metadata and concrete close handoff mechanics.
- Protocol adapters do not own reconnect or resume behavior.
- Authentication and logout remain token lifecycle behavior only.
- Domain modules do not read or mutate connection lifecycle state directly.

This decision does not add reconnect tokens, resume tokens, durable/distributed reconnect state, duplicate replacement socket close, close code mapping, close reason text, logout-triggered socket close, runtime session revocation, protocol session carriers, presence lifecycle, operations/admin disconnect, dependencies, or direct Nakama/Pitaya API compatibility.

## Nakama And Pitaya Mapping

Nakama informs the product pressure: realtime connection lifecycle must be explicit before presence, session visibility, and match runtime depend on it.

Pitaya informs the architecture pressure: acceptors, sessions, route handlers, and connection management remain separate surfaces.

vibit adapts those lessons by adding only server-observed epoch progression to its application registry. This does not copy either project's public API.

## Alternatives Considered

- Create another pure reconnect/epoch gate before implementation.
- Add reconnect or resume protocol messages now.
- Let WebSocket transport own reconnect behavior.
- Close older sockets automatically when a newer epoch is observed.
- Reuse the existing `closed` or `invalidated` states instead of adding `superseded`.
- Allow same-epoch reopen after terminal records.
- Jump directly to presence or protocol session carriers.

## Rationale

The useful middle ground is to make epoch ordering explicit without selecting product-facing reconnect behavior yet. `superseded` is a distinct lifecycle fact: the record is not merely closed by transport and not invalidated by policy. It became stale because a newer server-observed epoch exists for the same connection id.

Rejecting stale epochs after a newer epoch has been observed prevents regression in application state and gives future close, presence, and session-carrier work a stable monotonic primitive.

Keeping this as a Tier 2 slice follows the maintainer's velocity requirement while preserving control through a change spec, tests, manifests, and repository checks.

## Agent Reasoning Summary

The next step toward Nakama/Pitaya-class lifecycle behavior should not be a broad reconnect feature. The smallest safe step is server-observed epoch progression in the registry because later presence, match runtime, and protocol session carriers all need to know which connection record is current.

## Decision Weights

```yaml
decision_weights:
  development_velocity: high
  lifecycle_closure: high
  stale_connection_safety: high
  transport_application_separation: high
  protocol_surface_stability: high
  direct_api_compatibility: low
  reconnect_resume_now: low
confidence: high
```

## Consequences

- `runtime.reconnect_connection_epoch_functional_slice` becomes the repository check rule for this slice.
- The registry has a tested terminal `superseded` state for older active epochs.
- Stale epoch registration fails closed with `connection_epoch_stale`.
- Future protocol session carrier, presence, duplicate replacement, and reconnect/resume work can build on a clear epoch ordering rule.
- The work queue advances to protocol session carrier lifecycle work before presence expansion.

## Reversal Conditions

Revisit this decision if a future ADR replaces the server-observed connection id/epoch model, if distributed runtime requires a different monotonic epoch source, or if direct Nakama/Pitaya API compatibility is explicitly selected and requires a different lifecycle vocabulary.

## Follow-Up

- Define the protocol session carrier functional slice after epoch behavior is stable.
- Keep reconnect/resume tokens, durable reconnect, duplicate replacement close, and close code/reason mapping behind explicit future work.
