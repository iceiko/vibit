# Impact

Runtime impact:

- Adds `runtime/internal/app/authentication/verifier_key_env.go`.
- Adds `runtime/internal/app/authentication/verifier_key_env_test.go`.
- Allows only this application-owned adapter to read the five `VIBIT_AUTH_*` verifier key environment variables.
- Allows only this adapter to decode verifier key text using the Go standard library.
- Keeps the loader unwired from process startup.

Architecture impact:

- Marks `W-0099` and `M-027` completed.
- Opens `M-028` / `W-0100` as the next conservative gate-definition step for token and credential material generation implementation.
- Updates runtime checks so `verifier_key_env.go` is explicitly authorized while other packages remain forbidden from loading verifier keys.

No impact:

- No public API, command, query, event, permission, Protobuf, WebSocket, repository, or migration contract changes.
- No local secret file, dotenv, CLI secret input, KMS, cloud secret-manager, token generation, credential generation, digest computation, verifier comparison, authentication service, login, token validation, logout, cleanup, startup wiring, dependency, or production authentication behavior is added.
