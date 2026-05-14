# authentication Module Agent Guide 中文版

状态：Draft v0.1
说明：本文件是 `modules/authentication/AGENTS.md` 的简体中文译本。英文版本是权威版本。

## 何时使用本模块

当需求定义第一版已选姿态的 storage-neutral authentication persistence boundaries 时，使用本模块：

- `device_credential_login`
- opaque high-entropy access tokens
- PostgreSQL-backed credential and token verifier records

当前 module-owned runtime boundary 是：

```text
runtime/internal/modules/authentication/repository.go
```

它可以定义 future persistence adapters 所需的 credential record structs、token verifier record structs、mutation/query shapes、repository interfaces 和 validation helpers。

## 何时不要使用本模块

不要用本模块实现：

- Runtime login handlers。
- Token generation、parsing、issuance、validation、refresh 或 bearer-token acceptance。
- 超出 storage-neutral revocation mutation shapes 的 logout behavior。
- 超出 storage-neutral cleanup query shapes 的 cleanup jobs。
- 超出单独授权的 `M-015` platform boundary 的 PostgreSQL adapters。
- WebSocket routes、proof carriers 或 handshake authentication。
- Protobuf messages，或 generated authentication shape paths 下的 handwritten authentication behavior。
- Password hashing、JWT、OAuth、OIDC、provider SDKs、key-management、Redis-like token stores、S3 或 MinIO dependencies。
- Player account lifecycle storage。

如果需求涉及上述 surface，应先创建或更新单独 ratified change，再添加代码。

## Extension Points

- Authentication repository interface：`runtime/internal/modules/authentication/repository.go`。
- Authentication repository boundary tests：`runtime/internal/modules/authentication/repository_test.go`。
- Authentication PostgreSQL adapter source：`runtime/internal/platform/persistence/postgres/authentication_repository.go`。
- Authentication PostgreSQL adapter tests：`runtime/internal/platform/persistence/postgres/authentication_repository_test.go`。
- 已存在的 credential migration source：`runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`。
- 已存在的 token verifier migration source：`runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`。

Repository interface 是 storage-neutral domain code。它不得导入 PostgreSQL、WebSocket、Protobuf、OAuth、OIDC、JWT、password-hashing、Redis-like、S3 或 MinIO packages。

Repository interface 可以 normalize identifiers、digest byte slices、statuses、timestamps 和 storage mutation/query shapes。它不得创建 credential material、generate tokens、compare verifiers、parse bearer tokens、validate access tokens、open transactions、execute SQL、publish events 或调用 transport/protocol code。

已实现的 PostgreSQL adapter boundary 是 `runtime/internal/platform/persistence/postgres/authentication_repository.go`，focused tests 位于 `runtime/internal/platform/persistence/postgres/authentication_repository_test.go`。它使用 `NewAuthenticationRepositoryForUnitOfWork(executor)` 并实现 `authentication.Repository`；`UnitOfWork.NewAuthenticationRepository` 会从 caller-owned executor 创建它。`M-015` 只授权 platform-owned persistence adapter work；它仍不授权 runtime authentication behavior。

Runtime authentication implementation boundary planning 记录在 `docs/runtime-authentication-implementation-boundary.md` 和 `decisions/ADR-0036-runtime-authentication-implementation-boundary.md`。未来 runtime authentication 必须由 `runtime/internal/app` 下的 application boundary 拥有，必须通过 application unit-of-work boundary 使用本模块的 `authentication.Repository`，并在 domain dispatch 前把 validated proof 转换为 `RequestIdentity`。本模块不得吸收 token generation、verifier comparison、login execution、access-token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、generated authentication shape paths 下的 handwritten logic 或 authentication dependencies。

`M-018 Application Authentication Service Interface Boundary` 已完成。`W-0090` 已在 `docs/application-authentication-service-interface-boundary.md` 和 `ADR-0039` 中定义未来 application-owned authentication service interface boundary，且没有添加 service code 或 runtime authentication behavior。未来 service behavior 只能通过 application unit-of-work boundary 使用本模块的 `authentication.Repository`。当前 authentication-adjacent milestone 是 `M-019 Token And Credential Verifier Algorithm Redaction Boundary`，`W-0091` 是下一个 ready work item。它可以定义 verifier algorithm 和 redaction test expectations，但不得实现 token generation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository changes、migration schema changes 或 production authentication behavior。

## Forbidden Shortcuts

- 不要存储 raw credential 或 token material。
- 不要添加 `AuthService`、`Authenticator`、`TokenValidator`、`TokenIssuer`、`TokenVerifier`、`TokenRepository` 或 `CredentialRepository` implementation types。
- 不要在本模块下添加 PostgreSQL adapter files。
- 不要导入 `pgx`、`goose`、WebSocket、Protobuf、JWT、OAuth、OIDC、bcrypt、argon2、provider SDKs、Redis-like clients、S3 SDKs 或 MinIO clients。
- 编辑 repository interface boundary 时，不要添加 runtime authentication behavior。
- 未经单独 change spec 和 decision，不要改变已 ratified 的 credential 或 token verifier migration schemas。
- 不要让 player account lifecycle tables 拥有 authentication state。

## Required Tests

参见 `module.yaml` 中的 `tests.required`。

当前 module 状态下，focused tests 只限于 repository interface shape、closed status sets、required-field normalization、digest copying、timestamp UTC normalization，以及 storage-neutral mutation/query validation。PostgreSQL adapter tests 在 `M-015` 期间属于 `runtime/internal/platform/persistence/postgres/`；runtime authentication、WebSocket、Protobuf 和 generated-shape tests 只有在它们各自的 boundary 被单独 ratify 后才会成为必需。
