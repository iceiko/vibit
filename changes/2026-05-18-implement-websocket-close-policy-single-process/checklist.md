# Checklist

- [x] Select `implement_websocket_close_policy_single_process`.
- [x] Add `runtime/internal/app/connection/close_policy.go`.
- [x] Add `runtime/internal/app/connection/close_policy_test.go`.
- [x] Keep close policy application-owned.
- [x] Target only active bound registry records.
- [x] Mark records invalidated without concrete socket close.
- [x] Emit redacted close intents.
- [x] Preserve authentication, transport, protocol, reconnect, generated output, dependency, and direct compatibility deferrals.
- [x] Add ADR, change specs, conversation memory, manifests, AGENTS guidance, and checks.
- [x] Run required verification.
