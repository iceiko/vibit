# Verification

Verified:

- `node tools/vibit check work --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check change confirm-next-direction-after-authentication-design --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- Live PostgreSQL execution. This change does not add or modify persistence behavior.

Not applicable:

- Go runtime tests beyond repository `check all`, because this change does not add or modify Go runtime code.
- Protobuf generation, because this change does not add or modify `.proto` sources or generated Protobuf output.
- Migration apply or rollback verification, because this change does not add or modify migration sources.
