# Plan

1. Add an application-layer inventory repository provider and transactional inventory handler adapter.
2. Add focused tests for command unit-of-work repository binding and query pass-through behavior.
3. Add PostgreSQL bootstrap composition that uses `postgres.UnitOfWork.NewInventoryRepository`.
4. Add process startup store selection while keeping the in-memory path as the default.
5. Update architecture manifests, module metadata, runtime guides, and runbook documentation.
6. Run runtime and repository verification.
7. Mark `W-0020` complete and promote `W-0021` only if verification passes.
