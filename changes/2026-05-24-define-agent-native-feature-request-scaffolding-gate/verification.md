# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.agent_native_feature_request_scaffolding_gate`
- `node tools/vibit check change define-agent-native-feature-request-scaffolding-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go tests are not required because this gate changes no Go source and no runtime behavior.

