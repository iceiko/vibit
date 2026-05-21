# ADR-0092: Authenticated Gameplay E2E Functional Slice

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-define-authenticated-gameplay-e2e-slice/`

Related conversations:

- `conversations/2026-05-21-authenticated-gameplay-e2e-functional-slice.md`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/v0.1-alpha-goal.md`
- `AGENTS.md`
- `runtime/AGENTS.md`
- `README.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0091` selected authenticated gameplay E2E as the next alpha-enabling direction after local onboarding. The runtime already had the necessary pieces as separate slices: local onboarding, device credential login, runtime session metadata in login responses, first-message connection binding, protected inventory routes, protected presence query, and logout.

The missing evidence was that these pieces compose into one coherent local path. A developer-usable alpha needs that path before runbook, client-script, health/config, and acceptance-check work can be made accurate.

## Decision

Define and prove:

```text
local onboarding -> login -> connection binding -> protected inventory -> presence query -> logout
```

as a focused protocol-level Go E2E test over existing capabilities:

```text
runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
```

The test:

- calls `OnboardLocalPlayerWithDeviceCredential` directly as the local-only onboarding entrypoint;
- sends `runtime.authentication.AuthenticateWithDeviceCredential` through the existing Protobuf frame handler;
- registers a server-observed open connection in the in-memory connection registry and sends `runtime.authentication.BindConnection`;
- sends protected `inventory.GrantItem` and `inventory.GetInventory` frames through `AuthenticatedRequest`;
- sends protected `runtime.presence.GetPlayerPresence` through `AuthenticatedRequest`;
- sends `runtime.authentication.LogoutAccessToken`;
- confirms the same access token no longer satisfies a protected inventory request after logout, proving post-logout protected-route rejection;
- checks that raw credential/token text is not stored in the test authentication repository and is not leaked by binding/logout error surfaces.

This slice does not add production runtime behavior, protocol routes, Protobuf sources, generated output, migrations, repository interfaces, dependencies, public signup, production identity behavior, example clients, runbook refresh, release artifacts, broad product modules, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Documentation-only E2E definition without an executable proof.
- A local request-loop script instead of a Go test.
- A live PostgreSQL-only E2E proof.
- Adding a public onboarding protocol route before proving the path.
- Adding an example client before proving the lower-level request loop.
- Refreshing the runtime runbook before the flow was executable.

## Rationale

A focused Go E2E test is the narrowest proof because it can compose existing runtime layers without committing a developer-facing script or new protocol surface too early. It gives the runbook and future example client a stable behavioral target.

The direct local onboarding call is intentional. Onboarding is currently a local application service, not a public protocol route or production signup surface.

The presence assertion is intentionally scoped to online bound connection state. Runtime session id linkage inside presence snapshots remains deferred because route policy still uses request-level token proof and connection binding does not yet promote persisted session validation into active-connection policy.

Nakama informs the product pressure: developers expect a coherent authenticate-then-play loop. Pitaya informs the layering pressure: transport, protocol, application, and domain responsibilities should remain separated while the route composition is proven.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0184` was next ready. The safest implementation was to prove the selected path with existing pieces rather than add any new public capability. The test uses production application/protocol behavior where available and keeps only the fake repositories, entropy, IDs, and clock test-local.

## Decision Weights

```yaml
decision_weights:
  alpha_usability: high
  existing_capability_integration: high
  protocol_surface_restraint: high
  executable_verification: high
  runbook_readiness: medium
  production_signup_scope: low
  broad_product_expansion: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.authenticated_gameplay_e2e_functional_slice` becomes the repository check rule for this slice.
- `W-0184` completes the first authenticated gameplay E2E proof.
- Runbook refresh can now target a real flow instead of a planned flow.
- A future example client or request-loop script can follow the same sequence.
- Production signup, external identity, password login, account recovery, multi-device linking, chat/social/matchmaking/match runtime modules, health/config surfaces, release publishing, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if the runtime needs a different first alpha flow, if onboarding becomes public protocol behavior before alpha, if request-token proof is replaced by session/bound identity route policy, if presence must expose runtime session ids in the first alpha proof, or if the maintainer explicitly selects direct Nakama/Pitaya API compatibility.

## Follow-Up

- Refresh the runtime runbook around the proven alpha path.
- Add a minimal example client or request-loop script after the runbook target is accurate.
- Add health/readiness/version/config surfaces.
- Add an alpha acceptance checklist once the developer flow is stable.
