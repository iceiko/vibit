# Agent 操作指南

状态：Draft v0.1  
最后更新：2026-05-13
范围：仓库级 coding agent 操作指令  
权威来源：`CONSTITUTION.md`  
说明：本文件是 `AGENTS.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

这份指南把宪法转化为 agent 的工作规则。它不取代宪法。如果本指南与 `CONSTITUTION.md` 冲突，应遵循 `CONSTITUTION.md`，并更新本指南。

配套英文源文档是 `AGENTS.md`。英文文件是权威版本。

## 1. 项目身份

工作名：

```text
vibit
```

类别：

```text
Agent-Native Server Framework
```

定位：

```text
vibit is an open-source agent-native server framework for building backends that AI coding agents can understand, extend, verify, and maintain from first principles.
```

在本仓库中，“AI-native” 首先指 agent-native maintainability。它主要不是指添加 AI gameplay features 或 AI product features。

## 2. 必读内容

在进行非平凡变更前，先阅读：

- `CONSTITUTION.md`
- 本文件
- `.arch/` 下相关 architecture manifests，如果它们已经存在
- `modules/<module>/module.yaml` 中的相关 module manifest，如果它已经存在
- `modules/<module>/AGENTS.md` 中的相关 module guide，如果它已经存在
- `changes/` 下相关 change spec，如果该变更已有 spec

如果预期产物尚不存在，不要发明隐藏假设。要么把缺失产物作为本次变更的一部分创建出来，要么记录它目前尚不可用。

## 3. 当前仓库状态

本仓库当前是 pre-alpha，并正在推进 `v0.1 alpha`。

短期目标是第一个 developer-usable alpha，让外部开发者可以下载、本地运行、检查，并把它作为加入开发的入口。持久目标文档是 `docs/v0.1-alpha-goal.md`，配套简体中文译本是 `docs/v0.1-alpha-goal.zh-CN.md`。`ADR-0086` 记录该决策。

长期产品目标是 AI 时代的 Nakama/Pitaya-class open-source game/backend server framework。这表示按 vibit 的 agent-native maintainability 模型覆盖同级常用能力；它不表示 direct Nakama/Pitaya API compatibility。

现有基础：

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/README.md`
- `.arch/modules.yaml`
- `.arch/conventions.yaml`
- `.arch/protocol.yaml`
- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `.arch/dependencies.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `docs/module-manifest.md`
- `docs/module-manifest.zh-CN.md`
- `docs/change-spec.md`
- `docs/change-spec.zh-CN.md`
- `changes/_template/`
- `docs/conversation-log.md`
- `docs/conversation-log.zh-CN.md`
- `conversations/`
- `docs/agent-decision-record.md`
- `docs/agent-decision-record.zh-CN.md`
- `decisions/`
- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `docs/dependency-adoption.md`
- `docs/dependency-adoption.zh-CN.md`
- `docs/game-protocol.md`
- `docs/game-protocol.zh-CN.md`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/authentication-proof-token-session-contract-dimensions.zh-CN.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/credential-storage-external-identity-linking-boundaries.zh-CN.md`
- `docs/session-persistence-websocket-handshake-decision-gates.md`
- `docs/session-persistence-websocket-handshake-decision-gates.zh-CN.md`
- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `docs/workflow.md`
- `docs/workflow.zh-CN.md`
- `schema/`
- `rules/`

框架实现代码现在已位于 `runtime/`，generated output 已位于 `runtime/internal/generated/`，verification commands 已通过 `tools/vibit` 和 Go tests 存在。当某个 capability 或 verification path 尚不存在时，应记录它目前不可用，而不是假装已经运行。

当前 runtime readiness decisions 指向 Go 作为第一版 server runtime implementation language、WebSocket 作为第一版 gameplay/client protocol、Protobuf 作为第一版 wire message format、modular monolith single-process server model、contract-first commands/queries/events/errors/permissions，以及 `inventory` 作为优先的第一 proof slice。在创建 runtime implementation code 前，必须阅读 `.arch/runtime.yaml`、`ADR-0004` 到 `ADR-0010`，并注意 `ADR-0003` 已被 superseded。

第一版 game protocol framework 记录在 `.arch/protocol.yaml`、`docs/game-protocol.md`、`ADR-0015` 和 `ADR-0016` 中。它定义 WebSocket-framed Protobuf envelope，使用显式 `kind`、`module` 和 `name` routing fields，并包含 session metadata、game target scopes、server-authoritative message rules、error mapping 和 compatibility expectations。在新增 `.proto` files、WebSocket protocol handlers、generated protocol output 或 client/server protocol rules 前，必须阅读它。

第一批 protocol source files 是 `proto/vibit/protocol/v1/envelope.proto` 和 `proto/vibit/inventory/v1/inventory.proto`。Buf configuration 位于 `buf.yaml` 和 `buf.gen.yaml`。`ADR-0016` 记录 envelope 和 generation configuration decision。生成的 Go Protobuf output 仍计划位于 `runtime/internal/generated/proto/`；不要手工创建或编辑生成的 Go Protobuf files。

Generated output standard 是 `docs/generated-output.md`，配套简体中文译本是 `docs/generated-output.zh-CN.md`。`ADR-0017` 记录 generated output decision。在添加 generated files、generated output checks 或 generator behavior 前，应阅读这些文件。`runtime/internal/generated/proto/` 下的 Go Protobuf output 必须是 `*.pb.go`，必须包含 `protoc-gen-go` generated-code marker，并且必须能追溯到现有 `.proto` source。

Runtime protocol adapter boundary standard 是 `docs/runtime-protocol-adapter.md`，配套简体中文译本是 `docs/runtime-protocol-adapter.zh-CN.md`。`ADR-0018` 记录 boundary decision。在添加 WebSocket transport code、Protobuf runtime adapter code、application dispatch code 或 domain runtime handlers 前，应阅读这些文件。

Active game server reference alignment standard 是 `docs/reference-game-server-alignment.md`，配套简体中文译本是 `docs/reference-game-server-alignment.zh-CN.md`。`ADR-0019` 记录 Nakama 和 Pitaya 是主动参考基线。在新增 game server capability families、runtime subsystems、social/realtime features、matchmaking、match runtime、cluster/RPC work 或 operational surfaces 前，应阅读 `.arch/reference.yaml` 和该标准。Nakama 和 Pitaya 指导 capability planning；它们不覆盖 vibit 的 constitution、ADRs、manifests、generated boundaries 或 verification commands。

Authentication、token 和 session validation design standard 是 `docs/authentication-token-session-validation.md`，配套简体中文译本是 `docs/authentication-token-session-validation.zh-CN.md`。`ADR-0023` 记录这个 design boundary。在添加 authentication、token behavior、credential storage、external identity linking、session persistence、request identity trust changes、Protobuf envelope authentication changes、WebSocket handshake authentication、runtime player account handlers 或 WebSocket player routes 前，应阅读它。Metadata-only `player_id` 和 `session_id` 不是 authenticated proof。`runtime.authentication_token_session_boundary` 是该边界的 repository check rule。

Authentication proof 与 token/session contract dimensions standard 是 `docs/authentication-proof-token-session-contract-dimensions.md`，配套简体中文译本是 `docs/authentication-proof-token-session-contract-dimensions.zh-CN.md`。在定义或修改 authentication proof、token/session validation、actor kinds、validation statuses、proof statuses、failure classes、retryability、request identity handoff、session errors、session permissions 或 validation events 前，应阅读它。它只 ratify semantic vocabulary；不选择 login methods、token formats、credential storage、session persistence、Protobuf envelope behavior、WebSocket handshake behavior、runtime player handlers 或 WebSocket routes。

Credential storage 与 external identity linking boundary standard 是 `docs/credential-storage-external-identity-linking-boundaries.md`，配套简体中文译本是 `docs/credential-storage-external-identity-linking-boundaries.zh-CN.md`。在添加 credential storage、external identity linking、login methods、provider subjects、password hashing、OAuth、OIDC、provider SDKs、account linking、account recovery、account merge behavior 或相关 schema 前，应阅读它。该标准只定义 deferred capability families。它保持 `player_accounts` 和 `player_account_events` 作为 account lifecycle tables，并不授权 credential tables、external identity tables、provider dependencies、runtime lookup code、player lifecycle table changes 或 direct Nakama/Pitaya API compatibility。

Session persistence 与 WebSocket handshake decision-gates standard 是 `docs/session-persistence-websocket-handshake-decision-gates.md`，配套简体中文译本是 `docs/session-persistence-websocket-handshake-decision-gates.zh-CN.md`。在添加 session persistence、WebSocket handshake authentication、reconnect behavior、connection epoch behavior、token/session carriers、session-related Protobuf envelope changes、handshake/system messages 或 route-level authentication 前，应阅读它。它只把 request-level、first-message、handshake-level、every-request 和 hybrid validation 定义为未来选择；它不选择 production model、session store、session tables、envelope behavior 或 handshake behavior。

Login method 与 token format ratification standard 是 `docs/login-method-token-format-ratification.md`，配套简体中文译本是 `docs/login-method-token-format-ratification.zh-CN.md`。`ADR-0024` 记录该 ratification boundary。在选择第一批 login methods、token model、token format、proof carrier posture、token lifecycle semantics、credential/token/session schema gates 或 implementation queue 前，应阅读它。它只指导 comparison 与 ratification；不授权 runtime authentication、token parsing、credential storage、external identity linking、session persistence、Protobuf envelope changes、WebSocket handshake authentication、runtime player handlers 或 WebSocket routes。

Selected login/token boundary check standard 是 `docs/selected-login-token-boundary-checks.md`，配套简体中文译本是 `docs/selected-login-token-boundary-checks.zh-CN.md`。`ADR-0030` 记录 repository check decision。在添加 runtime authentication、token validation、token issuance、logout、refresh behavior、credential 或 token repositories、authentication Protobuf sources、generated authentication contract shapes、WebSocket proof carriers、authentication migrations、authentication dependencies，或改变已选 `device_credential_login` 加 opaque access-token 姿态前，应阅读它。`runtime.selected_login_token_boundary` 是该已选姿态的 repository check rule。

Credential record schema boundary standard 是 `docs/credential-record-schema-boundary.md`，配套简体中文译本是 `docs/credential-record-schema-boundary.zh-CN.md`。`ADR-0032` 记录 boundary decision。Credential migration source 现在位于 `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`；它不授权 repositories、PostgreSQL adapters、runtime lookup、handlers、routes、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior。

Token verifier record schema boundary standard 是 `docs/token-verifier-record-schema-boundary.md`，配套简体中文译本是 `docs/token-verifier-record-schema-boundary.zh-CN.md`。`ADR-0033` 记录 boundary decision。Token verifier migration source 现在位于 `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`；它不授权 repositories、PostgreSQL adapters、runtime token issuance、validation、logout、refresh、cleanup、handlers、routes、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior。

Authentication schema migration queue standard 是 `docs/authentication-schema-migration-queue.md`，配套简体中文译本是 `docs/authentication-schema-migration-queue.zh-CN.md`。`ADR-0034` 记录 queue decision，`ADR-0035` 记录 adapter implementation gate。Credential 与 token verifier migration sources、static checks、storage-neutral authentication repository interface boundary，以及 authentication PostgreSQL adapter boundary 现在都已存在。这个 queue 不授权 runtime credential lookup、token issuance、token validation、logout execution、refresh、cleanup jobs、handlers、routes、generated authentication shapes、Protobuf messages、WebSocket proof carriers、WebSocket handshake authentication、authentication dependencies 或 production authentication behavior。

`M-019 Token And Credential Verifier Algorithm Redaction Boundary` 已完成。`W-0091` 已在 `docs/token-credential-verifier-algorithm-redaction-boundary.md` 和 `ADR-0040` 中定义第一版 planned high-entropy verifier posture，且没有添加 runtime authentication behavior。第一版 planned verifier algorithm family 是 `vibit_hmac_sha256_v1`；未来 first-posture code 在后续 code gate 授权后可以使用 Go standard library packages `crypto/hmac`、`crypto/sha256`、`crypto/subtle`、`crypto/rand` 和 `encoding/base64`。Lookup digests 和 verifier digests 都不是 log-safe，verifier digest comparison 必须 constant-time，raw credential/token material 必须至少有 256 bits entropy。`runtime.token_credential_verifier_algorithm_redaction_boundary` 是已完成 verifier boundary 的 repository check rule。

`M-020 Secret Configuration And Verifier Key Loading Boundary` 已完成。`W-0092` 已在 `docs/secret-configuration-verifier-key-loading-boundary.md` 和 `ADR-0041` 中定义 future key loading posture，且没有添加 runtime authentication behavior。Future verifier key loading 由 `runtime/internal/app` 下的 application 拥有；第一版 local implementation 可在后续 code gate 授权后使用 process environment configuration 或 explicit runtime secret input；external KMS 或 secret-manager integration 仍然需要后续 dependency 和 operations gates。必须使用四个 separated logical verifier keys，`verifier_key_id` 默认不是 log-safe，production key configuration 无效时必须 fail closed，committed production-like secret values 仍然禁止。本 boundary 不得实现 secret loading、environment parsing、token generation、credential generation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository changes、migration schema changes、KMS integration、secret-manager integration 或 production authentication behavior。`runtime.secret_configuration_verifier_key_loading_boundary` 是该 boundary 的 repository check rule。

`M-021 Token And Credential Material Generation Boundary` 已完成。`W-0093` 已在 `docs/token-credential-material-generation-boundary.md` 和 `ADR-0042` 中定义 future raw device credential 与 opaque access-token material generation posture，且没有添加 runtime authentication behavior。Future material generation 由 `runtime/internal/app` 下的 application 拥有；第一版 device credential 与 access token 都是 server-issued 且 application-generated；raw material 必须是 32 cryptographically random bytes，至少 256 bits entropy；text presentation 使用 URL-safe unpadded Base64 或等价 encoding；raw material 只能 one-time client-visible，且不得存储。Repository handoff 保持 digest-only。Future first-posture generation helpers 可在后续 code gate 授权后使用 Go standard library `crypto/rand` 和 `encoding/base64`。本 boundary 不得实现 token generation、credential generation、secret loading、verifier digest computation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository changes、migration schema changes 或 production authentication behavior。`runtime.token_credential_material_generation_boundary` 是该 boundary 的 repository check rule。

`M-022 Verifier Digest Computation And Comparison Boundary` 已完成。`W-0094` 已在 `docs/verifier-digest-computation-comparison-boundary.md` 和 `ADR-0043` 中定义 future verifier digest computation 与 constant-time comparison posture，且没有添加 runtime authentication behavior。Future digest helpers 由 `runtime/internal/app` 下的 application 拥有；canonical HMAC input 带 version 且 length-prefixed；lookup digest equality 只能用于 record selection；validation 可跨 active 与 accepted previous key sets 计算 lookup candidates；stored `verifier_key_id` 用于选择 verifier digest computation 的 verifier key set；verifier digest comparison 必须 constant-time；invalid lookup、mismatch、unknown key id、unsupported algorithm、malformed proof 以及 expired 或 revoked proof 都收敛到同一个 public invalid-proof class。Future first-posture helpers 可在后续 code gate 授权后使用 Go standard library `crypto/hmac`、`crypto/sha256` 和 `crypto/subtle`。本 boundary 不得实现 verifier digest computation、verifier comparison、token generation、credential generation、secret loading、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository changes、migration schema changes 或 production authentication behavior。`runtime.verifier_digest_computation_comparison_boundary` 是该 boundary 的 repository check rule。

