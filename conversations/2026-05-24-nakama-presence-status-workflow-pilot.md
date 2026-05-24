# Conversation: Nakama Presence Status Workflow Pilot

Date: 2026-05-24
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/`

Related artifacts:

- `decisions/ADR-0129-nakama-aligned-presence-status-workflow-pilot.md`
- `docs/agent-native-feature-request-test-workflow.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-148/W-0220` completed the agent-native feature request and test workflow. The active next-ready item was `M-149/W-0221 Pilot Nakama-aligned feature request workflow`.

The maintainer asked for sustained progress toward Nakama. The workflow pilot needed one concrete Nakama-style product capability request and a bounded follow-up slice.

The repository already had a minimal protected self-presence query, server-owned active connection registry, first-message connection binding, and authenticated local alpha proof path.

## Maintainer Narrative

The maintainer said:

```text
继续推进，目标nakama 推进10步以上，推进10个小时以上，不要停止，我会离开10小时。
```

English summary: continue progressing toward the Nakama target for many steps and keep going while the maintainer is away.

## Agent Response Summary

The agent completed the workflow pilot instead of adding broad runtime behavior directly. The selected concrete requirement is:

```text
Prove self-presence/status through the authenticated local alpha flow, including online after connection binding and offline after close or invalidation.
```

No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, startup wiring, distributed runtime, broad product module, or direct Nakama/Pitaya API compatibility is added in this pilot.

## Decisions

- Select `presence_status_and_notifications` as the Nakama capability family for the first workflow pilot.
- Treat existing presence/status code as foundation to harden, not as a reason to jump to chat, groups, matchmaking, match runtime, or distributed runtime.
- Open `M-150/W-0222 Harden presence status local proof through close and offline cases` as the next-ready item.
- Keep Pitaya deferred as a future distributed architecture reference.
- Keep direct Nakama/Pitaya API compatibility out of scope.

## Artifacts

- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/`
- `decisions/ADR-0129-nakama-aligned-presence-status-workflow-pilot.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `.arch/modules.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `modules/storage/module.yaml`
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

- `W-0222` must decide whether the local alpha proof should be strengthened in `authenticated_gameplay_e2e_test.go`, lower-level connection/presence tests, or both.
- Future work must decide when to add presence subscriptions or status broadcast fanout.
- Future work must decide whether operations/admin inspection or example client ergonomics should follow the presence/status proof hardening slice.

## Follow-Up

- Advance `M-150/W-0222 Harden presence status local proof through close and offline cases`.
- Keep the AI-native requirement, test, implementation, verification, and memory workflow visible in every future Nakama-aligned feature slice.

## Redaction Notes

No secrets, GitHub tokens, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, DSNs with credentials, raw transport close reason text from a real user, or raw storage object values from a real user are recorded in this conversation log.
