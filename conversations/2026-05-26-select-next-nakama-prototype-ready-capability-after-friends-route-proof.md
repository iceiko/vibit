# Conversation: Select Next Nakama Prototype-Ready Capability After Friends Route Proof

Date: 2026-05-31

## Context

The maintainer asked the agent to continue and to commit and push, noting that the GitHub token is in a Git-ignored file. `M-170/W-0242` had just proved the protected friends relationship protocol route family locally and opened `M-171/W-0243` as the next selection slice.

## Maintainer Narrative

The maintainer's continuation request means to advance one `next_ready` work item unless blocked by an ask-first boundary, verification failure, or required maintainer confirmation.

## Agent Response Summary

The agent reviewed the Nakama-first roadmap, product maturity milestones, alpha developer flow, local alpha proof path, and friends relationship route proof. The agent selected the next bounded capability family as:

```text
admin_console_metrics_observability_and_operations
```

The selected follow-up is:

```text
M-172/W-0244 Define minimum operations inspection surface gate
```

The selection does not implement an operations inspection surface yet. It records that the current source-first alpha needs a minimum redacted operations inspection posture before expanding into groups, parties, chat, leaderboards, matchmaking, match runtime, SDK publication, hosted surfaces, or distributed runtime.

## Decisions

- Accept `ADR-0151`.
- Register `runtime.next_nakama_prototype_ready_capability_after_friends_route_proof`.
- Complete `M-171/W-0243`.
- Open `M-172/W-0244 Define minimum operations inspection surface gate` as next-ready.
- Keep Pitaya deferred as a future distributed architecture reference.
- Keep direct Nakama/Pitaya API compatibility out of scope.

## Artifacts

- `changes/2026-05-26-select-next-nakama-prototype-ready-capability-after-friends-route-proof/`
- `decisions/ADR-0151-select-next-nakama-prototype-ready-capability-after-friends-route-proof.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Should the future operations inspection gate start with repository-check introspection, local HTTP endpoints, CLI inspection commands, or a source-first static runbook posture?
- Which state categories should be inspectable first: server health, configuration posture, routes, players, sessions, tokens, connections, storage objects, friends relationships, realtime state, or migrations?
- Which identifiers and metadata are safe to show by default, and which must be redacted unless a later explicit debug gate authorizes exposure?

## Follow-Up

- Complete `W-0244`: define the minimum operations inspection surface gate.

## Redaction Notes

No raw access-token, device credential, verifier digest, lookup digest, verifier key, DSN, local secret, or GitHub token value is recorded here.
