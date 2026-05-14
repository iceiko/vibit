# Request

Define the secret configuration and verifier key loading boundary before adding secret loading, token material generation, credential material generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This change advances `W-0092` under `M-020`.

## Maintainer Intent

The maintainer asked the agent to continue the work queue for an extended period, with the existing rule that routine bounded work should proceed without unnecessary confirmation and true architecture forks should stop for discussion.

The project should remain self-bootstrapping and controlled. Authentication work should be prepared carefully enough that future agents can implement it without drifting into ad hoc secret handling, hidden operations dependencies, misplaced key loading, or leaky diagnostics.

## Required Outcome

- Define secret configuration ownership and future key loading posture.
- Define verifier key separation, key identifier handling, rotation expectations, fallback behavior, development/test posture, production requirements, and redaction requirements.
- State whether the first local implementation can use process environment configuration or requires an external secret-manager dependency adoption record.
- Preserve all runtime authentication behavior deferrals.
