# ADR-0133: Local Alpha Example Client Path Gate

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-define-local-alpha-example-client-path-gate/`

Related conversations:

- `conversations/2026-05-24-local-alpha-example-client-path-gate.md`

Related artifacts:

- `docs/local-alpha-example-client-path-gate.md`
- `docs/local-alpha-example-client-path-gate.zh-CN.md`
- `docs/agent-native-feature-request-test-workflow.md`
- `decisions/ADR-0132-select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof.md`
- `examples/README.md`
- `examples/local-alpha-request-loop.sh`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0132` selected `client_sdks_examples_and_developer_experience` as the next Nakama prototype-ready capability family. The current source-first alpha proves a meaningful local backend loop, but the most visible proof is still a shell wrapper over internal E2E tests.

The next implementation needs a clearer example-client or example-app path, but current boundaries matter: local onboarding is not a public client route, generated Protobuf output is not a public client package, and vibit does not yet publish an SDK. A live external client or generated client library would therefore require broader product and protocol decisions than this step should make.

## Decision

Accept `docs/local-alpha-example-client-path-gate.md` as the gate for the first clearer source-first local alpha example client path.

The selected first posture is:

```text
source_first_local_alpha_example_client_path
```

The future implementation candidates are:

```text
examples/local-alpha-client/README.md
examples/local-alpha-example-client.sh
runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
```

Open:

```text
M-154/W-0226 Implement local alpha example client path
```

as the next-ready work item.

This decision completes `M-153/W-0225`. It does not implement an example client, publish an SDK, generate client libraries, add dependencies, add runtime behavior, add protocol routes, add Protobuf source, change generated output, add migrations, add persistence, change startup wiring, change authentication/session behavior, add hosted deployments, create release artifacts, add Pitaya-style distributed architecture, or add direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement a live WebSocket client immediately.
- Publish a first SDK or generated client package.
- Add a public local onboarding route first.
- Add chat, groups, matchmaking, match runtime, or operations/admin work first.
- Keep relying only on the existing request-loop shell script.
- Reactivate Pitaya distributed architecture review.

## Rationale

A clearer example path is the highest-leverage next step because it makes the existing Nakama-style alpha foundations legible before the project broadens. A live external client would be premature because key client-facing pieces are intentionally not public yet. A source-first example path can still show the real local alpha loop while keeping redaction, internal package boundaries, and no-SDK posture intact.

The gate keeps implementation small enough for `W-0226`: docs and a wrapper entrypoint under `examples/`, plus optional focused runtime proof work under existing Protobuf E2E ownership. This aligns with the AI-native requirement/spec/test/implementation/verification workflow without expanding runtime behavior.

## Agent Reasoning Summary

The maintainer asked to keep advancing toward Nakama and to preserve AI-native development and testing as the product purpose. The current work item is a gate, so the correct action is to define allowed files, redaction rules, accepted existing routes, verification commands, and stop conditions before implementation. Pitaya remains deferred.

## Decision Weights

```yaml
decision_weights:
  nakama_developer_experience_alignment: high
  prototype_ready_value: high
  implementation_boundedness: high
  redaction_safety: high
  dependency_risk: none_in_this_step
  runtime_behavior_risk: none_in_this_step
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `docs/local-alpha-example-client-path-gate.md` and its Simplified Chinese translation are accepted.
- `runtime.local_alpha_example_client_path_gate` becomes the repository check rule for this gate.
- `M-153/W-0225` is completed.
- `M-154/W-0226 Implement local alpha example client path` becomes next-ready.
- The first example path remains source-first and local to the repository.
- SDK publication, generated client libraries, live external client guarantees, new protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, authentication/session behavior changes, hosted deployments, release artifacts, chat, groups, matchmaking, match runtime, operations/admin behavior, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a maintainer explicitly authorizes SDK publication before the example path implementation;
- public onboarding route work becomes a prerequisite for a useful example;
- alpha feedback shows operations/admin inspection is a stronger blocker than example readability;
- a later ADR chooses a public client package boundary or generated client output location.

## Follow-Up

- Complete `W-0226`: implement the local alpha example client path within the gate.
- Keep public SDKs, generated clients, hosted demos, protocol changes, runtime changes, dependencies, distributed runtime, and direct compatibility behind later bounded work items.
