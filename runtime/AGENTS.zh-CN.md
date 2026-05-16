# Go Runtime Agent 指南

状态：Draft v0.1
最后更新：2026-05-13
范围：`runtime/` Go server runtime workspace
权威来源：`../CONSTITUTION.md`、`../AGENTS.md` 和 `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
说明：本文件是 `runtime/AGENTS.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本指南适用于第一版 Go server runtime implementation。

## 1. 目的

`runtime/` 是 vibit 第一版 Go server runtime 的 Go module。

Go module path 是：

```text
github.com/iceiko/vibit/runtime
```

Runtime 的目标是通过一个小而长期维护的 backend slice 证明 vibit 的核心命题：

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

不要把这个 workspace 当成一次性 demo。

## 2. 必读内容

修改 `runtime/` 下文件前，先阅读：

- `../CONSTITUTION.md`
- `../AGENTS.md`
- `../.arch/runtime.yaml`
- `../.arch/dependencies.yaml`
- `../.arch/contracts.yaml`
- `../docs/generated-output.md`
- `../docs/runtime-protocol-adapter.md`
- persistence work 前阅读 `../docs/postgresql-persistence-boundary.md`
- live PostgreSQL verification work 前阅读 `../docs/postgresql-verification-environment.md`
- authentication、token、credential、external identity、session persistence、request identity trust、WebSocket handshake、player handler 或 player route work 前阅读 `../docs/authentication-token-session-validation.md`
- authentication proof、token/session validation、session error、session permission 或 validation event contract work 前阅读 `../docs/authentication-proof-token-session-contract-dimensions.md`
- `../docs/runtime-runbook.md`
- `../decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `../decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- persistence work 前阅读 `../decisions/ADR-0020-postgresql-persistence-boundary.md`
- player account persistence work 前阅读 `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`
- authentication/session design 或 implementation work 前阅读 `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`
- 受影响的 module manifest，例如 `../modules/inventory/module.yaml`
- `../changes/` 下相关 change spec

## 3. Package Ownership

使用以下 package boundaries：

- `cmd/vibit-server/`：process startup、configuration wiring 和 lifecycle。
- `internal/app/`：command/query dispatch、application composition 和 transaction orchestration。
- `internal/platform/transport/ws/`：WebSocket transport adapter 和 `github.com/coder/websocket` ownership。
- `internal/platform/protocol/protobuf/`：Protobuf framing、envelope conversion 和 wire message adaptation。
- `internal/platform/persistence/postgres/`：PostgreSQL adapter implementation 和 `github.com/jackc/pgx/v5` ownership。
- `internal/platform/migrations/`：migration tooling invocation 和 validation。
- `internal/platform/events/`：event recording 和 publication mechanisms。
- `internal/platform/tx/`：transaction boundary 和 unit-of-work interfaces。
- `internal/modules/<module>/`：手写 domain module runtime logic。
- `internal/generated/contracts/`：生成的 Go contract shapes。
- `internal/generated/proto/`：生成的 Go Protobuf files。
- `migrations/postgres/`：SQL-first PostgreSQL migration sources。

## 4. 依赖规则

Domain modules 不得直接 import 第三方 transport、protocol、persistence、migration、object-storage 或 framework dependencies。

允许的 owner packages：

- `github.com/coder/websocket`：仅 `internal/platform/transport/ws/`。
- `google.golang.org/protobuf`：仅 generated protocol packages 和 protocol adapter packages。
- `github.com/jackc/pgx/v5`：仅 `internal/platform/persistence/postgres/`。
- `github.com/pressly/goose/v3`：仅 `internal/platform/migrations/`。

未经检查 `../.arch/dependencies.yaml` 并创建所需 adoption record，不要添加新的 foundational dependencies。

## 5. Runtime Boundary Rules

Runtime protocol handoff 必须遵循 `../docs/runtime-protocol-adapter.md`。

WebSocket transport 读写 frames。Protobuf protocol adaptation 解码和编码 envelopes。Application dispatch 路由 commands 和 queries。Domain modules 执行 invariants。Generated packages 只提供 shapes。

WebSocket transport handlers 把 opaque frame bytes 交给注入的 protocol/application composition。它们不直接把 requests 适配成 commands 或 queries，也不得隐藏 business logic。

State-changing commands 应通过 `internal/app/` 进入，并在 application-owned unit of work 中运行。Repository mutations 和 domain event recording 应发生在同一个 unit of work 内。

当前 transaction skeleton 是 `internal/platform/tx.Runner`、`internal/platform/tx.UnitOfWork` 和 `internal/app.TransactionalDispatcher`。Application code 可以 import 这个 transaction boundary package，但不得 import persistence、migration、protocol 或 transport platform adapters。Query routes 默认应不经过 write unit of work。

