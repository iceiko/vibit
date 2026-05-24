# Feature Request

## Original Request

Advance toward a Nakama-class backend by piloting the new agent-native feature request scaffold on one bounded Nakama-style requirement. Use the pilot to turn the requirement into a spec, acceptance criteria, test plan, implementation boundaries, verification expectations, durable memory, and a bounded next follow-up without adding runtime behavior yet.

## Clarified Requirement

Use the new `tools/vibit scaffold feature` workflow on one concrete Nakama-style product request, select a bounded follow-up, and prove that the repository can move from a user requirement to spec, acceptance criteria, test plan, implementation boundary, verification record, and durable memory before coding.

The selected concrete request is:

```text
As a game developer building on vibit, I want a future player friendship relationship lifecycle so players can request, accept, reject, remove, and block social relationships through server-authoritative rules before broader groups, parties, chat, matchmaking, or match runtime features depend on social graph state.
```

## User-Visible Outcome

No runtime behavior changes in this pilot. The user-visible outcome is a completed source-first intake that future agents and contributors can inspect before implementing any friend system code. The next-ready work item becomes a gate for the friendship relationship lifecycle, not an implementation slice.

## Nakama Capability Mapping

- Capability family: `friends_groups_and_parties`
- Product intent: Nakama-class game backends commonly expose friend and social graph primitives before richer social surfaces such as groups, parties, chat targeting, invites, and matchmaking filters. vibit should define that capability through its own contract-first and agent-native workflow before adding behavior.
- API compatibility: This scaffold does not authorize direct Nakama API compatibility.

## Pitaya Status

Pitaya remains a deferred future distributed architecture reference. This request does not introduce frontend/backend roles, RPC, service discovery, groups, cluster routing, distributed session behavior, or Pitaya-style distributed architecture.

## Non-Goals

- Runtime friendship behavior.
- Friend protocol routes or Protobuf messages.
- Generated output.
- Friend repository interfaces, PostgreSQL adapters, migrations, indexes, or persistence.
- Startup wiring.
- Authentication/session behavior changes.
- Groups, parties, chat, matchmaking, match runtime, leaderboards, tournaments, economy, operations/admin behavior, SDK publication, hosted deployments, release artifacts, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Unknowns

- The eventual friend identity model should reuse validated player account identity, but exact request/response contract shapes remain for the W-0231 gate.
- The eventual persistence model, uniqueness constraints, block semantics, and event vocabulary remain for the W-0231 gate.
- The eventual client protocol route family remains deferred until semantic contracts are stable.

## Acceptance Criteria

- [x] The scaffolded change directory exists and contains `request.md`, `spec.yaml`, `impact.md`, `plan.md`, `checklist.md`, and `verification.md`.
- [x] The request records a concrete Nakama-style capability family: `friends_groups_and_parties`.
- [x] The request records Pitaya as deferred and rejects distributed topology work in this pilot.
- [x] Acceptance criteria, test expectations, implementation boundaries, verification commands, and durable memory expectations are written before any runtime implementation.
- [x] The selected follow-up is a bounded gate: `W-0231 Define friends relationship lifecycle gate`.
- [x] The pilot adds no runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, hosted deployment, SDK, distributed runtime, or direct compatibility scope.

## Test Expectations

- [x] Positive behavior tests are planned for the future gate: request, accept, reject, remove, block, idempotency, and list/read behavior.
- [x] Negative behavior tests are planned for the future gate: self-friendship, duplicate requests, blocked-player interactions, invalid state transitions, missing identity, invalid target, and redacted errors.
- [x] Permission/authentication tests are planned for the future gate: validated player identity must drive actor ownership and metadata-only identity must not be accepted.
- [x] Persistence/protocol/integration tests are deferred to later implementation gates and recorded as not applicable to this pilot.

## Redaction Notes

Do not record raw device credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private user data beyond explicit request text.
