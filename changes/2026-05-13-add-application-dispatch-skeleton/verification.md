# Verification

Verified:

- `node --check tools/vibit`
- `go test ./...` from `runtime/`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change add-application-dispatch-skeleton --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for GitHub tokens outside ignored local files:
  `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .`

Not verified:

- `go vet ./...` has not yet been run for this change.

Not applicable:

- Protobuf generation is not applicable because this change does not add or change `.proto` sources.
- Database migration verification is not applicable because no persistence or migrations are added.
