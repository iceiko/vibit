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

`M-019 Token And Credential Verifier Algorithm Redaction Boundary` 已完成。`W-0091` 已在 `docs/token-credential-verifier-algorithm-redaction-boundary.md` 和 `ADR-0040` 中定义第一版 planned high-entropy verifier posture，且没有添加 service code 或 runtime authentication behavior。更早的 application authentication service boundary 仍记录在 `docs/application-authentication-service-interface-boundary.md` 和 `ADR-0039`。第一版 planned verifier algorithm family 是 `vibit_hmac_sha256_v1`；未来 first-posture code 在后续 code gate 授权后可以使用 Go standard library packages `crypto/hmac`、`crypto/sha256`、`crypto/subtle`、`crypto/rand` 和 `encoding/base64`。本模块的 `authentication.Repository` 只能存取 already-computed digest material；它不得计算 HMAC、生成 token 或 credential material、比较 verifiers、加载 secret keys 或决定 authentication outcomes。

`M-020 Secret Configuration And Verifier Key Loading Boundary` 已完成。`W-0092` 已在 `docs/secret-configuration-verifier-key-loading-boundary.md` 和 `ADR-0041` 中定义 future key loading posture，且没有添加 service code 或 runtime authentication behavior。Future verifier key loading 由 `runtime/internal/app` 下的 application 拥有；第一版 local implementation 可在后续 code gate 授权后使用 process environment configuration 或 explicit runtime secret input；external KMS 或 secret-manager integration 仍然需要后续 dependency 和 operations gates。必须使用四个 separated logical verifier keys，`verifier_key_id` 默认不是 log-safe，production key configuration 无效时必须 fail closed，committed production-like secret values 仍然禁止。本模块不得 load verifier keys、parse environment variables、choose secret managers、rotate keys、compute HMACs、generate token 或 credential material、compare verifiers，或决定 authentication outcomes。

`M-021 Token And Credential Material Generation Boundary` 已完成。`W-0093` 已在 `docs/token-credential-material-generation-boundary.md` 和 `ADR-0042` 中定义 future raw device credential 与 opaque access-token material generation posture，且没有添加 service code 或 runtime authentication behavior。Future material generation 由 `runtime/internal/app` 下的 application 拥有；第一版 device credential 与 access token 都是 server-issued 且 application-generated；raw material 必须是 32 cryptographically random bytes，至少 256 bits entropy；text presentation 使用 URL-safe unpadded Base64 或等价 encoding；raw material 只能 one-time client-visible，且不得存储。本模块的 `authentication.Repository` 只能存取 already-computed digest material；不得 generate raw token 或 credential material、encode raw material、接收 raw material 用于 storage、compute digests、compare verifiers，或决定 authentication outcomes。`runtime.token_credential_material_generation_boundary` 是该 boundary 的 repository check rule。

`M-022 Verifier Digest Computation And Comparison Boundary` 已完成。`W-0094` 已在 `docs/verifier-digest-computation-comparison-boundary.md` 和 `ADR-0043` 中定义 future verifier digest computation 与 constant-time comparison posture，且没有添加 service code 或 runtime authentication behavior。Future digest computation and comparison 由 `runtime/internal/app` 下的 application 拥有；lookup digest equality 只能选择 candidate record；verifier digest comparison 必须 constant-time；invalid lookup、mismatch、unknown key id、unsupported algorithm、malformed proof 以及 expired 或 revoked proof 必须收敛到同一个 public invalid-proof class。本模块的 `authentication.Repository` 只能存取 already-computed digest material；不得 compute HMACs、choose verifier key sets、compare verifier digests、disclose lookup misses，或决定 authentication outcomes。`runtime.verifier_digest_computation_comparison_boundary` 是该 boundary 的 repository check rule。

`M-023 Authentication Service Implementation Readiness Gate` 已完成。`W-0095` 已在 `docs/authentication-service-implementation-readiness-gate.md` 和 `ADR-0044` 中定义 readiness gate，且没有添加 service code 或 runtime authentication behavior。Future service implementation 仍由 `runtime/internal/app` 下的 application 拥有，package candidate 是 `runtime/internal/app/authentication`。本模块仍是 storage-neutral repository boundary。它不得吸收 service orchestration、secret loading、material generation、digest computation、verifier comparison、login execution、token validation、protocol behavior、WebSocket behavior 或 production authentication decisions。`runtime.authentication_service_implementation_readiness_gate` 是该 gate 的 repository check rule。

`M-024 Local Verifier Key Configuration Loading Gate` 已完成。`W-0096` 已在 `docs/local-verifier-key-configuration-loading-gate.md` 和 `ADR-0045` 中定义该 gate，且没有添加 service code 或 runtime authentication behavior。`W-0097` 已在 `runtime/internal/app/authentication` 下实现 explicit in-memory verifier key set validation，不属于本模块。本模块仍是 storage-neutral，只能存储 already-computed verifier metadata。它不得 load verifier keys、parse environment variables、decode key text、hold key material、validate key sets、rotate keys、compute digests、compare verifiers，或决定 authentication outcomes。`runtime.local_verifier_key_configuration_loading_gate` 是该 gate 的 repository check rule。

