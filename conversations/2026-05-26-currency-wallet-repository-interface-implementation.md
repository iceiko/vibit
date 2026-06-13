# Conversation: Currency Wallet Repository Interface Implementation

Date: 2026-06-07
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-repository-interface/`

Related artifacts:

- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/modules/currency/repository_test.go`
- `modules/currency/module.yaml`
- `modules/currency/AGENTS.md`
- `modules/currency/AGENTS.zh-CN.md`
- `decisions/ADR-0204-currency-wallet-repository-interface-implementation.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `rules/check-rules.json`
- `tools/vibit`

## Context

`M-223/W-0295` completed the currency wallet repository boundary. It recorded `runtime/internal/modules/currency` as the future repository owner candidate and kept PostgreSQL adapters, SQL execution, runtime behavior, protocol shape, generated output, dependencies, authentication/session changes, rewards, inventory integration, purchases, hosted surfaces, SDKs, distributed runtime, and direct compatibility deferred.

The next-ready work item was `W-0296 Implement storage-neutral currency wallet repository interface`.

## Maintainer Narrative

The maintainer asked to continue multiple steps:

```text
继续20步
```

The standing product direction was Nakama/Hiro-aligned economy capability coverage, with Pitaya deferred as a future distributed architecture reference, and with AI-native requirement-to-spec-to-tests-to-implementation-to-verification as the product design goal.

## Agent Response Summary

The agent continued the selected bounded work item and implemented the storage-neutral currency wallet repository interface under `runtime/internal/modules/currency`.

The work added:

- currency wallet, balance, transaction, owner, actor, idempotency, amount, and version value types;
- closed owner kind, lifecycle state, transaction kind, actor kind, and conflict vocabularies;
- create, lookup, owner lookup, list balance, record grant, record spend, and list transaction repository input/result types;
- normalization helpers for records, list results, idempotency inputs, metadata JSON, actor identity, and expected versions;
- optimistic and idempotency conflict vocabulary and redacted repository errors;
- focused Go tests;
- currency module manifest and module AGENTS guides;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-224/W-0296`.
- Accept `ADR-0204`.
- Add `runtime.currency_wallet_repository_interface_implementation`.
- Register the `currency` module in `.arch/modules.yaml`.
- Keep the first implementation storage-neutral and module-owned.
- Select `M-225/W-0297 Define currency wallet PostgreSQL adapter gate` as the next-ready work item.

## Nakama, Hiro, And Pitaya Reference Basis

Nakama and Hiro guided the capability pressure: currency wallets are a common game/backend economy primitive for balances, grants, spends, and transaction records.

Pitaya remained deferred; no distributed topology, RPC routing, frontend/backend split, group broadcast, or service discovery behavior was added.

vibit adapted those lessons into its own model: an explicit module-owned repository interface with no direct public API compatibility and no runtime/protocol behavior in this slice.

## Artifacts

- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/modules/currency/repository_test.go`
- `modules/currency/module.yaml`
- `modules/currency/AGENTS.md`
- `modules/currency/AGENTS.zh-CN.md`
- `decisions/ADR-0204-currency-wallet-repository-interface-implementation.md`
- `changes/2026-05-26-implement-currency-wallet-repository-interface/`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- PostgreSQL adapter mapping remains deferred to `W-0297` and later implementation work.
- Runtime wallet creation, grant, spend, balance read, transaction read, and permission behavior remain deferred.
- Protocol routes and Protobuf messages remain deferred.
- Reward, inventory, purchase, payment, catalog, event/audit, reservation, settlement, refund, transfer, SDK, hosted surface, distributed runtime, and direct compatibility scope remain deferred.

## Follow-Up

- Define the currency wallet PostgreSQL adapter gate.
- Only after that gate, implement the adapter in a separate bounded slice.
- Only after repository and adapter boundaries, define runtime behavior and protocol routes.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, payment provider payloads, raw wallet records from a real user, or raw private economy data are recorded in this conversation log.
