# Verification

Verified:
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change ratify-first-token-format-and-proof-carrier-posture --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:
- None.

Not applicable:
- Runtime Go tests beyond default repository checks, because no runtime code changed.
- Live PostgreSQL verification, because no persistence behavior or migration changed.
