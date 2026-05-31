# Request

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Interpreted Work Item

Advance the current next-ready work item:

```text
M-179/W-0251 Implement Pitaya-aligned server-to-server RPC source-first map
```

## Scope

Implement a source-first repository inspection map for the Pitaya-aligned server-to-server RPC and remote-call vocabulary defined by `ADR-0158`.

## Non-Goals

- No RPC, remote call, or service discovery implementation.
- No frontend/backend server role implementation.
- No distributed runtime behavior.
- No distributed groups, room broadcast fanout, delivery guarantees, or cluster-safe session routing.
- No runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.