`M-023 Authentication Service Implementation Readiness Gate` 已完成。`W-0095` 已在 `docs/authentication-service-implementation-readiness-gate.md` 和 `ADR-0044` 中定义 readiness gate，且没有添加 runtime authentication behavior。Future service implementation 仍由 `runtime/internal/app` 下的 application 拥有，package candidate 是 `runtime/internal/app/authentication`。第一段 code slice 必须单独授权；recommended next gate 是 local verifier key configuration loading。本 gate 定义 prior boundaries、allowed 与 forbidden write areas、test classes、redaction expectations、Nakama/Pitaya capability mapping 和 deferrals。它不得实现 authentication service code、secret loading、token generation、credential generation、verifier digest computation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository changes、migration schema changes 或 production authentication behavior。`runtime.authentication_service_implementation_readiness_gate` 是该 gate 的 repository check rule。

`M-024 Local Verifier Key Configuration Loading Gate` 已完成。`W-0096` 已在 `docs/local-verifier-key-configuration-loading-gate.md` 和 `ADR-0045` 中定义该 gate，且没有添加 Go code 或 runtime authentication behavior。`W-0097` 已在 `runtime/internal/app/authentication` 下实现 explicit in-memory verifier key set validator，并覆盖 accepted input、copying、immutable accessors、invalid key sets 和 redacted errors 的测试。Environment parsing 和 Base64 text decoding 推迟到 `W-0098`，也就是 environment verifier key loader gate。它不得 parse process environment variables、读取 local secret files、wire startup、集成 KMS 或 cloud secret managers、generate tokens 或 credentials、compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、添加 authentication dependencies 或添加 production authentication behavior。`runtime.local_verifier_key_configuration_loading_gate` 是该 gate 的 repository check rule。

`M-026 Environment Verifier Key Loader Gate` 已完成。`W-0098` 已在 `docs/environment-verifier-key-loader-gate.md` 和 `ADR-0046` 中定义该 gate，且没有添加 Go code 或 runtime authentication behavior。Future process environment loader work 属于 `runtime/internal/app/authentication`，应使用 `verifier_key_env.go` 和 `verifier_key_env_test.go`，必须 decode required environment key text 并调用 `NewVerifierKeySet`，并保持 environment variable values 和 full concrete key set ids redacted。它不得 wire startup、读取 local secret files、parse `.env` files、接受 CLI secret input、集成 KMS 或 cloud secret managers、generate tokens 或 credentials、compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、添加 authentication dependencies 或添加 production authentication behavior。`runtime.environment_verifier_key_loader_gate` 是该 gate 的 repository check rule。

`M-027 Environment Verifier Key Loader Implementation` 已完成。`W-0099` 已实现 `runtime/internal/app/authentication/verifier_key_env.go` 和 `runtime/internal/app/authentication/verifier_key_env_test.go`。该 loader 接收 explicit lookup function，提供很小的 `os.LookupEnv` process adapter，decode Base64URL unpadded 与 standard padded Base64 key text，调用 `NewVerifierKeySet`，并返回 redacted typed errors。没有后续 bounded work item 授权时，不得把该 loader wire into process startup、读取 local secret files、parse `.env` files、接受 CLI secret input、集成 KMS 或 cloud secret managers、generate tokens 或 credentials、compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、添加 authentication dependencies 或添加 production authentication behavior。

`M-028 Token And Credential Material Generation Implementation Gate` 已完成。`W-0100` 已在 `docs/token-credential-material-generation-implementation-gate.md` 和 `ADR-0047` 中定义该 gate，且没有添加 Go code 或 runtime authentication behavior。Future material generation helper work 属于 `runtime/internal/app/authentication`，应使用 `material_generation.go` 和 `material_generation_test.go`，必须使用 32 random bytes、URL-safe unpadded Base64 presentation、explicit `io.Reader` entropy-source handoff、copied raw bytes、redacted errors 和 focused tests。它不得 compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。`runtime.token_credential_material_generation_implementation_gate` 是该 gate 的 repository check rule。

`M-029 Token And Credential Material Generation Helper Implementation` 已完成。`W-0101` 已实现 `runtime/internal/app/authentication/material_generation.go` 和 `runtime/internal/app/authentication/material_generation_test.go`。这些 helpers 暴露 `GenerateDeviceCredentialMaterial` 和 `GenerateAccessTokenMaterial`，接收 explicit `io.Reader`，使用 `io.ReadFull` 读取 32 bytes，编码 URL-safe unpadded Base64 text，保留 `MaterialKind`，返回 raw bytes copy，拒绝 nil readers、read errors、short reads、all-zero material 和 repeated-single-byte material，并返回 redacted typed errors。没有后续 bounded work item 授权时，不要把这些 helpers 扩展为 compute digests、compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。

`M-030 Verifier Digest Helper Implementation Gate` 已完成。`W-0102` 已在 `docs/verifier-digest-helper-implementation-gate.md` 和 `ADR-0048` 中定义 verifier digest helper implementation gate，且没有添加 Go code 或 runtime authentication behavior。Future digest helper work 属于 `runtime/internal/app/authentication`，应使用 `verifier_digest.go` 和 `verifier_digest_test.go`，必须用 versioned ASCII header、null separator、length-prefixed purpose label 和 raw material 构造 deterministic canonical input，用 already-validated `VerifierKeySet` 中匹配的 logical key 计算 HMAC-SHA-256，返回 copied 32-byte digest bytes，暴露 redacted errors，并添加 focused tests。它不得 compare verifiers、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。`runtime.verifier_digest_helper_implementation_gate` 是该 gate 的 repository check rule。

`M-031 Verifier Digest Computation Helper Implementation` 已完成。`W-0103` 已实现 `runtime/internal/app/authentication/verifier_digest.go` 和 `runtime/internal/app/authentication/verifier_digest_test.go`。这些 helpers 暴露 `ComputeCredentialLookupDigest`、`ComputeCredentialVerifierDigest`、`ComputeTokenLookupDigest` 和 `ComputeTokenVerifierDigest`，构造 canonical `vibit.auth.verifier.input.v1` HMAC input，使用 `VerifierKeySet` 中匹配的 logical key，通过 `ComputedDigest` 返回 copied 32-byte digest bytes，并使用 redacted typed errors。没有后续 bounded work item 授权时，不要把这些 helpers 扩展为 compare verifier digests、实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。下一个 next ready work item 是 verifier digest comparison helper gate。

`M-032 Verifier Digest Comparison Helper Gate` 已完成。`W-0104` 已在 `docs/verifier-digest-comparison-helper-gate.md` 和 `ADR-0049` 中定义 verifier digest comparison helper gate，且没有添加 Go code 或 runtime authentication behavior。Future comparison helper work 属于 `runtime/internal/app/authentication`，应使用 `verifier_comparison.go` 和 `verifier_comparison_test.go`，必须保持 `verifier_digest.go` computation-only，只比较 computed verifier digest bytes 和 stored verifier digest bytes，首选 `crypto/hmac.Equal`，仅允许 `crypto/subtle.ConstantTimeCompare` 作为等价 constant-time alternative，拒绝 lookup digest classes 和 malformed input，并避免在 comparison 中使用 raw material、lookup digests、key ids、database-only equality、protocol metadata 和 public failure details。它不得实现 authentication service behavior、暴露 protocol carriers、改变 repositories、改变 migrations、wire startup、添加 authentication dependencies 或添加 production authentication behavior。`runtime.verifier_digest_comparison_helper_gate` 是该 gate 的 repository check rule。下一个 next ready work item 是 verifier digest comparison helper implementation slice。

`M-033 Verifier Digest Comparison Helper Implementation` 已完成。`W-0105` 已实现 `runtime/internal/app/authentication/verifier_comparison.go` 和 `runtime/internal/app/authentication/verifier_comparison_test.go`。这些 helpers 暴露 `CompareCredentialVerifierDigest` 和 `CompareTokenVerifierDigest`，返回 `VerifierComparisonResult`，使用 `crypto/hmac.Equal`，只比较 `ComputedDigest` verifier bytes 和 stored verifier digest bytes，拒绝 lookup digest classes、wrong verifier classes、missing input、malformed computed input 和 malformed stored input，并保持 errors redacted。没有后续 bounded work item 授权时，不要把这些 helpers 扩展为 authentication service behavior、login execution、token validation、logout、refresh、cleanup、protocol carriers、repository calls、migration changes、startup wiring、authentication dependencies 或 production authentication behavior。下一个 next ready work item 是 authentication service behavior implementation gate。

`M-034 Authentication Service Behavior Implementation Gate` 已完成。`W-0106` 已在 `docs/authentication-service-behavior-implementation-gate.md` 和 `ADR-0050` 中定义该 gate，且没有添加 service code 或 runtime authentication behavior。它要求 future service behavior 继续由 `runtime/internal/app/authentication` 下的 application 拥有，必须只通过 application unit-of-work boundary 使用 `authentication.Repository`，必须按固定顺序组合 verifier key、material generation、digest computation 与 comparison helpers，并且必须折叠 public proof failures。`runtime.authentication_service_behavior_implementation_gate` 是该 gate 和 skeleton boundary 的 repository check rule。

`M-035 Authentication Service Behavior Skeleton` 已完成。`W-0107` 添加了 `runtime/internal/app/authentication/service.go` 和 `runtime/internal/app/authentication/service_test.go`，作为 skeleton-only application service shape。该 skeleton 定义 typed dependencies、request/result vocabulary、redacted internal failure classes、public error codes，以及 fail-closed 的 `AUTHENTICATION_NOT_IMPLEMENTED` 或 `AUTHENTICATION_REFRESH_NOT_SUPPORTED` 行为。它不执行 login、不 validate access tokens、不 issue tokens、不调用 repositories、不 revoke tokens、不 refresh tokens、不暴露 protocol carriers、不 wire startup、不添加 dependencies、不修改 migrations，也不添加 production authentication behavior。下一个 next ready work item 是 device credential login service behavior gate。

`M-036 Device Credential Login Service Behavior Gate` 已完成。`W-0108` 已在 `docs/device-credential-login-service-behavior-gate.md` 和 `ADR-0051` 中定义该 gate，且没有添加 login execution code 或 runtime authentication behavior。Future device credential login 必须继续由 `runtime/internal/app/authentication/service.go` 内的 application 拥有，必须把 `CredentialProof` 视为 server-issued Base64URL unpadded 32-byte high-entropy material，必须在 unit-of-work 前拒绝 missing 或 malformed proof，必须通过 unit-of-work capabilities 使用 existing authentication 和 player repositories，必须在 token generation 前 compare credential verifier digest，必须要求 active player account state，只能 store token digests，并且只能在 unit-of-work 成功后返回 raw access-token text。它必须保留 access-token validation、logout、refresh、cleanup、protocol carriers、startup wiring、repository changes、migrations、generated files、dependencies 和 broader production authentication behavior 的延期项。`runtime.device_credential_login_service_behavior_gate` 是该 gate 的 repository check rule。下一个 next ready work item 是 device credential login service behavior implementation slice。

`M-037 Device Credential Login Service Behavior Implementation` 已完成。`W-0109` 已在 `runtime/internal/app/authentication/service.go` 中实现 `AuthenticateWithDeviceCredential`，并在 `runtime/internal/app/authentication/service_test.go` 中添加 focused tests。该 service 在 unit-of-work 前校验 Base64URL unpadded 32-byte proof，使用 verifier helpers，通过 unit-of-work capabilities 获取 authentication 和 player repositories，要求 active credential 与 active player account state，只存储 opaque access-token digests，并且只在 token storage 和 unit-of-work commit 成功后返回 raw access-token text。它没有实现 access-token validation、logout、refresh、cleanup、protocol carriers、startup wiring、repository interface changes、migrations、generated files、dependencies 或 broader production authentication behavior。`runtime.device_credential_login_service_behavior_implementation` 是该 slice 的 repository check rule。下一个 next ready work item 是 access-token validation service behavior gate。

`M-038 Access Token Validation Service Behavior Gate` 已完成。`W-0110` 已在 `docs/access-token-validation-service-behavior-gate.md` 和 `ADR-0052` 中定义 future access-token validation behavior gate，且没有添加 validation execution code。Future validation behavior 必须继续由 `runtime/internal/app/authentication` 下的 application 拥有，必须把 `AccessToken` 视为 opaque Base64URL unpadded 32-byte proof，必须在 unit-of-work 前拒绝 missing 或 malformed proof，必须通过 unit-of-work capabilities 使用 authentication 和 player repositories，必须在 request identity 前 compare token verifier digest，必须要求 active player account state，必须在 session persistence 被 ratify 前保持 `SessionValidated` 为 false，并且必须折叠 public invalid-token failures。它必须保留 WebSocket proof carriers、handshake authentication、route protection、session persistence、logout、refresh、cleanup、protocol carriers、startup wiring、repository changes、migrations、generated files、dependencies 和 broader production authentication behavior 的延期项。`runtime.access_token_validation_service_behavior_gate` 是该 gate 的 repository check rule。下一个 next ready work item 是 access-token validation service behavior implementation slice。

`M-039 Access Token Validation Service Behavior Implementation` 已完成。`W-0111` 已在 `runtime/internal/app/authentication/service.go` 中实现 `ValidateAccessToken`，并在 `runtime/internal/app/authentication/service_test.go` 中添加 focused tests。该 service 会在 unit-of-work 前拒绝 missing 或 malformed opaque Base64URL unpadded 32-byte proof，用现有 helpers 计算 token lookup 和 verifier digests，通过 unit-of-work capabilities 获取 authentication 和 player repositories，检查 token kind、status、verifier posture、audience、issue time 和 expiration，在 request identity 前比较 token verifier digest，要求 active player account state，只在 unit-of-work 成功后返回 validated player identity，并保持 `SessionValidated` 为 false。它没有添加 protocol carriers、WebSocket handshake authentication、route protection、session persistence、logout、refresh、cleanup、startup wiring、repository interface changes、PostgreSQL adapter changes、migrations、generated files、dependencies 或 broader production authentication behavior。`runtime.access_token_validation_service_behavior_implementation` 是该 slice 的 repository check rule。

