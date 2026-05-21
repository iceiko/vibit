# ADR-0091: Next Alpha Direction Authenticated Gameplay E2E Slice

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-select-next-alpha-direction-after-local-onboarding/`

Related conversations:

- `conversations/2026-05-21-next-alpha-direction-authenticated-gameplay-e2e-slice.md`

Related artifacts:

- `docs/v0.1-alpha-goal.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `runtime/AGENTS.md`
- `README.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0182` implemented the first local-only onboarding/device credential issuance service. The runtime now has the separate pieces required for a first local alpha path: local onboarding, device credential login, runtime session metadata, first-message connection binding, protected inventory, protected presence query, and logout.

Those pieces are not yet proven or documented as one coherent developer flow. `docs/v0.1-alpha-goal.md` lists authenticated gameplay E2E as the recommended next direction after local onboarding.

This work item is a direction-selection step only. It does not implement the E2E flow, add protocol routes, change Protobuf sources, change generated output, add migrations, add dependencies, publish a release, add production signup, add external identity providers, add password login, add account recovery, add multi-device linking, add broad product modules, or add direct Nakama/Pitaya API compatibility.

## Decision

Select:

```text
define_authenticated_gameplay_e2e_slice
```

as the next alpha-enabling direction.

Create `M-112/W-0184` as the next ready bounded functional slice. The next slice should define and, if scoped by its change spec, prove the first authenticated gameplay E2E path:

```text
local onboarding -> login -> connection binding -> protected inventory -> presence query -> logout
```

This direction selection completes `M-111/W-0183`.

## Alternatives Considered

- Refresh the runtime runbook before proving the flow.
- Add a minimal example client or request-loop script before defining the E2E slice.
- Add health/readiness/version/config surfaces first.
- Add an alpha acceptance checklist first.
- Jump to chat, social, matchmaking, match runtime, SDK, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Rationale

Authenticated gameplay E2E is the next most useful alpha step because local onboarding removed the credential bootstrap blocker. A developer-usable alpha needs a single path that demonstrates account creation, authentication, bound connection identity, protected gameplay, visible presence, and logout working together.

The direction remains a bounded functional slice. The next work item should prefer using existing runtime capabilities and focused request-loop or test coverage before introducing new protocol surface.

Nakama informs the product pressure: a backend framework needs a coherent authenticate-then-play loop before broad gameplay services become meaningful.

Pitaya informs the layering pressure: the E2E path should prove route composition while preserving transport, protocol, application, and domain separation.

## Agent Reasoning Summary

The maintainer asked to continue after local onboarding was completed and verified. The alpha goal document names authenticated gameplay E2E as the recommended next direction, and the required primitives now exist as separate slices. The correct continuation is to queue a bounded E2E functional slice, not to expand directly into runbooks, clients, health surfaces, or broader product modules.

## Decision Weights

```yaml
decision_weights:
  alpha_usability: high
  existing_capability_integration: high
  developer_onboarding_value: high
  implementation_scope_control: high
  broad_product_expansion: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-111/W-0183` records the selected next alpha direction.
- `M-112/W-0184` becomes the next ready work item.
- The next work item should define and prove the authenticated gameplay E2E path using existing local onboarding, login, connection binding, protected inventory, presence query, and logout capabilities unless its change spec explicitly scopes a bounded addition.
- Runtime behavior, protocol expansion, generated output, migrations, dependencies, release publishing, production signup, external identity providers, password login, account recovery, multi-device behavior, broad product modules, and direct Nakama/Pitaya API compatibility remain deferred by this direction-selection step.

## Reversal Conditions

Revisit this decision if the maintainer redirects the alpha sequence, if local onboarding must be reworked before any E2E proof, if the runtime cannot be proven without a prerequisite runbook or health/config surface, or if direct Nakama/Pitaya API compatibility is explicitly selected.

## Follow-Up

- Complete `W-0184`: define and prove the authenticated gameplay E2E slice.
- Keep runtime runbook refresh, minimal example client/request-loop script, health/readiness/version/config surface, and alpha acceptance checklist queued after the E2E path is defined.
- Keep production signup, protocol expansion, generated output, migrations, dependencies, external identity, password login, account recovery, multi-device linking, broad product modules, release publishing, and direct compatibility behind explicit future work.
