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
- `node tools/vibit check memory --json`
- `node tools/vibit check change verify-agent-tooling-generated-contract-shapes --json`
- `cd runtime && go test ./...`

Not verified:

- `node tools/vibit check all --json` pending final rerun.
- `cd runtime && go vet ./...` pending final rerun.
- `git diff --check` pending final rerun.

Not applicable:

- Live PostgreSQL verification is not required because no runtime persistence behavior changed.