Query handlers 不应改变状态，默认不需要 write transaction。

在 vibit 采纳明确的 event delivery 或 outbox standard 前，transaction 外的 event publication 继续 deferred。

PostgreSQL persistence work 必须遵循 `../docs/postgresql-persistence-boundary.md`。Repository interfaces 保持 module-owned，`pgx` 保持在 `internal/platform/persistence/postgres/` 下，`goose` 保持在 `internal/platform/migrations/` 下，SQL migration sources 保持在 `migrations/postgres/` 下。

第一版 durable inventory implementation 中，`GrantItem` 必须使用 transaction-bound repository，并在读取当前 items、执行 capacity-sensitive mutation 前调用 `LockInventoryForMutation`。返回的 `MutationLock` 是 locked aggregate view，不是 transaction owner。Repositories 不得在 command flows 中偷偷开启独立 write transactions。

第一版 PostgreSQL inventory repository adapter 是 `internal/platform/persistence/postgres/inventory_repository.go`。使用 `NewInventoryRepositoryForUnitOfWork` 构造它，并传入由 application-owned unit of work 提供的 executor，例如 `pgx.Tx` 或兼容的 test executor。该 adapter 不得调用 `BEGIN`、`COMMIT` 或 `ROLLBACK`；transaction lifetime 属于 `internal/platform/tx` 和 `internal/app`。

PostgreSQL configuration 由 `internal/platform/persistence/postgres/config.go` 拥有。它读取 `VIBIT_POSTGRES_DSN`、`VIBIT_POSTGRES_MAX_CONNS` 和 `VIBIT_POSTGRES_MIN_CONNS`，构建 pgx pool configuration，并且普通 unit tests 不得要求 live PostgreSQL server。Connection strings 和 credentials 必须来自 environment 或显式 runtime input，不得 commit。

pgx-backed transaction runner 是 `internal/platform/persistence/postgres/runner.go`。它实现 `internal/platform/tx.Runner`，同时把 pgx transaction handles 保持在 PostgreSQL platform package 内部。它会 commit 成功的 command unit of work，rollback 失败的 callback unit of work，并提供 package-owned helpers，例如 `UnitOfWork.NewInventoryRepository`，供未来 persistent composition 使用。不要从 `internal/app/` 或 domain modules import PostgreSQL runner；persistent runtime wiring 必须发生在已批准的 composition boundary 中。

`GrantItemMutation` 携带 `event_id`、`occurred_at` 和 `reason`，这样 PostgreSQL adapter 可以在与 item quantity update 相同的 executor path 中持久化 `inventory_item_grants`。

第一版 inventory migration source 是 `migrations/postgres/000001_create_inventory_state.sql`。它创建 `inventory_accounts`、`inventory_items` 和 `inventory_item_grants`。当 migration sources 或 migration guidance 发生变化时，运行 `node ../tools/vibit check migrations`。设置 `VIBIT_POSTGRES_TEST_DSN` 后，opt-in live durable inventory request-loop verification 会覆盖 migration status 和 apply behavior。

Ratified player account PostgreSQL schema boundary 已记录在 `../docs/postgresql-persistence-boundary.md` 和 `../decisions/ADR-0022-player-account-postgresql-schema-boundary.md`。第一版 player account migration source 是 `migrations/postgres/000002_create_player_account_state.sql`。该 migration 只创建 `player_accounts` 和 `player_account_events` lifecycle state。它不得添加 credentials、password hashes、external identity links、access tokens、refresh tokens、runtime session rows、WebSocket connection state、request identity validation results、inventory state 或 permission grants。

Player account repository interface boundary 是 `internal/modules/player/repository.go`。它是 storage-neutral domain code，可以定义 account lifecycle structs、`Repository.CreatePlayerAccount`、`Repository.GetPlayerAccount`，以及 persistence adapters 所需的 durable mutation metadata。PostgreSQL adapter 是 `internal/platform/persistence/postgres/player_account_repository.go`，focused tests 位于 `internal/platform/persistence/postgres/player_account_repository_test.go`。它使用 `NewPlayerAccountRepositoryForUnitOfWork(executor)`，实现 `player.Repository`，从 application-owned unit of work 接收 executor，并且不得调用 `BEGIN`、`COMMIT` 或 `ROLLBACK`。`UnitOfWork.NewPlayerAccountRepository` 是 PostgreSQL package helper，不得向 application 或 domain packages 暴露 pgx。

Player account PostgreSQL adapter 不授权 runtime handlers、WebSocket routes、authentication、token behavior、credential storage、external identity linking 或 session persistence。除非后续 change ratify 更多行为，adapter 只能写入 `player_accounts`，为 `PlayerAccountCreated` 写入 `player_account_events`，并从 `player_accounts` 读取当前 lifecycle rows。

