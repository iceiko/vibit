# Conversation: Verifier Digest Computation And Comparison Boundary

Date: 2026-05-15
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-verifier-digest-computation-comparison-boundary/`

Related artifacts:

- `docs/verifier-digest-computation-comparison-boundary.md`
- `decisions/ADR-0043-verifier-digest-computation-comparison-boundary.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue the work queue for an extended period. Existing maintainer direction allows bounded professional decisions to proceed without unnecessary confirmation while true architecture forks still require discussion.

The next ready work item was `W-0094 Define verifier digest computation and comparison boundary`.

## Maintainer Narrative

The maintainer wants vibit to remain self-bootstrapping, controlled, and prepared before authentication runtime behavior exists. The project should do necessary preparation rather than rush to a demo or a fragile first implementation.

The maintainer previously emphasized that Nakama and Pitaya should remain active references for game server capability planning, while vibit's differentiator is agent-native maintainability through explicit boundaries, manifests, checks, and durable decisions.

## Agent Response Summary

The agent reviewed the verifier algorithm/redaction boundary, secret configuration/key loading boundary, and token/credential material generation boundary.

The agent concluded that the next safe step is to define digest computation and comparison posture without adding HMAC helper code. The boundary defines canonical digest input, purpose-label and logical-key use, lookup digest handoff, key-set selection for rotation, constant-time verifier comparison, public failure collapse, redaction expectations, and Go standard library dependency posture.

## Decisions

- Future verifier digest computation and comparison is application-owned under `runtime/internal/app`.
- The first canonical input version is `vibit.auth.verifier.input.v1`.
- Canonical input uses an ASCII version header, zero separator, uint16 big-endian purpose-label length, ASCII purpose label, uint16 big-endian raw material length, and raw material bytes.
- Lookup digest equality is allowed only for candidate record selection, not as authentication proof.
- Future validation must compare verifier digest bytes in constant time.
- Missing lookup, verifier mismatch, unknown key id, unsupported algorithm, malformed proof, and expired or revoked proof collapse to the same public invalid-proof class unless a later semantic standard allows more detail.
- Future first-posture digest helpers may use Go standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle`.
- Runtime behavior remains deferred.

## Artifacts

- `docs/verifier-digest-computation-comparison-boundary.md`
- `docs/verifier-digest-computation-comparison-boundary.zh-CN.md`
- `decisions/ADR-0043-verifier-digest-computation-comparison-boundary.md`
- `changes/2026-05-15-define-verifier-digest-computation-comparison-boundary/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Open Questions

- Exact Go helper names remain for a later implementation gate.
- Whether repository lookup should accept multiple digest candidates in one call remains for a later repository boundary update.
- Login execution, token validation, logout, cleanup, Protobuf messages, and WebSocket proof carriers remain deferred.
- Timing equalization strategy for missing records versus invalid verifier digests remains for the implementation gate.

## Follow-Up

- Advance an authentication service implementation readiness gate or the next explicitly ready preparation gate.
- Keep verifier digest helper code, verifier comparison code, token generation code, credential generation code, secret loading, login execution, token validation, logout execution, cleanup, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, and migration changes behind later gates.

## Redaction Notes

No real secrets, tokens, credentials, generated material, device identifiers, account details, environment variable values, verifier key values, digest bytes, HMAC inputs, or private data are stored in this conversation log.
