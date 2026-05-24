# Impact

## Scope

This is a planning and continuation-state slice. It selects the next bounded Nakama-first prototype-ready capability after the local alpha example client path.

## Selected Capability

```text
selected_nakama_capability_family: agent_native_requirement_test_implementation_workflow
selected_direction: define_agent_native_feature_request_scaffolding_gate
future_work_item: W-0228
```

## Product Impact

The selected next step strengthens vibit's core product purpose: a user states a backend requirement, and AI carries the work through spec, acceptance criteria, tests, implementation boundaries, verification, and durable memory.

This should precede broad new runtime capability families because the framework's differentiator is not only Nakama-class backend coverage. It is Nakama-class coverage made reliable for AI-native development and testing.

## Runtime Impact

No runtime behavior changes are made.

## Protocol Impact

No protocol routes, Protobuf sources, generated output, or protocol bridge changes are made.

## Persistence Impact

No migrations, repository changes, persistence adapters, or startup migration behavior are added.

## Dependency Impact

No dependencies are added.

## Compatibility Impact

No direct Nakama/Pitaya API compatibility is added. Nakama remains the primary product capability reference. Pitaya remains deferred as a future distributed architecture reference.

## Deferred Work

- The future scaffolding gate.
- Any future scaffolding implementation.
- Runtime behavior, protocol routes, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted deployments, release artifacts, chat, groups, matchmaking, match runtime, operations/admin behavior, distributed runtime, and direct compatibility.
