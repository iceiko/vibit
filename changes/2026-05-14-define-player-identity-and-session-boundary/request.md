# Request

## Original Request

Continue the current work queue after completing the durable inventory runtime and planning the next milestone.

## Clarified Requirement

Advance `W-0022 Define player identity and session boundary` by defining ownership between player identity, player accounts, authentication, WebSocket sessions, and request identity context before implementing authentication or player account persistence.

## User-Visible Outcome

Maintainers and future agents can see where player identity, account lifecycle, authentication, runtime session validation, transport connection metadata, request identity context, and inventory `player_id` references belong.

The repository should also expose the next bounded work item for creating the player module boundary without silently choosing a token format, authentication scheme, credential store, database migration, or WebSocket handshake contract.

## Non-Goals

- Do not implement authentication.
- Do not choose JWT, OAuth, OIDC, guest login, password login, device login, social login, or any external identity provider.
- Do not add player account database migrations.
- Do not change the Protobuf envelope shape.
- Do not change the WebSocket handshake contract.
- Do not move inventory ownership into the player module.
- Do not replace inventory permission behavior in this boundary-only step.

## Unknowns

- The concrete authentication scheme remains unselected.
- The token/session persistence model remains unselected.
- The player account database schema remains unselected.
- The WebSocket handshake authentication contract remains unselected.

## Acceptance Criteria

- [x] Add an English player identity and session boundary standard.
- [x] Add a Simplified Chinese translation for the standard.
- [x] Add an ADR recording the durable boundary decision.
- [x] Update architecture manifests so agents can find the standard.
- [x] Update inventory guidance to state how inventory consumes `player_id` without owning player accounts.
- [x] Update the work queue so `W-0022` is completed and the next bounded work item is visible.
- [x] Record verification.
