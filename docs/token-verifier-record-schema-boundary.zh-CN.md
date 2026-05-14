# Token Verifier Record Schema Boundary

Status: Draft v0.1
Last updated: 2026-05-14
Scope: 第一版 opaque access-token 姿态所需的 token verifier record schema boundary
Depends on: `docs/credential-token-session-schema-gates.md`, `docs/token-lifecycle-storage-implications.md`
Canonical decision: `ADR-0033`

配套英文源文档是 `docs/token-verifier-record-schema-boundary.md`。英文文件是权威版本。

## 1. Purpose

本标准 ratify vibit 在实现已选第一版 access-token posture 前所需的 token verifier record boundary：

```yaml
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
default_durable_target: PostgreSQL
schema_boundary_status: migration_source_added
```

本文定义 token verifier record semantics。W-0078 已为该 schema 添加 SQL migration source。它不添加 repository interfaces、PostgreSQL adapters、runtime token validation、token issuance、logout behavior、refresh behavior、cleanup jobs、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior。

## 2. Required Reading

应与以下文件一起阅读：

- `docs/credential-token-session-schema-gates.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/first-token-format-proof-carrier-posture.md`
- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/credential-record-schema-boundary.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/player-identity-session-boundary.md`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `ADR-0026`
- `ADR-0027`
- `ADR-0029`
- `ADR-0030`
- `ADR-0031`
- `ADR-0032`
- `ADR-0033`

Reference reading：

- Nakama authentication concepts：`https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts：`https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya features and session vocabulary：`https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya API and handler vocabulary：`https://pitaya.readthedocs.io/en/latest/API.html`

Nakama 和 Pitaya 仍然是 capability coverage、session vocabulary、token lifecycle pressure 和 handler context vocabulary 的 reference。它们不支配 vibit 的 token schema、public API、session model、transport behavior 或 generated boundaries。

## 3. Boundary Status

Token verifier record boundary 只作为 schema boundary 被 ratify：

```yaml
token_verifier_record_schema_boundary:
  status: migration_source_added
  default_durable_target: PostgreSQL
  future_logical_table: authentication_access_tokens
  owner: runtime.authentication
  migration_source: runtime/migrations/postgres/000004_create_authentication_access_tokens.sql
  migration_source_added: true
  repository_interface_added: false
  postgres_adapter_added: false
  runtime_validation_added: false
  token_issuance_added: false
  logout_added: false
  cleanup_added: false
  authentication_implemented: false
```

Logical table name 和 SQL migration source 已 ratify，目的是让 agents 在 repository 与 adapter work 开始前有稳定目标。只有应用 `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql` 后，该 table 才会在数据库中存在。

## 4. Ownership

未来 token verifier record 属于 runtime authentication boundary：

```text
runtime.authentication
```

未来 implementation ownership：

- Semantic token contract owner：`contracts/runtime/authentication/`
- Runtime validation 和 identity handoff owner：`runtime/internal/app`
- Future repository interface owner：后续 ratified authentication repository boundary
- Future PostgreSQL adapter owner：`runtime/internal/platform/persistence/postgres/`
- Future migration source owner：`runtime/migrations/postgres/`
- Future cleanup owner：后续 ratified authentication maintenance 或 token storage boundary

Player module 拥有 player account lifecycle state。它不拥有 access-token verifier records。

Player account lifecycle tables 仍然是：

```text
player_accounts
player_account_events
```

它们必须保持 credential-free、token-free、provider-subject-free、runtime-session-free、WebSocket-state-free 和 request-validation-free。

## 5. Logical Record

未来 token verifier record 表示一个 opaque access token 的 server-side verifier。

未来 logical record：

```yaml
record: authentication_access_token
fields:
  token_record_id: log_safe_identifier
  token_kind: access_token
  token_status: active | expired | revoked
  actor_kind: player
  player_id: player_account_reference
  credential_record_id: credential_record_reference
  token_lookup_digest: secret_adjacent_index_material
  token_verifier_digest: secret_verifier_material
  verifier_algorithm: versioned_non_plaintext_verifier
  verifier_version: integer
  verifier_key_id: secret_key_reference_not_secret_value
  audience: route_or_gameplay_audience_catalog
  issued_at: timestamp
  expires_at: timestamp
  revoked_at: nullable_timestamp
  revoked_reason: nullable_catalog_value
  replaced_by_token_record_id: nullable_log_safe_identifier
  last_validated_at: nullable_timestamp
  last_failed_validation_at: nullable_timestamp
  cleanup_after: nullable_timestamp
  created_at: timestamp
  updated_at: timestamp
```

