# ADR-0028: Authentication Contract Error Permission Surfaces

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-authentication-contract-error-permission-surfaces/`

Related conversations:

- `conversations/2026-05-14-authentication-contract-error-permission-surfaces.md`

Related artifacts:

- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/authentication-contract-error-permission-surfaces.zh-CN.md`
- `contracts/runtime/authentication/`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

M-013 has already ratified `device_credential_login`, opaque high-entropy access tokens, login-command response issuance, explicit request proof payloads, and the first lifecycle/storage posture.

Before implementation can start, future agents need stable semantic surfaces for login, token issuance, token validation, logout, refresh posture, error codes, permissions, and audit-oriented events.

Without those surfaces, an implementation agent might add handlers, Protobuf fields, WebSocket authentication, raw-token storage, or ad hoc errors before the architecture can verify ownership and failure behavior.

## Decision

Create a runtime `authentication` contract family owned by `runtime/internal/app`.

Register the family in `.arch/contracts.yaml` and source the semantic manifests under:

```text
contracts/runtime/authentication/
```

Ratify these command surfaces:

- `AuthenticateWithDeviceCredential`
- `ValidateAccessToken`
- `LogoutAccessToken`
- `RefreshAccessToken`

Ratify these audit-oriented event surfaces:

- `AuthenticationSucceeded`
- `AuthenticationFailed`
- `TokenIssued`
- `TokenValidationFailed`
- `TokenRevoked`
- `LogoutRequested`

Ratify the `authentication_errors` error catalog and `authentication_permissions` permission catalog.

Extend `tools/vibit` contract inspection and checks so runtime authentication contracts are inspectable through the same JSON tooling as runtime session contracts.

No runtime handler, route, Protobuf message, generated contract shape, credential table, token table, session table, migration, token behavior, or authentication implementation is authorized by this decision.

## Alternatives Considered

- Documenting the surfaces only in prose.
- Adding Protobuf messages immediately.
- Generating Go contract shape files immediately.
- Combining authentication with the existing session validation family.
- Omitting refresh from W-0070 because refresh tokens are not in the first implementation.
- Deferring audit event surfaces until implementation.

## Rationale

Source-level contract manifests are the right level for this step. They are machine-readable, inspectable, and checkable, but they do not force wire or runtime implementation before schema gates exist.

Keeping authentication separate from `session_validation` preserves the distinction between authentication proof and session validation. The current session family remains metadata-only for the existing request handoff, while the new authentication family records the future selected login and token surfaces.

Including `RefreshAccessToken` as a deferred semantic surface prevents ambiguity. Refresh is intentionally unsupported in the first posture, and renewal happens through reauthentication with `device_credential_login`.

Audit-oriented events are ratified as semantic surfaces so future implementation can avoid leaking raw credentials, raw tokens, token verifier hashes, provider secrets, or hidden validation details.

## Agent Reasoning Summary

The decision gives future agents a concrete, queryable contract index before implementation. It narrows the next schema and runtime work without crossing ask-first boundaries such as handlers, routes, Protobuf changes, migrations, or major dependencies.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_boundary_clarity: high
  protocol_stability: high
  implementation_deferral: high
  game_backend_capability_alignment: high
  dependency_load: low
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Agents can inspect authentication surfaces with `node tools/vibit inspect contracts --module runtime --json`.
- `node tools/vibit check contracts --json` now verifies the runtime authentication family.
- Authentication commands, events, errors, and permissions have stable semantic names before implementation.
- Runtime session validation remains a separate semantic family.
- Protobuf and WebSocket behavior remain unchanged.
- Generated output remains unchanged.
- W-0071 can focus on credential, token verifier, external identity, and session schema gates.

## Reversal Conditions

Revisit this decision if:

- The first implementation milestone selects a different authentication owner.
- A future protocol decision replaces explicit request proof payloads with a Protobuf envelope or handshake carrier.
- Refresh tokens become part of the first production posture.
- A security review requires different error or audit event naming.
- A compatibility goal with Nakama, Pitaya, or another framework explicitly changes the surface vocabulary.

## Follow-Up

- Define credential, token, and session schema gates in W-0071.
- Add repository checks for selected login/token boundaries in W-0072.
- Keep runtime authentication implementation deferred until the next selected implementation milestone.
