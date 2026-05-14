# Verification

Verified:

- `node tools/vibit check work --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check change confirm-next-direction-after-player-account-postgresql-persistence --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go runtime tests, because this change does not add or modify runtime Go behavior.
- Live PostgreSQL verification, because this change does not add or modify persistence behavior.
