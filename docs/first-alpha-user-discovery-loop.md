# First Alpha User Discovery Loop

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Bounded discovery loop for finding early developers after `v0.1.0-alpha.1`
Depends on: `docs/release-execution-final-authorization.md`, `docs/alpha-developer-flow.md`
Canonical decision: `ADR-0104`

The paired Simplified Chinese translation is `docs/first-alpha-user-discovery-loop.zh-CN.md`. The English file is authoritative.

This document defines the first user-discovery loop after the source-first `v0.1.0-alpha.1` release. It records who vibit should try to learn from first, which surfaces may be prepared, what feedback should be captured, what counts as signal, and where the loop must stop. It does not authorize public announcements beyond the existing GitHub release record, paid promotion, hosted deployments, release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publication, runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The user discovery record is:

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

## 2. Discovery Goal

The first loop should answer one question:

```text
Can technically capable developers understand why vibit exists, run the source alpha locally, and tell us which friction or missing capability blocks the next useful step?
```

This is not a growth campaign. The loop is a learning system for turning a small number of real developer reactions into concrete work items.

## 3. Target Developer Segments

Start with developers who are most likely to feel the problem vibit is trying to solve:

- Go game/backend developers evaluating server frameworks.
- Developers who have used or evaluated Nakama, Pitaya, Colyseus, Pomelo, Agones, or custom backend stacks.
- Engineers using AI coding agents on backend codebases and struggling with unsafe or context-heavy changes.
- Open-source contributors interested in explicit architecture, contracts, generated output, ADRs, and repository checks.
- Technical founders or solo developers building multiplayer or realtime backend prototypes who can tolerate source-first alpha setup.

Do not target production buyers, hosted-platform users, SDK-first developers, or teams requiring packaged binaries, container images, commercial support, direct Nakama/Pitaya API compatibility, or production hardening in this first loop.

## 4. Outreach Surfaces

This work item defines surfaces and copy posture only. It does not execute broad public announcements.

Allowed surfaces to prepare or point at:

- `README.md` and `README.zh-CN.md`.
- GitHub repository description and topics, if edited through a later explicit work item.
- GitHub release page for `v0.1.0-alpha.1`.
- `docs/releases/v0.1.0-alpha.1.md`.
- `docs/alpha-developer-flow.md`.
- `docs/alpha-acceptance-checklist.md`.
- A later GitHub issue or discussion intake surface after explicit authorization.
- Small direct maintainer-to-maintainer review requests after explicit authorization and without spam.

Deferred surfaces:

- Broad posts to social media, news aggregators, large community forums, newsletters, or paid promotion.
- Hosted demos.
- Binary, package, container, registry, checksum, signing, or provenance publication.
- Claims of production readiness, feature parity, SDK availability, or direct compatibility.

## 5. Feedback Capture

Every early-user conversation or issue should try to capture these fields:

```yaml
feedback_record:
  source: github_issue_or_direct_review_or_discussion
  developer_segment: go_game_backend | nakama_pitaya_evaluator | ai_agent_backend_engineer | oss_contributor | prototype_builder | other
  attempted_action: read_readme | clone_repo | run_checks | run_go_tests | run_request_loop | inspect_architecture | inspect_runtime | evaluate_contribution
  outcome: completed | blocked | abandoned | unclear
  first_blocker: prerequisite | command_failure | docs_confusion | concept_confusion | runtime_gap | missing_feature | trust_gap | packaging_gap | other
  exact_command_or_page: redacted_if_needed
  redaction_checked: true
  next_action_candidate: docs_fix | check_fix | runbook_fix | issue_template | runtime_work_item | protocol_work_item | roadmap_clarification | defer
```

Never ask users to share raw device credentials, raw access tokens, verifier keys, DSNs with credentials, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private environment files.

## 6. Review Questions

Use these questions when asking for feedback:

- What made vibit worth a look?
- Could you understand the README in under five minutes?
- Which command did you run first, and what happened?
- Did `node tools/vibit check all`, `go test ./...`, or `examples/local-alpha-request-loop.sh` fail?
- Was the agent-native workflow legible, or did it feel like internal process noise?
- Which missing capability would make vibit worth a second try?
- Would you open an issue or contribute a small slice? If not, what blocked that?

## 7. Success Signals

The first loop is successful when at least one of these happens:

- An external developer clones the repository and reports the result.
- A non-maintainer opens a useful issue or discussion.
- A maintainer receives concrete feedback tied to README, runbook, request-loop, setup, architecture clarity, or next capability.
- A failure in the try path is reproduced and converted into a bounded work item.
- A developer can accurately explain vibit's agent-native server framework premise after reading the README.

Secondary signals:

- GitHub stars, forks, or release views.
- Inbound questions about game backend roadmap scope.
- Specific comparison feedback against Nakama, Pitaya, Colyseus, Pomelo, Agones, or custom Go backends.

Do not treat vanity metrics alone as success.

## 8. Stop Conditions

Stop the loop or ask for maintainer authorization if any of these occur:

- A public announcement beyond the GitHub release record would be posted.
- Paid promotion is considered.
- A hosted deployment or demo is proposed.
- A binary, package, container, checksum, signing/provenance artifact, install script, registry publication, or SDK package is requested as part of this loop.
- Feedback requires runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, operations/admin expansion, authentication/session behavior, broad product module expansion, or direct Nakama/Pitaya API compatibility.
- Users report a security or secret-redaction risk.
- README, release notes, runbook, or request-loop instructions are materially wrong.
- The loop starts producing repeated low-signal outreach instead of concrete learning.

## 9. Next Work

The next bounded direction is:

```text
W-0197 Prepare first alpha feedback intake surfaces
```

That work should prepare a minimal, repository-owned place to receive early-user feedback, such as issue/discussion guidance or labels, without executing broad announcements, adding hosted deployment, or changing runtime behavior.

## 10. Reference Alignment

Nakama and Pitaya both benefit from clear contributor and user intake surfaces around their public repositories and documentation. vibit should learn from that posture while keeping this first loop smaller: source-first, honest about alpha limitations, and focused on agent-native maintainability feedback rather than broad feature parity claims.