`M-040 Next Direction Confirmation Gate` 已完成。`W-0112` 在 maintainer 授权继续推进后选择了 `expose_access_token_protocol_carrier_and_route_protection_gate`。`M-041 Access Token Protocol Carrier And Route Protection Gate` 已完成。`W-0113` 定义了 `docs/access-token-protocol-carrier-route-protection-gate.md`、`docs/access-token-protocol-carrier-route-protection-gate.zh-CN.md` 和 `ADR-0053`，但没有添加实现。该 gate 选择 request-level validation，并把未来 Protobuf payload wrapper candidate 定为 `vibit.authentication.v1.AuthenticatedRequest`；现有 Protobuf envelope 保持不变，WebSocket transport 保持 credential-neutral，并要求在 protected domain dispatch 前经过 application-owned route policy。`runtime.access_token_protocol_carrier_route_protection_gate` 是该 gate 的 repository check rule。

`M-042 Access Token Protocol Carrier And Route Protection Implementation` 已完成。`W-0114` 添加了 `proto/vibit/authentication/v1/authenticated_request.proto`，并通过 Buf 生成 `runtime/internal/generated/proto/vibit/authentication/v1/authenticated_request.pb.go`；同时添加 `runtime/internal/app/route_authentication.go`、`runtime/internal/app/authentication/route_validator.go` 以及 focused app/protocol/WebSocket tests。Protobuf adapter 会 unwrap `vibit.authentication.v1.AuthenticatedRequest`，在 protected domain dispatch 前通过 application route protection validate access-token proof，拒绝 protected routes 上的 metadata-only identity，保持 `SessionValidated` 为 false，保持现有 envelope route fields 不变，并保持 WebSocket transport credential-neutral。它没有添加 WebSocket handshake authentication、session persistence、startup wiring、repository changes、migrations、dependencies、logout、refresh、cleanup、token rotation 或 broader production authentication behavior。Work queue 现在 blocked 在 `M-043/W-0115`，即 next-direction confirmation gate。

`M-044 Runtime Authentication Startup Composition Gate` 已完成。`W-0116` 在 maintainer 选择 `wire_runtime_authentication_startup_composition` 且要求重点参考 Nakama 和 Pitaya 后，定义了 `docs/runtime-authentication-startup-composition-gate.md`、`docs/runtime-authentication-startup-composition-gate.zh-CN.md` 和 `ADR-0054`。该 gate 只允许在 `runtime/cmd/vibit-server` 做 process startup composition，第一条 composed path 仅限 `VIBIT_RUNTIME_STORE=postgres`。Memory store 仍是 metadata-only bootstrap behavior；WebSocket transport 保持 credential-neutral；现有 Protobuf envelope 不变；session persistence、WebSocket handshake authentication、authentication command routes、repository changes、migrations、dependencies、logout、refresh、cleanup、token rotation 和 broader production behavior 继续 deferred。`runtime.authentication_startup_composition_gate` 是该 gate 的 repository check rule。

`M-045 Runtime Authentication Startup Composition Implementation` 已完成。`W-0117` 在 `runtime/cmd/vibit-server/main.go` 中实现 startup composition，并在 `runtime/cmd/vibit-server/main_test.go` 中添加 focused tests。PostgreSQL runtime path 现在通过 existing environment loader 加载 verifier keys，用 PostgreSQL unit-of-work、`crypto/rand.Reader`、startup-owned clock、startup-owned token record id generator、默认或配置的 1h access-token lifetime、以及默认或配置的 `vibit_gameplay_runtime_requests` audience 构造 existing application authentication service，然后把 `app.NewRouteProtector(authentication.NewRouteAccessTokenValidator(service))` 注入 Protobuf frame handler。Memory runtime path 仍保持 metadata-only bootstrap behavior。没有 later bounded work item 授权时，不要把这个 slice 扩展为 session persistence、WebSocket handshake authentication、authentication command routes、repository changes、migrations、generated files、dependencies、logout、refresh、cleanup、token rotation、token validation audit mutation 或 direct Nakama/Pitaya API compatibility。Work queue 现在 blocked 在 `M-046/W-0118`，即 next-direction confirmation gate。

`M-046 Next Direction Confirmation Gate` 已完成。`W-0118` 在 startup composition 之后选择了 `add_authentication_command_protocol_messages_and_login_route_registration`，并继续把 Nakama 和 Pitaya 作为重点 reference baselines。`M-047 Authentication Command Protocol And Login Route Gate` 已完成。`W-0119` 定义了 `docs/authentication-command-protocol-login-route-gate.md`、`docs/authentication-command-protocol-login-route-gate.zh-CN.md` 和 `ADR-0055`。该 gate 只授权围绕现有 authentication service 暴露 public `runtime.authentication.AuthenticateWithDeviceCredential` Protobuf command route。`runtime.authentication_command_protocol_login_route_gate` 是该 gate 和 implementation slice 的 repository check rule。

`M-048 Authentication Command Protocol And Login Route Implementation` 已完成。`W-0120` 添加了 `proto/vibit/authentication/v1/authentication.proto`，并通过 Buf 生成 `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go`；同时添加 `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`、`runtime/internal/app/bootstrap/authentication.go`，只在 PostgreSQL startup composition 中注册 public login route，并为 `runtime.authentication.AuthenticateWithDeviceCredential` 添加 transaction-wrapper bypass，因为 authentication service 拥有自己的 unit-of-work。现有 Protobuf envelope 不变，WebSocket transport 保持 credential-neutral，memory durable login 仍不可用；repository interfaces、PostgreSQL adapters、migrations、dependencies、session persistence、WebSocket handshake authentication、logout、refresh、cleanup、token rotation、token validation audit mutation 和 direct Nakama/Pitaya API compatibility 继续 deferred。Work queue 现在 blocked 在 `M-049/W-0121`，即 next-direction confirmation gate。

`M-049 Next Direction Confirmation Gate` 已完成。`W-0121` 在 public login route 之后选择了 `ratify_session_persistence_and_websocket_handshake_authentication`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-050 Session Persistence And WebSocket Handshake Ratification` 已完成。`W-0122` 定义了 `docs/session-persistence-websocket-handshake-ratification.md`、`docs/session-persistence-websocket-handshake-ratification.zh-CN.md` 和 `ADR-0056`。Ratified current path 仍然是通过 `vibit.authentication.v1.AuthenticatedRequest` 做 request-level opaque access-token validation；WebSocket transport 继续 credential-neutral；现有 Protobuf envelope 不变；未来 connection-level identity deferred 到 first-message protocol/application binding gate；未来 durable session persistence deferred 到 PostgreSQL-first schema/repository/migration gate。该 ratification 不授权 session tables、migrations、repositories、dependencies、WebSocket handshake proof carriers、logout、refresh、cleanup、token rotation、reconnect/epoch behavior、memory durable authentication behavior 或 direct Nakama/Pitaya API compatibility。`runtime.session_persistence_websocket_handshake_ratification` 是 repository check rule。Work queue 现在 blocked 在 `M-051/W-0123`，即 next-direction confirmation gate。

`M-051 Next Direction Confirmation Gate` 已完成。`W-0123` 在 session 和 handshake ratification 之后选择了 `define_first_message_connection_binding_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-052 First Message Connection Binding Gate` 已完成。`W-0124` 定义了 `docs/first-message-connection-binding-gate.md`、`docs/first-message-connection-binding-gate.zh-CN.md` 和 `ADR-0057`。该 gate 选择 future `runtime.authentication.BindConnection` system route，并选择 `vibit.authentication.v1.BindConnectionRequest` 与 `vibit.authentication.v1.BindConnectionResponse` 作为 payload candidates。WebSocket transport 继续 credential-neutral；现有 Protobuf envelope 不变；request-level access-token validation 仍然是当前 protected-route path；future connection-bound identity 仍然需要 implementation gate。该 gate 不授权 Protobuf source changes、generated output、connection binding registries、route-policy use of bound identity、session persistence、repositories、migrations、dependencies、logout/revocation、reconnect/epoch behavior、memory durable authentication behavior 或 direct Nakama/Pitaya API compatibility。`runtime.first_message_connection_binding_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-053/W-0125`，即 next-direction confirmation gate。

`M-053 Next Direction Confirmation Gate` 已完成。`W-0125` 在 first-message connection binding gate 之后选择了 `define_first_message_connection_binding_implementation_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-054 First Message Connection Binding Implementation Gate` 已完成。`W-0126` 定义了 `docs/first-message-connection-binding-implementation-gate.md`、`docs/first-message-connection-binding-implementation-gate.zh-CN.md` 和 `ADR-0058`。该 gate 定义 future `runtime.authentication.BindConnection` bounded implementation slice，包括计划中的 `BindConnectionRequest`、`BindConnectionResponse`、`ConnectionBindingStatus`、Buf-generated Go output、Protobuf adapter bridge、application binding boundary、PostgreSQL startup composition、public error mapping 和 required tests。该 gate 仍然不授权 Go runtime behavior、Protobuf source changes、generated output、connection binding registries、route-policy use of bound identity、session persistence、repositories、migrations、dependencies、logout/revocation、reconnect/epoch behavior、memory durable authentication behavior 或 direct Nakama/Pitaya API compatibility。WebSocket transport 继续 credential-neutral；现有 Protobuf envelope 不变；request-level access-token validation 仍然是当前 protected-route path。`runtime.first_message_connection_binding_implementation_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-055/W-0127`，即 next-direction confirmation gate。

`M-055 Next Direction Confirmation Gate` 已完成。`W-0127` 选择了 `implement_first_message_connection_binding`，并继续把 Nakama 和 Pitaya 作为重点 reference baselines。`M-056 First Message Connection Binding Implementation` 已完成。`W-0128` 实现了 `ADR-0058` 授权的 bounded `runtime.authentication.BindConnection` system route slice。Authentication Protobuf source 现在包含 `BindConnectionRequest`、`BindConnectionResponse` 和 `ConnectionBindingStatus`，并通过 Buf 更新了 generated Go output。Application 层拥有 `ConnectionBinder`；Protobuf adapter 在普通 route protection 和 dispatch 之前处理该 system route；PostgreSQL startup composition 在 authentication service 已 compose 时注入 binder；WebSocket transport 只传递 server-observed connection id 和 epoch metadata，不解析 credential。该实现参考 Nakama 的“authenticated session material precedes authenticated realtime socket use”和 Pitaya 的 acceptor/session/handler separation，但不添加 direct Nakama/Pitaya API compatibility。它也不添加 durable session persistence、connection registries、让 bound identity 满足普通 protected route policy、WebSocket handshake authentication、transport credential carriers、repositories、migrations、dependencies、logout/revocation active-connection invalidation、reconnect/resume/duplicate replacement policy、memory durable authentication behavior、presence、rooms、parties、match runtime 或 broader game backend behavior。`runtime.first_message_connection_binding_implementation` 是 repository check rule。

`M-057 Next Direction Confirmation Gate` 已完成。`W-0129` 在 first-message connection binding implementation 之后选择了 `define_postgres_session_persistence_schema_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-058 PostgreSQL Session Persistence Schema Gate` 已完成。`W-0130` 定义了 `docs/postgres-session-persistence-schema-gate.md`、`docs/postgres-session-persistence-schema-gate.zh-CN.md` 和 `ADR-0059`。该 gate 选择 PostgreSQL 作为 first durable runtime session target，选择 `runtime_sessions` 作为 future logical table candidate，并选择 `runtime/migrations/postgres/000005_create_runtime_sessions.sql` 作为 future migration source candidate。它参考 Nakama 的 first-class session lifecycle 和 Pitaya 的 session-context separation，但不添加 SQL migration source、不创建 session tables、不添加 session repositories、不添加 PostgreSQL adapters、不创建 runtime session behavior、不把 `RequestIdentity.SessionValidated` 设为 true、不改变 route policy、不改变现有 Protobuf envelope、不添加 WebSocket handshake authentication、不解析 transport credential carriers、不添加 logout/revocation active-connection invalidation、不添加 reconnect/epoch behavior、不添加 dependencies，也不提供 direct Nakama/Pitaya API compatibility。`runtime.postgres_session_persistence_schema_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-059/W-0131`，即 next-direction confirmation gate。

`M-059 Next Direction Confirmation Gate` 已完成。`W-0131` 在 PostgreSQL session persistence schema gate 之后选择了 `implement_runtime_sessions_migration_source`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-060 Runtime Sessions Migration Source` 已完成。`W-0132` 添加了 `runtime/migrations/postgres/000005_create_runtime_sessions.sql` 和 `ADR-0060`。该 migration 创建 PostgreSQL `runtime_sessions` lifecycle table，包含 actor/player identity、session status、issued/expires/last_seen timestamps、可选 revocation fields，以及可选 `authentication_access_tokens(token_record_id)` linkage。它不存储 raw access-token text、raw credentials、token digests、credential digests、WebSocket connection state 或 connection registry rows。它也不添加 session repositories、PostgreSQL session adapters、runtime session validation、login 或 BindConnection 中的 session creation、route-policy use of session or bound identity、WebSocket handshake authentication、logout/revocation active-connection invalidation、reconnect/epoch behavior、dependencies、memory durable session behavior，也不提供 direct Nakama/Pitaya API compatibility。`runtime.runtime_sessions_migration_source` 是 repository check rule。Work queue 现在 blocked 在 `M-061/W-0133`，即 next-direction confirmation gate。

`M-061 Next Direction Confirmation Gate` 已完成。`W-0133` 在 runtime sessions migration source 之后选择了 `define_session_repository_boundary`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-062 Session Repository Boundary Gate` 已完成。`W-0134` 定义了 `docs/session-repository-boundary.md`、`docs/session-repository-boundary.zh-CN.md` 和 `ADR-0061`。该 gate 记录 `runtime/internal/app/session` 作为 future storage-neutral repository owner candidate，并记录 `runtime/internal/platform/persistence/postgres` 作为 future PostgreSQL adapter owner。它定义了 `CreateRuntimeSession`、`FindActiveSessionByID`、`UpdateRuntimeSessionLastSeen`、`MarkRuntimeSessionExpired`、`RevokeRuntimeSession` 等 candidate lifecycle capabilities，但不添加 Go repository code、PostgreSQL adapter behavior、runtime session validation、login 或 BindConnection 中的 session creation、route-policy use of persisted session or bound identity、WebSocket handshake authentication、transport credential carriers、Protobuf session messages、logout/revocation active-connection invalidation、reconnect/epoch behavior、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.session_repository_boundary` 是 repository check rule。Work queue 现在 blocked 在 `M-063/W-0135`，即 next-direction confirmation gate。

`M-063 Next Direction Confirmation Gate` 已完成。`W-0135` 在 session repository boundary 之后选择了 `implement_session_repository_interface`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-064 Session Repository Interface Implementation` 已完成。`W-0136` 添加了 `runtime/internal/app/session/repository.go`、`runtime/internal/app/session/repository_test.go` 和 `ADR-0062`。该 package 定义 storage-neutral runtime session lifecycle value types、active/expired/revoked status vocabulary、第一阶段 player actor session、`runtime/internal/app/session.Repository`、生命周期 query/mutation types，以及 creation、lookup、active lookup、last-seen update、expiration、revocation 和 bounded active-session listing 的 normalization helpers。它不添加 PostgreSQL session adapter、SQL query execution、unit-of-work factory wiring、runtime session creation 或 validation、`RequestIdentity.SessionValidated = true`、route-policy 使用 persisted session 或 bound identity、WebSocket handshake authentication、transport credential carrier、Protobuf session message、generated output、logout/revocation active-connection invalidation、reconnect/epoch behavior、dependency、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.session_repository_interface_implementation` 是该 slice 的 repository check rule。Work queue 现在 blocked 在 `M-065/W-0137`，即 next-direction confirmation gate。