`M-026 Environment Verifier Key Loader Gate` 已完成。`W-0098` 已在 `docs/environment-verifier-key-loader-gate.md` 和 `ADR-0046` 中定义该 gate，且没有添加 service code 或 runtime authentication behavior。Future process environment verifier key loading 属于 `runtime/internal/app/authentication`，不属于本模块，并且必须在该 package 调用 `NewVerifierKeySet`。本模块仍是 storage-neutral，只能存储 already-computed verifier metadata。它不得 parse environment variables、decode key text、hold key material、validate key sets、读取 local secret files、parse `.env` files、rotate keys、compute digests、compare verifiers，或决定 authentication outcomes。`runtime.environment_verifier_key_loader_gate` 是该 gate 的 repository check rule。

`M-027 Environment Verifier Key Loader Implementation` 已完成。`W-0099` 已在 `runtime/internal/app/authentication` 下实现 process environment verifier key loading，而不是在本模块中实现。本 authentication module 仍是 storage-neutral，不得 parse environment variables、decode key text、hold key material、validate key sets、读取 local secret files、parse `.env` files、rotate keys、compute digests、compare verifiers、wire startup 或决定 authentication outcomes。

`M-028 Token And Credential Material Generation Implementation Gate` 已完成。`W-0100` 已在 `docs/token-credential-material-generation-implementation-gate.md` 和 `ADR-0047` 中定义该 gate，且没有添加 service code 或 runtime authentication behavior。Future token and credential material helpers 属于 `runtime/internal/app/authentication`，不属于本模块。本模块仍是 storage-neutral，只能存储 already-computed digest material。它不得 generate raw token 或 credential material、encode raw material、hold raw material、accept raw material for storage、compute digests、compare verifiers、wire startup 或决定 authentication outcomes。`runtime.token_credential_material_generation_implementation_gate` 是该 gate 的 repository check rule。

`M-029 Token And Credential Material Generation Helper Implementation` 已完成。`W-0101` 已在 `runtime/internal/app/authentication` 下实现 token and credential material helpers，而不是在本模块中实现。本模块仍是 storage-neutral，只能存储 already-computed digest material。它不得 generate raw token 或 credential material、encode raw material、hold raw material、accept raw material for storage、compute digests、compare verifiers、wire startup 或决定 authentication outcomes。

`M-031 Verifier Digest Computation Helper Implementation` 已完成。`W-0103` 已在 `runtime/internal/app/authentication` 下实现 lookup and verifier digest computation helpers，而不是在本模块中实现。本模块仍是 storage-neutral，只能存储 already-computed digest material。它不得 compute HMACs、choose verifier key sets、compare verifier digests、generate raw token 或 credential material、hold raw material、wire startup 或决定 authentication outcomes。

`M-032 Verifier Digest Comparison Helper Gate` 已完成。`W-0104` 已在 `docs/verifier-digest-comparison-helper-gate.md` 和 `ADR-0049` 中定义 verifier digest comparison helper gate，且没有添加 service code 或 runtime authentication behavior。Future comparison helpers 属于 `runtime/internal/app/authentication`，不属于本模块，并应使用 `verifier_comparison.go` 和 `verifier_comparison_test.go`。本模块仍是 storage-neutral，只能存储 already-computed digest material。它不得 compare verifier digests、compare raw material、把 lookup digests comparison 当作 authentication proof、把 database-only equality 当作最终 proof、暴露 lookup misses、从 repository code 调用 comparison helpers，或决定 authentication outcomes。`runtime.verifier_digest_comparison_helper_gate` 是该 gate 的 repository check rule。

`M-033 Verifier Digest Comparison Helper Implementation` 已完成。`W-0105` 已在 `runtime/internal/app/authentication` 下实现 verifier comparison helpers，而不是在本模块中实现。本模块仍是 storage-neutral，只能存储 already-computed digest material。它不得 compare verifier digests、从 repository code 调用 comparison helpers、决定 mismatch/public failure mapping、执行 login、validate tokens、wire startup 或决定 authentication outcomes。

## Forbidden Shortcuts

- 不要存储 raw credential 或 token material。
- 不要在本模块中计算 HMAC、生成 verifier digests、加载 verifier keys 或比较 verifier material。
- 不要添加 `AuthService`、`Authenticator`、`TokenValidator`、`TokenIssuer`、`TokenVerifier`、`TokenRepository` 或 `CredentialRepository` implementation types。
- 不要在本模块下添加 PostgreSQL adapter files。
- 不要导入 `pgx`、`goose`、WebSocket、Protobuf、JWT、OAuth、OIDC、bcrypt、argon2、provider SDKs、Redis-like clients、S3 SDKs 或 MinIO clients。
- 编辑 repository interface boundary 时，不要添加 runtime authentication behavior。
- 未经单独 change spec 和 decision，不要改变已 ratified 的 credential 或 token verifier migration schemas。
- 不要让 player account lifecycle tables 拥有 authentication state。

## Required Tests

参见 `module.yaml` 中的 `tests.required`。

当前 module 状态下，focused tests 只限于 repository interface shape、closed status sets、required-field normalization、digest copying、timestamp UTC normalization，以及 storage-neutral mutation/query validation。PostgreSQL adapter tests 在 `M-015` 期间属于 `runtime/internal/platform/persistence/postgres/`；runtime authentication、WebSocket、Protobuf 和 generated-shape tests 只有在它们各自的 boundary 被单独 ratify 后才会成为必需。
