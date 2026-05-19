# Verification

Verified:

- node -c tools/vibit
- node tools/vibit check runtime --json
- node tools/vibit check module authentication --json
- node tools/vibit check work --json
- node tools/vibit check memory --json
- node tools/vibit check change define-postgres-session-persistence-schema-gate --json
- node tools/vibit check all --json
- git diff --check

Not verified:

- None.

Not applicable:

- Go tests are not required for this gate-only change because no Go runtime behavior is added.
- Buf generation is not required because no Protobuf source changes are added.
- Live PostgreSQL verification is not applicable because no migration source or runtime persistence behavior is added.
