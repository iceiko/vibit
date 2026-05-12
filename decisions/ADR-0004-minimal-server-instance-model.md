# ADR-0004: Minimal Server Instance Model

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
- `CONSTITUTION.md`

## Context

The project needs a first server shape that is small enough to implement and verify, but serious enough to prove agent-native server architecture. A distributed actor runtime, cluster manager, or service mesh would add many operational concerns before the project has proven module contracts, generation, and invariants.

The constitution already requires strong module boundaries, server authority, commands, queries, events, and generated structure.

## Decision

The first runtime will use a modular monolith single-process server instance.

The runtime should distinguish these responsibilities:

- Transport adapters convert external requests into commands and queries.
- Application dispatch routes commands and queries to module handlers.
- Domain modules own state, policies, invariants, and public contracts.
- A platform layer owns shared runtime concerns such as transport, storage adapters, event dispatch, generated clients, and verification hooks.

The first instance may support HTTP first. WebSocket, realtime sessions, distributed actors, multi-process clustering, and orchestration are deferred until contracts and module boundaries are proven in a smaller shape.

## Alternatives Considered

- Distributed actor model from the beginning.
- Microservices from the beginning.
- A pure library without a server process.
- A game-server-specific realtime runtime as the first shape.

## Rationale

A modular monolith keeps the first proof understandable while still forcing the important boundaries: commands, queries, events, permissions, errors, modules, generated files, and tests.

This shape is also more agent-friendly at the start. Agents can inspect one process, one repository, and explicit module boundaries without reasoning about deployment topology, network partitions, service discovery, or distributed consistency before those concerns have earned their place.

## Agent Reasoning Summary

The first server instance should test the architecture thesis, not the hardest deployment model. A modular monolith can still be strict about boundaries while keeping the context small enough for agents to modify safely.

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

- Initial runtime code should not scatter domain logic into transport handlers.
- Cross-module communication must still use declared commands, queries, events, public module APIs, or generated clients.
- Performance and distribution concerns are valid later, but they should not dominate the first proof slice.
- The first module tests should run without requiring distributed infrastructure.

## Reversal Conditions

Revisit this decision if the first real backend requirements cannot be expressed honestly inside a modular monolith without hiding distributed assumptions or weakening module boundaries.

## Follow-Up

- Define the minimal runtime package layout for transport, dispatch, platform services, contracts, generated code, and modules.
- Use the first proof slice to test whether module boundaries remain enforceable inside one process.
