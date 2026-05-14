# Credential、Token 与 Session Schema Gate

状态：Draft v0.1
最后更新：2026-05-14
范围：Credential record、external identity link、token verifier record、runtime session record、audit persistence 与 player account lifecycle separation 的未来 schema gate
依赖：`docs/token-lifecycle-storage-implications.md`、`docs/authentication-contract-error-permission-surfaces.md`
权威决策：`ADR-0029`

对应英文原文是 `docs/credential-token-session-schema-gates.md`。英文文件是权威版本。

## 1. 目的

本标准定义 vibit 在实现第一版 authentication 姿态前必须具备的 schema gate：

```yaml
login_method: device_credential_login
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_device_credential_login
default_durable_target: PostgreSQL
```

本文档只是 gate 定义，不是 database schema。

它不添加 tables、migrations、repository interfaces、PostgreSQL adapters、runtime lookup code、token behavior、credential behavior、session persistence、audit persistence、Protobuf fields、WebSocket handshake behavior、runtime player handlers 或 WebSocket routes。

## 2. Schema Gate 规则

schema gate 是一种必需的规划工件。未来任何 change 要为安全敏感能力添加持久化 schema 前，必须先满足对应 gate。

每一个未来的 persistent authentication concern 都必须按以下顺序通过 gate：

```yaml
schema_ratification: required_before_migration
migration_source: separate_future_change
repository_interface: separate_future_change
postgres_adapter: separate_future_change
live_verification: separate_future_change_or_explicit_deferral
runtime_wiring: separate_future_implementation_milestone
```

规则：

- schema gate 可以命名 required record semantics。
- schema gate 可以命名 forbidden shortcuts。
- schema gate 可以命名 required future decisions。
- schema gate 不得创建 migration。
- schema gate 不得暗示 implementation 已授权。
- schema gate 必须保持当前 player account lifecycle tables 不变。

## 3. Gate 矩阵

W-0071 的 gate 状态是：

```yaml
credential_record_schema_gate:
  required_for_first_posture: true
  status: ratified_no_schema_added
  boundary: docs/credential-record-schema-boundary.md
  decision: ADR-0032
token_verifier_record_schema_gate:
  required_for_first_posture: true
  status: ratified_no_schema_added
  boundary: docs/token-verifier-record-schema-boundary.md
  decision: ADR-0033
external_identity_link_schema_gate:
  required_for_first_posture: false
  status: deferred_no_schema_added
runtime_session_record_schema_gate:
  required_for_first_posture: false
  status: deferred_no_schema_added
audit_persistence_schema_gate:
  required_before_durable_authentication_audit: true
  status: defined_no_schema_added
player_account_lifecycle_schema:
  status: preserved
  credential_columns_added: false
  token_columns_added: false
  external_identity_columns_added: false
  session_columns_added: false
```

在 required credential 与 token verifier schema gates 被转化为明确的 schema ratification work，并继续完成 migrations、repositories、adapters、tests 和 runtime wiring 之前，第一版实现不得开始。

W-0074 已经 ratify credential record schema boundary，但没有添加 schema：

```yaml
credential_record_schema_boundary:
  status: migration_source_added
  standard: docs/credential-record-schema-boundary.md
  decision: ADR-0032
  future_logical_table: authentication_device_credentials
  migration_source: runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
  migration_added_now: true
  repository_added_now: false
  runtime_lookup_added_now: false
```

W-0075 已经 ratify token verifier record schema boundary，W-0078 已添加它的 migration source：

```yaml
token_verifier_record_schema_boundary:
  status: migration_source_added
  standard: docs/token-verifier-record-schema-boundary.md
  decision: ADR-0033
  future_logical_table: authentication_access_tokens
  migration_source: runtime/migrations/postgres/000004_create_authentication_access_tokens.sql
  migration_added_now: true
  repository_added_now: false
  runtime_validation_added_now: false
  token_issuance_added_now: false
  logout_added_now: false
  cleanup_added_now: false
```

W-0077 与 W-0078 已添加 credential 与 token verifier migration sources。Authentication migration static checks 仍然是在添加 repository interfaces、adapters 或 runtime authentication behavior 前的下一个必需 step。

## 4. Credential Record Gate

实现 `device_credential_login` 前必须有 credential records。

Gate 状态：

```yaml
required_for_first_posture: true
default_durable_target: PostgreSQL
schema_added_now: false
migration_added_now: false
repository_added_now: false
runtime_lookup_added_now: false
```

未来 credential schema ratification 必须定义：

