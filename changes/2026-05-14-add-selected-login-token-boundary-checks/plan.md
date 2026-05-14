# Plan

1. Add the English selected login/token boundary check standard.
2. Add the Simplified Chinese translation.
3. Add ADR-0030.
4. Add a conversation log.
5. Register `runtime.selected_login_token_boundary` in the rule catalog.
6. Extend `tools/vibit check runtime` with narrow static checks.
7. Register standards, decisions, and rule markers in architecture manifests and agent guides.
8. Mark W-0072 completed and promote W-0073.
9. Run repository verification.
10. Record verification results.

Non-goals:

- Do not implement runtime authentication.
- Do not add login handlers, WebSocket routes, or runtime player handlers.
- Do not add credential, token, external identity, runtime session, or audit schema.
- Do not add repositories or adapters.
- Do not add generated authentication contract shapes.
- Do not add authentication Protobuf wire messages.
- Do not change WebSocket handshake behavior.
- Do not add live service requirements.
