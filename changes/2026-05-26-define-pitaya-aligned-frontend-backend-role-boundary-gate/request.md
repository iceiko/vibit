# Request

Date: 2026-05-31

## User Request

```text
继续推进20步，向pitaya靠拢
```

## Interpreted Work Item

```text
M-176/W-0248 Define Pitaya-aligned frontend/backend role boundary gate
```

## Scope

Define a gate-only Pitaya-aligned frontend/backend role vocabulary boundary after the source-first Pitaya vocabulary map.

The slice must not implement frontend/backend server roles, distributed runtime behavior, server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.

No frontend/backend server role implementation is authorized by this request.

## Success Criteria

- A standard defines how `frontend_server` and `backend_server` vocabulary may be used for future planning.
- `ADR-0156` accepts the boundary.
- `runtime.pitaya_aligned_frontend_backend_role_boundary_gate` is registered.
- W-0248 is completed and W-0249 is opened as a source-first role map follow-up.
