# ADR-0107: Prototype-Ready Local Development Path Gate

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-prototype-ready-local-development-path-gate/`

Related conversations:

- `conversations/2026-05-22-prototype-ready-local-development-path-gate.md`

Related artifacts:

- `docs/prototype-ready-local-development-path-gate.md`
- `docs/prototype-ready-local-development-path-gate.zh-CN.md`
- `docs/prototype-ready-foundation-execution-plan.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0106` selected the prototype-ready local development path gate as the first execution slice after the source-first alpha and feedback intake surface. The maintainer's product direction is to move vibit beyond an inspectable alpha skeleton toward a foundation developers can actually use for prototypes.

The current local path has important pieces: README, runtime runbook, status endpoints, a redacted request-loop script, repository checks, PostgreSQL migrations, and an authenticated gameplay E2E proof. The gap is that setup, configuration, migration, startup, and example-flow expectations are not yet packaged as one repeatable prototype-author path.

## Decision

Define `docs/prototype-ready-local-development-path-gate.md` as the gate for the local development path.

Select the next implementation slice as:

```text
W-0200 Implement prototype-ready local development path package
```

The next slice may package the local path using docs, scripts, examples, placeholder templates, static checks, and focused verification of existing behavior. It must ask first before changing production runtime behavior, protocol source, generated output, SQL migrations, repository interfaces, dependencies, release artifacts, hosted deployment surfaces, broad operations/admin behavior, authentication/session semantics, broad product modules, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

The repository check rule for this decision is:

```text
runtime.prototype_ready_local_development_path_gate
```

## Decision Record

```yaml
prototype_ready_local_development_path_gate: defined
completed_work_item: W-0199
decision: ADR-0107
check_rule: runtime.prototype_ready_local_development_path_gate
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
source_execution_plan: docs/prototype-ready-foundation-execution-plan.md
gate_standard: docs/prototype-ready-local-development-path-gate.md
future_implementation_work_item: W-0200
future_implementation_direction: prototype_ready_local_development_path_package
supported_prerequisites_recorded: true
startup_expectations_recorded: true
migration_expectations_recorded: true
configuration_secret_posture_recorded: true
example_flow_shape_recorded: true
allowed_future_write_areas_recorded: true
verification_expectations_recorded: true
stop_conditions_recorded: true
planning_only: true
gate_only: true
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

## Alternatives Considered

- Jump directly into runtime behavior changes for a new local onboarding or setup route.
- Add Docker Compose or container packaging as the first local path.
- Add a hosted demo before local setup is repeatable.
- Treat the existing E2E request-loop script as enough for prototype-ready developer experience.
- Move directly to storage objects or realtime messaging before reducing setup friction.

## Rationale

The highest-leverage prototype-ready step is to make the existing source-first runtime easier to start, configure, migrate, and exercise locally. That work should happen before broadening product capabilities because every later feature will depend on developers being able to run a trustworthy local path.

The gate keeps W-0200 productive but bounded: documentation, scripts, examples, static checks, and verification packaging are allowed; production runtime semantics, protocol, migrations, dependencies, release surfaces, hosted deployments, and compatibility promises remain outside scope.

## Agent Reasoning Summary

The agent treated `W-0199` as a gate-only slice. The decision narrows the next implementation to local development ergonomics while preserving all previously deferred product and release surfaces. The next step can improve how users experience the current alpha without claiming production readiness or changing server behavior.

## Decision Weights

```yaml
decision_weights:
  developer_activation: high
  prototype_usefulness: high
  setup_friction_reduction: high
  runtime_scope_change: none
  artifact_scope_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0199` completes as the gate for local development path packaging.
- `W-0200` becomes the next ready work item.
- The next implementation should be a source-first local path package, not a runtime feature expansion.
- Storage objects, realtime messaging/server push, failure/concurrency verification, and operations inspection remain planned later families.
- No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, hosted deployment, operations/admin expansion, authentication/session behavior change, broad product module, extra release artifact, public announcement, paid promotion, or direct Nakama/Pitaya API compatibility is added by this decision.

## Reversal Conditions

Revisit this decision if early feedback shows that a different first prototype-ready blocker is more urgent than local setup and example ergonomics, or if the maintainer explicitly authorizes runtime, protocol, migration, dependency, release, hosted deployment, or product-scope work before the local development path package.

## Follow-Up

- Implement the prototype-ready local development path package in `W-0200`.
- Keep the package source-first and local.
- Preserve the secret redaction posture.
- Ask before any runtime, protocol, migration, dependency, release artifact, hosted deployment, public outreach, broad product, authentication/session, or direct compatibility expansion.
