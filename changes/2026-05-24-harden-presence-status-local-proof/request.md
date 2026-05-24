# Request

## Original Request

```text
继续推进，目标nakama 推进10步以上，推进10个小时以上，不要停止，我会离开10小时。
```

## Clarified Requirement

Complete `M-150/W-0222 Harden presence status local proof through close and offline cases`.

The previous workflow pilot selected a Nakama-style self-presence/status requirement. This slice must harden the existing local alpha proof so that a developer can see online presence after authenticated connection binding and offline presence after transport close or policy invalidation.

## Selected User Requirement

```text
As a developer using the local alpha request flow, I want the server to prove a player's self-presence/status through authenticated requests, including online after binding and offline after close or invalidation, so I can trust the foundation before broader Nakama-style realtime, social, and multiplayer features build on it.
```

## Nakama Capability Mapping

Capability family:

```text
presence_status_and_notifications
```

Nakama-style product value:

- Presence/status is a core online-service capability for game backends.
- The proof should be trustworthy before chat, groups, parties, matchmaking, match runtime, or broader realtime features depend on it.
- This slice hardens proof using existing vibit-native protocol/runtime surfaces, not direct Nakama API compatibility.

## User-Visible Outcome

Developers and future agents now have focused Go tests proving:

- self-presence reports online from an active bound connection;
- self-presence reports offline after a transport close;
- self-presence reports offline after close-policy invalidation;
- the authenticated local alpha Protobuf request flow can observe those states without new routes or wire messages.

## Non-Goals

- Add new presence protocol routes or Protobuf messages.
- Add generated output, migrations, dependencies, persistence, or startup wiring.
- Add presence subscriptions, status broadcast fanout, chat, groups, parties, matchmaking, match runtime, leaderboards, economy, SDKs, operations, or distributed runtime.
- Add Pitaya-style cluster/RPC/frontend-backend/service-discovery work.
- Add hosted deployments, release artifacts, public announcements, or paid promotion.
- Add direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- `runtime/internal/app/presence/presence_test.go` proves online self-presence and offline self-presence after close and invalidation.
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go` proves the same selected behavior through the existing authenticated local alpha Protobuf request flow.
- No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, persistence, startup wiring, or direct compatibility scope is added.
- `ADR-0130` records the hardening decision.
- Repository checks include `runtime.presence_status_local_proof_hardening`.
- `M-151/W-0223 Strengthen authenticated gameplay failure-path verification` is opened as the next-ready bounded follow-up.
