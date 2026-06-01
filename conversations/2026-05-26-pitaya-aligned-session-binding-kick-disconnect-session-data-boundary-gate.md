# Conversation: Pitaya-Aligned Session Binding, Kick/Disconnect, And Session Data Boundary Gate

Date: 2026-06-01
Work item: W-0268
Decision: ADR-0176
Check rule: runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate

## Context

`W-0268` followed `ADR-0175`, which selected the session binding, kick/disconnect, and session data boundary gate after the acceptor and connection lifecycle source-first map.

## Maintainer Narrative

The maintainer asked to continue advancing toward Pitaya alignment in bounded steps, with commits and pushes after completed increments.

## Agent Response Summary

- Continued from `W-0268 Define Pitaya-aligned session binding, kick/disconnect, and session data boundary gate`.
- Used the existing acceptor and connection lifecycle gate pattern from `W-0265`.
- Defined a gate-only standard for session binding, kick/disconnect, and session data vocabulary.
- Registered `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate`.
- Completed `M-196/W-0268`.
- Opened `M-197/W-0269 Implement Pitaya-aligned session binding, kick/disconnect, and session data source-first map` as next-ready.

## Decisions

- Accepted `ADR-0176`.
- Registered `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate`.
- Kept the follow-up implementation source-first only through `W-0269`.

## Artifacts

- `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md`
- `docs/pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.zh-CN.md`
- `decisions/ADR-0176-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for W-0268. W-0269 remains the next-ready source-first map slice.

## Follow-Up

- Implement `W-0269` as a source-first repository inspection map.
- Keep session binding behavior, kick/disconnect behavior, session data persistence, protocol, generated output, persistence, dependencies, hosted, SDK, distributed runtime, and direct compatibility behind later bounded work items.

## Redaction Notes

No ignored credential file contents were read or printed.

This slice did not add session binding behavior, kick/disconnect behavior, session data behavior, session data persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, handler routing behavior, handler pipeline behavior, backend route targeting, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, metrics endpoints, tracing pipelines, service discovery implementation, RPC implementation, remote calls, frontend/backend server role implementation, distributed runtime behavior, cluster-safe session routing behavior, hosted deployment, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

## Verification Record

Expected RED checks before implementation:

- `node tools/vibit check change define-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate --json` failed because the change directory did not exist.
- `node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate` failed because the rule id was not registered.

Expected verification after implementation:

- `node -c tools/vibit`
- `node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate`
- `node tools/vibit inspect next --json`
- `node tools/vibit check change define-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`
