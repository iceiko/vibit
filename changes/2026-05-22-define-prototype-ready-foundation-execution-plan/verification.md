# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect next`
- `node tools/vibit check change define-prototype-ready-foundation-execution-plan --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

Known accepted warning:

- `runtime.identity_boundary`

Not verified:

- None.

Not applicable:

- Go runtime tests, because this change does not modify Go runtime code.
- Live PostgreSQL verification, because this change does not modify runtime behavior or persistence behavior.

