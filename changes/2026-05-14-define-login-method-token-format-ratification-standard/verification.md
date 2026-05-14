# Verification

Verified:
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-login-method-token-format-ratification-standard --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for concrete GitHub token patterns, excluding `.git`, `.vibit.local.env`, and `node_modules`.

Not verified:
- None.

Not applicable:
- Runtime Go tests, because no runtime code changed.
- Live PostgreSQL verification, because no persistence behavior or migration changed.
