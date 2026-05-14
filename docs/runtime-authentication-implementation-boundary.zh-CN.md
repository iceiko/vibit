# Runtime Authentication Implementation Boundary 中文版

状态：Draft v0.1
最后更新：2026-05-14
范围：authentication PostgreSQL adapter 之后的第一版 runtime authentication implementation boundary planning
依赖：`docs/authentication-token-session-validation.md`、`docs/login-method-token-format-ratification.md`、`docs/selected-login-token-boundary-checks.md`、`docs/postgresql-persistence-boundary.md`
权威决策：`ADR-0036`
说明：本文件是 `docs/runtime-authentication-implementation-boundary.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

## 1. 目的

本文定义 vibit 在已经完成以下基础之后，未来 runtime authentication implementation 的边界：

- 已 ratify 第一种 login method：`device_credential_login`。
- 已 ratify opaque high-entropy access tokens。
- 已 ratify explicit request proof payload 作为第一版 request proof posture。
- 已添加 credential 与 token verifier 的 PostgreSQL migration sources。
- 已添加 storage-neutral `authentication.Repository` interface。
- 已实现 authentication PostgreSQL persistence adapter。

本文是 planning 与 boundary standard。它不实现 login、token generation、token validation、verifier comparison、logout execution、refresh、cleanup jobs、Protobuf messages、WebSocket proof carriers、generated authentication shapes、authentication dependencies 或 production authentication behavior。

## 2. 核心规则

Runtime authentication 由 application 拥有。

第一版实现必须在 domain module 收到 authenticated identity 之前通过 application boundary：

```text
protocol decoded request
-> application authentication boundary
-> module-owned authentication repository interface
-> platform-owned PostgreSQL adapter through an application unit of work
-> application-owned request identity handoff
-> domain dispatch
```

Transport、protocol adapter、domain module、player repository、generated file 或 PostgreSQL adapter 都不能成为 authentication proof、token validation、verifier comparison 或 permission decisions 的所有者。

当前 runtime 仍然只有 metadata-only identity。Metadata-only `player_id`、`session_id` 和 `connection_id` 仍然不是 authenticated proof。

## 3. Ownership Split

### Application Service Boundary

Owner：

```text
runtime/internal/app
```

未来在 bounded implementation work item 授权代码之后的职责：

- 在 protocol decoding 之后 orchestrate device credential login。
- 在 production-sensitive domain dispatch 之前 orchestrate access-token validation。
- Orchestrate presented-token logout。
- 把 authentication outcomes 转换为 `RequestIdentity`。
- 把 authentication failures 映射为已注册 application errors。
- 通过 module-owned `authentication.Repository` 请求 repository operations。
- 让会改变状态的 authentication operations 在 application unit of work 内运行。

不得：

- import WebSocket transport packages。
- import generated Protobuf packages。
- 存储 raw token 或 credential material。
- 拥有 PostgreSQL SQL text 或 driver handles。
- 把 verifier comparison 藏在 generic dispatch 中。
- 把 metadata-only identity 当作 proof。

### Authentication Module Boundary

Owner：

```text
runtime/internal/modules/authentication
```

当前职责：

- 拥有 storage-neutral authentication repository interfaces 和 record shapes。
- 拥有 credential 与 token verifier records 的封闭 status vocabulary。
- 维护 repository inputs 和 outputs 的 copying、normalization 与 UTC timestamp invariants。

未来职责：

- 当 runtime authentication behavior 开始时，继续拥有 storage-neutral repository interfaces。

不得：

- 生成 tokens。
- 比较 credential 或 token verifiers。
- 解析 bearer tokens。
- 执行 login、logout 或 cleanup jobs。
- import PostgreSQL、WebSocket、Protobuf、JWT、OAuth、OIDC、password-hashing、provider SDK 或 Redis-like dependencies。

### PostgreSQL Adapter Boundary

Owner：

```text
runtime/internal/platform/persistence/postgres
```

当前职责：

- 使用 caller-supplied executors 实现 `authentication.Repository`。
- 持久化并读取已 ratified 的 credential 与 token verifier records。
- 保持 no-transaction-control behavior。
- 把 pgx details 留在 platform package 内。

不得：

- 生成 raw credential 或 token material。
- 比较 verifier digests。
- 解释 token proof。
- 执行 authentication decisions。
- 直接发出 domain 或 audit events。
- 知晓 WebSocket、Protobuf、request identity 或 permission behavior。

### Protocol Adapter Boundary

Owner：

```text
runtime/internal/platform/protocol/protobuf
```

未来在 protocol gate 授权后的职责：

- 把已 ratified 的 Protobuf authentication request 和 response messages 转换为 application route requests。
- 把 application authentication errors 转换为 public-safe Protobuf error envelopes。

不得：

- 选择 proof carrier semantics。
- 在 application-owned boundaries 之外解析或验证 tokens。
- 把当前 `Session` metadata 当作 proof。
- 生成 tokens 或比较 verifiers。

### WebSocket Transport Boundary

Owner：

```text
runtime/internal/platform/transport/ws
```

当前职责：

- 接受 connections。
- 读取和写入 opaque binary frames。
- 把 frame bytes 委托给 injected handlers。

不得：

- 在后续 WebSocket proof-carrier decision 授权之前，把 `Authorization`、`Bearer`、`Cookie` 或 `Sec-WebSocket-Protocol` 当作 authentication proof 读取。
- 解析 credentials 或 tokens。
- 通过 connection metadata 绑定 player identity。
- 拥有 domain permission decisions。

## 4. Token Material Lifecycle Placeholder

Opaque access-token generation 仍然是单独 gate。

未来 token generation 必须定义：

- Raw token generation owner。
- Entropy source 和 minimum entropy verification。
- Text encoding。
- Redaction tests。
- One-time client presentation rules。
- Verifier digest derivation。
- 通过 `authentication.Repository` 的 storage behavior。
- Rotation 和 revocation behavior。

在该 gate 存在之前：

- 任何代码不得为 authentication tokens 调用 `crypto/rand` 或等价机制。
- raw token 不得出现在 logs、tests、change specs、conversation logs 或 database rows 中。
- token string 不得复制到 Protobuf `Session` metadata 中。

## 5. Verifier Comparison Boundary

Verifier comparison 与 repository persistence 是单独 gate。

未来 verifier comparison 必须定义：

- device credential proof 使用哪种 verifier algorithm。
- access-token proof 使用哪种 verifier algorithm。
- Go standard library 是否足够，或是否需要 external dependency adoption record。
- constant-time comparison requirements。
- invalid、expired、revoked、unavailable、malformed proof 的 error mapping。
- 能证明 raw secrets 不被存储或回显的 tests。

Repositories 可以存储和读取 verifier digests。它们不得决定 presented proof 是否有效。

## 6. Request Identity Handoff

未来 token validation 必须在 domain dispatch 前产生 application-owned request identity。

目标 handoff：

```yaml
owner: runtime/internal/app
input: explicit_request_proof_payload_after_protocol_decode
output: RequestIdentity
required_success_markers:
  actor_kind: player
  validation_status: authentication_proven_or_session_validated
  player_id_validated: true
  session_validated: false_until_session_persistence_is_ratified
