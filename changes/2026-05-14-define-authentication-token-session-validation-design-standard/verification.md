# Verification

Verified:

- `node tools/vibit check architecture --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change define-authentication-token-session-validation-design-standard --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go runtime tests are not required because this change does not modify Go runtime code.
