# ADR-0086: v0.1 Alpha Goal And Long-Term Product Target

Status: Accepted
Date: 2026-05-20
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-20-define-v0-1-alpha-goal/`

Related conversations:

- `conversations/2026-05-20-v0-1-alpha-goal.md`

Related artifacts:

- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`

## Context

The maintainer may pause the project because of token constraints and wants the repository to preserve enough product intent that a future Codex session started directly in this directory can continue without relying on the current chat context.

The project has advanced beyond constitutional design. It now has a real Go runtime foundation and a long series of bounded authentication, session, connection, logout, close, reconnect, session-carrier, and presence slices. The immediate next ready work item remains `M-106/W-0178`, a presence protocol query functional slice.

The maintainer clarified two product targets:

- short term: make a `v0.1 alpha` that real users can download, run, try, and use as a contribution entry point;
- long term: build an AI-era Pitaya or Nakama.

## Decision

Adopt `v0.1 alpha` as the short-term release target and `ai_era_nakama_pitaya_class_server_framework` as the long-term product target.

`v0.1 alpha` means the first developer-usable release: a single-node, PostgreSQL-backed, WebSocket + Protobuf game backend runtime with enough onboarding, authentication, protected gameplay, presence, logout, documentation, checks, and example flow for external developers to try the project and contribute.

The long-term target means Nakama/Pitaya-class capability coverage adapted to vibit's agent-native architecture. It does not mean direct Nakama or Pitaya API compatibility.

The current continuation queue stays active at `M-106/W-0178`. The next `continue` request should still advance the next ready work item through `.arch/work-items.yaml`.

## Alternatives Considered

- Keep only the existing Nakama/Pitaya parity roadmap and avoid a release-specific target.
- Declare the project already alpha because runtime code exists.
- Jump directly into broad social, matchmaking, match runtime, SDK, or distributed runtime work.
- Make direct Nakama/Pitaya API compatibility the long-term target.
- Replace the work-item continuation queue with a release checklist.

## Rationale

The project needs a shorter horizon than full Nakama/Pitaya product parity. A concrete alpha target keeps future work focused on making the existing runtime usable by external developers instead of endlessly expanding internal architecture.

The existing `.arch/work-items.yaml` continuation queue is already working. Replacing it would add process churn. The better path is to add a release target above the queue: keep advancing one bounded work item at a time, but prefer work that shortens the path to a runnable alpha.

Direct API compatibility would create a different product and likely distort vibit's contract-first, agent-native design. The durable target is same-class usefulness, not copied surface area.

## Agent Reasoning Summary

The maintainer needs the repository to be self-resuming after a pause. The most useful durable artifact is a short goal standard plus README update and manifest pointers, not another implementation gate. This keeps the next step concrete while preserving the broader product ambition.

## Decision Weights

```yaml
decision_weights:
  external_developer_usability: high
  continuation_without_chat_context: high
  development_velocity: high
  nakama_pitaya_class_ambition: high
  direct_api_compatibility: low
  broad_scope_before_alpha: low
confidence: high
```

## Consequences

- `docs/v0.1-alpha-goal.md` becomes the durable short-term release target.
- `README.md` becomes more user/contributor facing and no longer describes the project as only in constitutional design.
- Future agents should treat `W-0178` as the next immediate step and `v0.1 alpha` as the destination behind the next sequence of work.
- Alpha work should prioritize onboarding, authenticated E2E flow, presence query, runbook accuracy, example client/request loop, health/readiness/version/config surfaces, and acceptance checks.
- Full Nakama/Pitaya product-class parity remains long-term and remains governed by `docs/nakama-pitaya-product-parity-roadmap.md`.

## Reversal Conditions

Revisit this decision if the maintainer chooses a different first public release shape, decides to stop pursuing external developer usability, selects direct Nakama/Pitaya API compatibility, or replaces the work-item continuation workflow with a release-driven planner.

## Follow-Up

- Complete `W-0178`, the presence protocol query functional slice.
- Add first local onboarding/device credential issuance.
- Prove login -> connection binding -> protected inventory -> presence query -> logout end to end.
- Refresh the runtime runbook and add a minimal example client or request-loop script.
- Add a health/readiness/version/config check surface and a later alpha acceptance check.
