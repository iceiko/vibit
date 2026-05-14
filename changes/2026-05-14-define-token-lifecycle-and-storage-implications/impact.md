# Impact

## Affected Modules

- `player`: future login and token issuance depends on player account lifecycle state, but lifecycle tables remain credential-free and token-free.
- `runtime`: future token validation must run through application-owned validation before domain dispatch.

No module implementation changes are made.

## Module Ownership Impact

No ownership changes.

The change reinforces that:

- Token lifecycle behavior is not owned by WebSocket transport or Protobuf adapters.
- Domain modules do not parse, validate, store, revoke, rotate, clean up, or audit tokens.
- Player account lifecycle persistence does not store credentials, token verifiers, external identity links, runtime sessions, or WebSocket state.

## Public Contract Impact

No contracts are added, changed, or removed.

W-0070 will define future semantic contract, error, permission, and audit surfaces.

## Event Impact

No event changes.

The change defines audit implications for future work without adding audit events.

## Permission Impact

No permission changes.

Future login, logout, token validation, token revocation, cleanup, and administrative lifecycle permissions remain deferred to W-0070.

## Data And Migration Impact

No data changes and no migrations.

The change defines future credential and token verifier storage gates and preserves external identity and session storage deferral.

## Test Impact

No runtime tests are required because no runtime code changes.

Future implementation must add focused tests for expiration, revocation, logout, rotation, replay, redaction, cleanup, and ownership boundaries.

## Documentation Impact

Adds:

- `docs/token-lifecycle-storage-implications.md`
- `docs/token-lifecycle-storage-implications.zh-CN.md`
- `decisions/ADR-0027-token-lifecycle-and-storage-implications.md`
- `conversations/2026-05-14-token-lifecycle-storage-implications.md`

Updates:

- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Compatibility Risks

No API, event, wire, or data compatibility risk in this change.

The main risk is that a future agent treats lifecycle posture as implementation permission. The standard and manifests explicitly forbid implementation until later gates are complete.
