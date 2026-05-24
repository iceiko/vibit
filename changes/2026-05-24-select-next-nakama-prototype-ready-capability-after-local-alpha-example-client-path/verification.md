# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.agent_native_feature_request_scaffolding_selection`
- `node tools/vibit check change select-next-nakama-prototype-ready-capability-after-local-alpha-example-client-path --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go tests are not required for this selection slice because no Go source or runtime behavior changes are made.

Notes:

- `node tools/vibit check runtime --json` passed with the existing non-blocking `runtime.identity_boundary` warning on `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
