# Prototype-Ready Foundation Execution Plan

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Execution plan from source-first alpha to prototype-ready vibit foundation
Depends on: `docs/product-maturity-milestones.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/first-alpha-feedback-intake-surfaces.md`, `docs/alpha-developer-flow.md`
Canonical decision: `ADR-0106`

The paired Simplified Chinese translation is `docs/prototype-ready-foundation-execution-plan.zh-CN.md`. The English file is authoritative.

This document defines the first execution plan for moving vibit from `v0.1.0-alpha.1` source-first alpha toward a prototype-ready game/backend foundation. It is a planning artifact. It does not implement runtime behavior, add protocol routes, add Protobuf source or generated output, add migrations, add dependencies, broaden operations/admin behavior, add hosted deployments, create release artifacts, run public announcements, run paid promotion, change authentication/session behavior, add broad product modules, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The prototype-ready execution record is:

```yaml
prototype_ready_foundation_execution_plan: defined
completed_work_item: W-0198
decision: ADR-0106
check_rule: runtime.prototype_ready_foundation_execution_plan
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
product_maturity_milestones_standard: docs/product-maturity-milestones.md
execution_plan_standard: docs/prototype-ready-foundation-execution-plan.md
execution_plan_standard_translation: docs/prototype-ready-foundation-execution-plan.zh-CN.md
recommended_sequence_recorded: true
candidate_work_items_recorded: true
maturity_stage_mapping_recorded: true
nakama_pitaya_capability_mapping_recorded: true
success_criteria_recorded: true
stop_conditions_recorded: true
selected_first_execution_slice: prototype_ready_local_development_path_gate
next_work_item: W-0199
next_direction: prototype_ready_local_development_path_gate
planning_only: true
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Interpretation

`v0.1.0-alpha.1` proves that vibit is a real source-first alpha: it has a local authenticated request loop, PostgreSQL-backed state, WebSocket + Protobuf transport, local onboarding, protected inventory and presence paths, runtime sessions, connection binding, logout, runbooks, checks, release notes, and feedback intake.

Prototype-ready means something narrower and more useful than production-ready:

- a developer can start from the repository and build a serious small multiplayer or realtime backend prototype;
- the local setup, migration, configuration, and example path are clear enough to repeat;
- at least one shared online-service capability exists beyond the current proof slices;
- runtime lifecycle behavior has enough verification that prototype authors can trust the core loop;
- missing production concerns are visible instead of accidental.

The next stage is still source-first and may remain single-node. It should feel useful, not merely inspectable.

## 3. Recommended Sequence

The first prototype-ready sequence is:

1. Complete the local development path gate.
2. Implement the local development path and richer example flow after the gate authorizes exact files and behavior.
3. Define the first general storage-object behavior beyond inventory.
4. Implement the smallest storage-object functional slice.
5. Define the first server push, stream, broadcast, or realtime messaging vocabulary.
6. Implement the smallest realtime messaging or server-push functional slice.
7. Strengthen concurrency and failure-path verification around login, protected requests, logout, reconnect, and database failures.
8. Define and add the minimum operations inspection surface needed for serious prototype use.
9. Run a prototype-ready exit review against external feedback and the maturity criteria.

The sequence is intentionally ordered around developer usefulness. A shared service such as storage or messaging is much more valuable once the project has a repeatable local development path and an example that exercises more than one isolated request.

## 4. First Execution Slice

The first selected slice is:

```text
W-0199 Define prototype-ready local development path gate
```

Purpose: define the gate for a repeatable local development path that can later make setup, migrations, configuration, and example flow credible for prototype authors.

The gate should record:

- supported local prerequisites;
- startup and migration expectations;
- local secret/configuration expectations;
- example client or example app shape;
- which later runtime, script, docs, or test files may be changed;
- verification expectations;
- stop conditions before any behavior, protocol, migration, dependency, release, or hosted deployment expansion.

This plan chooses the gate first because developer trust starts with a reproducible path. If a user cannot reliably start the server, run migrations, configure secrets safely, and run a meaningful example, higher-level capabilities such as storage objects, push, chat, or matchmaking remain hard to evaluate.

## 5. Candidate Work Families

The prototype-ready execution plan may open later bounded work items in these families:

- `prototype_ready_local_development_path`: local setup, migration, configuration, and example ergonomics.
- `storage_objects_and_durable_game_state`: general durable object behavior beyond inventory.
- `server_push_streams_or_realtime_messaging`: first outbound realtime vocabulary.
- `failure_and_concurrency_verification`: tests for lifecycle and database edge cases.
- `minimal_operations_inspection`: local inspection for players, sessions, tokens, active connections, and runtime state.
- `feedback_triage_to_work_items`: conversion of real alpha feedback into bounded work.

Each family must still follow vibit's normal order: requirement, gate or spec, contract, generated shape where applicable, logic, tests, checks, docs, and memory.

## 6. Maturity Mapping

The plan maps to Stage 2 `prototype_ready_foundation` as follows:

- Setup friction is handled by the local development path gate and later implementation.
- Example/client ergonomics are handled by the same local path, because the example should demonstrate a realistic multi-request flow.
- General durable game state is handled by the storage-object family.
- Realtime usefulness is handled by server push, streams, broadcast, or messaging.
- Trust in the existing core loop is handled by failure and concurrency verification.
- Operational confidence is handled by minimal inspection surfaces.
- User discovery is handled by the existing feedback intake loop and later triage.

The plan does not claim Stage 3 single-node production-candidate readiness. Stage 3 still requires stronger security review, packaging, observability, upgrade posture, operational runbooks, and failure-mode hardening.

## 7. Nakama/Pitaya Mapping

Nakama pressure:

- storage objects and durable game state;
- presence, status, notifications, chat, and streams;
- SDK/example ergonomics;
- operational inspection;
- clear account/session lifecycle.

Pitaya pressure:

- connection lifecycle and session binding;
- handler route clarity;
- push, groups, broadcast, and stream vocabulary;
- Go-first local runtime ergonomics;
- later frontend/backend and RPC topology, after single-process semantics are stable.

This plan uses Nakama and Pitaya as capability baselines only. It does not add direct API compatibility, copied routes, copied schemas, copied clustering behavior, or a compatibility promise.

## 8. Success Criteria

The prototype-ready foundation track is successful when:

- a technically capable external developer can start from source and complete a repeatable local prototype path;
- the example flow demonstrates more than isolated one-off requests;
- at least one shared online-service capability beyond inventory exists or is explicitly selected as the next implementation;
- setup, configuration, migrations, and secret redaction are clear enough for repeated local use;
- failure and concurrency risks in the authenticated request loop are covered by focused tests or recorded as blockers;
- at least one real feedback item is triaged into a bounded work item or explicitly deferred;
- remaining production gaps are visible and stage-mapped.

## 9. Stop Conditions

Stop and ask for maintainer authorization if executing the plan would require:

- runtime behavior changes before a specific implementation work item authorizes them;
- protocol route changes or Protobuf source/generated output changes;
- migrations or repository/storage adapter changes;
- dependencies;
- broad operations/admin behavior;
- authentication/session semantic changes;
- broad product module expansion;
- direct Nakama/Pitaya API compatibility;
- hosted deployments or demos;
- release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, or SDK packages;
- public announcements beyond the GitHub release record;
- paid promotion;
- handling or disclosure of secrets.

## 10. Next Work

The next bounded direction is:

```text
W-0199 Define prototype-ready local development path gate
```

That work should define the gate for making vibit easier to start, configure, migrate, and exercise locally as the first product-useful step toward Stage 2.

## 11. Verification

The repository should verify:

- this execution plan and its translation exist;
- `ADR-0106` records the decision;
- `.arch` manifests mark `W-0198` completed and open `W-0199`;
- README, alpha goal, developer flow, acceptance checklist, AGENTS guides, and product roadmap point at the new next work;
- runtime, protocol, generated output, migration, dependency, operations/admin, authentication/session, product module, hosted deployment, release artifact, public announcement, paid promotion, and direct compatibility deferrals remain preserved.

