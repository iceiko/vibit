# Impact

Allowed impact:

- Add `node tools/vibit inspect pitaya-metrics-tracing --json`.
- Add ADR-0183, conversation memory, change artifacts, rule catalog metadata, and runtime check coverage.
- Mark `M-203/W-0275` complete and open `M-204/W-0276 Select next Pitaya-aligned direction after metrics and tracing map` as next-ready.
- Update architecture manifests, module continuation metadata, and public continuation docs.

Forbidden impact:

- Runtime endpoint behavior.
- Metrics endpoints.
- Tracing pipelines or span runtime behavior.
- Observability pipelines.
- Dashboards or admin console behavior.
- Player/session/token inspectors.
- Event/audit tables.
- Protocol messages or routes.
- Protobuf sources.
- Generated output.
- Repository interfaces, PostgreSQL adapters, migrations, persistence, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.
