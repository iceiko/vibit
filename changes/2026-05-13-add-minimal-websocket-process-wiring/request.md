# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0008` by adding minimal runtime process wiring for the first WebSocket Protobuf application request loop. The runtime process should mount the composed WebSocket handler at `/v1/ws`, keep business logic out of `cmd/vibit-server`, and provide a narrow verification path.

## User-Visible Outcome

Maintainers can start the first Go runtime server and reach the planned gameplay WebSocket endpoint:

```bash
cd runtime
go run ./cmd/vibit-server
```

The process mounts:

```text
/v1/ws
```

## Non-Goals

- Do not add authentication or session validation.
- Do not add PostgreSQL persistence or transaction wiring.
- Do not add MinIO or object storage.
- Do not add generated route registration.
- Do not change the Protobuf envelope shape.
- Do not add a broad REST API surface.

## Unknowns

- None for this bounded step.

## Acceptance Criteria

- [x] Mount the composed WebSocket handler at `/v1/ws` from runtime process startup.
- [x] Keep business logic out of `cmd/vibit-server`.
- [x] Reuse the existing transport, Protobuf composition, app dispatcher, and inventory handler boundaries.
- [x] Add a narrow integration test or documented manual run path.
- [x] Preserve existing import and layer boundaries.
