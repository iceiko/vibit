# Verification

Verified:
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change compare-token-format-and-carrier-options --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:
- None.

Not applicable:
- Runtime Go tests, because no runtime code changed.
- Live PostgreSQL verification, because no persistence behavior or migration changed.
