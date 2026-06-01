# Impact: Pitaya-Aligned Serializer And Message Forwarding Source-First Map

This change adds a source-first repository inspection map for Pitaya-aligned serializer and message forwarding vocabulary.

## Runtime Impact

No runtime behavior is added.

The map reports current vibit serializer and forwarding-adjacent concepts:

- protocol bridge;
- Protobuf envelope encoding;
- generated payload encoding and decoding;
- outbound message handling;
- target-scope metadata;
- absent internal forwarding envelope;
- single-process WebSocket delivery handoff.

It maps those concepts to future serializer and message forwarding vocabulary without changing dispatch, transport, protocol, domain, persistence, or startup behavior.

## Tooling Impact

`tools/vibit` gains:

- `node tools/vibit inspect pitaya-serializer-forwarding --json`;
- `runtime.pitaya_aligned_serializer_message_forwarding_source_first_map`;
- repository checks for W-0263 artifacts, inspection output markers, deferrals, redaction flags, and W-0264 next-ready state.

## Documentation And Memory Impact

The change records `ADR-0171`, a conversation memory entry, work item completion, runtime/reference/convention/contract manifest markers, and continuation docs pointing to W-0264.

## Explicit Non-Impact

- No protocol shape.
- No generated output.
- No persistence.
- No dependencies.
- No authentication/session behavior change.
- No route handler implementation.
- No handler routing, pipeline, middleware, serializer, forwarding, or backend route targeting behavior.
- No distributed runtime, service discovery, RPC, remote call, frontend/backend role, cluster-safe session routing, group, broadcast, hosted, SDK, release, or direct compatibility surface.
