# Plan

1. Add `proto/vibit/player/v1/player.proto` aligned with the existing player semantic contracts.
2. Generate Go Protobuf output with `buf generate`.
3. Update `modules/player/module.yaml`.
4. Update `.arch/contracts.yaml`, `.arch/protocol.yaml`, `.arch/runtime.yaml`, and `.arch/work-items.yaml`.
5. Update relevant standards and Simplified Chinese translations.
6. Update `tools/vibit` runtime identity-boundary checks for the new ratified player wire-source phase.
7. Run verification.

## Boundaries

Do not add player runtime handlers, WebSocket routes, authentication, token behavior, credential storage, player account persistence, session persistence, database migrations, Protobuf envelope changes, or WebSocket handshake changes.

Do not hand-edit generated output.
