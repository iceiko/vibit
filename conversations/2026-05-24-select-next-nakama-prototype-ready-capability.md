# Conversation: Select Next Nakama Prototype-Ready Capability

Date: 2026-05-24

## Context

The maintainer asked the agent to continue toward the Nakama target for many steps. `M-151/W-0223` had just strengthened authenticated gameplay failure-path verification and opened `M-152/W-0224` as the next selection slice.

## Maintainer Narrative

The maintainer clarified that Nakama should be the primary product reference and that the product purpose is AI-native development plus AI-native testing: user requirements should become specs, acceptance criteria, tests, implementation, verification, and durable project memory through AI assistance.

## Agent Response Summary

The agent reviewed the Nakama-first roadmap, product maturity milestones, alpha developer flow, examples directory, local request-loop proof, and authenticated failure-path proof. The agent selected the next bounded capability family as:

```text
client_sdks_examples_and_developer_experience
```

The selected follow-up is:

```text
M-153/W-0225 Define local alpha example client path gate
```

The selection does not implement an example client yet. It records that the current source-first alpha needs a clearer client-like or example-app path before expanding into chat, friends, groups, matchmaking, match runtime, operations/admin breadth, SDK publication, hosted surfaces, or distributed runtime.

## Decisions

- Accept `ADR-0132`.
- Register `runtime.next_nakama_prototype_ready_capability_selection`.
- Complete `M-152/W-0224`.
- Open `M-153/W-0225 Define local alpha example client path gate` as next-ready.
- Keep Pitaya deferred as a future distributed architecture reference.
- Keep direct Nakama/Pitaya API compatibility out of scope.

## Artifacts

- `changes/2026-05-24-select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof/`
- `decisions/ADR-0132-select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Should the future example path be a small Go client, a command script over a client helper, or a docs-plus-test fixture package?
- Should the first implementation demonstrate only the authenticated request loop, or also the existing realtime outbound frame observation?

## Follow-Up

- Complete `W-0225`: define the local alpha example client path gate.

## Redaction Notes

No raw access-token, device credential, verifier digest, lookup digest, verifier key, DSN, local secret, or GitHub token value is recorded here.
