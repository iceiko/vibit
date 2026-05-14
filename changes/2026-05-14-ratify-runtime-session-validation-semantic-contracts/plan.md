# Plan

1. Add runtime session semantic contract files under `contracts/runtime/session/`.
2. Register the runtime session contracts in `.arch/contracts.yaml`.
3. Update `.arch/runtime.yaml` with the ratified semantic contract state.
4. Update `docs/player-account-session-contracts.md` and `docs/player-account-session-contracts.zh-CN.md`.
5. Update `tools/vibit` so `check contracts` validates runtime session contracts without requiring a domain module manifest.
6. Update `.arch/work-items.yaml` to complete `W-0034`.
7. Record verification.

## Boundaries

Do not add authentication, token behavior, credential storage, session persistence, database migrations, Protobuf envelope changes, WebSocket handshake changes, player runtime handlers, or generated output.
