# Checklist

- [x] Add application-owned realtime runtime service.
- [x] Add focused realtime service tests.
- [x] Require validated service/admin sender authority.
- [x] Reject metadata-only and player-authored sender attempts before registry resolution.
- [x] Resolve active bound connection and player-current-connection targets through the connection registry.
- [x] Keep stream subscriber delivery future-only.
- [x] Keep delivery results and errors redacted.
- [x] Avoid WebSocket writes, Protobuf, generated output, protocol bridge, startup wiring, persistence, migrations, dependencies, auth/session changes, route-protection changes, and direct compatibility.
- [x] Add `ADR-0123`.
- [x] Add conversation and change records.
- [x] Update `.arch` manifests and next-ready state.
- [x] Register `runtime.first_server_push_realtime_messaging_runtime_slice`.
- [x] Update docs and module guidance.
- [x] Run verification commands.
- [x] Commit and attempt push.
