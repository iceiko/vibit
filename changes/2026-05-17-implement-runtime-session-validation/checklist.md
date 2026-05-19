# Checklist

- [x] Added `runtime/internal/app/runtime_session_validator.go`.
- [x] Added `runtime/internal/app/runtime_session_validator_test.go`.
- [x] Kept validation application-owned.
- [x] Used only `runtime/internal/app/session.Repository`.
- [x] Required already validated player identity before repository lookup.
- [x] Used `FindActiveSessionByID` for active persisted session lookup.
- [x] Set `SessionValidated = true` only after durable validation succeeds.
- [x] Collapsed lookup, expiration, revocation, malformed row, mismatch, nil repository, and repository errors to a stable invalid-session public reason.
- [x] Avoided session creation and last-seen updates.
- [x] Kept route policy unchanged.
- [x] Kept WebSocket transport credential-neutral.
- [x] Kept Protobuf envelope and generated output unchanged.
- [x] Kept direct Nakama/Pitaya API compatibility out of scope.
- [x] Ran verification.
