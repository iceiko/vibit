# Conversation: Pitaya-Aligned Acceptor Connection Lifecycle Boundary Gate

Date: 2026-06-01

## Context

The maintainer asked to continue moving toward Pitaya. `W-0264` selected `define_pitaya_aligned_acceptor_connection_lifecycle_boundary_gate` as the next bounded direction and opened `W-0265` as next-ready.

## Maintainer Narrative

Continue advancing the repository toward Pitaya-class architecture vocabulary while preserving source-first, agent-native maintainability and committing/pushing each completed slice.

## Agent Response Summary

Defined a gate-only Pitaya-aligned acceptor and connection lifecycle boundary. The gate maps existing single-process WebSocket acceptor, connection id/epoch metadata, first-message binding, active connection registry, close handoff, and presence lifecycle surfaces without adding behavior.

## Decisions

- Accepted `ADR-0173`.
- Registered `runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate`.
- Completed `M-193/W-0265`.
- Opened `M-194/W-0266 Implement Pitaya-aligned acceptor and connection lifecycle source-first map` as next-ready.

## Artifacts

- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.zh-CN.md`
- `decisions/ADR-0173-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

None for this gate.

## Follow-Up

- Complete `W-0266`: implement the source-first Pitaya-aligned acceptor and connection lifecycle map.
- Keep acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, metrics/tracing behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility deferred.

## Redaction Notes

No ignored credential file contents were read or recorded. No secrets, raw token values, verifier keys, database credentials, transport payloads, or private environment values are included.
