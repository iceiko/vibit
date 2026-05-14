# Request

## Original Request

Continue advancing until no safe work item remains. The current next-ready work item is `W-0032 Ratify minimal player account semantic contracts`.

## Clarified Requirement

Ratify the minimal player account semantic contracts recommended by `docs/player-account-session-contracts.md` without adding runtime code, persistence, authentication, token behavior, credential storage, Protobuf messages, or WebSocket handshake changes.

## User-Visible Outcome

The player module will no longer be only a boundary placeholder for account lifecycle. It will have registered semantic contracts for creating and reading player accounts, plus the first account-created event, error catalog, and permission catalog.

## Non-Goals

- Do not implement Go player runtime code.
- Do not add generated code.
- Do not add Protobuf player messages.
- Do not add player database schema or migrations.
- Do not choose authentication scheme, login method, token format, credential storage, or session persistence.
- Do not change WebSocket handshake behavior.
- Do not copy Nakama or Pitaya public API shape.

## Acceptance Criteria

- Add player account contract manifests under `contracts/player/`.
- Register those contracts in `.arch/contracts.yaml`.
- Update `modules/player/module.yaml`.
- Preserve ask-first boundaries.
- Complete `W-0032`.
- Add a conservative next-ready work item only if safe.
- Run contract, protocol, work, change, and repository checks.
