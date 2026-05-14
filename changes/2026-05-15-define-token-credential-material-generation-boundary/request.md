# Request

Define the token and credential material generation boundary before adding token generation code, credential generation code, secret loading, verifier digest computation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This change advances `W-0093` under `M-021`.

## Maintainer Intent

The maintainer asked the agent to continue the work queue for an extended period, with the existing rule that routine bounded work should proceed without unnecessary confirmation and true architecture forks should stop for discussion.

The project should remain self-bootstrapping and controlled. Authentication work should be prepared carefully enough that future agents can implement generation behavior without inventing raw secret storage, weak entropy, client metadata credentials, or hidden dependencies.

## Required Outcome

- Define future token and credential material generation ownership.
- Define raw material entropy enforcement, text encoding, one-time presentation, storage prohibition, repository handoff, redaction, test fixture posture, and dependency posture.
- State whether future first-posture generation helpers can use Go standard library `crypto/rand` and `encoding/base64` without an external dependency adoption record.
- Preserve all runtime authentication behavior deferrals.
