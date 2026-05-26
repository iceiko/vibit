# Plan

1. Read the repository constitution, module instructions, friends lifecycle/persistence/migration/repository boundary artifacts, and current next work item.
2. Add `runtime/internal/modules/friends` with storage-neutral repository value types, lifecycle/status/conflict vocabulary, normalizers, redacted errors, and the repository interface.
3. Add focused tests for storage neutrality, closed vocabularies, canonical pair normalization, actor handoff, version validation/copying, result copying, redacted errors, and forbidden fields.
4. Add `modules/friends` manifest and module AGENTS guides.
5. Add change spec, ADR, conversation memory, check rule, and architecture manifest updates.
6. Advance `.arch/work-items.yaml` from `W-0235` to `W-0236 Define friends relationship PostgreSQL adapter gate`.
7. Run focused and full verification.
8. Commit and push if verification passes.
