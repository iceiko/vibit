# Request

## Original Request

```text
继续推进
```

## Current Planning Context

The maintainer has clarified that vibit should actively reference Nakama and Pitaya for game server capability planning, while preserving vibit's Agent-Native differentiator.

The near-term sequence recorded in `docs/reference-game-server-alignment.md` says the next runtime work should define inventory repository and policy interfaces, then implement command/query handlers through the application dispatcher, before PostgreSQL persistence and WebSocket transport.

## Clarified Requirement

Add the first handwritten inventory runtime boundary in Go.

This change should:

- Keep inventory business behavior under `runtime/internal/modules/inventory/`.
- Define vibit-owned repository and policy interfaces before PostgreSQL or WebSocket implementation.
- Implement the first `GrantItem` command handler and `GetInventory` query handler in terms of those interfaces.
- Provide an application dispatch registration function so app dispatch can route inventory command/query requests without knowing inventory internals.
- Add focused tests for successful grants, invalid quantity, capacity rejection, permission rejection, event emission, query immutability, and dispatcher integration.

## User-Visible Outcome

Future agents can see where inventory behavior belongs and how it connects to application dispatch.

The project moves from protocol/app skeleton toward the first complete backend slice without prematurely adding WebSocket transport or PostgreSQL adapters.

## Non-Goals

- Do not add PostgreSQL persistence yet.
- Do not add WebSocket transport yet.
- Do not add new third-party dependencies.
- Do not modify generated Protobuf files.
- Do not implement player/account/session ownership.
- Do not implement item catalog validation.
- Do not add distributed runtime, RPC, groups, or cluster work.

## Unknowns

- Exact transaction/unit-of-work interface remains deferred until persistence adapter work.
- Exact Protobuf-to-domain payload mapping remains deferred to a protocol adapter bridge change.
- Generated Go contract shapes remain planned but not implemented in this change.
- Persistence error taxonomy remains deferred until the PostgreSQL repository adapter exists.

## Acceptance Criteria

- [ ] Inventory runtime package exists under `runtime/internal/modules/inventory/`.
- [ ] Repository and policy interfaces are declared inside the inventory module boundary.
- [ ] `GrantItem` validates positive quantity, checks permission, checks capacity, mutates repository state, and emits exactly one `ItemGranted` event on success.
- [ ] `GetInventory` checks permission and reads without mutation.
- [ ] Inventory routes can be registered with `runtime/internal/app.Dispatcher`.
- [ ] Tests cover command, query, event, permission, capacity, and dispatcher integration behavior.
- [ ] Repository, runtime, and all checks pass.
