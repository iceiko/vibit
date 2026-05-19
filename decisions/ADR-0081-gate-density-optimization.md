# ADR-0081: Gate Density Optimization

Status: Accepted
Date: 2026-05-19
Scope: Workflow governance, milestone gating, development velocity

## Context

The vibit project has accumulated 102 milestones, of which approximately 50% are pure gate/documentation milestones that do not produce implementation code. While this approach ensures architectural safety for security-critical boundaries (cryptography, verifiers, credential schema), it creates excessive overhead for functional implementation tasks (transport features, connection policy, route registration).

At the current gate density, each functional feature requires 2-3 milestones (gate → direction confirmation → implementation), consuming significant agent compute and wall-clock time without proportional safety gains.

## Decision

Adopt a three-tier gate strategy:

- **Tier 1 — Security-Critical** (Two-Step: Gate + Implementation): Applies to cryptography, verifier algorithms, Protobuf wire format, credential schema, token lifecycle, secret configuration. Retains the existing two-step process.
- **Tier 2 — Functional Implementation** (Single Step): Applies to transport features, application policy, registry behavior, route registration, protocol bridge, session lifecycle, connection lifecycle. Single implementation work item with boundary in change spec.
- **Tier 3 — Lightweight** (Direct Implementation): Applies to documentation, tooling, migration source, translation, check rules. Direct implementation without change spec unless non-trivial.

Direction confirmation milestones are no longer mandatory. Direction is managed through `ask_first` fields and `recommended_direction` on continuation semantics.

## Alternatives Considered

1. **Keep the current gate-heavy approach**: Every functional feature requires gate + direction confirmation + implementation. Rejected because the overhead is disproportionate for non-security work, and the project is still in early development where velocity matters.

2. **Remove all gates entirely**: Direct implementation for everything. Rejected because security-critical boundaries (cryptography, credential schema, token lifecycle) genuinely benefit from separate threat-model and boundary documentation before implementation.

3. **Two-tier system (security vs everything else)**: Simpler but loses the distinction between functional features (which still benefit from embedded boundary definitions) and lightweight changes (which need almost no ceremony).

## Rationale

The three-tier system preserves full protection for security-critical code while eliminating ~25% of milestone overhead for functional work. The key insight is that transport close handoff, connection policy, and route registration do not have the same threat model as verifier algorithms or credential schema — they don't need separate gate milestones.

## Agent Reasoning Summary

Analysis of M-001 through M-102 showed that gate-only milestones consistently produced documentation without reducing implementation risk for non-security features. The time spent on direction confirmation milestones could be eliminated by using `ask_first` fields, which already existed but were underutilized. The three-tier boundary aligns with the project's own classification of work: security primitives vs functional behavior vs documentation.

## Decision Weights

| Factor | Weight | Direction |
|--------|--------|-----------|
| Development velocity | High | Favor fewer gates |
| Security boundary integrity | Critical | Maintain full Tier 1 gates |
| Agent context efficiency | High | Favor embedded boundaries |
| Consistency with existing work | Medium | Do not retroactively change |

## Consequences

- Estimated ~25% reduction in milestone overhead for future work.
- Security-critical boundaries retain full protection.
- Functional features can be implemented in a single bounded step.
- Existing completed milestones are not retroactively changed.
- Chinese translation for feature-level documents is relaxed to asynchronous/deferred.
- Universal prohibitions (Nakama/Pitaya compatibility, dependency adoption) are consolidated in `AGENTS.md` Section 12 rather than repeated per-milestone.

## Reversal Conditions

- If a Tier 2 implementation produces a security regression that a separate gate would have caught, escalate that category to Tier 1.
- If the project reaches external users and needs stricter change control, consider adding review gates.
- If agent implementations repeatedly misinterpret embedded boundaries, reintroduce separate gate milestones.

## Follow-Up

- Apply Tier 2 strategy to the next functional milestone (M-103 and beyond).
- Monitor whether embedded boundaries in change specs provide sufficient agent context.
- Do not retroactively reclassify completed milestones.
