# ADR-0089: Local Onboarding Device Credential Issuance Gate

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-define-local-onboarding-device-credential-issuance-gate/`

Related conversations:

- `conversations/2026-05-21-local-onboarding-device-credential-issuance-gate.md`

Related artifacts:

- `docs/local-onboarding-device-credential-issuance-gate.md`
- `docs/local-onboarding-device-credential-issuance-gate.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `modules/player/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `modules/player/AGENTS.md`
- `modules/player/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0180` selected `define_local_onboarding_device_credential_issuance_gate` as the next alpha-enabling direction after the protected presence protocol query. The alpha goal requires a local player onboarding path that creates a player account, issues the first device credential, presents the raw credential once, and stores only verifier digests.

This capability is security-sensitive. It crosses player account creation, server-generated credential material, verifier digest computation, digest-only credential storage, unit-of-work ordering, one-time secret presentation, redaction, and future developer ergonomics.

The runtime already has the prerequisites for a bounded future implementation: player and authentication repositories, PostgreSQL adapters, verifier key configuration, credential material generation, digest helpers, login behavior, access-token validation, session creation, protocol login, protected gameplay routes, presence query, and logout. The missing boundary is first credential issuance.

## Decision

Define the local onboarding device credential issuance gate in:

```text
docs/local-onboarding-device-credential-issuance-gate.md
docs/local-onboarding-device-credential-issuance-gate.zh-CN.md
```

The future implementation candidate is a local-only application service method under:

```text
runtime/internal/app/authentication
```

with a candidate method name:

```text
OnboardLocalPlayerWithDeviceCredential
```

Future implementation must create an active player account and active device credential record in the same application-owned unit of work, using:

- `CreatePlayerAccount`
- `StoreCredential`
- `GenerateDeviceCredentialMaterial`
- `ComputeCredentialLookupDigest`
- `ComputeCredentialVerifierDigest`
- `VerifierKeySet`

The raw device credential may be returned only once, only after successful commit, and must never be stored.

This ADR does not implement onboarding, generate or display credentials, create player accounts through a new runtime flow, write credential records through a new runtime flow, expose a public protocol route, change Protobuf sources, change generated output, change migrations, add dependencies, publish a release, add production signup, add external identity providers, add password login, add account recovery, add multi-device linking, or adopt direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement onboarding immediately.
- Make the existing `AuthenticateWithDeviceCredential` route create accounts when `AccountCreationIntent` allows create.
- Add a public WebSocket Protobuf onboarding route first.
- Add a CLI or HTTP onboarding surface first.
- Seed credentials directly in SQL fixtures or migrations.
- Make onboarding issue an access token and runtime session directly.
- Defer onboarding and jump to authenticated gameplay E2E.

## Rationale

Immediate implementation would be too risky because the flow returns secret material and writes durable identity records. A gate makes ordering, ownership, one-time presentation, digest-only storage, and redaction explicit before code exists.

Keeping the first surface local-only removes the alpha blocker without committing to production signup or public protocol semantics. Developers need a way to obtain the first credential; they do not yet need OAuth, password login, account recovery, multi-device linking, abuse controls, or a public registration API.

The existing login route should remain proof-based. Turning `AccountCreationIntent` into account creation behavior would combine credential issuance with authentication and make public signup semantics implicit. The safer first posture is: local onboarding creates the credential, then normal login authenticates with it.

Nakama guides the need for a practical account/session/token entry path. Pitaya guides layer separation. vibit adapts both by keeping credential issuance application-owned and transport/protocol-neutral until a later route gate explicitly changes that surface.

## Agent Reasoning Summary

The maintainer asked to continue. `.arch/work-items.yaml` identified `W-0181` as the next ready security-critical gate. The correct continuation is to define the local onboarding boundary and queue a bounded implementation slice, not to write runtime onboarding behavior directly.

## Decision Weights

```yaml
decision_weights:
  alpha_blocker_removal: high
  credential_security: high
  one_time_secret_presentation: high
  unit_of_work_atomicity: high
  implementation_boundedness: high
  developer_usability: high
  public_api_commitment: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `docs/local-onboarding-device-credential-issuance-gate.md` and its Simplified Chinese translation define the gate.
- `runtime.local_onboarding_device_credential_issuance_gate` becomes the repository check rule.
- `M-109/W-0181` is completed as a gate-only milestone.
- `M-110/W-0182` becomes the next ready implementation slice.
- The next implementation can focus on a local application service method without protocol, migration, dependency, or release churn.
- Public signup, protocol onboarding routes, generated output, migrations, dependencies, external identity, password login, account recovery, multi-device linking, release publishing, and direct Nakama/Pitaya API compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- The maintainer chooses a public protocol signup route before local onboarding.
- The project decides first credentials must be provisioned outside the runtime repository.
- A production signup requirement becomes mandatory before alpha.
- The credential record or player account repository boundary changes.
- A future ADR adopts direct Nakama or Pitaya public API compatibility for account creation.

## Follow-Up

- Implement `W-0182`: the bounded local onboarding device credential issuance application service slice.
- After onboarding exists, prove the authenticated gameplay end-to-end path: onboarding -> login -> connection binding -> protected inventory -> presence query -> logout.
- Refresh the runtime runbook and add a minimal example client or request-loop script after the local flow is executable.
