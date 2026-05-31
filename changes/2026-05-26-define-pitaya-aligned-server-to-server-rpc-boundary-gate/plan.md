# Plan

## Steps

1. Define `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md` and its Simplified Chinese translation.
2. Register `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`.
3. Record ADR-0158 and conversation memory for W-0250.
4. Update work-item, runtime, reference, convention, contract, module, README, AGENTS, alpha, maturity, and roadmap memory so W-0250 is complete and W-0251 is next-ready.
5. Run repository verification and update this change's checklist and verification records.

## Constraints

- Keep W-0250 gate-only.
- Do not add server-to-server RPC implementation, remote call behavior, service discovery, frontend/backend server role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.
