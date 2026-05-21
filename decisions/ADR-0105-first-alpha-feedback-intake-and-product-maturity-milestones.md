# ADR-0105: First Alpha Feedback Intake And Product Maturity Milestones

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-prepare-first-alpha-feedback-intake-surfaces/`

Related conversations:

- `conversations/2026-05-22-first-alpha-feedback-intake-surfaces.md`

Related artifacts:

- `docs/first-alpha-feedback-intake-surfaces.md`
- `docs/first-alpha-feedback-intake-surfaces.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `.github/ISSUE_TEMPLATE/alpha-feedback.yml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`v0.1.0-alpha.1` is live as a source-first alpha. The repository has a real authenticated gameplay loop, but it is still not a production server distribution. After the agent explained the current maturity state, the maintainer asked to turn that judgment into durable milestones and continue advancing toward a real product and productivity foundation.

The immediate queued work item, `W-0197`, is to prepare feedback intake surfaces before any broader outreach. This is the correct place to connect early feedback with the larger maturity path: source alpha, prototype-ready foundation, single-node production-candidate foundation, and Nakama/Pitaya-class product.

## Decision

Prepare the first feedback intake surface as a GitHub issue form:

```text
.github/ISSUE_TEMPLATE/alpha-feedback.yml
```

Define the feedback intake standard in `docs/first-alpha-feedback-intake-surfaces.md`.

Define product maturity milestones in `docs/product-maturity-milestones.md`:

- Stage 1: source-first alpha, reached by `v0.1.0-alpha.1`.
- Stage 2: prototype-ready foundation, the next product stage.
- Stage 3: single-node production-candidate foundation, planned.
- Stage 4: Nakama/Pitaya-class product, long-term target.

Set the next direction to:

```text
W-0198 Define prototype-ready foundation execution plan
```

Do not execute public announcements beyond the GitHub release record, run paid promotion, add hosted deployments, create release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, or SDK packages, change runtime behavior, add protocol routes or generated output, add migrations or dependencies, broaden operations/admin behavior, change authentication/session behavior, add broad product modules, or add direct Nakama/Pitaya API compatibility.

The repository check rule for this decision is `runtime.first_alpha_feedback_intake_surfaces`.

## Decision Record

```yaml
first_alpha_feedback_intake_surfaces: prepared
product_maturity_milestones: defined
completed_work_item: W-0197
decision: ADR-0105
check_rule: runtime.first_alpha_feedback_intake_surfaces
release_identifier: v0.1.0-alpha.1
feedback_intake_surface: .github/ISSUE_TEMPLATE/alpha-feedback.yml
stage_1_source_first_alpha: reached
stage_2_prototype_ready_foundation: next_product_stage
stage_3_single_node_production_candidate_foundation: planned
stage_4_nakama_pitaya_class_product: long_term_target
next_direction: prototype_ready_foundation_execution_plan
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
```

## Alternatives Considered

- Wait for organic feedback without a structured issue form.
- Open a GitHub Discussion area first.
- Start public announcements before triage and redaction guidance exist.
- Resume runtime feature work before recording product maturity milestones.
- Treat the current alpha as production-adjacent and ask broadly for production users.

## Rationale

The current alpha is useful for technical inspection, but it needs structured feedback before broader outreach. A GitHub issue form is the smallest repository-owned surface that can collect setup friction, prototype blockers, production-readiness concerns, and product-class gaps without adding external services or release artifacts.

The maturity milestones prevent two common failures: overstating the current source alpha, and losing the long-term product ambition. They let maintainers triage feedback by stage while still keeping the next product step concrete: move toward prototype-ready foundation.

Nakama and Pitaya remain product-class reference baselines, but this decision does not add API compatibility or product modules.

## Agent Reasoning Summary

The agent treated `W-0197` as a developer-experience and product-roadmap slice. The work prepares a feedback issue form, records redaction-safe triage, preserves all release/runtime/protocol deferrals, and turns the maintainer's product-stage intent into durable milestones.

## Decision Weights

```yaml
decision_weights:
  user_feedback_capture: high
  redaction_safety: high
  product_stage_clarity: high
  prototype_ready_momentum: high
  runtime_scope_change: none
  artifact_scope_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0197` is completed as the durable first alpha feedback intake surface record.
- Early feedback now has a structured GitHub issue form.
- Feedback is mapped to product maturity buckets.
- The repository now records source alpha, prototype-ready foundation, single-node production-candidate foundation, and Nakama/Pitaya-class product stages.
- The next ready work item becomes `W-0198 Define prototype-ready foundation execution plan`.
- No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, hosted deployment, operations/admin expansion, authentication/session behavior change, broad product module, extra release artifact, public announcement, paid promotion, or direct Nakama/Pitaya API compatibility is added by this decision.

## Reversal Conditions

Revisit this decision if early feedback shows the issue form is too rigid, GitHub Discussions are a better first surface, the maturity buckets confuse users, or the next product stage should be changed before the prototype-ready execution plan.

## Follow-Up

- Define the prototype-ready foundation execution plan.
- Keep feedback redaction-safe.
- Convert actionable feedback into bounded work items.
- Ask for explicit maintainer authorization before broad outreach, hosted deployment, release artifacts, or direct compatibility claims.