Authentication、token 和 session validation design boundary 记录在 `../docs/authentication-token-session-validation.md` 和 `../decisions/ADR-0023-authentication-token-session-validation-design-boundary.md`。它分离 authentication proof、login methods、tokens、credentials、external identity links、runtime sessions、request identity、WebSocket handshake authentication、player account lifecycle、transport connection metadata 和 Protobuf envelope metadata。当前 `MetadataOnlySessionValidator` 是 non-authenticated bootstrap path。不要把 metadata-only `player_id`、`session_id` 或 `connection_id` 当作 production proof；未经单独 ratify，不要添加 authentication runtime code、token parsing、credential lookup、external identity linking、session persistence、Protobuf envelope authentication changes、WebSocket handshake authentication、runtime player handlers 或 WebSocket routes。`runtime.authentication_token_session_boundary` 是该边界的 repository check rule。

Authentication proof 与 token/session contract dimensions 记录在 `../docs/authentication-proof-token-session-contract-dimensions.md`。Actor kinds、validation statuses、proof statuses、failure classes、retryability、request identity handoff、session error metadata、session permission metadata 和 validation event metadata 应使用该标准。这些 dimensions 只是 semantic vocabulary，不授予实现 login methods、token formats、credential lookup、session persistence、Protobuf envelope changes、WebSocket handshake changes、runtime player handlers 或 WebSocket routes 的权限。

Credential storage 与 external identity linking boundaries 记录在 `../docs/credential-storage-external-identity-linking-boundaries.md`。在添加 credential storage、external identity linking、login methods、provider subjects、password hashing、OAuth、OIDC、provider SDKs、account linking、recovery flows、merge behavior 或相关 schema 前，应使用该标准。该边界保持 `player_accounts` 和 `player_account_events` 作为 lifecycle-only tables，并不授权 credential tables、external identity tables、provider dependencies、runtime lookup code、player lifecycle table changes 或 direct Nakama/Pitaya API compatibility。

Session persistence 与 WebSocket handshake decision gates 记录在 `../docs/session-persistence-websocket-handshake-decision-gates.md`。在添加 session persistence、WebSocket handshake authentication、reconnect behavior、connection epoch behavior、token/session carriers、session-related Protobuf envelope changes、handshake/system messages 或 route-level authentication 前，应使用该标准。它不选择 request-level、first-message、handshake-level、every-request 或 hybrid validation 作为 production model。它不授权 session tables、session store、envelope changes 或 handshake authentication behavior。

Login method 与 token format ratification 记录在 `../docs/login-method-token-format-ratification.md` 和 `../decisions/ADR-0024-login-method-token-format-ratification-boundary.md`。在选择第一批 login methods、token model、token format、proof carrier posture、token lifecycle semantics、credential/token/session schema gates 或 implementation queue 前，应使用该标准。它不授权 runtime authentication、token parsing、credential lookup、external identity linking、session persistence、Protobuf envelope changes、WebSocket handshake authentication、runtime player handlers 或 WebSocket routes。

Selected login/token boundary checks 记录在 `../docs/selected-login-token-boundary-checks.md` 和 `../decisions/ADR-0030-selected-login-token-boundary-checks.md`。已选第一版姿态是 `device_credential_login`、opaque high-entropy access tokens、login-command token issuance、explicit request proof payloads、第一版不支持 refresh token，并为 credential 与 token verifier records 使用 PostgreSQL schema gates。在已 ratified 的 credential 与 token verifier migration sources 和 metadata-only generated authentication contract shapes 之后，该姿态仍然 deferred runtime behavior。在未来 bounded implementation milestone 授权之前，不要添加 token validators、token issuers、credential repositories、token repositories、authentication Protobuf sources、additional authentication migrations、WebSocket proof carriers 或 authentication dependencies。`runtime.selected_login_token_boundary` 是该已选姿态的 repository check rule。

Credential record schema boundary 记录在 `../docs/credential-record-schema-boundary.md` 和 `../decisions/ADR-0032-credential-record-schema-boundary.md`。Credential migration source 现在位于 `migrations/postgres/000003_create_authentication_device_credentials.sql`，用于已 ratified 的 `authentication_device_credentials` semantics。除非后续 bounded work 明确授权 repositories、adapters、runtime lookup、handlers、routes、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior，否则 runtime code 必须保持不变。

Token verifier record schema boundary 记录在 `../docs/token-verifier-record-schema-boundary.md` 和 `../decisions/ADR-0033-token-verifier-record-schema-boundary.md`。Token verifier migration source 现在位于 `migrations/postgres/000004_create_authentication_access_tokens.sql`，用于已 ratified 的 `authentication_access_tokens` semantics。除非后续 bounded work 明确授权 repositories、adapters、runtime token issuance、validation、logout、refresh、cleanup、handlers、routes、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior，否则 runtime code 必须保持不变。

