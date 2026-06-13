# Currency Wallet Protocol Route Gate

Status: Accepted v0.1
Last updated: 2026-06-08
Scope: Gate-only boundary for future client-facing currency wallet protocol routes after application runtime behavior
Depends on: `docs/currency-wallet-runtime-behavior-gate.md`, `decisions/ADR-0208-currency-wallet-runtime-behavior-implementation.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/generated-output.md`, `docs/bound-identity-route-policy-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0209`

The paired Simplified Chinese translation is `docs/currency-wallet-protocol-route-gate.zh-CN.md`. The English file is authoritative.

This document defines the currency wallet protocol route gate. It is a gate artifact. It does not add protocol route implementation, Protobuf source, generated output, startup wiring, runtime handlers, repository interface changes, PostgreSQL adapter changes, migration changes, dependencies, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, payment behavior, reservation behavior, settlement behavior, refund behavior, transfer behavior, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The currency wallet protocol route gate record is:

```yaml
currency_wallet_protocol_route_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0301
decision: ADR-0209
check_rule: runtime.currency_wallet_protocol_route_gate
source_runtime_behavior_implementation_decision: ADR-0208
source_runtime_behavior_implementation: runtime/internal/app/currency/service.go
source_runtime_behavior_tests: runtime/internal/app/currency/service_test.go
source_runtime_behavior_gate_decision: ADR-0207
source_repository_interface_decision: ADR-0204
repository_interface: runtime/internal/modules/currency.Repository
future_protocol_source_candidate: proto/vibit/currency/v1/currency.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/currency_bridge.go
future_protocol_bridge_test_candidate: runtime/internal/platform/protocol/protobuf/currency_bridge_test.go
future_application_handler_candidate: runtime/internal/app/bootstrap/currency.go
future_application_handler_test_candidate: runtime/internal/app/bootstrap/currency_test.go
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed_as_proof: false
client_supplied_wallet_id_allowed_for_lookup: false
client_supplied_actor_id_allowed_as_proof: false
server_authoritative_grant_policy_required: true
first_owner_kind: player
first_actor_kinds:
  - player
  - system
first_payload_package: vibit.currency.v1
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
protocol_route_gate_only: true
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
runtime_handler_added: false
startup_wiring_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
reward_integration_added: false
inventory_integration_added: false
purchase_behavior_added: false
currency_catalog_table_added: false
currency_wallet_events_table_added: false
payment_behavior_added: false
reservation_behavior_added: false
settlement_behavior_added: false
refund_behavior_added: false
transfer_behavior_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_protocol_route_implementation_work_item: W-0302
future_protocol_route_implementation_direction: implement_currency_wallet_protocol_route
```

## 2. Purpose

`W-0300` implemented application-owned currency wallet behavior under `runtime/internal/app/currency`. The next useful boundary is not route code or `.proto` generation. The next useful boundary is a protocol route gate that records how future WebSocket/Protobuf exposure should call that service without moving wallet behavior into transport, generated files, persistence adapters, reward systems, purchase systems, or payment integrations.

Nakama motivates the product surface: durable wallets, balances, grants, spends, and transaction history are common game backend economy capabilities. vibit should cover that capability class.

Pitaya motivates the architecture posture: acceptors, sessions, route handlers, serializers, and backend services should remain separated. vibit adapts that by keeping WebSocket transport credential-neutral, keeping Protobuf payload bridging explicit, and invoking application-owned route handlers that call application-owned currency wallet services.

This gate records the future protocol shape before implementation:

- candidate route names;
- candidate request and response message shapes;
- route protection and identity handoff posture;
- protocol adapter, application handler, and startup ownership;
- generated-output expectations;
- public error mapping and redaction expectations;
- local proof expectations;
- Nakama/Pitaya reference mapping;
- stop conditions that keep implementation and generated artifacts out of this slice.

## 3. Future Route Surface

The first route family should expose own-player wallet operations only:

```yaml
candidate_routes:
  - kind: command
    module: currency
    name: EnsurePlayerWallet
    route_id: currency.EnsurePlayerWallet
    service_method: EnsurePlayerWallet
  - kind: query
    module: currency
    name: GetOwnWallet
    route_id: currency.GetOwnWallet
    service_method: GetOwnWallet
  - kind: query
    module: currency
    name: ListOwnWalletBalances
    route_id: currency.ListOwnWalletBalances
    service_method: ListOwnWalletBalances
  - kind: command
    module: currency
    name: GrantCurrency
    route_id: currency.GrantCurrency
    service_method: GrantCurrency
  - kind: command
    module: currency
    name: SpendCurrency
    route_id: currency.SpendCurrency
    service_method: SpendCurrency
  - kind: query
    module: currency
    name: ListOwnWalletTransactions
    route_id: currency.ListOwnWalletTransactions
    service_method: ListOwnWalletTransactions
```

