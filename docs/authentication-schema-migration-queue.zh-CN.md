# Authentication Schema Migration Queue

Status: Draft v0.1
Last updated: 2026-05-14
Scope: credential 与 token verifier schema boundaries 之后的 authentication schema migration、repository、adapter、redaction 与 verification 计划队列
Depends on: `docs/credential-record-schema-boundary.md`, `docs/token-verifier-record-schema-boundary.md`, `docs/postgresql-persistence-boundary.md`, `docs/postgresql-verification-environment.md`
Canonical decision: `ADR-0034`

配套英文源文档是 `docs/authentication-schema-migration-queue.md`。英文文件是权威版本。

## 1. Purpose

本标准定义两个必需 authentication schema boundaries 都 ratify 之后的下一个 bounded work queue：

```yaml
credential_record_schema_boundary: ratified_no_schema_added
token_verifier_record_schema_boundary: ratified_no_schema_added
default_durable_target: PostgreSQL
implementation_authorized_now: false
```

目标是让未来 authentication storage path 可确定地推进，而不是从 schema ratification 直接跳到宽泛的 authentication implementation。

本文规划 migration order、repository-interface gates、PostgreSQL adapter gates、redaction checks、live PostgreSQL verification expectations 和 milestone closeout。它不添加 SQL migration source、tables、repository interfaces、PostgreSQL adapters、runtime credential lookup、token issuance、token validation、logout、refresh、cleanup、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior。

## 2. Required Reading

应与以下文件一起阅读：

- `docs/credential-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/selected-login-token-boundary-checks.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/postgresql-persistence-boundary.md`
- `docs/postgresql-verification-environment.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `ADR-0029`
- `ADR-0032`
- `ADR-0033`
- `ADR-0034`

Reference reading：

- Nakama authentication concepts：`https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts：`https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya features and session vocabulary：`https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya API and handler vocabulary：`https://pitaya.readthedocs.io/en/latest/API.html`

Nakama 和 Pitaya 只作为 game-backend capability coverage 与 vocabulary reference。它们不支配 vibit 的 migration order、repository shape、public API 或 transport behavior。

## 3. Queue Rule

Authentication persistence 必须通过可追踪的 queue steps 推进：

```yaml
schema_boundary:
  credential: completed_by_W_0074
  token_verifier: completed_by_W_0075
migration_queue_planning:
  completed_by: W_0076
migration_sources:
  credential: separate_future_work
  token_verifier: separate_future_work
schema_static_checks:
  separate_future_work_or_migration_work
repository_interfaces:
  separate_future_work_after_migrations
postgres_adapters:
  separate_future_work_after_repository_interfaces
redaction_and_live_verification:
  separate_future_work_after_adapters
milestone_closeout:
  separate_future_work
runtime_authentication:
  blocked_until_later_milestone
```

Agents 不得在一个宽泛 change 中合并 migration source creation、repository interface design、adapter implementation、runtime validation、protocol behavior 和 authentication behavior。

## 4. Planned Work Items

本 planning step 之后的 M-014 planned queue 是：

| Work item | Title | Scope |
| --- | --- | --- |
| `W-0077` | Add credential PostgreSQL migration source | 只创建已 ratified 的 `authentication_device_credentials` SQL source。 |
| `W-0078` | Add token verifier PostgreSQL migration source | Credential migration 存在之后，只创建已 ratified 的 `authentication_access_tokens` SQL source。 |
| `W-0079` | Add authentication migration static checks | 强化 authentication migration naming、ownership、forbidden raw secret columns 和 player lifecycle table separation 的 repository checks。 |
| `W-0080` | Define authentication repository interface boundary | 为 credential 与 token verifier storage 定义 storage-neutral interfaces 和 mutations，不加 adapters 或 runtime behavior。 |
| `W-0081` | Define authentication PostgreSQL adapter boundary | 预留 adapter source paths、transaction expectations、SQL operation scope 和 focused tests，不实现 adapter。 |
| `W-0082` | Close credential and token verifier schema ratification milestone | 按 exit criteria 回顾 M-014，并打开下一个 confirmation 或 implementation gate。 |

这个 queue 有意不包含 runtime login、token issuance、token validation、logout、cleanup、Protobuf、WebSocket 或 dependency implementation work。

## 5. Migration Order

第一版 authentication migration order 是：

```text
000003_create_authentication_device_credentials.sql
000004_create_authentication_access_tokens.sql
```

Rationale：

- Credential records 必须先存在，token verifier records 才能引用 `credential_record_id`。
- Player account lifecycle tables 已经存在于 `000002_create_player_account_state.sql`。
- Token verifier records 依赖 credential schema boundary 来支持 rotation 和 previous-token revocation linkage。
- 分离 credential 与 token verifier migrations 让每个 SQL source 更小，更容易由 agents 验证。

未来 migration work 只有在现有 migration files 已占用这些编号时才可以调整 exact sequence numbers。如果编号变化，change spec 必须解释实际 sequence，并保持 credential migration 在 token verifier migration 之前。

## 6. Credential Migration Gate

`W-0077` 只能添加一个 SQL-first migration source：

```text
runtime/migrations/postgres/000003_create_authentication_device_credentials.sql
```

Allowed scope：

- 创建 `authentication_device_credentials`。
- 添加 `docs/credential-record-schema-boundary.md` 要求的 indexes 与 constraints。
- 添加 `-- +goose Up`、`-- +goose Down` 和 `-- Module: runtime.authentication`。
- 保持 player account lifecycle tables 不变。

Forbidden scope：

- Token verifier table creation。
- Repository interfaces 或 adapters。
- Runtime credential lookup。
- Login behavior。
- Password hashing、OAuth、OIDC、provider SDK、JWT、key-management 或 Redis-like dependencies。
- 可隐藏 credentials、tokens、provider payloads 或 transport state 的 generic metadata columns。

