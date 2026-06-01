# Impact

## In Scope

- Add a gate-only acceptor and connection lifecycle standard and translation.
- Add `ADR-0173`.
- Register `runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate`.
- Mark `W-0265` complete and open `W-0266` as next-ready.
- Update architecture manifests, module manifests, public continuation docs, and agent guides.

## Out Of Scope

- No acceptor behavior.
- No TCP acceptors.
- No WebSocket behavior changes.
- No connection lifecycle behavior changes.
- No session binding behavior.
- No kick/disconnect behavior.
- No protocol shape changes.
- No generated output.
- No repository, PostgreSQL adapter, migration, or dependency changes.
- No metrics endpoints, tracing pipelines, dashboards, hosted surfaces, SDKs, release artifacts, or direct compatibility.

## Risk

The main risk is vocabulary being mistaken for implementation permission. The gate mitigates that by recording explicit no-behavior flags and a narrow source-first map follow-up.
