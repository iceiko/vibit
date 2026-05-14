# Credential Record Schema Boundary

Status: Draft v0.1
Last updated: 2026-05-14
Scope: 第一版 `device_credential_login` 姿态所需的 credential record schema boundary
Depends on: `docs/credential-token-session-schema-gates.md`
Canonical decision: `ADR-0032`

配套英文源文档是 `docs/credential-record-schema-boundary.md`。英文文件是权威版本。

## 1. Purpose

本标准 ratify vibit 在实现已选第一版 login method 前所需的 credential record boundary：

```yaml
login_method: device_credential_login
credential_kind: high_entropy_installation_credential
default_durable_target: PostgreSQL
schema_boundary_status: migration_source_added
credential_migration_source: runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
```

本文定义未来 record semantics。它不添加 SQL migration source、tables、repository interfaces、PostgreSQL adapters、runtime lookup、login handlers、token behavior、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior。

## 2. Required Reading

应与以下文件一起阅读：

- `docs/credential-token-session-schema-gates.md`
- `docs/first-login-method-set.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/player-identity-session-boundary.md`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `ADR-0025`
- `ADR-0029`
- `ADR-0030`
- `ADR-0031`
- `ADR-0032`

Reference reading：

- Nakama authentication concepts：`https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts：`https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya features and session vocabulary：`https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya API and handler vocabulary：`https://pitaya.readthedocs.io/en/latest/API.html`

Nakama 和 Pitaya 仍然是 capability coverage 与 vocabulary reference。它们不支配 vibit 的 credential schema、public API、session model 或 generated boundaries。

## 3. Boundary Status

Credential record boundary 只作为 schema boundary 被 ratify：

```yaml
credential_record_schema_boundary:
  status: migration_source_added
  default_durable_target: PostgreSQL
  future_logical_table: authentication_device_credentials
  owner: runtime.authentication
  migration_source: runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
  migration_source_added: true
  repository_interface_added: false
  postgres_adapter_added: false
  runtime_lookup_added: false
  authentication_implemented: false
```

Logical table name 已 ratify，并且 SQL source 现在存在于 `runtime/migrations/postgres/` 下。该 migration source 不授权 repositories、adapters、runtime lookup、handlers、routes、generated authentication output、Protobuf messages、WebSocket behavior、authentication dependencies 或 production authentication behavior。

## 4. Ownership

未来 credential record 属于 runtime authentication boundary：

```text
runtime.authentication
```

未来 implementation ownership：

- Semantic login contract owner：`contracts/runtime/authentication/`
- Runtime validation 和 identity handoff owner：`runtime/internal/app`
- Future repository interface owner：后续 ratified authentication repository boundary
- Future PostgreSQL adapter owner：`runtime/internal/platform/persistence/postgres/`
- Future migration source owner：`runtime/migrations/postgres/`

Player module 拥有 player account lifecycle state。它不拥有 credential verifier records。

Player account lifecycle tables 仍然是：

```text
player_accounts
player_account_events
```

它们必须保持 credential-free、token-free、provider-subject-free、runtime-session-free、WebSocket-state-free 和 request-validation-free。

## 5. Logical Record

未来 credential record 表示一个 durable high-entropy installation credential，并绑定到一个 player account。

未来 logical record：

```yaml
record: authentication_device_credential
fields:
  credential_record_id: log_safe_identifier
  player_id: player_account_reference
  credential_kind: device_credential_login
  credential_status: active | disabled | revoked | replaced
  credential_lookup_digest: secret_adjacent_index_material
  credential_verifier_digest: secret_verifier_material
  verifier_algorithm: versioned_non_plaintext_verifier
  verifier_version: integer
  verifier_key_id: secret_key_reference_not_secret_value
  client_instance_id_digest: optional_privacy_sensitive_correlation_material
  created_at: timestamp
  updated_at: timestamp
  last_verified_at: nullable_timestamp
  disabled_at: nullable_timestamp
  disabled_reason: nullable_catalog_value
  revoked_at: nullable_timestamp
  revoked_reason: nullable_catalog_value
  replaced_at: nullable_timestamp
  replaced_by_credential_record_id: nullable_log_safe_identifier
```

未来 SQL migration 可以细化 exact PostgreSQL types、constraints 和 index names，但必须保留这些 semantics，除非后续 ADR supersede 本 boundary。

