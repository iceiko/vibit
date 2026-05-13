# Conversation: Application Dispatch Skeleton

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-add-application-dispatch-skeleton/`

Related artifacts:

- `runtime/internal/app/`
- `docs/runtime-protocol-adapter.md`
- `.arch/runtime.yaml`
- `tools/vibit`

## Context

The previous runtime slice introduced generated Go Protobuf output, application-owned handoff types, and a Protobuf protocol adapter that converts generated envelopes into `RouteRequest` values.

The maintainer asked to continue after that slice was committed and pushed.

## Maintainer Narrative

The maintainer continues to prefer deliberate progress toward a long-lived, agent-native server framework rather than a disposable demo. Earlier guidance also emphasized that branch points should be discussed, while professional sub-decisions inside ratified directions may be handled by the agent.

## Agent Response Summary

The agent chose the next narrow runtime step as application dispatch, not WebSocket transport or persistence.

This keeps the flow moving from protocol handoff toward domain behavior while preserving the layer boundaries already ratified by ADR-0018.

The agent intentionally avoided inventory business logic, PostgreSQL repositories, WebSocket transport, generated route registration, authentication, and transaction wiring in this change.

## Decisions

- Add a pure application-layer dispatcher under `runtime/internal/app/`.
- Keep route registration explicit until generated registration exists.
- Support command and query route requests first.
- Treat events, errors, system messages, acknowledgements, heartbeats, inputs, and state updates as non-dispatchable by the application dispatcher for now.
- Strengthen runtime checks so app/domain packages cannot import platform adapters or generated Protobuf packages.

## Artifacts

- Added application dispatch types and tests under `runtime/internal/app/`.
- Updated `.arch/runtime.yaml` to record the application dispatch skeleton state.
- Updated `docs/runtime-protocol-adapter.md` and its Simplified Chinese translation.
- Updated runtime and repository guides to distinguish current dispatch skeleton work from deferred transport, persistence, and business handler work.
- Added `runtime.layer_boundary` to the rule catalog and runtime checks.

## Open Questions

- Should generated route registration be created before the first inventory runtime handler?
- Should the first inventory handler use a test-only in-memory repository behind the future repository interface, or wait for PostgreSQL repository interfaces first?
- What exact transaction interface should command handlers receive once PostgreSQL persistence begins?

## Follow-Up

- Define the first inventory runtime handler boundary after repository and policy interfaces are declared.
- Decide whether generated route registration should come before handwritten inventory handler implementation.
- Keep WebSocket transport deferred until dispatch and module handler behavior are covered by tests.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
