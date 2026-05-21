# ADR-0090: Local Onboarding Device Credential Issuance Implementation

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-implement-local-onboarding-device-credential-issuance/`

Related conversations:

- `conversations/2026-05-21-local-onboarding-device-credential-issuance-implementation.md`

Related artifacts:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `modules/player/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0089` defined the local onboarding device credential issuance gate. The gate authorized only a local application service slice under `runtime/internal/app/authentication`, with no public signup surface, no protocol route, no Protobuf or generated output changes, no migrations, no repository interface changes, no dependencies, no access-token issuance from onboarding, and no runtime session creation from onboarding.

The repository already had the required building blocks: server-issued material generation, credential lookup and verifier digest helpers, storage-neutral player account and authentication repositories, PostgreSQL unit-of-work capabilities, and focused authentication service tests.

The short-term `v0.1 alpha` target needs a controlled local way for a developer to create the first player account and obtain the first device credential before using the existing login route.

## Decision

Implement:

```text
implement_local_onboarding_device_credential_issuance
```

as a local-only application service method:

```text
OnboardLocalPlayerWithDeviceCredential
```

The implementation:

- rejects missing display name before opening a unit of work,
- generates server-issued device credential material with `GenerateDeviceCredentialMaterial`,
- computes credential lookup and verifier digests with `ComputeCredentialLookupDigest` and `ComputeCredentialVerifierDigest`,
- opens the existing application-owned unit of work,
- obtains player and authentication repositories through unit-of-work capabilities,
- generates player id, player account event id, and credential record id through injected generators,
- creates an active player account with `CreatePlayerAccount`,
- stores an active digest-only device credential record with `StoreCredential`,
- returns raw device credential text only after unit-of-work success,
- does not issue an access token or runtime session from onboarding,
- does not change the existing `AuthenticateWithDeviceCredential` login route account-creation behavior.

The runtime startup composition now provides the service with random local onboarding id generators and a device credential random source. This is dependency composition only; it does not expose onboarding through WebSocket, Protobuf, HTTP, CLI, process startup auto-creation, or any public route.

## Alternatives Considered

- Keep the onboarding method deferred after the gate.
- Add a public protocol onboarding route immediately.
- Make the existing login route create accounts when `AccountCreationIntentCreate` is present.
- Return an access token or runtime session directly from onboarding.
- Let the caller supply the device credential material.
- Store raw credential text or bytes for later display.
- Add repository interface or migration changes in this slice.
- Defer startup dependency composition until a later local tool exists.

## Rationale

The local application service is the narrowest slice that removes the alpha blocker without turning local onboarding into production signup. It gives tests and later local tooling a deterministic service boundary while keeping route exposure, CLI UX, abuse controls, identity providers, account recovery, multi-device behavior, and production signup behind later decisions.

Creating the player account and credential record in one unit of work avoids orphaned accounts or orphaned credentials. Digest-only storage preserves the credential verifier posture, and returning raw device credential text only after commit prevents callers from treating failed writes as issued credentials.

Nakama demonstrates why a game backend needs a path to create/authenticate a player before gameplay. vibit adapts that capability as local-only onboarding first, not direct Nakama API compatibility.

Pitaya reinforces the separation between transport/protocol handlers and application behavior. vibit keeps onboarding in the application service instead of WebSocket transport or Protobuf routing.

## Agent Reasoning Summary

The previous gate already defined the ownership, ordering, redaction, and deferral rules. The implementation composes existing helpers and repositories rather than inventing new storage or protocol surfaces. The most important safety choice is keeping onboarding local and service-only while preserving the existing login route as proof-based authentication, not account creation.

## Decision Weights

```yaml
decision_weights:
  alpha_entrypoint_progress: high
  credential_redaction_safety: high
  transaction_boundary_clarity: high
  repository_interface_stability: high
  protocol_surface_restraint: high
  production_signup_scope: low
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.local_onboarding_device_credential_issuance_implementation` becomes the repository check rule for this slice.
- The application authentication service can create a local player account and issue the first device credential.
- Local tooling can later call this service without changing core service semantics.
- Public protocol onboarding remains deferred.
- Authenticated gameplay E2E work can now use local onboarding -> login -> bind connection -> protected gameplay -> presence query -> logout as the intended alpha path.
- The work queue moves to the next alpha direction selection point.

## Reversal Conditions

Revisit this decision if public signup becomes the first required onboarding surface, if account creation must be driven through a protocol route before local tooling, if id generation must be database-owned, if credential rotation/replacement becomes part of first credential issuance, if direct Nakama/Pitaya API compatibility is selected, or if production abuse controls require a different boundary before any onboarding service exists.

## Follow-Up

- Select the next alpha-enabling direction after local onboarding implementation.
- Likely next candidates include authenticated gameplay E2E, runtime runbook refresh, minimal example client/request-loop script, health/readiness/version/config surface, and alpha acceptance checklist.
- Keep public signup, protocol onboarding route, generated output, migrations, dependencies, access-token issuance from onboarding, runtime session creation from onboarding, production identity providers, password login, account recovery, multi-device linking, release publishing, and direct Nakama/Pitaya API compatibility behind explicit future work.
