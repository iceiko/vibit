# ADR-0030: Selected Login Token Boundary Checks

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-add-selected-login-token-boundary-checks/`

Related conversations:

- `conversations/2026-05-14-selected-login-token-boundary-checks.md`

Related artifacts:

- `docs/selected-login-token-boundary-checks.md`
- `docs/selected-login-token-boundary-checks.zh-CN.md`
- `docs/login-method-token-format-ratification.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/credential-token-session-schema-gates.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

M-013 has selected a first authentication posture: `device_credential_login`, opaque high-entropy access tokens, login-command token issuance, explicit request proof payloads, no refresh token in the first implementation, and PostgreSQL as the default durable target for future credential and token verifier schema gates.

Those decisions are specific enough that future agents may be tempted to implement concrete behavior immediately. The risky shortcuts are predictable: token parsing in WebSocket transport, credential lookup in player persistence, token columns in player account tables, runtime authentication generated output before generation is ratified, Protobuf authentication messages before protocol impact is ratified, and new auth dependencies without adoption records.

W-0072 adds a narrow repository check so the selected posture remains contract-only until a later bounded implementation milestone grants explicit permission.

## Decision

Add the `runtime.selected_login_token_boundary` repository check rule.

The rule statically checks:

- Selected posture markers.
- Implementation-deferral markers.
- Authentication contract blocked-until markers.
- Refresh-token deferral markers.
- Generated Go authentication contract shape deferral.
- Runtime authentication Protobuf source deferral.
- WebSocket carrier deferral.
- Credential/token/session/external-identity/audit migration deferral.
- Player account lifecycle storage separation.
- Forbidden authentication dependency additions.
- Machine-readable check output with stable `rule_id` and repository-relative `artifact` paths.

The rule runs as part of `node tools/vibit check runtime --json` and therefore as part of `node tools/vibit check all --json`.

This decision does not implement runtime authentication, token behavior, credential storage, token storage, session persistence, external identity linking, generated output, Protobuf messages, WebSocket route behavior, WebSocket handshake authentication, migrations, repositories, PostgreSQL adapters, or live external-service checks.

## Alternatives Considered

- Rely on the broader `runtime.authentication_token_session_boundary` rule only.
- Add prose standards without a machine-checkable rule.
- Wait until implementation begins before adding boundary checks.
- Add a broad static-analysis dependency.
- Require live PostgreSQL or external identity provider checks.
- Implement the first login behavior while adding the check.

## Rationale

The broader authentication/token/session rule protects design boundaries, but M-013 has now selected a narrower posture. A narrower posture creates narrower failure modes that a repository check can catch earlier and explain more directly.

Static checks are appropriate here because the current goal is to prevent unapproved files, dependencies, schema names, carriers, generated outputs, and implementation vocabulary from appearing. Live services would not improve this boundary and would make default verification harder for agents.

Nakama and Pitaya remain useful references, but their authentication and session API shapes must not leak into vibit by convenience. The check keeps the selected posture vibit-native until implementation work explicitly maps reference ideas into contracts, schemas, tests, and runtime code.

## Agent Reasoning Summary

The safe next step is a machine-readable circuit breaker for the chosen login/token direction. It keeps the project self-bootstrapping and implementation-ready without letting selection turn into accidental runtime code.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_boundary_clarity: high
  implementation_deferral: high
  machine_verifiability: high
  default_verification_cost: low
  dependency_load: low
  game_backend_reference_alignment: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future agents receive a specific failing rule when they add selected login/token implementation too early.
- `node tools/vibit check runtime --json` now verifies the selected posture boundary in addition to the broader authentication/token/session boundary.
- Runtime authentication remains blocked behind future schema, repository, adapter, generated-output, protocol, test, and implementation gates.
- The first implementation milestone must update this rule when it intentionally authorizes a concrete slice.
- Default verification remains local and does not require PostgreSQL, Redis-like services, external providers, Docker, Podman, or network access.

## Reversal Conditions

Revisit this decision if:

- The first implementation milestone authorizes a concrete auth slice and the check becomes too strict for approved files.
- A future security review requires different token or credential boundary names.
- Protobuf envelope or WebSocket carrier behavior is ratified and needs a narrower allowlist.
- A future dependency adoption record authorizes a JWT, OAuth/OIDC, password-hashing, provider SDK, key-management, or Redis-like dependency.
- Direct compatibility with Nakama, Pitaya, or another reference framework is explicitly ratified.

## Follow-Up

- Close M-013 if all ratification artifacts and checks now satisfy the milestone criteria.
- Create the next implementation gate without starting production authentication implicitly.
- Keep `runtime.selected_login_token_boundary` machine-readable and narrow as future authentication work is approved.
