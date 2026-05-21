# ADR-0088: Next Alpha Direction Local Onboarding Device Credential Issuance Gate

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-select-local-onboarding-device-credential-issuance-gate/`

Related conversations:

- `conversations/2026-05-21-next-alpha-direction-local-onboarding-device-credential-issuance-gate.md`

Related artifacts:

- `docs/v0.1-alpha-goal.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`

## Context

`W-0178` completed the protected self-player presence protocol query. `node tools/vibit inspect next` then reported no `next_ready` work item, while `M-108 Next Alpha Direction Selection` was active.

`docs/v0.1-alpha-goal.md` identifies first local onboarding/device credential issuance as the next preferred alpha step after presence query. That capability is required before a new developer can create or obtain a first credential, log in, and exercise the authenticated runtime path locally.

The selected work is security-sensitive. It may eventually create player account state, generate raw device credential material for one-time client presentation, compute credential digests, and store only verifier records. `ADR-0082` requires security-critical work to use a separate gate before implementation.

## Decision

Select:

```text
define_local_onboarding_device_credential_issuance_gate
```

as the next alpha-enabling direction.

Create `M-109/W-0181` as a `next_ready` gate-only work item. The gate must define the local onboarding/device credential issuance boundary before any implementation code is added.

This direction selection completes `M-108/W-0180`.

## Alternatives Considered

- Implement local onboarding/device credential issuance immediately.
- Start the authenticated gameplay end-to-end slice next.
- Refresh the runtime runbook before onboarding exists.
- Add a minimal example client/request loop before onboarding exists.
- Add health/readiness/version/config surfaces before onboarding exists.
- Add an alpha acceptance checklist before the runnable alpha path exists.

## Rationale

Onboarding is the next blocking alpha capability because the runtime already has login, protected routes, connection binding, presence query, and logout, but a new local developer still lacks a documented bounded way to obtain the first device credential.

Immediate implementation would be too risky. The flow crosses credential material generation, digest computation, player account creation, repository mutation ordering, redaction, one-time raw secret presentation, and future developer ergonomics. A gate keeps the next implementation slice explicit and checkable.

The end-to-end gameplay slice should follow onboarding because it needs a credential issuance path as its starting point.

## Agent Reasoning Summary

The maintainer asked to continue after the queue reached an explicit direction-selection milestone. The alpha goal document gives a clear next product step, and existing security boundaries require a gate before credential issuance code. The correct continuation is therefore to queue the gate, not to write onboarding behavior directly.

## Decision Weights

```yaml
decision_weights:
  alpha_blocker_removal: high
  credential_security: high
  implementation_safety: high
  developer_usability: high
  immediate_runtime_behavior: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-108/W-0180` records the selected next alpha direction.
- `M-109/W-0181` becomes the next ready work item.
- The next implementation-relevant work must define a local onboarding/device credential issuance gate.
- Runtime behavior, protocol messages, generated output, migrations, dependencies, release publishing, and direct Nakama/Pitaya API compatibility remain unchanged.

## Reversal Conditions

Revisit this decision if the maintainer chooses to prioritize a different alpha blocker, removes local onboarding from the alpha goal, selects an external credential provisioning model, or decides that first local onboarding must be provided outside the runtime repository.

## Follow-Up

- Complete `W-0181`: define the local onboarding/device credential issuance gate.
- After the gate, implement the bounded local onboarding/device credential issuance slice.
- Then prove the authenticated gameplay end-to-end path: onboarding -> login -> connection binding -> protected inventory -> presence query -> logout.
