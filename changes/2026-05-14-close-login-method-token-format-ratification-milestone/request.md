# Request

Close `M-013 Login Method And Token Format Ratification` after verifying its completion criteria.

The closeout must preserve implementation deferral and create the next bounded gate. It must not start production authentication, add runtime player handlers, add WebSocket routes, change Protobuf envelope behavior, change WebSocket handshake behavior, add migrations, or add new major dependencies.

## Completion Review

`M-013` completion criteria are satisfied:

- First login methods were compared against Nakama capability coverage and vibit security/agent-native constraints in `docs/first-login-method-candidates.md`.
- The first login-method set was ratified as `device_credential_login` in `docs/first-login-method-set.md` and `ADR-0025`.
- Token format and carrier options were compared in `docs/token-format-carrier-options.md`.
- The first token posture was ratified as opaque high-entropy access token, login-command token issuance, and explicit request proof payload in `docs/first-token-format-proof-carrier-posture.md` and `ADR-0026`.
- Token lifecycle and storage implications were recorded in `docs/token-lifecycle-storage-implications.md` and `ADR-0027`.
- Authentication contract, error, permission, and audit surfaces were defined in `docs/authentication-contract-error-permission-surfaces.md` and `ADR-0028`.
- Credential, token verifier, external identity, runtime session, and audit schema gates were defined in `docs/credential-token-session-schema-gates.md` and `ADR-0029`.
- Selected login/token boundary checks were added in `docs/selected-login-token-boundary-checks.md`, `ADR-0030`, and `runtime.selected_login_token_boundary`.
- No implementation boundary was crossed accidentally.

## Next Gate

The next milestone is `M-014 Credential And Token Verifier Schema Ratification`.

The first next work item is `W-0074 Define credential record schema boundary`.

This next gate continues the selected posture by ratifying persistent schema before migrations, repositories, adapters, runtime lookup, handlers, routes, Protobuf changes, WebSocket changes, generated authentication shapes, or authentication implementation.
