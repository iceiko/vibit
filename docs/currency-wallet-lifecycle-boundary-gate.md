# Currency Wallet Lifecycle Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-07
Scope: Gate for future currency wallet lifecycle semantics before persistence, protocol, runtime behavior, rewards, purchases, or broader economy features
Depends on: `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/agent-native-feature-request-test-workflow.md`
Canonical decision: `ADR-0200`

The paired Simplified Chinese translation is `docs/currency-wallet-lifecycle-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines the currency wallet lifecycle semantic gate. It is a gate artifact. It does not add runtime currency wallet behavior, balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend execution, audit/event tables, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The currency wallet lifecycle gate record is:

```yaml
currency_wallet_lifecycle_boundary_gate: defined
completed_work_item: W-0292
decision: ADR-0200
check_rule: runtime.currency_wallet_lifecycle_boundary_gate
selection_decision: ADR-0199
gate_standard: docs/currency-wallet-lifecycle-boundary-gate.md
gate_standard_translation: docs/currency-wallet-lifecycle-boundary-gate.zh-CN.md
selected_capability_family: economy_inventory_rewards_currencies_and_progression
selected_module_candidate: currency
primary_product_reference: Nakama
secondary_product_reference: Hiro
pitaya_reference_status: deferred_future_architecture_reference
semantic_gate_only: true
future_persistence_schema_gate_work_item: W-0293
future_persistence_schema_gate_direction: define_currency_wallet_persistence_schema_gate
permission_required: validated_player_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

Currency wallets are a core economy primitive in Nakama-class and Hiro-class game backends. They can later support rewards, grants, spends, purchases, progression, live operations, inventory pricing, leaderboards, tournaments, quests, and player support workflows. vibit adopts the product capability family, not any external public API shape.

The vibit posture is:

- contract-first lifecycle semantics before implementation;
- server-authoritative balance mutations;
- validated player identity for player-owned wallet reads and player-initiated spends;
- service-authoritative posture for future grants and administrative adjustments;
- idempotency and transaction boundaries before any balance mutation exists;
- detailed balances, ledger records, transaction ids, wallet ids, and player ids treated as not log-safe by default;
- direct external API compatibility rejected unless a later ADR authorizes it.

Pitaya remains deferred as a future distributed architecture reference. This gate must not pull in distributed routing, frontend/backend roles, RPC, cluster groups, service discovery, or server-to-server messaging.

## 3. Future Semantic Scope

The future currency wallet lifecycle must cover:

```yaml
semantic_scope:
  - currency_catalog_read
  - wallet_read
  - balance_read
  - grant_currency
  - spend_currency
  - balance_change_recording
  - transaction_read
```

The lifecycle is player-wallet-oriented and server-authoritative. The future domain owner is a currency/economy capability boundary, not WebSocket transport, protocol adapters, authentication, storage objects, inventory, friends, rewards, purchases, leaderboards, matchmaking, match runtime, operations dashboards, or distributed runtime.

Transfers, reservations, settlement, refunds, player-to-player exchange, store purchases, paid currency, and live-operations grants remain future decisions. They may need separate gates after the base wallet lifecycle and persistence schema are stable.

## 4. Future Contract Vocabulary

The future command vocabulary is:

```yaml
commands:
  - GrantCurrency
  - SpendCurrency
```

The future query vocabulary is:

```yaml
queries:
  - GetCurrencyWallet
  - ListCurrencyBalances
  - GetCurrencyTransaction
```

The future event vocabulary is:

```yaml
events:
  - CurrencyGranted
  - CurrencySpent
  - CurrencyBalanceChanged
  - CurrencyTransactionRecorded
```

The future error vocabulary is:

```yaml
errors:
  - CURRENCY_INVALID_CODE
  - CURRENCY_WALLET_NOT_FOUND
  - CURRENCY_INSUFFICIENT_BALANCE
  - CURRENCY_INVALID_AMOUNT
  - CURRENCY_DUPLICATE_TRANSACTION
  - CURRENCY_INVALID_TRANSITION
  - CURRENCY_METADATA_IDENTITY_NOT_AUTHENTICATED
```

This vocabulary is semantic planning. It does not create contract source files, generated shapes, protocol payloads, routes, repositories, migrations, adapters, or runtime handlers.

## 5. Identity And Permissions

Future player-owned wallet reads and player-initiated spends require:

```yaml
permission: validated_player_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_identity_source: validated_request_identity
```

Future grants, administrative adjustments, reward handoffs, purchase settlement, refunds, and live-operations actions require a later service-authoritative permission model before implementation. This gate does not define that model beyond rejecting client-authoritative grants.

Rules:

- The player actor is derived from validated request identity, not from a client-supplied actor id.
- `player_id` and `session_id` metadata are not authentication proof.
- Player-supplied target wallet ids are not authority to mutate another wallet.
- Grants are never client-authoritative.
- Public failures must collapse details when wallet existence, player existence, currency catalog configuration, or private balance information would leak.
- Wallet records, player ids, wallet ids, transaction ids, detailed balances, idempotency keys, and ledger details are not log-safe by default.

## 6. Lifecycle State And Invariants

The future public wallet lifecycle vocabulary is:

```yaml
wallet_lifecycle_states:
  - active
  - suspended
  - closed
```

The future currency balance posture is:

```yaml
balance_posture:
  amount_unit: integer_minor_unit
  negative_balance_allowed_by_default: false
  currency_code_source: future_currency_catalog
  mutation_idempotency_required: true
  mutations_transactional: true
```

The first lifecycle invariants are:

- Currency codes must be validated against a future currency catalog before mutation.
- Amounts must be positive for grants and spends.
- Balance reads must not require mutation authority.
- Balances must never become negative unless a later ADR explicitly authorizes overdraft behavior.
- Every future mutation must carry an idempotency key or equivalent duplicate-prevention key.
- Every future mutation must record a durable transaction or ledger fact before behavior is considered implemented.
- Grant and spend behavior must be transactional with balance changes and transaction recording.
- Spend applies only when the future wallet is active and the balance is sufficient.
- Grant applies only through server/service-authoritative paths defined by a later permission gate.
- Suspended or closed wallet behavior must be explicitly defined before runtime implementation.
- Public query output must not expose internal ledger details beyond the future route contract.

## 7. Future Persistence Gate

The next bounded work item is:

```text
W-0293 Define currency wallet persistence schema gate
```

That follow-up should define table candidates, currency catalog posture, wallet account identity, balance row identity, ledger or transaction record posture, idempotency key uniqueness, indexes, retention, redaction, and repository/adapter boundaries before migration source exists.

This gate intentionally does not decide:

- exact table names;
- whether the first persistent model uses balances plus ledger, ledger-derived balances, or another pattern;
- currency catalog table ownership;
- idempotency key storage shape;
- transaction id format;
- reservation or settlement tables;
- reward, purchase, inventory, or live-operations integration;
- repository interface shape;
- PostgreSQL adapter SQL;
- protocol routes or payloads.

## 8. Future Test Expectations

Future behavior tests must be planned before implementation.

Positive tests:

- grant currency through a service-authoritative path;
- spend currency for a validated player wallet;
- read one player-owned wallet;
- list player-owned currency balances;
- read a currency transaction by permitted owner;
- reject duplicate mutations by idempotency key without double-applying balance changes.

Negative tests:

- invalid currency code;
- zero or negative mutation amount;
- insufficient balance;
- missing wallet;
- suspended or closed wallet transition;
- duplicate transaction with conflicting payload;
- metadata-only identity.

Permission and authentication tests:

- player-owned reads and spends require validated player identity;
- client-supplied actor id is ignored or rejected;
- metadata-only `player_id` and `session_id` are rejected as proof;
- client-authoritative grant attempts are rejected before mutation.

Persistence and transaction tests:

- schema and repository tests must be defined after `W-0293` and before migration/adapter/runtime implementation;
- mutation transitions must be transactional;
- balance changes and transaction records must be consistent within the future unit-of-work boundary;
- duplicate idempotency keys must not double-apply mutations.

Failure and redaction tests:

- public errors do not leak private wallet, balance, catalog, or transaction details where privacy requires collapse;
- logs do not expose raw credentials, tokens, verifier keys, digests, transport metadata, wallet ids, player ids, detailed balances, idempotency keys, or ledger internals.

Concurrency tests:

- simultaneous grant/spend/idempotency conflicts must have explicit expected outcomes before runtime implementation.

Integration and end-to-end tests:

- deferred until protocol routes and runtime handlers are authorized.

## 9. Non-Authorization

This gate does not authorize:

- runtime currency wallet behavior;
- balance tables;
- wallet transaction behavior;
- reward integration;
- inventory integration;
- purchase behavior;
- grant execution;
- spend execution;
- transfer execution;
- reservation or settlement behavior;
- audit/event tables;
- protocol routes;
- Protobuf source;
- generated output;
- migrations;
- repository interfaces;
- PostgreSQL adapters;
- startup wiring;
- authentication/session behavior changes;
- dependencies;
- SDK publication;
- hosted deployment;
- release artifacts;
- distributed runtime behavior;
- direct Nakama/Pitaya API compatibility.

Any future work that needs one of those surfaces must create a separate bounded work item and pass its own repository check.