Authentication schema migration queue planning 记录在 `../docs/authentication-schema-migration-queue.md` 和 `../decisions/ADR-0034-authentication-schema-migration-queue.md`。Adapter implementation gate 记录在 `../decisions/ADR-0035-authentication-postgresql-adapter-implementation-gate.md`。Credential 与 token verifier migration sources、static checks、storage-neutral authentication repository interface boundary，以及 authentication PostgreSQL adapter boundary 现在都已存在。除非后续 bounded work 明确授权 runtime credential lookup、token issuance、token validation、logout execution、refresh、cleanup jobs、handlers、routes、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior，否则 runtime authentication code 必须保持不变。

`M-019 Token And Credential Verifier Algorithm Redaction Boundary` 已完成。`W-0091` 已在 `../docs/token-credential-verifier-algorithm-redaction-boundary.md` 和 `../decisions/ADR-0040-token-credential-verifier-algorithm-redaction-boundary.md` 中定义第一版 planned high-entropy verifier posture，且没有添加 runtime authentication behavior。第一版 planned verifier algorithm family 是 `vibit_hmac_sha256_v1`；未来 first-posture code 在后续 code gate 授权后可以使用 Go standard library packages `crypto/hmac`、`crypto/sha256`、`crypto/subtle`、`crypto/rand` 和 `encoding/base64`。Lookup digests 和 verifier digests 都不是 log-safe，verifier digest comparison 必须 constant-time，raw credential/token material 必须至少有 256 bits entropy。`runtime.token_credential_verifier_algorithm_redaction_boundary` 是已完成 verifier boundary 的 repository check rule。

`M-020 Secret Configuration And Verifier Key Loading Boundary` 已完成。`W-0092` 已在 `../docs/secret-configuration-verifier-key-loading-boundary.md` 和 `../decisions/ADR-0041-secret-configuration-verifier-key-loading-boundary.md` 中定义 future key loading posture，且没有添加 runtime authentication behavior。Future verifier key loading 由 `internal/app` 下的 application 拥有；第一版 local implementation 可在后续 code gate 授权后使用 process environment configuration 或 explicit runtime secret input；external KMS 或 secret-manager integration 仍然需要后续 dependency 和 operations gates。必须使用四个 separated logical verifier keys，`verifier_key_id` 默认不是 log-safe，production key configuration 无效时必须 fail closed，committed production-like secret values 仍然禁止。不要基于本 boundary 实现 secret loading、environment parsing、token generation、credential generation、verifier comparison、login、token validation、logout execution、refresh、cleanup jobs、Protobuf messages、WebSocket behavior、authentication dependencies、repository changes、migration schema changes、KMS integration、secret-manager integration 或 production authentication behavior。`runtime.secret_configuration_verifier_key_loading_boundary` 是该 boundary 的 repository check rule。

`M-021 Token And Credential Material Generation Boundary` 已完成。`W-0093` 已在 `../docs/token-credential-material-generation-boundary.md` 和 `../decisions/ADR-0042-token-credential-material-generation-boundary.md` 中定义 future raw device credential 与 opaque access-token material generation posture，且没有添加 runtime authentication behavior。Future material generation 由 `internal/app` 下的 application 拥有；第一版 device credential 与 access token 都是 server-issued 且 application-generated；raw material 必须是 32 cryptographically random bytes，至少 256 bits entropy；text presentation 使用 URL-safe unpadded Base64 或等价 encoding；raw material 只能 one-time client-visible，且不得存储。Repository handoff 保持 digest-only。Future first-posture generation helpers 可在后续 code gate 授权后使用 Go standard library `crypto/rand` 和 `encoding/base64`。不要基于本 boundary 实现 token generation、credential generation、secret loading、verifier digest computation、verifier comparison、login、token validation、logout execution、refresh、cleanup jobs、Protobuf messages、WebSocket behavior、authentication dependencies、repository changes、migration schema changes 或 production authentication behavior。`runtime.token_credential_material_generation_boundary` 是该 boundary 的 repository check rule。

`M-022 Verifier Digest Computation And Comparison Boundary` 已完成。`W-0094` 已在 `../docs/verifier-digest-computation-comparison-boundary.md` 和 `../decisions/ADR-0043-verifier-digest-computation-comparison-boundary.md` 中定义 future verifier digest computation 与 constant-time comparison posture，且没有添加 runtime authentication behavior。Future digest helpers 由 `internal/app` 下的 application 拥有；canonical HMAC input 带 version 且 length-prefixed；lookup digest equality 只能用于 record selection；validation 可跨 active 与 accepted previous key sets 计算 lookup candidates；stored `verifier_key_id` 用于选择 verifier digest computation 的 verifier key set；verifier digest comparison 必须 constant-time；invalid lookup、mismatch、unknown key id、unsupported algorithm、malformed proof 以及 expired 或 revoked proof 都收敛到同一个 public invalid-proof class。Future first-posture helpers 可在后续 code gate 授权后使用 Go standard library `crypto/hmac`、`crypto/sha256` 和 `crypto/subtle`。不要基于本 boundary 实现 verifier digest computation、verifier comparison、token generation、credential generation、secret loading、login、token validation、logout execution、refresh、cleanup jobs、Protobuf messages、WebSocket behavior、authentication dependencies、repository changes、migration schema changes 或 production authentication behavior。`runtime.verifier_digest_computation_comparison_boundary` 是该 boundary 的 repository check rule。

