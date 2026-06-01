# Conversation: Select Next Pitaya-Aligned Direction After Cluster-Safe Session Routing Map

Date: 2026-06-01
Work item: W-0258
Decision: ADR-0166
Check rule: runtime.next_pitaya_aligned_direction_after_cluster_safe_session_routing_map

## Context

Continue moving toward Pitaya alignment and preserve commit/push discipline.

`W-0257` had completed `node tools/vibit inspect pitaya-sessions --json`, registered `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map`, and opened `W-0258 Select next Pitaya-aligned direction after cluster-safe session routing map`.

The expected W-0258 rule and change directory were absent before implementation:

```text
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_cluster_safe_session_routing_map
Unknown rule_id: runtime.next_pitaya_aligned_direction_after_cluster_safe_session_routing_map
```

```text
node tools/vibit check change select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map --json
change directory does not exist: changes/select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map
```

## Maintainer Narrative

The maintainer asked to continue moving toward Pitaya alignment and to keep commit and push discipline.

## Agent Response Summary

The agent selected the next bounded Pitaya-aligned direction after the completed cluster-safe session routing map.

The selected direction is `define_pitaya_aligned_route_handler_pipeline_boundary_gate`, which opens `W-0259 Define Pitaya-aligned route handler pipeline boundary gate`.

## Decisions

The completed Pitaya source-first maps cover distributed runtime vocabulary, frontend/backend roles, server-to-server RPC, service discovery, distributed group/broadcast, and cluster-safe session routing.

The next useful Pitaya reference surface is route handler, pipeline, serializer, and forwarding vocabulary.

This W-0258 selection is planning-only. It does not implement route handlers, pipelines, serializers, forwarding, runtime behavior, protocol routes, Protobuf sources, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.

W-0258 does not implement route handlers, pipelines, serializers, or forwarding.

## Artifacts

- `decisions/ADR-0166-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map.md`
- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-cluster-safe-session-routing-map/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `rules/check-rules.json`
- `tools/vibit`
- Repository navigation docs and module guide updates for the W-0259 next-ready state.

## Open Questions

No open questions for W-0258. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0259 Define Pitaya-aligned route handler pipeline boundary gate`.

## Redaction Notes

No credentials, tokens, raw verifier material, lookup digests, verifier digests, or local secret values were recorded.
