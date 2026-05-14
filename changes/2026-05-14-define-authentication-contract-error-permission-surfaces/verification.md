# Verification

Verified:
- `node tools/vibit inspect next --json`
- `node tools/vibit check contracts --json`
- `node tools/vibit inspect contracts --module runtime --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-authentication-contract-error-permission-surfaces --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan excluding `.git`, `.vibit.local.env`, and `node_modules`

Not verified:
- None for this design and contract-surface change.

Not applicable:
- Go runtime behavior tests beyond default repository checks, because no runtime authentication behavior changed.
- Live PostgreSQL verification, because no schema, migration, repository, adapter, or persistence behavior changed.
- Protobuf generation, because no Protobuf source or generated Protobuf output changed.
