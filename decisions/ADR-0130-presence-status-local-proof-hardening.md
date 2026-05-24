# ADR-0130: Presence Status Local Proof Hardening

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-harden-presence-status-local-proof/`

Related conversations:

- `conversations/2026-05-24-presence-status-local-proof-hardening.md`

Related artifacts:

- `runtime/internal/app/presence/presence_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
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

`ADR-0129` piloted the agent-native workflow on a Nakama-style self-presence/status request and opened `M-150/W-0222` as the bounded follow-up. The repository already had a presence query, connection registry, close policy, and authenticated local alpha request flow, but the local proof did not explicitly cover offline behavior after close or invalidation.

Presence/status is a core Nakama-style online-service capability. Before broader realtime, chat, social, matchmaking, or match runtime work depends on it, the alpha foundation should prove that server-owned connection state makes a player online only while active bound connection records remain active.

Nakama capability family:

```text
presence_status_and_notifications
```

Pitaya remains deferred as a future distributed architecture reference.

## Decision

Harden the self-presence/status local proof using tests only.

The application-level proof adds or strengthens coverage for:

```text
online from active bound connection
offline after transport close
offline after close-policy invalidation
```

The local alpha proof uses the existing authenticated Protobuf request flow and existing presence query route to prove:

```text
bound authenticated player -> presence online
transport closed connection -> presence offline
bound authenticated player -> presence online
close policy invalidation -> presence offline
```

Open:

```text
M-151/W-0223 Strengthen authenticated gameplay failure-path verification
```

as the next-ready work item.

This decision does not add runtime behavior beyond tests, protocol routes, Protobuf source, generated output, migrations, dependencies, startup wiring, persistence, presence subscriptions, status broadcast fanout, chat, groups, matchmaking, match runtime, distributed runtime, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add a new presence route specifically for offline or status checks.
- Add status subscription or broadcast fanout.
- Wire close policy directly into logout or WebSocket close behavior in this slice.
- Treat the existing lower-level registry tests as enough and skip local alpha flow proof.
- Jump to chat, groups, matchmaking, or match runtime now that presence exists.

## Rationale

The selected work item asked for proof hardening, not feature expansion. Existing runtime surfaces already represent the behavior: bound active connections make presence online, and terminal connection states make presence offline. Tests at both the application service and authenticated protocol-flow boundary give useful confidence without widening product scope.

New protocol routes, subscriptions, fanout, persistence, and direct compatibility are broader product decisions. They should remain behind later bounded work items.

## Agent Reasoning Summary

The maintainer asked to continue toward Nakama. The active work item required applying the previous workflow pilot to the selected presence/status proof. The codebase showed that most behavior already existed, so the correct small step was to strengthen tests and durable memory rather than add new runtime behavior.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  alpha_foundation_confidence: high
  testability: high
  bounded_scope: high
  reuse_existing_protocol_surface: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-150/W-0222` is completed.
- `M-151/W-0223` is next-ready.
- Presence/status local alpha proof now covers online, close/offline, and invalidation/offline.
- The proof uses existing runtime and Protobuf surfaces.
- Broad presence subscriptions, status fanout, chat, groups, matchmaking, match runtime, persistence, distributed runtime, SDKs, operations, and direct compatibility remain deferred.
- `runtime.presence_status_local_proof_hardening` checks the hardening records and next-ready state.

## Reversal Conditions

Revisit this decision if:

- the local alpha proof needs a real external client fixture instead of in-process Protobuf frame tests;
- future close handoff or WebSocket transport behavior changes the right proof boundary;
- presence/status expands into subscription or broadcast semantics and requires a new protocol gate;
- alpha users report a more urgent Nakama-style blocker.

## Follow-Up

- Complete `W-0223`: strengthen authenticated gameplay failure-path verification.
- Keep new protocol routes, Protobuf, generated output, migrations, dependencies, subscriptions, broadcasts, chat, groups, matchmaking, match runtime, persistence, distributed runtime, and direct compatibility behind later bounded work items.
