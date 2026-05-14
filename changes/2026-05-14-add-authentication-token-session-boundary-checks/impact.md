# Impact

## Affected Modules

Affected boundaries:

- `runtime`: repository checks now enforce the authentication/token/session design boundary.
- `player`: checks protect player account lifecycle persistence, repositories, Protobuf shapes, and runtime handler deferral.

No runtime authentication behavior is added.

## Module Ownership Impact

No ownership is moved.

The change reinforces existing ownership:

- `runtime/internal/app` owns request identity and session validation handoff.
- `runtime/internal/platform/transport/ws` owns WebSocket transport mechanics only.
- `runtime/internal/platform/protocol/protobuf` owns envelope conversion only.
- `modules/player` and `runtime/internal/modules/player` own player identity and account lifecycle only.
- Domain modules consume `RequestIdentity` but do not validate tokens or credentials.

## Public Contract Impact

No command, query, event, error, permission, or Protobuf contract is added or changed.

The new rule checks that authentication/token/session behavior does not appear before future contracts or standards ratify it.

## Data And Migration Impact

No migration is added.

The check verifies that player account lifecycle migrations and repositories do not introduce credentials, password hashes, provider subjects, tokens, runtime sessions, WebSocket state, request identity rows, or validation-result persistence.

## Test Impact

No Go test files are added. `node tools/vibit check runtime --json` runs the existing Go tests after the static checks.

## Documentation Impact

Updated:

- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Compatibility Risks

No runtime protocol, API, or data compatibility risk is expected because this change is check-only.

The main risk is false positives from overbroad static pattern checks. The implemented checks use targeted token, credential, external identity, session persistence, and handshake markers and preserve the existing metadata-only `player_id`/`session_id` fields.
