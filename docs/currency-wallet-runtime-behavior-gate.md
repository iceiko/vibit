# Currency Wallet Runtime Behavior Gate

Status: Accepted v0.1
Last updated: 2026-06-07
Scope: Gate-only boundary for future application-owned currency wallet runtime behavior after the PostgreSQL adapter
Depends on: `docs/currency-wallet-lifecycle-boundary-gate.md`, `docs/currency-wallet-repository-boundary.md`, `docs/currency-wallet-postgresql-adapter-gate.md`, `runtime/internal/modules/currency/repository.go`, `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`, `docs/runtime-protocol-adapter.md`, `docs/bound-identity-route-policy-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0207`

The paired Simplified Chinese translation is `docs/currency-wallet-runtime-behavior-gate.zh-CN.md`. The English file is authoritative.

This document defines the currency wallet runtime behavior gate. It is a gate artifact. It does not add runtime behavior implementation, runtime handlers, startup wiring, protocol routes, Protobuf source, generated output, repository interface changes, PostgreSQL adapter changes, migration changes, dependencies, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, payment behavior, reservation behavior, settlement behavior, refund behavior, transfer behavior, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The currency wallet runtime behavior gate record is:

```yaml
currency_wallet_runtime_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0299
decision: ADR-0207
check_rule: runtime.currency_wallet_runtime_behavior_gate
source_postgresql_adapter_implementation_decision: ADR-0206
source_postgresql_adapter: runtime/internal/platform/persistence/postgres/currency_wallet_repository.go
source_repository_interface_decision: ADR-0204
repository_interface: runtime/internal/modules/currency.Repository
repository_interface_source: runtime/internal/modules/currency/repository.go
future_runtime_owner_candidate: runtime/internal/app
future_currency_application_package_candidate: runtime/internal/app/currency
future_runtime_service_source_candidate: runtime/internal/app/currency/service.go
future_runtime_service_test_candidate: runtime/internal/app/currency/service_test.go
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
owner_identity_source: validated_request_identity_player_id
actor_identity_source: validated_request_identity_or_server_operation
first_owner_kind: player
first_actor_kinds:
  - player
  - system
route_policy_requirement: request_token_required
service_application_owner: runtime/internal/app
repository_handoff: unit_of_work_currency_wallet_repository_factory
unit_of_work_handoff_required: true
runtime_behavior_gate_only: true
runtime_behavior_added: false
runtime_handlers_added: false
startup_wiring_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
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
direct_nakama_pitaya_api_compatibility_added: false
future_runtime_behavior_implementation_work_item: W-0300
future_runtime_behavior_implementation_direction: implement_currency_wallet_runtime_behavior
```

## 2. Purpose

`W-0298` implemented the PostgreSQL adapter for `runtime/internal/modules/currency.Repository`. The next useful boundary is not a protocol route or reward integration. The next useful boundary is the runtime behavior gate that defines how application code may later turn a validated player request or server-authoritative operation into currency wallet repository operations.

This gate records the future behavior shape before implementation:

- application ownership for the service;
- wallet owner identity derivation from validated request identity;
- actor derivation for player-initiated and server-authoritative operations;
- command/query posture for wallet creation, balance reads, transaction reads, grants, and spends;
- permission and route-policy posture;
- validation and conflict mapping expectations;
- unit-of-work and repository handoff;
- idempotency and redaction rules;
- fake-service test expectations;
- stop conditions that keep protocol, generated output, authentication/session changes, reward, inventory, purchase, catalog, event/audit, payment, and distributed runtime scope out of this slice.

Nakama motivates durable wallets, balances, grants, spends, and transaction history as a core economy capability. Pitaya motivates keeping route/session context, handlers, and persistence responsibilities separated. vibit adapts those references through explicit application-owned behavior and checks, not direct public API compatibility.

## 3. Ownership

Future runtime behavior is application-owned:

```yaml
future_runtime_owner_candidate: runtime/internal/app
future_currency_application_package_candidate: runtime/internal/app/currency
future_runtime_service_source_candidate: runtime/internal/app/currency/service.go
future_runtime_service_test_candidate: runtime/internal/app/currency/service_test.go
repository_interface_owner: runtime/internal/modules/currency
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
player_account_owner: runtime/internal/modules/player
```

