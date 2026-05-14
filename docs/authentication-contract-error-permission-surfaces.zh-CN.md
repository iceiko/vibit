# Authentication Contract、错误与权限表面

状态：草案 v0.1
最后更新：2026-05-14
范围：所选第一版登录与 token 姿态的语义 contract、错误、权限和审计表面
依赖：`docs/token-lifecycle-storage-implications.md`
权威决策：`ADR-0028`

英文文件 `docs/authentication-contract-error-permission-surfaces.md` 是权威版本。本文是面向中文读者的人类可读翻译。

## 1. 目的

本文定义 vibit 在实现所选第一版 authentication 姿态之前必须具备的语义表面：

```yaml
login_method: device_credential_login
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_device_credential_login
```

目标是在 runtime code 存在之前，让未来 authentication 工作对 agent 可读、可检查、可验证。

本文不实现 authentication、token 生成、token 验证、logout、refresh、credential lookup、token storage、audit persistence、runtime handler、WebSocket route、Protobuf message、generated output、migration 或 schema change。

## 2. 问题

Authentication 很容易让框架代码变得难以被 agent 安全修改。

如果没有显式语义表面，未来 agent 可能会：

- 把 credential parsing 放进 WebSocket transport。
- 把 token validation 放进 Protobuf adapter 或 domain module。
- 把 metadata-only `player_id` 或 `session_id` 当作 proof。
- 添加临时 error code。
- 在 permission 与 failure model 存在之前添加 route handler。
- 认为 access-token issuance 隐含 refresh token behavior。
- 把 raw token 或 credential material 存进日志、测试、event 或数据库记录。

解决方式是在现在定义 contract、error、permission 和 audit surfaces，同时继续把实现阻挡在 schema、repository、test 和 runtime milestone gates 之后。

## 3. 规则

未来 authentication implementation 必须从以下目录中的已登记 semantic contract sources 开始：

```text
contracts/runtime/authentication/
```

runtime authentication family 登记在：

```text
.arch/contracts.yaml
```

在添加 authentication behavior 前，agent 必须使用：

```bash
node tools/vibit inspect contracts --module runtime --json
node tools/vibit check contracts --json
```

Runtime authentication 仍由 application 层拥有：

```text
runtime/internal/app
```

除非后续 ADR 改变所有权，否则其他层不得拥有 credential parsing、access-token validation 或第一版 logout behavior。

## 4. 已选表面

第一批选定语义表面为：

```yaml
commands:
  - AuthenticateWithDeviceCredential
  - ValidateAccessToken
  - LogoutAccessToken
  - RefreshAccessToken
events:
  - AuthenticationSucceeded
  - AuthenticationFailed
  - TokenIssued
  - TokenValidationFailed
  - TokenRevoked
  - LogoutRequested
errors:
  - authentication_errors
permissions:
  - authentication_permissions
queries: []
```

这些只是 semantic source contracts。它们不会创建 generated Go contract shapes、Protobuf messages、runtime handlers、WebSocket routes、database tables 或 migrations。

## 5. Command 表面

### `AuthenticateWithDeviceCredential`

目的：

```text
使用所选 high-entropy device credential login method 认证 player，并颁发 opaque access token。
```

Contract source:

```text
contracts/runtime/authentication/commands/AuthenticateWithDeviceCredential.yaml
```

该 command 只能在未来 gates 定义 credential schema、token verifier schema、repository boundaries、redaction tests 和 runtime authentication implementation milestone 后实现。

它不得把原始操作系统 device ID、player ID、session ID、connection ID 或其他 metadata 当作 credential proof。

### `ValidateAccessToken`

目的：

```text
在 domain dispatch 前，把显式 opaque access-token proof 验证成 application-owned request identity。
```

Contract source:

```text
contracts/runtime/authentication/commands/ValidateAccessToken.yaml
```

已选 request proof carrier 是 explicit request payload proof。当前 Protobuf `Session` metadata 仍然只是 metadata-only，不能通过重新解释变成 proof。

### `LogoutAccessToken`

目的：

```text
吊销当前呈递的 opaque access token。
```

Contract source:

```text
contracts/runtime/authentication/commands/LogoutAccessToken.yaml
```

第一版 logout scope 是：

```yaml
logout_scope: presented_access_token
```

这不授权 logout-all-sessions、account-wide token revocation、credential-wide token revocation、WebSocket close behavior 或 admin revocation。

### `RefreshAccessToken`

目的：

```text
保留 token refresh 语义表面，同时让 refresh token 不进入第一版实现。
```

Contract source:

```text
contracts/runtime/authentication/commands/RefreshAccessToken.yaml
```

Refresh token 不属于第一版姿态。第一版 renewal method 仍然是：

```yaml
renewal_method: reauthenticate_with_device_credential_login
```

Refresh contract 的存在，是为了让 agent 明确看到 refresh 被有意考虑并在第一版中有意不支持，而不是被遗忘。

## 6. 错误表面

错误目录是：

```text
contracts/runtime/authentication/errors/authentication_errors.yaml
```

第一批必需错误族包括：

```yaml
missing_proof:
  - AUTHENTICATION_PROOF_MISSING
  - AUTHENTICATION_TOKEN_MISSING
malformed_proof:
  - AUTHENTICATION_PROOF_MALFORMED
  - AUTHENTICATION_TOKEN_MALFORMED
invalid_proof:
  - AUTHENTICATION_CREDENTIAL_INVALID
  - AUTHENTICATION_TOKEN_INVALID
expired_proof:
  - AUTHENTICATION_TOKEN_EXPIRED
revoked_proof:
  - AUTHENTICATION_TOKEN_REVOKED
unsupported_proof:
  - AUTHENTICATION_REFRESH_NOT_SUPPORTED
actor_disabled:
  - AUTHENTICATION_ACCOUNT_DISABLED
validator_unavailable:
  - AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
  - AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
not_implemented:
  - AUTHENTICATION_NOT_IMPLEMENTED
```

