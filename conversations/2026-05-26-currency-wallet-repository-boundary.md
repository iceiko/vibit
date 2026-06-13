# Conversation: Currency Wallet Repository Boundary

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-repository-boundary/`

Related artifacts:

- `docs/currency-wallet-repository-boundary.md`
- `docs/currency-wallet-repository-boundary.zh-CN.md`
- `decisions/ADR-0203-currency-wallet-repository-boundary.md`
- `runtime/migrations/postgres/000008_create_currency_wallets.sql`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`
- `runtime/internal/modules/currency.Repository`

## Context

`W-0294` added the currency wallet migration source and opened `W-0295 Define currency wallet repository boundary`.

## Maintainer Narrative

The maintainer asked to continue 20 steps. The current next-ready item allowed defining only the repository boundary and required stopping before repository interfaces, PostgreSQL adapters, runtime behavior, protocol shape, generated output, dependencies, hosted surfaces, SDKs, distributed runtime, or direct compatibility.

## Agent Response Summary

The agent treated `W-0295` as a repository-boundary gate only. It defined the future storage-neutral repository vocabulary for currency wallets, balances, transactions, idempotency, conflicts, redaction, and PostgreSQL adapter expectations.

The agent did not create `runtime/internal/modules/currency`, Go repository interfaces, PostgreSQL adapters, SQL execution code, runtime wallet behavior, grant/spend execution, protocol routes, Protobuf sources, generated output, dependencies, startup wiring, authentication/session behavior changes, catalog tables, event/audit tables, rewards, purchases, inventory integration, hosted surfaces, SDKs, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Decisions

- Accepted `ADR-0203`.
- Registered `runtime.currency_wallet_repository_boundary`.
- Defined `docs/currency-wallet-repository-boundary.md`.
- Added the paired Simplified Chinese translation.
- Opened `M-224/W-0296 Implement storage-neutral currency wallet repository interface` as the next-ready work item.

## Artifacts

- `docs/currency-wallet-repository-boundary.md`
- `docs/currency-wallet-repository-boundary.zh-CN.md`
- `decisions/ADR-0203-currency-wallet-repository-boundary.md`
- `changes/2026-05-26-define-currency-wallet-repository-boundary/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Exact Go repository interface names and method signatures are deferred to `W-0296`.
- PostgreSQL adapter SQL mapping remains deferred.
- Runtime grant/spend behavior remains deferred.
- Protocol routes and Protobuf payloads remain deferred.
- Currency catalog and event/audit tables remain deferred.
- Rewards, purchases, inventory integration, reservations, settlement, refunds, and transfers remain deferred.

## Follow-Up

- Complete `W-0296 Implement storage-neutral currency wallet repository interface`.
- Keep PostgreSQL adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, rewards, purchases, inventory integration, event/audit tables, hosted surfaces, SDKs, distributed runtime, and direct compatibility behind later bounded work items.

## Redaction Notes

No raw credentials, raw access tokens, verifier keys, digests, DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, payment secrets, private wallet ids, detailed balances, idempotency keys, external references, or transaction metadata were recorded beyond explicit planning vocabulary.