第一版 credential record 不 ratify generic `metadata` column。新增 metadata 需要未来 schema decision，因为任意 JSON fields 很容易被 agents 用来隐藏 credentials、provider payloads、tokens 或 transport state。

## 6. Identifier Rules

`credential_record_id` 是唯一适合普通 logs、change specs、ADRs、tests 和 conversation logs 的 credential identifier。

Rules：

- `credential_record_id` 必须 globally unique。
- `credential_record_id` 不得从 raw credential material 派生。
- `credential_record_id` 不得暴露 device model、platform account、provider subject、IP address、user agent、player name 或 player display name。
- `credential_lookup_digest` 不是 log-safe。
- `credential_verifier_digest` 不是 log-safe。
- `client_instance_id_digest` 默认不是 log-safe。
- `player_id` 不是 secret proof，不能单独满足 authentication。

## 7. Verifier Semantics

Raw credential proof 永远不存储。

第一版 verifier posture：

```yaml
raw_credential_storage: forbidden
raw_os_device_id_as_credential: forbidden
credential_lookup_digest_required: true
credential_verifier_digest_required: true
verifier_algorithm_versioned: true
plaintext_comparison: forbidden
password_hashing_dependency_required: false
external_provider_dependency_required: false
```

该 credential 预期是 high entropy。因此本 boundary 不 ratify bcrypt、Argon2、OAuth、OIDC、JWT、provider SDKs 或 key-management dependencies。

未来 implementation 仍必须在 code 存在前定义 exact verifier algorithm。如果后续 implementation work item ratify exact algorithm 和 secret configuration boundary，implementation 可以使用 standard-library cryptographic primitives。如果未来设计需要 password-like 或 low-entropy credential，则必须使用单独的 credential boundary 和 dependency adoption record。

## 8. Lifecycle States

Ratified states：

| State | Meaning | Login allowed |
| --- | --- | --- |
| `active` | 如果 linked player account 也是 active 且 token issuance gates 通过，credential 可以 authenticate。 | Yes |
| `disabled` | Credential 被 policy 或 operations 临时阻断。 | No |
| `revoked` | Credential 永久失效。 | No |
| `replaced` | Credential 已被 rotated 或 replacement credential supersede。 | No |

State rules：

- `revoked` 和 `replaced` 是 terminal。
- `disabled` 只能通过未来明确授权的 administrative 或 recovery flow 变回 `active`。
- 即使 credential 是 `active`，linked player account disabled 或 deleted 时也必须阻断 login。
- Credential status 不得覆盖 player account lifecycle state。
- Disabled 或 deleted player account 不得因为 credential login 被静默 re-enable。

## 9. Player Relationship

第一版 credential record relationship：

```yaml
one_credential_record_belongs_to_one_player: true
credential_player_id_mutable: false
credential_can_move_between_players: false
one_player_active_device_credentials_first_posture: at_most_one
historical_records_per_player_allowed: true
multi_device_linking: deferred
account_recovery: deferred
account_merge: deferred
```

第一版姿态允许同一 player 有历史 credential records，因为 rotation 和 revocation 需要 audit trail。它不授权一个 player 同时拥有多个 active device credentials。Multi-device linking、account recovery 和 account merge 需要后续 account-linking decisions。

当 first login 创建 player account 时，player account creation 和 credential record creation 必须在同一个 application-owned unit of work 中保持 atomic。当 login authenticate existing account 时，credential verification 和 token verifier creation 必须保持一致的 failure behavior，且不能暴露 credential、player account 或 token record 是否存在。

## 10. Uniqueness And Index Rules

未来 migration work 必须保留这些 uniqueness semantics：

```yaml
unique:
  - credential_record_id
  - credential_lookup_digest
conditional_unique:
  - at_most_one_active_device_credential_per_player_for_first_posture
foreign_key_like_relationship:
  - player_id references player account lifecycle identity
```

未来 SQL source 可以把 `player_id` 实现为真实 foreign key，但前提是这不破坏 module ownership、migration order 或 test isolation。该 relationship 是 database foreign key 还是 application-enforced reference，必须在 migration work item 中明确。

## 11. Rotation And Replacement

Credential rotation 会创建新的 credential record，并把 previous active credential 标记为 `replaced`。

Rules：

