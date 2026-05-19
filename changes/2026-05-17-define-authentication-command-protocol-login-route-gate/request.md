# Request

## Original Request

The maintainer asked the agent to recommend the next ten steps and continue according to that recommendation.

The selected next direction is:

```text
add_authentication_command_protocol_messages_and_login_route_registration
```

## Clarified Requirement

Define a gate for public authentication command protocol messages and device credential login route registration before implementation.

The gate must specify how `AuthenticateWithDeviceCredential` may be exposed through Protobuf payloads, application route registration, the protocol bridge, and PostgreSQL startup composition while preserving transport neutrality and session deferrals.

## User-Visible Outcome

The repository has a written standard and ADR that authorize a bounded implementation slice for public device credential login route exposure.

## Non-Goals

- Do not implement protocol messages in the gate-only step.
- Do not generate Protobuf output in the gate-only step.
- Do not register the login route in the gate-only step.
- Do not add session persistence.
- Do not add WebSocket handshake authentication.
- Do not add HTTP `Authorization`, Bearer, cookie, query string, or WebSocket subprotocol carriers.
- Do not change repository interfaces, PostgreSQL adapters, migrations, or dependencies.
- Do not add logout, refresh, cleanup, token rotation, token validation audit mutation, or broader production authentication behavior.

## Acceptance Criteria

- [x] English and Simplified Chinese gate docs exist.
- [x] `ADR-0055` records the authentication command protocol and login route decision.
- [x] Manifests record the gate, selected public login route, planned protocol source, planned generated output, startup/composition boundaries, and deferrals.
- [x] Repository check rule `runtime.authentication_command_protocol_login_route_gate` exists.
- [x] `W-0119` is completed and `W-0120` is available for implementation.
- [x] Verification is recorded.
