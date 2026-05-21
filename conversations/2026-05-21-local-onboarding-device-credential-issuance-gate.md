# Conversation: Local Onboarding Device Credential Issuance Gate

Date: 2026-05-21
Related change: `changes/2026-05-21-define-local-onboarding-device-credential-issuance-gate/`
Related decision: `ADR-0089`

## Context

The maintainer asked to continue. The work queue reported `M-109/W-0181 Define local onboarding device credential issuance gate` as the next ready item.

The step was security-critical because future onboarding will create player account state, generate raw device credential material, compute credential verifier digests, store credential records, and present the raw credential once. The change therefore defined a gate before implementation.

## Maintainer Narrative

The maintainer's instruction was to continue the alpha path by advancing the next ready work item. No separate implementation behavior was requested in this conversation.

## Agent Response Summary

The agent defined a gate-only standard for local onboarding device credential issuance, recorded the decision, updated the work queue, and added repository checks. The work completed `W-0181` and made `M-110/W-0182 Implement local onboarding device credential issuance` the next ready item.

No runtime onboarding behavior was implemented. No credential was generated or displayed, and no player account or credential record was created through a new runtime flow.

## Decisions

- `ADR-0089` records the local onboarding device credential issuance gate.
- The future first surface is local-only application service behavior under `runtime/internal/app/authentication`, not public signup.
- Future behavior must create an active player account and active device credential record in the same application-owned unit of work.
- Future credential storage must remain digest-only.
- Future raw device credential text may be returned only once, and only after the unit of work commits.
- The existing `runtime.authentication.AuthenticateWithDeviceCredential` login route must not create accounts from `AccountCreationIntent`.

## Artifacts

- `docs/local-onboarding-device-credential-issuance-gate.md`
- `docs/local-onboarding-device-credential-issuance-gate.zh-CN.md`
- `decisions/ADR-0089-local-onboarding-device-credential-issuance-gate.md`
- `changes/2026-05-21-define-local-onboarding-device-credential-issuance-gate/`
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

- The implementation slice still needs to choose the concrete local request/result names, identifier generator shape, and test fakes within the gate's constraints.
- Public signup, CLI exposure, WebSocket routing, and production identity provider behavior remain deferred.

## Follow-Up

The next ready work item is `M-110/W-0182 Implement local onboarding device credential issuance`, bounded to the gate's local application service behavior.

## Redaction Notes

No raw credential text, raw credential bytes, credential lookup digest, credential verifier digest, HMAC input or output, verifier key material, database credential, or full concrete verifier key id was recorded.
