# Impact

## Affected Modules

- `player`: future selected login methods may create or authenticate player accounts.
- `runtime`: future authentication proof must become runtime-owned request identity before domain dispatch.

No module implementation changes are made.

## Module Ownership Impact

No ownership changes.

The comparison reinforces that:

- Player account lifecycle does not own credentials or tokens.
- Runtime application dispatch owns future validation handoff.
- WebSocket transport remains credential-neutral.
- Protobuf envelope metadata remains metadata-only.

## Public Contract Impact

No contracts are added, changed, or removed.

The comparison recommends future contract work for device credential login.

## Event Impact

No event changes.

## Permission Impact

No permission changes.

The comparison recommends future permission catalog entries before implementation.

## Data And Migration Impact

No data changes and no migrations.

The recommendation implies future credential storage schema work, but does not add it.

## Test Impact

No runtime tests are required because no runtime code changes.

## Documentation Impact

Adds:

- `docs/first-login-method-candidates.md`
- `docs/first-login-method-candidates.zh-CN.md`
- `conversations/2026-05-14-first-login-method-candidate-comparison.md`

Updates:

- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Compatibility Risks

No API, event, or data compatibility risk in this change.

The main risk is future agents treating the recommendation as implementation permission. The document explicitly forbids that.
