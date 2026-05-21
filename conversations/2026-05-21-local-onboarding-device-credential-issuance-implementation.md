# Conversation: Local Onboarding Device Credential Issuance Implementation

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0090-local-onboarding-device-credential-issuance-implementation.md`
Related change: `changes/2026-05-21-implement-local-onboarding-device-credential-issuance/`

## Context

The maintainer asked to continue. `.arch/work-items.yaml` identified `M-110/W-0182 Implement local onboarding device credential issuance` as the next ready item.

`W-0181` had already defined the gate in `docs/local-onboarding-device-credential-issuance-gate.md` and `ADR-0089`.

## Maintainer Narrative

The maintainer's instruction was to keep advancing the next bounded work item. No separate request was made to expose onboarding publicly or to skip the gate's constraints.

## Agent Response Summary

The agent implemented the local-only onboarding service method under `runtime/internal/app/authentication`.

The service validates the local request, generates server-issued device credential material, computes credential lookup and verifier digests, creates an active player account, stores an active digest-only device credential record, and returns the raw device credential text only after unit-of-work success.

The implementation did not add a protocol route, Protobuf source, generated output, migration, repository interface, dependency, public signup surface, access-token issuance from onboarding, runtime session creation from onboarding, or direct Nakama/Pitaya API compatibility.

## Decisions

- Accept `ADR-0090`.
- Add `runtime.local_onboarding_device_credential_issuance_implementation` as the repository check rule.
- Complete `M-110/W-0182`.
- Move the work queue to `M-111/W-0183 Next Alpha Direction Selection`.
- Keep existing `runtime.authentication.AuthenticateWithDeviceCredential` proof-based and non-creating.

## Artifacts

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `decisions/ADR-0090-local-onboarding-device-credential-issuance-implementation.md`
- `changes/2026-05-21-implement-local-onboarding-device-credential-issuance/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `modules/player/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The next alpha direction still needs selection.
- Public onboarding route exposure remains future work.
- Local CLI or example tooling that calls the service remains future work.
- Authenticated gameplay E2E remains future work.

## Follow-Up

Advance `M-111/W-0183` to select the next alpha-enabling direction after local onboarding implementation.

## Redaction Notes

No real raw credential text, raw credential bytes, credential lookup digest, credential verifier digest, HMAC input or output, verifier key material, database credential, access token, or full concrete verifier key id was recorded.
