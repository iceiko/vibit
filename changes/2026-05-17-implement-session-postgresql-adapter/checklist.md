# Checklist

- [x] Added `runtime/internal/platform/persistence/postgres/session_repository.go`.
- [x] Added `runtime/internal/platform/persistence/postgres/session_repository_test.go`.
- [x] Added `UnitOfWork.NewSessionRepository()`.
- [x] Kept the adapter persistence-only.
- [x] Kept token proof validation out of the adapter.
- [x] Kept `RequestIdentity.SessionValidated` false.
- [x] Kept WebSocket transport credential-neutral.
- [x] Kept Protobuf envelope and generated output unchanged.
- [x] Kept direct Nakama/Pitaya API compatibility out of scope.
- [x] Ran verification.
