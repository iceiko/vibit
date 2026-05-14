# Impact

## Affected Modules

`player` is affected because the standard clarifies that player account lifecycle storage must not contain credentials, provider subjects, tokens, sessions, WebSocket state, or request validation results.

`runtime` is affected because future authentication and session validation work must treat credential storage and external identity linking as separate boundaries.

## Module Ownership Impact

No ownership moves.

The player module remains the owner of stable `player_id` and account lifecycle. Future credential storage and external identity linking owners remain planned and unimplemented.

## Public Contract Impact

No command, query, event, error, or permission contract is added or changed.

The change defines future artifact gates for login method contracts, credential schema boundaries, provider namespace contracts, provider subject contracts, link/unlink contracts, audit events, errors, permissions, and verification before implementation.

## Data And Migration Impact

No migration is added.

The change explicitly preserves `player_accounts` and `player_account_events` as account lifecycle tables only.

## Protocol Impact

No Protobuf source, generated Protobuf output, envelope field, message kind, WebSocket route, or WebSocket handshake behavior is changed.

## Reference Alignment Impact

Nakama multiple authentication methods are mapped into deferred vibit login-method families.

Pitaya session binding and session context vocabulary are retained as architecture vocabulary, but no direct public API compatibility is adopted.

## Compatibility Risks

This is a boundary standard and manifest update. It should not break runtime behavior.

The main risk is implying that listed login-method families are selected. The standard labels them as deferred capability coverage only.

## Test Impact

No Go runtime code changes are required.

Repository verification should cover runtime boundary checks, architecture checks, work queue consistency, memory/change-spec checks, player module checks, and the full repository check.
