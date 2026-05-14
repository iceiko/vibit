# Plan

## Files To Create

- `changes/2026-05-14-add-authentication-token-session-boundary-checks/`

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`

## Generated Artifacts

None.

## Handwritten Logic

Add a static runtime repository check named `runtime.authentication_token_session_boundary`.

The check should validate:

- Required standard and ADR artifacts.
- Manifest and guide references to the authentication/token/session boundary.
- Implementation-status markers that keep authentication, token behavior, credential storage, external identity linking, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player handlers, and WebSocket routes deferred.
- Metadata-only validator markers.
- Protobuf source and generated-output absence of credential/token/external-identity/auth-result fields.
- WebSocket transport absence of credential/token/handshake-auth handling.
- Runtime module, Protobuf adapter, player repository, and player account migration absence of credential/token/external-identity/session-persistence behavior.

## Tests

No new Go tests are required because this is repository tooling and documentation work. Existing Go tests run through `node tools/vibit check runtime --json`.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect rule runtime.authentication_token_session_boundary --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check work --json
node tools/vibit check change add-authentication-token-session-boundary-checks --json
node tools/vibit check all --json
node tools/vibit inspect next --json
git diff --check
```

## Rollback Notes

Rollback means removing the new check function and rule catalog entry, restoring the authentication design manifests to `W-0057`/`W-0058`, and removing this change spec.