未来 SQL migration 可以细化 exact PostgreSQL types、constraints 和 index names，但必须保留这些 semantics，除非后续 ADR supersede 本 boundary。

第一版 token verifier record 不 ratify generic `metadata` column。新增 metadata 需要未来 schema decision，因为任意 JSON fields 很容易被 agents 用来隐藏 raw tokens、claims、transport state、provider payloads、device fingerprints 或 private request details。

## 6. Identifier Rules

`token_record_id` 是唯一适合普通 logs、change specs、ADRs、tests 和 conversation logs 的 token identifier。

Rules：

- `token_record_id` 必须 globally unique。
- `token_record_id` 不得从 raw token material 派生。
- `token_record_id` 不得暴露 `player_id`、`credential_record_id`、session metadata、route name、IP address、user agent、provider subject 或 token issuance time。
- `token_lookup_digest` 不是 log-safe。
- `token_verifier_digest` 不是 log-safe。
- Raw access-token text 永远不是 log-safe。
- `player_id` 不是 secret proof，不能单独满足 authentication。
- `credential_record_id` 作为 credential identifier 是 log-safe，但不得被当成 token proof。

## 7. Verifier Semantics

Raw access token 永远不存储。

第一版 verifier posture：

```yaml
raw_token_storage: forbidden
token_lookup_digest_required: true
token_verifier_digest_required: true
verifier_algorithm_versioned: true
plaintext_comparison: forbidden
minimum_entropy_bits: 256
jwt_or_claim_parsing_required: false
signing_dependency_required: false
redis_like_store_required: false
```

该 token 预期是 high entropy。因此本 boundary 不 ratify JWT、JWK、OAuth、OIDC、signing、key-management、Redis-like token/session stores、provider SDKs、bcrypt 或 Argon2 dependencies。

未来 implementation 仍必须在 code 存在前定义 exact verifier algorithm。如果后续 implementation work item ratify exact algorithm、pepper/secret configuration boundary、digest format 和 comparison behavior，implementation 可以使用 standard-library cryptographic primitives。

## 8. Lifecycle States

Ratified states：

| State | Meaning | Request validation allowed |
| --- | --- | --- |
| `active` | 如果未过期且 linked actor 与 credential gates 通过，token verifier 可以 authenticate request。 | Yes |
| `expired` | Token verifier 超过 validity window。 | No |
| `revoked` | Token verifier 在过期前或过期后被明确 invalidated。 | No |

State rules：

- `expired` 和 `revoked` 对第一版姿态都是 terminal。
- 即使 `token_status` 尚未物化为 `expired`，`active` token 在 `expires_at` 之后也无效。
- 未来 validator 可以把 expiration 当成 computed state，而不是急切更新每一行 expired row。
- Revocation 必须在 production-sensitive domain dispatch 前生效。
- Disabled 或 deleted player account 必须阻断 token validation，即使 token record 仍是 active。
- 如果 validator 在授权实现中检查 credential state，那么 disabled、revoked 或 replaced credential 必须阻断 token validation。

## 9. Actor And Credential Relationship

第一版 token verifier relationship：

```yaml
one_token_record_represents_one_access_token: true
one_access_token_belongs_to_one_actor: true
first_actor_kind: player
player_id_mutable: false
credential_record_id_required_for_first_posture: true
token_can_move_between_players: false
token_can_move_between_credentials: false
refresh_token_storage: forbidden_for_first_posture
runtime_session_binding: deferred
websocket_connection_binding: deferred
```

第一版姿态要求 `credential_record_id` linkage，以便后续 implementation 在 successful login rotation 时可以 revoke same credential 的 previous active tokens。该 linkage 不是 token proof，也不会单独授权 credential lookup、token validation 或 logout behavior。

