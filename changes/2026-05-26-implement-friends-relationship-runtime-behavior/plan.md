# Plan

1. Read the friends runtime behavior gate, repository interface, PostgreSQL adapter handoff, application identity shape, and unit-of-work patterns.
2. Add focused failing tests for the application-owned friends service.
3. Implement `runtime/internal/app/friends/service.go` with typed requests, results, statuses, errors, identity derivation, unit-of-work repository handoff, conflict mapping, and actor-relative status conversion.
4. Keep protocol routes, Protobuf source, generated output, dependencies, migrations, startup wiring, event/audit tables, and broad social scope out of the implementation.
5. Add ADR-0147, conversation memory, and completion records.
6. Update architecture manifests, module guide/manifest, README/alpha docs, rules, and tooling checks so W-0239 is recorded as completed.
7. Run focused Go tests, repository checks, and whitespace verification.
