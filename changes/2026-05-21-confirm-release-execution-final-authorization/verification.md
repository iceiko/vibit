# Verification

Status: Verified

## Required Commands

```bash
git tag --list 'v0.1*'
git ls-remote --tags origin 'refs/tags/v0.1*'
GitHub release lookup for v0.1.0-alpha.1
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change confirm-release-execution-final-authorization --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
git status --short --branch
```

## Results

- `git tag --list 'v0.1*'` returned no local matching tags.
- `git ls-remote --tags origin 'refs/tags/v0.1*'` returned no remote matching tags.
- GitHub release lookup for `v0.1.0-alpha.1` returned HTTP `404`, so no GitHub release record was found before execution.
- `node -c tools/vibit` passed.
- `node tools/vibit inspect next` passed and reports `M-124/W-0196 Define first alpha user discovery loop` as next-ready.
- `node tools/vibit check change confirm-release-execution-final-authorization --json` passed.
- `node tools/vibit check work --json` passed.
- `node tools/vibit check runtime --json` passed with one pre-existing warning: `runtime.identity_boundary` for `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- `node tools/vibit check memory --json` passed.
- `node tools/vibit check schemas --json` passed.
- `node tools/vibit check all --json` passed with the same pre-existing `runtime.identity_boundary` warning.
- `cd runtime && go test ./cmd/vibit-server` passed.
- `cd runtime && go test ./...` passed.
- `examples/local-alpha-request-loop.sh` passed.
- `git diff --check` passed.
- `git status --short --branch` showed expected local changes before commit and no unexpected branch divergence.

## Release Execution Results

- Main push succeeded for commit `30e868803796e261651069eeafb8b5e526ac2260`.
- Annotated Git tag `v0.1.0-alpha.1` was created and pushed.
- Remote tag `refs/tags/v0.1.0-alpha.1^{}` resolves to `30e868803796e261651069eeafb8b5e526ac2260`.
- GitHub Release `v0.1.0-alpha.1` was created with id `327144724`.
- GitHub Release URL: `https://github.com/iceiko/vibit/releases/tag/v0.1.0-alpha.1`.
- GitHub source archive is available through the release/tag page.