Rules:

- The route names must stay vibit-native and must not copy Nakama route paths or Pitaya route naming.
- `EnsurePlayerWallet`, `GrantCurrency`, and `SpendCurrency` are commands.
- `GetOwnWallet`, `ListOwnWalletBalances`, and `ListOwnWalletTransactions` are queries.
- The first route family is only for the validated player wallet owner. It must not expose arbitrary owner ids or arbitrary wallet id lookup.
- `GrantCurrency` remains server-authoritative in service behavior. Future route implementation must explicitly decide whether it is local-proof/dev-only, admin-only, system-initiated through application composition, or excluded from public client registration. It must not become an unauthenticated or client-authoritative mint route.
- Client payloads may provide currency code, positive amount, idempotency key/scope, reason code, external reference, metadata JSON, expected wallet version, expected balance version, list limit, currency pagination cursor, and transaction pagination cursor where the service already has vocabulary for them.
- Client payloads must not provide owner ids, wallet ids for lookup, actor ids as proof, raw access tokens, credential material, lookup digests, verifier digests, SQL details, payment provider payloads, catalog ids with purchase semantics, reward ids with reward execution semantics, inventory mutation fields, or direct external API compatibility markers.
- Reward execution, inventory integration, purchase behavior, currency catalogs, event/audit streams, payment providers, reserves, settlement, refunds, transfers, operations/admin inspection, script hooks, SDK/client libraries, hosted deployments, distributed runtime routing, and direct compatibility remain deferred.
- Future route implementation must register routes explicitly. No catch-all currency route or reflective handler is allowed.

## 4. Protocol Shape

The first currency wallet protocol source candidate is:

```text
proto/vibit/currency/v1/currency.proto
```

The first generated output candidate is:

```text
runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go
```

The first Protobuf package candidate is:

```text
vibit.currency.v1
```

Candidate messages:

```yaml
messages:
  CurrencyWallet:
    fields:
      wallet_id: string
      owner_kind: string
      lifecycle_state: string
      wallet_version: int64
      created_at: string
      updated_at: string
      state_changed_at: string
  CurrencyWalletBalance:
    fields:
      currency_code: string
      balance_amount: int64
      balance_version: int64
      created_at: string
      updated_at: string
  CurrencyWalletTransaction:
    fields:
      transaction_id: string
      currency_code: string
      transaction_kind: string
      amount_delta: int64
      balance_after: int64
      actor_kind: string
      reason_code: string
      external_reference: string
      metadata_json: string
      created_at: string
  EnsurePlayerWalletRequest:
    fields: {}
  EnsurePlayerWalletResponse:
    fields:
      wallet: CurrencyWallet
      status: string
  GetOwnWalletRequest:
    fields: {}
  GetOwnWalletResponse:
    fields:
      wallet: CurrencyWallet
      status: string
  ListOwnWalletBalancesRequest:
    fields:
      limit: int32
      after_currency_code: string
  ListOwnWalletBalancesResponse:
    fields:
      balances: repeated CurrencyWalletBalance
      next_currency_code: string
      status: string
  GrantCurrencyRequest:
    fields:
      currency_code: string
      amount: int64
      idempotency_key: string
      idempotency_scope: string
      reason_code: string
      external_reference: string
      metadata_json: string
      expected_wallet_version: int64
      expected_balance_version: int64
  GrantCurrencyResponse:
    fields:
      transaction: CurrencyWalletTransaction
      status: string
  SpendCurrencyRequest:
    fields:
      currency_code: string
      amount: int64
      idempotency_key: string
      idempotency_scope: string
      reason_code: string
      external_reference: string
      metadata_json: string
      expected_wallet_version: int64
      expected_balance_version: int64
  SpendCurrencyResponse:
    fields:
      transaction: CurrencyWalletTransaction
      status: string
  ListOwnWalletTransactionsRequest:
    fields:
      currency_code: string
      limit: int32
      after_transaction_id: string
      after_transaction_time: string
  ListOwnWalletTransactionsResponse:
    fields:
      transactions: repeated CurrencyWalletTransaction
      next_transaction_id: string
      next_transaction_time: string
      status: string
```

Rules:

