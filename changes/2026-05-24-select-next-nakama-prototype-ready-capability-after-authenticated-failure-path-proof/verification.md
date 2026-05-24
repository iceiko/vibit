# Verification

Date: 2026-05-24

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.next_nakama_prototype_ready_capability_selection
node tools/vibit check change select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Status: Verified.

Notes:

- This selection slice adds no runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, persistence, startup wiring, SDK, hosted deployment, release artifact, or direct compatibility scope.
- Runtime test coverage remains exercised through `node tools/vibit check runtime --json`, which runs `go test ./...` under `runtime`.
