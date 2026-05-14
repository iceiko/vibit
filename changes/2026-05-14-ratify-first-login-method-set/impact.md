# Impact

## Affected Modules

- `player`: future selected login behavior may create or authenticate player accounts after later gates.
- `runtime`: future authentication proof must become request identity through application-owned validation.

No module implementation changes are made.

## Module Ownership Impact

No ownership changes.

The ratification reinforces that:

- Player account lifecycle does not own credentials or tokens.
- Runtime application dispatch owns future validation handoff.
- WebSocket transport remains credential-neutral.
- Protobuf envelope metadata remains metadata-only.

## Public Contract Impact

No contracts are added, changed, or removed.

The ratification creates a future need for semantic login contracts in W-0070.

## Event Impact

No event changes.

Future authentication or audit events remain deferred.

## Permission Impact

No permission changes.

Future login, token refresh, logout, linking, and revocation permissions remain deferred to W-0070.

## Data And Migration Impact

No data changes and no migrations.

The selected method requires future credential storage gates, but does not add tables or migrations in this change.

## Test Impact

No runtime tests are required because no runtime code changes.

Repository checks must verify the documentation, manifests, memory, work queue, and architecture state.

## Documentation Impact

Adds:

- `docs/first-login-method-set.md`
- `docs/first-login-method-set.zh-CN.md`
- `decisions/ADR-0025-first-login-method-set.md`
- `conversations/2026-05-14-first-login-method-set-ratification.md`

Updates:

- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Compatibility Risks

No API, event, or data compatibility risk in this change.

The main risk is future agents treating login-method ratification as implementation permission. The ratification document and manifests explicitly forbid that.
