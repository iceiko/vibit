# Impact

Allowed impact:

- Add `node tools/vibit inspect pitaya-dashboard-admin --json`.
- Add ADR-0186, conversation memory, change artifacts, rule catalog metadata, and runtime check coverage.
- Mark `M-206/W-0278` complete and open `M-207/W-0279 Select next Pitaya-aligned direction after dashboard/admin operations map` as next-ready.
- Update architecture manifests, module continuation metadata, and public continuation docs.

Forbidden impact:

- Runtime endpoint behavior.
- Metrics endpoints.
- Tracing pipelines or span runtime behavior.
- Observability pipelines.
- Dashboards or admin console behavior.
- Player/session/token inspectors.
- Event/audit tables.
- Admin users, roles, permissions, or authentication behavior.
- Protocol messages or routes.
- Protobuf sources.
- Generated output.
- Repository interfaces, PostgreSQL adapters, migrations, persistence, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.
