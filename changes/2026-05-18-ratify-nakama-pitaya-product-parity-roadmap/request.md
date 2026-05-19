# Request

## Original Request

The maintainer clarified that vibit should become a Nakama/Pitaya-class product and cover all common functionality, then asked to replan the development route and continue.

## Clarified Requirement

Ratify a repository-visible roadmap that treats common Nakama/Pitaya capability coverage as a product target, while preserving vibit's agent-native standards and direct API compatibility deferrals.

## User-Visible Outcome

Maintainers and future agents can inspect the product target, capability families, phase order, and next recommended direction from repository artifacts instead of relying on conversation memory.

## Non-Goals

- Do not implement protocol logout behavior in this change.
- Do not implement concrete WebSocket close handoff.
- Do not implement reconnect, session carrier, presence, chat, groups, matchmaking, match runtime, admin console, SDK, cluster, or RPC behavior.
- Do not add dependencies.
- Do not add direct Nakama/Pitaya API compatibility.

## Unknowns

- Which client SDK should be first.
- Whether admin console starts as CLI, HTTP API, or web UI.
- Whether any future direct compatibility surface will be intentionally adopted.

## Acceptance Criteria

- [x] The roadmap standard and Simplified Chinese translation exist.
- [x] ADR-0078 records the product parity decision.
- [x] `.arch/reference.yaml` records parity capability families and phase order.
- [x] The work queue leaves a concrete next ready work item for protocol logout route gate.
- [x] Repository checks cover the roadmap markers.
