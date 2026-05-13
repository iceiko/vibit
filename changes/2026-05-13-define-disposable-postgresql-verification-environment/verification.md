# Verification

Verified:

- `node tools/vibit --help`
- `node tools/vibit check postgres-env --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-disposable-postgresql-verification-environment --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub token patterns in tracked changes.

Not verified:

- Live PostgreSQL verification commands. This change defines the opt-in environment standard and static check; it intentionally does not connect to PostgreSQL.

Not applicable:

- Live PostgreSQL verification; this change defines the environment standard but does not add live integration commands.
