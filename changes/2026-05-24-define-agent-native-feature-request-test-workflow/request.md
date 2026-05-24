# Request

## Original Request

```text
继续推进，目标nakama
```

The preceding maintainer direction selected Nakama as the current reference target and clarified that vibit's architecture should serve AI-native development and AI-native testing: a user states a backend requirement, and AI handles the rest of the development and verification workflow.

## Clarified Requirement

Complete `M-148/W-0220 Define agent-native feature request and test workflow`.

Define the repository standard for turning user backend requirements into bounded specs, Nakama capability mapping, acceptance criteria, test plans, tests or explicit test rationale, implementation boundaries, verification records, and durable project memory.

## User-Visible Outcome

Future Nakama-aligned feature work has a clear workflow:

```text
user requirement
-> requirement spec
-> Nakama capability mapping
-> acceptance criteria
-> test plan
-> tests
-> implementation boundaries
-> verification
-> durable memory
```

## Non-Goals

- Add runtime behavior.
- Add protocol routes or Protobuf messages.
- Add generated output.
- Add migrations or dependencies.
- Add startup wiring.
- Add chat, groups, matchmaking, match runtime, leaderboards, economy, operations, SDKs, or other broad product modules.
- Add Pitaya-style distributed architecture, RPC, frontend/backend roles, service discovery, groups, or broadcast fanout.
- Add hosted deployments, release artifacts, public announcements, or paid promotion.
- Add direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- `docs/agent-native-feature-request-test-workflow.md` and the Simplified Chinese translation define the workflow.
- `ADR-0128` records the decision.
- The change spec and conversation log preserve maintainer intent.
- `.arch/` manifests and relevant guides record `W-0220` completion and `W-0221` as next-ready.
- Repository checks include `runtime.agent_native_feature_request_test_workflow`.
- The workflow requires specs, acceptance criteria, test plans, tests or explicit test rationale, implementation boundaries, verification, and durable memory for non-trivial user-facing feature work.
- The workflow keeps Nakama as the primary product reference and Pitaya deferred.
- No runtime/protocol/generated/migration/dependency/direct compatibility scope is added.