- 不要为了 rotation 原地覆盖 verifier material。
- 通过 `replaced_by_credential_record_id` 保留 lineage。
- 如果 rotation 改变 old 和 new records，必须在一个 unit of work 中完成。
- 如果 token verifier records 已存在，rotation 对 active access tokens 的影响必须遵循 token verifier schema boundary。
- Presented credential revocation 不得要求 raw credential storage。

第一版 schema boundary 为 rotation 做准备。它不实现 rotation commands 或 runtime behavior。

## 12. Redaction

以下内容禁止出现在 logs、errors、traces、tests、fixtures、ADRs、change specs 和 conversation logs 中：

- Raw credential proof。
- 作为 proof 的 raw operating-system device ID。
- Credential lookup digest。
- Credential verifier digest。
- Server-side verifier secrets 或 peppers。
- 完整 client instance identifiers。
- Raw access tokens。
- Token verifier hashes。
- Provider secrets。

可谨慎使用：

- `credential_record_id`
- `player_id`
- lifecycle state names
- non-secret reason catalog values

Test fixtures 必须使用 synthetic values，不得包含真实 device identifiers、真实 credentials、真实 tokens 或复制的 production data。

## 13. Forbidden Shortcuts

Agents 不得：

- 仅凭本标准添加 `authentication_device_credentials` SQL migration source。
- 仅凭本标准添加 credential repository interfaces。
- 仅凭本标准添加 PostgreSQL credential adapters。
- 仅凭本标准添加 runtime credential lookup 或 login handlers。
- 在 `player_accounts` 或 `player_account_events` 中存储 credentials。
- 在 credential records 中存储 provider subjects。
- 在 credential records 中存储 access tokens、refresh tokens、runtime sessions、WebSocket connection state 或 request validation results。
- 把当前 Protobuf `Session` fields 当作 credential proof。
- 把 credential proof 放入 WebSocket handshake headers、cookies、subprotocols、query parameters、routes、request IDs、player IDs、session IDs 或 connection IDs。
- 仅凭本标准添加 password hashing、OAuth、OIDC、provider SDK、JWT、key-management、Redis-like 或其他 major authentication dependencies。
- 复制 Nakama 或 Pitaya public API shapes。

## 14. Required Future Gates

在实现 credential storage 前，未来 work 必须完成这些 gates：

```yaml
credential_record_schema_boundary: completed_by_W_0074
token_verifier_record_schema_boundary: required_before_authentication_implementation
authentication_schema_migration_queue: required_before_migration
credential_migration_source: completed_by_W_0077
credential_migration_source_path: runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
credential_repository_interface: separate_future_work
credential_postgres_adapter: separate_future_work
redaction_tests: separate_future_or_implementation_work
runtime_authentication_wiring: separate_future_implementation_milestone
```

Migration gate 必须证明 SQL source 只创建已 ratify 的 credential structures，并且不修改 player account lifecycle tables。

## 15. Reference Alignment

### Nakama

Nakama 展示了 mature game backend 可以支持 device-style authentication、automatic account creation posture、session tokens、refresh、logout 和 multiple linked authentication methods。

vibit 第一版 login method 只适配 low-friction device-login capability。它拒绝把 raw public device identifiers 当作 credential proof，推迟 multiple linked authentication methods，并把 session-token 与 refresh-token behavior 放在 credential record 之外。

### Pitaya

Pitaya 展示了有用的 session 与 handler vocabulary：connection acceptors、handler context 中可访问的 sessions、ID binding、frontend/backend session differences 和 route-aware message handling。

vibit 保持这些概念与 credential storage 分离。Credential records 是 persistent authentication verifier records，不是 transport sessions、handler context state 或 connection binding records。

## 16. Verification

本标准的默认 verification：

```bash
node tools/vibit inspect next --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check memory --json
node tools/vibit check change define-credential-record-schema-boundary --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

本 work 不添加 migration source、tables、repositories、adapters 或 runtime behavior，因此不需要 live PostgreSQL verification。

没有 Go runtime behavior 变化，因此不需要 Go tests。

## 17. Follow-Up

Next work：

```text
W-0075 Define token verifier record schema boundary
```

Token verifier record boundary 必须在 migration planning 开始前定义 opaque access-token verifier semantics、statuses、expiration、revocation、credential-token linkage、retention、cleanup 和 replay-sensitive failure classes。
