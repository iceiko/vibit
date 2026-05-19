# ADR-0078: Nakama And Pitaya Product Parity Roadmap

Status: Accepted
Date: 2026-05-18
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-18-confirm-next-direction-after-websocket-close-policy-implementation/`
- `changes/2026-05-18-ratify-nakama-pitaya-product-parity-roadmap/`

Related conversations:

- `conversations/2026-05-18-nakama-pitaya-product-parity-roadmap.md`

Related artifacts:

- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The project already used Nakama and Pitaya as reference baselines, but the previous posture could be interpreted as design inspiration only. The maintainer clarified that vibit's intended target is stronger: vibit should become a Nakama/Pitaya-class product and cover their common feature surface.

At the same time, vibit must preserve its original thesis. The project should not become a direct clone, a compatibility wrapper, or a transport-heavy server that hides domain behavior in handlers. The product target must be expressed in a way that future agents can inspect, verify, and use for prioritization.

The work queue was blocked at `M-095/W-0167` after the WebSocket close policy implementation. One candidate direction was `expand_core_game_backend_modules_after_nakama_pitaya_review`. The clarified product goal makes that direction the right next step, but it must be made precise before feature expansion begins.

## Decision

Select:

```text
expand_core_game_backend_modules_after_nakama_pitaya_review
```

Implement it as a roadmap ratification slice:

```text
ratify_nakama_pitaya_product_parity_roadmap
```

Add:

```text
docs/nakama-pitaya-product-parity-roadmap.md
docs/nakama-pitaya-product-parity-roadmap.zh-CN.md
```

Define the repository check rule:

```text
runtime.reference_product_parity_roadmap
```

The ratified target is:

```text
parity_goal: nakama_pitaya_same_class_common_capability_coverage
api_compatibility_goal: false
direct_nakama_pitaya_api_compatibility_added: false
implementation_authorized_by_this_standard: roadmap_only
```

The roadmap defines first-class capability families for identity/auth/session, connection lifecycle, storage, presence/status/notifications, chat/realtime messaging, friends/groups/parties, leaderboards/tournaments, economy/rewards/progression, matchmaking, match runtime, server runtime hooks/RPC, operations, client SDK/developer experience, and distributed runtime.

The near-term execution priority remains lifecycle closure before high-level module expansion:

1. `define_protocol_logout_route_gate`
2. `define_transport_close_handoff_gate`
3. `define_reconnect_connection_epoch_gate`
4. `define_protocol_session_carrier_gate`
5. `define_presence_lifecycle_gate`
6. `strengthen_operations_observability_and_admin_tooling`

## Alternatives Considered

- Continue directly to `define_protocol_logout_route_gate` without changing the product roadmap.
- Continue directly to presence, chat, groups, parties, matchmaking, or match runtime.
- Add a broad aspirational paragraph to `README.md` without a checkable standard.
- Attempt direct Nakama API compatibility.
- Attempt direct Pitaya cluster/RPC behavior before single-process lifecycle closure.
- Install an external agent workflow tool and rely on it instead of repository-native planning.

## Rationale

The clarified target affects prioritization across the whole project. Without a checkable roadmap, future work could keep improving local lifecycle gates while losing sight of the broader game backend product surface.

The roadmap is intentionally standard-only. It changes direction and work method, but it does not authorize new runtime behavior. That keeps the repository honest: product parity becomes a durable target, while concrete capabilities still need their own contracts, gates, implementations, tests, and checks.

The near-term lifecycle-first sequence follows both references. Nakama makes account/session/socket lifecycle a first-class product concern, and Pitaya separates acceptors, sessions, handlers, groups, RPC, and cluster roles. Higher-level features such as presence, chat, parties, matchmaking, and match runtime should not depend on ambiguous logout, close, reconnect, or session-carrier semantics.

## Agent Reasoning Summary

The maintainer's clarification changes the roadmap but not vibit's engineering constitution. The correct response is not to rush into feature modules. The correct response is to make the product parity target inspectable, then resume implementation through the lowest unstable shared lifecycle surfaces.

## Decision Weights

```yaml
decision_weights:
  product_goal_clarity: high
  nakama_pitaya_alignment: high
  agent_context_reduction: high
  lifecycle_before_feature_expansion: high
  direct_api_compatibility: low
  immediate_runtime_behavior: low
  distributed_runtime_now: low
confidence: high
```

## Consequences

- Product parity with common Nakama/Pitaya capability families becomes an explicit roadmap target.
- Future major work must map to a roadmap family.
- New high-level feature modules should wait until the runtime lifecycle foundation is less ambiguous.
- `runtime.reference_product_parity_roadmap` becomes the repository check rule for this ratification.
- The next concrete recommended direction is `define_protocol_logout_route_gate`.

## Reversal Conditions

Revisit this decision if the maintainer narrows vibit back to a smaller proof-of-concept framework, if a direct compatibility goal is explicitly selected, if a new reference baseline replaces Nakama/Pitaya, or if the project chooses to become backend-general instead of game-backend-first.

## Follow-Up

- Define protocol logout route behavior.
- Define concrete transport close handoff.
- Define reconnect and connection epoch behavior.
- Define protocol session carrier behavior.
- Define presence lifecycle after the connection/session/logout foundation is stable.
- Keep the parity matrix updated as new modules are added.
