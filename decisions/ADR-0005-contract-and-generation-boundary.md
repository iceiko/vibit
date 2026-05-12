# ADR-0005: Contract And Generation Boundary

Status: Accepted  
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-runtime-readiness-decisions/`

Related conversations:

- `conversations/2026-05-12-runtime-readiness-decisions.md`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `CONSTITUTION.md`
- `docs/module-manifest.md`

## Context

The constitution requires public behavior to be specified before implementation and generated files to be immutable to non-system agents. The project now needs a more concrete boundary for the first runtime slice so agents know which artifacts to edit and which artifacts to regenerate.

The exact schema source format is not finalized yet.

## Decision

The first runtime slice will treat these as public contract types:

- Commands
- Queries
- Events
- Errors
- Permissions

Contracts must be declared before handwritten behavior uses them.

Generated output should include repeatable shape such as typed contract objects, dispatch or route wiring, generated clients or call helpers, and test fixtures. Handwritten code should live in declared extension points such as command handlers, query handlers, domain policies, repositories, and invariant tests.

Generated files are immutable to non-system agents. If generated output is wrong, agents should change the source contract, manifest, template, or generator. A generated file override requires an explicit change spec or Agent Decision Record with the `generated_file_override` permission concept.

The first schema source may be YAML or JSON Schema. The final source format should be chosen when the first command, query, and event are specified.

## Alternatives Considered

- Handwrite all runtime files first and introduce generation later.
- Generate both structure and business logic.
- Use OpenAPI as the only contract source.
- Choose Protobuf before the first slice.

## Rationale

Agents need an obvious split between safe handwritten behavior and generated structure. Without that split, an agent fixing a small bug may rewrite routing, clients, or boilerplate in ways that appear locally correct but damage the framework.

Keeping contracts broader than HTTP APIs also matters. vibit must model backend behavior, not only request/response routes. Commands, queries, events, errors, and permissions are the first contract surface because they map directly to module ownership and server authority.

## Agent Reasoning Summary

The contract boundary should make agents edit the smallest meaningful source artifact. Public behavior belongs in contracts; repeated shape belongs in generated output; business rules belong in handwritten extension points.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- The first runtime proof must include at least one command, query, event, error, and permission.
- The first generator should be small and traceable rather than broad and clever.
- Agents must not hand-edit generated contract output during ordinary feature work.
- The framework needs a visible source-to-generated mapping before generated output becomes large.

## Reversal Conditions

Revisit this decision if the first proof slice shows that the chosen contract categories are too heavy, too vague, or insufficient to model real backend behavior.

## Follow-Up

- Choose the first schema source format when writing the initial inventory contracts.
- Define generated file markers and manifest declarations before generating runtime files.
- Add verification that generated files are declared before the generated surface grows.
