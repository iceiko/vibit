# Plan

1. Add focused tests for storage-neutral repository vocabulary, input normalization, JSON value copying, redacted errors, and forbidden material.
2. Implement `runtime/internal/modules/storage.Repository`, value types, conflict types, repository errors, and normalizers.
3. Add the storage module manifest and paired AGENTS guides.
4. Record `ADR-0113` and the conversation log.
5. Update `.arch/work-items.yaml` so `W-0205` is completed and `W-0206` becomes next-ready.
6. Update architecture manifests, README, alpha, product maturity, roadmap, and agent guide pointers.
7. Add `runtime.storage_objects_repository_interface_implementation` to `rules/check-rules.json` and `tools/vibit`.
8. Run static checks, Go tests, and diff hygiene checks.

## Generated Artifacts

None.

## Handwritten Logic

Add storage-neutral Go value types, repository interface methods, redacted errors, and normalization helpers. The code must not import PostgreSQL, WebSocket, generated Protobuf, platform adapters, S3 SDKs, or MinIO clients.

## Verification Commands

- `go test ./internal/modules/storage`
- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.storage_objects_repository_interface_implementation`
- `node tools/vibit check change implement-storage-objects-repository-interface --json`
- `node tools/vibit check module storage --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `git diff --check`

## Rollback Or Migration Notes

Reversal means deleting `runtime/internal/modules/storage`, removing the storage module manifest and guides, removing ADR-0113 and the `W-0205` manifest/check entries, and reopening `M-133/W-0205`. No data migration rollback is required because this change adds no migration or SQL execution behavior.
