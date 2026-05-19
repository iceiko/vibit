# Request

## Original Request

The maintainer asked, in Chinese, for the agent to recommend the next ten steps and then execute those ten steps. The selected direction was:

```text
add_authentication_command_protocol_messages_and_login_route_registration
```

The maintainer also required strong reference to Nakama and Pitaya.

## Clarified Requirement

Implement the bounded slice authorized by `ADR-0055` and `docs/authentication-command-protocol-login-route-gate.md`:

- Add the public device credential login Protobuf command messages.
- Generate Go Protobuf output through Buf.
- Bridge the Protobuf login request/response to the existing application authentication service types.
- Register the explicit public login route in application bootstrap.
- Wire the route only in PostgreSQL startup composition where the authentication service exists.
- Bypass the outer transactional dispatcher for the login route because the authentication service owns its own unit-of-work.

## User-Visible Outcome

The PostgreSQL runtime path can now accept a public `runtime.authentication.AuthenticateWithDeviceCredential` command over the existing WebSocket-framed Protobuf envelope and return the existing service's opaque access token in the command response.

## Non-Goals

- Do not add session persistence.
- Do not add WebSocket handshake authentication.
- Do not add HTTP `Authorization`, Bearer, cookie, query-string, or WebSocket subprotocol credential carriers.
- Do not change repository interfaces, PostgreSQL adapters, migrations, or dependencies.
- Do not add logout, refresh, cleanup, token rotation, token validation audit mutation, or previous-token revocation.
- Do not add memory durable authentication behavior.
- Do not add direct Nakama/Pitaya public API compatibility.

## Acceptance Criteria

- [x] `proto/vibit/authentication/v1/authentication.proto` defines the login request and response messages.
- [x] Generated Go Protobuf output is produced through Buf.
- [x] Protocol bridge maps login request/response payloads without changing the existing envelope.
- [x] Application bootstrap registers the explicit public login route.
- [x] PostgreSQL startup composition registers the login route only when the authentication service is composed.
- [x] The public login route bypasses the outer transactional dispatcher.
- [x] Login errors do not leak credential proof or access-token text.
- [x] Protected-route behavior is preserved.
- [x] WebSocket transport remains credential-neutral.
- [x] Verification is recorded.
