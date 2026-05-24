# Conversation: Agent-Native Feature Request Scaffolding Selection

Date: 2026-05-24

## Context

The maintainer asked the agent to continue toward the Nakama target. Earlier in the thread, the maintainer clarified that Pitaya should be deferred and that vibit's product purpose is AI-native development and AI-native testing: a user requirement should become a spec, acceptance criteria, tests, implementation, verification, and durable project memory through AI assistance.

`M-154/W-0226` implemented the source-first local alpha example client path and opened `M-155/W-0227` as the next selection slice.

## Maintainer Narrative

The maintainer wants Nakama as the primary product reference. Pitaya is deferred. The product purpose remains AI-native development and AI-native testing. The framework architecture and implementation should serve the workflow where users describe backend requirements and AI performs the bounded engineering work, including tests.

## Agent Response Summary

The agent selected `agent_native_requirement_test_implementation_workflow` as the next Nakama-first prototype-ready capability family after the local alpha example path.

The selected follow-up direction is:

```text
define_agent_native_feature_request_scaffolding_gate
```

The selection opens:

```text
M-156/W-0228 Define agent-native feature request scaffolding gate
```

as the next-ready bounded follow-up.

## Decisions

- Accept `ADR-0135`.
- Register `runtime.agent_native_feature_request_scaffolding_selection`.
- Complete `M-155/W-0227`.
- Open `M-156/W-0228 Define agent-native feature request scaffolding gate` as next-ready.
- Keep Nakama as the primary product capability reference.
- Keep Pitaya deferred as a future distributed architecture reference.
- Keep direct Nakama/Pitaya API compatibility out of scope.
- Do not add scaffolding implementation, runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted surfaces, release artifacts, broad chat/social/matchmaking/match runtime behavior, operations/admin behavior, or distributed runtime in this selection slice.

## Artifacts

- `changes/2026-05-24-select-next-nakama-prototype-ready-capability-after-local-alpha-example-client-path/`
- `decisions/ADR-0135-agent-native-feature-request-scaffolding-selection.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Should the future scaffolding implementation be primarily a `tools/vibit` command, template files, documentation, or a mixed approach?
- Should the first scaffolding implementation create draft change artifacts only, or should it also add repository checks for missing required workflow phases?

## Follow-Up

- Complete `W-0228`: define the agent-native feature request scaffolding gate.

## Redaction Notes

No raw access-token, device credential, verifier digest, lookup digest, verifier key, DSN, local secret, local environment value, user private request content beyond the maintainer's explicit product direction, or GitHub token value is recorded here.