Rules:

- Future service behavior may live under `runtime/internal/app/currency` or an equivalent application-owned package ratified by the implementation slice.
- The service may call `runtime/internal/modules/currency.Repository` only through application or unit-of-work dependencies.
- The service may check player account existence and account state only through existing application-owned player repository capabilities when the implementation slice authorizes that dependency handoff.
- The service must not import PostgreSQL adapter packages, SQL row types, migration packages, WebSocket transport packages, generated Protobuf packages, reward packages, inventory packages, purchase packages, payment provider SDKs, catalog packages, event/audit packages, SDK packages, or distributed runtime packages.
- The currency module remains the owner of storage-neutral wallet value types, normalizers, lifecycle vocabulary, transaction vocabulary, actor vocabulary, idempotency vocabulary, and repository error vocabulary.
- The PostgreSQL adapter remains persistence-only and must not derive request identity, route policy, public protocol errors, reward decisions, inventory decisions, purchase decisions, or payment decisions.
- Transport and protocol adapters must not own currency wallet permission, economy policy, or business behavior.

## 4. Request Identity And Owner Derivation

The first runtime behavior posture is player-wallet owned:

```yaml
first_owner_kind: player
owner_identity_source: validated_request_identity_player_id
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed_as_proof: false
client_supplied_wallet_id_allowed_for_lookup: false
server_wallet_lookup_by_owner_required: true
```

Rules:

- A future player-visible wallet operation must derive `currency.CurrencyWalletOwner{Kind: player, ID: identity.PlayerID}` from a validated `app.RequestIdentity`.
- `RequestIdentity.Status` must be `validated`.
- `RequestIdentity.ActorKind` must be `player`.
- `RequestIdentity.PlayerIDValidated` must be true.
- `RequestIdentity.PlayerID` must be non-empty and must match the actor identity when both are present.
- Metadata-only `player_id` from envelope/session metadata must never satisfy this gate.
- A persisted `session_id` alone must never become proof.
- Client payloads must not choose another owner id in the first posture.
- Client payloads should not choose arbitrary wallet ids for player-owned operations in the first posture; the service should resolve the wallet by server-derived owner when possible.

This gate does not change `RequestIdentity`, access-token validation, bound connection identity, durable runtime session validation, or WebSocket handshake behavior. It only records the identity requirements that future currency behavior must require before repository access.

## 5. Actor Derivation

Currency transactions require an actor in addition to wallet ownership:

```yaml
actor_identity_source: validated_request_identity_or_server_operation
first_actor_kinds:
  - player
  - system
operation_actor_kind_reserved: true
client_supplied_actor_id_allowed_as_proof: false
```

Rules:

- Player-initiated wallet reads and spends must derive actor identity from validated request identity.
- Server-authoritative grants may use a service-owned `system` actor when the implementation slice explicitly exposes that dependency and reason-code posture.
- `operation` actors remain reserved for later operations/admin behavior and are not enabled by this gate.
- Client payloads must not provide actor ids as proof.
- Actor ids, wallet ids, owner ids, idempotency keys, reason codes, external references, and transaction ids are not log-safe by default.

## 6. Future Runtime Behavior Shape

The future first implementation may expose an application service with these candidate operations:

```yaml
candidate_operations:
  - ensure_player_wallet
  - get_own_wallet
  - list_own_wallet_balances
  - grant_currency
  - spend_currency
  - list_own_wallet_transactions
```

Recommended first posture:

- `ensure_player_wallet` creates or finds the active wallet for the validated player through a caller-owned unit of work.
- `get_own_wallet` reads the wallet for the validated player.
- `list_own_wallet_balances` lists balances for the validated player's wallet with bounded pagination.
- `grant_currency` records a server-authoritative grant with an idempotency key, reason code, and optional metadata, without integrating reward, inventory, purchase, catalog, or payment behavior.
- `spend_currency` records a player-authorized spend with an idempotency key, positive amount input, and insufficient-balance conflict mapping, without integrating purchase, reservation, settlement, refund, transfer, or paid-currency behavior.
- `list_own_wallet_transactions` lists transactions for the validated player's wallet with bounded pagination and optional currency filter.