`M-065 Next Direction Confirmation Gate` 已完成。`W-0137` 在 session repository interface implementation 之后选择了 `define_session_postgresql_adapter_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-066 Session PostgreSQL Adapter Gate` 已完成。`W-0138` 定义了 `docs/session-postgresql-adapter-gate.md`、`docs/session-postgresql-adapter-gate.zh-CN.md` 和 `ADR-0063`。该 gate 记录 `runtime/internal/platform/persistence/postgres` 作为 `runtime/internal/app/session.Repository` 的 future adapter owner，并定义 future SQL shape、transaction handoff、error mapping、redaction 和 adapter test requirements，但不添加 PostgreSQL session adapter files、SQL execution、unit-of-work factory wiring、runtime session creation 或 validation、`RequestIdentity.SessionValidated = true`、route-policy 使用 persisted session 或 bound identity、WebSocket handshake authentication、transport credential carriers、Protobuf session messages、generated output、logout/revocation active-connection invalidation、reconnect/epoch behavior、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.session_postgresql_adapter_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-067/W-0139`，即 next-direction confirmation gate。

`M-067 Next Direction Confirmation Gate` 已完成。`W-0139` 在 session PostgreSQL adapter gate 之后选择了 `implement_session_postgresql_adapter`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-068 Session PostgreSQL Adapter Implementation` 已完成。`W-0140` 添加了 `runtime/internal/platform/persistence/postgres/session_repository.go`、`runtime/internal/platform/persistence/postgres/session_repository_test.go`、`UnitOfWork.NewSessionRepository()` 和 `ADR-0064`。该 adapter 基于 `runtime_sessions` 实现 `runtime/internal/app/session.Repository`，通过 session package 规范化输入和返回记录，将 PostgreSQL errors 映射为 redacted typed errors，并保持 transaction-neutral。它不在 login 或 BindConnection 创建 session，不验证 runtime session，不设置 `RequestIdentity.SessionValidated = true`，不改变 WebSocket handshake authentication，不添加 transport credential carriers，不添加 Protobuf session messages，不改变现有 envelope，不让 persisted session 或 bound identity 满足 route policy，不在 logout/revocation 时 invalidates active connections，不添加 reconnect/epoch behavior、cleanup jobs、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.session_postgresql_adapter_implementation` 是 repository check rule。Work queue 现在 blocked 在 `M-069/W-0141`，即 next-direction confirmation gate。

`M-069 Next Direction Confirmation Gate` 已完成。`W-0141` 在 session PostgreSQL adapter implementation 之后选择了 `define_runtime_session_validation_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-070 Runtime Session Validation Gate` 已完成。`W-0142` 定义了 `docs/runtime-session-validation-gate.md`、`docs/runtime-session-validation-gate.zh-CN.md` 和 `ADR-0065`。该 gate 规定未来 runtime session validation 由 `runtime/internal/app` 下的 application 拥有，要求在信任 persisted `runtime_sessions` row 之前必须已有 validated actor identity，定义 active/expired/revoked 和 actor-mismatch 的对外失败折叠，并保持 request identity handoff 显式。它不实现 runtime session validation，不设置 `RequestIdentity.SessionValidated = true`，不在 login 或 BindConnection 创建 session，不改变 route policy，不改变 WebSocket handshake authentication，不添加 transport credential carriers，不添加 Protobuf session messages，不改变现有 envelope，不添加 logout/revocation active-connection invalidation，不添加 reconnect/epoch behavior、cleanup jobs、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.runtime_session_validation_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-071/W-0143`，即 next-direction confirmation gate。

`M-071 Next Direction Confirmation Gate` 已完成。`W-0143` 在 runtime session validation gate 之后选择了 `implement_runtime_session_validation`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-072 Runtime Session Validation Implementation` 已完成。`W-0144` 添加了 `runtime/internal/app/runtime_session_validator.go`、`runtime/internal/app/runtime_session_validator_test.go` 和 `ADR-0066`。Application-owned `PersistentSessionValidator` 使用 `runtime/internal/app/session.Repository.FindActiveSessionByID`，在信任 persisted runtime session row 前要求 already validated player identity，验证 active/unexpired actor-player match，将对外 invalid-session failure 折叠为稳定 redacted reason，并且只在 durable validation 成功后设置 `RequestIdentity.SessionValidated = true`。它没有接入 startup 或 route policy，不创建 session，不更新 `last_seen_at`，不改变 WebSocket 或 Protobuf behavior，也不添加 logout/revocation active-connection behavior、reconnect/epoch behavior、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.runtime_session_validation_implementation` 是 repository check rule。Work queue 现在 blocked 在 `M-073/W-0145`，即 next-direction confirmation gate。

`M-073 Next Direction Confirmation Gate` 已完成。`W-0145` 在 runtime session validation implementation 之后选择了 `define_session_creation_composition_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-074 Session Creation Composition Gate` 已完成。`W-0146` 定义了 `docs/session-creation-composition-gate.md`、`docs/session-creation-composition-gate.zh-CN.md` 和 `ADR-0067`。该 gate 规定 future durable runtime session creation 由 `runtime/internal/app` 下的 application 拥有，识别 `AuthenticateWithDeviceCredential` 为第一条 future login-time composition candidate，记录 future `session.Repository.CreateRuntimeSession` 只能通过 unit-of-work capabilities 使用，保持 `access_token_record_id` 为 private server metadata，并定义 session id、lifetime、redaction 和 future test expectations。它不实现 session creation，不修改 authentication service behavior，不生成 session ids，不在 login 或 BindConnection 创建 sessions，不改变 runtime session validation 或 route policy，不改变 WebSocket 或 Protobuf behavior，也不添加 logout/revocation active-connection behavior、reconnect/epoch behavior、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.session_creation_composition_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-075/W-0147`，即 next-direction confirmation gate。

`M-075 Next Direction Confirmation Gate` 已完成。`W-0147` 在 session creation composition gate 之后选择了 `implement_session_creation_composition`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-076 Session Creation Composition Implementation` 已完成。`W-0148` 更新了 `runtime/internal/app/authentication/service.go`、`runtime/internal/app/authentication/service_test.go`、`runtime/cmd/vibit-server/main.go`、`runtime/cmd/vibit-server/main_test.go` 和 `ADR-0068`。成功的 device-credential login 现在会在同一个 unit of work 中存储 access-token record，并创建一条 linked to `access_token_record_id` 的 active durable runtime session，session id 由 server-owned generator 生成，第一版 lifetime 与 access-token 对齐。它不通过 Protobuf 暴露 session ids，不改变现有 envelope，不改变 route policy，不在 token validation 期间设置 `SessionValidated` true，不改变 WebSocket handshake authentication，不添加 transport credential carriers，不添加 logout/revocation active-connection behavior，不添加 reconnect/epoch behavior、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.session_creation_composition_implementation` 是 repository check rule。Work queue 现在 blocked 在 `M-077/W-0149`，即 next-direction confirmation gate。

`M-077 Next Direction Confirmation Gate` 已完成。`W-0149` 在 session creation composition implementation 之后选择了 `define_bound_identity_route_policy_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-078 Bound Identity Route Policy Gate` 已完成。`W-0150` 定义了 `docs/bound-identity-route-policy-gate.md`、`docs/bound-identity-route-policy-gate.zh-CN.md` 和 `ADR-0069`。该 gate 规定未来 route-policy 使用 request-token、bound-connection、session-validated、bound-session identity 时，由 `runtime/internal/app` 下的 application 拥有，并且必须 route-scoped、fail-closed、redacted。推荐的第一版实现姿态是：普通 protected domain routes 继续使用 request-level access-token proof；bound identity 和 session-validated identity 只是显式 future policy families。它不实现 route-policy 使用 bound 或 session identity，不移除 per-request token proof，不改变 WebSocket handshake authentication，不添加 transport credential carriers，不通过 Protobuf 暴露 session ids，不改变现有 envelope，不添加 logout/revocation active-connection behavior，不添加 reconnect/epoch behavior、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.bound_identity_route_policy_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-079/W-0151`，即 next-direction confirmation gate。

`M-079 Next Direction Confirmation Gate` 已完成。`W-0151` 在 bound identity route policy gate 之后选择了 `implement_bound_identity_route_policy`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-080 Bound Identity Route Policy Implementation` 已完成。`W-0152` 更新了 `runtime/internal/app/route_authentication.go`、`runtime/internal/app/route_authentication_test.go` 和 `ADR-0070`。Application route protector 现在有显式 route policy families：`public`、`request_token_required`、`bound_connection_required`、`session_validated_required` 和 `bound_session_required`。`runtime.authentication.AuthenticateWithDeviceCredential` 仍然是显式 public route；普通 protected domain routes 仍默认要求 request-level access-token proof，并且 token validation 后仍会清掉 `SessionValidated`。Bound/session identity 只能满足显式分类的 routes；metadata-only identity 对所有 protected policy families 都会被拒绝；bound-session routes 要求多个 identity source 一致。该实现不接入 WebSocket handshake authentication、transport credential carriers、Protobuf session carriers、现有 envelope changes、connection registries、frame handling 中的 persistent session validation、logout/revocation active-connection behavior、reconnect/epoch behavior、dependencies、memory durable session behavior、broader game backend behavior 或 direct Nakama/Pitaya API compatibility。`runtime.bound_identity_route_policy_implementation` 是 repository check rule。Work queue 现在 blocked 在 `M-081/W-0153`，即 next-direction confirmation gate。

`M-081 Next Direction Confirmation Gate` 已完成。`W-0153` 在 bound identity route policy implementation 之后选择了 `define_logout_revocation_active_connection_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-082 Logout Revocation Active Connection Gate` 已完成。`W-0154` 定义了 `docs/logout-revocation-active-connection-gate.md`、`docs/logout-revocation-active-connection-gate.zh-CN.md` 和 `ADR-0071`。该 gate 规定未来 logout/revocation active-connection policy 由 `runtime/internal/app` 下的 application 拥有，把 presented-token logout、runtime session revocation 和 active socket invalidation 保持为分离的 future decisions，推荐 presented-token logout 作为第一版 future scope，并要求在定位 open sockets 前先有显式 connection registry。它不实现 `LogoutAccessToken`，不 revoke tokens，不 revoke runtime sessions，不关闭 WebSocket connections，不添加 active connection registries，不添加 WebSocket close policy，不添加 Protobuf logout routes，不添加 protocol session carriers，不改变现有 envelope，不添加 reconnect/epoch behavior、cleanup jobs、dependencies、memory durable session behavior 或 direct Nakama/Pitaya API compatibility。`runtime.logout_revocation_active_connection_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-083/W-0155`，即 next-direction confirmation gate。

`M-083 Next Direction Confirmation Gate` 已完成。`W-0155` 在 logout/revocation active-connection gate 之后选择了 `define_logout_access_token_behavior_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-084 Logout Access Token Behavior Gate` 已完成。`W-0156` 定义了 `docs/logout-access-token-behavior-gate.md`、`docs/logout-access-token-behavior-gate.zh-CN.md` 和 `ADR-0072`。该 gate 规定未来 `LogoutAccessToken` behavior 由 `runtime/internal/app/authentication` 下的 application 拥有，第一版 logout scope 保持为 `presented_access_token_only`，要求 revocation 前先完成 lookup digest 和 verifier digest comparison，并要求 unit-of-work commit 后才可以返回 success。`runtime.logout_access_token_behavior_gate` 是 repository check rule。

`M-085 Next Direction Confirmation Gate` 已完成。`W-0157` 在 logout access-token behavior gate 之后选择了 `implement_logout_access_token_behavior`。`M-086 Logout Access Token Behavior Implementation` 已完成。`W-0158` 在 `runtime/internal/app/authentication/service.go` 实现了 `LogoutAccessToken`，并在 `runtime/internal/app/authentication/service_test.go` 增加聚焦测试。实现会在 unit-of-work 之前拒绝 missing 或 malformed opaque access-token proof，计算 token lookup digest，要求 active access-token posture，在 revocation 前比较 verifier digest，用 `logout_presented_access_token` 调用一次 `RevokeToken`，并且只在 commit 后返回 `LogoutStatusRevoked`。它不 revoke runtime sessions，不关闭 WebSocket connections，不添加 connection registries，不添加 WebSocket close policy，不添加 Protobuf logout routes，不添加 protocol session carriers，不改变现有 envelope，不添加 refresh、logout-all、admin revocation、reconnect/epoch behavior、cleanup jobs、dependencies、memory durable session behavior、broader game backend behavior 或 direct Nakama/Pitaya API compatibility。`runtime.logout_access_token_behavior_implementation` 是 repository check rule。Work queue 随后 blocked 在 `M-087/W-0159`，即 next-direction confirmation gate。

`M-087 Next Direction Confirmation Gate` 已完成。`W-0159` 在 logout access-token behavior implementation 之后选择了 `define_active_connection_registry_gate`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-088 Active Connection Registry Gate` 已完成。`W-0160` 定义了 `docs/active-connection-registry-gate.md`、`docs/active-connection-registry-gate.zh-CN.md` 和 `ADR-0074`。该 gate 规定 future active connection registry behavior 由 `runtime/internal/app/connection` 下的 application 拥有，第一版 posture 是 single-process、in-memory、non-durable，registry records 表示 server-observed connection state 和 validated identity linkage，而不是 client proof，并且 WebSocket transport 继续 credential-neutral。它不实现 registry，不关闭 WebSocket connections，不添加 kick/disconnect behavior，不 revoke runtime sessions，不添加 WebSocket close policy，不添加 Protobuf logout routes，不添加 protocol session carriers，不改变现有 envelope，不添加 reconnect/epoch behavior，不添加 durable/distributed registry storage、dependencies、memory durable session behavior、broader game backend behavior 或 direct Nakama/Pitaya API compatibility。`runtime.active_connection_registry_gate` 是 repository check rule。Work queue 现在 blocked 在 `M-089/W-0161`，即 next-direction confirmation gate。

`M-089 Next Direction Confirmation Gate` 已完成。`W-0161` 在 active connection registry gate 之后选择了 `implement_active_connection_registry_single_process`，并继续把 Nakama 和 Pitaya 作为 reference baselines。`M-090 Active Connection Registry Single Process Implementation` 已完成。`W-0162` 添加了 `runtime/internal/app/connection/registry.go`、`runtime/internal/app/connection/registry_test.go` 和 `ADR-0075`。该 registry 由 application 拥有，第一版是 single-process、in-memory、non-durable；它登记 server-observed connection id 和 epoch，绑定已验证 player identity 以及可选 runtime session id、access-token record id，能把记录标记为 closed 或 invalidated，能按 connection id/epoch 查找记录，并能按 player/session/token record 列出 active bound records。它不 wire startup 或 transport handoff，不关闭 WebSocket connections，不添加 kick/disconnect behavior，不 revoke runtime sessions，不替换 duplicate connections，不添加 reconnect/epoch behavior，不添加 Protobuf logout routes，不添加 protocol session carriers，不改变现有 envelope，不添加 durable/distributed registry storage、dependencies、memory durable session behavior、broader game backend behavior 或 direct Nakama/Pitaya API compatibility。`runtime.active_connection_registry_single_process_implementation` 是 repository check rule。Work queue 现在 blocked 在 `M-091/W-0163`，即 next-direction confirmation gate。