- Credential record owner。
- Record lifecycle states。
- Account creation 与 account lookup 的关系。
- 一个 player 是否可以有多个 device credential records。
- 一个 credential 是否可以在 players 之间迁移。
- 哪个 identifier 可以安全写日志。
- 存储哪种 secret verifier。
- 禁止存储哪些 raw credential material。
- 哪些 fields 是 unique。
- 哪些 fields 是 mutable。
- Rotation、replacement、revocation 与 disabled-credential behavior。
- Disabled 或 deleted player account state 如何阻止 login。
- 哪些操作必须和 player account lifecycle changes 保持原子性。
- Logs、errors、traces、tests、fixtures 和 conversation logs 的 redaction rules。
- Abuse-control 与 retryability classes。

第一版姿态的必需约束：

```yaml
credential_kind: device_credential_login
raw_os_device_id_as_credential: forbidden
raw_credential_storage: forbidden
credential_verifier_storage: required
credential_record_id_log_safe: required
player_account_lifecycle_tables_store_credentials: forbidden
```

本 gate 不选择 password model、OAuth provider、OIDC issuer、platform identity provider、provider SDK、key-management dependency、password-hashing dependency 或 external identity linking behavior。

## 5. Token Verifier Record Gate

实现 opaque access-token validation 前必须有 token verifier records。

Gate 状态：

```yaml
required_for_first_posture: true
default_durable_target: PostgreSQL
schema_added_now: false
migration_added_now: false
repository_added_now: false
runtime_validation_added_now: false
redis_like_store_selected: false
```

未来 token verifier schema ratification 必须定义：

- Token record owner。
- Non-plaintext verifier storage。
- Verifier algorithm 与 comparison rules。
- 可以安全写日志的 token record identifier。
- Actor kind 与 actor identifier relationship。
- Player access token 与 player account 的关系。
- 支持 rotation 与 presented-token logout 所需的 credential-token linkage。
- Audience 与 route eligibility semantics。
- 最小状态：`active`、`expired`、`revoked`。
- `issued_at`、`expires_at` 与 `revoked_at` semantics。
- Revocation reason retention。
- Replacement 或 rotation lineage。
- Expired 与 revoked record retention。
- Cleanup owner 与 cleanup trigger。
- Replay-sensitive failure classes。
- Logs、errors、traces、tests、fixtures 和 conversation logs 的 redaction rules。

第一版姿态的必需约束：

```yaml
raw_token_storage: forbidden
minimum_entropy_bits: 256
access_token_ttl: 1h
refresh_token_storage: forbidden_for_first_posture
session_token_vocabulary: deferred_until_session_persistence
token_record_id_log_safe: required
```

第一版 token verifier schema 必须能够 revoke presented access token。它也应该保留足够的 credential linkage，以便在未来授权实现时支持撤销同一 credential 下的 previous active tokens。

## 6. External Identity Link Gate

第一版 `device_credential_login` 姿态不要求 external identity storage。

Gate 状态：

```yaml
required_for_first_posture: false
status: deferred
schema_added_now: false
migration_added_now: false
provider_dependency_added_now: false
```

未来 external identity schema ratification 必须定义：

- Provider namespace semantics。
- Provider subject semantics。
- Provider subjects 是 globally unique 还是 provider-scoped。
- 一个 account 是否可以有多个 provider links。
- 一个 provider subject 是否可以映射到多个 vibit accounts。
- Link、unlink、conflict、recovery 和 merge behavior。
- Provider metadata retention 与 redaction。
- 哪些 events 只是 audit-only，哪些可以 client-visible。
- 采用哪些 provider dependencies，以及它们允许存在于哪里。

本 gate 保持 provider identity 与 credential storage、token verifier storage、runtime sessions 和 player account lifecycle storage 分离。

## 7. Runtime Session Record Gate

第一版 access-token 姿态不要求 runtime session persistence。

Gate 状态：

```yaml
required_for_first_posture: false
status: deferred
schema_added_now: false
migration_added_now: false
session_store_selected_now: false
websocket_connection_binding_selected_now: false
```

未来 runtime session schema ratification 必须定义：

- Runtime sessions 是否作为 durable records 存在。
- PostgreSQL 是否足够，还是需要 Redis-like store。
- Session identifier semantics。
- Session token semantics，如果有。
- Access-token 与 session 的关系，如果有。
- WebSocket connection binding。
- Reconnect behavior。
- Connection epoch behavior。
- Expiration、revocation 与 cleanup。
- Validation 是 per request、first message、WebSocket handshake，还是 hybrid model。
- 是否改变 Protobuf envelope fields。

