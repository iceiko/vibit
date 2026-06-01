# Request

Define `W-0259 Define Pitaya-aligned route handler pipeline boundary gate`.

The work must stay gate-only. No route handler implementation, handler routing behavior, pipeline middleware behavior, serializer behavior, message forwarding behavior, backend route targeting, runtime behavior, protocol shape changes, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility may be added.

## User Requirement

Continue toward Pitaya alignment after the cluster-safe session routing source-first map by bounding route handler, handler pipeline, serializer, and forwarding vocabulary before implementation.

## Scope

In scope:

- Add the route handler pipeline boundary standard and Simplified Chinese translation.
- Add ADR-0167.
- Add change artifacts and conversation memory.
- Register `runtime.pitaya_aligned_route_handler_pipeline_boundary_gate`.
- Complete W-0259 and open W-0260 as a source-first route handler pipeline map.

Out of scope:

- Route handler implementation.
- Handler routing behavior.
- Handler pipeline behavior or pipeline middleware behavior.
- Serializer behavior.
- Message forwarding behavior.
- Backend route targeting.
- Runtime behavior, protocol messages or routes, Protobuf source, generated output, persistence, dependencies, hosted surfaces, SDK publication, release artifacts, or direct compatibility.
