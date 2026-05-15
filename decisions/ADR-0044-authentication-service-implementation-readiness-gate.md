# ADR-0044: Authentication Service Implementation Readiness Gate

Status: Accepted
Date: 2026-05-15
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-authentication-service-implementation-readiness-gate/`

Related conversations:

- `conversations/2026-05-15-authentication-service-implementation-readiness-gate.md`

Related artifacts:

- `docs/authentication-service-implementation-readiness-gate.md`
- `docs/authentication-service-implementation-readiness-gate.zh-CN.md`
- `docs/runtime-authentication-implementation-boundary.md`
- `docs/application-authentication-service-interface-boundary.md`
- `docs/token-credential-verifier-algorithm-redaction-boundary.md`
- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `docs/token-credential-material-generation-boundary.md`
- `docs/verifier-digest-computation-comparison-boundary.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Context

The repository now has design boundaries for runtime authentication ownership, service interface shape, verifier algorithm posture, secret configuration, raw material generation, and verifier digest computation/comparison. The next risk is starting code without a single implementation entry gate.

Without a readiness gate, a future agent may implement the first authentication service slice while implicitly choosing package paths, secret loading behavior, helper names, test seams, redaction expectations, repository changes, protocol carriers, or dependency posture.

## Decision

Define an authentication service implementation readiness gate before adding application authentication service behavior.

The recommended future service package candidate is:

```text
runtime/internal/app/authentication
```

The first code slice must still be separately authorized. This readiness gate does not implement code.

The recommended first code gate is local verifier key configuration loading because later material generation, digest helpers, and service behavior need validated key material and should not parse secret configuration inside login logic.

The readiness gate preserves separate gates for secret loading, material generation, verifier digest helpers, verifier comparison, service behavior, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket proof carriers, repository changes, migration changes, major dependencies, and production authentication behavior.

## Alternatives Considered

- Start authentication service code immediately after digest boundary ratification.
- Implement token and credential generation before secret configuration loading.
- Implement digest helpers before deciding package ownership and test entry criteria.
- Let the first implementation agent choose the service package path.
- Combine Protobuf/WebSocket proof carriers with the first service code slice.
- Open direct Nakama or Pitaya API compatibility as the next milestone.

## Rationale

Authentication code will cross application orchestration, repositories, secret configuration, token material, verifier digests, request identity handoff, and future protocol carriers. A readiness gate reduces the chance that these choices leak into implementation as undocumented assumptions.

Choosing `runtime/internal/app/authentication` as the candidate package keeps behavior application-owned without polluting generic app dispatch packages or moving authentication decisions into the domain repository.

Secret configuration loading should be first because it provides validated key material for generation and digest helpers. It is smaller than login behavior, testable without PostgreSQL, and preserves protocol and transport deferrals.

Nakama and Pitaya remain references for capability coverage and server vocabulary. vibit adapts them through application-owned proof validation and request identity handoff, not by copying public APIs or making transport handlers validate proof.

## Agent Reasoning Summary

The repository is close to implementation but still benefits from one readiness consolidation step. This gate lets future agents start code with explicit entry criteria, write boundaries, tests, and sequencing instead of re-reading every prior standard and inferring the code queue.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  implementation_control: high
  security_boundary_integrity: high
  testability: high
  protocol_deferral: high
  dependency_minimization: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- `docs/authentication-service-implementation-readiness-gate.md` becomes the standard for the first authentication service implementation entry criteria.
- `runtime.authentication_service_implementation_readiness_gate` becomes the repository check rule for this gate.
- The next recommended implementation gate is local verifier key configuration loading.
- Runtime authentication behavior remains unimplemented.

## Reversal Conditions

Revisit this decision if:

- The first implementation must live outside `runtime/internal/app/authentication` for a concrete Go packaging reason.
- A security review requires a different first implementation order.
- Protocol proof carriers become a prerequisite for any meaningful local implementation test.
- A future repository boundary must change before key configuration loading can be useful.

## Follow-Up

- Define or implement the local verifier key configuration loading gate as the next bounded work item.
- Keep authentication service behavior, token generation, credential generation, digest computation, verifier comparison, login execution, token validation, logout execution, cleanup, Protobuf messages, and WebSocket proof carriers behind later bounded work items.
