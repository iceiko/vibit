# Impact

Runtime impact:

- No Go code is added by the gate-only step.
- The gate authorizes a future implementation to wire existing authentication validation into `runtime/cmd/vibit-server`.

Architecture impact:

- Startup composition ownership is assigned to `runtime/cmd/vibit-server`.
- The first composed runtime store is limited to `VIBIT_RUNTIME_STORE=postgres`.
- The memory runtime store remains a metadata-only bootstrap path.
- Nakama informs the authenticated session/token capability expectation.
- Pitaya informs transport/session/route/handler separation.

Data impact:

- No migrations are added or changed.
- No repository interfaces are changed.
- No PostgreSQL adapters are changed.

Compatibility impact:

- Existing Protobuf envelope remains unchanged.
- WebSocket handshake remains unchanged.
- Existing memory startup remains available.

Security impact:

- PostgreSQL startup must fail closed when verifier key configuration is missing or invalid.
- Raw verifier keys, access tokens, credential proof, lookup digests, verifier digests, and full verifier key ids remain redacted.
