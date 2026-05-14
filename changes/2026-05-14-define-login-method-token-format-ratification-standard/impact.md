# Impact

## Affected Modules

- `player`: referenced because first login methods may eventually create, link, or authenticate player accounts.
- `runtime`: referenced because token/session validation and request identity handoff are runtime-owned.

No module implementation changes are made.

## Module Ownership Impact

No ownership changes.

The standard reinforces existing ownership:

- Player account lifecycle remains separate from credentials, tokens, external identity links, sessions, and WebSocket state.
- Runtime application dispatch remains the future owner of normalized request identity and session validation handoff.
- WebSocket transport remains credential-neutral.
- Protobuf envelope fields remain metadata-only until a future protocol decision changes them.

## Public Contract Impact

No commands, queries, events, errors, or permissions are added.

The standard plans future contract work for login, token refresh, logout, validation, errors, permissions, and audit events.

## Event Impact

No event changes.

## Permission Impact

No permission changes.

## Data And Migration Impact

No data changes and no migrations.

The standard requires future schema gates before credential, token, external identity, or session tables may be added.

## Test Impact

No runtime tests are required because no runtime code changes.

Repository checks and documentation consistency checks are required.

## Documentation Impact

Adds:

- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `decisions/ADR-0024-login-method-token-format-ratification-boundary.md`
- `conversations/2026-05-14-login-method-token-format-ratification-standard.md`

Updates:

- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Compatibility Risks

No API, event, or data compatibility risk in this change.

The main risk is future agents treating the standard as implementation permission. The standard explicitly says it is not.
