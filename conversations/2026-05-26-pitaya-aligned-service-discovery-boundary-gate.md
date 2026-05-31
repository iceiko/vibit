# Conversation: Pitaya-Aligned Service Discovery Boundary Gate

Date: 2026-05-31
Status: Accepted
Related change: `changes/2026-05-26-define-pitaya-aligned-service-discovery-boundary-gate/`
Related decision: `ADR-0160`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. The active continuation queue after `W-0251` was `M-180/W-0252 Define Pitaya-aligned service discovery boundary gate`.

`W-0251` had already implemented `node tools/vibit inspect pitaya-rpc --json`, accepted `ADR-0159`, registered `runtime.pitaya_aligned_server_to_server_rpc_source_first_map`, and opened the service discovery boundary gate as next-ready.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Service discovery is useful future architecture vocabulary, not permission to add a distributed runtime.

## Agent Response Summary

The agent treated W-0252 as a gate-only work item. It added a service discovery boundary standard and translation, accepted ADR-0160, registered the `runtime.pitaya_aligned_service_discovery_boundary_gate` check rule, completed W-0252, and opened W-0253 as the source-first service discovery map follow-up.

RED checks confirmed the rule and change artifacts were initially absent:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_service_discovery_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_service_discovery_boundary_gate

node tools/vibit check change define-pitaya-aligned-service-discovery-boundary-gate --json
# change directory does not exist
```

## Decisions

- `ADR-0160` defines the Pitaya-aligned service discovery boundary gate.
- The allowed vocabulary is `service_discovery`, `service_registry`, `service_instance`, and `service_selector`.
- Current vibit behavior remains static single-process startup composition and direct in-process dispatch.
- W-0253 is the next-ready follow-up for a source-first service discovery map.

## Artifacts

- `docs/pitaya-aligned-service-discovery-boundary-gate.md`
- `docs/pitaya-aligned-service-discovery-boundary-gate.zh-CN.md`
- `decisions/ADR-0160-pitaya-aligned-service-discovery-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-service-discovery-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

No runtime implementation question is answered by this gate. A later bounded work item must separately choose any service discovery model, registry storage, selector behavior, node identity, membership, heartbeat, topology, dependency, RPC, or remote-call implementation.

## Follow-Up

- `M-181/W-0253 Implement Pitaya-aligned service discovery source-first map`

## Redaction Notes

The gate exposes no raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, PostgreSQL DSNs, database payloads, local secret file contents, node credentials, or transport metadata.
