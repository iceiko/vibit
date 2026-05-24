# ADR-0132: Select Next Nakama Prototype-Ready Capability After Authenticated Failure Path Proof

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof/`

Related conversations:

- `conversations/2026-05-24-select-next-nakama-prototype-ready-capability.md`

Related artifacts:

- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/product-maturity-milestones.md`
- `docs/alpha-developer-flow.md`
- `docs/prototype-ready-local-development-path-package.md`
- `examples/README.md`
- `examples/local-alpha-request-loop.sh`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0131` completed authenticated gameplay failure-path verification and opened `M-152/W-0224` as a selection slice. The local alpha now proves login, connection binding, protected inventory, presence, storage objects, logout, revoked-token rejection, presence offline behavior, and core authentication failure paths.

The next prototype-ready decision should choose the most useful Nakama-first gap without prematurely adding broad social modules, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted releases, or distributed runtime.

The strongest current gap is developer experience: the source-first alpha has a redacted request-loop script and internal E2E tests, but it does not yet have a clearer client-like or small example-app path that lets a developer see how the existing backend capabilities fit together.

Nakama capability family:

```text
client_sdks_examples_and_developer_experience
```

Pitaya remains deferred as a future distributed architecture reference.

## Decision

Select:

```text
define_local_alpha_example_client_path_gate
```

as the next direction.

Open:

```text
M-153/W-0225 Define local alpha example client path gate
```

as the next-ready work item.

The follow-up gate must define the acceptable source-first example-client or example-app shape, redaction rules, test expectations, and stop conditions before implementation. This decision does not implement an example client, publish an SDK, generate client libraries, add dependencies, add protocol routes, change runtime behavior, create release artifacts, add hosted deployments, or add direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Define the minimum operations inspection surface.
- Prove notification or realtime outbound delivery through a new flow.
- Define chat or group messaging.
- Define friends, groups, parties, leaderboards, matchmaking, or match runtime gates.
- Reactivate Pitaya distributed architecture review.

## Rationale

The existing foundation is now useful enough to demonstrate, but the demonstration path is still too internal. A clearer example path will help external developers and future agents evaluate the current Nakama-style foundations before the project broadens into chat, social, competitive, or multiplayer systems.

Operations inspection is important, but it belongs closer to single-node production-candidate maturity. Chat, groups, matchmaking, and match runtime are product-class features, but adding them before a readable client/example path would increase capability breadth while leaving the user evaluation path weak.

## Agent Reasoning Summary

The maintainer asked to keep advancing toward Nakama and to keep AI-native development and testing as the product purpose. After the authenticated failure-path proof, the next highest-leverage step is to make the current source-first alpha easier to understand as a prototype foundation. The selected gate keeps implementation bounded and continues the agent-native requirement/spec/test/verification flow.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  prototype_ready_value: high
  developer_experience_value: high
  implementation_boundedness: high
  reuse_existing_runtime_surface: high
  avoids_premature_product_breadth: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-152/W-0224` is completed.
- `M-153/W-0225` is next-ready.
- The selected capability family is `client_sdks_examples_and_developer_experience`.
- The next work should define a gate for a source-first local alpha example client/app path.
- Runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted deployments, release artifacts, chat, groups, matchmaking, match runtime, operations/admin behavior, distributed runtime, and direct compatibility remain deferred.
- `runtime.next_nakama_prototype_ready_capability_selection` checks the selection records and next-ready state.

## Reversal Conditions

Revisit this decision if:

- external alpha users report operations/admin inspection as a stronger prototype blocker than client/example clarity;
- an explicit maintainer decision authorizes a broader product module before example-client work;
- the future gate finds that a useful example path requires protocol or runtime work that should be split first.

## Follow-Up

- Complete `W-0225`: define the local alpha example client path gate.
- Keep example implementation, SDK publication, generated client libraries, hosted demos, protocol changes, runtime changes, dependencies, distributed runtime, and direct compatibility behind later bounded work items.
