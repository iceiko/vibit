# Request

Implement `W-0272`, a source-first Pitaya-aligned runtime observability inspection map after the runtime observability boundary gate.

The change must:

- Accept `ADR-0180`.
- Register `runtime.pitaya_aligned_runtime_observability_source_first_map`.
- Add `node tools/vibit inspect pitaya-observability --json`.
- Report the runtime observability boundary vocabulary, current operations inspection, health/readiness/version/config posture, route inventory, verification posture, redaction posture, source surfaces, and deferrals.
- Open `W-0273 Select next Pitaya-aligned direction after runtime observability map` as next-ready.

No runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, transport behavior changes, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility may be added.
