# Conversation: Authentication Contract Error Permission Surfaces

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-authentication-contract-error-permission-surfaces/`

Related artifacts:

- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/authentication-contract-error-permission-surfaces.zh-CN.md`
- `decisions/ADR-0028-authentication-contract-error-permission-surfaces.md`
- `contracts/runtime/authentication/`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue advancing the work queue unless a truly necessary decision required human input.

After W-0069 defined token lifecycle and storage implications, W-0070 required defining the contract, error, permission, and audit surfaces needed before authentication implementation.

## Maintainer Narrative

The maintainer has emphasized that vibit should be a long-term production-minded Agent-Native server framework rather than a demo. Routine technical details inside an already ratified milestone may be decided by the agent using professional judgment, while true product and architecture branch points should be brought back for confirmation.

## Agent Response Summary

The agent defined a new runtime `authentication` contract family owned by `runtime/internal/app`.

The family registers semantic commands for `AuthenticateWithDeviceCredential`, `ValidateAccessToken`, `LogoutAccessToken`, and a deferred `RefreshAccessToken` surface. It registers audit-oriented event surfaces for authentication success/failure, token issuance, token validation failure, token revocation, and logout request. It also defines an authentication error catalog and permission catalog.

The agent extended `tools/vibit` so runtime authentication contracts are inspectable and checkable through JSON tooling alongside runtime session contracts.

No runtime authentication implementation, handler, WebSocket route, Protobuf message, generated contract shape, table, migration, token behavior, credential lookup, or session persistence was added.

## Decisions

- Register `contracts/runtime/authentication/` as the semantic source directory for authentication surfaces.
- Own the first authentication surfaces in `runtime/internal/app`.
- Define `AuthenticateWithDeviceCredential` as the selected first login command surface.
- Define `ValidateAccessToken` as the explicit request proof validation surface.
- Define `LogoutAccessToken` as presented-access-token logout only.
- Define `RefreshAccessToken` as a deferred surface, with first renewal still handled through reauthentication.
- Define public-safe authentication error codes and permission dimensions.
- Define audit-oriented event surfaces without audit persistence.
- Preserve Protobuf, WebSocket, runtime handler, route, schema, migration, dependency, and generated-output deferral.

## Artifacts

- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/authentication-contract-error-permission-surfaces.zh-CN.md`
- `decisions/ADR-0028-authentication-contract-error-permission-surfaces.md`
- `contracts/runtime/authentication/`
- `changes/2026-05-14-define-authentication-contract-error-permission-surfaces/`

## Open Questions

- Exact credential table fields and token verifier table fields.
- Whether future token validation is implemented as request-level validation only or later combined with session persistence.
- Whether future audit events become durable database rows, logs, metrics, or an internal event bus.
- Whether admin token revocation enters the first runtime implementation milestone or remains an operations/admin milestone.

## Follow-Up

- Advance W-0071 to define credential, token, external identity, and session schema gates.
- Do not implement authentication runtime behavior until schema gates, checks, tests, and an implementation milestone exist.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
