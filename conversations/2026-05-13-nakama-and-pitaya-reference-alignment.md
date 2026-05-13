# Conversation: Nakama And Pitaya Reference Alignment

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-align-with-nakama-and-pitaya/`

Related artifacts:

- `docs/reference-game-server-alignment.md`
- `.arch/reference.yaml`
- `decisions/ADR-0019-nakama-and-pitaya-reference-baseline.md`

## Context

The project had just added the first Go runtime handoff and application dispatch skeletons. The maintainer then clarified that vibit should actively and repeatedly reference Nakama and Pitaya because the broad server capability surface is essentially the same; vibit's primary distinction is being Agent-Native.

## Maintainer Narrative

The maintainer said vibit should actively and frequently refer to Nakama and Pitaya. The intended functionality is in the same game server problem class, while the key difference is that vibit is Agent-Native.

The maintainer also asked for a review and a plan for future direction based on this alignment.

## Agent Response Summary

The agent treated this as a long-term product and architecture signal rather than a request to immediately copy APIs or implement a new feature.

The agent added a reference alignment standard that uses Nakama and Pitaya as active capability baselines, while preserving vibit's own constitutional constraints: contract-first behavior, explicit module ownership, generated shape, architecture checks, and agent-operable change workflows.

## Decisions

- Treat Nakama as the strongest reference for broad all-in-one game backend product capabilities.
- Treat Pitaya as the strongest reference for Go server framework, routing, session, cluster, and RPC architecture vocabulary.
- Use both references as planning baselines, not governing standards.
- Keep vibit's differentiator as Agent-Native maintainability and verification.
- Avoid distributed clustering until the modular monolith proof slice has proven its boundaries.

## Artifacts

- Added `docs/reference-game-server-alignment.md` and `docs/reference-game-server-alignment.zh-CN.md`.
- Added `.arch/reference.yaml`.
- Added `decisions/ADR-0019-nakama-and-pitaya-reference-baseline.md`.
- Updated README, AGENTS, and `.arch/README` intake paths.

## Open Questions

- Which future capability should become the first non-inventory module: player account/session, wallet/currency, match, presence, or chat?
- Should generated route registration be built before expanding into account/session modules?
- How much Nakama-style API ergonomics should vibit expose versus purely manifest-generated APIs?
- When should Pitaya-style distributed routing become a first-class runtime phase?

## Follow-Up

- Use the reference capability matrix before adding new modules or runtime subsystems.
- Keep the next runtime step focused on handler/repository/policy boundaries before WebSocket transport.
- Revisit the capability roadmap after the first inventory command/query/event slice is fully implemented.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
