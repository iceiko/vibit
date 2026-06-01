# Verification

## RED Checks

The initial checks failed because the W-0256 rule and change artifacts were absent:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate

node tools/vibit check change define-pitaya-aligned-cluster-safe-session-routing-boundary-gate --json
# change directory does not exist
```

## Required Verification

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate`
- `node tools/vibit check change define-pitaya-aligned-cluster-safe-session-routing-boundary-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Result

Verified for this commit:

- `node -c tools/vibit`: passed.
- `node tools/vibit inspect next --json`: next-ready `W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map`.
- `node tools/vibit inspect rule runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate`: passed and reported `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate`.
- `node tools/vibit check change define-pitaya-aligned-cluster-safe-session-routing-boundary-gate --json`: passed 13, warnings 0, failures 0.
- `node tools/vibit check work --json`: passed 1554, warnings 0, failures 0.
- `node tools/vibit check runtime --json`: passed 21644, warnings 1, failures 0. The warning is the existing `runtime.identity_boundary` warning on `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- `node tools/vibit check memory --json`: passed 4388, warnings 0, failures 0.
- `node tools/vibit check schemas --json`: passed 4787, warnings 0, failures 0.
- `node tools/vibit check all --json`: passed 311 subchecks, warnings 1, failures 0.
- `git diff --check`: passed.