错误必须是 public-safe 的，并且不得包含 raw credential material、raw token values、token prefixes、verifier hashes、password hashes、provider secrets 或隐藏验证细节。

## 7. 权限表面

权限目录是：

```text
contracts/runtime/authentication/permissions/authentication_permissions.yaml
```

第一批权限：

```yaml
authentication_device_credential_login:
  dimension: unauthenticated_login_entrypoint
authentication_access_token_validate:
  dimension: validation_infrastructure_permission
authentication_access_token_logout:
  dimension: player_token_lifecycle_permission
authentication_access_token_refresh:
  dimension: deferred_token_lifecycle_permission
authentication_admin_revoke_token:
  dimension: deferred_admin_permission
```

这些权限本身不授予 domain module authority。Domain modules 在 validation 之后消费 normalized request identity，并使用自己的 permission policies。

Metadata-only identity 仍不足以作为 production permission basis。

## 8. 审计事件表面

第一批 authentication audit-oriented event surfaces 是：

```yaml
AuthenticationSucceeded:
  reason: successful selected login proof
AuthenticationFailed:
  reason: rejected or unavailable selected login proof
TokenIssued:
  reason: opaque access-token verifier record issued
TokenValidationFailed:
  reason: access-token proof rejected or unavailable
TokenRevoked:
  reason: token verifier record revoked
LogoutRequested:
  reason: presented-token logout requested
```

这些 event 只是语义表面。它们不添加 event bus、client event stream、audit table、durable audit persistence 或 runtime publication behavior。

Raw credentials、raw tokens、token verifier hashes、password hashes、provider secrets 和完整 provider payloads 禁止出现在这些 event 中。

## 9. Account Linking 表面

Account linking 不属于第一版已选姿态的范围。

状态：

```yaml
account_linking: deferred
external_identity_storage_required_for_first_posture: false
```

未来 account linking 需要 provider subject semantics、link/unlink permissions、conflict behavior、recovery behavior、merge behavior、schema gates、audit events 和 tests。

## 10. Protocol 与 Transport 影响

本标准不改变：

- Protobuf envelope fields。
- 当前 Protobuf `Session` metadata semantics。
- WebSocket handshake authentication。
- WebSocket transport behavior。
- First system-message authentication。
- Runtime player handlers。
- WebSocket routes。

在后续 protocol decision ratify wire shape 之前，已选 request proof carrier 仍然是 semantic explicit request payload proof。

## 11. Storage 影响

本标准不添加 storage。

未来实现仍然受 W-0071 schema gates 阻挡，涉及：

- Credential records。
- Token verifier records。
- Optional session records，如果后续姿态需要。
- External identity records，如果后续 login/linking 姿态需要。

Player account lifecycle tables 仍必须保持 credential-free、token-free、external-identity-free 和 session-free。

## 12. Agent 影响

本标准让下一步实现边界可检查：

```bash
node tools/vibit inspect contracts --module runtime --json
node tools/vibit inspect contract --module runtime --type command --id AuthenticateWithDeviceCredential --json
node tools/vibit inspect contract --module runtime --type command --id ValidateAccessToken --json
node tools/vibit check contracts --json
```

Agent 在编辑 implementation code 前，应先阅读相关 command、event、error 和 permission manifests。

## 13. Human 影响

人类可以在不阅读未来 Go code 的情况下，用稳定词汇讨论第一版 authentication slice。

代价是仓库会在 runtime behavior 存在之前增加更多 manifest files。这是有意选择：本项目更重视显式、可检查的边界，而不是隐含 security behavior。

## 14. Verification 路径

本标准的默认 verification：

```bash
node tools/vibit inspect next --json
node tools/vibit check contracts --json
node tools/vibit inspect contracts --module runtime --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check agent-tooling --json
node tools/vibit check memory --json
node tools/vibit check work --json
node tools/vibit check change define-authentication-contract-error-permission-surfaces --json
node tools/vibit check all --json
git diff --check
```

本步骤只添加设计与 contract source，不改 Go runtime behavior，因此不要求 Go tests。

## 15. Migration 路径

未来实现应按以下顺序推进：

1. 定义 credential、token 和 session schema gates。
2. 为所选 login/token boundaries 添加 repository checks。
3. 关闭 M-013。
4. 选择 implementation milestone。
5. 在 runtime behavior 前添加 schema、migrations 和 repository boundaries。
6. 在 `runtime/internal/app` 后面添加 runtime authentication interfaces。
7. 只有在单独 Protobuf 或 WebSocket 决策需要时，才添加 protocol/wire behavior。

## 16. 非授权

本标准不授权：

- Runtime authentication code。
- Login handlers。
- Token generation、parsing、validation、refresh、revocation、rotation、replay handling 或 storage。
- Credential tables。
- External identity tables。
- Token tables。
- Session tables。
- Migrations。
- Generated contract shapes。
- Protobuf messages 或 generated Protobuf output。
- Password hashing、JWT、OAuth、OIDC、provider SDK、Redis-like、cryptography、key-management 或 major authentication dependencies。
- Protobuf envelope changes。
- WebSocket handshake authentication。
- First system-message authentication。
- Runtime player handlers。
- WebSocket routes。
- 把 metadata-only `player_id`、`session_id`、`connection_id` 或 `connection_epoch` 当作 proof。

## 17. 后续

下一步工作：

```text
W-0071 Define credential token session schema gates
```

W-0071 必须定义未来 schema gates，但不添加 migrations、repository implementations、runtime lookup code、handlers 或 routes。
