# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change close-authentication-token-session-validation-design-milestone --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

Not verified:

- Live PostgreSQL execution. This change does not add or modify persistence behavior, and default repository checks do not require a running PostgreSQL server.

Not applicable:

- Go runtime tests beyond repository `check runtime`, because this change does not add or modify Go runtime code.
- Protobuf generation, because this change does not add or modify `.proto` sources or generated Protobuf output.
- Migration apply or rollback verification, because this change does not add or modify migration sources.
- WebSocket request-loop verification for authentication, player account handlers, or routes, because those behaviors remain deferred behind the confirmation gate.
