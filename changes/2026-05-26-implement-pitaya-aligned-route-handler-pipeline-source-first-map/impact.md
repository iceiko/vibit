# Impact: Pitaya-Aligned Route Handler Pipeline Source-First Map

This change adds a source-first repository inspection map for Pitaya-aligned route handler pipeline vocabulary.

## Runtime Impact

No runtime behavior is added.

The map reports current vibit route-flow concepts:

- protocol envelope;
- route request handoff;
- application dispatch;
- transactional dispatch;
- protocol bridge;
- outbound message handling;
- target-scope metadata.

It maps those concepts to future route handler pipeline vocabulary without changing dispatch, transport, protocol, domain, persistence, or startup behavior.

## Tooling Impact

`tools/vibit` gains:

- `node tools/vibit inspect pitaya-routes --json`;
- `runtime.pitaya_aligned_route_handler_pipeline_source_first_map`;
- repository checks for W-0260 artifacts, inspection output markers, deferrals, redaction flags, and W-0261 next-ready state.

## Documentation And Memory Impact

The change records `ADR-0168`, a conversation memory entry, work item completion, runtime/reference/convention/contract manifest markers, and continuation docs pointing to W-0261.

## Explicit Non-Impact

- No protocol shape.
- No generated output.
- No persistence.
- No dependencies.
- No authentication/session behavior change.
- No route handler implementation.
- No handler routing, pipeline, middleware, serializer, forwarding, or backend route targeting behavior.
- No distributed runtime, service discovery, RPC, remote call, frontend/backend role, cluster-safe session routing, group, broadcast, hosted, SDK, release, or direct compatibility surface.
