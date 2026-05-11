# vibit

vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.

Status: constitutional design phase

## What This Project Means By Agent-Native

Agent-native does not primarily mean that the server has AI features.

It means the server architecture is designed so AI coding agents can work inside it reliably:

- The architecture is explicit instead of tribal knowledge.
- Module ownership is declared instead of guessed.
- Public behavior is contract-first.
- Repeatable structure is generated.
- Business rules are tested as invariants.
- Cross-module communication is bounded.
- Change workflow is documented and verifiable.
- Documentation is written for both humans and agents.

AI gameplay features such as NPC agents, memory, model routing, tool calling, and simulations may become extensions later. They are not the foundation.

## Why This Exists

Many existing server codebases were built for human maintainers with local context, long memory, and implicit team conventions. AI coding agents can help in those codebases, but they often lose force when architecture rules are hidden, module boundaries are weak, tests are incomplete, or public contracts are unclear.

vibit starts from a different premise:

> The next generation of long-lived server software should be designed so agents can safely understand, modify, verify, and extend it.

The goal is not to make agents magically smarter. The goal is to make the codebase more legible, bounded, generated, contract-driven, and testable.

## Current Documents

- `CONSTITUTION.md`: canonical project constitution
- `CONSTITUTION.zh-CN.md`: Simplified Chinese translation
- `AGENTS.md`: repository-level operating guide for coding agents
- `AGENTS.zh-CN.md`: Simplified Chinese translation

English documents are canonical. Simplified Chinese translations are maintained for human readers and early project discussion.

## Intended Direction

vibit should evolve toward:

- Architecture manifests under `.arch/`
- Module manifests at `modules/<module>/module.yaml`
- Module-level agent guides at `modules/<module>/AGENTS.md`
- Contract-first commands, queries, events, errors, permissions, and migrations
- Generated scaffolds for repeatable framework structure
- Architecture checks that verify dependency, contract, event, and generated-file rules
- Change specs under `changes/<date>-<change-id>/`

The first serious prototype should prove one claim:

> Given a new backend requirement, an AI coding agent can identify the affected module, update the correct contracts, generate the correct structure, implement the behavior, add tests, run verification, and update documentation without damaging unrelated architecture.

## Early Reference Domain

A small game backend is the recommended first demonstration domain because it naturally contains state, permissions, events, consistency rules, and long-lived modules.

Suggested modules:

- Player accounts
- Inventory
- Currency
- Rewards
- Tasks or quests
- Match sessions

The demo should emphasize maintainability and agent workflow over feature count.

## Governance

Project decisions are governed by `CONSTITUTION.md`.

Before changing constitutional principles, ratifying the name, introducing a major architectural pattern, or making a breaking standard change, read the constitution and record the motivation, alternatives, compatibility impact, and migration path.

## Name

`vibit` is the product name.

The intended category phrase is:

```text
agent-native server framework
```

Before final ratification, the name should be checked against major public registries and platforms for obvious conflicts.