Rules:

- Runtime behavior must use server-derived owner and actor identity.
- Runtime behavior must validate currency code, positive amount, idempotency fields, reason code, external reference, metadata JSON, expected wallet version, expected balance version, list limit, and pagination cursor before repository calls when possible.
- Runtime behavior must not expose cross-owner wallet reads, arbitrary wallet id lookup, admin wallet inspection, player-to-player transfer, paid-currency purchase, reward inventory integration, catalog lookup, exchange rates, reserves, settlement, refunds, or payment provider behavior in the first implementation.
- Runtime behavior must not add public protocol routes or generated output unless a later protocol gate authorizes them.

## 7. Candidate Application Service Shape

The first implementation slice should define a small application-owned service boundary. Candidate inputs and outputs:

```yaml
candidate_request_fields:
  - request_identity
  - currency_code
  - amount
  - idempotency_key
  - idempotency_scope
  - reason_code
  - external_reference
  - metadata_json
  - expected_wallet_version
  - expected_balance_version
  - list_limit
  - after_currency_code
  - after_transaction_id
  - after_transaction_time
candidate_result_fields:
  - wallet
  - balance
  - balances
  - transaction
  - transactions
  - next_currency_code
  - next_transaction_id
  - next_transaction_time
  - public_error_code
```

Rules:

- The service should accept already-normalized application identity, not raw tokens, cookies, headers, WebSocket subprotocol values, or envelope proof carriers.
- The service should call currency module normalizers before repository handoff.
- The service should keep transaction metadata JSON, idempotency keys, wallet ids, owner ids, actor ids, reason codes, external references, and platform errors out of default errors and logs.
- The service should expose stable public error codes or classes for runtime handlers to map later.
- The service should not add route registration, Protobuf conversion, startup composition, reward composition, inventory composition, purchase composition, or command/query dispatch wiring in the gate slice.

## 8. Validation Rules

Future runtime behavior must enforce validation before persistence:

```yaml
validation_required:
  - validated_player_identity
  - server_derived_wallet_owner
  - active_wallet_required_for_mutation
  - currency_code_non_empty_length_bounded
  - amount_positive_for_grant_and_spend_requests
  - idempotency_key_non_empty_length_bounded
  - idempotency_scope_non_empty_length_bounded
  - reason_code_length_bounded
  - external_reference_length_bounded
  - metadata_json_top_level_object_when_present
  - expected_wallet_version_positive_when_present
  - expected_balance_version_positive_when_present
  - list_limit_bounded
  - pagination_cursor_bounded
```

Rules:

- Currency code and amount validation should reuse currency module normalization rules unless a future contract explicitly tightens protocol-visible syntax.
- Grant and spend requests should accept positive amount inputs; transaction deltas remain repository-owned.
- Metadata JSON is not log-safe and must remain copied or immutable across boundaries.
- Missing expected version behavior must be explicit in implementation tests.
- Invalid input must fail before repository mutation when possible.
- Repository unavailable errors must remain redacted.

## 9. Permission And Route Policy Posture

The first route-policy posture is conservative:

```yaml
route_policy_requirement: request_token_required
public_currency_wallet_routes_allowed: false
bound_connection_required_by_this_gate: false
session_validated_required_by_this_gate: false
bound_session_required_by_this_gate: false
```

Candidate permission families for later public contracts:

- read own wallet;
- list own wallet balances;
- list own wallet transactions;
- spend own wallet balance through a ratified service command;
- receive server-authoritative currency grants.

Rules:

- Currency wallet routes must be protected routes.
- The first posture should use the existing `request_token_required` route protection family unless a later route-policy ADR changes named routes.
- Public routes must not read or mutate currency wallets.
- Bound connection identity and durable session validation may remain available for future route families, but this gate does not require them or change ordinary protected route behavior.
- Metadata-only identity must fail closed.

## 10. Conflict And Error Mapping

Future runtime behavior must map currency repository errors into stable application classes:

