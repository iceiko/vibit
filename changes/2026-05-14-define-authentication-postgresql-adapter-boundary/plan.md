# Plan

1. Add authentication PostgreSQL adapter boundary details to the persistence standard and Chinese translation.
2. Update authentication module manifest and guides with planned adapter source, test path, constructor, and unit-of-work helper.
3. Update architecture manifests to record boundary-defined and adapter-not-added state.
4. Update runtime checks so the boundary is machine-verifiable after the authentication repository interface exists.
5. Mark `W-0081` completed and make `W-0082` next ready.
6. Run focused and full verification.

## Rollback

Remove the boundary markers, persistence-standard section, guide updates, and runtime check changes. No database rollback is needed because no migration source changes.
