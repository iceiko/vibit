# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0011` by extending the inventory repository boundary so `GrantItem` can lock the player inventory aggregate before reading current items and applying a capacity-sensitive mutation.

## User-Visible Outcome

Maintainers and agents can see, in code and documentation, that the inventory write path follows this order:

- Validate request shape.
- Check grant permission.
- Lock the inventory aggregate for `player_id`.
- Read current inventory through the locked view.
- Enforce capacity policy.
- Apply the grant through the locked view.
- Emit `ItemGranted`.

## Non-Goals

- Do not add PostgreSQL adapter code yet.
- Do not add SQL migrations yet.
- Do not introduce `pgx` or `goose` imports into the inventory module.
- Do not move transaction ownership into inventory repositories.
- Do not replace the accepted inventory account row lock model with advisory locks or another concurrency model.
- Do not add MinIO or object storage.
- Do not change public command, query, event, permission, error, WebSocket, or Protobuf contracts.

## Unknowns

- The PostgreSQL adapter will define the concrete SQL row-lock statement after the first migration exists.
- The transaction boundary skeleton is still deferred to `W-0012`.

## Acceptance Criteria

- [x] Add a module-owned lock boundary for inventory command mutations.
- [x] Make `GrantItem` use the locked mutation view before reading inventory and applying the grant.
- [x] Keep transaction ownership outside inventory repositories.
- [x] Keep third-party persistence dependencies outside the inventory module.
- [x] Update in-memory runtime support and tests.
- [x] Update English guidance and paired Simplified Chinese translations.
