# Selected Login And Token Boundary Checks

状态：Draft v0.1
最后更新：2026-05-14
范围：针对第一版已选 login method、opaque token 姿态、schema gates、dependency deferral、generated output deferral、protocol deferral 与 runtime implementation deferral 的 repository checks
依赖：`docs/login-method-token-format-ratification.md`、`docs/token-lifecycle-storage-implications.md`、`docs/authentication-contract-error-permission-surfaces.md`、`docs/credential-token-session-schema-gates.md`
权威决策：`ADR-0030`

对应英文原文是 `docs/selected-login-token-boundary-checks.md`。英文文件是权威版本。

## 1. 目的

本标准定义保护 vibit 第一版已选 authentication 姿态的窄范围 repository checks：

```yaml
login_method: device_credential_login
access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_device_credential_login
default_durable_target: PostgreSQL
implementation_authorized: false
```

这些检查存在的原因是：姿态已经选定，但 runtime behavior 仍未授权。选择不等于实现。

本标准不添加 runtime authentication、login handlers、token generation、token parsing、token validation、logout behavior、refresh behavior、credential lookup、credential storage、token storage、session persistence、external identity linking、Protobuf messages、generated Go contract shapes、WebSocket route behavior、WebSocket handshake authentication、migrations、repositories、PostgreSQL adapters、audit persistence、provider dependencies、signing dependencies、password-hashing dependencies、key-management dependencies 或 Redis-like dependencies。

## 2. 规则

repository check rule 是：

```text
runtime.selected_login_token_boundary
```

默认命令是：

```bash
node tools/vibit check runtime --json
```

该规则也必须通过：

```bash
node tools/vibit check all --json
```

该检查是静态和本地的。它不得要求 live PostgreSQL、Docker、Podman、cloud services、OAuth providers、OIDC providers、platform identity providers、Redis-like services 或 network access。

## 3. 检查保护什么

该检查保护这些边界：

- 已选 login method 仍然只是 semantic-contract-only。
- Opaque access-token behavior 仍然只是 semantic-contract-only。
- Refresh tokens 在第一版实现姿态中仍然 unsupported。
- Credential 与 token verifier schema gates 仍然是 defined but not implemented。
- Runtime session persistence 仍然 deferred。
- External identity linking 仍然 deferred。
- Authentication generated output 仍然 deferred。
- Authentication Protobuf source 仍然 deferred。
- WebSocket transport 仍然 credential-neutral。
- 当前 Protobuf `Session` metadata 仍然只是 metadata-only，而不是 proof。
- Player account lifecycle storage 仍然 credential-free、token-free、external-identity-free、session-free、WebSocket-state-free。
- 新的 authentication dependencies 仍然 blocked，直到 adoption records 与 implementation gates 授权。

## 4. 必需静态信号

该检查要求以下静态信号继续存在：

```yaml
selected_posture:
  login_method: device_credential_login
  token_format: opaque_high_entropy_token
  token_issuance_carrier: login_command_response_token
  request_proof_carrier: explicit_request_proof_payload
  access_token_ttl: 1h
  refresh_token: not_in_first_implementation
  renewal_method: reauthenticate_with_device_credential_login
  logout_scope_first_posture: presented_access_token
schema_gates:
  credential_record_schema_gate_status: ratified_no_schema_added
  credential_record_schema_boundary: docs/credential-record-schema-boundary.md
  token_verifier_record_schema_gate_status: ratified_no_schema_added
  token_verifier_record_schema_boundary: docs/token-verifier-record-schema-boundary.md
  external_identity_link_schema_gate_status: deferred_no_schema_added
  runtime_session_record_schema_gate_status: deferred_no_schema_added
implementation_status:
  runtime_authentication_implemented: false
  token_behavior_implemented: false
  credential_storage_implemented: false
  token_storage_schema_added: false
  credential_storage_schema_added: false
  external_identity_storage_schema_added: false
  session_storage_schema_added: false
  migration_sources_added: false
  repository_interfaces_added: false
  postgres_adapters_added: false
  runtime_lookup_added: false
  websocket_handshake_authentication_changed: false
  runtime_player_handlers_added: false
  websocket_routes_added: false
  protobuf_envelope_changed: false
  generated_contract_shapes_added: false
```

这些信号不能替代 schemas 或 tests。它们是 tripwires，用来迫使未来 implementation work 明确改变 architecture state。

## 5. 禁止的捷径

