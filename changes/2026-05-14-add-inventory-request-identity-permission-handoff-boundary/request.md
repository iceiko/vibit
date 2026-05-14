# Request

## Original Request

Advance `W-0027`: add the inventory request identity permission handoff boundary.

## Clarified Requirement

Inventory permission policies must be able to inspect the application-owned request identity context without inventory owning authentication, player accounts, sessions, tokens, credentials, or the player module.

## User-Visible Outcome

Agents can see where request identity reaches inventory permission decisions, and tests prove metadata-only identity is not treated as authenticated proof for privileged grants.

## Non-Goals

- No production authentication implementation.
- No token parsing, credential lookup, session persistence, or player account lookup.
- No player module dependency from inventory.
- No Protobuf envelope, WebSocket handshake, command, query, or event contract shape changes.
- No service/admin authorization model beyond placeholder-safe context fields.

## Unknowns

- The production authorization model for service, admin, player self-read, and privileged grants remains deferred.
- The concrete authentication scheme, token format, session store, and player account schema remain deferred.

## Acceptance Criteria

- [x] Inventory permission policy receives request identity context through a vibit-owned runtime type.
- [x] Existing bootstrap static permission behavior remains explicit and unchanged for current local/runtime paths.
- [x] Metadata-only identity is not treated as authenticated proof for privileged grants by the new identity-aware policy boundary.
- [x] Inventory remains independent from the player module.
- [x] Focused Go tests cover permission handoff and metadata-only denial/neutral semantics.