未来 service 或 admin tokens 需要单独的 actor-kind decision，因为它们的 lifetime、audience、permission、storage 和 audit rules 与 player access tokens 不同。

## 10. Expiration

第一版 access-token TTL 仍然是：

```yaml
access_token_ttl: 1h
```

未来 token verifier records 必须保留：

```yaml
issued_at_required: true
expires_at_required: true
expires_after_issued_at: true
expired_proof_failure_class: required
expired_record_retention: required
```

Expired token 必须以 `expired_proof` 失败，而不是 missing proof 或 malformed proof。Expired records 可以为 replay analysis、logout idempotency、abuse investigation 和 audit correlation 临时保留。

## 11. Revocation, Logout, And Rotation

第一版 logout scope 仍然是：

```yaml
logout_scope_first_posture: presented_access_token
```

Revocation semantics：

- Logout 只 revoke presented access token。
- Admin revocation 仍 deferred 到未来 permission surface。
- Logout-all-devices 仍 deferred。
- Disabled 或 deleted account 的 forced revocation 仍 deferred 到 account policy 和 audit work。
- Revoked tokens 必须与 malformed、missing、invalid、expired proof 区分失败。
- Revocation reason 必须作为 non-secret catalog value 保留。

Rotation semantics：

```yaml
rotation_on_successful_login: required_when_implementation_exists
previous_active_token_for_same_credential: revoke_when_repository_supports_it
replaced_by_token_record_id: optional_lineage_field
automatic_background_rotation: deferred
refresh_token_rotation: deferred
```

Successful login 应 issue new token，并在 schema、repository、implementation gates 授权后 revoke same credential 的 previous active access tokens。这不需要 runtime session persistence。

## 12. Retention And Cleanup

Production token storage 启用前需要 cleanup，但本标准不添加 cleanup job。

第一版 retention posture：

```yaml
active_records_retained_until_expiration_or_revocation: true
expired_record_default_retention_recommendation: 7d
revoked_record_default_retention_recommendation: 7d
cleanup_target: expired_and_revoked_token_verifier_records
cleanup_trigger_first_posture: explicit_maintenance_command_or_scheduled_runtime_job_deferred
cleanup_owner: future_authentication_or_token_storage_boundary
```

未来 migration 可以使用 `cleanup_after` 让 retention 显式化。如果不使用，则 repository 或 cleanup boundary 必须定义 records 如何进入 deletion eligibility。

Cleanup 不得删除 active tokens。Cleanup 不得隐藏在 request validation 中。生产使用前，cleanup 必须 idempotent、concurrency-aware 且 auditable。

## 13. Replay-Sensitive Failure Classes

未来 validation 必须保留这些 failure classes：

```yaml
failure_classes:
  - missing_proof
  - malformed_proof
  - unsupported_proof
  - invalid_proof
  - expired_proof
  - revoked_proof
  - actor_disabled
  - validator_unavailable
```

Replay-sensitive rules：

- Malformed token 不得触发会泄露 valid token 是否存在的 lookup behavior。
- Invalid token 不得暴露 lookup digest 是否命中但 verifier digest 失败。
- Expired 和 revoked tokens 可以在 proof format 被接受后产生不同的 stable failure classes。
- Failure responses 不得在 ratified error surface 之外泄露 player account、credential record 或 token record 是否存在。
- Rate limiting 和 abuse controls 仍是 future work，但 schema 和 error design 不得阻碍它们。

## 14. Uniqueness And Index Rules

未来 migration work 必须保留这些 uniqueness semantics：

```yaml
unique:
  - token_record_id
  - token_lookup_digest
indexes:
  - player_id
  - credential_record_id
  - token_status
  - expires_at
  - cleanup_after
foreign_key_like_relationships:
  - player_id references player account lifecycle identity
  - credential_record_id references authentication_device_credentials
```

未来 SQL source 可以把 `player_id` 和 `credential_record_id` 实现为真实 foreign keys，但前提是这不破坏 module ownership、migration order 或 test isolation。每个 relationship 是 database foreign key 还是 application-enforced reference，必须在 migration work item 中明确。

## 15. Redaction

