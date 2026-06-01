# Conversation: Pitaya-Aligned Route Handler Pipeline Boundary Gate

Date: 2026-06-01
Work item: W-0259
Decision: ADR-0167
Check rule: runtime.pitaya_aligned_route_handler_pipeline_boundary_gate

## Context

The maintainer asked to continue moving toward Pitaya alignment with commit and push discipline.

`W-0258` selected `define_pitaya_aligned_route_handler_pipeline_boundary_gate` as the next bounded Pitaya-aligned direction after the cluster-safe session routing source-first map.

The expected W-0259 rule and change directory were absent before implementation:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_route_handler_pipeline_boundary_gate
Unknown rule_id: runtime.pitaya_aligned_route_handler_pipeline_boundary_gate
```

```text
node tools/vibit check change define-pitaya-aligned-route-handler-pipeline-boundary-gate --json
change directory does not exist: changes/define-pitaya-aligned-route-handler-pipeline-boundary-gate
```

## Maintainer Narrative

The maintainer wanted continued Pitaya alignment, but the repository must keep route handler, pipeline, serializer, and forwarding vocabulary separated from runtime behavior until a later bounded work item authorizes implementation.

## Agent Response Summary

The agent defined a gate-only Pitaya-aligned route handler pipeline boundary and opened a source-first map follow-up.

The gate maps current vibit protocol envelope, route request, application dispatch, transactional dispatch, protocol bridge, and outbound message surfaces to deferred route handler pipeline vocabulary.

## Decisions

The allowed future vocabulary is `route_handler`, `route_key`, `handler_dispatch`, `handler_pipeline`, `pipeline_step`, `serializer_boundary`, `message_forwarding`, and `route_target`.

W-0259 does not implement route handlers, pipelines, serializers, or forwarding.

The follow-up is `W-0260 Implement Pitaya-aligned route handler pipeline source-first map`.

## Artifacts

- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md`
- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.zh-CN.md`
- `decisions/ADR-0167-pitaya-aligned-route-handler-pipeline-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-route-handler-pipeline-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `rules/check-rules.json`
- `tools/vibit`
- Repository navigation docs and module guide updates for the W-0260 next-ready state.

## Open Questions

No open questions for W-0259. The next work item is bounded and ready.

## Follow-Up

Proceed to `W-0260 Implement Pitaya-aligned route handler pipeline source-first map`.

## Redaction Notes

No credentials, tokens, raw verifier material, lookup digests, verifier digests, or local secret values were recorded.
