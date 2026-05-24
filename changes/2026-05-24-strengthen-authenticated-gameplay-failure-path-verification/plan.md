# Plan

## Files To Create

1. `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/request.md`
2. `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/spec.yaml`
3. `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/impact.md`
4. `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/plan.md`
5. `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/checklist.md`
6. `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/verification.md`
7. `decisions/ADR-0131-authenticated-gameplay-failure-path-verification.md`
8. `conversations/2026-05-24-authenticated-gameplay-failure-path-verification.md`

## Files To Edit

1. `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
2. `.arch/work-items.yaml`
3. `.arch/runtime.yaml`
4. `.arch/reference.yaml`
5. `.arch/contracts.yaml`
6. `.arch/conventions.yaml`
7. `.arch/modules.yaml`
8. `modules/storage/module.yaml`
9. `modules/storage/AGENTS.md`
10. `modules/storage/AGENTS.zh-CN.md`
11. `README.md`
12. `README.zh-CN.md`
13. `AGENTS.md`
14. `AGENTS.zh-CN.md`
15. `runtime/AGENTS.md`
16. `runtime/AGENTS.zh-CN.md`
17. `docs/v0.1-alpha-goal.md`
18. `docs/v0.1-alpha-goal.zh-CN.md`
19. `docs/alpha-developer-flow.md`
20. `docs/alpha-developer-flow.zh-CN.md`
21. `docs/alpha-acceptance-checklist.md`
22. `docs/alpha-acceptance-checklist.zh-CN.md`
23. `docs/product-maturity-milestones.md`
24. `docs/product-maturity-milestones.zh-CN.md`
25. `docs/nakama-pitaya-product-parity-roadmap.md`
26. `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
27. `tools/vibit`
28. `rules/check-rules.json`

## Generated Artifacts

None. This slice does not add Protobuf sources, generated Go output, generated contract shapes, migrations, release artifacts, or dependency metadata.

## Handwritten Logic

Add test-only local alpha E2E coverage for authenticated gameplay failure paths. The production runtime, protocol, storage, authentication, session, and route-policy code paths remain unchanged.

## Tests

Add `TestAuthenticatedGameplayFailurePathsLocalAlphaFlow` in `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`.

The test uses the existing local alpha E2E fixture and covers:

- missing authenticated wrapper on protected inventory;
- malformed authenticated wrapper on protected inventory;
- malformed access-token text;
- unknown well-formed access-token text;
- expired token record;
- revoked token after logout;
- missing authenticated wrapper on protected presence;
- raw credential/token redaction in error envelopes.

## Verification Commands

```sh
cd runtime && go test ./internal/platform/protocol/protobuf
cd runtime && go test ./...
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.authenticated_gameplay_failure_path_verification
node tools/vibit check change strengthen-authenticated-gameplay-failure-path-verification --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback Or Migration Notes

Rollback removes the E2E test additions, the W-0223 records, ADR-0131, rule registration, and W-0224 next-ready update. No data, protocol, generated output, dependency, or migration rollback is required.
