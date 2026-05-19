# Conversation: Nakama And Pitaya Product Parity Roadmap

Date: 2026-05-18
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-implementation/`
- `changes/2026-05-18-ratify-nakama-pitaya-product-parity-roadmap/`

Related artifacts:

- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- `decisions/ADR-0078-nakama-pitaya-product-parity-roadmap.md`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The project had just completed the single-process WebSocket close policy implementation and was blocked at `M-095/W-0167`, a next-direction confirmation gate.

## Maintainer Narrative

The maintainer clarified in Chinese that the previous goal summaries were correct but incomplete:

> 我们是要做和nakama/pitaya同一级别的产品，要覆盖它们所有常用的功能。

The maintainer then requested:

> 按照刚才的这个重新规划项目开发的路线和方式。然后继续推进。

## Agent Response Summary

The agent treated the clarification as a product-scope change, not as a request to immediately implement a large feature module. The agent selected the existing `expand_core_game_backend_modules_after_nakama_pitaya_review` direction and implemented it as a precise roadmap ratification slice.

The roadmap defines Nakama/Pitaya-class common capability coverage as a product target while preserving vibit's agent-native constraints and non-compatibility posture.

## Decisions

- Product parity with common Nakama/Pitaya capability families is now an explicit target.
- Direct Nakama/Pitaya API compatibility remains out of scope unless a future compatibility ADR selects it.
- The next concrete implementation direction should finish runtime lifecycle closure before presence, chat, social, matchmaking, or match runtime expansion.
- The recommended next direction is `define_protocol_logout_route_gate`.

## Artifacts

- Added `docs/nakama-pitaya-product-parity-roadmap.md`.
- Added `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`.
- Added `ADR-0078`.
- Updated architecture manifests, rules, checks, and work queue state.

## Open Questions

- Which client SDK should be first when the server lifecycle and first product modules are stable.
- Whether admin console should begin as CLI inspection, HTTP admin API, or web UI.
- How much direct compatibility, if any, should ever be offered for existing Nakama/Pitaya users.

## Follow-Up

- Define protocol logout route behavior.
- Define concrete transport close handoff.
- Define reconnect and connection epoch behavior.
- Define protocol session carrier behavior.
- Start presence only after the connection/session/logout foundation is stable.

## Redaction Notes

No secrets, tokens, account details, or private user data were included.
