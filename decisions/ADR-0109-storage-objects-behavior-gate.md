# ADR-0109: Storage Objects Behavior Gate

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-behavior-gate/`

Related conversations:

- `conversations/2026-05-22-storage-objects-behavior-gate.md`

Related artifacts:

- `docs/storage-objects-behavior-gate.md`
- `docs/storage-objects-behavior-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0108` packaged the source-first local development path. The next prototype-ready product gap is a general durable game-state surface beyond the module-local inventory proof slice.

Nakama treats storage objects as a common game backend capability. Pitaya reinforces that handler behavior and persistence should remain separated. vibit needs the same class of durable game-state usefulness while preserving explicit ownership, contract-first behavior, generated boundaries, route-scoped permissions, repository adapters, redaction, and bounded work items.

## Decision

Define the storage objects behavior gate before implementation.

The first posture is player-owned small JSON storage objects addressed by:

```text
owner_kind + owner_id + collection + key
```

The gate records ownership, scope/key posture, read/write semantics, permission posture, optimistic conflict semantics, protocol expectations, data expectations, verification expectations, and stop conditions.

The repository check rule for this decision is:

```text
runtime.storage_objects_behavior_gate
```

The next bounded work item is:

```text
W-0202 Define storage objects persistence schema gate
```

## Decision Record

```yaml
storage_objects_behavior_gate: defined
completed_work_item: W-0201
decision: ADR-0109
check_rule: runtime.storage_objects_behavior_gate
source_package_decision: ADR-0108
source_package_standard: docs/prototype-ready-local-development-path-package.md
gate_standard: docs/storage-objects-behavior-gate.md
target_stage: prototype_ready_foundation
reference_capability_family: storage_objects_and_durable_game_state
first_scope_posture: player_owned_small_json_objects
object_identity_tuple: owner_kind_owner_id_collection_key
ownership_posture_recorded: true
scope_key_posture_recorded: true
read_write_semantics_recorded: true
permission_posture_recorded: true
conflict_semantics_recorded: true
protocol_expectations_recorded: true
data_expectations_recorded: true
verification_expectations_recorded: true
stop_conditions_recorded: true
future_schema_gate_work_item: W-0202
future_schema_gate_direction: storage_objects_persistence_schema_gate
gate_only: true
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
repository_interface_changed: false
storage_adapter_changed: false
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

- Implement a storage object module immediately.
- Start with SQL migration source before defining behavior.
- Copy Nakama storage object routes and permission values.
- Use arbitrary binary/blob storage as the first storage object behavior.
- Defer storage objects and start chat, friends, leaderboards, or match runtime first.

## Rationale

General durable game state is the next small product-useful capability after the local development path package. It unlocks prototype behavior without requiring a full social, economy, or match-runtime module.

The behavior gate keeps the first slice narrow: player-owned, small JSON object records, protected by existing authenticated request identity, with version-aware writes. That is broad enough for prototypes and narrow enough to preserve architecture boundaries before schema, repository, protocol, and runtime implementation work.

## Agent Reasoning Summary

The agent treated `W-0201` as a gate-only product capability definition. It selected the smallest useful storage-object posture, mapped it to Nakama/Pitaya reference pressure without copying public APIs, and advanced the queue to a persistence schema gate so the next step remains bounded and reviewable.

## Decision Weights

```yaml
decision_weights:
  prototype_usefulness: high
  product_capability_progress: high
  implementation_scope: none
  protocol_scope_change: none
  migration_scope_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0201` completes as the storage objects behavior gate.
- The first planned storage object posture is player-owned small JSON durable game-state objects.
- The next work item becomes `W-0202 Define storage objects persistence schema gate`.
- No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, repository interface, storage adapter, hosted deployment, broad operations/admin surface, authentication/session behavior, broad module implementation, release artifact, public announcement, paid promotion, or direct Nakama/Pitaya API compatibility is added.

## Reversal Conditions

Revisit this decision if early prototype authors need a different first storage scope, if JSON object payloads are too loose for agent-native verification, if version semantics need a different conflict posture before schema work, or if the maintainer explicitly chooses another shared online service as the next prototype unlock.

## Follow-Up

- Define the storage objects persistence schema gate in `W-0202`.
- Keep protocol, migration, repository, adapter, and runtime implementation behind later bounded work items.
- Preserve ask-first boundaries for direct compatibility, broad public outreach, release artifacts, hosted deployment, and product-scope expansion.
