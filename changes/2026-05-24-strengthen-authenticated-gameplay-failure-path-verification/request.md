# Request

## Original Request

```text
继续推进，目标nakama 推进10步以上，推进10个小时以上，不要停止，我会离开10小时。
```

## Clarified Requirement

Complete `M-151/W-0223 Strengthen authenticated gameplay failure-path verification`.

The previous presence/status hardening slice opened this bounded follow-up. This slice must strengthen the authenticated local alpha gameplay proof around failure paths that a Nakama-class backend foundation must get right before broader services depend on protected requests.

## Selected User Requirement

```text
As a developer using the local alpha request flow, I want protected gameplay requests to fail closed when authentication proof is missing, malformed, invalid, expired, or revoked, and I want those failures to avoid leaking credential or token material.
```

## Nakama Capability Mapping

Capability family:

```text
identity_authentication_sessions
```

Nakama-style product value:

- Protected gameplay requests must require a valid authenticated identity handoff before inventory, presence, storage, social, realtime, matchmaking, or match runtime behavior depends on them.
- Failure paths are part of the product surface because developers need predictable errors, fail-closed routing, and redaction before building a prototype on the server.
- This slice uses vibit-native runtime and protocol surfaces. It does not pursue direct Nakama API compatibility.

## User-Visible Outcome

Developers and future agents now have a local alpha E2E proof that protected gameplay routes fail closed for:

- missing authenticated request wrapper;
- malformed authenticated request wrapper;
- malformed access-token text;
- unknown but well-formed access-token text;
- expired access token;
- revoked access token after logout;
- protected presence without authenticated proof;
- response redaction for raw device credential and raw access-token material.

## Non-Goals

- Add new authentication, session, gameplay, inventory, presence, or storage protocol routes.
- Add Protobuf source, generated output, migrations, dependencies, persistence, or startup wiring.
- Change login, logout, access-token validation, token refresh, cleanup, session validation, or route policy behavior.
- Add chat, friends, groups, parties, matchmaking, match runtime, leaderboards, economy, SDKs, operations/admin behavior, hosted deployments, release artifacts, public announcements, or paid promotion.
- Add Pitaya-style cluster/RPC/frontend-backend/service-discovery work.
- Add direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go` includes `TestAuthenticatedGameplayFailurePathsLocalAlphaFlow`.
- The test uses existing authenticated local alpha fixture surfaces and does not add production runtime behavior.
- Missing, malformed, invalid, expired, and revoked proof paths produce stable public authentication errors.
- Protected gameplay routes remain behind authenticated proof before successful domain behavior is returned.
- Error responses do not leak raw access-token or device credential text.
- `ADR-0131` records the verification decision.
- Repository checks include `runtime.authenticated_gameplay_failure_path_verification`.
- `M-152/W-0224 Select next Nakama prototype-ready capability after authenticated failure-path proof` is opened as the next-ready bounded follow-up.
