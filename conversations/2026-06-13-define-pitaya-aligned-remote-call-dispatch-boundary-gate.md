# Conversation: Pitaya-Aligned Remote Call Dispatch Boundary Gate

Date: 2026-06-13

## Context

The maintainer asked to commit and push current modifications, keep the ignored local key file out of Git, and continue ten steps toward Pitaya alignment. The approved scope was source-first Pitaya service dispatch alignment only.

## Maintainer Narrative

The maintainer confirmed the ten-step plan with `确认`. The work should move vibit toward Pitaya-class concepts while preserving vibit's agent-native maintainability model and avoiding runtime compatibility commitments.

## Agent Response Summary

The agent recorded `W-0305`, accepted `ADR-0213`, and registered `runtime.pitaya_aligned_remote_call_dispatch_boundary_gate`. Defined gate-only vocabulary for Pitaya-aligned remote call dispatch.

No runtime behavior, protocol route, Protobuf source, generated output, persistence, dependency, distributed runtime implementation, or direct Nakama/Pitaya API compatibility was added.

## Decisions

- `ADR-0213` records this source-first boundary or map.
- `runtime.pitaya_aligned_remote_call_dispatch_boundary_gate` is the repository check rule for this slice.
- `W-0305` is completed in `.arch/work-items.yaml`.

## Artifacts

- `changes/2026-06-13-define-pitaya-aligned-remote-call-dispatch-boundary-gate`
- `decisions/ADR-0213-pitaya-aligned-remote-call-dispatch-boundary-gate.md`
- `docs/pitaya-aligned-remote-call-dispatch-boundary-gate.md`
- `docs/pitaya-aligned-remote-call-dispatch-boundary-gate.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

None for this bounded source-first slice. Future runtime behavior remains behind later explicit work items and ADRs.

## Follow-Up

Opened W-0306 to implement the Pitaya-aligned remote call dispatch source-first map.

## Redaction Notes

The ignored local environment file was not copied into project memory. Raw credentials, access tokens, lookup digests, verifier digests, verifier keys, DSNs with credentials, HTTP headers, cookies, query strings, WebSocket transport metadata, route payloads, event payloads, and local secret file contents remain out of this log.

## Verification Markers

```yaml
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
distributed_runtime_implementation_added: false
direct_nakama_pitaya_api_compatibility_added: false
```
