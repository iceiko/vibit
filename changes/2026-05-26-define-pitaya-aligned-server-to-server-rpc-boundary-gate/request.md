# Request

Date: 2026-05-31

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Interpreted Work Item

```text
M-178/W-0250 Define Pitaya-aligned server-to-server RPC boundary gate
```

## Scope

Define a gate-only Pitaya-aligned server-to-server RPC and remote-call vocabulary boundary after the source-first Pitaya role map.

The slice must not implement server-to-server RPC, remote calls, service discovery, frontend/backend server roles, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.

No RPC, remote call, or service discovery implementation is authorized by this request.

## Success Criteria

- A standard defines how `server_to_server_rpc` and `remote_call` vocabulary may be used for future planning.
- `ADR-0158` accepts the boundary.
- `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate` is registered.
- W-0250 is completed and W-0251 is opened as a source-first RPC map follow-up.
