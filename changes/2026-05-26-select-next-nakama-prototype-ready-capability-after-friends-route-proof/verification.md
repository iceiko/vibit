# Verification

Date: 2026-05-31

Initial RED check:

```bash
node tools/vibit inspect rule runtime.next_nakama_prototype_ready_capability_after_friends_route_proof
```

Result before implementation:

```text
Unknown rule_id: runtime.next_nakama_prototype_ready_capability_after_friends_route_proof
```

Required final commands:

```bash
node -c tools/vibit
node tools/vibit inspect rule runtime.next_nakama_prototype_ready_capability_after_friends_route_proof
node tools/vibit inspect next --json
node tools/vibit check change select-next-nakama-prototype-ready-capability-after-friends-route-proof --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Final verification:

```text
node -c tools/vibit
node tools/vibit inspect rule runtime.next_nakama_prototype_ready_capability_after_friends_route_proof
node tools/vibit inspect next --json
node tools/vibit check change select-next-nakama-prototype-ready-capability-after-friends-route-proof --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Final result:

```text
node -c tools/vibit: passed
inspect rule: passed; rule is registered in rules/check-rules.json
inspect next: passed; M-172 / W-0244 is the next-ready repository direction
check change: passed; 13 passed, 0 warnings, 0 failures
check work: passed; 1476 passed, 0 warnings, 0 failures
check runtime: passed; 19022 passed, 1 warning, 0 failures
check memory: passed; 4076 passed, 0 warnings, 0 failures
check schemas: passed; 4502 passed, 0 warnings, 0 failures
check all: passed; 298 subchecks passed, 1 warning, 0 failures
git diff --check: passed
```

Status: Verified.

Known warning:

- `runtime.identity_boundary` still warns on the existing `runtime/internal/platform/persistence/postgres/authentication_repository.go` identity boundary; this warning predates this selection slice and remains non-blocking.

Notes:

- This selection slice adds no runtime behavior, operations/admin endpoint, metrics endpoint, observability pipeline, protocol route, Protobuf source, generated output, migration, dependency, persistence, startup wiring, SDK, hosted deployment, release artifact, or direct compatibility scope.
- Runtime test coverage remains exercised through `node tools/vibit check runtime --json`, which runs `go test ./...` under `runtime`.
