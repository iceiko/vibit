# Agent Decision Record Standard

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: `decisions/`

This document defines Agent Decision Records for vibit.

Agent Decision Records are ADR-style records for decisions made by humans and agents that affect architecture, standards, modules, generated files, or long-term maintainability.

## 1. Purpose

Conversation logs preserve how the discussion unfolded.

Change specs preserve how a change is executed.

Agent Decision Records preserve the durable rationale for decisions that future agents must respect.

They exist because agent-written code can look locally correct while violating long-term design intent.

## 2. Location

Decision records live under:

```text
decisions/
```

Template:

```text
decisions/_template/adr-agent.md
```

Decision files should use:

```text
decisions/ADR-0001-short-kebab-case-title.md
```

## 3. What To Record

Record decisions that affect:

- Constitutional principles
- Architecture standards
- Module ownership
- Public commands, queries, events, errors, or permissions
- Generated file conventions
- Verification gates
- Implementation language
- Server instance architecture
- Long-term maintainability tradeoffs

Do not create a decision record for every tiny implementation detail.

## 4. Public Rationale, Not Hidden Chain-of-Thought

Decision records must contain public rationale.

They should not contain private chain-of-thought, hidden reasoning dumps, or unverifiable internal monologue.

The standard format is:

- Context
- Decision
- Alternatives considered
- Rationale
- Agent reasoning summary
- Confidence
- Consequences
- Links

The "Agent reasoning summary" should be concise and inspectable. It should explain why the decision is appropriate without exposing hidden reasoning.

## 5. Links From Other Artifacts

When a decision affects a module, the module should reference the decision ID.

Example:

```yaml
decisions:
  - ADR-0001
```

When a decision affects a change, the change spec should reference it.

Example:

```yaml
decisions:
  - ADR-0001
```

The decision record is the right place for durable rationale. `module.yaml` should contain links and concise metadata, not long reasoning text.

## 6. Confidence And Weight

Decision records may include a confidence level:

```text
low
medium
high
```

They may also include decision weights:

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: medium
  implementation_cost: low
  reversibility: medium
  long_term_maintainability: high
```

These weights are not mathematical truth. They are compact public metadata that helps later agents understand what mattered.

## 7. Generated Immutability

Generated files are immutable to non-system agents.

Rules:

- Generated files must be declared in manifests.
- Ordinary agents must not hand-edit generated files.
- If generated output is wrong, change the source schema, template, or generator.
- A generated file override requires an explicit decision record or change spec.
- The reason for any override must be recorded.

Initial permission concept:

```text
generated_file_override
```

This is a standards-level permission until the framework has a real permission system.

## 8. Required Sections

Each decision record should include:

```text
# ADR-0001: Title

Status:
Date:
Decision Makers:
Related changes:
Related conversations:
Related artifacts:

## Context

## Decision

## Alternatives Considered

## Rationale

## Agent Reasoning Summary

## Decision Weights

## Consequences

## Reversal Conditions

## Follow-Up
```

## 9. Agent Rules

Agents should create or update a decision record when:

- A decision shapes future architecture.
- A decision rejects a plausible alternative.
- A decision affects generated files or module boundaries.
- A decision will be hard for a future agent to infer from code alone.

Agents must not use decision records as a place to dump lengthy private reasoning.

Keep decision records concise, public, and linked.
