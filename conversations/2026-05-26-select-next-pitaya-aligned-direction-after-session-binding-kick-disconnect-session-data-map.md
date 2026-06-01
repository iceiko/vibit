# Conversation: Select Next Pitaya-Aligned Direction After Session Binding Kick Disconnect Session Data Map

Date: 2026-06-01

## Context

The maintainer asked to continue advancing toward Pitaya alignment. The repository next-ready item was `W-0270 Select next Pitaya-aligned direction after session binding, kick/disconnect, and session data map`, opened by `ADR-0177` after `node tools/vibit inspect pitaya-session-lifecycle --json`.

The RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_session_binding_kick_disconnect_session_data_map
# Unknown rule_id: runtime.next_pitaya_aligned_direction_after_session_binding_kick_disconnect_session_data_map

node tools/vibit check change select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map --json
# change directory does not exist
```

## Maintainer Narrative

Continue advancing in bounded steps toward Pitaya-class architecture coverage while preserving commit/push discipline and avoiding direct compatibility or runtime behavior expansion.

## Agent Response Summary

The agent selected `define_pitaya_aligned_runtime_observability_boundary_gate` as the next bounded Pitaya-aligned direction. This completes W-0270, accepts `ADR-0178`, registers `runtime.next_pitaya_aligned_direction_after_session_binding_kick_disconnect_session_data_map`, and opens `M-199/W-0271 Define Pitaya-aligned runtime observability boundary gate` as next-ready.

## Decisions

- Accept `ADR-0178`.
- Complete W-0270.
- Open W-0271 as next-ready.
- Keep runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, session binding behavior, kick/disconnect behavior, session data behavior and persistence, transport behavior changes, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility deferred.

## Artifacts

- `decisions/ADR-0178-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for W-0270. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0271 Define Pitaya-aligned runtime observability boundary gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, or concrete operational payloads are recorded.
