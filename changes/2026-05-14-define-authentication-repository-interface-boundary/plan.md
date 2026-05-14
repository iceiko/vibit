# Plan

1. Add `modules/authentication/` manifest, guides, and standard module directories.
2. Add `runtime/internal/modules/authentication/repository.go` with storage-neutral repository boundary types.
3. Add focused tests in `runtime/internal/modules/authentication/repository_test.go`.
4. Update `tools/vibit` to allow this specific repository boundary while keeping runtime authentication behavior blocked.
5. Update architecture manifests and agent guides so `W-0080` is completed and `W-0081` becomes next ready.
6. Run focused and full verification.

## Rollback

Remove the authentication module manifest, guides, runtime repository package, focused tests, and the related manifest/check updates. No database rollback is needed because no migration source changes.
