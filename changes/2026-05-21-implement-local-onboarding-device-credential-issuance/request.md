# Request

Implement `W-0182`: local onboarding device credential issuance.

The maintainer asked to continue. Per `docs/workflow.md`, continuation advances one `next_ready` work item unless blocked. `.arch/work-items.yaml` identified `W-0182` as the current next-ready item.

## Scope

- Implement only the local application service behavior authorized by `docs/local-onboarding-device-credential-issuance-gate.md`.
- Create an active player account and active device credential record in one application-owned unit of work.
- Generate server-issued device credential material with an explicit entropy reader.
- Compute credential lookup and verifier digests with existing helpers.
- Store only digest material and metadata.
- Return raw device credential text only once and only after unit-of-work success.
- Add focused tests for ordering, digest-only storage, fail-closed behavior, redaction, and unchanged login-route account creation behavior.

## Out Of Scope

- Public signup.
- WebSocket, HTTP, CLI, or Protobuf onboarding route exposure.
- Protobuf source or generated output changes.
- Migrations or repository interface changes.
- New dependencies.
- Access-token issuance from onboarding.
- Runtime session creation from onboarding.
- Changing existing login route account-creation behavior.
- External identity providers, password login, account recovery, account merge, or multi-device linking.
- Release publishing.
- Direct Nakama/Pitaya API compatibility.
