# ADR-0003: First Reference Runtime Language

Status: Superseded
Date: 2026-05-12  
Decision Makers: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-runtime-readiness-decisions/`
- `changes/2026-05-12-ratify-go-websocket-protobuf-runtime/`

Related conversations:

- `conversations/2026-05-12-runtime-readiness-decisions.md`

Related artifacts:

- `.arch/runtime.yaml`
- `README.md`
- `AGENTS.md`
- `decisions/ADR-0008-go-server-runtime-language.md`

Superseded by:

- `ADR-0008`

## Context

vibit needs enough runtime readiness to build the first backend slice without preventable churn. The framework should remain an agent-native server framework, not a language-specific experiment. At the same time, the first implementation needs a real language and package layout so contracts, generation, tests, and verification can become executable.

The existing repository tooling is a Node.js standard-library CLI. No production runtime code exists yet.

## Decision

This decision is superseded.

The earlier decision selected TypeScript on Node.js as the first reference implementation. That selection was made without enough explicit maintainer confirmation for a major architecture choice.

The ratified first server runtime implementation language is now Go, as recorded in `ADR-0008`.

The broader vibit architecture standard remains portable at the manifest, contract, generation, and verification levels. JavaScript/Node.js may remain useful for repository tooling, but it is not the server runtime direction unless a future ADR explicitly says so.

Initial dependency policy is conservative: start with the smallest dependency set that can support contracts, generation, tests, and server behavior. Do not add a major framework dependency until a change spec explains why the first slice needs it.

## Alternatives Considered

- Go: strong server runtime and deployment story, but less aligned with existing repository tooling and more verbose for early schema and generator experiments.
- Rust: strong correctness and performance story, but higher implementation cost for early agent-driven iteration.
- C#/.NET: strong typing and server tooling, but heavier initial project shape for the current repository.
- Java or Kotlin: mature backend ecosystem, but higher ceremony for the first proof slice.
- Continue without choosing a language: preserves optionality but blocks the first executable proof.

## Rationale

TypeScript would have given agents readable types, fast local feedback, practical schema tooling, and a broad backend ecosystem. It also aligned with the current `tools/vibit` implementation, which reduced bootstrapping friction.

That reasoning was not sufficient for a server runtime decision because it weighted existing tooling convenience too heavily. The server runtime is a long-term product architecture choice, not a side effect of the bootstrap CLI implementation.

Go is now preferred because vibit is intended to become a long-lived server framework with a serious network runtime, WebSocket-first gameplay protocol support, Protobuf wire messages, and production-oriented deployment expectations.

## Agent Reasoning Summary

The useful lesson from this superseded decision is procedural: major runtime choices must be discussed with the maintainer before being accepted. Existing tooling can inform a decision, but it cannot silently ratify the server architecture.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: high
  long_term_maintainability: high
confidence: superseded
```

## Consequences

- Do not use this ADR as authority for new server implementation work.
- Treat the TypeScript runtime slice and npm package baseline as historical artifacts that were removed from the mainline direction.
- Keep the Node.js CLI classified as tooling only.
- Use `ADR-0008` for the first server runtime language.

## Reversal Conditions

This decision has already been reversed by maintainer discussion and `ADR-0008`.

## Follow-Up

- Keep this ADR for historical traceability.
- Do not relabel removed TypeScript runtime work as the server runtime direction.
- Update conversation logs and change specs when future agents encounter this superseded context.
