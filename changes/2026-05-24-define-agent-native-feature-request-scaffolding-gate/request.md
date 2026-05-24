# Request

## Original Request

Continue toward Nakama as the primary reference. The product purpose is AI-native development and AI-native testing: after a user states a requirement, AI should help write the requirement, tests, implementation, verification, and durable memory. Keep Pitaya deferred for now.

## Clarified Requirement

Define the agent-native feature request scaffolding gate before implementing scaffolding. The gate must make the future source-first feature intake path explicit enough that a later implementation can create the correct change artifacts before coding starts.

## User-Visible Outcome

The repository records a bounded gate for future feature-request scaffolding:

- user requirement intake;
- source-first change artifact creation;
- Nakama capability mapping;
- acceptance criteria;
- test planning;
- implementation boundaries;
- verification expectations;
- durable memory updates.

The next work item becomes `M-157/W-0229 Implement agent-native feature request scaffolding`.

## Non-Goals

- Implementing scaffolding in this slice.
- Adding `tools/vibit scaffold feature` in this slice.
- Adding `changes/_template/feature-request/` in this slice.
- Adding runtime behavior.
- Adding protocol routes, Protobuf source, or generated output.
- Adding migrations, persistence, repository interfaces, adapters, dependencies, or startup wiring.
- Publishing SDKs or generated client libraries.
- Adding hosted demos, release artifacts, public announcements, or paid promotion.
- Adding chat, groups, matchmaking, match runtime, operations/admin behavior, delivery guarantees, stream subscriptions, or broadcast fanout.
- Adding Pitaya-style distributed architecture.
- Adding direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- `docs/agent-native-feature-request-scaffolding-gate.md` and `docs/agent-native-feature-request-scaffolding-gate.zh-CN.md` define the gate.
- `ADR-0136` accepts the gate.
- `runtime.agent_native_feature_request_scaffolding_gate` is registered and checks the gate artifacts.
- `M-156/W-0228` is completed.
- `M-157/W-0229` is opened as the next-ready implementation slice.
- The gate records required future scaffold artifacts: `request.md`, `spec.yaml`, `impact.md`, `plan.md`, `checklist.md`, and `verification.md`.
- The gate preserves Nakama as the primary product reference and Pitaya as a deferred future architecture reference.
- The gate does not add scaffolding implementation, runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, hosted surfaces, SDKs, distributed runtime, or direct compatibility.