Runtime authentication implementation boundary standard 是 `docs/runtime-authentication-implementation-boundary.md`，配套简体中文译本是 `docs/runtime-authentication-implementation-boundary.zh-CN.md`。`ADR-0036` 记录该 boundary decision。未来 runtime authentication 由 `runtime/internal/app` 下的 application boundary 拥有；它必须通过 application unit-of-work boundary 使用 `authentication.Repository`，保持 PostgreSQL adapter 只负责 persistence，并在 domain dispatch 前把 validated proof 转换为 `RequestIdentity`。Token generation、verifier comparison、login execution、access-token validation、logout execution、cleanup jobs、Protobuf authentication messages、WebSocket proof carriers、generated authentication shapes 和 authentication dependencies 仍然是单独 gates。`runtime.authentication_implementation_boundary` 是该边界的 repository check rule。

`M-092 WebSocket Close Policy Gate` 已完成。`W-0164` 定义了 `docs/websocket-close-policy-gate.md`、`docs/websocket-close-policy-gate.zh-CN.md` 和 `ADR-0076`。Future WebSocket close policy 由 `runtime/internal/app` 下的 application 拥有；active connection registry 仍是 target state，不是 policy；WebSocket transport 只能在未来 narrow concrete close handoff 中执行 application policy 产生的 redacted close intent。不要把这个 gate 扩展为 transport close handoff code、close codes、close reason text、kick/disconnect behavior、logout-triggered socket close、runtime session revocation、duplicate replacement、reconnect/epoch behavior、Protobuf logout routes、protocol session carriers、现有 envelope changes、WebSocket handshake authentication、transport credential carriers、durable/distributed registry storage、dependencies、broader game backend behavior 或 direct Nakama/Pitaya API compatibility。`runtime.websocket_close_policy_gate` 是 repository check rule。

`M-094 WebSocket Close Policy Single Process Implementation` 已完成。`W-0166` 添加了 `runtime/internal/app/connection/close_policy.go`、`runtime/internal/app/connection/close_policy_test.go` 和 `ADR-0077`。该 policy 由 application 拥有，第一版是 single-process、registry-backed；它只通过 active bound registry records 按 connection id/epoch、player id、runtime session id 或 access-token record id 定位目标；它把匹配记录标记为 invalidated，并产出带有 `mark_invalidated_only` 的 redacted `CloseIntent`。它不关闭 concrete WebSocket sockets，不添加 transport close handoff，不选择 close codes 或 reason text，不添加 protocol close messages，不改变 logout behavior，不 revoke runtime sessions，不替换 duplicate connections，不添加 reconnect/epoch behavior，不添加 Protobuf logout routes，不添加 protocol session carriers，不改变 generated output，不添加 durable/distributed registry storage、dependencies、broader game backend behavior 或 direct Nakama/Pitaya API compatibility。`runtime.websocket_close_policy_single_process_implementation` 是 repository check rule。Work queue 现在 blocked 在 `M-095/W-0167`，即 next-direction confirmation gate。

Authentication generated contract shape timing standard 是 `docs/authentication-generated-contract-shape-timing.md`，配套简体中文译本是 `docs/authentication-generated-contract-shape-timing.zh-CN.md`。`ADR-0038` 记录 timing decision。Generated Go authentication contract shapes 现在已在 runtime authentication implementation boundary 之后、service interfaces 之前存在，使用 `contracts/runtime/authentication/` 作为 source，并以 `runtime/internal/generated/contracts/runtime/authentication/` 作为 output root。Generated files 保持 immutable 且 metadata-only。

Application authentication service interface boundary standard 是 `docs/application-authentication-service-interface-boundary.md`，配套简体中文译本是 `docs/application-authentication-service-interface-boundary.zh-CN.md`。`ADR-0039` 记录该 boundary decision。未来 authentication service interfaces 由 `runtime/internal/app` 下的 application 拥有；generated authentication shapes 用于指导 service-level request/result vocabulary；service behavior 只能通过 application unit-of-work boundary 使用 `authentication.Repository`；validated proof 必须在 domain dispatch 前转换为 `RequestIdentity`。该 boundary 不授权 Go service code 或 runtime authentication behavior。`runtime.application_authentication_service_interface_boundary` 是该 boundary 的 repository check rule。

`ADR-0037` 关闭 runtime authentication implementation boundary planning milestone，并打开 generated authentication contract shape timing gate。`ADR-0038` 完成 timing decision。`W-0089` 完成 generator/check support 加 metadata-only generated authentication shape output。`ADR-0039` 和 `W-0090` 完成 service-interface boundary step。`ADR-0040` 和 `W-0091` 完成 verifier algorithm/redaction step。`ADR-0041` 和 `W-0092` 完成 secret configuration/verifier key loading preparation step。`ADR-0042` 和 `W-0093` 完成 material generation preparation step。`ADR-0043` 和 `W-0094` 完成 verifier digest computation and comparison preparation step。`ADR-0044` 和 `W-0095` 完成 implementation readiness step。`ADR-0045` 和 `W-0096` 完成 local verifier key configuration loading gate。`W-0097` 完成 explicit in-memory verifier key set validator implementation slice。`ADR-0046` 和 `W-0098` 完成 environment verifier key loader gate。`W-0099` 完成 environment verifier key loader implementation slice。`ADR-0047` 和 `W-0100` 完成 token and credential material generation implementation gate。`W-0101` 完成 token and credential material generation helper implementation slice。`ADR-0048` 和 `W-0102` 完成 verifier digest helper implementation gate。`W-0103` 完成 verifier digest computation helper implementation slice。`ADR-0049` 和 `W-0104` 完成 verifier digest comparison helper gate。`W-0105` 完成 verifier digest comparison helper implementation slice。`ADR-0050` 和 `W-0106` 完成 authentication service behavior implementation gate。

Work continuation standard 是 `docs/workflow.md`，配套简体中文译本是 `docs/workflow.zh-CN.md`。机器可读 work queue 是 `.arch/work-items.yaml`。当 maintainer 说“continue”或“继续”时，应理解为推进一个 `next_ready` work item，除非当前被 blocked 或需要 confirmation。当 maintainer 要求继续多步时，应按顺序最多推进相应数量的 work items，并在 blockers、verification failures、ask-first boundaries 或 maintainer redirect 处停止。

当前可执行工具：

```bash
node tools/vibit --help
node tools/vibit check all
node tools/vibit check all --json
node tools/vibit check schemas
node tools/vibit check schemas --json
node tools/vibit check memory
node tools/vibit check memory --json
node tools/vibit check contracts
node tools/vibit check contracts --json
node tools/vibit check protocol
node tools/vibit check protocol --json
node tools/vibit check generated
node tools/vibit check generated --json
node tools/vibit check agent-tooling
node tools/vibit check agent-tooling --json
node tools/vibit check migrations
node tools/vibit check migrations --json
node tools/vibit check postgres-env
node tools/vibit check postgres-env --json
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit check work
node tools/vibit check work --json
node tools/vibit inspect module <module>
node tools/vibit inspect boundary --from <module> --to <module>
node tools/vibit inspect contract --module <module> --type <type> --id <id>
node tools/vibit inspect contracts
node tools/vibit inspect contracts --module <module>
node tools/vibit inspect generated
node tools/vibit inspect next
node tools/vibit inspect reference
node tools/vibit inspect change <change-id>
node tools/vibit inspect work
node tools/vibit inspect work --json
node tools/vibit inspect memory
node tools/vibit inspect rule <rule-id>
node tools/vibit inspect rules
node tools/vibit inspect rules --category <category>
node tools/vibit check architecture
node tools/vibit check architecture --json
node tools/vibit check change <change-id>
node tools/vibit check change <change-id> --json
node tools/vibit check module <module>
node tools/vibit check module <module> --json
node tools/vibit generate module <module>
node tools/vibit generate contract-shapes <module|all>
```

当 CLI tooling 可用时，默认使用 `node tools/vibit check all` 作为仓库验证命令。

当前 CLI 只是 Node.js standard-library tooling。不要把 CLI implementation language 当成 server runtime language。

当 agent 在 intake、verification 或 handoff 阶段需要机器可读检查结果时，使用 `--json`。

每条 JSON check result item 都应包含 `rule_id` 和 `artifact`。`artifact`、`path`、`source` 和 `output` 等 JSON fields 中的 repository-relative paths 必须在所有平台使用 forward slash，包括 Windows。把 `check all --json` 视为紧凑总览；需要完整细节时，对具体失败检查单独运行 `--json`。

当新增或修改 conversation logs 或 Agent Decision Records 时，使用 `node tools/vibit check memory`。

当新增或修改 contract source files 或 `.arch/contracts.yaml` 时，使用 `node tools/vibit check contracts`。

在创建或修改 `.proto` files、生成的 Protobuf output 或 protocol generation rules 前，使用 `node tools/vibit check protocol`。当 protocol sources 仍处于 planned 状态时，缺失 `.proto` files 可以通过；但一旦 `.proto` file 存在，它就必须与已登记的 command、query 和 event contracts 对齐。

当新增或修改 generated files、module manifest 中的 `generated` declarations、generated output standards 或 Go Protobuf generated output 时，使用 `node tools/vibit check generated`。

当新增或修改 agent-facing inspection、generation 或 verification commands 时，使用 `docs/agent-tooling.md` 和 `node tools/vibit check agent-tooling`。

当新增或修改 SQL migration sources、migration ownership manifests、migration guidance 或 persistence migration standards 时，使用 `node tools/vibit check migrations`。该检查会验证 PostgreSQL migration naming、goose markers、SQL-first boundaries、owning-module traces，以及第一版 inventory migration 的 table references。

当新增或修改 disposable PostgreSQL verification environment standards、live PostgreSQL verification guidance 或 persistence verification environment manifests 时，使用 `node tools/vibit check postgres-env`。这是静态标准检查；它不得连接 PostgreSQL，也不得要求 Docker、Podman、cloud PostgreSQL 或其他 service manager。

当新增或修改 runtime module behavior、runtime adapter boundaries、runtime guidance 或 tests 时，使用 `node tools/vibit check runtime`。在 Go runtime 尚不存在前，该检查应以 not applicable 的方式通过，因为 runtime implementation 尚未开始。当 `runtime/go.mod` 已存在但 Go source files 尚不存在时，该检查应验证 ADR-0014 skeleton 和 ADR-0018 runtime protocol adapter boundary，并且不运行 `go test` 也可以通过。一旦 Go source files 存在，runtime checks 必须要求 Go test files 和本地 Go toolchain。

当 runtime check 在 authentication、token、credential、external identity、session persistence、Protobuf envelope authentication、WebSocket handshake authentication、runtime player handler 或 WebSocket route 边界上失败时，使用 `node tools/vibit inspect rule runtime.authentication_token_session_boundary --json`。

当 runtime check 在已选 `device_credential_login`、opaque access-token、explicit request proof payload、generated authentication shape metadata boundary、authentication Protobuf deferral、WebSocket carrier deferral、schema-gate、migration、repository、adapter 或 dependency boundary 上失败时，使用 `node tools/vibit inspect rule runtime.selected_login_token_boundary --json`。

在解释 continuation request 前，使用 `node tools/vibit inspect work`。当 `.arch/work-items.yaml`、workflow docs 或 work item state 发生变化时，使用 `node tools/vibit check work`。默认 continuation unit 是一个 work item。

当 immediate continuation step 不清楚时，使用 `node tools/vibit inspect next --json`。

当 agent 在 contract、generator 或 runtime planning work 前需要完整 registered contract index 时，使用 `node tools/vibit inspect contracts --json`。

当 agent 在 intake 阶段需要以 JSON 读取单个 contract 的 registry entry、source summary、module manifest declaration 和 consistency status 时，使用 `node tools/vibit inspect contract --module <module> --type <type> --id <id>`。

在修改 generated output standards、generated output checks 或 generator behavior 前，使用 `node tools/vibit inspect generated --json`。

在根据 Nakama 或 Pitaya reference context 规划新的 game server capability families 前，使用 `node tools/vibit inspect reference --json`。

使用 `node tools/vibit generate contract-shapes all` 从 semantic contract manifests 重新生成 Go contract shape files。不要手改 `runtime/internal/generated/contracts/` 下的文件。

当 change spec 已存在，并且 agent 在 intake 或 handoff 阶段需要结构化了解它的文件、metadata、affected modules 和 verification state 时，使用 `node tools/vibit inspect change <change-id>`。

当 agent 在决定完整阅读哪些 artifacts 之前，需要结构化索引 change specs、conversation logs 和 Agent Decision Records 时，使用 `node tools/vibit inspect memory`。

使用 `rules/check-rules.json` 解读 check result 中的 `rule_id`。

当只需要单条 rule metadata 时，使用 `node tools/vibit inspect rule <rule-id>`。

使用 `node tools/vibit inspect rules --category <category>` 按 category 发现 rules。

使用 `.arch/runtime.yaml` 作为 runtime readiness 的机器可读 intake 入口。它链接了约束语言、服务器实例模型、contract 与 generation boundary、client protocol、dependency adoption，以及第一 proof slice 的 ADR。

使用 `.arch/protocol.yaml` 作为 game protocol framework decisions 的机器可读 intake 入口。它链接 `ADR-0015`，并定义第一版 WebSocket Protobuf envelope、route fields、session model、target scopes、authority rules、error model 和第一版 inventory slice protocol scope。

修改 protocol envelope、inventory Protobuf source、Buf generation configuration、generated output checks 或生成的 Go Protobuf output path 前，先阅读 `ADR-0016`、`ADR-0017`、`buf.yaml`、`buf.gen.yaml`、`proto/README.md` 和 `docs/generated-output.md`。

修改 WebSocket transport、Protobuf protocol adaptation、application dispatch、generated code 和 domain modules 之间的 runtime code 前，先阅读 `ADR-0018` 和 `docs/runtime-protocol-adapter.md`。

在添加 foundational dependencies 前，使用 `.arch/dependencies.yaml` 作为机器可读 intake 入口。Adoption records 使用 `docs/dependency-adoption.md` 和 `docs/_templates/dependency-adoption.md`。

在添加 live PostgreSQL migration checks、repository integration tests、transaction-runner integration tests 或 persistent-runtime end-to-end checks 前，使用 `docs/postgresql-verification-environment.md`。Live PostgreSQL verification 通过 `VIBIT_POSTGRES_TEST_DSN` 选择性启用；默认 repository checks 不得要求运行中的 database。