- The existing `proto/vibit/protocol/v1/envelope.proto` must remain unchanged unless a later protocol ADR explicitly changes envelope semantics.
- Time values should use RFC3339Nano UTC text when exposed.
- Optional expected-version mapping must preserve the service's optional expected-version vocabulary. Future implementation must make field absence vs `0` semantics explicit in tests.
- Amount inputs are positive integer minor units for grant and spend requests. Transaction deltas remain service/repository-owned and must not let clients submit negative spend deltas.
- `metadata_json` is not log-safe. It must not appear in default error messages, logs, route policy diagnostics, or test names.
- Wallet ids, owner ids, transaction ids, actor ids, idempotency keys, reason codes, external references, balance amounts, metadata JSON, and payment-adjacent fields are not log-safe by default.
- The protocol shape must not include `owner_id`, raw `wallet_id` lookup, client-supplied actor id, raw access tokens, credential material, lookup digests, verifier digests, SQL details, payment provider payloads, catalog purchase fields, reward execution fields, inventory mutation fields, private transport metadata, or direct external API compatibility markers.

## 5. Route Protection And Identity Handoff

The first route-policy posture is:

```yaml
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed_as_proof: false
client_supplied_wallet_id_allowed_for_lookup: false
client_supplied_actor_id_allowed_as_proof: false
server_authoritative_grant_policy_required: true
```

Rules:

- Future currency wallet routes must be protected gameplay routes.
- Future handlers must receive a validated `app.RequestIdentity` from the existing protected-route flow.
- Metadata-only `player_id` or `session_id` from envelope/session metadata must never become wallet owner proof.
- Client payloads must not choose owner ids, wallet ids for lookup, or actor ids as proof.
- The service remains responsible for rejecting invalid identity before unit-of-work access or repository access.
- The WebSocket transport remains credential-neutral.
- This gate does not change authentication, token validation, session persistence, first-message binding, WebSocket handshake behavior, bound-identity policy, route-protection semantics, or permission schemas.

## 6. Future Route Flow

Future implementation must preserve this sequence:

```yaml
future_route_flow:
  - receive WebSocket/Protobuf envelope through existing request path
  - apply protected-route authenticated wrapper policy
  - obtain validated authenticated request identity
  - decode vibit.currency.v1 request payload
  - reject payload owner proof, wallet lookup proof, and actor proof
  - map payload fields to runtime/internal/app/currency service request
  - call application-owned currency wallet service
  - map service result to vibit.currency.v1 response payload
  - map service public errors to protocol error responses
  - keep transport, Protobuf bridge, application handler, service, repository, and PostgreSQL adapter ownership separated
```

Rules:

- WebSocket transport remains credential-neutral and currency-neutral.
- The Protobuf bridge should map payloads and response shapes only; it must not own wallet behavior, minting policy, spend policy, permission decisions, repository calls, or payment logic.
- Application bootstrap handlers should own route registration and service invocation.
- The application service remains the owner of identity checks, validation handoff, repository conflict mapping, actor derivation, idempotency handoff, and public service errors.
- PostgreSQL adapters remain persistence-only.

## 7. Public Error Mapping

Future route implementation should map service public errors without leaking internal details:

```yaml
public_error_mapping:
  CURRENCY_WALLET_INVALID_REQUEST: invalid_request
  CURRENCY_WALLET_UNAUTHENTICATED: unauthenticated
  CURRENCY_WALLET_NOT_FOUND: not_found
  CURRENCY_WALLET_ALREADY_EXISTS: already_exists
  CURRENCY_WALLET_NOT_ACTIVE: wallet_not_active
  CURRENCY_WALLET_INSUFFICIENT_BALANCE: insufficient_balance
  CURRENCY_WALLET_DUPLICATE_TRANSACTION: duplicate_transaction
  CURRENCY_WALLET_VERSION_MISMATCH: version_mismatch
  CURRENCY_WALLET_UNAVAILABLE: unavailable
```

Rules:

- Public protocol errors may expose only stable public codes/classes and retryability posture if a later ADR authorizes that field.
- Not-found, owner mismatch, inactive-wallet, insufficient-balance, and duplicate/idempotency cases must not leak cross-player wallet existence or private transaction details.
- Internal repository errors, SQL details, wallet ids, owner ids, actor ids, transaction ids, idempotency keys, reason codes, external references, metadata JSON, access-token material, credential material, lookup digests, verifier digests, payment provider payloads, and transport metadata must remain out of default logs and error messages.
- Authentication and route protection failures must use existing protected-route semantics. This gate does not invent a new proof carrier.

## 8. Generated Output Posture

Future generated output must follow `docs/generated-output.md`.

Rules:

