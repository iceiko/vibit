# Verification

Verified:

- node -c tools/vibit
- node tools/vibit check migrations --json
- node tools/vibit check runtime --json
- node tools/vibit check module authentication --json
- node tools/vibit check work --json
- node tools/vibit check memory --json
- node tools/vibit check change confirm-next-direction-after-postgres-session-persistence-schema-gate --json
- node tools/vibit check change implement-runtime-sessions-migration-source --json
- node tools/vibit check all --json
- git diff --check

Not verified:

- None.

Not applicable:

- Go tests are not required for this migration-source-only change because no Go runtime behavior is added.
- Buf generation is not required because no Protobuf source changes are added.
- Live PostgreSQL execution is not required by this slice because repository and adapter behavior remain deferred.