```yaml
candidate_public_error_classes:
  - CURRENCY_WALLET_INVALID_REQUEST
  - CURRENCY_WALLET_NOT_FOUND
  - CURRENCY_WALLET_ALREADY_EXISTS
  - CURRENCY_WALLET_NOT_ACTIVE
  - CURRENCY_BALANCE_NOT_FOUND
  - CURRENCY_AMOUNT_INVALID
  - CURRENCY_INSUFFICIENT_BALANCE
  - CURRENCY_TRANSACTION_DUPLICATE
  - CURRENCY_TRANSACTION_CONFLICT
  - CURRENCY_WALLET_VERSION_MISMATCH
  - CURRENCY_BALANCE_VERSION_MISMATCH
  - CURRENCY_WALLET_FORBIDDEN
  - CURRENCY_WALLET_UNAVAILABLE
```

Rules:

- Not-found, owner mismatch, closed wallet, and suspended wallet cases must avoid cross-player existence leaks.
- Insufficient balance may be public as a stable conflict class, but balance internals, wallet ids, SQL details, driver errors, DSNs, credentials, token material, verifier digests, and route proof carriers must not leak.
- Duplicate idempotency may be public only as a stable duplicate or conflicting-duplicate class.
- Repository `storage_unavailable` errors must map to a retryable or unavailable class without exposing platform internals.
- Permission failure must happen before repository access when the request identity is not validated.
- Runtime behavior must not add authentication/token/session failure detail beyond existing application route-protection classes.

## 11. Unit-Of-Work And Repository Handoff

Future runtime behavior should use the existing application transaction boundary:

```yaml
repository_handoff: unit_of_work_currency_wallet_repository_factory
unit_of_work_handoff_required: true
caller_owns_transaction_control: true
postgres_adapter_transaction_control_allowed: false
```

Rules:

- The application service should obtain `currency.Repository` through a unit-of-work capability or equivalent app-owned dependency ratified by the implementation slice.
- Wallet creation and transaction mutations should execute inside the caller-supplied unit of work.
- The service may compose player account lookup and currency repository calls in one unit of work only after the implementation slice explicitly records that composition.
- The PostgreSQL adapter must not start, commit, or roll back transactions itself.
- Runtime behavior must not bypass the repository interface by importing the PostgreSQL adapter directly.

## 12. Test Expectations

Future runtime behavior implementation should add focused fake-service tests before implementation:

```yaml
future_test_expectations:
  - rejects_missing_or_unvalidated_request_identity_before_unit_of_work
  - derives_wallet_owner_from_validated_player_identity
  - rejects_metadata_only_player_id_as_proof
  - ensures_player_wallet_through_repository_handoff
  - lists_balances_for_own_wallet_only
  - records_server_authoritative_grant_with_idempotency
  - records_player_spend_with_insufficient_balance_mapping
  - lists_transactions_for_own_wallet_only
  - maps_repository_conflicts_to_public_errors
  - redacts_wallet_owner_actor_idempotency_and_metadata_details
  - keeps_protocol_routes_and_generated_output_absent
```

Default verification for this gate remains static repository checks. Future implementation verification should include focused Go tests under `runtime/internal/app/currency`.

## 13. Stop Conditions

Stop and create a new bounded work item before adding any of the following:

- runtime behavior implementation;
- runtime handlers or route dispatch wiring;
- startup composition;
- protocol routes;
- Protobuf source;
- generated output;
- repository interface changes;
- PostgreSQL adapter changes;
- migration changes;
- dependency additions;
- authentication/session behavior changes;
- reward integration;
- inventory integration;
- purchase behavior;
- catalog tables;
- event/audit tables;
- payment provider behavior;
- reservation, settlement, refund, transfer, or paid-currency behavior;
- hosted deployment;
- SDK publication;
- release artifacts;
- public announcements or paid promotion;
- Pitaya-style distributed runtime;
- direct Nakama/Pitaya API compatibility.

## 14. Verification

Required repository verification for this gate:

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.currency_wallet_runtime_behavior_gate
node tools/vibit check change define-currency-wallet-runtime-behavior-gate --json
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

The accepted pre-existing warning is `runtime.identity_boundary` on `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
