# ADR-0008: Go Server Runtime Language

Status: Accepted
Date: 2026-05-12
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-12-ratify-go-websocket-protobuf-runtime/`

Related conversations:

- `conversations/2026-05-12-go-websocket-protobuf-direction.md`

Related artifacts:

- `.arch/runtime.yaml`
- `README.md`
- `AGENTS.md`
- `decisions/ADR-0003-first-reference-runtime-language.md`
- `decisions/ADR-0009-websocket-protobuf-client-protocol.md`
- `decisions/ADR-0010-foundational-dependency-adoption.md`

Supersedes:

- `ADR-0003`

## Context

vibit needs a first server runtime language before real runtime implementation begins. A previous decision selected TypeScript on Node.js, but the maintainer later clarified that this had not been explicitly discussed or approved as the server runtime direction.

The maintainer proposed Go as the more appropriate long-term server language and emphasized that vibit must not be treated as a demo. It is intended to become a long-maintained framework for agent-native server architecture.

The existing `tools/vibit` CLI is implemented with Node.js standard-library APIs. That CLI remains repository tooling. Its implementation language must not determine the server runtime language.

## Decision

Go is the ratified first server runtime implementation language for vibit.

The broader vibit architecture remains portable at the manifest, contract, generation, and verification levels. Future language implementations may exist, but the first long-term server runtime path is Go.

JavaScript/Node.js may continue to be used for bootstrap tooling when it is already present and dependency-light. It is not the server runtime direction unless a future ADR explicitly changes that.

The first Go runtime should use Go modules. Exact workspace layout, package names, and dependency choices must be introduced by a dedicated change spec before implementation.

## Alternatives Considered

- Keep TypeScript/Node.js as the first server runtime.
- Use Go as the first server runtime.
- Use Rust for stronger compile-time guarantees and performance.
- Use C#/.NET or Java/Kotlin for mature backend ecosystems.
- Delay the language decision until after more tooling work.

## Rationale

Go is a strong fit for a serious server framework because it has a stable toolchain, straightforward deployment, readable concurrency primitives, explicit error handling, broad production adoption, and mature networking support.

Go also works well with WebSocket servers and Protobuf message handling without forcing a large application framework. That fits vibit's agent-native premise: keep module ownership, contracts, generated boundaries, and verification visible rather than hidden inside framework magic.

TypeScript remains useful for tooling and agent readability, but choosing the server runtime only because the bootstrap CLI uses Node.js would couple a long-term product decision to an incidental early implementation detail.

## Agent Reasoning Summary

The runtime language should optimize for long-lived server architecture, protocol implementation, operational clarity, and agent-safe maintainability. Go gives vibit a serious production server path while still allowing manifests, contracts, and generators to remain language-portable.

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

- New server runtime implementation work should target Go first.
- The TypeScript runtime slice and npm package baseline are removed from the mainline direction.
- `tools/vibit` remains a Node.js repository tool, not a server runtime.
- Future generated server shapes should target Go unless a change spec declares a narrower tooling-only target.
- Go package layout, test strategy, and dependency choices must be explicit before runtime implementation starts.

## Reversal Conditions

Revisit this decision if early Go runtime work shows that Go prevents agent-safe generation, clear module boundaries, reliable contract verification, or maintainable protocol implementation.

## Follow-Up

- Define the first Go workspace layout before writing runtime code.
- Choose foundational WebSocket and Protobuf dependencies through an adoption record.
- Add Go-aware runtime verification only when Go runtime files exist.
- Keep `.arch/runtime.yaml` as the machine-readable intake point for this decision.
