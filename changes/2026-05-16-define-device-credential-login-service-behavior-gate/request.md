# Request

Define the device credential login service behavior gate as the next bounded work item after the authentication service skeleton.

The gate must not implement login execution. It should make the future login sequence, proof shape, repository handoff, helper ordering, token issuance posture, public error collapse, redaction rules, and required tests explicit before code is allowed to execute real login behavior.

## Non-Goals

- Implementing device credential login.
- Issuing access tokens.
- Validating access tokens.
- Changing service method signatures.
- Changing repository interfaces, PostgreSQL adapters, migrations, Protobuf messages, WebSocket carriers, startup wiring, generated files, or dependencies.
- Adding broader production authentication behavior.
