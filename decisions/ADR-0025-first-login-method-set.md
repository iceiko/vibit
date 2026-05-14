# ADR-0025: First Login Method Set

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-compare-first-login-method-candidates/`
- `changes/2026-05-14-ratify-first-login-method-set/`

Related conversations:

- `conversations/2026-05-14-first-login-method-candidate-comparison.md`
- `conversations/2026-05-14-first-login-method-set-ratification.md`

Related artifacts:

- `docs/first-login-method-candidates.md`
- `docs/first-login-method-candidates.zh-CN.md`
- `docs/first-login-method-set.md`
- `docs/first-login-method-set.zh-CN.md`
- `docs/login-method-token-format-ratification.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

M-013 exists to ratify first login methods, token format, proof carrier behavior, lifecycle semantics, schema gates, and repository checks before authentication implementation.

W-0065 compared device credential, guest/anonymous, custom ID, email/password, external provider, and service-auth login families. The comparison recommended `device_credential_login` because it gives a low-friction game login path without adding OAuth, OIDC, password hashing, provider SDKs, WebSocket handshake authentication, Protobuf envelope changes, or runtime player routes.

The repository still has no production authentication implementation. Metadata-only `player_id`, `session_id`, and `connection_id` remain unauthenticated.

## Decision

Ratify the first login-method set as:

```yaml
first_login_method_set:
  - device_credential_login
```

`device_credential_login` means the client proves possession of a high-entropy installation credential. The credential is secret proof material, not a raw operating-system device ID or other public metadata.

The method is production-capable only after the required contract, token, credential, storage, redaction, test, and repository-check gates are complete.

The method may create a player account or authenticate an existing account only after account creation and credential lookup policy are ratified. Account linking, account recovery, anonymous upgrade, external identity linking, service authentication, email/password login, guest login, and provider login remain deferred.

This decision does not implement authentication, add credential storage, add token behavior, add session persistence, change the Protobuf envelope, change WebSocket handshake behavior, add runtime player handlers, or add WebSocket routes.

## Alternatives Considered

- `guest_anonymous_login`
- `custom_id_login`
- `email_password_login`
- `external_provider_login`
- `service_authentication`
- Metadata-only `player_id` or `session_id`
- Direct Nakama account/auth API compatibility
- Pitaya session binding as implementation API

## Rationale

The first player login method should be small enough for agents to implement correctly after gates are complete, but serious enough to become a production path.

`device_credential_login` is the best first fit because it keeps complexity local:

- No major provider dependency is required by the selection itself.
- No password workflow is required.
- No WebSocket handshake authentication is required.
- No Protobuf envelope change is required.
- Credential storage can be separated from player account lifecycle tables.
- The proof can flow through an application-owned boundary before domain dispatch.

Nakama supports device-style authentication as part of broad game backend capability coverage. vibit adapts that low-friction capability without copying Nakama's public API or token behavior.

Pitaya provides useful vocabulary around session context and handler binding. vibit adapts the separation by keeping validation outside transport and domain handlers.

## Agent Reasoning Summary

The selected method gives future agents the narrowest production-minded implementation path: define a login contract, define credential storage and redaction, define token/session behavior, add checks, then implement through application-owned validation. Other login methods remain valuable, but each introduces wider product, dependency, or trust-boundary work.

## Decision Weights

```yaml
decision_weights:
  production_safety: medium
  game_onboarding_ergonomics: high
  agent_context: high
  dependency_load: low
  storage_complexity: medium
  abuse_and_recovery_load: medium
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future login implementation work must target `device_credential_login` first unless a later ADR supersedes this decision.
- Future agents must not add guest, custom ID, email/password, provider, or service-auth login by convenience.
- The selected method must not treat raw device IDs, player IDs, session IDs, connection IDs, or envelope metadata as proof.
- Credential storage remains separate from player account lifecycle persistence.
- Token format and carrier behavior remain separate decisions for W-0067 and W-0068.
- Runtime authentication implementation remains deferred until contracts, schemas, checks, tests, and ownership rules are complete.

## Reversal Conditions

Revisit this decision if:

- A security review shows high-entropy installation credentials are unacceptable as the first player login method.
- The product direction requires cross-device identity before any device credential path is useful.
- Token/session ratification proves that another login method produces a simpler and safer first production slice.
- A compatibility goal with Nakama, Pitaya, or another external framework is explicitly ratified.

## Follow-Up

- Compare token format and carrier options.
- Ratify the first token format and proof carrier posture.
- Define token lifecycle and storage implications.
- Define authentication contract, error, and permission surfaces.
- Define credential, token, and session schema gates.
- Add repository checks for selected login/token boundaries.