metadata_only_allowed_as_proof: false
```

规则：

- Domain modules 消费 `RequestIdentity`；它们不验证 tokens。
- 只有 validation boundary 存在之后，Inventory permission policy 才能使用 validated identity。
- 当前 `MetadataOnlySessionValidator` 仍然只是 bootstrap path，不授予 production player-owned permissions。
- Runtime session persistence 仍然 deferred。

## 7. Error、Permission 与 Audit Mapping

Runtime authentication implementation 必须使用以下目录下已注册的 semantic surfaces：

```text
contracts/runtime/authentication/
```

第一版必须映射：

- `AuthenticateWithDeviceCredential` failures 映射到 `authentication_errors`。
- `ValidateAccessToken` failures 映射到 missing、malformed、invalid、expired、revoked、unsupported、unavailable 或 not-implemented failure classes。
- `LogoutAccessToken` 必须区分 presented-token revocation success、missing proof、invalid proof、expired proof、revoked proof 和 store unavailability。
- 如果暴露 `RefreshAccessToken`，第一版必须保持 unsupported，并映射到 `AUTHENTICATION_REFRESH_NOT_SUPPORTED`。
- raw credential 与 raw token material 绝不能出现在 public errors 中。

Permission surfaces 在实现前仍然只是 semantic：

- `authentication_device_credential_authenticate`
- `authentication_access_token_validate`
- `authentication_access_token_logout`
- `authentication_access_token_refresh`

Audit event publication 与 audit persistence 仍然是单独 gates。Runtime authentication code 只有在 event publication 和 storage path 已 ratify 后，才可以准备 event facts。

## 8. Implementation Queue

第一版 runtime authentication implementation queue 必须保持为可单独 review 的 gates。

推荐顺序：

1. 添加或细化 runtime authentication boundary checks。
2. 决定 runtime code 之前是否需要 generated Go authentication contract shapes。
3. Ratify token 与 credential verifier algorithms、redaction tests 和 dependency posture。
4. 添加 application-owned authentication service interfaces 和 tests，但不添加 protocol 或 transport behavior。
5. 只在已 ratified 的 application boundary 内实现 token material generation。
6. 只在已 ratified 的 application boundary 内实现 credential verifier comparison。
7. 实现 `AuthenticateWithDeviceCredential` execution。
8. 实现 `ValidateAccessToken` execution 和 request identity handoff。
9. 为 presented access token 实现 `LogoutAccessToken` execution。
10. 定义并实现 cleanup job behavior。
11. 在 protocol impact 已 ratify 后添加 Protobuf authentication messages 和 bridge code。
12. 只有在后续 decision 选择 WebSocket carrier 后，才添加 WebSocket proof-carrier behavior。

后续 ADR 可以修改这个顺序，但任何 gate 都不得静默合并到另一个 gate 中。

## 9. Nakama 与 Pitaya Alignment

Nakama 仍然是以下能力的参考：

- Account authentication。
- Session token issuance。
- Token expiration。
- Token revocation 与 logout。
- Realtime socket 绑定 authenticated actors。

Pitaya 仍然是以下词汇的参考：

- 传递给 handlers 的 session context。
- Frontend acceptor 与 backend handler separation。
- Route handler identity context。
- Session binding concepts。

vibit 会把这些概念改造成自己的边界。它不复制 Nakama 或 Pitaya 的 public API shapes，不默认把 authentication 放进 WebSocket acceptors，也不让 route handlers 直接验证 tokens。

## 10. Verification Path

该边界的 repository check rule 是：

```text
runtime.authentication_implementation_boundary
```

触及此边界的变更应运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check architecture --json
node tools/vibit check all --json
```

如果存在 change spec，还应运行：

```bash
node tools/vibit check change <change-id> --json
```

Implementation changes 还必须运行相关 Go tests。除非某次变更明确添加或修改无法静态验证的 database behavior，否则本 boundary standard 不要求 live PostgreSQL。

## 11. Non-Goals

本标准不授权：

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
- Generated authentication shapes。
- Authentication dependencies。
- Authentication audit persistence。
- Runtime session persistence。
- 修改 `authentication.Repository`。
- 修改已 ratified 的 migration schemas。