## 7. Token Verifier Migration Gate

`W-0078` 只能添加一个 SQL-first migration source：

```text
runtime/migrations/postgres/000004_create_authentication_access_tokens.sql
```

Allowed scope：

- 创建 `authentication_access_tokens`。
- 添加 `docs/token-verifier-record-schema-boundary.md` 要求的 indexes 与 constraints。
- 添加 `-- +goose Up`、`-- +goose Down` 和 `-- Module: runtime.authentication`。
- 保持 player account lifecycle tables 不变。

Forbidden scope：

- Credential table changes，除非为了保持 migration ordering 且明确说明理由。
- Repository interfaces 或 adapters。
- Runtime token issuance、validation、logout、refresh 或 cleanup。
- JWT、signing、OAuth、OIDC、provider SDK、key-management、password-hashing 或 Redis-like dependencies。
- 可隐藏 raw tokens、claims、transport state、provider payloads 或 private request details 的 generic metadata columns。

## 8. Static Check Gate

`W-0079` 应更新 repository checks，使 authentication migrations 在 repositories 和 adapters 存在前就可检查。

Expected checks：

- Authentication migration files 使用 deterministic sequence numbers。
- Authentication migration files 包含 goose Up 和 Down sections。
- Authentication migration files 包含 `-- Module: runtime.authentication`。
- `authentication_device_credentials` 和 `authentication_access_tokens` 是 M-014 migration work 引入的唯一 authentication tables。
- Raw credential 与 raw token columns 被禁止。
- `player_accounts` 和 `player_account_events` 保持 credential-free 与 token-free。
- JSON metadata columns 继续不存在，除非后续 ADR 明确 ratify。
- Repository-relative JSON check output 在所有平台都使用 forward slashes。

这个 check gate 默认必须保持 static 和 local。

## 9. Repository Interface Gate

`W-0080` 只能在 migration sources 和 static checks 存在之后定义 storage-neutral interfaces。

Expected ownership：

```text
runtime/internal/modules/authentication/
```

Expected interface families：

- Credential lookup and mutation repository boundary。
- Token verifier creation、lookup、revocation、expiration query 和 cleanup eligibility boundary。
- Unit-of-work expectations，确保未来 implementation 授权后 credential、player account 与 token mutations 可以保持 atomic。

Forbidden scope：

- PostgreSQL adapter implementation。
- Runtime login handlers。
- Token issuance 或 validation。
- Logout behavior。
- Cleanup jobs。
- WebSocket routes 或 proof carriers。
- Protobuf messages。
- Authentication dependencies。

## 10. PostgreSQL Adapter Boundary Gate

`W-0081` 可以在 storage-neutral interfaces 存在之后定义 adapter boundaries。

Expected ownership：

```text
runtime/internal/platform/persistence/postgres/
```

Expected boundary content：

- Reserved source and test paths。
- Constructor names。
- Caller-owned executor and transaction rules。
- SQL operation scope。
- Error mapping expectations。
- Fake-executor test expectations。
- Optional live PostgreSQL verification expectations。

Forbidden scope：

- Runtime authentication implementation。
- WebSocket route 或 handshake behavior。
- Protobuf behavior。
- Token generation 或 cryptographic verifier implementation。
- External provider 或 Redis-like dependencies。

## 11. Verification Expectations

Default verification 保持 local：

```bash
node tools/vibit check migrations --json
node tools/vibit check runtime --json
node tools/vibit check all --json
```

Live PostgreSQL verification 继续通过以下环境变量 opt-in：

```text
VIBIT_POSTGRES_TEST_DSN
VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1
```

未来 migration 与 adapter work 应记录以下之一：

```text
Verified: live PostgreSQL verification ran against VIBIT_POSTGRES_TEST_DSN
Not verified: live PostgreSQL verification skipped because VIBIT_POSTGRES_TEST_DSN was not set
Not applicable: no migration, table, repository, adapter, or persistence behavior changed
```

W-0076 不需要 live PostgreSQL verification，因为本 planning step 不添加 schema 或 runtime behavior。

## 12. Redaction Expectations

未来 authentication migration、repository、adapter 和 test work 必须保持 redaction。

以下内容禁止出现在 logs、errors、traces、fixtures、ADRs、change specs 和 conversation logs 中：

- Raw credential proof。
- Raw access tokens。
- Credential lookup 或 verifier digests。
- Token lookup 或 verifier digests。
- Server-side verifier secrets 或 peppers。
- Authorization headers。
- Cookie contents。
- WebSocket subprotocol token carriers。
- URL query token carriers。
- Provider secrets。

可谨慎使用：

- `credential_record_id`
- `token_record_id`
- `player_id`
- lifecycle state names
- non-secret reason catalog values
- stable failure class names

## 13. Non-Authorization

本标准不授权：

- Authentication migration source。
- Credential 或 token tables。
- Authentication repository interfaces。
- PostgreSQL authentication adapters。
- Runtime credential lookup。
- Token generation、issuance、parsing、validation、refresh、revocation、rotation、replay handling、cleanup 或 storage behavior。
- Login handlers。
- Runtime player handlers。
- WebSocket routes。
- Protobuf messages 或 generated Protobuf output。
- WebSocket proof carriers。
- WebSocket handshake authentication。
- Password hashing、JWT、OAuth、OIDC、provider SDK、Redis-like、cryptography、key-management 或 major authentication dependencies。
- 把 metadata-only `player_id`、`session_id`、`connection_id` 或 `connection_epoch` 当作 proof。

## 14. Follow-Up

下一步工作：

```text
W-0077 Add credential PostgreSQL migration source
```

下一项工作只能在声明边界内添加 credential migration source 和相关 static manifest/check updates。
