# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check module player --json`
- `node tools/vibit check change define-credential-storage-external-identity-linking-boundaries --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go runtime tests are not required because this change does not modify Go runtime code.
