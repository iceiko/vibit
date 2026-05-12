# ADR-0010: Foundational Dependency Adoption

Status: Accepted
Date: 2026-05-12
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-12-ratify-go-websocket-protobuf-runtime/`

Related conversations:

- `conversations/2026-05-12-go-websocket-protobuf-direction.md`

Related artifacts:

- `.arch/runtime.yaml`
- `AGENTS.md`
- `README.md`

## Context

The maintainer clarified that vibit can integrate mature, high-star, frequently used open-source projects for protocol and networking. The maintainer also emphasized that vibit is not a demo and must be treated as a long-maintained framework.

Foundational dependencies can accelerate implementation, but they can also hide architecture, create lock-in, or make agent reasoning harder if adopted casually.

## Decision

Foundational runtime dependencies require an explicit adoption record before they become part of the server architecture.

A foundational dependency is any external library or tool that shapes transport, protocol generation, persistence, dispatch, module loading, lifecycle management, observability, testing strategy, or generated code conventions.

The adoption record may be an ADR or a dedicated dependency section in a change spec. It must evaluate:

- Maintenance activity.
- License compatibility.
- API stability.
- Production adoption signals, including stars and usage when relevant.
- Replaceability and abstraction boundary.
- Agent readability.
- Testability.
- Generated-code compatibility.
- Operational fit for Go, WebSocket, and Protobuf.

High stars and frequent usage are useful signals, not sufficient reasons by themselves.

Domain modules must not directly depend on third-party transport, protocol, persistence, or framework libraries. Platform adapters own those dependencies and expose vibit-owned interfaces to modules.

## Alternatives Considered

- Avoid all external dependencies until the framework is mature.
- Freely adopt popular dependencies whenever they speed implementation.
- Create a heavy dependency review process before any package can be used.
- Use a lightweight adoption record for foundational dependencies.

## Rationale

vibit should benefit from mature ecosystem projects where they reduce risk and improve long-term quality. WebSocket, Protobuf, and network runtime work should not be hand-rolled without reason.

At the same time, vibit's core claim is agent-native maintainability. If dependencies hide important boundaries or let domain modules couple to framework details, agents will lose the clear architecture context this project is meant to provide.

A lightweight but explicit adoption record gives future agents durable context without blocking ordinary implementation work.

## Agent Reasoning Summary

Use mature open-source components, but put them behind vibit-owned boundaries and record why they were chosen. This keeps the framework serious without surrendering its architecture to dependency convenience.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: medium
  implementation_cost: medium
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- The first WebSocket library must be selected through an adoption record.
- The first Protobuf tooling stack must be selected through an adoption record.
- Foundation dependencies should live behind platform adapters or generation tooling, not inside domain modules.
- Dependency choices should include a migration or replacement path when practical.
- Agents must not add major external framework dependencies only because they are popular.

## Reversal Conditions

Revisit this decision if the adoption process becomes too heavy for small changes or too weak to prevent architecture drift.

## Follow-Up

- Add a dependency adoption document or template if repeated dependency choices need more structure.
- Use this standard before selecting WebSocket and Protobuf libraries.
- Add architecture checks for forbidden third-party dependency imports once Go code exists.
