# Verification

Verified:

- `node tools/vibit inspect next --json`
- `node tools/vibit inspect contracts --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect reference --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check work --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check change verify-m007-tooling-hardening-state --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification is not required because this change does not alter persistence behavior.
