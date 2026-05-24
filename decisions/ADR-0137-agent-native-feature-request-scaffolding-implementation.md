# ADR-0137: Agent-Native Feature Request Scaffolding Implementation

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-implement-agent-native-feature-request-scaffolding/`

Related conversations:

- `conversations/2026-05-24-agent-native-feature-request-scaffolding-implementation.md`

Related artifacts:

- `changes/_template/feature-request/`
- `docs/agent-native-feature-request-scaffolding.md`
- `docs/agent-native-feature-request-scaffolding.zh-CN.md`
- `docs/agent-native-feature-request-scaffolding-gate.md`
- `docs/agent-native-feature-request-test-workflow.md`
- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`

## Context

`ADR-0136` defined the feature request scaffolding gate and opened `M-157/W-0229 Implement agent-native feature request scaffolding`.

The maintainer's product direction is Nakama-first and AI-native: users should be able to state a backend requirement, and AI should help carry it through specification, acceptance criteria, tests, implementation, verification, and durable memory. The repository already has the workflow standard, but the next step is an executable source-first scaffold.

## Decision

Implement source-first feature request scaffolding with:

```text
changes/_template/feature-request/
tools/vibit scaffold feature
```

The scaffold creates or guides these artifacts:

```text
request.md
spec.yaml
impact.md
plan.md
checklist.md
verification.md
```

The command shape is:

```bash
node tools/vibit scaffold feature <change-id> --request <text> [--summary <text>] [--date YYYY-MM-DD] [--dry-run]
```

The command refuses to overwrite an existing change directory. It performs deterministic template replacement only. It does not generate runtime code, protocol sources, generated output, migrations, dependencies, SDKs, hosted artifacts, or direct compatibility shims.

Open:

```text
M-158/W-0230 Pilot scaffolded Nakama feature request intake
```

as the next-ready work item.

## Alternatives Considered

- Add templates only, without a command.
- Add a command only, with embedded templates.
- Delay scaffolding and move directly to a broad Nakama capability such as chat, groups, matchmaking, match runtime, or operations/admin.
- Add a richer interactive scaffold with prompts and schema validation.

## Rationale

Templates alone are useful, but a deterministic command makes the workflow executable for agents and contributors. Embedded-only templates would hide the scaffold shape from review. A richer interactive scaffold is premature; the first version should be simple, diffable, and safe.

The scaffold directly supports Nakama-first capability growth by forcing future feature work to record requirement, capability mapping, acceptance criteria, tests, boundaries, verification, and memory before coding.

Pitaya remains deferred because this work concerns repository workflow, not distributed topology, frontend/backend server roles, RPC, service discovery, groups, or cluster routing.

## Consequences

- `M-157/W-0229` is completed.
- `M-158/W-0230` is next-ready.
- `runtime.agent_native_feature_request_scaffolding_implementation` checks the scaffold, docs, command, manifests, rule catalog, and deferrals.
- Future non-trivial user-facing backend feature work can begin with `tools/vibit scaffold feature`.
- Runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, generated client libraries, hosted deployments, release artifacts, chat, groups, matchmaking, match runtime, operations/admin behavior, distributed runtime, and direct Nakama/Pitaya API compatibility remain deferred.

## Agent Reasoning Summary

The user asked to continue toward Nakama and emphasized that AI-native development and testing are the product purpose. This implementation turns that purpose into an executable intake path without broadening server behavior.

## Decision Weights

```yaml
decision_weights:
  ai_native_product_purpose: high
  nakama_product_alignment: high
  future_feature_velocity: high
  test_workflow_clarity: high
  implementation_boundedness: high
  direct_api_compatibility: low
confidence: high
```

## Reversal Conditions

Revisit this decision if:

- the scaffolded artifact shape proves insufficient for external users;
- future agents need structured schema validation beyond the current text templates;
- a broader product management workflow supersedes source-first change artifacts.

## Follow-Up

- Complete `W-0230`: use the scaffold on one bounded Nakama feature request intake to prove the workflow before further module expansion.
- Keep runtime behavior, protocol changes, generated output, migrations, dependencies, SDK publication, hosted demos, distributed runtime, and direct compatibility behind later bounded work items.
