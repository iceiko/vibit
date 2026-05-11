# ADR-0001: Agent-Native Maintainability Is The Foundation

Status: Accepted  
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-agent-decision-records-and-inspect/`

Related conversations:

- `conversations/2026-05-12-founding-session.md`
- `conversations/2026-05-12-external-ai-feedback-on-traceability-tooling-and-immutability.md`

Related artifacts:

- `CONSTITUTION.md`
- `AGENTS.md`
- `.arch/`
- `docs/agent-decision-record.md`

## Context

The project began from a clarification: "AI-native server" does not primarily mean a server with AI gameplay features. It means a server architecture designed from first principles so AI coding agents can understand, extend, verify, and maintain it.

The maintainer emphasized that many existing server codebases are hard for agents to work on because architecture rules, conventions, and boundaries are implicit.

## Decision

vibit is an agent-native server framework whose foundation is agent-native maintainability.

AI gameplay features may be added later, but they are application-layer capabilities, not the architectural foundation.

## Alternatives Considered

- Build a conventional game server framework with optional AI features.
- Build an AI gameplay backend focused first on NPC agents, memory, and model routing.
- Build a generic backend framework without special agent-oriented constraints.

## Rationale

Agent-native maintainability is the distinct project thesis.

It gives the project a clear reason to exist beyond existing game server frameworks and AI feature frameworks.

It also justifies stricter requirements:

- Explicit architecture
- Machine-readable manifests
- Contract-first behavior
- Generated structure
- Bounded changes
- Verification gates
- Conversation logs
- Agent Decision Records

## Agent Reasoning Summary

The project should optimize for the work agents repeatedly struggle with: finding the correct boundary, preserving long-term design intent, updating contracts before code, and proving changes with checks. AI gameplay features do not solve those problems by themselves.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: low
  long_term_maintainability: high
confidence: high
```

## Consequences

- Standards and tooling come before runtime feature count.
- The first CLI checks architecture and manifests rather than running a game server.
- Server language and instance architecture remain open until standards can evaluate them.
- Future features must be judged by whether they improve agent maintainability or clearly belong to application-layer AI capabilities.

## Reversal Conditions

This decision should only be revisited if the project intentionally changes category away from agent-native server architecture.

## Follow-Up

- Add schema validation for manifests.
- Add more inspect commands with JSON output.
- Create explicit decision records for server implementation language and server instance architecture when ready.
