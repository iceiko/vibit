# Request

Define the verifier digest computation and comparison boundary before adding verifier digest helper code, verifier comparison code, token generation code, credential generation code, secret loading, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This change advances `W-0094` under `M-022`.

## Maintainer Intent

The maintainer asked the agent to continue the work queue for an extended period, with the existing rule that routine bounded work should proceed without unnecessary confirmation and true architecture forks should stop for discussion.

The project should remain self-bootstrapping and controlled. Authentication work should be prepared carefully enough that future agents can implement digest computation and verifier comparison without inventing inconsistent HMAC input framing, weakening comparison behavior, or hiding authentication decisions in repositories.

## Required Outcome

- Define future verifier digest computation and constant-time comparison ownership.
- Define canonical byte input construction, purpose-label use, logical key use, key-id selection, lookup digest handoff, verifier digest comparison, and failure redaction expectations.
- State whether future first-posture digest helpers can use Go standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle` without an external dependency adoption record.
- Preserve all runtime authentication behavior deferrals.
