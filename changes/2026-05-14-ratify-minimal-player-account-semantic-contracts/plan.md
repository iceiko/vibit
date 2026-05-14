# Plan

1. Add player account command, query, event, error, and permission contract manifests.
2. Register the contracts in `.arch/contracts.yaml`.
3. Update `modules/player/module.yaml` so the module declares and references the new contracts.
4. Update the player account/session standard to record ratified minimal account contracts.
5. Complete `W-0032` and add the next conservative work item.
6. Run repository verification.

## Boundaries

Do not add runtime implementation, generated code, Protobuf source, migrations, authentication, token behavior, credential storage, session persistence, or WebSocket handshake changes.
