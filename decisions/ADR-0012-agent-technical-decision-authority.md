# ADR-0012: Agent Technical Decision Authority

Status: Accepted
Date: 2026-05-13
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-13-adopt-first-runtime-dependencies/`

Related conversations:

- `conversations/2026-05-13-technical-decision-authority-and-runtime-dependencies.md`

Related artifacts:

- `AGENTS.md`
- `.arch/conventions.yaml`

## Context

The project previously required the agent to discuss and confirm branch-producing foundational dependency choices before accepting concrete packages.

That was appropriate after an earlier mismatch: runtime language and protocol direction were inferred too aggressively. The maintainer then clarified a more precise rule: professional technical questions may be evaluated and decided by the agent after the maintainer grants that authority.

The maintainer's clarification:

```text
这些你按照专业评估来定。这些问题越专业越好，你自己评估就可以。以后这种问题就是要来跟我确认，我允许你自己评估决定，你才能自己评估决定。
```

## Decision

After explicit maintainer authorization, an agent may professionally evaluate and decide technical sub-decisions that are necessary to continue within an already ratified direction.

Examples include:

- Choosing between mature libraries for an already ratified dependency slot.
- Selecting verification tooling that enforces existing standards.
- Choosing a low-level adapter implementation detail behind a vibit-owned interface.
- Deferring a dependency when the first runtime slice does not need it yet.

The agent must still ask the maintainer before:

- Changing constitutional principles.
- Ratifying or replacing the project name.
- Changing the server runtime language, primary protocol direction, persistence direction, or core project thesis.
- Introducing a major new architecture pattern.
- Expanding product scope or changing the intended user-visible promise.
- Accepting a meaningful licensing, hosting, cost, operations, or vendor-lock-in commitment.
- Redefining module ownership or moving data ownership between modules.
- Making breaking public command, query, event, or data changes.
- Weakening validation, permission, generated-file, or verification rules.
- Adding a major external framework dependency that shapes application architecture rather than living behind an adapter.

Technical decisions made under delegated authority must be recorded in a change spec, architecture manifest, dependency adoption record, or ADR when they affect long-term maintainability.

## Alternatives Considered

- Require maintainer confirmation for every branch point.
- Let agents decide all technical and product direction questions without asking.
- Keep confirmation required only for constitutional changes.
- Use bounded technical delegation with explicit ask-first categories.

## Rationale

vibit needs expert technical momentum. Requiring maintainer confirmation for every library, tool, and internal adapter choice would slow development and create unnecessary conversation overhead.

However, vibit's central goal is agent-native maintainability under human direction. Agents should not silently reinterpret product intent, major architecture, licensing risk, cost commitments, or module ownership.

Bounded technical delegation keeps the project moving while preserving the maintainer's control over direction.

## Agent Reasoning Summary

The agent can decide professional engineering details once the maintainer grants that authority, but it must keep those decisions inspectable and bounded. Human confirmation remains required when a decision changes what vibit is, who it serves, how it is governed, or what commitments it imposes.

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

- Agents may continue through professional technical branch points after authorization instead of repeatedly asking.
- Agents must still document durable technical decisions.
- Ask-first rules remain active for product, governance, licensing-risk, operations, scope, ownership, and breaking-contract decisions.
- Future AGENTS guides should distinguish delegated technical decisions from maintainer-owned direction decisions.

## Reversal Conditions

Revisit this decision if agents start making choices that conflict with maintainer intent, or if the ask-first boundary becomes too vague for reliable development.

## Follow-Up

- Update AGENTS guides with the delegated technical decision rule.
- Link this ADR from architecture convention metadata.
