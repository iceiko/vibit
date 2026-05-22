# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect next`
- `node tools/vibit inspect rule runtime.storage_objects_behavior_gate`
- `node tools/vibit check change define-storage-objects-behavior-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `git diff --check`
- Secret scan for obvious GitHub token patterns excluding `.git/`, `.vibit.local.env`, and `node_modules/`

Known accepted warning:

- `runtime.identity_boundary`

Not applicable:

- Live PostgreSQL verification, because this gate does not modify runtime or persistence behavior.
- Release artifact verification, because this change creates no binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, SDK packages, hosted deployments, public announcements, or paid promotion.
