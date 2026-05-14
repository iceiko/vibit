# Request

## Original Request

The maintainer asked the agent to continue ten steps according to professional judgment, stopping only for truly necessary decisions:

```text
按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。
```

## Clarified Requirement

Advance `W-0064` by defining the login method and token format ratification standard for M-013.

The work must:

- Define how first login methods are compared and ratified.
- Define how token model, token format, proof carrier, lifecycle, and storage implications are compared and ratified.
- Map Nakama and Pitaya reference concepts into vibit-native terms.
- Plan the remaining M-013 queue.
- Avoid adding implementation code or crossing authentication/token/session boundaries.

## User-Visible Outcome

The repository gains a durable standard, ADR, manifest references, and work queue that future agents can follow before implementing authentication.

## Non-Goals

- Implement production authentication.
- Select or implement final login methods in this work item.
- Select or implement final token format in this work item.
- Add token parsing, signing, refresh, revocation, storage, or validation.
- Add credential tables, external identity tables, token tables, session tables, or migrations.
- Add password hashing, JWT, OAuth, OIDC, provider SDK, cryptography, key-management, Redis-like, or other major dependencies.
- Change Protobuf envelope behavior.
- Change WebSocket handshake authentication behavior.
- Add runtime player account handlers or WebSocket routes.

## Unknowns

- Which login methods will be selected first.
- Which token format will be selected first.
- Which proof carrier will be selected first.
- Which schema gates are required before implementation.

## Acceptance Criteria

- The English standard exists.
- The Simplified Chinese translation exists.
- ADR-0024 records the durable decision.
- Manifests and agent guides reference the new standard.
- M-013 has a bounded follow-up work queue.
- Verification results are recorded.
