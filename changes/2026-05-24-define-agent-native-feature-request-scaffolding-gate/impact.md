# Impact

## Affected Areas

```yaml
affected_modules:
  - runtime
  - reference
  - repository_workflow
  - developer_experience
selected_direction: define_agent_native_feature_request_scaffolding_gate
next_direction: implement_agent_native_feature_request_scaffolding
future_work_item: W-0229
decision: ADR-0136
check_rule: runtime.agent_native_feature_request_scaffolding_gate
```

## Module Ownership

The gate is repository-level workflow and tooling posture. It does not assign ownership to a game-domain module and does not change runtime module boundaries.

Future ownership candidates:

- `tools/vibit` for a narrow scaffold command;
- `changes/_template/feature-request/` for source-first templates;
- repository checks for template/scaffold shape validation.

## Public Contract Impact

No public commands, queries, events, permissions, errors, Protobuf messages, or protocol routes are added or changed.

## Data And Migration Impact

No schema, migration, persistence, repository interface, adapter, or storage ownership changes are made.

## Runtime Impact

No runtime behavior changes are made. The gate does not wire startup, change authentication/session behavior, add delivery guarantees, add route protection, or change any protocol-visible path.

## Developer Experience Impact

The gate creates a durable standard for future feature request scaffolding. The intended future effect is to make the product loop executable:

```text
user requirement
-> request.md
-> spec.yaml
-> acceptance criteria
-> test plan
-> implementation boundaries
-> verification.md
-> durable memory
```

## Nakama/Pitaya Alignment

Nakama remains the primary product reference for capability planning. This gate supports Nakama-class growth by making feature intake and tests explicit before broadening modules.

Pitaya remains deferred as a future distributed architecture reference. This gate does not introduce frontend/backend roles, RPC, service discovery, groups, cluster routing, or distributed session behavior.

## Compatibility Risks

Compatibility risk is low because this slice is gate-only and text/check driven. The main risk is future overreach: treating this gate as permission to generate runtime code or direct compatibility shims. The gate and check rule explicitly reject that.