`M-023 Authentication Service Implementation Readiness Gate` 已完成。`W-0095` 已在 `../docs/authentication-service-implementation-readiness-gate.md` 和 `../decisions/ADR-0044-authentication-service-implementation-readiness-gate.md` 中定义 readiness gate，且没有添加 runtime authentication behavior。Future service implementation 仍由 `internal/app` 下的 application 拥有，package candidate 是 `internal/app/authentication`。第一段 code slice 必须单独授权；recommended next gate 是 local verifier key configuration loading。本 gate 定义 prior boundaries、allowed 与 forbidden write areas、test classes、redaction expectations、Nakama/Pitaya capability mapping 和 deferrals。不要基于本 gate 实现 authentication service code、secret loading、token generation、credential generation、verifier digest computation、verifier comparison、login、token validation、logout execution、refresh、cleanup jobs、Protobuf messages、WebSocket behavior、authentication dependencies、repository changes、migration schema changes 或 production authentication behavior。`runtime.authentication_service_implementation_readiness_gate` 是该 gate 的 repository check rule。

`M-024 Local Verifier Key Configuration Loading Gate` 已完成。`W-0096` 已在 `../docs/local-verifier-key-configuration-loading-gate.md` 和 `../decisions/ADR-0045-local-verifier-key-configuration-loading-gate.md` 中定义该 gate，且没有添加 runtime authentication behavior。`W-0097` 已在 `internal/app/authentication` 下实现 explicit in-memory verifier key set validator，并覆盖 accepted input、copying、immutable accessors、invalid key sets 和 redacted errors 的测试。Environment parsing 和 Base64 text decoding 推迟到 `W-0098`，也就是 environment verifier key loader gate。不要基于本 gate parse process environment variables、读取 local secret files、wire startup、集成 KMS 或 cloud secret managers、generate tokens 或 credentials、compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、添加 authentication dependencies 或添加 production authentication behavior。`runtime.local_verifier_key_configuration_loading_gate` 是该 gate 的 repository check rule。

`M-026 Environment Verifier Key Loader Gate` 已完成。`W-0098` 已在 `../docs/environment-verifier-key-loader-gate.md` 和 `../decisions/ADR-0046-environment-verifier-key-loader-gate.md` 中定义该 gate，且没有添加 Go code 或 runtime authentication behavior。Future process environment loader work 属于 `internal/app/authentication`，应使用 `verifier_key_env.go` 和 `verifier_key_env_test.go`，必须 decode required environment key text 并调用 `NewVerifierKeySet`，并保持 environment variable values 和 full concrete key set ids redacted。不要基于本 gate wire startup、读取 local secret files、parse `.env` files、接受 CLI secret input、集成 KMS 或 cloud secret managers、generate tokens 或 credentials、compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、添加 authentication dependencies 或添加 production authentication behavior。`runtime.environment_verifier_key_loader_gate` 是该 gate 的 repository check rule。

`M-027 Environment Verifier Key Loader Implementation` 已完成。`W-0099` 已实现 `internal/app/authentication/verifier_key_env.go` 和 `internal/app/authentication/verifier_key_env_test.go`。该 loader 接收 explicit lookup function，提供很小的 `os.LookupEnv` process adapter，decode Base64URL unpadded 与 standard padded Base64 key text，调用 `NewVerifierKeySet`，并返回 redacted typed errors。没有后续 bounded work item 授权时，不得把该 loader wire into process startup、读取 local secret files、parse `.env` files、接受 CLI secret input、集成 KMS 或 cloud secret managers、generate tokens 或 credentials、compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、添加 authentication dependencies 或添加 production authentication behavior。

`M-028 Token And Credential Material Generation Implementation Gate` 已完成。`W-0100` 已在 `../docs/token-credential-material-generation-implementation-gate.md` 和 `../decisions/ADR-0047-token-credential-material-generation-implementation-gate.md` 中定义该 gate，且没有添加 Go code 或 runtime authentication behavior。Future helper work 属于 `internal/app/authentication`，应使用 `material_generation.go` 和 `material_generation_test.go`，必须使用 32 random bytes、URL-safe unpadded Base64 presentation、explicit `io.Reader` entropy-source handoff、copied raw bytes、redacted errors 和 focused tests。不要基于本 gate compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。`runtime.token_credential_material_generation_implementation_gate` 是该 gate 的 repository check rule。