- `proto/vibit/currency/v1/currency.proto` may be added only by a later implementation slice.
- `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go` may be added only as generated output from Buf/protoc.
- Generated Go output must contain the `protoc-gen-go` generated-code marker and trace to the source `.proto`.
- Agents must not hand-edit generated Go Protobuf files.
- This gate does not change `buf.yaml`, `buf.gen.yaml`, or generated output.

## 9. Ownership

Future implementation must preserve these owners:

```yaml
currency_service_owner: runtime/internal/app/currency
application_handler_owner: runtime/internal/app/bootstrap
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
currency_repository_interface_owner: runtime/internal/modules/currency
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
websocket_transport_owner: runtime/internal/platform/transport/ws
startup_owner: runtime/cmd/vibit-server
generated_output_owner: runtime/internal/generated/proto/vibit/currency/v1
protobuf_source_owner: proto/vibit/currency/v1
```

Rules:

- Currency wallet runtime behavior remains in `runtime/internal/app/currency`.
- Protocol bridge code may map payload fields only.
- Persistence code remains currency adapter behavior only.
- Generated output must be produced from `.proto` sources and must not be hand-edited.
- Startup wiring, route registration, and generated output remain behind the later implementation work item.

## 10. Nakama Reference Mapping

Nakama reference mapping:

```yaml
nakama_reference_mapping:
  capability_family: economy_wallets_and_ledgers
  mapped_capabilities:
    - wallet_create_or_get
    - wallet_balance_read
    - wallet_balance_list
    - currency_grant
    - currency_spend
    - wallet_transaction_history
  direct_api_compatibility: false
```

Nakama informs the useful capability class. vibit does not copy Nakama route paths, API names, runtime script APIs, storage model names, permission semantics, economy ledger semantics, or public API compatibility.

## 11. Pitaya Reference Mapping

Pitaya reference mapping:

```yaml
pitaya_reference_mapping:
  architecture_pressure:
    - acceptor_session_handler_separation
    - serializer_adapter_separation
    - backend_service_boundary
  distributed_architecture_status: deferred
  direct_api_compatibility: false
```

Pitaya informs layering pressure. This gate does not add Pitaya-style distributed topology, frontend/backend split, RPC, groups, service discovery, distributed wallet routing, or direct Pitaya API compatibility.

## 12. Required Future Tests

Future implementation tests should cover:

- route registration for all selected route ids;
- command/query kind mapping;
- Protobuf request and response bridge mapping;
- optional expected-version mapping;
- owner derivation from validated request identity;
- rejection of metadata-only identity through the existing protected-route wrapper;
- rejection of client-supplied owner id, wallet id lookup proof, and actor proof;
- grant route policy cannot become unauthenticated or client-authoritative minting;
- spend route maps positive amount input to service-owned spend behavior;
- balance and transaction pagination mapping;
- metadata JSON redaction;
- public `CURRENCY_WALLET_*` error to protocol error mapping;
- redaction of private wallet, balance, transaction, idempotency, token, credential, SQL, payment, and transport details;
- no route behavior in WebSocket transport or PostgreSQL adapter packages;
- generated-output traceability if Protobuf source is added;
- local proof expectations for ensure, get, list balances, spend, and transaction listing when a later proof slice authorizes it.

## 13. Stop Conditions

Stop and create a separate work item before adding:

- protocol route implementation;
- `proto/vibit/currency/v1/currency.proto`;
- generated Go Protobuf output;
- protocol bridge implementation;
- application bootstrap handlers;
- startup route registration;
- new dependencies;
- migration changes;
- repository interface changes;
- PostgreSQL adapter changes;
- authentication/session behavior changes;
- route-protection semantic changes;
- reward integration;
- inventory integration;
- purchase behavior;
- currency catalog tables;
- event/audit tables;
- payment behavior;
- reservation behavior;
- settlement behavior;
- refund behavior;
- transfer behavior;
- operations/admin behavior;
- SDK publication;
- generated client libraries;
- hosted deployments;
- release artifacts;
- public announcements;
- paid promotion;
- Pitaya-style distributed architecture;
- direct Nakama/Pitaya API compatibility.

## 14. Verification

The repository check rule for this gate is:

```text
runtime.currency_wallet_protocol_route_gate
```

Recommended verification after this gate:

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.currency_wallet_protocol_route_gate
node tools/vibit check change define-currency-wallet-protocol-route-gate --json
node tools/vibit check module currency --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check protocol --json
node tools/vibit check schemas --json
node tools/vibit check memory --json
node tools/vibit check all --json
cd runtime && go test ./...
git diff --check
```

Go tests are not required by this gate because no Go runtime behavior is added, but a full runtime test run remains useful before closing a development turn.
