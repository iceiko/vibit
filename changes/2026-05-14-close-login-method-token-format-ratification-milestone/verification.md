# Verification

Verified:
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change close-login-method-token-format-ratification-milestone --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan excluding `.git`, `.vibit.local.env`, and `node_modules`

Not verified:
- None for this workflow closeout.

Not applicable:
- Runtime authentication tests, because no runtime authentication behavior was added.
- Live PostgreSQL verification, because no schema, migration, repository, adapter, or persistence behavior changed.
- Protobuf generation, because no Protobuf source or generated Protobuf output changed.
