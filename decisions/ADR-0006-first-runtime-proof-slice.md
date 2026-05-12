# ADR-0006: First Runtime Proof Slice

Status: Accepted  
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-runtime-readiness-decisions/`

Related conversations:

- `conversations/2026-05-12-runtime-readiness-decisions.md`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/modules.yaml`
- `modules/inventory/module.yaml`

## Context

The project needs to prepare carefully before implementation, but preparation must remain tied to a concrete proof. The constitution already defines the preferred proof unit:

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

The repository already contains a draft `inventory` module manifest. Inventory is a small game-backend domain with state, permissions, events, capacity rules, and common edge cases.

## Decision

The first runtime proof slice should use the `inventory` module unless a stronger requirement appears before implementation begins.

The slice should prove the complete agent-native workflow with the smallest meaningful capability:

- Define a bounded requirement and change spec.
- Declare one command, one query, one event, one error, and one permission.
- Generate repeatable structure from contracts and manifests.
- Implement handwritten behavior only in declared extension points.
- Add focused tests for command behavior, query behavior, emitted events, permissions or errors, and invariants.
- Run one default verification command that includes architecture, schema, module, and runtime checks.
- Update documentation and translations.

## Alternatives Considered

- Player accounts as the first slice.
- Currency as the first slice.
- Rewards or quests as the first slice.
- Build a generic transport example without a domain module.

## Rationale

Inventory is concrete without being too broad. It can demonstrate ownership, validation, state transitions, events, errors, permissions, and invariants without requiring account lifecycle, payments, matchmaking, or long-running sessions.

The slice should be small enough to finish but complete enough to expose whether vibit's standards actually help agents build and maintain server behavior.

## Agent Reasoning Summary

The first runtime proof should test the whole workflow, not maximize feature count. Inventory gives a practical domain where module boundaries and contract-first development can be verified quickly.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future preparation should be judged by whether it changes how the inventory proof slice will be built, verified, or maintained.
- The first runtime implementation should avoid broad framework abstractions that the inventory slice does not need.
- The existing `modules/inventory/module.yaml` should be updated before implementation starts.
- Other candidate modules remain useful later, but they should not distract from the first proof.

## Reversal Conditions

Revisit this decision if the maintainer chooses a different first product requirement, or if inventory proves too weak to test the framework's contract, generation, and verification claims.

## Follow-Up

- Before runtime implementation, refine the inventory module manifest with real ownership, contracts, invariants, generated files, and required tests.
- Create a dedicated change spec for the first inventory runtime slice.