在该 gate 被 ratify 之前，当前 Protobuf `Session` fields 仍然只是 metadata-only，不能变成 proof。

## 8. Audit Persistence Gate

W-0070 定义了语义层面的 audit-oriented events。W-0071 不添加 durable audit persistence。

未来 durable audit schema ratification 必须定义：

- Authentication audit 存在 module tables、event log、outbox、operational logs，还是其他已批准 store。
- 哪些 semantic events 会变成 durable rows。
- 哪些 identifiers 可以安全存储。
- 哪些 fields 明确禁止。
- Retention 与 cleanup policy。
- Query 与 inspection requirements。
- Audit writes 是否和 credential 或 token verifier mutations 共享同一个 unit of work。

禁止出现在 audit payload 中的材料：

- Raw credential values。
- Raw access-token values。
- Token verifier hashes。
- Password hashes。
- Provider secrets。
- Full provider payloads。
- WebSocket connection secrets。

## 9. Player Account Lifecycle 保持不变

当前 player account lifecycle storage 仍然是：

```text
player_accounts
player_account_events
```

这些 tables 仍然只用于 lifecycle storage。

禁止放入 player lifecycle storage 的内容：

- Credential columns 或 rows。
- Password hashes。
- External identity provider subjects。
- Access-token state。
- Refresh-token state。
- Runtime session state。
- WebSocket connection state。
- Request validation results。
- Raw authentication proof。

未来 authentication work 可以通过显式 repository boundaries 引用 player accounts，但不得把 player lifecycle storage 改造成 authentication storage。

## 10. 未来工作拆分

未来实现至少必须这样拆分：

```yaml
credential_schema_ratification: separate_work
token_verifier_schema_ratification: separate_work
credential_and_token_migration_sources: completed_by_W_0077_and_W_0078
repository_interfaces: separate_work
postgres_adapters: separate_work
redaction_and_schema_tests: separate_work
runtime_authentication_wiring: separate_work
protocol_or_websocket_changes: separate_decision_if_needed
```

Agent 不得在一个宽泛 change 中同时合并 schema ratification、migration creation、repository implementation 和 runtime authentication behavior，除非未来 work item 明确授权这种范围。

## 11. Reference Alignment

### Nakama

Nakama 仍然是 accounts、authentication methods、session tokens、refresh、logout、expiry 和 revocation vocabulary 的 capability reference。

vibit 会吸收 capability coverage，但不会复制 Nakama 的 public API、token/session model 或 storage shape。Refresh tokens 与 session token vocabulary 仍然 deferred。

### Pitaya

Pitaya 仍然是 sessions、handler context、connection binding、frontend/backend roles 与 realtime server routing 的 vocabulary reference。

vibit 会吸收 connection acceptors 与 application identity 分离的思想。它不会把 credential 或 token validation 放进 WebSocket acceptors，也不会复制 Pitaya 的 public API shape。

## 12. Verification 路径

本标准的默认 verification：

```bash
node tools/vibit inspect next --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check agent-tooling --json
node tools/vibit check memory --json
node tools/vibit check work --json
node tools/vibit check change define-credential-token-session-schema-gates --json
node tools/vibit check all --json
git diff --check
```

本变更只定义 gate，不改 Go runtime behavior，因此不要求 Go runtime tests。

## 13. 非授权

本标准不授权：

- Credential tables。
- Token tables。
- External identity tables。
- Runtime session tables。
- Audit tables。
- Migrations。
- Repository interfaces。
- PostgreSQL adapters。
- Runtime credential lookup。
- Token generation、parsing、validation、refresh、revocation、rotation、replay handling、cleanup 或 storage。
- Login handlers。
- Runtime player handlers。
- WebSocket routes。
- Protobuf messages 或 generated Protobuf output。
- WebSocket handshake authentication。
- First system-message authentication。
- Password hashing、JWT、OAuth、OIDC、provider SDK、Redis-like、cryptography、key-management 或 major authentication dependencies。
- 把 metadata-only `player_id`、`session_id`、`connection_id` 或 `connection_epoch` 当作 proof。

## 14. 后续

下一步工作：

```text
W-0079 Add authentication migration static checks
```

W-0079 应强化已 ratified 的 authentication migration sources 的本地 checks，且不添加 repositories、adapters、runtime token validation、generated output、Protobuf messages、WebSocket behavior 或 authentication dependencies。
