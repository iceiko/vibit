# Conversation: Nakama-First AI-Native Feature Workflow Direction

Date: 2026-05-24
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-23-confirm-next-alpha-direction-after-realtime-outbound-delivery-slice/`

Related artifacts:

- `decisions/ADR-0127-nakama-first-ai-native-requirement-test-workflow-direction.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/reference-game-server-alignment.md`
- `docs/v0.1-alpha-goal.md`
- `docs/product-maturity-milestones.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-146/W-0218` completed the realtime protocol and WebSocket outbound delivery implementation slice. The next-ready item was `W-0219 Confirm next alpha direction after realtime outbound delivery slice`.

The project previously used Nakama and Pitaya together as active references for product capability and game server architecture planning. The maintainer clarified that this is now too broad for the near-term direction.

## Maintainer Narrative

The maintainer said:

```text
那这样吧，现在就是参考那卡马吧，以后再说皮塔亚的事。然后还有一点就是，我们要研发跟测试都是AI原生，也就是用户说了一个需求之后，就AI把事情都帮他办好，这就是我们的这个产品的设计目的，所有的架构跟实现都是为这点服务的。你按这个把计划再调整一下。
```

English summary: use Nakama as the current reference and leave Pitaya for later. Development and testing should be AI-native. After a user states a requirement, AI should handle the work. That is the product design purpose, and architecture plus implementation should serve it.

## Agent Response Summary

The agent completed the current direction-selection slice instead of implementing runtime behavior. The selected next direction is:

```text
define_agent_native_feature_request_test_workflow
```

The change completes `M-147/W-0219`, accepts `ADR-0127`, updates the roadmap posture to Nakama-first, defers Pitaya as a future distributed architecture reference, and opens `M-148/W-0220 Define agent-native feature request and test workflow` as the next-ready item.

## Decisions

- Make Nakama the primary product capability reference.
- Defer Pitaya as a future architecture reference for distributed runtime concerns.
- Treat AI-native development and AI-native testing as the product purpose.
- Make the expected feature workflow:

```text
user requirement -> AI-written bounded requirement spec -> AI-written acceptance criteria -> AI-written test plan -> AI-written or updated tests -> AI implementation -> AI-run verification -> AI-updated docs/manifests/ADRs/change records
```

- Open `W-0220` to define the workflow standard.
- Do not add runtime behavior, protocol changes, generated output, migrations, dependencies, broad product modules, Pitaya-style distributed architecture, hosted deployment, release artifacts, public announcements, paid promotion, or direct compatibility in this slice.

## Reference Basis

Nakama is the current product reference because it provides the broad capability target for an open-source game/backend server framework.

Pitaya is deferred because its strongest current value is distributed architecture vocabulary, and that should not pull cluster/RPC/frontend-backend work into the prototype-ready foundation before the AI-native feature workflow is defined.

Direct Nakama/Pitaya API compatibility remains deferred.

## Artifacts

- `decisions/ADR-0127-nakama-first-ai-native-requirement-test-workflow-direction.md`
- `changes/2026-05-23-confirm-next-alpha-direction-after-realtime-outbound-delivery-slice/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/reference-game-server-alignment.md`
- `docs/v0.1-alpha-goal.md`
- `docs/product-maturity-milestones.md`
- `README.md`
- `AGENTS.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- `W-0220` must decide the required artifacts for AI-native requirement, acceptance, test plan, implementation, verification, and memory records.
- `W-0220` must decide whether all non-trivial user-facing feature requests require tests first, tests alongside implementation, or explicit not-applicable rationale.
- Future work must decide when Pitaya is reactivated for distributed architecture planning, if ever.

## Follow-Up

- Advance `W-0220 Define agent-native feature request and test workflow`.
- Keep Nakama-first product planning and Pitaya-deferred architecture posture visible in future work.

## Redaction Notes

No secrets, GitHub tokens, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, DSNs with credentials, or raw storage object values from a real user are recorded in this conversation log.
