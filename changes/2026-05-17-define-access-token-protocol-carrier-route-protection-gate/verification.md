# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-access-token-protocol-carrier-route-protection-gate --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not applicable:

- Go runtime tests are not required because this gate does not add or modify Go runtime code.
- Protobuf generation is not required because this gate does not add `.proto` sources.
- Live PostgreSQL verification is not required because this gate does not change persistence behavior.
