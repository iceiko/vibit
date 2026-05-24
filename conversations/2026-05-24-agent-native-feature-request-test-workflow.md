# Conversation: Agent-Native Feature Request Test Workflow

Date: 2026-05-24
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/`

Related artifacts:

- `docs/agent-native-feature-request-test-workflow.md`
- `docs/agent-native-feature-request-test-workflow.zh-CN.md`
- `decisions/ADR-0128-agent-native-feature-request-test-workflow.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The maintainer asked to continue with the Nakama target. Earlier in the same direction, the maintainer clarified that referencing both Nakama and Pitaya could be counterproductive for the near term, and selected Nakama as the current product reference while deferring Pitaya.

The maintainer also clarified the product purpose: development and testing should be AI-native. A user should state a backend requirement, and AI should turn it into requirements, tests, implementation, verification, and durable project records.

The current next-ready item was `M-148/W-0220 Define agent-native feature request and test workflow`.

## Maintainer Narrative

The maintainer asked to continue toward Nakama:

```text
继续推进，目标nakama
```

Earlier maintainer direction in the same planning thread:

```text
那这样吧，现在就是参考那卡马吧，以后再说皮塔亚的事。然后还有一点就是，我们要研发跟测试都是AI原生，也就是用户说了一个需求之后，就AI把事情都帮他办好，这就是我们的这个产品的设计目的，所有的架构跟实现都是为这点服务的。你按这个把计划再调整一下。
```

English summary: use Nakama as the current reference and leave Pitaya for later. Development and testing should be AI-native. After a user states a requirement, AI should handle the requirement, tests, implementation, verification, and durable records.

## Agent Response Summary

The agent completed `M-148/W-0220` as a standards, artifact, and repository-check slice. The change defines the first agent-native feature request and test workflow and keeps the implementation posture bounded.

No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, startup wiring, distributed runtime, broad product module, or direct Nakama/Pitaya API compatibility is added.

The workflow requires:

- user requirement capture;
- bounded requirement spec;
- Nakama capability mapping;
- acceptance criteria;
- test plan;
- tests or explicit test rationale;
- implementation boundaries;
- verification records;
- durable memory updates.

## Decisions

- Make Nakama the primary product capability reference for this workflow.
- Keep Pitaya deferred as a future distributed architecture reference.
- Require non-trivial user-facing feature work to create a durable change folder with:

```text
request.md
spec.yaml
impact.md
plan.md
checklist.md
verification.md
```

- Require the workflow phases:

```text
user_requirement
requirement_spec
nakama_capability_mapping
acceptance_criteria
test_plan
tests
implementation_boundaries
verification
durable_memory
```

- Require tests or an explicit bounded test rationale for feature work.
- Do not authorize runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, startup wiring, social/realtime product modules, matchmaking, match runtime, Pitaya-style cluster/RPC work, or direct Nakama/Pitaya API compatibility in this slice.
- Open `M-149/W-0221 Pilot Nakama-aligned feature request workflow` as the next-ready item.

## Artifacts

- `docs/agent-native-feature-request-test-workflow.md`
- `docs/agent-native-feature-request-test-workflow.zh-CN.md`
- `decisions/ADR-0128-agent-native-feature-request-test-workflow.md`
- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `modules/storage/module.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
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

- `W-0221` must pilot the workflow against a concrete Nakama-aligned feature request.
- Future work must decide which Nakama-aligned capability family is the best pilot target.
- Future work must decide when, if ever, Pitaya should be reactivated for distributed architecture planning.

## Follow-Up

- Advance `M-149/W-0221 Pilot Nakama-aligned feature request workflow`.
- Keep Nakama-first capability mapping visible in future feature specs.
- Keep the AI-native requirement, test, implementation, verification, and memory workflow central to architecture planning.

## Redaction Notes

No secrets, GitHub tokens, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, DSNs with credentials, or raw storage object values from a real user are recorded in this conversation log.
