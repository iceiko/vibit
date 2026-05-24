# Request

Original maintainer direction:

```text
Continue advancing. Target Nakama.
```

Clarified requirement:

Select the next bounded Nakama-first prototype-ready capability after the source-first local alpha example client path, while keeping vibit's product purpose focused on AI-native development and AI-native testing.

User-visible outcome:

Future contributors and agents should know the next safe step after the example path. The selected step should improve the path where a user states a backend requirement and AI turns it into a bounded spec, acceptance criteria, tests, implementation, verification, and durable memory.

Non-goals:

- Implement the selected capability in this selection slice.
- Add runtime behavior.
- Add protocol routes or Protobuf messages.
- Add generated output.
- Add migrations, persistence, or dependencies.
- Publish SDKs or generated client libraries.
- Add hosted demos, release artifacts, public announcements, or paid promotion.
- Add chat, groups, matchmaking, match runtime, broad operations/admin behavior, or distributed runtime.
- Add direct Nakama/Pitaya API compatibility.

Unknowns:

- The exact scaffolding interface for future requirement intake remains to be defined by the follow-up gate.
- Whether the future implementation should be documentation-only, `tools/vibit`-driven, template-driven, or a mix remains deferred.

Acceptance criteria:

- The selected capability family is recorded as `agent_native_requirement_test_implementation_workflow`.
- The selected follow-up direction is `define_agent_native_feature_request_scaffolding_gate`.
- The selection explains why this step should precede broad chat, groups, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted surfaces, distributed runtime, or direct compatibility.
- A bounded follow-up work item is opened as `M-156/W-0228 Define agent-native feature request scaffolding gate`.
- The selection preserves ask-first boundaries for runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted deployments, release artifacts, Pitaya-style distributed architecture, and direct compatibility.
