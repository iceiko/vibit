# ADR-0108: Prototype-Ready Local Development Path Package

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-prototype-ready-local-development-path-package/`

Related conversations:

- `conversations/2026-05-22-prototype-ready-local-development-path-package.md`

Related artifacts:

- `docs/prototype-ready-local-development-path-package.md`
- `docs/prototype-ready-local-development-path-package.zh-CN.md`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `examples/local.prototype.env.example`
- `.gitignore`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0107` defined the gate for a repeatable source-first local development path. The current alpha already has a real authenticated gameplay proof, status endpoints, repository checks, Go tests, runtime runbook, and release records, but the practical setup path needed a single prototype-author package.

The maintainer's current direction is to move vibit from a source-first alpha toward a foundation developers can use for real prototypes. That requires better setup, configuration, migration, redaction, example-flow, and verification ergonomics before broadening product capabilities.

## Decision

Implement the prototype-ready local development path package as docs, examples, placeholder local configuration, `.gitignore` guardrails, static checks, and contributor entry-point updates.

The package records:

- a quick source checkout verification path;
- supported local prerequisites;
- private local configuration posture;
- a redacted placeholder env template;
- explicit migration and startup expectations;
- the authenticated request-loop proof;
- local status surfaces;
- verification commands;
- stop conditions.

The repository check rule for this decision is:

```text
runtime.prototype_ready_local_development_path_package
```

The next bounded work item is:

```text
W-0201 Define storage objects behavior gate
```

## Decision Record

```yaml
prototype_ready_local_development_path_package: implemented
completed_work_item: W-0200
decision: ADR-0108
check_rule: runtime.prototype_ready_local_development_path_package
gate_decision: ADR-0107
gate_standard: docs/prototype-ready-local-development-path-gate.md
package_standard: docs/prototype-ready-local-development-path-package.md
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
package_scope: docs_scripts_examples_static_checks
quick_source_path_recorded: true
prerequisite_check_recorded: true
redacted_local_configuration_template_added: true
gitignore_local_secret_guard_added: true
migration_expectations_packaged: true
runtime_startup_guidance_packaged: true
example_flow_script_recorded: true
request_loop_proof_recorded: true
verification_commands_recorded: true
stop_conditions_recorded: true
next_work_item: W-0201
next_direction: storage_objects_behavior_gate
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

- Add Docker Compose or a containerized quickstart first.
- Add a public local onboarding protocol route.
- Add a CLI that creates local users and credentials.
- Implement storage objects immediately before reducing local path friction.
- Treat the existing request-loop script as enough without a package standard.

## Rationale

The package should reduce friction without changing runtime semantics. A source-first path is the right near-term tradeoff because it lets prototype authors inspect, run, and verify the current alpha while preserving explicit boundaries for runtime behavior, protocol, migrations, dependencies, release artifacts, hosted deployment, and broader product scope.

The next product capability should be storage objects, but only after a gate defines behavior, ownership, permissions, protocol, data, compatibility, and verification boundaries.

## Agent Reasoning Summary

The agent treated `W-0200` as a package implementation inside the W-0199 gate. The change records user-facing local setup and verification ergonomics, adds a redacted env template and `.gitignore` guardrails, and advances the queue to a storage-objects behavior gate without implementing that product module.

## Decision Weights

```yaml
decision_weights:
  developer_activation: high
  setup_friction_reduction: high
  prototype_usefulness: high
  runtime_scope_change: none
  artifact_scope_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0200` completes as the first prototype-ready local development package.
- Contributor entry points now link the package standard and examples.
- Private local env files are ignored by default.
- The next work item becomes `W-0201 Define storage objects behavior gate`.
- No production runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, hosted deployment, broad operations/admin surface, authentication/session semantics, broad product module, release artifact, public announcement, paid promotion, or direct Nakama/Pitaya API compatibility is added.

## Reversal Conditions

Revisit this decision if early users cannot follow the source-first local path, if the package accidentally hides manual PostgreSQL or secret setup requirements, or if the maintainer explicitly authorizes runtime, protocol, migration, dependency, release, hosted deployment, public outreach, or product-scope expansion.

## Follow-Up

- Define the storage objects behavior gate in `W-0201`.
- Keep local development docs and examples redacted.
- Preserve ask-first boundaries before adding protocol, data, dependencies, release artifacts, hosted deployment, or broader product capability.
