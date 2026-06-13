# Conversation: Currency Wallet Migration Source

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-add-currency-wallet-migration-source/`

Related artifacts:

- `runtime/migrations/postgres/000008_create_currency_wallets.sql`
- `decisions/ADR-0202-currency-wallet-migration-source.md`
- `docs/currency-wallet-persistence-schema-gate.md`
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

## Context

`W-0293` defined the currency wallet persistence schema gate. That gate selected PostgreSQL, future `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` table candidates, player-owned wallet identity, balance identity, transaction facts, idempotency uniqueness, catalog deferral, timestamp checks, redaction posture, and future repository/adapter ownership candidates.

The work queue identified `W-0294 Add currency wallet migration source` as next-ready. The maintainer asked to continue 20 steps, with continuation meaning advance next-ready work items unless an ask-first boundary, verification failure, or required maintainer confirmation blocks progress.

## Maintainer Narrative

The maintainer said:

```text
继续20步
```

The request was interpreted as permission to continue through bounded next-ready items in dependency order. `W-0294` allowed only the migration source, so the slice stayed migration-source-only.

## Agent Response Summary

The agent added `runtime/migrations/postgres/000008_create_currency_wallets.sql`. The migration creates only the PostgreSQL source for `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`, with player-owned wallet identity, balance row identity, lifecycle state constraints, non-negative balances, transaction kind direction checks, durable idempotency uniqueness, and lookup indexes.

The agent also recorded `ADR-0202`, added change-spec artifacts, updated work-item and architecture manifests, extended `tools/vibit` static checks, and advanced the queue to `W-0295 Define currency wallet repository boundary`.

The agent did not add currency wallet repository interfaces, PostgreSQL adapters, SQL execution code, runtime wallet behavior, grant/spend execution, reward integration, inventory integration, purchase behavior, protocol routes, Protobuf sources, generated output, dependencies, startup wiring, authentication/session behavior changes, hosted surfaces, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Decisions

- Complete `M-222/W-0294` as a migration-source-only slice.
- Accept `ADR-0202`.
- Add `runtime.currency_wallet_migration_source` as the repository check rule.
- Add only `runtime/migrations/postgres/000008_create_currency_wallets.sql`.
- Advance the queue to `M-223/W-0295 Define currency wallet repository boundary`.
- Continue to defer currency wallet repositories, adapters, protocol, runtime behavior, catalog tables, event/audit tables, rewards, purchases, inventory integration, SDKs, hosted surfaces, distributed runtime, and direct compatibility.

## Artifacts

- `changes/2026-05-26-add-currency-wallet-migration-source/`
- `runtime/migrations/postgres/000008_create_currency_wallets.sql`
- `decisions/ADR-0202-currency-wallet-migration-source.md`
- `conversations/2026-05-26-currency-wallet-migration-source.md`
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

## Open Questions

- Storage-neutral repository interface shape remains deferred.
- PostgreSQL adapter mapping and unit-of-work exposure remain deferred.
- Runtime grant/spend command/query behavior and conflict handling remain deferred.
- Protocol route and Protobuf payload shape remain deferred.
- Currency catalog ownership remains deferred.
- Event/audit history, rewards, purchases, inventory integration, reservations, settlement, refunds, and transfers remain separate future capability work.

## Follow-Up

- Start `W-0295 Define currency wallet repository boundary`.
- Keep the next slice gate-only unless explicitly authorized to implement repository interfaces.
- Run repository checks and record any accepted warnings explicitly.

## Redaction Notes

No secrets, raw device credentials, raw access tokens, lookup digests, verifier digests, verifier keys, DSNs, cookies, query strings, WebSocket subprotocol values, remote addresses, payment provider secrets, full payment provider payloads, private wallet ids, detailed balances, idempotency keys, transaction metadata, or GitHub tokens are recorded in this conversation log.
