# Request

## Original Request

```text
继续推进,另外我们要主动经常的参考nakama和pitaya 我们实际上要实现的功能和它是一样的，我们主要是区别是我们是 Agent Native。你也可以针对这个进行一次回顾，并对我们的计划和未来的推进方向等等进行一次规划。
```

## Clarified Requirement

Record Nakama and Pitaya as active reference baselines for vibit's long-term game server capability surface.

This change should:

- Clarify that vibit should eventually cover a comparable game backend/server framework capability surface.
- Preserve vibit's differentiator: Agent-Native architecture, contracts, manifests, generated structure, verification, and agent-operable workflows.
- Add a planning standard that future agents should consult before creating modules, runtime features, protocol features, cluster features, or persistence work.
- Update project entry points so the reference alignment is visible during intake.

## User-Visible Outcome

Future agents should not treat `inventory` as the product's entire ambition.

They should understand that vibit is aiming at the same broad game server problem class as Nakama and Pitaya, while choosing different architecture constraints so AI coding agents can maintain and extend it safely.

## Non-Goals

- Do not copy Nakama or Pitaya APIs directly.
- Do not add new runtime dependencies.
- Do not implement new server features in this change.
- Do not change the ratified Go/WebSocket/Protobuf/PostgreSQL direction.
- Do not introduce distributed clustering before the modular monolith proof slice is healthy.
- Do not weaken existing Agent-Native constraints for faster feature parity.

## Unknowns

- Exact public API names for future modules remain undecided.
- The first production-grade account/session/auth shape remains undecided.
- Cluster and distributed routing remain deferred until single-process module boundaries are proven.
- Whether future compatibility layers should resemble Nakama/Pitaya developer ergonomics remains open.

## Acceptance Criteria

- [ ] A new English reference alignment document exists.
- [ ] A Simplified Chinese translation exists.
- [ ] A new ADR records Nakama/Pitaya as active reference baselines without making them governing standards.
- [ ] `.arch/` has a machine-readable reference alignment manifest.
- [ ] README and agent intake docs point to the new standard.
- [ ] Conversation memory records the maintainer intent.
- [ ] Repository verification passes.
