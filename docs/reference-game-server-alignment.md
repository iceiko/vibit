# Reference Game Server Alignment

Status: Draft v0.1
Last updated: 2026-05-13
Scope: Active reference baseline for vibit's game server capability planning

This document records how vibit should use mature game server projects as references.

The paired Simplified Chinese translation is `docs/reference-game-server-alignment.zh-CN.md`. The English file is authoritative.

`ADR-0078` and `docs/nakama-pitaya-product-parity-roadmap.md` refine this reference baseline into an explicit product target: vibit should become a Nakama/Pitaya-class game backend framework by covering common capability families while preserving vibit's agent-native constraints. This document still defines reference roles; the product parity roadmap defines phased execution.

## 1. Purpose

vibit is not only an inventory proof slice.

vibit should eventually cover the same broad problem class as mature game backend and game server frameworks such as Nakama and Pitaya:

- Player identity and sessions.
- Social and realtime communication.
- Storage and durable game state.
- Matchmaking and match/session lifecycle.
- Realtime multiplayer and authoritative server behavior.
- Leaderboards, rewards, currencies, and other metagame systems.
- Operational visibility and production maintainability.
- Scalable routing and, later, distributed server topology.

The difference is not the existence of those game server capabilities. The difference is that vibit must make those capabilities Agent-Native:

- Explicit module ownership.
- Contract-first public behavior.
- Generated repeatable structure.
- Machine-readable manifests.
- Narrow runtime boundaries.
- Durable decision records.
- Verification commands that agents can run before and after changes.

Nakama documents a broad game server product surface that includes user accounts, authentication, storage, friends, groups, chat, leaderboards, tournaments, matchmaking, realtime multiplayer, and authoritative match runtime concepts. Pitaya documents a Go game server framework shape around client acceptors, sessions, route handlers, remote calls, groups, serializers, and cluster vocabulary. vibit should learn from both surfaces while making every adopted capability explicit, generated where appropriate, and machine-checkable.

Current product target markers:

```text
nakama_pitaya_product_parity_roadmap: ratified
decision: ADR-0078
check_rule: runtime.reference_product_parity_roadmap
parity_goal: nakama_pitaya_same_class_common_capability_coverage
api_compatibility_goal: false
recommended_next_direction: define_protocol_logout_route_gate
```

## 2. Reference Roles

### Nakama

Use Nakama as the primary reference for a broad game backend product capability surface.

Reference areas:

- Authentication, users, accounts, and sessions.
- Social systems such as friends, groups, chat, parties, leaderboards, and tournaments.
- Storage objects and server-side runtime customization.
- Match listing, matchmaker, realtime multiplayer, and authoritative match logic.
- Dashboard, metrics, and operational visibility.
- Economy, rewards, currencies, and LiveOps-style feature families when looking at the broader Heroic Labs ecosystem.

Nakama should guide what a useful general game backend must eventually support.

Nakama must not become a governing API shape for vibit unless a future ADR explicitly adopts a compatibility surface.

### Pitaya

Use Pitaya as the primary reference for Go game server framework architecture vocabulary.

Reference areas:

- Client connection acceptors such as WebSocket and TCP.
- User sessions and session binding.
- Handler routing for client messages.
- Remote calls for server-to-server communication.
- Frontend and backend server roles in cluster mode.
- Groups for broadcast and multicast use cases such as rooms.
- Message forwarding, serializers, RPC services, and cluster service discovery.

Pitaya should guide how Go game server frameworks separate transport, session, route, server role, RPC, and group concerns.

Pitaya must not force vibit into distributed clustering before the modular monolith proof slice is healthy.

## 3. Capability Matrix

The following matrix is a planning tool. It is not a promise that every capability must be implemented immediately.

| Capability | Reference | Vibit Direction |
| --- | --- | --- |
| Accounts and authentication | Nakama | `player` or `identity` module with explicit auth/session contracts. |
| Sessions and connection identity | Nakama, Pitaya | Platform session adapter plus app-owned session context; no transport shortcut identity. |
| Storage objects | Nakama | PostgreSQL-backed module state first; object storage only for large artifacts. |
| Inventory | Common game backend need | First proof module; must prove contract -> generated shape -> handler -> tests. |
| Currency and wallets | Nakama/Hiro capability family | Future `currency` module with transactional invariants. |
| Rewards and claims | Nakama/Hiro capability family | Future `reward` module with eligibility, idempotency, and event tests. |
| Friends, groups, parties | Nakama | Future social modules with explicit membership ownership and events. |
| Chat and realtime messaging | Nakama | Future realtime module; not hidden in WebSocket transport. |
| Presence and status | Nakama, Pitaya | Future platform/application capability with explicit lifecycle semantics. |
| Matchmaking | Nakama | Future `matchmaking` module; query and criteria contracts before implementation. |
| Match/session lifecycle | Nakama, Pitaya | Future `match` module; authoritative match behavior separate from transport. |
| Rooms and broadcast groups | Pitaya | Future group/room abstraction; target scopes already reserve `room` and `match`. |
| Leaderboards and tournaments | Nakama | Future competitive modules; ranking and reset rules must be contract-first. |
| Authoritative realtime simulation | Nakama | Future match runtime; server remains authoritative. |
| Cluster frontend/backend split | Pitaya | Deferred until single-process boundaries are proven. |
| Server-to-server RPC | Pitaya | Deferred; must not bypass module contracts when introduced. |
| Dashboard and operations | Nakama | Future admin/inspection surface; do not mix with gameplay protocol. |
| Metrics and observability | Nakama, Pitaya | Deferred dependency decision; must be platform-owned. |

