# ADR-0087: Presence Protocol Query Functional Slice

Status: Accepted
Date: 2026-05-20
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-20-define-presence-protocol-query-functional-slice/`

Related conversations:

- `conversations/2026-05-20-presence-protocol-query-functional-slice.md`

Related artifacts:

- `proto/vibit/presence/v1/presence.proto`
- `runtime/internal/generated/proto/vibit/presence/v1/presence.pb.go`
- `runtime/internal/app/presence/presence.go`
- `runtime/internal/app/presence/presence_test.go`
- `runtime/internal/platform/protocol/protobuf/presence_bridge.go`
- `runtime/internal/platform/protocol/protobuf/presence_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `rules/check-rules.json`

## Context

`ADR-0085` added the first server-owned presence lifecycle primitive: the application-owned active connection registry can derive player online/offline presence from active bound connection records, and PostgreSQL startup composition wires WebSocket open/close lifecycle plus successful first-message binding into that registry.

The short-term `v0.1 alpha` target in `ADR-0086` calls out a basic observable presence query as part of the first developer-usable runtime path. Before this decision, the server had real presence state but no protected client-visible query route for it.

`ADR-0082` classifies this as a Tier 2 functional slice: the boundary belongs in the change spec, and implementation can proceed in the same work item without a separate pure confirmation gate.

## Decision

Select:

```text
define_presence_protocol_query_functional_slice
```

as a Tier 2 functional slice and implement the smallest protected presence query directly.

The implementation:

- adds `runtime.presence.GetPlayerPresence` as a query route,
- adds `vibit.presence.v1.GetPlayerPresenceRequest` and `GetPlayerPresenceResponse`,
- keeps the route protected by the existing `AuthenticatedRequest` request-level access-token wrapper,
- makes the first query self-only: a validated player may query only their own presence snapshot,
- reads the existing application-owned registry-backed presence snapshot,
- returns online/offline status, active connection count, bounded active connection metadata, runtime session ids, `last_seen_at`, and `observed_at`,
- does not expose access-token record ids through the protocol response,
- registers the query in the PostgreSQL runtime composition where the presence lifecycle registry exists.

## Boundaries

This ADR keeps these boundaries:

- Presence query behavior is application-owned under `runtime/internal/app/presence`.
- Presence snapshot state remains application-owned under `runtime/internal/app/connection`.
- Protobuf protocol adapter code owns only wire/domain payload mapping.
- WebSocket transport remains credential-neutral and does not own presence semantics.
- Authentication service behavior remains token lifecycle behavior and does not own presence queries.
- The first query uses request-token proof, not client-supplied `player_id` or `session_id` metadata.

This decision does not add presence subscriptions, presence broadcasts, chat, friends, groups, parties, matchmaking, match runtime, cluster behavior, SDK behavior, operations/admin behavior, reconnect tokens, resume tokens, logout-triggered socket close, runtime session revocation, durable/distributed presence, dependencies, broad product modules, or direct Nakama/Pitaya API compatibility.

## Nakama And Pitaya Mapping

Nakama informs the product pressure: presence/status is a basic game backend capability that should be visible to authenticated clients before higher-level realtime/social features are built.

Pitaya informs the layering pressure: acceptors, session/binding context, route handlers, and connection management should remain separate.

vibit adapts those lessons by exposing one protected self-presence query over its existing registry-backed state. It does not copy either project's public API.

## Alternatives Considered

- Add presence subscriptions or broadcasts immediately.
- Add durable or distributed presence storage.
- Let clients query arbitrary players' presence.
- Trust envelope `session.player_id` as query authority.
- Put presence in the authentication service or WebSocket transport.
- Expose access-token record ids through the protocol response.
- Delay protocol presence until chat or social modules exist.

## Rationale

The smallest useful protocol-visible presence feature is a protected self-query. It proves the full route path from authenticated request wrapper to application handler to Protobuf response without opening broader social, subscription, or distributed runtime scope.

Self-only semantics are deliberately conservative. They let the alpha runtime expose real server-owned presence while avoiding friend graph, privacy policy, admin visibility, and cross-player lookup decisions.

Keeping access-token record ids out of the wire response preserves token lifecycle internals as server metadata. Runtime session ids are already protocol-visible after `ADR-0084`; including them in the presence response is consistent with the current session carrier posture while still not making them proof.

## Agent Reasoning Summary

The maintainer asked to continue moving toward a developer-usable Nakama/Pitaya-class baseline. After presence lifecycle state existed, the next bounded product-parity step was to expose that state through one protected query rather than jumping to subscriptions, chat, social modules, matchmaking, or match runtime.

## Decision Weights

```yaml
decision_weights:
  development_velocity: high
  alpha_user_observability: high
  route_protection_safety: high
  transport_application_separation: high
  durable_distributed_scope: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `runtime.presence_protocol_query_functional_slice` becomes the repository check rule for this slice.
- The runtime has a protected client-visible presence query over server-owned presence state.
- The query proves a useful authenticated gameplay route outside inventory/authentication.
- Future alpha E2E work can include onboarding, login, connection binding, protected inventory, presence query, and logout.
- Cross-player presence, subscriptions, broadcasts, durable/distributed presence, social modules, matchmaking, and match runtime remain separate future work.

## Reversal Conditions

Revisit this decision if future privacy policy selects cross-player presence visibility, if distributed runtime replaces single-process registry semantics, if presence subscriptions need a different source of truth, or if direct Nakama/Pitaya API compatibility is explicitly selected and requires different route or payload semantics.

## Follow-Up

- Select the next alpha-enabling work item after `W-0178`.
- Keep first local onboarding/device credential issuance, authenticated gameplay E2E, runtime runbook refresh, example client/request loop, health/readiness/version/config surface, and alpha acceptance checks in the near-term queue.
- Keep subscriptions, broadcasts, chat, social modules, matchmaking, match runtime, durable/distributed presence, reconnect/resume tokens, logout-triggered close, runtime session revocation, and direct compatibility behind explicit future work.
