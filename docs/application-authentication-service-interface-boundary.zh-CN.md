# Application Authentication Service Interface Boundary

状态：Draft v0.1
最后更新：2026-05-15
范围：未来由 application 拥有的 runtime authentication service interface boundary
依赖：`docs/runtime-authentication-implementation-boundary.md`、`docs/authentication-generated-contract-shape-timing.md`、`docs/authentication-contract-error-permission-surfaces.md`
Canonical decision: `ADR-0039`

对应英文原文是 `docs/application-authentication-service-interface-boundary.md`。英文文件是权威版本。

## 1. 目的

本文定义未来 runtime authentication behavior 必须契合的 service-interface boundary。

它位于 metadata-only generated authentication contract shapes 之后、handwritten runtime authentication behavior 之前。它的作用是让未来 application service shape 对 Agent 可预测，但不添加 login execution、token generation、verifier comparison、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository interface changes 或 migration schema changes。

这只是 service-interface boundary。它不添加 Go service code。

## 2. 核心规则

Application authentication service interfaces 由 application 拥有。

未来 owner 是：

```text
runtime/internal/app
```

第一版实现必须通过 application-level request/result vocabulary 暴露 authentication behavior，然后 domain dispatch 才能接收 validated identity。

目标 ownership：

```text
protocol-decoded request proof
-> application authentication service interface
-> application unit of work
-> authentication.Repository
-> persistence-only PostgreSQL adapter
-> application-owned RequestIdentity handoff
-> domain dispatch
```

Transport、Protobuf adapter、domain module、generated contract shape、repository 或 PostgreSQL adapter 都不能成为 service-interface owner。

## 3. Planned Package Boundary

未来 package 可以直接位于 `runtime/internal/app` 下，也可以位于 `runtime/internal/app` 拥有的子 package，例如：

```text
runtime/internal/app/authentication
```

后续 implementation work item 必须在添加代码前选择具体 Go package path。

在后续 service-interface implementation gate 之后，允许的 interface-level dependencies：

- Go standard library types。
- Application-owned `RequestIdentity`、`RouteKey`、`Session`、`ApplicationResult` 和 application error vocabulary。
- 由 `runtime/internal/platform/tx` 表示的 application unit-of-work boundary。
- 通过 application unit of work 获取的 module-owned `authentication.Repository`。

禁止的 interface-level dependencies：

- WebSocket transport packages。
- Generated Protobuf packages。
- PostgreSQL driver packages。
- Migration tooling。
- JWT、OAuth、OIDC、password-hashing、provider SDK、Redis-like token/session store、S3 或 MinIO dependencies。
- 把 generated authentication contract shape packages 当作 runtime registry 或 behavior owner。

Generated authentication contract shapes 可以影响命名和 mapping tables，但未来 service interfaces 不能要求 domain modules 或 transport handlers 导入 generated authentication shape packages。

## 4. Service Vocabulary

以下路径中的 generated authentication contract shapes：

```text
runtime/internal/generated/contracts/runtime/authentication/
```

用于指导未来 service request/result vocabulary。

第一版 planned service surface 是：

```yaml
service_boundary:
  owner: runtime/internal/app
  status: boundary_defined_no_code
  commands:
    AuthenticateWithDeviceCredential:
      request_vocabulary:
        - credential_proof
        - requested_player_id
        - client_instance_id
        - account_creation_intent
      result_vocabulary:
        - authentication_status
        - actor_kind
        - player_id
        - access_token
        - token_type
        - issued_at
        - expires_at
        - token_record_id
    ValidateAccessToken:
      request_vocabulary:
        - access_token
        - route_kind
        - route_module
        - route_name
        - connection_id
        - connection_epoch
      result_vocabulary:
        - validation_status
        - proof_status
        - actor_kind
        - actor_id
        - player_id
        - player_id_validated
        - session_validated
        - token_record_id
    LogoutAccessToken:
      request_vocabulary:
        - access_token
        - logout_reason
      result_vocabulary:
        - revoked
        - logout_scope
        - token_record_id
        - revoked_at
    RefreshAccessToken:
      first_posture: unsupported_reserved
      required_error: AUTHENTICATION_REFRESH_NOT_SUPPORTED
```

这些 vocabulary 不是 public wire schema。Protobuf authentication messages 仍然 deferred。

## 5. Redaction Expectations

未来 service interfaces 必须先对 secret fields 分类，再实现。

Secret input fields：

- `credential_proof`
- `access_token`

Secret internal material：

- Raw credential material。
- Raw token material。
- Credential verifier digest inputs。
- Token verifier digest inputs。
- 当未来标准将其标记为 confidential 时的 verifier keys 或 key identifiers。

规则：