`M-029 Token And Credential Material Generation Helper Implementation` 已完成。`W-0101` 已实现 `internal/app/authentication/material_generation.go` 和 `internal/app/authentication/material_generation_test.go`。这些 helpers 暴露 `GenerateDeviceCredentialMaterial` 和 `GenerateAccessTokenMaterial`，接收 explicit `io.Reader`，使用 `io.ReadFull` 读取 32 bytes，编码 URL-safe unpadded Base64 text，保留 `MaterialKind`，返回 raw bytes copy，拒绝 nil readers、read errors、short reads、all-zero material 和 repeated-single-byte material，并返回 redacted typed errors。不要把这些 helpers 扩展为 compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。

`M-031 Verifier Digest Computation Helper Implementation` 已完成。`W-0103` 已实现 `internal/app/authentication/verifier_digest.go` 和 `internal/app/authentication/verifier_digest_test.go`。这些 helpers 暴露 `ComputeCredentialLookupDigest`、`ComputeCredentialVerifierDigest`、`ComputeTokenLookupDigest` 和 `ComputeTokenVerifierDigest`，构造 canonical `vibit.auth.verifier.input.v1` HMAC input，使用 `VerifierKeySet` 中匹配的 logical key，通过 `ComputedDigest` 返回 copied 32-byte digest bytes，并使用 redacted typed errors。不要把这些 helpers 扩展为 compare verifier digests、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。下一个 next ready work item 是 verifier digest comparison helper gate。

`M-032 Verifier Digest Comparison Helper Gate` 已完成。`W-0104` 已在 `../docs/verifier-digest-comparison-helper-gate.md` 和 `../decisions/ADR-0049-verifier-digest-comparison-helper-gate.md` 中定义该 gate，且没有添加 Go code 或 runtime authentication behavior。Future comparison helper work 属于 `internal/app/authentication`，应使用 `verifier_comparison.go` 和 `verifier_comparison_test.go`，必须保持 `verifier_digest.go` computation-only，只比较 computed verifier digest bytes 和 stored verifier digest bytes，首选 `crypto/hmac.Equal`，仅允许 `crypto/subtle.ConstantTimeCompare` 作为等价 constant-time alternative，拒绝 lookup digest classes 和 malformed input，并避免在 comparison 中使用 raw material、lookup digests、key ids、database-only equality、protocol metadata 和 public failure details。不要从本 gate 实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。`runtime.verifier_digest_comparison_helper_gate` 是该 gate 的 repository check rule。

`M-033 Verifier Digest Comparison Helper Implementation` 已完成。`W-0105` 已实现 `internal/app/authentication/verifier_comparison.go` 和 `internal/app/authentication/verifier_comparison_test.go`。这些 helpers 暴露 `CompareCredentialVerifierDigest` 和 `CompareTokenVerifierDigest`，返回 `VerifierComparisonResult`，使用 `crypto/hmac.Equal`，只比较 `ComputedDigest` verifier bytes 和 stored verifier digest bytes，拒绝 lookup digest classes、wrong verifier classes、missing input、malformed computed input 和 malformed stored input，并保持 errors redacted。没有后续 bounded work item 授权时，不要把这些 helpers 扩展为 authentication service behavior、login execution、token validation、logout、refresh、cleanup、protocol carriers、repository calls、migration changes、startup wiring、authentication dependencies 或 production authentication behavior。下一个 next ready work item 是 authentication service behavior implementation gate。

`M-034 Authentication Service Behavior Implementation Gate` 已完成。`W-0106` 已在 `../docs/authentication-service-behavior-implementation-gate.md` 和 `../decisions/ADR-0050-authentication-service-behavior-implementation-gate.md` 中定义该 gate，且没有添加 service code 或 runtime authentication behavior。Future service behavior 仍由 `internal/app/authentication` 下的 application 拥有，必须只通过 application unit-of-work boundary 使用 `authentication.Repository`，必须按固定顺序组合 verifier key、material generation、digest computation 和 comparison helpers，必须折叠 public proof failures，并且必须在 production-sensitive domain dispatch 前把 validated token proof 转换为 `RequestIdentity`。`runtime.authentication_service_behavior_implementation_gate` 是该 gate 和 skeleton boundary 的 repository check rule。