以下内容禁止出现在 logs、errors、traces、tests、fixtures、ADRs、change specs 和 conversation logs 中：

- Raw access tokens。
- Token lookup digest。
- Token verifier digest。
- Server-side verifier secrets 或 peppers。
- Authorization header contents。
- Cookie contents。
- WebSocket subprotocol token carriers。
- URL query token carriers。
- 包含 token text 的完整 request proof payloads。
- Credential lookup digest。
- Credential verifier digest。
- Provider secrets。

可谨慎使用：

- `token_record_id`
- `credential_record_id`
- `player_id`
- lifecycle state names
- non-secret reason catalog values
- stable failure class names

Test fixtures 必须使用 synthetic values，不得包含真实 tokens、credentials、device identifiers、provider payloads、authorization headers、cookies 或复制的 production data。

## 16. Forbidden Shortcuts

Agents must not：

- 在 `W-0078` 或后续明确 migration work item 之外添加或修改 `authentication_access_tokens` SQL migration source。
- 仅凭本标准添加 token repository interfaces。
- 仅凭本标准添加 PostgreSQL token adapters。
- 仅凭本标准添加 runtime token issuance、validation、logout、refresh 或 cleanup。
- 在 `player_accounts` 或 `player_account_events` 中存储 tokens。
- 在 credential records 中存储 token verifier records。
- 在 token verifier records 中存储 raw access tokens、refresh tokens、runtime sessions、WebSocket connection state 或 request validation results。
- 把当前 Protobuf `Session` fields 当成 access-token proof。
- 把 token proof 放入 WebSocket handshake headers、cookies、subprotocols、query parameters、routes、request IDs、player IDs、session IDs 或 connection IDs。
- 仅凭本标准添加 JWT、signing、OAuth、OIDC、provider SDK、key-management、Redis-like、password hashing 或其他 major authentication dependencies。
- 复制 Nakama 或 Pitaya public API shapes。

## 17. Required Future Gates

在 token storage 可以实现前，未来工作必须完成这些 gates：

```yaml
credential_record_schema_boundary: completed_by_W_0074
token_verifier_record_schema_boundary: completed_by_W_0075
authentication_schema_migration_queue: required_before_migration
token_verifier_migration_source: completed_by_W_0078
authentication_repository_interface: separate_future_work
token_postgres_adapter: separate_future_work
verifier_algorithm_and_secret_configuration: separate_future_work
redaction_tests: separate_future_or_implementation_work
runtime_token_validation_wiring: separate_future_implementation_milestone
logout_behavior: separate_future_implementation_milestone
cleanup_behavior: separate_future_maintenance_milestone
```

Migration gate 必须证明 SQL source 只创建已 ratified 的 token verifier structures，并且不修改 player account lifecycle tables。

## 18. Reference Alignment

### Nakama

Nakama 说明成熟 game backend 通常会提供 access tokens 或 session tokens、expiration、refresh、revocation、logout 和多个 authentication entry points。

vibit 只借鉴 lifecycle pressure，不复制 public API shape。第一版 vibit posture 使用 opaque high-entropy access tokens、server-side verifier storage、无 refresh token、login-command issuance 和 explicit request proof payloads。

### Pitaya

Pitaya 展示了 connection acceptors、handler context 中的 sessions、ID binding、frontend/backend session 差异和 route-aware message handling 等有用 session 与 handler vocabulary。

vibit 将这些概念与 token verifier storage 分离。Token verifier records 是 persistent authentication proof verifier records，不是 transport sessions、handler context state 或 WebSocket connection binding records。

## 19. Verification

本标准的默认 verification：

```bash
node tools/vibit inspect next --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check memory --json
node tools/vibit check change define-token-verifier-record-schema-boundary --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

对于 source-only migration，live PostgreSQL verification 是可选的，并应由 migration work item 记录。默认不强制要求它，因为 repository 与 adapter behavior 仍然 deferred。

由于没有 Go runtime behavior changes，因此不需要 Go tests。

## 20. Follow-Up

Next work：

```text
W-0079 Add authentication migration static checks
```

下一个 gate 应先强化已 ratified 的 credential 与 token verifier migration sources 的本地 checks，然后再定义 repository interfaces 或 adapters。
