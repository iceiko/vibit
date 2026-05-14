# Impact

## Affected Modules

- `player`: future token subject depends on player account identity after credential and account policy success.
- `runtime`: future token validation must produce request identity through application-owned validation before domain dispatch.

No module implementation changes are made.

## Module Ownership Impact

No ownership changes.

The ratification reinforces that:

- WebSocket transport remains credential-neutral.
- Protobuf envelope metadata remains metadata-only.
- Application dispatch owns future token validation handoff.
- Domain modules do not issue, parse, validate, or store tokens.
- Player account lifecycle does not own tokens or credentials.

## Public Contract Impact

No contracts are added, changed, or removed.

The ratification creates future needs for login response and authenticated request proof contracts in W-0070.

## Event Impact

No event changes.

Future authentication, logout, revocation, and audit events remain deferred.

## Permission Impact

No permission changes.

Future login, token refresh, logout, token validation, linking, and revocation permissions remain deferred to W-0070.

## Data And Migration Impact

No data changes and no migrations.

The selected opaque token posture requires a future verifier storage or equivalent lookup boundary before implementation.

## Test Impact

No runtime tests are required because no runtime code changes.

Repository checks must verify the documentation, manifests, memory, work queue, and architecture state.

## Documentation Impact

Adds:

- `docs/first-token-format-proof-carrier-posture.md`
- `docs/first-token-format-proof-carrier-posture.zh-CN.md`
- `decisions/ADR-0026-first-token-format-and-proof-carrier-posture.md`
- `conversations/2026-05-14-first-token-format-proof-carrier-ratification.md`

Updates:

- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Compatibility Risks

No API, event, wire, or data compatibility risk in this change.

The main risk is future agents treating token ratification as implementation permission. The ratification document and manifests explicitly forbid that.
