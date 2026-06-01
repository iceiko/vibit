# Conversation: Pitaya-Aligned Session Binding Kick Disconnect Session Data Source-First Map

Date: 2026-06-01

## Context

The maintainer asked to continue advancing toward Pitaya alignment and to preserve commit/push discipline. The next-ready work item was `W-0269 Implement Pitaya-aligned session binding, kick/disconnect, and session data source-first map`, opened by `ADR-0176` after `W-0268` defined the Pitaya-aligned session binding, kick/disconnect, and session data boundary gate.

The initial RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect pitaya-session-lifecycle --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map

node tools/vibit check change implement-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map --json
# change directory does not exist
```

## Maintainer Narrative

Continue advancing in bounded steps toward Pitaya-class capability coverage while keeping direct compatibility, credentials, and production behavior changes out of this slice.

## Agent Response Summary

Accepted `ADR-0177` and implemented a source-first inspection map through:

- `node tools/vibit inspect pitaya-session-lifecycle --json`
- `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map`
- `changes/2026-05-26-implement-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map/`

The map reports allowed session binding, kick/disconnect, and session data vocabulary, current single-process mappings, redaction posture, and explicit deferrals without adding runtime behavior.

## Decisions

- Complete W-0269 with a source-first repository inspection map.
- Register `runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_source_first_map`.
- Expose `node tools/vibit inspect pitaya-session-lifecycle --json`.
- Open `M-198/W-0270 Select next Pitaya-aligned direction after session binding, kick/disconnect, and session data map` as next-ready.
- Keep session binding behavior, kick/disconnect behavior, session data behavior and persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, protocol changes, generated output, persistence, dependencies, metrics/tracing behavior, hosted surfaces, SDKs, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0177-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map/`
- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- Repository navigation docs and module guide updates for the W-0270 next-ready state.

## Open Questions

No open questions for W-0269. The selected follow-up is bounded and ready.

`W-0270 Select next Pitaya-aligned direction after session binding, kick/disconnect, and session data map` is now next-ready.

## Follow-Up

Proceed to `W-0270 Select next Pitaya-aligned direction after session binding, kick/disconnect, and session data map`.

Keep session binding behavior, kick/disconnect behavior, session data behavior and persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, runtime behavior, protocol shape, generated output, persistence, dependencies, metrics/tracing behavior, hosted surfaces, SDKs, and direct compatibility deferred until later bounded work items.

## Redaction Notes

No ignored credential file contents were read or printed. Push credentials remain loaded only through the ignored local environment file during push. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, or local secret values are recorded.
