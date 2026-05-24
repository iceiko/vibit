# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit scaffold feature scaffold-smoke --date 2026-05-24 --request "Smoke test feature request scaffold." --summary "Smoke test feature request scaffold." --dry-run`
- `node tools/vibit scaffold feature implement-agent-native-feature-request-scaffolding --date 2026-05-24 --request "<redacted in command history as normal request text>" --summary "Implement source-first feature request scaffolding."`
- `node tools/vibit inspect rule runtime.agent_native_feature_request_scaffolding_implementation`
- `node tools/vibit check change implement-agent-native-feature-request-scaffolding --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Separate Go runtime tests are not required for this slice because it changes no Go runtime source or runtime behavior. `node tools/vibit check runtime --json` still runs the repository Go tests.
- Protocol generated output checks are not separately required beyond repository checks because this slice adds no Protobuf source or generated output.
- Migration checks are not separately required beyond repository checks because this slice adds no migration source.

Redaction:

- No raw credentials, raw access tokens, verifier keys, digests, DSNs with credentials, GitHub tokens, transport metadata, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private user data were recorded.
