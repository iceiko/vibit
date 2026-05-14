# Request

## Original Request

Continue the current work queue.

## Clarified Requirement

Advance `W-0023 Add player module manifest boundary` by creating the `player` module manifest and module-level agent guides as boundary artifacts only.

## User-Visible Outcome

Maintainers and future agents can inspect `modules/player/module.yaml` and understand that `player` owns stable player identity and player account lifecycle vocabulary, while authentication providers, token formats, credentials, session storage, player account migrations, Protobuf messages, WebSocket handshake changes, and inventory ownership remain out of scope.

The module registry should now list `player` as a first-class module. Contract and protocol checks should also understand that `player` is currently a boundary-only module with no public contracts yet.

## Non-Goals

- Do not implement authentication.
- Do not choose JWT, OAuth, OIDC, password login, guest login, device login, social login, token format, credential storage, or external identity provider behavior.
- Do not add player account database migrations or persistent schemas.
- Do not add player runtime Go code.
- Do not add player public commands, queries, events, errors, or permissions.
- Do not add player Protobuf messages.
- Do not change the WebSocket handshake contract.
- Do not move inventory state or inventory dependencies under the player module.

## Unknowns

- Concrete authentication scheme remains unselected.
- Token/session persistence model remains unselected.
- Player account database schema remains unselected.
- WebSocket handshake authentication contract remains unselected.
- Public player account/session API remains unselected.

## Acceptance Criteria

- [x] Create `modules/player/module.yaml`.
- [x] Create `modules/player/AGENTS.md`.
- [x] Create `modules/player/AGENTS.zh-CN.md`.
- [x] Register `player` in `.arch/modules.yaml`.
- [x] Record `player` as boundary-only in `.arch/contracts.yaml` without inventing public contracts.
- [x] Update verification tooling so boundary-only modules do not require public contracts or Protobuf sources.
- [x] Complete `W-0023` and expose the next bounded work item.
- [x] Record verification.