## 4. Phased Roadmap

### Phase 0: Agent-Native Foundation

Current phase.

Goals:

- Constitution, AGENTS guides, change specs, conversation logs, ADRs.
- Architecture manifests and checks.
- Go/WebSocket/Protobuf/PostgreSQL direction.
- Runtime handoff, protocol adapter, and application dispatch skeleton.
- First inventory semantic and wire contracts.

Exit criteria:

- A small backend change can move through spec, contract, generated shape, handwritten logic, tests, verification, and docs.

### Phase 1: First Vertical Game Backend Slice

Goals:

- Inventory command/query/event handler boundary.
- Repository and policy interfaces.
- PostgreSQL persistence for the first durable module state.
- Focused tests for invariants, errors, events, and repository behavior.
- Protocol response mapping through the existing dispatcher.

Reference pull:

- Nakama storage and custom server logic concepts.
- Pitaya handler routing separation.

### Phase 2: Player, Session, And Transport

Goals:

- Player/account/session module boundaries.
- WebSocket transport adapter.
- Session validation and connection lifecycle.
- Protocol error response mapping.
- Minimal playable client/server request flow.

Reference pull:

- Nakama authentication and session model.
- Pitaya session binding and WebSocket acceptor separation.

### Phase 3: Core Game Backend Modules

Goals:

- Currency/wallet.
- Rewards/claims.
- Friends/groups/party.
- Presence/status.
- Chat or realtime messaging.
- Leaderboards.

Reference pull:

- Nakama social, competitive, and metagame capabilities.

### Phase 4: Multiplayer And Match Runtime

Goals:

- Matchmaking criteria and ticket lifecycle.
- Match/session lifecycle.
- Room/group broadcast semantics.
- Authoritative match loop contracts.
- Reconnect/replay decisions.

Reference pull:

- Nakama matchmaker and authoritative/relayed multiplayer distinction.
- Pitaya groups and route/handler model.

### Phase 5: Distributed Runtime

Goals:

- Frontend/backend server role split.
- Server-to-server RPC.
- Service discovery.
- Distributed groups/rooms.
- Cluster-safe session and routing semantics.

Reference pull:

- Pitaya cluster architecture.

This phase must not begin until single-process module, transaction, protocol, and verification boundaries are stable.

## 5. Agent Rules

Agents must:

- Consult this document before proposing new game server modules or runtime subsystems.
- Check whether a proposed capability maps to a known Nakama/Pitaya capability family.
- Preserve vibit's Agent-Native constraints even when matching reference functionality.
- Record whether a reference pattern is adopted, adapted, deferred, or rejected, and why.
- Prefer adding a small enforceable manifest/check over adding broad aspirational text.

Agents must not:

- Copy external APIs without an explicit compatibility ADR.
- Add a Nakama-like or Pitaya-like feature directly in transport handlers.
- Add cluster/RPC/service-discovery work before the modular monolith proof slice is stable.
- Treat feature parity as more important than explicit contracts, module ownership, or verification.
- Use references as an excuse to bypass vibit's generated-file and boundary rules.

## 6. Next Planning Implications

The next runtime work should still avoid a premature WebSocket server if handler, repository, and policy boundaries are unclear.

Recommended near-term sequence:

1. Define inventory runtime repository and policy interfaces.
2. Implement the first inventory command/query handler through the application dispatcher.
3. Add PostgreSQL migration and repository behavior after the interfaces are stable.
4. Add WebSocket transport once the request -> dispatch -> handler -> result path is tested.
5. Add player/session/auth boundaries before treating WebSocket connection identity as durable player identity.

This sequence keeps the project aligned with Nakama/Pitaya functionality while preserving vibit's primary difference.

## 7. References

- Nakama documentation home: `https://heroiclabs.com/docs/`
- Nakama getting started: `https://heroiclabs.com/docs/nakama/getting-started/`
- Nakama multiplayer concepts: `https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Nakama GitHub repository: `https://github.com/heroiclabs/nakama`
- Pitaya overview: `https://pitaya.readthedocs.io/en/latest/overview.html`
- Pitaya features: `https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya GitHub repository: `https://github.com/topfreegames/pitaya`
