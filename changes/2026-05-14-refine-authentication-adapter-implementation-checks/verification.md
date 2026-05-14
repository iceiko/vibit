# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change refine-authentication-adapter-implementation-checks --json`
- `node tools/vibit inspect next --json`

Pending for final repository state:

- `cd runtime && go test ./...`
- `node tools/vibit check all --json`
- `git diff --check`
- secret scan before commit

Not applicable:

- Live PostgreSQL verification; this change updates static checks and adds no adapter implementation or database behavior.
