# Request

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Interpreted Work Item

Advance the current next-ready work item:

```text
M-175/W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map
```

## Scope

Implement a source-first repository inspection map for the Pitaya-aligned distributed runtime vocabulary defined by `ADR-0154`.

## Non-Goals

- No distributed runtime implementation.
- No frontend/backend server role implementation.
- No server-to-server RPC or remote call behavior.
- No service discovery.
- No distributed groups, room broadcast fanout, delivery guarantees, or cluster-safe session routing.
- No runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.
