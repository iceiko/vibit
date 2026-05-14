# Request

## Original Request

The maintainer asked to continue the current work queue:

```text
继续推进
```

The active next-ready item was:

```text
W-0058 Add authentication token session boundary architecture checks
```

## Clarified Requirement

Add repository checks that enforce the authentication, token, and session validation design boundary from `docs/authentication-token-session-validation.md` and `ADR-0023`.

The checks must preserve metadata-only `player_id` and `session_id` as unauthenticated, keep credentials/tokens/external identities/session persistence out of player account lifecycle persistence, and prevent WebSocket transport, Protobuf adapters, domain modules, and player repositories from silently owning authentication behavior before ratification.

## User-Visible Outcome

`node tools/vibit check runtime --json` should include `runtime.authentication_token_session_boundary`, and `node tools/vibit inspect rule runtime.authentication_token_session_boundary --json` should expose the new rule metadata.

`node tools/vibit inspect next --json` should advance from `W-0058` to `W-0059`.

## Non-Goals

- Do not implement production authentication.
- Do not choose a concrete login method.
- Do not choose token format, token carrier, signing, refresh, expiration, revocation, rotation, or token storage behavior.
- Do not add credential storage, password hashing, OAuth, OIDC, external identity linking, or session-store dependencies.
- Do not add runtime authentication code, token parsing, credential lookup, external identity linking, or session persistence.
- Do not change Protobuf envelope behavior.
- Do not change WebSocket handshake authentication behavior.
- Do not add runtime player account handlers.
- Do not add WebSocket routes.
- Do not require live external services for default repository checks.
- Do not declare metadata-only identity sufficient for production permissions.

## Unknowns

None for this check-only work item.

## Acceptance Criteria

- [x] Rule catalog contains `runtime.authentication_token_session_boundary`.
- [x] Runtime check includes the new rule.
- [x] The rule checks standard, ADR, manifest, guide, metadata-only validator, Protobuf, generated-output, WebSocket, player repository, and migration boundaries where statically checkable.
- [x] The checks do not add runtime authentication code or dependencies.
- [x] Manifests and guides reference the new rule.
- [x] English and Simplified Chinese documentation are updated.
- [x] Work queue advances to `W-0059`.
- [x] Verification is recorded.
