# Plan

1. Confirm `W-0237` is next-ready.
2. Read the adapter gate, friends repository interface, friends migration source, and existing PostgreSQL adapter patterns.
3. Establish RED by adding focused PostgreSQL adapter tests before implementation.
4. Implement `FriendRelationshipRepository` under `runtime/internal/platform/persistence/postgres`.
5. Add unit-of-work repository handoff without startup/runtime/protocol wiring.
6. Register the check rule `runtime.friends_relationship_postgresql_adapter_implementation`.
7. Add `ADR-0145` and conversation log.
8. Add change spec files.
9. Update work item, architecture, module, public next-work, and agent-guide manifests to complete `W-0237` and open `W-0238`.
10. Run verification commands and record results.
