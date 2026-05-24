# Plan

1. Review `ADR-0132`, the agent-native workflow standard, the local alpha package, and existing examples.
2. Define the allowed source-first local alpha example client path.
3. Record ownership, file candidates, accepted existing routes, redaction rules, verification expectations, and stop conditions.
4. Add `ADR-0133`.
5. Mark `M-153/W-0225` completed and open `M-154/W-0226`.
6. Update runtime, reference, contracts, conventions, module, alpha, roadmap, and agent guide records.
7. Add repository check coverage for `runtime.local_alpha_example_client_path_gate`.
8. Run focused and repository-wide verification.

## Selected Path

The gate selects a source-first repository-local path:

```text
examples/local-alpha-client/README.md
examples/local-alpha-example-client.sh
runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
```

The future implementation may add a clearer docs/script entrypoint and, if needed, a focused local alpha example-flow test under the existing Protobuf E2E test boundary.

## Verification

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.local_alpha_example_client_path_gate
node tools/vibit check change define-local-alpha-example-client-path-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

No Go tests are required for this gate unless runtime source changes are added.

## Rollback

Rollback removes the W-0225 gate records, ADR-0133, rule registration, and W-0226 next-ready update. No runtime, protocol, generated output, migration, dependency, data, release, hosted, SDK, or direct compatibility rollback is required.
