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

本仓库当前处于宪法和标准设计阶段。

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

框架实现代码、generators、modules 和 verification commands 可能尚不存在。如果它们不存在，应记录 verification 当前不可用，而不是假装已经运行。

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

Runtime authentication implementation boundary standard 是 `docs/runtime-authentication-implementation-boundary.md`，配套简体中文译本是 `docs/runtime-authentication-implementation-boundary.zh-CN.md`。`ADR-0036` 记录该 boundary decision。未来 runtime authentication 由 `runtime/internal/app` 下的 application boundary 拥有；它必须通过 application unit-of-work boundary 使用 `authentication.Repository`，保持 PostgreSQL adapter 只负责 persistence，并在 domain dispatch 前把 validated proof 转换为 `RequestIdentity`。Token generation、verifier comparison、login execution、access-token validation、logout execution、cleanup jobs、Protobuf authentication messages、WebSocket proof carriers、generated authentication shapes 和 authentication dependencies 仍然是单独 gates。`runtime.authentication_implementation_boundary` 是该边界的 repository check rule。

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

## 11. 禁止事项

禁止：

- 把 AI gameplay features 当作本项目基础
- 为方便绕过 module boundaries
- 把业务逻辑藏在 transport handlers
- 添加未登记 public events
- 添加未登记 permissions
- 添加无类型 cross-module payloads
- 在没有声明边界的情况下做大范围仓库编辑
- 无记录地手工编辑 generated files
- 英文公共文档发生实质变更后，让中文译本静默落后
- 没有运行验证却声称已验证

## 12. 新增标准时

新增标准应说明：

- 要解决的问题
- 引入的规则
- 该规则为什么帮助 agents
- 对人类的影响
- 预期产物
- 验证路径
- 从现有工作的迁移路径

优先选择可以被执行和检查的小标准，而不是无法检查的宏大表述。

## 13. 新增实现代码时

不要一开始就把框架代码分散到整个仓库。

从能证明核心命题的最小完整切片开始：

```text
requirement -> spec -> contract -> generated shape -> handwritten logic -> tests -> verification -> docs
```

一个好的第一版实现目标，应包含一个小而完整的后端领域，例如 player accounts、inventory、currency、rewards、quests 或 match sessions。

## 14. 自举控制

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

## 15. 交接要求

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
