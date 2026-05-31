# Request

## Original Request

```text
继续推进20步，向pitaya靠拢
```

## Clarified Requirement

Advance the active next-ready work item, `M-173/W-0245 Implement minimum operations inspection source-first surface`, and use the slice to move toward Pitaya by recording Pitaya architecture vocabulary as deferred, source-first inspection context. Do not bypass `ADR-0152` by adding distributed runtime behavior.

## User-Visible Outcome

Maintainers and agents can run:

```bash
node tools/vibit inspect operations --json
```

The command reports the current local alpha operations posture, redaction posture, route families, verification posture, and Pitaya deferred architecture mapping.

## Non-Goals

- No operations/admin endpoint behavior.
- No metrics endpoint.
- No observability pipeline.
- No dashboard.
- No runtime behavior change.
- No protocol route or Protobuf source.
- No generated output.
- No repository interface or PostgreSQL adapter change.
- No migration.
- No dependency.
- No authentication/session behavior change.
- No event/audit table.
- No SDK publication.
- No hosted deployment.
- No distributed runtime implementation.
- No direct Nakama/Pitaya API compatibility.

## Unknowns

None for this slice. Future Pitaya vocabulary reactivation details are left to W-0246.

## Acceptance Criteria

- [x] `node tools/vibit inspect operations --json` emits JSON.
- [x] The output includes the ten `ADR-0152` minimum inspectable categories.
- [x] The output includes existing local alpha status endpoints and source surfaces.
- [x] The output includes current route families without exposing credentials or sensitive identifiers.
- [x] The output includes redaction flags for credential, token, digest, key, DSN, transport metadata, sensitive identifier, database payload, and local secret classes.
- [x] The output includes Pitaya deferred architecture mapping.
- [x] Repository checks register and verify `runtime.minimum_operations_inspection_source_first_surface_implementation`.
- [x] W-0245 is completed and W-0246 is next-ready.
