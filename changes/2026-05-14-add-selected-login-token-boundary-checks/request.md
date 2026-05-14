# Request

Add repository checks for the selected login/token boundaries.

The change must enforce the ratified `device_credential_login` and opaque access-token posture as semantic-contract-only until future schema, repository, adapter, generated-output, protocol, test, and runtime implementation gates explicitly authorize concrete behavior.

The change must not implement authentication behavior, add login handlers, add token parsing or validation, add credential/token/session schema, add repositories or adapters, change Protobuf or WebSocket behavior, or require live external services by default.
