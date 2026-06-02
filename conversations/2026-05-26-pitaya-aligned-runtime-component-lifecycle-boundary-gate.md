# Conversation: Pitaya-Aligned Runtime Component Lifecycle Boundary Gate

Date: 2026-06-02
Work item: W-0280
Decision: ADR-0188
Rule: runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate

## Context

The maintainer asked to continue advancing toward Pitaya alignment in bounded steps and to commit and push completed increments. W-0279 selected `define_pitaya_aligned_runtime_component_lifecycle_boundary_gate`, accepted ADR-0187, and opened W-0280 as the next-ready gate work item.

The RED checks confirmed the expected missing surfaces:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate

node tools/vibit check change define-pitaya-aligned-runtime-component-lifecycle-boundary-gate --json
# change directory does not exist: changes/define-pitaya-aligned-runtime-component-lifecycle-boundary-gate
```

## Maintainer Narrative

Continue advancing in bounded steps toward Pitaya-class architecture coverage while preserving commit/push discipline and avoiding direct compatibility or runtime behavior expansion.

## Agent Response Summary

The agent defined the Pitaya-aligned runtime component lifecycle boundary gate. This completes W-0280, accepts ADR-0188, registers `runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate`, and opens M-209/W-0281 as next-ready.

Keep runtime component lifecycle behavior deferred until a later bounded work item explicitly authorizes it. Keep handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, distributed runtime behavior, hosted surfaces, SDKs, release artifacts, and direct compatibility deferred.

## Decisions

- Accept `ADR-0188`.
- Complete W-0280.
- Open W-0281 as next-ready.
- Define `runtime.pitaya_aligned_runtime_component_lifecycle_boundary_gate`.
- Keep runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, and direct compatibility deferred.

## Artifacts

- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
- `docs/pitaya-aligned-runtime-component-lifecycle-boundary-gate.zh-CN.md`
- `decisions/ADR-0188-pitaya-aligned-runtime-component-lifecycle-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-runtime-component-lifecycle-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No open questions for W-0280. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0281 Implement Pitaya-aligned runtime component lifecycle source-first map`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, component lifecycle payloads, handler registration payloads, dashboard payloads, admin console payloads, or concrete operational payloads are recorded.
