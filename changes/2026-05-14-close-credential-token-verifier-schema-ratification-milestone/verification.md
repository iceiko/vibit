# Verification

Verified:

- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check migrations --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change close-credential-token-verifier-schema-ratification-milestone --json`

Pending for final repository state:

- `cd runtime && go test ./...`
- `node tools/vibit check all --json`
- `git diff --check`
- secret scan before commit

Not applicable:

- Live PostgreSQL verification; this closeout adds no migration, adapter implementation, or runtime persistence behavior.
