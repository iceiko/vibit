# Verification

Status: Verified

## Commands

```bash
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change add-health-readiness-version-config-surface --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Results

- `cd runtime && go test ./cmd/vibit-server`: passed.
- `cd runtime && go test ./...`: passed.
- `node -c tools/vibit`: passed.
- `node tools/vibit inspect next`: passed and reports `M-116/W-0188 Add alpha acceptance checklist` as the next ready work item.
- `node tools/vibit check change add-health-readiness-version-config-surface --json`: passed.
- `node tools/vibit check work --json`: passed.
- `node tools/vibit check runtime --json`: passed with the known pre-existing `runtime.identity_boundary` warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go` mentioning credential dependency behind a ratified boundary.
- `node tools/vibit check memory --json`: passed.
- `node tools/vibit check schemas --json`: passed.
- `node tools/vibit check all --json`: passed with the same known `runtime.identity_boundary` warning.
- `git diff --check`: passed.
