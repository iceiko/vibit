# Conversation: Presence Protocol Query Functional Slice

Date: 2026-05-20
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0087-presence-protocol-query-functional-slice.md`
Related changes:

- `changes/2026-05-20-define-presence-protocol-query-functional-slice/`

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
- `tools/vibit`

## Context

The repository had completed `M-105/W-0177`, the presence lifecycle functional slice. The runtime could derive server-owned player online/offline presence from active bound connection registry records, but that state was not yet visible through a protected protocol query.

The short-term `v0.1 alpha` goal also identified a basic presence query as part of the first developer-usable runtime path.

## Maintainer Narrative

The maintainer asked to continue:

```text
继续推进
```

Under `docs/workflow.md`, that means advancing exactly one `next_ready` work item unless blocked. `.arch/work-items.yaml` identified `W-0178 define_presence_protocol_query_functional_slice` as next ready.

## Agent Response Summary

The agent advanced W-0178 as a Tier 2 functional slice. The implementation did not add another pure confirmation gate.

The slice added the first protected presence query:

- `runtime.presence.GetPlayerPresence` is a query route.
- `vibit.presence.v1.GetPlayerPresenceRequest` and `GetPlayerPresenceResponse` define the wire payloads.
- Existing `AuthenticatedRequest` request-token proof protects the route.
- The first query is self-only and rejects cross-player presence requests.
- The application handler reads the existing server-owned active connection registry snapshot.
- The response returns status, connection count, bounded active connection metadata, runtime session ids, and timestamps.
- Access-token record ids remain server-only metadata and are not returned in the protocol response.

## Decisions

- Complete `M-106/W-0178`.
- Accept `ADR-0087`.
- Add `runtime.presence_protocol_query_functional_slice` as the repository check rule.
- Keep presence subscriptions, broadcasts, chat, friends, groups, parties, matchmaking, match runtime, cluster, SDK, operations/admin behavior, reconnect/resume tokens, logout-triggered socket close, runtime session revocation, durable/distributed presence, dependencies, direct Nakama/Pitaya API compatibility, and broad product behavior deferred.

## Artifacts

- `proto/vibit/presence/v1/presence.proto`
- `runtime/internal/generated/proto/vibit/presence/v1/presence.pb.go`
- `runtime/internal/app/presence/presence.go`
- `runtime/internal/app/presence/presence_test.go`
- `runtime/internal/platform/protocol/protobuf/presence_bridge.go`
- `runtime/internal/platform/protocol/protobuf/presence_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `changes/2026-05-20-define-presence-protocol-query-functional-slice/`
- `decisions/ADR-0087-presence-protocol-query-functional-slice.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- The next alpha-enabling work item is not yet selected in this slice.
- First local onboarding/device credential issuance remains a near-term alpha gap.
- A full authenticated gameplay E2E path remains future work.
- Presence subscriptions and broadcasts remain deferred.
- Durable/distributed presence remains deferred.
- Direct Nakama/Pitaya API compatibility remains deferred.

## Follow-Up

Select the next bounded alpha-enabling work item after W-0178, likely from the `docs/v0.1-alpha-goal.md` preferred queue unless the maintainer redirects.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, close reason text, remote addresses, headers, or GitHub tokens are recorded in this conversation log.
