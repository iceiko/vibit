# ADR-0003: First Reference Runtime Language

Status: Accepted  
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-runtime-readiness-decisions/`

Related conversations:

- `conversations/2026-05-12-runtime-readiness-decisions.md`

Related artifacts:

- `.arch/runtime.yaml`
- `README.md`
- `AGENTS.md`

## Context

vibit needs enough runtime readiness to build the first backend slice without preventable churn. The framework should remain an agent-native server framework, not a language-specific experiment. At the same time, the first implementation needs a real language and package layout so contracts, generation, tests, and verification can become executable.

The existing repository tooling is a Node.js standard-library CLI. No production runtime code exists yet.

## Decision

The first reference implementation will use TypeScript on Node.js.

The broader vibit architecture standard remains language-neutral. TypeScript is the first proof vehicle, not a permanent restriction on future runtimes.

Initial dependency policy is conservative: start with the smallest dependency set that can support contracts, generation, tests, and server behavior. Do not add a major framework dependency until a change spec explains why the first slice needs it.

## Alternatives Considered

- Go: strong server runtime and deployment story, but less aligned with existing repository tooling and more verbose for early schema and generator experiments.
- Rust: strong correctness and performance story, but higher implementation cost for early agent-driven iteration.
- C#/.NET: strong typing and server tooling, but heavier initial project shape for the current repository.
- Java or Kotlin: mature backend ecosystem, but higher ceremony for the first proof slice.
- Continue without choosing a language: preserves optionality but blocks the first executable proof.

## Rationale

TypeScript gives agents readable types, fast local feedback, practical schema tooling, and a broad backend ecosystem. It also aligns with the current `tools/vibit` implementation, which reduces bootstrapping friction.

The decision is deliberately scoped to the first reference implementation. A successful TypeScript slice should prove the agent-native workflow. It should not prevent future adapters, generators, or runtime implementations in other languages.

## Agent Reasoning Summary

The most useful next step is to make one runtime path concrete while keeping the architectural standard portable. TypeScript balances agent readability, code generation friendliness, test ergonomics, and current repository momentum better than the alternatives for the first proof.

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

- The first runtime package layout should assume TypeScript and Node.js.
- Runtime scaffolding should avoid dependencies that make the framework hard to inspect or generate.
- Documentation should describe TypeScript as the first reference implementation, not as the whole project identity.
- Future language implementations may exist if they preserve the same contracts, manifests, generated boundaries, and verification expectations.

## Reversal Conditions

Revisit this decision if the first runtime slice shows that TypeScript prevents clear contracts, reliable verification, simple deployment, or agent-safe maintainability.

## Follow-Up

- Define the first TypeScript package layout when runtime scaffolding starts.
- Choose the package manager and test runner only when they are needed by executable runtime code.
- Keep `.arch/runtime.yaml` updated as runtime decisions become more precise.
