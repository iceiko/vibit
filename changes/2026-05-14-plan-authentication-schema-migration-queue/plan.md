# Plan

1. Inspect the current M-014 milestone and next-ready work item.
2. Add the English authentication schema migration queue standard.
3. Add the Simplified Chinese translation.
4. Record ADR-0034.
5. Record the conversation log.
6. Update architecture manifests and agent guides.
7. Mark W-0076 completed and add W-0077 through W-0082 in dependency order.
8. Run focused repository checks and full verification.
9. Commit and push if verification passes.

## Resulting Queue Shape

- `W-0076`: Plan authentication schema migration queue.
- `W-0077`: Add credential PostgreSQL migration source.
- `W-0078`: Add token verifier PostgreSQL migration source.
- `W-0079`: Add authentication migration static checks.
- `W-0080`: Define authentication repository interface boundary.
- `W-0081`: Define authentication PostgreSQL adapter boundary.
- `W-0082`: Close credential and token verifier schema ratification milestone.

This order keeps schema source, static checks, repository interfaces, and adapter boundaries separate. Runtime authentication implementation remains outside M-014.
