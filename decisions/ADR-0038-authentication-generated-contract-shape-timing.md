# ADR-0038: Authentication Generated Contract Shape Timing

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-decide-authentication-generated-contract-shape-timing/`

Related conversations:

- `conversations/2026-05-14-authentication-generated-contract-shape-timing.md`

Related artifacts:

- `docs/authentication-generated-contract-shape-timing.md`
- `docs/authentication-generated-contract-shape-timing.zh-CN.md`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`

## Context

Runtime authentication has semantic contracts, selected login/token posture, credential and token verifier schema sources, an authentication repository interface, a PostgreSQL adapter, and an application-owned implementation boundary.

Generated Go contract shapes already exist for inventory and player semantic contracts. Runtime authentication has not generated shapes yet because earlier gates deliberately deferred generated authentication output.

Before application authentication service interfaces appear, agents need a machine-readable shape of the authentication contract surface that is easier to inspect than prose and safer than starting from runtime code.

## Decision

Generated Go authentication contract shapes should be introduced after this timing decision and before application authentication service interfaces or runtime authentication behavior.

This decision defines the planned source and output boundary only. It does not generate files.

Allowed source:

```text
contracts/runtime/authentication/
```

Allowed planned output root:

```text
runtime/internal/generated/contracts/runtime/authentication/
```

Allowed planned shape:

```text
runtime/internal/generated/contracts/runtime/authentication/<contract-type>/<ContractID>.go
```

The generated files must be metadata-only, trace to the semantic contract source and generator, expose the runtime family `authentication`, and remain immutable to non-system agents.

The next bounded work item may update generator/check support and generate authentication contract shapes. It must still avoid runtime authentication behavior, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository interface changes, and migration schema changes.

## Alternatives Considered

- Skip generated authentication contract shapes and design service interfaces directly.
- Generate authentication shapes in the timing decision step.
- Reuse the non-runtime path `runtime/internal/generated/contracts/runtime/<type>/<id>.go` without a family segment.
- Wait until after runtime authentication behavior exists.
- Generate Protobuf authentication messages before Go contract shapes.

## Rationale

The generated shape step gives agents a stable, inspectable bridge from semantic contracts to future implementation planning. It is useful before service interfaces because it reduces naming drift and keeps the contract surface visible without writing behavior.

The family segment is required because `runtime` can own multiple semantic families, including `session` and `authentication`. A family-aware output path keeps generated runtime surfaces understandable and avoids future collisions.

Deferring actual generation from this decision keeps the change bounded. The next work can update generator and check support, then generate files under explicit verification.

Nakama and Pitaya both reinforce the need for stable authentication and session vocabulary in game server frameworks. vibit adapts that lesson through semantic contracts and generated metadata rather than copying public APIs.

## Agent Reasoning Summary

The safest next sequencing is contract source, implementation boundary, generated metadata, then service interface design. Generated authentication contract shapes should arrive before handwritten application behavior, but the tooling support and generated files should be a separate, verifiable slice.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  generated_output_traceability: high
  security_boundary_clarity: high
  implementation_deferral: high
  tooling_cost: medium
  protocol_stability: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- `docs/authentication-generated-contract-shape-timing.md` becomes the standard for the next authentication generated-shape step.
- `docs/generated-output.md` now records the runtime family-aware path extension.
- Runtime authentication generated shape files remain absent in this change.
- The next work item can explicitly authorize generator/check support and generated files.
- Application service interfaces and runtime authentication behavior remain deferred.

## Reversal Conditions

Revisit this decision if:

- Runtime family-aware paths make Go package layout impractical.
- The generator cannot support runtime contract families without mixing generated metadata and runtime behavior.
- A future protocol decision requires authentication Protobuf messages before Go contract metadata.
- A security review requires additional redaction or classification fields before generated authentication metadata can be committed.

## Follow-Up

- Advance the next work item to update generator/check support and generate authentication contract shapes under the family-aware path.
- Keep runtime authentication service interfaces behind a later gate.
- Keep Protobuf messages, WebSocket proof carriers, authentication dependencies, and runtime behavior behind later gates.

