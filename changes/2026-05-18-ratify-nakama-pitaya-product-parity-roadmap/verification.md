# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change ratify-nakama-pitaya-product-parity-roadmap --json`
- `node tools/vibit check all --json`
- `git diff --check`
- `node tools/vibit inspect next --json`

Not verified:

- None.

Not applicable:

- Go runtime tests, because this change does not alter Go runtime behavior.