使用 `.arch/reference.yaml` 作为 Nakama/Pitaya reference alignment 的机器可读 intake 入口。Nakama 是 broad game backend product capability surface 的主要参考。Pitaya 是 Go game server framework architecture vocabulary 的主要参考。改造 reference patterns 时必须保留 vibit 的 Agent-Native constraints，并记录为什么采纳、改造或拒绝某个 reference pattern。

使用 `docs/authentication-token-session-validation.md`、`docs/authentication-proof-token-session-contract-dimensions.md`、`docs/credential-storage-external-identity-linking-boundaries.md`、`docs/session-persistence-websocket-handshake-decision-gates.md`、`docs/login-method-token-format-ratification.md`、`ADR-0023`、`ADR-0024` 和 `runtime.authentication_token_session_boundary`，再修改 authentication proof、login methods、token behavior、credential storage、external identity linking、runtime session persistence、request identity trust、Protobuf envelope authentication behavior、WebSocket handshake authentication、runtime player handlers 或 WebSocket routes。该设计标准分离 authentication proof、login methods、tokens、credentials、external identity links、runtime sessions、request identity、transport metadata、envelope metadata 和 player account lifecycle。该 dimensions standard ratify actor kinds、validation statuses、proof statuses、failure classes、retryability、request identity handoff、session error dimensions、session permission dimensions 和 validation event dimensions。Credential/external identity boundary standard 保持 player lifecycle tables 不包含 credential 和 provider subject，并延后 login-method families、credential schema、provider subject semantics、account linking、recovery、merge behavior 和 dependencies。Session/handshake gates standard 会保持 request-level、first-message、handshake-level、every-request 和 hybrid validation 为未来选择，直到它们被单独选择。Login/token ratification standard 定义如何 comparison 和选择第一批 login-method set、token model、proof carrier posture、lifecycle semantics、schema gates、checks 和 implementation queue，但不授予 implementation permission。

使用 `.arch/work-items.yaml` 作为 continuation 的机器可读 intake 入口。`W-0007` 这样的 Work item IDs 是执行步骤；ADR IDs 仍然是架构决策；change spec IDs 仍然是具体执行记录；Git hashes 仍然是 repository snapshots；versions 仍然是 release identifiers。

使用 `docs/v0.1-alpha-goal.md` 作为短期 release target intake。当前 release state 是 source-first `v0.1.0-alpha.1`，first alpha user discovery loop、feedback intake surface、product maturity milestones、prototype-ready execution plan、local development path gate、local development path package 和 storage objects behavior gate 已完成。目标是 `v0.1 alpha`，product maturity milestones 位于 `docs/product-maturity-milestones.md`，prototype-ready execution plan 位于 `docs/prototype-ready-foundation-execution-plan.md`，local development path gate 位于 `docs/prototype-ready-local-development-path-gate.md`，local development path package 位于 `docs/prototype-ready-local-development-path-package.md`，storage objects behavior gate 位于 `docs/storage-objects-behavior-gate.md`，next work item 是 `W-0202 storage_objects_persistence_schema_gate`。Alpha acceptance checklist 位于 `docs/alpha-acceptance-checklist.md`，packaged local alpha developer flow 位于 `docs/alpha-developer-flow.md`，release publishing decision gate 位于 `docs/release-publishing-decision-gate.md`，release execution preparation gate 位于 `docs/release-execution-preparation-gate.md`，release execution authorization gate 位于 `docs/release-execution-authorization-gate.md`，release execution maintainer decision 位于 `docs/release-execution-maintainer-decision.md`，release identifier plan 位于 `docs/release-identifier-artifact-plan.md`，final authorization 位于 `docs/release-execution-final-authorization.md`，first alpha user discovery loop 位于 `docs/first-alpha-user-discovery-loop.md`，first alpha feedback intake 位于 `docs/first-alpha-feedback-intake-surfaces.md`。不要把 source-first alpha authorization、user discovery loop、feedback intake、prototype-ready execution plan、local development path gate、local development path package 或 storage objects behavior gate 理解为可以添加 release binaries、packages、containers、checksums、signing/provenance artifacts、hosted deployments、direct Nakama/Pitaya API compatibility、runtime behavior、protocol routes、generated output、migrations、dependencies、repository interfaces、storage adapters、GitHub release record 之外的 public announcements、paid promotion，或跳过 `.arch/work-items.yaml`。

修改 Go runtime 文件前，先阅读 `ADR-0014`。第一版 Go module 位于 `runtime/go.mod`，module path 为 `github.com/iceiko/vibit/runtime`。Process startup 放在 `runtime/cmd/vibit-server/`，application dispatch 和 composition 放在 `runtime/internal/app/`，platform adapters 放在 `runtime/internal/platform/`，手写 domain module logic 放在 `runtime/internal/modules/<module>/`，生成的 Go contract shapes 放在 `runtime/internal/generated/contracts/`，生成的 Go Protobuf output 放在 `runtime/internal/generated/proto/`，SQL-first PostgreSQL migrations 放在 `runtime/migrations/postgres/`，Protobuf source files 放在仓库根目录的 `proto/vibit/<module>/v1/`。

## 4. 文档规则

英文是项目权威文档语言。

每个面向公众的文档都应该有：

- 英文源文档
- 简体中文可读译本

命名示例：

```text
CONSTITUTION.md
CONSTITUTION.zh-CN.md
AGENTS.md
AGENTS.zh-CN.md
docs/<name>.md
docs/<name>.zh-CN.md
.arch/README.md
.arch/README.zh-CN.md
```

规则：

- 英文源文档发生实质变更时，应在同一次变更中更新中文译本。
- 如果同一次变更无法更新译本，必须明确标记译本已过期。
- 机器可读标识符保持英文。
- 代码标识符、模块名、命令、事件、权限、错误码应使用英文，除非存在强领域理由。
- 翻译应保留意义。不要为了逐字翻译牺牲清晰度。

## 5. 标准变更工作流

每个非平凡 feature、bug fix、migration、refactor 或 standard change 都应遵循：

1. 澄清 requirement。
2. 识别 affected modules 和 contracts。
3. 当变更足够大、需要持久上下文时，编写或更新 change spec。
4. 当 public behavior 改变时，先更新 schemas、manifests 或 contracts，再实现。
5. 当 generators 存在时，用它们生成重复结构。
6. 只在声明过的边界内实现。
7. 新增或更新聚焦的 tests。
8. 运行相关 verification commands。
9. 更新文档和译本。
10. 记录已验证和未验证的内容。

对于早期纯设计变更，涉及代码、测试、生成器和验证的步骤可以不适用，但必须明确说明。

## 6. 架构规则

优先选择这样的设计：

- 给 agents 更少歧义上下文
- 创建更强 module boundaries
- 让 behavior 更容易验证
- 让 contracts 显式
- 减少 hidden coupling
- 支持 code generation
- 对人类开发者仍然实用

不要让架构规则只存在于维护者记忆里。如果一条规则重要，它最终应体现为文档、schema、manifest、test、generator 或 architecture check。

## 7. 模块规则

当 modules 存在时，每个 module 应声明：

- 它拥有什么
- 它不拥有什么
- Public commands
- Public queries
- Published events
- Subscribed events
- Allowed dependencies
- Forbidden dependencies
- Invariants
- Required tests
- Generated files
- Handwritten extension points

其他 modules 不能直接访问某个 module 的内部实现。跨模块通信应通过 commands、queries、events、public module APIs 或 generated clients 完成。

`modules/<module>/module.yaml` 应以 `docs/module-manifest.md` 作为源标准。

`changes/<date>-<change-id>/` 应以 `docs/change-spec.md` 作为源标准。

`conversations/` 应以 `docs/conversation-log.md` 作为源标准。

当 maintainer 引入 product intent、拒绝一种解释、命名一个概念或做出架构决策时，应在 conversation log 中保留该上下文。提交前必须脱敏 secrets。

`decisions/` 应以 `docs/agent-decision-record.md` 作为源标准。

当某个 decision 影响长期 architecture、generated file conventions、module ownership 或拒绝了一个合理替代方案时，应创建或更新 Agent Decision Record。Rationale 应简洁、公开；不要存储隐藏 chain-of-thought。

Generated files 对 non-system agents 不可变。如果 generated output 错了，应修改 source schema、template 或 generator，除非 change spec 或 decision record 显式授予 `generated_file_override`。

对于 server runtime，Go 是第一版 implementation language。WebSocket 是第一版 gameplay/client protocol。Protobuf 是第一版 wire message format。PostgreSQL 是第一版 authoritative durable relational store。S3-compatible object storage 是计划中的 object-storage abstraction，MinIO 是本地/自托管方向的优先候选，但必须先经过 dependency adoption record。Domain modules 不得直接依赖第三方 transport、protocol、persistence、object-storage 或 framework libraries；这些依赖由 platform adapters 通过 vibit-owned interfaces 承载。

第一批已接受的 Go runtime dependencies 记录在 `ADR-0013` 和 `.arch/dependencies.yaml` 中：

- `github.com/coder/websocket` 用于 platform WebSocket transport adapter。
- `google.golang.org/protobuf` 和 `google.golang.org/protobuf/cmd/protoc-gen-go` 用于 Go Protobuf runtime 和 generation。
- Buf CLI 用于 Protobuf linting、breaking checks、formatting 和 generation orchestration。
- `github.com/jackc/pgx/v5` 用于 PostgreSQL platform persistence adapters。
- `github.com/pressly/goose/v3` 用于 SQL-first migration tooling。
- 先使用 Go standard-library `testing`；目前尚未采纳外部 test framework。

只有声明过的 owner layers 可以直接 import 或 invoke 这些已接受依赖。Domain runtime logic 和 domain modules 必须使用 vibit-owned interfaces、generated contracts、repositories 和 adapters。

Goose migrations 应以 SQL-first 为默认。Go migrations 需要 change spec 解释为什么 SQL 不足，并且不得成为隐藏 domain business logic 的位置。

新增 Go runtime code 时，遵循 ADR-0014 的 package boundary：

- `runtime/cmd/vibit-server/` 拥有 startup、configuration wiring 和 process lifecycle。
- `runtime/internal/app/` 拥有 command/query dispatch、application service composition 和 transaction orchestration。
- `runtime/internal/platform/transport/ws/` 拥有 `github.com/coder/websocket`。
- `runtime/internal/platform/protocol/protobuf/` 拥有 Protobuf framing 和 envelope conversion。
- `runtime/internal/platform/persistence/postgres/` 拥有 `github.com/jackc/pgx/v5`。
- `runtime/internal/platform/migrations/` 拥有 `github.com/pressly/goose/v3` invocation 和 migration validation。
- `runtime/internal/platform/events/` 拥有 event recording 和 publication mechanisms。
- `runtime/internal/platform/tx/` 拥有 unit-of-work 和 transaction boundary interfaces。
- `runtime/internal/modules/<module>/` 只拥有手写 domain behavior。

Runtime protocol handoff 必须遵循 `docs/runtime-protocol-adapter.md`：WebSocket transport 读写 frames，Protobuf adapter 转换 envelopes 和 payloads，application dispatch 路由 commands 和 queries，domain modules 执行 invariants，generated packages 只提供 shapes。

State-changing commands 应通过 application dispatch 进入，并在 application-owned unit of work 中运行。Command 产生的 domain events 应记录在同一个 unit of work 中。Query handlers 不应改变状态，默认不需要 write transaction。在明确 event delivery 或 outbox decision 前，transaction 外的 event publication 继续 deferred。

在新增 persistence implementation 前，agents 必须先声明或更新相关 repository interfaces、migration expectations、transaction boundaries 和 storage verification path。未通过遵循 `ADR-0010` 和 `ADR-0011` 的 change spec 或 adoption record，不要添加 PostgreSQL drivers、migration tools、S3 SDKs 或 MinIO clients。

在添加 foundational dependencies 前，agents 必须检查 `.arch/dependencies.yaml`。在 adoption record 尚未记录 problem solved、license、maintenance activity、abstraction boundary、allowed owners、forbidden owners、replacement path 和 verification path 前，不要把 dependency slot 改为 `accepted`。

在新增 game server capability families 或 runtime subsystems 前，agents 必须检查 `.arch/reference.yaml` 和 `docs/reference-game-server-alignment.md`。先把 proposal 映射到相关 Nakama/Pitaya capability family，再让实现顺序符合 vibit 的 contract-first、manifest-first、generated 和 checkable architecture。没有明确 compatibility ADR，不要复制外部 APIs。在 modular monolith proof slice 稳定前，不要加入 Pitaya-style cluster/RPC/service-discovery work。

Decision authority 以 `ADR-0012` 为准。在 maintainer 明确授权后，agents 可以在既有已确认方向内，按专业评估决定 technical sub-decisions。但修改 constitutional principles、product direction、runtime language、primary protocol direction、persistence direction、major architecture patterns、module ownership、breaking contracts、validation 或 permission 强度，以及接受 licensing-risk、hosting、cost、operations 或 vendor-lock-in commitments 之前，仍必须询问 maintainer。

`schema/` 应以 `docs/schema-validation.md` 作为源标准。

当修改 module manifests、change specs、Agent Decision Records 或 tool JSON output 的结构时，必须更新对应 schema file，并运行 `node tools/vibit check schemas`。

## 8. 契约规则

公共行为应先定义，再实现。

带有 contract 属性的产物可以包括：

- API schemas
- Command schemas
- Query schemas
- Event schemas
- Error catalogs
- Permission catalogs
- Database migration schemas
- Generated clients

规则：

- Public contracts 必须先声明再使用。
- 对兼容性敏感的 contracts 必须 versioned。
- Breaking changes 必须显式说明。
- Generated output 必须能追踪到 source schema。
- 不要手工编辑 generated contract output，除非本次变更的目标就是 generator 本身。

## 9. 测试与验证

测试是架构的一部分，不是收尾步骤。

当实现代码存在时，相关验证可以包括：

- Unit tests
- Contract tests
- Invariant tests
- Integration tests
- Migration tests
- Replay tests
- Architecture checks
- Generator checks
- Documentation consistency checks

本仓库尚未定义最终 verification commands。在此之前，按以下方式记录验证：

```text
Verified: <commands or checks run>
Not verified: <reason>
Not applicable: <reason>
```

没有运行验证时，绝不能声称变更已经验证。

## 10. 先询问

以下情况应先询问人类维护者：

- 修改宪法原则
- 正式确定或替换项目名
- 重新定义 module ownership
- 引入新的架构模式
- 做 breaking API、command、query 或 event changes
- 修改 generated file conventions
- 删除 tests
- 削弱 validation 或 permission checks
- 在 modules 之间迁移 data ownership
- 接受有实质影响的 licensing-risk、hosting、cost、operations 或 vendor-lock-in commitments
- 修改 server runtime language、primary protocol direction、persistence direction 或 core project thesis
- 添加重大的外部框架依赖

## 11. 分层 Gate 策略

不是所有 work item 都需要同样重的 gate 成本。使用以下分层，在安全性和推进效率之间取得平衡：

