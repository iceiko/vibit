# ADR-0002: Productive Bootstrapping Over Meta-Tooling Drift

Status: Accepted  
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-bootstrapping-control/`

Related conversations:

- `conversations/2026-05-12-bootstrapping-control.md`

Related artifacts:

- `CONSTITUTION.md`
- `AGENTS.md`
- `tools/vibit`
- `changes/`
- `conversations/`

## Context

vibit deliberately began by bootstrapping standards, memory, checks, and inspect commands before selecting a server runtime. This was useful because the project thesis is agent-native maintainability, not only runtime feature delivery.

The maintainer then asked whether the project could drift in the future and how to keep self-bootstrapping within a controllable scope.

## Decision

vibit will treat self-bootstrapping as a means, not the product.

New standards, inspect commands, check commands, schemas, generators, and workflow rules should be added only when they directly support a concrete future server change, a runtime vertical slice, a module boundary, a contract, generated shape, test, verification path, or meaningful agent context reduction.

When the repository has enough standards and tooling to attempt a small backend capability, the default next step should be a runtime vertical slice rather than more meta-tooling.

## Alternatives Considered

- Continue adding meta-tooling until the framework feels complete.
- Stop all tooling work immediately and start runtime code without further governance.
- Add an automated hard limit on the number of standards or CLI commands.

## Rationale

The current standards and CLI are useful because they reduce ambiguity and create verification paths. However, the same pattern can become self-referential if every concern leads to another meta-tool.

The project should stay empirical: a standard is valuable when it helps an agent safely deliver or verify a real server capability.

A lightweight governance rule is enough right now. A hard automated limit would be premature because the project has not yet learned which checks matter during runtime implementation.

## Agent Reasoning Summary

The next highest-value proof is no longer another general-purpose inspect command. It is a small end-to-end backend slice that tests whether the existing standards actually help agents implement server behavior. Future tooling should earn its place by making that slice safer, smaller, or more verifiable.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: low
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Agents should justify new meta-tooling against an expected runtime or module implementation need.
- Additional standards should remain small and verifiable.
- The next major project step should favor a runtime decision and first backend vertical slice.
- Exceptions are allowed, but they must be recorded in a change spec or Agent Decision Record.

## Reversal Conditions

Revisit this decision if runtime work repeatedly fails because missing meta-tooling creates preventable ambiguity, unsafe edits, or untestable behavior.

## Follow-Up

- Decide the first runtime language and minimal server instance architecture.
- Build a small end-to-end backend slice, preferably inside the existing `inventory` module.
