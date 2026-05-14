# Impact

## Architecture Impact

This change advances the M-014 authentication schema ratification queue by adding the second planned authentication migration source.

The migration source is owned by `runtime.authentication` and keeps authentication storage separate from player account lifecycle storage. It references `player_accounts(player_id)` and `authentication_device_credentials(credential_record_id)` without modifying those tables.

## Runtime Impact

No Go runtime behavior changes.

The change does not add token generation, token validation, logout, refresh, cleanup, credential lookup, handlers, routes, WebSocket proof carriers, Protobuf messages, generated authentication shapes, authentication dependencies, or production authentication behavior.

## Data Impact

The source-only migration creates one future table when applied:

- `authentication_access_tokens`

The table stores non-plaintext token lookup and verifier digests only. It does not store raw access tokens, refresh tokens, JWT claims, session records, WebSocket state, generic metadata, provider payloads, authorization headers, cookies, or credential proof material.

Live PostgreSQL apply/rollback remains opt-in through the existing disposable DSN workflow.

## Agent Impact

Agents now have a concrete SQL target for token verifier records before repository and adapter work begins.

The change also updates static checks so later agents can distinguish the authorized token verifier migration from forbidden runtime authentication, token behavior, session storage, external identity storage, or audit persistence work.

## Reference Alignment

Nakama and Pitaya remain references for game backend authentication/session capability pressure and vocabulary only.

This change does not copy their public API shape, session model, token carrier model, route model, or storage layout.
