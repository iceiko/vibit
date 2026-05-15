# Conversation: Verifier Digest Computation Helper Implementation

Date: 2026-05-15  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-15-implement-verifier-digest-computation-helpers/`

Related artifacts:

- `runtime/internal/app/authentication/verifier_digest.go`
- `runtime/internal/app/authentication/verifier_digest_test.go`
- `docs/verifier-digest-helper-implementation-gate.md`
- `decisions/ADR-0048-verifier-digest-helper-implementation-gate.md`
- `.arch/work-items.yaml`

## Context

The verifier digest helper implementation gate was completed in `W-0102`. The next ready work item was `W-0103`, implementing helper-only verifier digest computation under the already-defined boundary.

## Maintainer Narrative

The maintainer asked the agent to continue development and had previously emphasized that small, non-branching implementation steps should proceed without stopping for confirmation, while real decision boundaries should be raised explicitly.

The maintainer expects the project to remain self-bootstrapping and controllable through bounded work items, architecture manifests, change specs, checks, and conversation memory.

## Agent Response Summary

The agent advanced `W-0103` as a helper-only implementation change.

The implementation adds application-owned verifier digest computation helpers under `runtime/internal/app/authentication`. The helpers build deterministic canonical input, compute HMAC-SHA-256 using the matching logical key from an already-validated `VerifierKeySet`, return copied digest bytes through `ComputedDigest`, and expose redacted typed errors.

The change updates repository checks so this exact helper file is authorized while other files remain blocked from digest computation and verifier comparison.

The change does not add verifier comparison, authentication service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf authentication messages, WebSocket proof carriers, repository changes, migration changes, startup wiring, authentication dependencies, or production authentication behavior.

## Decisions

- Implement digest computation now because `W-0103` authorizes the helper-only slice.
- Keep verifier digest comparison as a later gate because comparison is closer to proof validation and authentication behavior.
- Keep the next ready work item conservative: define the verifier digest comparison helper gate before adding comparison code.

## Artifacts

- `runtime/internal/app/authentication/verifier_digest.go`
- `runtime/internal/app/authentication/verifier_digest_test.go`
- `changes/2026-05-15-implement-verifier-digest-computation-helpers/`
- `conversations/2026-05-15-verifier-digest-computation-helper-implementation.md`

## Open Questions

- None for this helper-only implementation.

## Follow-Up

- Define the verifier digest comparison helper gate in `W-0104`.
- Keep authentication service behavior, protocol carriers, repositories, migrations, startup wiring, dependencies, and production authentication behavior behind later bounded work items.

## Redaction Notes

No secrets, tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, account details, or private data are recorded in this conversation log.
