# Conversation: Token And Credential Material Generation Boundary

Date: 2026-05-15
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-token-credential-material-generation-boundary/`

Related artifacts:

- `docs/token-credential-material-generation-boundary.md`
- `decisions/ADR-0042-token-credential-material-generation-boundary.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue the work queue for an extended period. Existing maintainer direction allows routine bounded work to proceed without unnecessary confirmation, while true architecture forks still require discussion.

The next ready work item was `W-0093 Define token and credential material generation boundary`.

## Maintainer Narrative

The maintainer wants vibit to remain self-bootstrapping and controlled before authentication behavior exists. The project should prepare enough standards that future agents can implement a real server slice without turning security, protocol, or repository choices into hidden assumptions.

The maintainer previously emphasized that Nakama and Pitaya should remain active references for game server capability planning, while vibit's differentiator is agent-native maintainability through explicit boundaries, manifests, checks, and durable decisions.

## Agent Response Summary

The agent reviewed the first login method set, token lifecycle posture, verifier algorithm boundary, secret configuration boundary, and authentication service interface boundary.

The agent concluded that the next safe step is to define material generation posture without adding generation code. The boundary selects server-issued application generation for the first device credential and access token posture, requires 256-bit raw material, defines URL-safe text presentation, forbids raw material storage, preserves one-time client presentation, and keeps generation helpers behind a later code gate.

## Decisions

- Future token and credential material generation is application-owned under `runtime/internal/app`.
- The first device credential source is server-issued application generation.
- The first access-token source is server-issued application generation.
- Raw credential and token material must have at least 256 bits of entropy.
- The first raw byte shape is 32 cryptographically random bytes.
- Text presentation is URL-safe unpadded Base64 or equivalent.
- Raw material is one-time client-visible and must not be stored.
- Future first-posture generation helpers may use Go standard library `crypto/rand` and `encoding/base64`.
- Runtime behavior remains deferred.

## Artifacts

- `docs/token-credential-material-generation-boundary.md`
- `docs/token-credential-material-generation-boundary.zh-CN.md`
- `decisions/ADR-0042-token-credential-material-generation-boundary.md`
- `changes/2026-05-15-define-token-credential-material-generation-boundary/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Open Questions

- Exact Go helper names remain for a later implementation gate.
- Verifier digest computation and constant-time comparison remain a separate boundary.
- Login execution, token validation, logout, cleanup, Protobuf messages, and WebSocket proof carriers remain deferred.
- Client storage guidance for the one-time device credential remains for client/protocol documentation.

## Follow-Up

- Advance a bounded verifier digest computation and comparison boundary, or the next explicitly ready preparation gate.
- Keep token generation code, credential generation code, secret loading, verifier digest computation, verifier comparison, login execution, token validation, logout execution, cleanup, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, and migration changes behind later gates.

## Redaction Notes

No real secrets, tokens, credentials, generated material, device identifiers, account details, environment variable values, verifier key values, or private data are stored in this conversation log.
