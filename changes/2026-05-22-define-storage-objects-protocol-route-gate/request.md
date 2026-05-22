# Request

## Original Request

继续推进。注意要贴合nakama pitaya

## Clarified Requirement

Advance the next-ready work item `W-0210 Define storage objects protocol route gate`. Define a gate-only standard for future storage object protocol routes after application runtime behavior, explicitly aligning the capability with Nakama-style storage objects and Pitaya-style route/session/handler separation, without adding route implementation, Protobuf source, generated output, startup wiring, runtime handlers, or direct compatibility.

## User-Visible Outcome

Maintainers and agents should see:

- `docs/storage-objects-protocol-route-gate.md`
- `docs/storage-objects-protocol-route-gate.zh-CN.md`
- `ADR-0118`
- `runtime.storage_objects_protocol_route_gate`
- updated manifests showing `W-0210` completed and the next step moved to storage object protocol route implementation.

## Non-Goals

- No protocol route implementation.
- No Protobuf source creation.
- No generated output.
- No startup wiring.
- No runtime handlers.
- No repository interface changes.
- No PostgreSQL adapter changes.
- No migrations.
- No dependencies.
- No authentication/session behavior changes.
- No hosted deployments.
- No release artifacts.
- No public announcements or paid promotion.
- No public ACLs, admin search, batch writes, JSON patch, merge semantics, TTL, script hooks, group/guild/party/room/match storage scopes, large object/blob storage, or S3-compatible object storage.
- No direct Nakama/Pitaya API compatibility.

## Unknowns

- Exact Protobuf field numbering and optional-field semantics remain deferred to the implementation slice.
- Whether future route implementation uses generated contract-shape helpers remains deferred.
- Live PostgreSQL verification remains optional and not required for this gate.

## Acceptance Criteria

- [x] English and Simplified Chinese gate documents are added.
- [x] ADR-0118 records the decision.
- [x] Nakama and Pitaya reference mapping is explicit.
- [x] Candidate route names and message shapes are recorded.
- [x] Protected-route policy and metadata-only identity refusal are recorded.
- [x] Stop conditions preserve route implementation, Protobuf/generated output, startup wiring, repository/adapter/migration/dependency/authentication/session, release, blob/S3, and direct compatibility deferrals.
- [x] Repository manifests and check rules recognize the gate.