### Tier 1 — Security-Critical（两步：Gate + Implementation）

适用于 cryptography、verifier algorithms、Protobuf wire format、credential schema、token lifecycle、secret configuration。

- 需要先用独立 gate milestone 定义边界，再进入实现。
- Gate 必须包含 threat model、redaction rules、fail-closed behavior 和 dependency posture。
- Implementation 作为后续独立 bounded work item 推进。

### Tier 2 — Functional Implementation（单步）

适用于 transport features、application policy、registry behavior、route registration、protocol bridge、session lifecycle、connection lifecycle。

- 使用单个 implementation work item，并把边界定义嵌入 change spec。
- 不要求单独 gate milestone。
- 仍必须包含 focused tests、`vibit check all` 和 verification record。

### Tier 3 — Lightweight（直接实现）

适用于 documentation、translation、简单 check rules、小型 tooling edits，以及已经 ratify 的机械性 migration-source updates。

- 非平凡时才需要 change spec；简单变更可以直接实现。
- 不要把新的 data ownership、新 schema semantics、新 dependencies 或新的 runtime behavior 归为 Tier 3。
- 通过 `vibit check all` 或相关静态检查完成验证即可。

### Direction Confirmation

- Direction confirmation milestones 不再要求作为独立 work item。
- 方向通过 work item 的 `ask_first` 字段和 continuation semantics 的 `recommended_direction` 管理。
- 如果方向选择足够重大，用 ADR 记录。
- 当推荐方向已经明确，且工作属于 Tier 2 时，应直接进入有边界的功能切片，而不是创建纯 confirmation milestone。

## 12. 通用禁止事项

以下禁止事项适用于所有 milestones 和 work items。不要在单个 milestone 的 `non_goals` 或 work item 的 `ask_first` 中反复重复：

禁止：

- 把 AI gameplay features 当作本项目基础
- 为方便绕过 module boundaries
- 把业务逻辑藏在 transport handlers
- 添加未登记 public events
- 添加未登记 permissions
- 添加无类型 cross-module payloads
- 在没有声明边界的情况下做大范围仓库编辑
- 无记录地手工编辑 generated files
- 英文核心文档发生实质变更后，让中文译本静默落后
- 没有运行验证却声称已验证
- 没有明确 ADR 就添加 direct Nakama/Pitaya API compatibility
- 没有 dependency adoption record 就添加 dependencies

## 13. 新增标准时

新增标准应说明：

- 要解决的问题
- 引入的规则
- 该规则为什么帮助 agents
- 对人类的影响
- 预期产物
- 验证路径
- 从现有工作的迁移路径

优先选择可以被执行和检查的小标准，而不是无法检查的宏大表述。

## 14. 新增实现代码时

不要一开始就把框架代码分散到整个仓库。

从能证明核心命题的最小完整切片开始：

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

一个好的第一版实现目标，应包含一个小而完整的后端领域，例如 player accounts、inventory、currency、rewards、quests 或 match sessions。

## 15. 自举控制

Self-bootstrapping 只有在它能改善通向可工作服务器框架的路径时才有价值。

在新增 standard、inspect command、check command、schema、generator 或 workflow rule 前，先确认它直接支持以下至少一项：

- 下一个 runtime vertical slice
- 具体 module boundary
- Public contract 或 generated shape
- Test 或 verification path
- 针对预期实现任务减少 agent context

如果收益主要只是让 tooling 本身更完整，应推迟。

当仓库已经有足够 tooling 可以尝试一个小的端到端后端能力时，优先做 runtime readiness work，而不是继续增加 meta-tooling，然后再实现 runtime slice。

Runtime readiness 只应回答让第一个 slice 成立所必需的决策：

- Implementation language 和 package layout
- Minimal server instance model
- 第一个 module 和 capability boundary
- Contract format
- Generated files 与 handwritten files 的边界
- 最小 test 和 verification strategy
- Persistence 和 migration 假设

当这些选择仍然含糊时，不要仓促进入实现。但当 readiness work 不再改变第一个 slice 如何构建、验证或维护时，也不要继续扩张准备工作。

例外情况应记录在 change spec 或 Agent Decision Record 中。

## 16. 交接要求

每次变更结束时，都要给下一个 agent 或 human 留下足够上下文。

记录：

- 改了什么
- 为什么改
- 哪些文件变更
- 哪些 contracts 或 standards 变更
- 已验证什么
- 未验证什么
- 还有哪些开放问题

如果工作未完成，说明下一个具体动作。

## 16. 当前产品同级路线图

`M-096 Nakama Pitaya Product Parity Roadmap` 已完成。`W-0168` 添加了 `docs/nakama-pitaya-product-parity-roadmap.md`、`docs/nakama-pitaya-product-parity-roadmap.zh-CN.md` 和 `ADR-0078`。维护者明确：vibit 应成为 Nakama/Pitaya-class game backend product，并覆盖 common capability families，而不只是松散参考这些项目。Product parity 表示常用能力覆盖和可用性同级；它不表示 direct API compatibility。`runtime.reference_product_parity_roadmap` 是 repository check rule。

`M-097 Protocol Logout Route Gate` 已完成。`W-0169` 添加了 `docs/protocol-logout-route-gate.md`、`docs/protocol-logout-route-gate.zh-CN.md` 和 `ADR-0079`。该 gate 用 `access_token_in_logout_request_payload` 和 `explicit_service_validated_token_lifecycle_route` posture 定义未来 `runtime.authentication.LogoutAccessToken` route semantics。它授权了 bounded W-0170 implementation slice，但不添加 socket close、runtime session revocation、active connection invalidation、reconnect behavior、protocol session carriers、dependencies 或 direct Nakama/Pitaya API compatibility。`runtime.protocol_logout_route_gate` 是 repository check rule。

`M-098 Protocol Logout Route Implementation` 已完成。`W-0170` 把已有 `LogoutAccessToken` service behavior 暴露为显式 `runtime.authentication.LogoutAccessToken` command route。该 slice 在 `proto/vibit/authentication/v1/authentication.proto` 添加了 `LogoutAccessTokenRequest` 和 `LogoutAccessTokenResponse`，通过 Buf 重新生成 Go Protobuf output，添加 protocol bridge mapping、application bootstrap route registration 和 handler、PostgreSQL startup registration、transaction bypass，以及 focused tests。Logout 会拒绝 `AuthenticatedRequest` wrapper，因此 proof 只来自 `LogoutAccessTokenRequest.access_token`。它不关闭 sockets，不 revoke runtime sessions，不 invalidate active connection records，不添加 reconnect/epoch behavior、protocol session carriers、dependencies 或 direct Nakama/Pitaya API compatibility。`runtime.protocol_logout_route_implementation` 是 repository check rule。

`M-099 Next Direction Confirmation After Protocol Logout Route Implementation` 已完成。`W-0171` 选择了 `define_transport_close_handoff_gate`，因为 protocol logout 已经可见，application close policy 也能 invalidate registry records，但还没有用于 concrete WebSocket socket close 的窄 handoff。该选择吸收 Nakama 的 lifecycle separation 和 Pitaya 的 acceptor/session/handler/connection-management layering。当前 work queue active 在 `M-100/W-0172`，next ready work item 是 `define_transport_close_handoff_gate`。不要在 direction-confirmation step 中实现 concrete socket close、close codes、close reason text、logout-triggered close、runtime session revocation、reconnect/epoch behavior、protocol session carriers、presence、chat、social modules、matchmaking、match runtime、dependencies 或 direct Nakama/Pitaya API compatibility。

`M-100 Transport Close Handoff Gate` 已完成。`W-0172` 添加了 `docs/transport-close-handoff-gate.md`、`docs/transport-close-handoff-gate.zh-CN.md` 和 `ADR-0080`。该 gate 定义未来窄 application-to-WebSocket concrete close handoff，保持 close decisions application-owned，保持 WebSocket transport credential-neutral 和 policy-neutral，并选择 server-observed `connection_id + epoch` 作为第一版 handoff target。它不实现 concrete socket close、close codes、close reason text、logout-triggered close、runtime session revocation、reconnect/epoch behavior、protocol session carriers、operations/admin disconnect、dependencies、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.transport_close_handoff_gate` 是 repository check rule。当前 work queue active 在 `M-101/W-0173`，next ready work item 是 `confirm_next_direction_after_transport_close_handoff_gate`。

`M-101 Next Direction Confirmation After Transport Close Handoff Gate` 已完成。`W-0173` 选择了 `implement_transport_close_handoff_single_process` 作为下一步 lifecycle-closure direction。当前 work queue active 在 `M-102/W-0174`，next ready work item 是 `implement_transport_close_handoff_single_process`。该 implementation 只能通过 WebSocket transport-owned handoff target server-observed `connection_id + epoch`，同时保持 application-owned close policy 和 transport credential neutrality。不要在此 slice 中选择 close codes、close reason text、logout-triggered socket close、runtime session revocation、reconnect/epoch behavior、protocol session carriers、operations/admin disconnect、dependencies、direct Nakama/Pitaya API compatibility 或 broad product modules。

`M-102 Transport Close Handoff Single Process Implementation` 已完成。`W-0174` 添加了 `runtime/internal/platform/transport/ws/close_handoff.go`、`runtime/internal/platform/transport/ws/close_handoff_test.go` 和 `ADR-0081`。WebSocket transport 现在拥有 single-process in-memory accepted socket table，并通过 server-observed `connection_id + epoch` 暴露 `RequestClose`，返回 close requested、socket not found、epoch mismatch、already closed 和 close failed 等 redacted outcomes。该实现不解析 credentials，不改变 Protobuf envelope/logout/session behavior，不选择 close codes 或 reason text，也不添加 logout-triggered socket close、runtime session revocation、reconnect behavior、protocol session carriers、operations/admin disconnect、dependencies、direct Nakama/Pitaya API compatibility 或 broad product modules。当前 work queue active 在 `M-103/W-0175`，next ready work item 是 `define_reconnect_connection_epoch_functional_slice`。

`M-103 Reconnect Connection Epoch Functional Slice` 已完成。`W-0175` 添加了 `ADR-0083`，并在 `runtime/internal/app/connection` 中实现第一版 application-owned、server-observed connection epoch progression behavior。Active connection registry 现在会把同一 connection id 下更早的 active epochs 标记为 `superseded`，记录 `superseded_at` 和 `superseded_by_epoch`，在较新 epoch 已存在后用 `connection_epoch_stale` 拒绝 stale 或 repeated epochs，保留 superseded records 供 lifecycle inspection，并把它们排除在 active target lists 之外。该 slice 不添加 reconnect tokens、resume routes、Protobuf changes、protocol session carriers、logout-triggered socket close、runtime session revocation、presence、operations/admin disconnect、dependencies、direct Nakama/Pitaya API compatibility 或 broad product modules。当前 work queue active 在 `M-104/W-0176`，next ready work item 是 `define_protocol_session_carrier_functional_slice`。

`M-104 Protocol Session Carrier Functional Slice` 已完成。`W-0176` 添加了 `ADR-0084`，并复用现有 `Envelope.Session` metadata 作为第一版 protocol-visible runtime session carrier。成功的 `runtime.authentication.AuthenticateWithDeviceCredential` 响应现在会在 response envelope session metadata 中携带 server-created runtime session id 和 authenticated player id，且不改变 Protobuf sources 或 generated output。Response envelopes 可以从 already validated application identity 派生 session metadata，但 metadata-only identity 仍保持 metadata-only，`session_id` 不是 proof。该 slice 不添加 reconnect tokens、resume routes、WebSocket handshake authentication、logout-triggered socket close、runtime session revocation、presence behavior、operations/admin disconnect、dependencies、direct Nakama/Pitaya API compatibility 或 broad product modules。当前 work queue active 在 `M-105/W-0177`，next ready work item 是 `define_presence_lifecycle_functional_slice`。

`M-105 Presence Lifecycle Functional Slice` 已完成。`W-0177` 添加了 `ADR-0085`、`runtime/internal/app/connection` 中的 registry-backed player presence snapshots、`runtime/internal/platform/transport/ws` 中的 credential-neutral WebSocket lifecycle observer，以及 `runtime/cmd/vibit-server` 下的 PostgreSQL startup composition adapters。第一版 presence behavior 从 active bound server-owned connection registry records 推导 online/offline state；成功的 first-message connection binding 可以把 validated player identity 写入 registry。该 slice 不添加 Protobuf presence messages、generated output、protocol presence queries、subscriptions、broadcasts、chat、friends、groups、parties、matchmaking、match runtime、operations/admin behavior、durable/distributed presence、reconnect/resume tokens、logout-triggered close、runtime session revocation、dependencies、direct Nakama/Pitaya API compatibility 或 broad product modules。当前 work queue active 在 `M-106/W-0178`，next ready work item 是 `define_presence_protocol_query_functional_slice`。

`M-107 v0.1 Alpha Goal And Long Term Product Target` 已作为 documentation and roadmap slice 完成。`W-0179` 添加了 `docs/v0.1-alpha-goal.md`、`docs/v0.1-alpha-goal.zh-CN.md`、`ADR-0086`，以及定义 `v0.1 alpha` 短期 developer-usable release target 和 AI-era Nakama/Pitaya-class server framework 长期目标的 conversation/change records。它没有发布 release、添加 runtime behavior、添加 protocol messages、改变 generated output、修改 migrations、添加 dependencies 或选择 direct Nakama/Pitaya API compatibility。当前 work queue 仍 active 在 `M-106/W-0178`，next ready work item 仍是 `define_presence_protocol_query_functional_slice`。

`M-108 Next Alpha Direction Selection` 已完成。`W-0180` 添加了 `ADR-0088`，并在 protected presence query 之后选择 `define_local_onboarding_device_credential_issuance_gate` 作为下一步 alpha-enabling direction。这个方向选择没有添加 runtime behavior、protocol messages、generated output、migrations、dependencies、release artifacts 或 direct Nakama/Pitaya API compatibility。由于 onboarding/device credential issuance 会触及 credential material、verifier digests、one-time raw secret presentation、player account creation 和 repository mutation ordering，下一步是 gate-only 的 `M-109/W-0181` work item，而不是直接 implementation。

`M-109 Local Onboarding Device Credential Issuance Gate` 已完成。`W-0181` 添加了 `docs/local-onboarding-device-credential-issuance-gate.md`、`docs/local-onboarding-device-credential-issuance-gate.zh-CN.md` 和 `ADR-0089`。该 gate 定义未来 local-only application service，owner 是 `runtime/internal/app/authentication`；它可以在一个 unit of work 中创建 active player account 和 active device credential record，只存 credential digests，并且只在 commit 成功后一次性返回 raw device credential text。它不授权 public signup、protocol routes、generated output、migrations、dependencies、onboarding 直接签发 access-token、onboarding 直接创建 runtime session、external identity providers、password login、account recovery、multi-device linking、release publishing、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.local_onboarding_device_credential_issuance_gate` 是 repository check rule。

