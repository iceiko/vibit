# Request

Define the token and credential verifier algorithm redaction boundary before adding application authentication service code, token material generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, migration schema changes, or production authentication behavior.

This change advances `W-0091` under `M-019`.

## Maintainer Intent

The maintainer asked the agent to continue the work queue for an extended period, with the existing rule that routine bounded work should proceed without unnecessary confirmation and true architecture forks should stop for discussion.

The project should remain self-bootstrapping and controlled. Authentication work should be prepared carefully enough that future agents can implement it without drifting into ad hoc security, repository, protocol, or transport choices.

## Required Outcome

- Define verifier algorithm posture for device credential lookup/verifier material and opaque access-token lookup/verifier material.
- Define entropy, token encoding, digest classification, verifier key identifier handling, constant-time comparison, dependency posture, and redaction test expectations.
- State whether the first planned verifier posture can use Go standard library packages or requires an external dependency adoption record.
- Preserve all runtime authentication behavior deferrals.
