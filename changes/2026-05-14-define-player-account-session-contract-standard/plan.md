# Plan

1. Add `docs/player-account-session-contracts.md`.
2. Add `docs/player-account-session-contracts.zh-CN.md`.
3. Update runtime/reference manifests to point to the new standard.
4. Complete `W-0031` and add the next conservative `next_ready` item.
5. Verify the repository.

## Boundaries

This change must not implement authentication, choose token behavior, add credentials, add migrations, add session persistence, change Protobuf, change WebSocket handshake behavior, or copy Nakama/Pitaya API shape.
