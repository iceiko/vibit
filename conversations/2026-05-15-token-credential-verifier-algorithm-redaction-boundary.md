# Conversation: Token And Credential Verifier Algorithm Redaction Boundary

Date: 2026-05-15
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-token-credential-verifier-algorithm-redaction-boundary/`

Related artifacts:

- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `decisions/ADR-0040-token-credential-verifier-algorithm-redaction-boundary.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue work for an extended period. Existing maintainer direction allows routine professional sequencing to proceed without stopping for unnecessary confirmation, while true architecture forks still require discussion.

The next ready work item was `W-0091 Define token and credential verifier algorithm redaction boundary`.

## Maintainer Narrative

The maintainer wants the project to stay self-bootstrapping and controlled before runtime authentication behavior is added. The goal is not to rush a demo, but to make the architecture explicit enough that future agents can continue implementation without drifting into ad hoc security, protocol, repository, or transport choices.

The maintainer previously emphasized that Nakama and Pitaya should remain active references for game server capability coverage, while vibit's differentiator is agent-native maintainability and stronger agent-facing boundaries.

## Agent Response Summary

The agent reviewed the credential record schema boundary, token verifier record schema boundary, token lifecycle posture, application authentication service interface boundary, authentication module manifest, runtime manifests, and existing check rules.

The agent concluded that the next safe step is to ratify a verifier algorithm posture and redaction-test boundary without adding runtime authentication code. The selected first posture is standard-library HMAC-SHA-256 over high-entropy credential and token material, with separate lookup and verifier digest classes, purpose labels, key identifiers, constant-time verifier comparison expectations, and redaction requirements.

## Decisions

- The first planned verifier algorithm family is `vibit_hmac_sha256_v1`.
- Future first-posture implementation may use Go standard library packages `crypto/hmac`, `crypto/sha256`, `crypto/subtle`, `crypto/rand`, and `encoding/base64`.
- No external cryptography, password-hashing, JWT, OAuth, OIDC, KMS, provider, or Redis-like token/session dependency is required for the first high-entropy posture.
- Raw access-token and raw device credential material must have at least 256 bits of entropy.
- Lookup digests are secret-adjacent index material and are not log-safe.
- Verifier digests are secret verifier material and require constant-time comparison.
- Raw proof, digest material, verifier keys, peppers, and full verifier key identifiers must not appear in public errors, logs, audit-safe facts, change specs, ADRs, conversation logs, or documentation examples.
- Runtime behavior remains deferred.

## Artifacts

- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.zh-CN.md`
- `decisions/ADR-0040-token-credential-verifier-algorithm-redaction-boundary.md`
- `changes/2026-05-15-define-token-credential-verifier-algorithm-redaction-boundary/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Open Questions

- Secret configuration, verifier key loading, and rotation still need a dedicated boundary before verifier code.
- The exact Go package path for verifier helpers remains for a later implementation gate.
- Protobuf authentication messages and WebSocket proof carriers remain deferred.

## Follow-Up

- Advance the secret configuration and verifier key loading boundary or the next explicitly ready preparation gate.
- Keep token generation, verifier comparison, login execution, token validation, logout execution, cleanup, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, and migration changes behind later gates.

## Redaction Notes

No real secrets, tokens, credentials, device identifiers, account details, or private data are stored in this conversation log. Synthetic sentinel strings may appear only in standards and tests as redaction fixtures.