`M-110 Local Onboarding Device Credential Issuance Implementation` 已完成。`W-0182` 添加了 `OnboardLocalPlayerWithDeviceCredential`、local onboarding request/result vocabulary、显式 device credential entropy 和 id generator dependencies、startup dependency composition、focused tests，以及 `ADR-0090`。该 service 用显式 entropy reader 生成 server-issued device credential material，用现有 helpers 计算 credential lookup 和 verifier digests，在同一个 unit of work 中创建 player account 和 digest-only credential record，只在 commit 后返回 raw credential text，不从 onboarding 签发 access tokens 或 runtime sessions，并保持现有 login route 不创建账号。它没有添加 public protocol routes、Protobuf sources 或 generated output、migrations、repository interface changes、dependencies、production signup、external identity、password login、recovery、multi-device linking、release artifacts、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.local_onboarding_device_credential_issuance_implementation` 是 repository check rule。

`M-111 Next Alpha Direction Selection After Local Onboarding` 已完成。`W-0183` 添加了 `ADR-0091`，并选择 `define_authenticated_gameplay_e2e_slice` 作为 local onboarding 之后的下一步 alpha-enabling direction。这个 direction selection 没有实现 authenticated gameplay E2E、添加 protocol routes、添加 generated output、改变 migrations、添加 dependencies、发布 release、选择 direct Nakama/Pitaya API compatibility、添加 production signup、external identity providers、password login、account recovery、multi-device linking 或 broad product modules。

`M-112 Authenticated Gameplay E2E Functional Slice` 已完成。`W-0184` 添加了 `ADR-0092` 和 `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`，作为 focused proof，证明 local onboarding、protocol login、first-message connection binding、protected inventory grant/read、protected presence query、logout，以及 logout 后 protected route rejection 能组合成一条 local alpha path。它没有添加 production runtime behavior、protocol routes、Protobuf sources、generated output、migrations、repository interface changes、dependencies、release artifacts、production signup、broad product modules 或 direct Nakama/Pitaya API compatibility。`runtime.authenticated_gameplay_e2e_functional_slice` 是 repository check rule。

`M-113 Runtime Runbook Alpha Path Refresh` 已完成。`W-0185` 已围绕已证明的 local alpha path 刷新 `docs/runtime-runbook.md` 和 `docs/runtime-runbook.zh-CN.md`，包括 memory vs PostgreSQL runtime posture、verifier key handling、focused authenticated gameplay E2E proof、application-service-only local onboarding 和 redaction expectations。它没有添加 runtime behavior、改变 startup configuration semantics、添加 protocol routes、改变 Protobuf sources 或 generated output、添加 migrations、添加 dependencies、发布 release、添加 production signup、broad product modules 或 direct Nakama/Pitaya API compatibility。`runtime.runtime_runbook_alpha_path_refresh` 是 repository check rule。

`M-114 Minimal Example Client Or Request Loop` 已完成。`W-0186` 添加了 `examples/local-alpha-request-loop.sh` 和 `ADR-0094`，作为 focused authenticated gameplay E2E proof 上的 minimal local alpha request-loop script。该 script 打印 redacted path summary 和 Go test status，并且没有添加 runtime behavior、改变 startup configuration semantics、添加 public protocol onboarding、改变 Protobuf sources 或 generated output、添加 migrations、添加 dependencies、发布 release、选择 production signup、添加 broad product modules 或 direct Nakama/Pitaya API compatibility。`runtime.minimal_example_client_or_request_loop` 是 repository check rule。

`M-115 Health Readiness Version Config Surface` 已完成。`W-0187` 在 `runtime/cmd/vibit-server` 添加了最小 JSON `/healthz`、`/readyz`、`/version` 和 `/configz` endpoints，并添加 focused tests 与 `ADR-0095`。该 surface 只用于 local alpha troubleshooting；它报告 redacted runtime posture，且不得暴露 verifier keys、raw credentials、raw tokens、DSNs、digests、headers、cookies、query strings、subprotocol values、remote addresses 或 concrete transport metadata。它没有添加 broad operations/admin behavior、observability dependencies、authentication/session behavior changes、startup configuration semantic changes、Protobuf sources、generated output、migrations、release artifacts、broad product modules 或 direct Nakama/Pitaya API compatibility。`runtime.health_readiness_version_config_surface` 是 repository check rule。

`M-116 Alpha Acceptance Checklist` 已完成。`W-0188` 添加了 `docs/alpha-acceptance-checklist.md`、`docs/alpha-acceptance-checklist.zh-CN.md` 和 `ADR-0096`，作为本地 v0.1 alpha acceptance checklist。它覆盖 repository intake、prerequisites、migration posture、local configuration、local onboarding posture、login、connection binding、protected inventory、presence query、logout、status endpoints、checks、redaction、contribution entry points 和 release deferrals。它没有发布 `v0.1 alpha`、添加 release packaging、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.alpha_acceptance_checklist` 是 repository check rule。

`ADR-0082` 采纳了分层 gate density strategy。Security-critical work 仍保持 gate plus implementation 两步；transport features、application policy、registry behavior、route registration、protocol bridges、session lifecycle、connection lifecycle 等 functional work 通常应作为单个 bounded work item 推进，并把边界嵌入 change spec；轻量 docs/tooling/translation/check-rule work 在简单时可以直接实现。未来 Tier 2 work 不再强制使用独立 direction confirmation milestone；用 `ask_first`、`recommended_direction` 和重大方向变化的 ADR 管理方向。`M-103/W-0175` 是该策略的第一个前向应用，应作为 bounded reconnect/connection epoch functional slice 推进，而不是纯 confirmation step。

`M-117 Package Alpha Developer Flow` 已完成。`W-0189` 添加了 `docs/alpha-developer-flow.md`、`docs/alpha-developer-flow.zh-CN.md` 和 `ADR-0097`，作为 packaged local alpha developer journey。它连接 README intake、v0.1 alpha goal、alpha acceptance checklist、runtime runbook、redacted request-loop script、status endpoints、PostgreSQL manual setup posture、redaction rules、verification commands 和下一步 contribution path。它没有发布 `v0.1 alpha`、创建 release tags 或 binaries、添加 hosted deployments、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.package_alpha_developer_flow` 是 repository check rule。

`M-118 Release Publishing Decision Gate` 已完成。`W-0190` 添加了 `docs/release-publishing-decision-gate.md`、`docs/release-publishing-decision-gate.zh-CN.md` 和 `ADR-0098`，作为 release publishing decision gate。它定义 release-publishing prerequisites、release artifact boundaries、verification requirements、stop conditions、redaction expectations 和下一步 release execution preparation direction。它没有发布 `v0.1 alpha`、创建 release tags、binaries、archives、containers、packages、checksums、hosted deployments、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.release_publishing_decision_gate` 是 repository check rule。

`M-119 Release Execution Preparation Gate` 已完成。`W-0191` 添加了 `docs/release-execution-preparation-gate.md`、`docs/release-execution-preparation-gate.zh-CN.md` 和 `ADR-0099`，作为 release execution preparation gate。它定义 release execution planning inputs、release-note input boundaries、artifact plan boundaries、maintainer approval points、verification requirements、rollback notes、stop conditions、redaction expectations 和下一步 release execution authorization direction。它没有发布 `v0.1 alpha`、创建 release tags、binaries、archives、containers、packages、checksums、provenance files、hosted deployments、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.release_execution_preparation_gate` 是 repository check rule。

`M-120 Release Execution Authorization Gate` 已完成。`W-0192` 添加了 `docs/release-execution-authorization-gate.md`、`docs/release-execution-authorization-gate.zh-CN.md` 和 `ADR-0100`，作为 release execution authorization gate。它定义 final go/no-go criteria、required verification state、release identifier review、artifact authorization boundaries、maintainer approval requirements、authorization outcome、stop conditions、redaction expectations 和 blocked next release execution maintainer decision。它没有发布 `v0.1 alpha`、选择或创建 release identifiers 或 tags、创建 release binaries、archives、containers、packages、checksums、provenance files、hosted deployments、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。`runtime.release_execution_authorization_gate` 是 repository check rule。

`M-121 Release Execution Maintainer Decision` 已完成。`W-0193` 把 maintainer decision 记录为 `go_to_release_identifier_artifact_plan`，添加了 `docs/release-execution-maintainer-decision.md`、`docs/release-execution-maintainer-decision.zh-CN.md` 和 `ADR-0101`，并添加 `runtime.release_execution_maintainer_decision` check coverage。该 decision 只允许 release execution path 继续进入 planning。它没有批准最终 release identifier、没有创建或授权 release tags、没有创建或授权 release binaries、archives、containers、packages、checksums、provenance files、hosted deployments、GitHub release records、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。

`M-122 Release Identifier And Artifact Plan` 已完成。`W-0194` 添加了 `docs/release-identifier-artifact-plan.md`、`docs/release-identifier-artifact-plan.zh-CN.md` 和 `ADR-0102`，并添加 `runtime.release_identifier_artifact_plan` check coverage。该 plan 提出 `v0.1.0-alpha.1`，记录 2026-05-21 未发现 local tag、remote origin tag 或 GitHub release record conflict，并定义 source-first future surface：Git tag、GitHub release record、release notes 和 hosting-platform source archive。它没有发布 `v0.1 alpha`、没有选择 identifier for execution、没有创建或 push release tags、没有创建 release binaries、archives、containers、packages、checksums、provenance files、hosted deployments、GitHub release records、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。

`M-123 Release Execution Final Authorization` 已完成。`W-0195` 记录了 `v0.1.0-alpha.1` 的最终 maintainer `go` authorization，添加了 `docs/release-execution-final-authorization.md`、`docs/release-execution-final-authorization.zh-CN.md`、`ADR-0103`、release notes 和 `runtime.release_execution_final_authorization` check coverage。该授权允许创建并 push Git tag `v0.1.0-alpha.1`，创建 GitHub Release `v0.1.0-alpha.1`，并且只发布 GitHub source archive。它不授权 release binaries、packages、containers、checksums、signing/provenance artifacts、hosted deployments、install scripts、registry publication、GitHub release record 之外的 public announcements、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。

`M-124 First Alpha User Discovery` 已完成。`W-0196` 添加了 `docs/first-alpha-user-discovery-loop.md`、`docs/first-alpha-user-discovery-loop.zh-CN.md`、`ADR-0104` 和 `runtime.first_alpha_user_discovery_loop` check coverage。该 loop 记录 target developer segments、outreach surfaces、feedback capture fields、review questions、success signals 和 stop conditions。它不授权 GitHub release record 之外的 public announcements、paid promotion、hosted deployments、additional release artifacts、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct Nakama/Pitaya API compatibility 或 broad product modules。

`M-125 First Alpha Feedback Intake Surfaces` 已完成。`W-0197` 添加了 `.github/ISSUE_TEMPLATE/alpha-feedback.yml`、`docs/first-alpha-feedback-intake-surfaces.md`、`docs/first-alpha-feedback-intake-surfaces.zh-CN.md`、`docs/product-maturity-milestones.md`、`docs/product-maturity-milestones.zh-CN.md`、`ADR-0105` 和 `runtime.first_alpha_feedback_intake_surfaces` check coverage。它记录 source-first alpha 已达到、prototype-ready foundation 是下一产品阶段、single-node production-candidate foundation 已规划、Nakama/Pitaya-class product 是长期目标。它不授权 broad announcements、paid promotion、hosted deployments、additional release artifacts、runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、direct compatibility 或 broad product modules。

`M-126 Prototype Ready Foundation Execution Plan` 已完成。`W-0198` 添加了 `docs/prototype-ready-foundation-execution-plan.md`、`docs/prototype-ready-foundation-execution-plan.zh-CN.md`、`ADR-0106` 和 `runtime.prototype_ready_foundation_execution_plan` check coverage。它记录 Stage 2 execution sequence、candidate work families、maturity-stage mapping、Nakama/Pitaya capability mapping、success criteria 和 stop conditions，并选择 `prototype_ready_local_development_path_gate` 作为第一项 execution slice。它不授权 runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、hosted deployments、release artifacts、public announcements、paid promotion、broad operations/admin behavior、authentication/session behavior changes、broad product modules 或 direct compatibility。

`M-127 Prototype Ready Local Development Path Gate` 已完成。`W-0199` 添加了 `docs/prototype-ready-local-development-path-gate.md`、`docs/prototype-ready-local-development-path-gate.zh-CN.md`、`ADR-0107` 和 `runtime.prototype_ready_local_development_path_gate` check coverage。它记录 supported prerequisites、startup expectations、migration expectations、configuration and secret posture、example-flow shape、allowed future write areas、verification expectations 和 stop conditions。它不授权 runtime behavior、protocol routes、Protobuf sources 或 generated output、migrations、dependencies、hosted deployments、release artifacts、public announcements、paid promotion、broad operations/admin behavior、authentication/session behavior changes、broad product modules 或 direct compatibility。

`M-128 Prototype Ready Local Development Path Package` 已完成。`W-0200` 添加了 `docs/prototype-ready-local-development-path-package.md`、`docs/prototype-ready-local-development-path-package.zh-CN.md`、`examples/README.md`、`examples/README.zh-CN.md`、`examples/local.prototype.env.example`、`.gitignore` local env guardrails、`ADR-0108` 和 `runtime.prototype_ready_local_development_path_package` check coverage。它在 W-0199 gate 内打包 setup、migration、configuration/secret redaction、example-flow 和 verification ergonomics，且不改变 runtime behavior、protocol/generated output/migrations/dependencies、release artifacts、hosted deployment、public announcements、paid promotion、broad product modules、authentication/session semantics 或 direct compatibility。

`M-129 Storage Objects Behavior Gate` 已完成。`W-0201` 定义了 `docs/storage-objects-behavior-gate.md`、`docs/storage-objects-behavior-gate.zh-CN.md` 和 `ADR-0109`。该 gate 选择第一版 player-owned small JSON storage objects，以 `owner_kind + owner_id + collection + key` 寻址，记录 read/write、permission、optimistic conflict、protocol、data、verification 和 stop-condition expectations，并添加 `runtime.storage_objects_behavior_gate` check coverage。`M-130 Storage Objects Persistence Schema Gate` active，`W-0202 Define storage objects persistence schema gate` 是 next-ready。除非后续 bounded work 明确授权，不得添加 storage objects migrations、repository interfaces、storage adapters、runtime behavior、protocol routes、Protobuf sources 或 generated output、dependencies、authentication/session changes、hosted deployments、release artifacts、public announcements、paid promotion、broad product modules、large object/blob storage、S3-compatible object storage 或 direct Nakama/Pitaya compatibility。

后续 major work 必须映射到一个 roadmap family：identity/auth/session、connection lifecycle、storage、presence/status/notifications、chat/realtime messaging、friends/groups/parties、leaderboards/tournaments、economy/progression、matchmaking、match runtime、server runtime hooks/RPC、operations、SDK/developer experience 或 distributed runtime。近期优先级现在是在 source-first release、feedback intake、execution plan 和 local path gate 之后实现 prototype-ready local development path package。在该 package 记录前，不要直接跳到 chat、social modules、matchmaking、match runtime、SDKs、distributed runtime、public announcements、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。
