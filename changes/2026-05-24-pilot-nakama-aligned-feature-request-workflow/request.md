# Request

## Original Request

```text
继续推进，目标nakama 推进10步以上，推进10个小时以上，不要停止，我会离开10小时。
```

## Clarified Requirement

Complete `M-149/W-0221 Pilot Nakama-aligned feature request workflow`.

Apply the agent-native feature request and test workflow to one concrete Nakama-aligned product capability request. Select a bounded follow-up slice that can be implemented and verified without broad runtime expansion.

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

- Presence/status is a core online-service capability in a game backend.
- The current alpha already has connection lifecycle, binding, protected self-presence query, and E2E proof coverage.
- The next useful step is to harden the proof around close and offline behavior before adding subscriptions, broadcasts, chat, friends, groups, matchmaking, match runtime, persistence, or direct compatibility.

## User-Visible Outcome

Future agents can continue with a bounded `W-0222` implementation slice:

```text
Harden presence status local proof through close and offline cases.
```

The implementation slice should use existing presence, connection, protocol, and local alpha request-flow surfaces where possible.

## Non-Goals

- Add runtime behavior in this pilot slice.
- Add new protocol routes or Protobuf messages in this pilot slice.
- Add generated output, migrations, dependencies, or startup wiring.
- Add presence subscriptions, status broadcast fanout, chat, groups, parties, matchmaking, match runtime, leaderboards, economy, operations, SDKs, or distributed runtime.
- Add Pitaya-style cluster/RPC/frontend-backend/service-discovery work.
- Add hosted deployments, release artifacts, public announcements, or paid promotion.
- Add direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- The pilot records the user requirement, Nakama capability family, acceptance criteria, test plan, implementation boundaries, verification plan, and durable memory expectations.
- The selected capability family is `presence_status_and_notifications`.
- The selected follow-up direction is `presence_status_local_proof_hardening`.
- `M-150/W-0222 Harden presence status local proof through close and offline cases` is opened as next-ready.
- `ADR-0129` records the pilot decision.
- Repository checks include `runtime.nakama_aligned_feature_request_workflow_pilot`.
- No runtime/protocol/generated/migration/dependency/direct compatibility scope is added in this pilot slice.
