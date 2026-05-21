# ADR-0104: First Alpha User Discovery Loop

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-define-first-alpha-user-discovery-loop/`

Related conversations:

- `conversations/2026-05-21-first-alpha-user-discovery-loop.md`

Related artifacts:

- `docs/first-alpha-user-discovery-loop.md`
- `docs/first-alpha-user-discovery-loop.zh-CN.md`
- `docs/release-execution-final-authorization.md`
- `docs/alpha-developer-flow.md`
- `README.md`
- `README.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`v0.1.0-alpha.1` is live as a source-first alpha release. `W-0195` refreshed the README so a developer can understand vibit, try the local alpha path, and inspect the agent-native project memory. The maintainer explicitly said the project now needs to find users.

The next risk is not lack of internal planning. The next risk is talking to the wrong people, asking vague questions, or converting early attention into unbounded scope. `W-0196` therefore defines the first discovery loop before any public announcements beyond the GitHub release record.

## Decision

Define the first alpha user discovery loop as a bounded learning loop.

The loop targets:

- Go game/backend developers.
- Developers who have used or evaluated Nakama, Pitaya, Colyseus, Pomelo, Agones, or custom backends.
- Engineers using AI coding agents on backend codebases.
- Open-source contributors interested in explicit architecture and checkable workflows.
- Prototype builders who can tolerate source-first alpha setup.

The loop records:

- outreach surfaces that may be prepared or pointed at;
- feedback fields to capture;
- review questions;
- success signals;
- stop conditions;
- the next work item, `W-0197 Prepare first alpha feedback intake surfaces`.

Do not authorize public announcements beyond the GitHub release record, paid promotion, hosted deployments, additional release artifacts, runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, broad product module expansion, or direct Nakama/Pitaya API compatibility.

The repository check rule for this decision is `runtime.first_alpha_user_discovery_loop`.

## Decision Record

```yaml
first_alpha_user_discovery_loop: defined
completed_work_item: W-0196
decision: ADR-0104
check_rule: runtime.first_alpha_user_discovery_loop
release_identifier: v0.1.0-alpha.1
source_first_alpha_release_available: true
user_discovery_loop_defined: true
target_developer_segments_recorded: true
outreach_surfaces_recorded: true
feedback_capture_recorded: true
success_signals_recorded: true
stop_conditions_recorded: true
public_announcements_beyond_github_release_authorized_by_this_decision: false
paid_promotion_authorized_by_this_decision: false
hosted_deployment_authorized_by_this_decision: false
additional_release_artifacts_authorized_by_this_decision: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
direct_nakama_pitaya_api_compatibility_added: false
next_direction: first_alpha_feedback_intake_surfaces
```

## Alternatives Considered

- Start broad public announcements immediately.
- Add hosted demos or packaged artifacts before feedback intake is ready.
- Resume runtime feature work before observing first external friction.
- Treat stars, forks, or views as the main success measure.
- Ask for feedback without a redaction and stop-condition posture.

## Rationale

The source alpha is now visible enough to learn from, but the project still needs precise signal more than broad attention. A bounded loop focuses outreach on developers likely to understand game/backend server frameworks and the AI-maintainability problem.

Preparing intake before announcements reduces the chance that useful feedback disappears into chats or gets converted into unbounded scope. Keeping runtime and artifact deferrals explicit preserves the source-first alpha posture.

Nakama and Pitaya remain reference baselines for capability class. In this slice, their relevance is contributor/user intake discipline, not API compatibility or feature expansion.

## Agent Reasoning Summary

The agent treated `W-0196` as a planning and repository-check item. The work records who to learn from first, how to capture feedback, and where to stop. The next step is an intake-surface preparation item, not broad announcement execution.

## Decision Weights

```yaml
decision_weights:
  user_learning_focus: high
  outreach_scope_control: high
  source_alpha_honesty: high
  redaction_safety: high
  runtime_scope_change: none
  artifact_scope_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0196` is completed as the durable first alpha user discovery loop record.
- The first target developer segments, surfaces, feedback fields, success signals, and stop conditions are documented.
- The next ready work item becomes `W-0197 Prepare first alpha feedback intake surfaces`.
- Public announcements beyond the GitHub release record still require later explicit authorization.
- No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, hosted deployment, operations/admin expansion, authentication/session behavior change, broad product module, extra release artifact, or direct Nakama/Pitaya API compatibility is added by this decision.

## Reversal Conditions

Revisit this decision if early feedback shows that the target segments are wrong, the README or try path is materially misleading, the feedback intake surface must be changed before outreach, or the maintainer authorizes a different discovery strategy.

## Follow-Up

- Prepare the first alpha feedback intake surfaces.
- Keep all feedback records redacted.
- Ask for explicit maintainer authorization before any broad announcement, hosted deployment, or additional release artifact.
