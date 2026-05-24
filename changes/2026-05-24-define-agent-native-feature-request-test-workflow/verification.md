# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.agent_native_feature_request_test_workflow`
- `node tools/vibit check change define-agent-native-feature-request-test-workflow --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go runtime tests are not required for this slice because it does not change runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, startup wiring, or Go code.

Expected warnings:

- `runtime.identity_boundary` may continue to warn on `runtime/internal/platform/persistence/postgres/authentication_repository.go`; this is a known pre-existing repository warning.

