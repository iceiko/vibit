# Plan

## Steps

1. Add `node tools/vibit inspect pitaya-roles --json` as a source-first inspection command.
2. Register `runtime.pitaya_aligned_frontend_backend_role_source_first_map`.
3. Record ADR-0157 and conversation memory for W-0249.
4. Update work-item, runtime, reference, convention, contract, module, README, AGENTS, alpha, maturity, and roadmap memory so W-0249 is complete and W-0250 is next-ready.
5. Run repository verification and update this change's checklist and verification records.

## Constraints

- Keep W-0249 source-first only.
- Do not add frontend/backend role implementation, distributed runtime behavior, server-to-server RPC, remote calls, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.