`M-035 Authentication Service Behavior Skeleton` 已完成。`W-0107` 添加了 `internal/app/authentication/service.go` 和 `internal/app/authentication/service_test.go`，作为 skeleton-only service shape。该 skeleton 定义 `Service`、`ServiceDependencies`、`UnitOfWorkRunner`、`ServiceError`、reserved authentication operations 的 request/result vocabulary、public error codes 和 redacted internal failure classes。`AuthenticateWithDeviceCredential`、`ValidateAccessToken` 与 `LogoutAccessToken` fail closed 为 `AUTHENTICATION_NOT_IMPLEMENTED`；`RefreshAccessToken` fail closed 为 `AUTHENTICATION_REFRESH_NOT_SUPPORTED`。没有后续 bounded work item 授权时，不要把该 skeleton 扩展为执行 login、validate access tokens、issue tokens、调用 repositories、revoke tokens、refresh tokens、暴露 protocol carriers、wire startup、添加 dependencies、修改 migrations 或添加 production authentication behavior。下一个 next ready work item 是 device credential login service behavior gate。

已实现的 authentication PostgreSQL adapter boundary 是 `internal/platform/persistence/postgres/authentication_repository.go`，focused tests 位于 `internal/platform/persistence/postgres/authentication_repository_test.go`。它的 constructor 是 `NewAuthenticationRepositoryForUnitOfWork(executor)`，`UnitOfWork.NewAuthenticationRepository` 会从 caller-owned executor 创建 `authentication.Repository`。该 adapter 仍然只能是 persistence-adapter-only。

Runtime authentication implementation boundary planning 记录在 `../docs/runtime-authentication-implementation-boundary.md` 和 `../decisions/ADR-0036-runtime-authentication-implementation-boundary.md`。未来 runtime authentication 由 `internal/app` 下的 application boundary 拥有；它必须通过 application unit-of-work boundary 使用 `authentication.Repository`，保持 PostgreSQL adapter 只负责 persistence，并在 domain dispatch 前把 validated proof 转换为 `RequestIdentity`。Token generation、verifier comparison、login execution、access-token validation、logout execution、cleanup jobs、Protobuf authentication messages、WebSocket proof carriers、generated authentication shapes 和 authentication dependencies 仍然是单独 gates。`runtime.authentication_implementation_boundary` 是该边界的 repository check rule。

Authentication generated contract shape timing 记录在 `../docs/authentication-generated-contract-shape-timing.md` 和 `../decisions/ADR-0038-authentication-generated-contract-shape-timing.md`。Source 是 `contracts/runtime/authentication/`，output root 是 `runtime/internal/generated/contracts/runtime/authentication/`。Generated authentication shapes 是 metadata-only 且 immutable；service interfaces 和 runtime behavior 仍然单独 gated。

Application authentication service interface boundary 记录在 `../docs/application-authentication-service-interface-boundary.md` 和 `../decisions/ADR-0039-application-authentication-service-interface-boundary.md`。未来 authentication service interfaces 由 `internal/app` 下的 application 拥有；generated authentication shapes 用于指导 service-level request/result vocabulary；service behavior 只能通过 application unit-of-work boundary 使用 `authentication.Repository`；validated proof 必须在 domain dispatch 前转换为 `RequestIdentity`。该 boundary 不授权 Go service code 或 runtime authentication behavior。`runtime.application_authentication_service_interface_boundary` 是该 boundary 的 repository check rule。

`../decisions/ADR-0037-close-runtime-auth-boundary-and-open-generated-shape-gate.md` 关闭 M-016 并打开 M-017。`ADR-0038` 完成 timing decision。`W-0089` 完成 generator/check support 加 metadata-only generated authentication shape output。`ADR-0039` 和 `W-0090` 完成 service-interface boundary step。`ADR-0040` 和 `W-0091` 完成 verifier algorithm/redaction step。`ADR-0041` 和 `W-0092` 完成 secret configuration/verifier key loading preparation step。`ADR-0042` 和 `W-0093` 完成 material generation preparation step。`ADR-0043` 和 `W-0094` 完成 verifier digest computation and comparison preparation step。`ADR-0044` 和 `W-0095` 完成 implementation readiness step。`ADR-0045` 和 `W-0096` 完成 local verifier key configuration loading gate。`W-0097` 完成 explicit in-memory verifier key set validator implementation slice。`ADR-0046` 和 `W-0098` 完成 environment verifier key loader gate。`W-0099` 完成 environment verifier key loader implementation slice。`ADR-0047` 和 `W-0100` 完成 token and credential material generation implementation gate。`W-0101` 完成 token and credential material generation helper implementation slice。`ADR-0048` 和 `W-0102` 完成 verifier digest helper implementation gate。`W-0103` 完成 verifier digest computation helper implementation slice。`ADR-0049` 和 `W-0104` 完成 verifier digest comparison helper gate。`W-0105` 完成 verifier digest comparison helper implementation slice。`ADR-0050` 和 `W-0106` 完成 authentication service behavior implementation gate。

