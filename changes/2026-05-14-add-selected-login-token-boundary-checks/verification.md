# Verification

Verified:
- `node tools/vibit inspect next --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-selected-login-token-boundary-checks --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan excluding `.git`, `.vibit.local.env`, and `node_modules`

Not verified:
- None for this tooling and standards change.

Not applicable:
- Live PostgreSQL verification, because no schema, migration, repository, adapter, or persistence behavior changed.
- Protobuf generation, because no Protobuf source or generated Protobuf output changed.
- Runtime authentication tests, because no runtime authentication behavior was added.
