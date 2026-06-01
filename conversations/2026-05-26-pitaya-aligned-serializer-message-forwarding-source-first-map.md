# Conversation: Pitaya-Aligned Serializer And Message Forwarding Source-First Map

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-implement-pitaya-aligned-serializer-message-forwarding-source-first-map/`
Related decision: `ADR-0171`

## Context

The maintainer asked to continue pushing toward Pitaya alignment with commit and push discipline. The active continuation queue after `W-0262` was `M-191/W-0263 Implement Pitaya-aligned serializer and message forwarding source-first map`.

`W-0262` had already defined the Pitaya-aligned serializer and message forwarding boundary gate, accepted `ADR-0170`, registered `runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate`, and opened the source-first serializer and message forwarding map as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Serializer and message forwarding vocabulary should become inspectable as future architecture vocabulary, not implemented as runtime behavior.

## Agent Response Summary

The agent treated W-0263 as a source-first inspection-map work item. It added `node tools/vibit inspect pitaya-serializer-forwarding --json`, accepted ADR-0171, registered the `runtime.pitaya_aligned_serializer_message_forwarding_source_first_map` check rule, completed W-0263, and opened W-0264 as the next Pitaya-aligned direction selection follow-up.

RED checks confirmed the command, rule, and change artifacts were initially absent:

```text
node tools/vibit inspect pitaya-serializer-forwarding --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_serializer_message_forwarding_source_first_map

node tools/vibit check change implement-pitaya-aligned-serializer-message-forwarding-source-first-map --json
# change directory does not exist
```

## Decisions

- `ADR-0171` implements the Pitaya-aligned serializer and message forwarding source-first map.
- The inspection command is `node tools/vibit inspect pitaya-serializer-forwarding --json`.
- The allowed vocabulary is `serializer_boundary`, `serializer_format`, `encode_boundary`, `decode_boundary`, `message_forwarding`, `forwarding_target`, `forwarding_envelope`, and `delivery_handoff`.
- Current vibit behavior remains Protobuf envelope ownership inside the protocol adapter, generated Protobuf payload bridge functions, single-process outbound message handling, metadata-only target scopes, no internal forwarding envelope, and single-process WebSocket delivery.
- W-0264 is the next-ready follow-up for selecting the next Pitaya-aligned direction.

## Artifacts

- `tools/vibit`
- `rules/check-rules.json`
- `decisions/ADR-0171-pitaya-aligned-serializer-message-forwarding-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-serializer-message-forwarding-source-first-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- Repository navigation docs and module guide updates for the W-0264 next-ready state.

## Open Questions

No runtime implementation question is answered by this source-first map. A later bounded work item must separately choose any route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, protocol carrier, persistence, dependency, service discovery, RPC, remote-call, frontend/backend role, cluster-safe session routing, or distributed runtime implementation.

## Follow-Up

- `M-192/W-0264 Select next Pitaya-aligned direction after serializer and message forwarding map`

## Redaction Notes

The inspection output exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, transport metadata, route payloads, or local secret values.
