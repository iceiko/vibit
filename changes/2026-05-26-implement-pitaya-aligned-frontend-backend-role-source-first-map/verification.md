# Verification

Status: Passed with one pre-existing runtime warning.

Final verification was run on 2026-05-31.

```sh
node -c tools/vibit
```

Result: passed with exit code 0.

```sh
node tools/vibit inspect pitaya-roles --json
```

Result: passed. Output reported `kind: pitaya_roles_inspection`, `work_item: W-0249`, `check_rule: runtime.pitaya_aligned_frontend_backend_role_source_first_map`, `implementation_decision: ADR-0157`, allowed roles `frontend_server` and `backend_server`, current single-process role mapping, all W-0249 deferrals as false-added surfaces, and next-ready `M-178/W-0250`.

```sh
node tools/vibit inspect next --json
```

Result: passed. Output reported `status: ready`, current milestone `M-178`, and next-ready `W-0250 Define Pitaya-aligned server-to-server RPC boundary gate`.

```sh
node tools/vibit inspect rule runtime.pitaya_aligned_frontend_backend_role_source_first_map
```

Result: passed. Output returned the rule catalog entry for `runtime.pitaya_aligned_frontend_backend_role_source_first_map`.

```sh
node tools/vibit check change implement-pitaya-aligned-frontend-backend-role-source-first-map --json
```

Result: passed with summary `passed: 13`, `warnings: 0`, `failures: 0`.

```sh
node tools/vibit check work --json
```

Result: passed with summary `passed: 1512`, `warnings: 0`, `failures: 0`; summary counts reported `work_items: 250`, `completed: 249`, `next_ready: 1`, and next-ready `W-0250`.

```sh
node tools/vibit check runtime --json
```

Result: passed with summary `passed: 20059`, `warnings: 1`, `failures: 0`; summary counts reported `go_files: 154` and `test_files: 60`. The warning is the pre-existing accepted warning:

```text
runtime.identity_boundary
runtime/internal/platform/persistence/postgres/authentication_repository.go mentions credential dependency; keep it behind an explicit ratified boundary
```

```sh
node tools/vibit check memory --json
```

Result: passed with summary `passed: 4220`, `warnings: 0`, `failures: 0`; summary counts reported `conversations: 202` and `decisions: 157`.

```sh
node tools/vibit check schemas --json
```

Result: passed with summary `passed: 4634`, `warnings: 0`, `failures: 0`.

```sh
node tools/vibit check all --json
```

Result: passed with summary `subchecks: 304`, `passed: 304`, `warnings: 1`, `failures: 0`. The warning aligns with the known `runtime.identity_boundary` warning recorded above.

```sh
git diff --check
```

Result: passed with exit code 0.
