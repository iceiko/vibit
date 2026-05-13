# Plan

1. Create the `W-0021` change spec.
2. Add a PostgreSQL-owned `*sql.DB` helper for migration runner integration.
3. Add an opt-in live PostgreSQL Protobuf request-loop integration test.
4. Update manifests and bilingual docs with the new verification command and skip semantics.
5. Run Go tests, Go vet, focused skipped/live test output, repository checks, whitespace checks, and secret scan.
6. Mark `W-0021` completed if verification succeeds and live unavailability is explicitly recorded.
7. Do not close `M-002` unless live PostgreSQL verification actually runs or the maintainer confirms closure with unavailable live verification.
