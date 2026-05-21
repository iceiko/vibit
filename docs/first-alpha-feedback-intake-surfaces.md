# First Alpha Feedback Intake Surfaces

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Repository-owned intake surfaces for early alpha user feedback
Depends on: `docs/first-alpha-user-discovery-loop.md`, `docs/product-maturity-milestones.md`
Canonical decision: `ADR-0105`

The paired Simplified Chinese translation is `docs/first-alpha-feedback-intake-surfaces.zh-CN.md`. The English file is authoritative.

This document prepares the first repository-owned feedback intake surface for `v0.1.0-alpha.1`. It records the intake format, feedback fields, redaction expectations, triage labels, product maturity mapping, and stop conditions before any broader outreach. It does not execute public announcements beyond the GitHub release record, run paid promotion, add hosted deployments, create release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, or SDK packages, change runtime behavior, add protocol routes or generated output, add migrations or dependencies, broaden operations/admin behavior, change authentication/session behavior, add broad product modules, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The feedback intake record is:

```yaml
first_alpha_feedback_intake_surfaces: prepared
completed_work_item: W-0197
decision: ADR-0105
check_rule: runtime.first_alpha_feedback_intake_surfaces
release_identifier: v0.1.0-alpha.1
source_first_alpha_release_available: true
user_discovery_loop_defined: true
product_maturity_milestones_defined: true
feedback_intake_surface: .github/ISSUE_TEMPLATE/alpha-feedback.yml
feedback_intake_standard: docs/first-alpha-feedback-intake-surfaces.md
feedback_intake_standard_translation: docs/first-alpha-feedback-intake-surfaces.zh-CN.md
issue_template_added: true
feedback_fields_recorded: true
redaction_expectations_recorded: true
triage_labels_recorded: true
product_maturity_mapping_recorded: true
next_direction: prototype_ready_foundation_execution_plan
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
additional_release_artifacts_authorized: false
hosted_deployment_added: false
runtime_behavior_added: false
authentication_session_behavior_changed: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
product_module_expansion_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Intake Surface Choice

The first intake surface is a GitHub issue form:

```text
.github/ISSUE_TEMPLATE/alpha-feedback.yml
```

Rationale:

- It is repository-owned and visible to contributors.
- It works without adding a hosted deployment, binary package, SDK, container, or external service.
- It can guide users away from secret leakage.
- It can map feedback directly to the product maturity stages.
- It gives maintainers a stable place to convert friction into bounded work items.

GitHub Discussions may be added later, but that is not required for the first feedback surface. A discussion forum would be useful for open-ended comparison or design feedback, but the first loop needs structured reports more than conversation volume.

## 3. Intake Format

Each alpha feedback issue should capture:

- developer segment;
- attempted action;
- outcome;
- first blocker;
- product maturity bucket;
- commands or pages involved, redacted as needed;
- expected next useful step;
- whether the reporter checked the redaction guidance.

The issue template intentionally avoids asking for logs by default. If logs are necessary, ask the reporter to remove secrets before sharing anything.

## 4. Redaction Expectations

Never ask users to paste:

- raw device credentials;
- raw access tokens;
- verifier key values;
- PostgreSQL DSNs with credentials;
- lookup digests or verifier digests;
- HMAC input or output;
- headers;
- cookies;
- query strings;
- WebSocket subprotocol values;
- remote addresses;
- private `.env`, shell profile, or local environment files;
- GitHub tokens or any other API tokens.

Safe snippets should be shortened and redacted. Prefer command names, error class names, file paths, and documentation page names over full environment output.

## 5. Triage Labels

Recommended labels for feedback issues:

```text
alpha-feedback
source-alpha-friction
prototype-ready-gap
production-candidate-gap
product-class-gap
redaction-review
needs-reproduction
candidate-work-item
deferred
```

These labels are documented, but this work item does not require creating repository labels through GitHub API calls. Label creation can be done later by a maintainer or a separate repository-admin work item.

## 6. Product Maturity Mapping

Map incoming feedback to one maturity bucket:

- `source_alpha_friction`: the reporter cannot understand, clone, check, test, run the request loop, or follow the current local alpha path.
- `prototype_ready_gap`: the reporter can inspect the alpha but cannot build a serious prototype because a shared online-service, example, client, setup, or local development capability is missing.
- `production_candidate_gap`: the reporter is blocked by packaging, security, operations, migration, observability, performance, reliability, or failure-mode concerns.
- `product_class_gap`: the reporter needs broader product capability such as social systems, chat, leaderboards, matchmaking, match runtime, SDKs, admin console, extensibility, or distributed runtime.
- `out_of_scope_for_now`: the request is valid but cannot be acted on without later authorization or a later maturity stage.

This mapping keeps early feedback from being flattened into an unordered feature list. The next planned direction is `prototype_ready_foundation_execution_plan`, so prototype-blocking feedback should receive special attention without pretending the current release is production-ready.

## 7. Triage Posture

Use this default triage flow:

1. Confirm the issue does not contain secrets.
2. Assign a product maturity bucket.
3. Identify whether the report is setup friction, concept confusion, command failure, runtime gap, missing capability, trust gap, packaging gap, or roadmap question.
4. Reproduce locally when the report contains a concrete command failure.
5. Convert actionable findings into a bounded work item or explicitly defer them.
6. Link the issue to the relevant doc, ADR, change spec, or future work item.

Do not promise production readiness, direct Nakama/Pitaya compatibility, hosted demos, SDK availability, binary/container publication, or broad feature parity in issue replies unless a later explicit work item authorizes that claim.

## 8. Stop Conditions

Stop and ask for maintainer authorization if feedback handling would require:

- public announcements beyond the GitHub release record;
- paid promotion;
- hosted deployment or demo;
- release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, or SDK packages;
- runtime behavior changes;
- protocol route changes;
- Protobuf source or generated output changes;
- migrations;
- dependencies;
- broad operations/admin behavior;
- authentication/session behavior changes;
- broad product module expansion;
- direct Nakama/Pitaya API compatibility;
- disclosure or handling of secrets.

## 9. Next Work

The next bounded direction is:

```text
W-0198 Define prototype-ready foundation execution plan
```

That work should take the maturity milestones seriously and define the smallest path from source-first alpha toward prototype-ready foundation. It may use early feedback if available, but it should not wait indefinitely before defining the first execution plan.

## 10. Verification

The repository should verify:

- the issue form exists;
- this standard and its translation exist;
- product maturity milestones are recorded;
- `ADR-0105` records the decision;
- `.arch` manifests point to `W-0198`;
- runtime, protocol, generated output, migration, dependency, release artifact, hosted deployment, announcement, paid promotion, authentication/session, broad product, and direct compatibility deferrals remain preserved.
