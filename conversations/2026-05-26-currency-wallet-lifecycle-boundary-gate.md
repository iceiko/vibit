# Conversation: Currency Wallet Lifecycle Boundary Gate

Date: 2026-06-07
Work item: W-0292
Milestone: M-220
Decision: ADR-0200
Change: `changes/2026-05-26-define-currency-wallet-lifecycle-boundary-gate`
Rule: `runtime.currency_wallet_lifecycle_boundary_gate`

## Context

Continue repository work from the current `next_ready` item for up to 20 steps.

`W-0291 Select next Pitaya-aligned direction after startup and shutdown hook map` completed ADR-0199, registered `runtime.next_pitaya_aligned_direction_after_startup_shutdown_hook_map`, selected `define_currency_wallet_lifecycle_boundary_gate`, and opened `W-0292` as the next-ready work item.

The RED checks confirmed the expected missing W-0292 surfaces:

```text
node tools/vibit inspect rule runtime.currency_wallet_lifecycle_boundary_gate
# Unknown rule_id: runtime.currency_wallet_lifecycle_boundary_gate

node tools/vibit check change define-currency-wallet-lifecycle-boundary-gate --json
# change directory does not exist
```

`node tools/vibit inspect next --json` reported `W-0292` as next-ready.

## Maintainer Narrative

The maintainer asked to continue 20 steps. The repository queue selected `W-0292 Define currency wallet lifecycle boundary gate`.

## Agent Response Summary

The agent defined the currency wallet lifecycle boundary gate as a semantic-only planning artifact for the future currency/economy module candidate.

This completes `W-0292`, accepts `ADR-0200`, registers `runtime.currency_wallet_lifecycle_boundary_gate`, and opens `M-221/W-0293 Define currency wallet persistence schema gate` as next-ready.

## Decisions

- Accept `ADR-0200`.
- Complete `W-0292`.
- Open `W-0293` as next-ready.
- Define future currency catalog read, wallet read, balance read, grant currency, spend currency, balance change recording, and transaction read semantics.
- Require validated player identity for future player-owned wallet reads and player-initiated spends.
- Keep grants service-authoritative and reject client-authoritative grants.
- Treat wallet ids, player ids, transaction ids, detailed balances, idempotency keys, and ledger details as not log-safe by default.
- Keep currency wallet behavior, balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend execution, audit/event tables, protocol routes, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime, and direct compatibility deferred.

## Artifacts

- `docs/currency-wallet-lifecycle-boundary-gate.md`
- `docs/currency-wallet-lifecycle-boundary-gate.zh-CN.md`
- `decisions/ADR-0200-currency-wallet-lifecycle-boundary-gate.md`
- `changes/2026-05-26-define-currency-wallet-lifecycle-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Open Questions

No open questions for `W-0292`. The selected next work item is bounded and ready.

## Follow-Up

Proceed to `W-0293 Define currency wallet persistence schema gate`.

## Redaction Notes

No ignored credential file contents were read or printed. No secrets, credentials, raw access tokens, verifier material, lookup digests, verifier digests, DSNs with credentials, database payloads, transport payloads, local secret values, route payloads, session data payloads, dashboard payloads, admin console payloads, component lifecycle payloads, handler registration payloads, component module payloads, startup hook payloads, shutdown hook payloads, currency wallet payloads, wallet transaction payloads, idempotency keys, ledger details, or concrete operational payloads are recorded.
