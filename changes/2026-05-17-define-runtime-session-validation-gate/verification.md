# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check change define-runtime-session-validation-gate --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None yet.

Not applicable:

- Go tests are not required because this gate adds no runtime code.