在未来 bounded work item 授权实现之前，agent 不得添加：

- `AuthenticateWithDeviceCredential` runtime behavior。
- `ValidateAccessToken` runtime behavior。
- `LogoutAccessToken` runtime behavior。
- `RefreshAccessToken` runtime behavior。
- `AuthService`、`Authenticator`、`TokenValidator`、`TokenIssuer`、`TokenVerifier`、`TokenRepository` 或 `CredentialRepository` runtime implementation。
- Token 或 credential random generation code。
- Runtime code 中的 bearer-token parsing 或 acceptance。
- WebSocket `Authorization`、`Bearer`、`Cookie`、`Sec-WebSocket-Protocol` 或 handshake header authentication behavior。
- `proto/vibit/runtime/` 下的 runtime authentication Protobuf source。
- `runtime/internal/generated/contracts/runtime/authentication/` 下的 generated Go authentication contract shapes。
- Credential、token、refresh-token、runtime-session、external-identity、provider-subject 或 authentication-audit migrations。
- 在 player account lifecycle tables 中存储 credential、token、external identity、session、request validation 或 WebSocket state 的变更。
- 用于 authentication 的 JWT、OAuth、OIDC、password-hashing、provider SDK、key-management 或 Redis-like dependencies。

## 6. 当前允许的工件

以下工件是允许的，因为它们是 semantic 或 gate artifacts，而不是 runtime implementation：

- `contracts/runtime/authentication/**`
- `docs/authentication-contract-error-permission-surfaces.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/credential-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/first-login-method-set.md`
- `docs/first-token-format-proof-carrier-posture.md`
- 标记 implementation is deferred 的 architecture manifest markers。
- 解释边界的 agent-facing guides。

Generated Go contract shapes 仍然 deferred。Protobuf wire messages 仍然 deferred。Runtime handlers 仍然 deferred。

## 7. 机器可读输出

每条 check item 必须包含：

```yaml
rule_id: runtime.selected_login_token_boundary
artifact: <repo-relative-forward-slash-path>
```

JSON output 中的 repository-relative paths 必须在所有平台使用 forward slash，包括 Windows。

当 agent 需要 rule metadata 时，应使用：

```bash
node tools/vibit inspect rule runtime.selected_login_token_boundary --json
```

当 agent 需要 actionable results 时，应使用：

```bash
node tools/vibit check runtime --json
```

## 8. 与 Nakama 和 Pitaya 的关系

Nakama 和 Pitaya 仍然是 game backend capability 与 Go game-server vocabulary 的主动参考。

本检查不复制它们的 authentication API shapes。它保护 vibit 的 agent-native 路径，避免过早把参考框架行为导入 transport、Protobuf、player persistence 或 domain modules。

Reference alignment：

```yaml
nakama:
  device_style_login: adapted_as_capability_reference
  session_token_refresh_logout_vocabulary: adapted_as_lifecycle_reference
  direct_api_compatibility: rejected_for_now
pitaya:
  handler_session_context: adapted_as_request_identity_vocabulary
  session_binding: deferred
  transport_owned_authentication: rejected
```

## 9. 验证路径

凡是触碰该边界的变更，应运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check agent-tooling --json
node tools/vibit check memory --json
node tools/vibit check work --json
node tools/vibit check all --json
```

如果存在 change spec，还应运行：

```bash
node tools/vibit check change <change-id> --json
```

Live PostgreSQL verification 不是本检查的要求，因为 W-0072 不添加 schema、migrations、repositories、adapters 或 runtime behavior。

## 10. 未来迁移路径

未来 implementation milestone 可以有意改变这个边界。该工作必须：

- 创建或更新 change spec。
- 更新英文标准和简体中文译文。
- 更新相关 ADR 或创建新的 ADR。
- 更新 `.arch/runtime.yaml`、`.arch/conventions.yaml`、`.arch/contracts.yaml`、`.arch/reference.yaml` 和相关 module manifests。
- 先 ratify schema，再添加 migrations。
- 先 ratify repository interfaces，再添加 adapters。
- 先 ratify generated output，再生成 shapes。
- 先 ratify Protobuf impact，再添加 wire messages。
- 先 ratify WebSocket carrier behavior，再改变 handshake 或 route behavior。
- 添加 focused tests。
- 更新 `runtime.selected_login_token_boundary`，让它只阻断未批准的捷径，并允许新批准的 slice。
