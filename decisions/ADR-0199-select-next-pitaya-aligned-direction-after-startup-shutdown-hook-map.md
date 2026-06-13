# ADR-0199: Select Next Pitaya-Aligned Direction After Startup And Shutdown Hook Map

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map/`

Related conversations:

- `conversations/2026-05-26-select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map.md`

Related artifacts:

- `decisions/ADR-0198-pitaya-aligned-startup-shutdown-hook-source-first-map.md`
- `decisions/ADR-0197-pitaya-aligned-startup-shutdown-hook-boundary-gate.md`
- `decisions/ADR-0196-select-next-pitaya-aligned-direction-after-component-discovery-module-loading-map.md`
- `docs/reference-game-server-alignment.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0198` completed the source-first Pitaya-aligned startup and shutdown hook map through `node tools/vibit inspect pitaya-startup-shutdown --json`.

The repository has now completed a bounded Pitaya-aligned vocabulary pass for distributed runtime vocabulary, frontend/backend roles, server-to-server RPC, service discovery, distributed groups and broadcast, cluster-safe session routing, route handler pipelines, serializer and message forwarding, acceptor and connection lifecycle, session binding and data, observability, metrics and tracing, dashboard/admin operations, runtime component lifecycle, handler module registration, component discovery and module loading, and startup and shutdown hooks. Those maps are source-first or gate-only and do not authorize concrete Pitaya runtime behavior.

The roadmap then returns to core game backend modules. The Phase 3 candidate list includes Currency and wallets, rewards and claims, friends/groups/parties, chat, leaderboards, tournaments, and quests/progression. Currency and wallets are a useful next boundary because they require explicit invariants before any balance table, repository, transaction behavior, protocol shape, or runtime behavior is safe to add.

## Decision

Select `define_currency_wallet_lifecycle_boundary_gate` as the next bounded direction after the startup and shutdown hook source-first map.

Register `runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map` as the repository check rule.

Complete `M-219/W-0291` and open `M-220/W-0292 Define currency wallet lifecycle boundary gate` as next-ready.

This decision does not add currency wallet behavior. It also does not add balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend behavior, audit/event tables, startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, runtime endpoint behavior, dashboard behavior, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Continue adding more Pitaya-aligned runtime architecture vocabulary after startup and shutdown hooks.
- Define reward and claim lifecycle boundaries first.
- Define leaderboard and tournament lifecycle boundaries first.
- Define groups, parties, or chat boundaries first.
- Implement a currency or wallet module immediately.
- Add balance tables, wallet transaction behavior, or economy protocol messages immediately.

## Rationale

The startup and shutdown hook source-first map is a natural stopping point for the Pitaya vocabulary pass. Continuing into concrete lifecycle behavior would cross an implementation boundary, and continuing indefinitely with architecture vocabulary would delay Phase 3 product capability coverage.

Currency and wallets are listed in the reference matrix as a future `currency` module with transactional invariants. A boundary gate is the smallest safe follow-up: it can define vocabulary, ownership, invariants, Nakama/Hiro capability mapping, and deferrals before any persistent balance model, repository interface, adapter, protocol route, generated shape, runtime behavior, reward integration, dependency, or direct compatibility is introduced.

## Agent Reasoning Summary

The active work item is a selection-only continuation. The correct continuation is to record the selected follow-up direction, update repository memory, register the check rule, and open W-0292. Currency wallet behavior, balance tables, wallet transaction behavior, rewards, inventory integration, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime, startup/shutdown hook behavior, lifecycle execution, dependency containers, and direct compatibility remain deferred until later gates and implementation slices explicitly authorize them.

## Decision Weights

```yaml
decision_weights:
  phase_3_product_capability_value: high
  currency_wallet_invariant_clarity: high
  implementation_boundedness: high
  pitaya_vocabulary_completion_value: high
  nakama_hiro_alignment_value: high
  runtime_behavior_risk: none_in_this_step
  persistence_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

W-0292 becomes the next-ready work item and must define only a currency wallet lifecycle boundary gate.

This selection does not authorize currency wallet behavior, balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend behavior, audit/event tables, startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency containers, runtime endpoints, dashboards, admin console behavior, protocol changes, generated output, persistence changes, dependencies, hosted surfaces, SDK publication, release artifacts, distributed runtime behavior, or direct compatibility.

## Reversal Conditions

Revisit this decision if:

- a later ADR chooses a different Phase 3 core game backend module before currency and wallets;
- currency and wallet vocabulary needs to be split into separate gates for currency catalog, wallet ledger, transaction idempotency, grant/spend rules, inventory integration, reward integration, or audit/event behavior;
- the roadmap reclassifies currency and wallets behind rewards, leaderboards, groups, parties, chat, matchmaking, or another product capability family;
- a maintainer explicitly pauses Phase 3 product capability coverage and reopens Pitaya architecture vocabulary first.

## Follow-Up

- Complete `W-0292`: define the currency wallet lifecycle boundary gate.
- Keep currency wallet behavior, balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend behavior, audit/event tables, protocol shape, generated output, persistence, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, startup/shutdown hook behavior, lifecycle execution, dependency containers, and direct compatibility behind later bounded work items.
