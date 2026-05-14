# Plan

1. Read current work queue, protocol, runtime, module, and reference alignment artifacts.
2. Define player identity, account, authentication, runtime session, transport connection, and request identity context vocabulary.
3. Add the English player identity/session boundary standard.
4. Add the Simplified Chinese translation.
5. Add an ADR recording the durable decision.
6. Update architecture manifests and inventory guidance.
7. Complete `W-0022` and add the next bounded work item.
8. Run repository verification and record results.

## Ask-First Boundaries

Stop and ask before:

- Choosing a concrete authentication scheme.
- Choosing a token format.
- Adding JWT, OAuth, OIDC, password hashing, credential storage, or external identity provider dependencies.
- Adding player account migrations.
- Changing the Protobuf envelope.
- Changing the WebSocket handshake contract.
- Moving inventory ownership into the player module.
