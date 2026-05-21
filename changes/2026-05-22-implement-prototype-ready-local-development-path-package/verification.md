# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect next`
- `node tools/vibit inspect rule runtime.prototype_ready_local_development_path_package`
- `node tools/vibit check change implement-prototype-ready-local-development-path-package --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `examples/local-alpha-request-loop.sh`
- `git diff --check`
- Secret scan for obvious GitHub token patterns excluding `.git/`, `.vibit.local.env`, and `node_modules/`

Known accepted warning:

- `runtime.identity_boundary`

Not verified:

- Live PostgreSQL verification, because this package does not modify runtime or persistence behavior and default checks must not require live PostgreSQL or private secrets.

Not applicable:

- Release artifact verification, because this change creates no binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, SDK packages, hosted deployments, public announcements, or paid promotion.