- Raw credential material and raw access token material must be redacted from logs、errors、conversation logs、change specs、test names、table rows 和 audit records。
- Public errors 必须使用已注册的 authentication error codes，而不是 raw proof details。
- 未来 tests 必须断言 redacted values 不会出现在 application errors 或 audit-safe facts 中。
- Token issuance 可以在后续 ratified response carrier 中把 raw access token 一次性展示给 client。这个例外必须由 token generation gate 测试并记录。

## 6. Unit-Of-Work And Repository Boundary

未来 service behavior 只能通过 application unit-of-work boundary 使用 `authentication.Repository`。

允许的形状：

```text
application service method
-> runner.WithinUnitOfWork(...)
-> UnitOfWork.NewAuthenticationRepository(...)
-> authentication.Repository
```

规则：

- Application service 拥有 orchestration。
- Repository 只拥有 storage-neutral record operations。
- PostgreSQL adapter 只拥有 SQL persistence。
- 对于会改变状态的 authentication operations，application service 不能绕开 unit-of-work boundary。
- Repository 不能 generate tokens、compare verifiers、validate proof、parse bearer tokens、publish events 或 return public protocol responses。
- PostgreSQL adapter 不能做 authentication decisions。

## 7. Request Identity Handoff

未来 `ValidateAccessToken` service behavior 必须在 production-sensitive domain dispatch 前，把 proven proof 转换为 `RequestIdentity`。

目标 handoff：

```yaml
request_identity_handoff:
  owner: runtime/internal/app
  input: ValidateAccessToken result
  output: RequestIdentity
  required_success_markers:
    status: validated
    actor_kind: player
    actor_id: player_id
    player_id_validated: true
    session_validated: false
  metadata_only_allowed_as_proof: false
```

`MetadataOnlySessionValidator` 仍然只是 bootstrap path，不满足这个 boundary。

Domain modules 必须消费 `RequestIdentity`；它们不能 parse、validate 或 compare authentication proof。

## 8. Error, Permission, And Audit Handoff

未来 service interfaces 必须映射到现有 semantic contracts。

Errors：

- `AuthenticateWithDeviceCredential` 使用 `AUTHENTICATION_PROOF_MISSING`、`AUTHENTICATION_PROOF_MALFORMED`、`AUTHENTICATION_CREDENTIAL_INVALID`、`AUTHENTICATION_ACCOUNT_DISABLED`、`AUTHENTICATION_RATE_LIMITED`、`AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE`、`AUTHENTICATION_TOKEN_STORE_UNAVAILABLE` 和 `AUTHENTICATION_NOT_IMPLEMENTED`。
- `ValidateAccessToken` 使用 missing、malformed、invalid、expired、revoked、unavailable、disabled-account 和 not-implemented token errors。
- `LogoutAccessToken` 使用 missing、malformed、invalid、token-store-unavailable 和 not-implemented token errors。
- `RefreshAccessToken` 在第一姿态中映射为 `AUTHENTICATION_REFRESH_NOT_SUPPORTED`。

Permissions：

- `authentication_device_credential_login`
- `authentication_access_token_validate`
- `authentication_access_token_logout`
- `authentication_access_token_refresh`

Audit handoff：

- `AuthenticationSucceeded`
- `AuthenticationFailed`
- `TokenIssued`
- `TokenValidationFailed`
- `TokenRevoked`
- `LogoutRequested`

Audit publication 和 audit persistence 仍然是 separate gates。Service interface 以后可以定义 audit-safe facts，但在单独的 storage path 被 ratify 前，不能存储 audit events。

## 9. Nakama And Pitaya Alignment

Nakama 仍然是 account authentication、session token issuance、token expiration、logout、revocation 和 realtime authenticated actor binding 的 capability reference。

Pitaya 仍然是 handler context、frontend acceptor separation、backend handler separation 和 route identity context 的 vocabulary reference。

vibit 通过 application-owned service boundary 吸收这些经验。它不复制 Nakama 或 Pitaya public APIs，也不会让 transport handlers 或 domain handlers 负责 token validation。

## 10. Verification Path

此 boundary 的 repository check rule 是：

```text
runtime.application_authentication_service_interface_boundary
```

涉及此 boundary 的变更应运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

如果存在 change spec，还应运行：

```bash
node tools/vibit check change <change-id> --json
```

除非后续 work item 添加或改变 Go runtime code，否则这个 boundary-only standard 不要求运行 runtime Go tests。

## 11. Non-Goals

本标准不授权：

- Go application authentication service code。
- Runtime login behavior。
- Access-token generation。
- Access-token validation。
- Credential verifier comparison。
- Token verifier comparison。
- Logout execution。
- Refresh-token behavior。
- Cleanup jobs。
- Protobuf authentication messages。
- WebSocket proof carriers。
- WebSocket handshake authentication。
- Authentication dependencies。
- Authentication audit persistence。
- Runtime session persistence。
- 对 `authentication.Repository` 的变更。
- 对已 ratified migration schemas 的变更。
