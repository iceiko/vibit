# Impact

## Affected Modules

- `player`: future token validation may bind authenticated proof to player accounts after later gates.
- `runtime`: future token validation must produce request identity before domain dispatch.

No module implementation changes are made.

## Module Ownership Impact

No ownership changes.

The comparison reinforces that:

- Tokens are not player account lifecycle rows.
- Runtime application validation owns request identity handoff.
- WebSocket transport remains credential-neutral.
- Protobuf envelope metadata remains metadata-only.

## Public Contract Impact

No contracts are added, changed, or removed.

The comparison recommends future token issuance and validation contract work.

## Event Impact

No event changes.

Future token or session audit events remain deferred.

## Permission Impact

No permission changes.

Future token validation, refresh, logout, and revocation permissions remain deferred.

## Data And Migration Impact

No data changes and no migrations.

The recommendation implies future token storage or verifier schema gates, but does not add storage.

## Protocol Impact

No Protobuf source or generated output changes.

The recommendation explicitly defers Protobuf envelope extension and rejects current `Session` metadata as proof.

## Test Impact

No runtime tests are required because no runtime code changes.

## Documentation Impact

Adds:

- `docs/token-format-carrier-options.md`
- `docs/token-format-carrier-options.zh-CN.md`
- `conversations/2026-05-14-token-format-carrier-comparison.md`

Updates:

- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Compatibility Risks

No API, event, data, Protobuf, or WebSocket compatibility risk in this change.

The main risk is future agents treating the recommendation as implementation permission. The document explicitly forbids that.
