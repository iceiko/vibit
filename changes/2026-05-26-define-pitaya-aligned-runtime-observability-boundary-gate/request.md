# Request

Define `W-0271`, a gate-only Pitaya-aligned runtime observability boundary gate after the session binding, kick/disconnect, and session data direction selection.

The change must:

- Accept `ADR-0179`.
- Define the English standard and Simplified Chinese translation.
- Register `runtime.pitaya_aligned_runtime_observability_boundary_gate`.
- Map current minimum operations inspection, existing health/readiness/version/config endpoint summaries, route inventory, repository verification, redaction posture, and deferred operations surfaces to future observability vocabulary.
- Open `W-0272 Implement Pitaya-aligned runtime observability source-first map` as next-ready.

No runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, transport behavior changes, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility may be added.
