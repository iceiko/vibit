# Plan

## Implementation

- Add `runtime/internal/app/realtime/service.go`.
- Add `runtime/internal/app/realtime/service_test.go`.
- Accept `ADR-0123`.
- Record the conversation and change spec.
- Update `.arch` manifests so `W-0215` is completed and `W-0216` is next-ready.
- Register `runtime.first_server_push_realtime_messaging_runtime_slice`.
- Update `tools/vibit` to check the new runtime slice and allow the realtime service after `W-0215` completes.
- Refresh continuation docs and storage-module guidance.

## Runtime Behavior

- Validate message intent kind, target kind, payload type, payload bytes, and sender authority before recipient resolution.
- Resolve `connection_id_and_epoch` and `player_current_connections` against active bound connection registry records.
- Keep `stream_subscribers` valid vocabulary but future-only.
- Return copied, redacted delivery intents without transport sockets or proof material.

## Non-Goals

- No WebSocket outbound writer.
- No socket writes.
- No `proto/vibit/realtime/v1/realtime.proto`.
- No generated Go Protobuf output.
- No protocol bridge.
- No application bootstrap handler.
- No startup wiring.
- No persistence, migrations, dependencies, auth/session behavior changes, route-protection changes, delivery guarantees, offline inboxes, stream subscription persistence, matchmaking, match runtime, broad social behavior, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Verification

- Format the new Go files.
- Run focused realtime tests.
- Run all Go runtime tests.
- Run `tools/vibit` syntax, inspect, change, work, runtime, memory, schema, and all checks.
- Run `git diff --check`.
