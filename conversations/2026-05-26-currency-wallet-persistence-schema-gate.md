# Conversation: Currency Wallet Persistence Schema Gate

Date: 2026-05-26

## Context

The maintainer asked to continue 20 steps from the current repository queue. `W-0292` completed the currency wallet lifecycle boundary gate and opened `W-0293 Define currency wallet persistence schema gate`.

## Maintainer Narrative

The maintainer wants sustained progress on Nakama/Hiro-class economy capability coverage while keeping vibit's architecture optimized for agent-native, requirement-to-test development. The work should advance in bounded gates and avoid implementation before schema and contract posture are explicit.

## Agent Response Summary

The agent treated `W-0293` as a schema gate only. It defined future PostgreSQL currency wallet persistence posture before SQL migration source, repositories, adapters, runtime behavior, protocol routes, generated output, dependencies, or broader economy features.

The selected posture uses future `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` table candidates, player-owned wallet identity, balance row identity, transaction or ledger facts, idempotency uniqueness, catalog deferral, timestamp and redaction posture, and future repository/adapter ownership candidates.

## Decisions

- Accepted `ADR-0201`.
- Registered `runtime.currency_wallet_persistence_schema_gate`.
- Selected `runtime/migrations/postgres/000008_create_currency_wallets.sql` as the future migration source candidate.
- Opened `M-222/W-0294 Add currency wallet migration source` as the next-ready work item.

## Artifacts

- `docs/currency-wallet-persistence-schema-gate.md`
- `docs/currency-wallet-persistence-schema-gate.zh-CN.md`
- `decisions/ADR-0201-currency-wallet-persistence-schema-gate.md`
- `changes/2026-05-26-define-currency-wallet-persistence-schema-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Exact SQL checks and index names are deferred to `W-0294`.
- Currency catalog table ownership is deferred.
- Repository interface shape is deferred.
- PostgreSQL adapter behavior is deferred.
- Runtime grant/spend behavior and protocol routes are deferred.
- Event/audit tables, rewards, purchases, inventory integration, reservations, settlement, refunds, and transfers are deferred.

## Follow-Up

- Complete `W-0294 Add currency wallet migration source`.
- Keep repository interfaces, adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, rewards, purchases, inventory integration, event/audit tables, hosted surfaces, SDKs, distributed runtime, and direct compatibility behind later bounded work items.

## Redaction Notes

No raw credentials, raw access tokens, verifier keys, digests, DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, payment secrets, private wallet ids, detailed balances, idempotency keys, or transaction metadata were recorded beyond explicit planning vocabulary.
