# Verification

Verified:

- `node -c tools/vibit`
- `node -e "JSON.parse(require('fs').readFileSync('rules/check-rules.json','utf8'))"`
- `node tools/vibit inspect rule runtime.identity_boundary`
- `node tools/vibit check runtime --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-metadata-only-identity-repository-checks --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification is not required for repository boundary checks.
- Authentication, token, credential, session store, player account migration, Protobuf envelope, and WebSocket handshake verification are not applicable because this change does not implement or change those surfaces.