第一版显式 PostgreSQL migration runner 是 `internal/platform/migrations/postgres.go`。它拥有 `github.com/pressly/goose/v3`，接收调用方提供的 `*sql.DB` 和 migration source filesystem 或 directory，列出 SQL migration sources，报告结构化 status，并且只在被显式调用时应用 pending migrations。未经 change spec 授权，不要把它接入普通 `cmd/vibit-server` startup。

Live PostgreSQL verification 受 `../docs/postgresql-verification-environment.md` 约束。它通过 `VIBIT_POSTGRES_TEST_DSN` 选择性启用；普通 unit tests、`node ../tools/vibit check runtime` 和默认 repository checks 不得要求运行中的 PostgreSQL server。当 live PostgreSQL check 因为没有 disposable DSN 而跳过时，必须显式记录。

## 6. Generated Files

Generated files 对 non-system agents 不可变。

如果 generated output 错了，应修改 source contract、schema、template 或 generator。除非 change spec 或 Agent Decision Record 明确授予 `generated_file_override`，不要手工编辑 generated files。

`internal/generated/proto/` 下的 Go Protobuf generated output 必须通过已接受的 Buf 和 `protoc-gen-go` 路径从 `../proto/` sources 生成。该目录下的文件必须是带有 `protoc-gen-go` marker 和 source trace 的 generated `*.pb.go` files，或者是在 generation 尚未运行时使用的临时 `.gitkeep` placeholders。

不要把 handwritten runtime code 放在 `internal/generated/proto/` 或 `internal/generated/contracts/` 下。

## 7. 当前状态

这个 runtime workspace 现在已经有第一批 generated Protobuf output、第一段窄 runtime handoff slice、第一版 WebSocket transport adapter、一个用于 command 和 query routes 的小型 application dispatch skeleton、第一版 transaction boundary skeleton、带 command-safe mutation lock 的第一版 inventory repository/policy/handler runtime boundary、第一版 PostgreSQL configuration parser、第一版 pgx-backed transaction runner adapter、第一版 PostgreSQL inventory repository adapter、第一条 inventory Protobuf/domain payload bridge、第一条 application-error-to-Protobuf-error-envelope mapper、第一版 frame-to-Protobuf-to-application composition adapter、用于 Protobuf command/query tests 的 package-local request-loop test fixture、挂载 `/v1/ws` 的 minimal process wiring、显式 PostgreSQL inventory runtime composition path，以及 opt-in live PostgreSQL durable inventory request-loop verification test。

这个 workspace 已经有 documented PostgreSQL persistence boundary、transaction skeleton、PostgreSQL configuration parser、pgx-backed transaction runner、第一版 inventory migration source、第一版显式 migration apply/status runner、第一版 PostgreSQL repository adapter、显式 runtime store selection、已添加第一版 migration source、storage-neutral repository interface、focused PostgreSQL adapter implementation 和 PostgreSQL unit-of-work factory helper 的 ratified player account PostgreSQL lifecycle schema boundary，已 ratified 的 authentication credential 与 token verifier migration sources 及 static checks，storage-neutral authentication repository interface boundary，已实现的 authentication PostgreSQL adapter，已经 ratify 但尚未实现 runtime authentication 的 authentication/token/session validation design boundary，runtime authentication implementation boundary planning standard，以及 metadata-only generated authentication contract shapes，并有只有设置 `VIBIT_POSTGRES_TEST_DSN` 才会运行 live branch 的 live verification test。`VIBIT_RUNTIME_STORE=memory` 仍然是默认值。提供 `VIBIT_POSTGRES_DSN` 时，`VIBIT_RUNTIME_STORE=postgres` 会启用 PostgreSQL-backed inventory composition。这个 workspace 仍然没有实现 generated route registration、generated protocol bridge creation、production authentication/session validation、runtime player account handlers、WebSocket player routes、automatic startup migrations 或 catalog-driven error retryability。

第一版手动 process run path 是：

```bash
cd runtime
go run ./cmd/vibit-server
```

第一版显式 persistent process run path 是：

```bash
cd runtime
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

普通 server startup 不会自动 apply migrations。

第一条 opt-in live durable inventory verification command 是：

```bash
cd runtime
VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

如果未设置 `VIBIT_POSTGRES_TEST_DSN`，该 test 会 skip，并记录 live PostgreSQL verification 不可用。第一版 live execution 已经在 local Termux PostgreSQL 18.2 上通过。

## 8. 验证

从仓库根目录运行 repository verification：

```bash
node tools/vibit check runtime
node tools/vibit check generated
node tools/vibit check migrations
node tools/vibit check postgres-env
node tools/vibit check all
```

当 Go source files 存在且本地 Go toolchain 可用时，runtime verification 应包括：

```bash
go test ./...
go vet ./...
```

当 Go toolchain 不可用或测试没有运行时，不要声称已经完成 Go test verification。
