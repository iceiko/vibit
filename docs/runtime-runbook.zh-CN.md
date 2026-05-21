# Runtime Runbook 中文版

状态：Draft v0.2
最后更新：2026-05-21
范围：第一版 Go runtime process startup、local alpha path 和 manual verification
说明：本文件是 `docs/runtime-runbook.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本 runbook 记录如何启动第一版 vibit Go runtime process，以及当前 local authenticated gameplay alpha path 的形态。

## Current Runtime Surface

当前 runtime process 挂载一个 gameplay WebSocket endpoint：

```text
/v1/ws
```

该 endpoint 期望接收 binary WebSocket messages，其中 payload 是 `vibit.protocol.v1.Envelope` Protobuf bytes。

Text WebSocket messages 会被 transport adapter 拒绝。该 endpoint 不接受 JSON。

该 process 也挂载小型 JSON troubleshooting endpoints：

```text
/healthz
/readyz
/version
/configz
```

`/healthz` 报告 process 是否存活。`/readyz` 报告 ready status、选中的 runtime store posture 和 WebSocket path。`/version` 报告 pre-alpha runtime version。`/configz` 只报告 redacted posture：runtime store、WebSocket path、local alpha request-loop script path、PostgreSQL configuration 是否存在、所选 composition 是否要求 authentication configuration，以及 `secrets_redacted: true`。

这些 endpoints 是 local alpha troubleshooting surfaces，不是 production operations API、admin console、metrics backend、protocol route 或 release packaging surface。

当前有两种 startup composition：

```text
VIBIT_RUNTIME_STORE=memory
```

Memory store 是默认 bootstrap path。它只 wire 最初的 in-memory inventory request loop：

```text
WebSocket binary frame
-> Protobuf envelope
-> application dispatch
-> inventory command or query handler
-> Protobuf response envelope
-> WebSocket binary frame
```

它不会 wire 当前 authentication service、route protection、connection binding、logout、runtime sessions 或 presence query。

```text
VIBIT_RUNTIME_STORE=postgres
```

PostgreSQL store 是当前 alpha runtime composition。它会 wire：

- PostgreSQL-backed inventory state。
- Player account lifecycle persistence。
- Device credential login 和 opaque access-token issuance。
- Login 期间的 runtime session creation。
- 通过 `AuthenticatedRequest` 做 request-level access-token route protection。
- 通过 `runtime.authentication.BindConnection` 做 first-message connection binding。
- 通过 `runtime.authentication.LogoutAccessToken` 做 logout。
- Registry-backed presence lifecycle，以及 protected `runtime.presence.GetPlayerPresence` query。

Local onboarding 现在作为 application service method `OnboardLocalPlayerWithDeviceCredential` 存在，但它不是 public WebSocket、Protobuf、HTTP、CLI 或 startup auto-creation surface。它已经通过 tests 证明，并预留给未来 local tooling 使用。

## Start The Server

从 Go runtime module 启动：

```bash
cd runtime
go run ./cmd/vibit-server
```

默认 listen address 是：

```text
:8080
```

可以这样覆盖：

```bash
VIBIT_ADDR=:9090 go run ./cmd/vibit-server
```

默认 runtime store 是 in memory：

```text
VIBIT_RUNTIME_STORE=memory
```

如果要启动显式 PostgreSQL-backed inventory composition path，需要同时提供 store selector 和 PostgreSQL DSN：

```bash
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

PostgreSQL runtime path 也会 wire 当前 authentication、token、runtime session、route protection、connection binding、logout 和 presence-lifecycle composition。它需要 authentication verifier key environment variables：

```text
VIBIT_AUTH_VERIFIER_KEY_SET_ID
VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY
VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY
VIBIT_AUTH_TOKEN_LOOKUP_KEY
VIBIT_AUTH_TOKEN_VERIFIER_KEY
```

Verifier key values 必须是 runtime loader 接受的 Base64 text。不要提交 local verifier keys。

Verifier key requirements：

- `VIBIT_AUTH_VERIFIER_KEY_SET_ID` 必须非空。
- 每个 logical verifier key 解码后至少 32 bytes。
- 四个 logical verifier keys 必须互不相同。
- Weak repeated-byte keys 会被拒绝。
- Key values 和 concrete key set ids 都不是 log-safe。按 secrets 处理。

Runtime loader 接受 URL-safe unpadded Base64 和 standard padded Base64 key text。Verifier keys 应保存在 local environment configuration 或显式 local secret source 中；不要提交到 repository、shell history、ADR、change records、test fixtures 或 runbook examples。

可选 authentication settings：

```text
VIBIT_AUTH_ACCESS_TOKEN_TTL
VIBIT_AUTH_TOKEN_AUDIENCE
```

可选 PostgreSQL pool settings：

```text
VIBIT_POSTGRES_MAX_CONNS
VIBIT_POSTGRES_MIN_CONNS
```

普通 server startup 不会自动 apply migrations。在 fresh database 上使用 PostgreSQL store path 前，必须显式 apply 或 verify migrations。

## Local Alpha Flow

现在已经证明的 local alpha path 是：

```text
local onboarding
-> device credential login
-> connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

可执行证明是：

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

Minimal local alpha request-loop script 是：

```bash
examples/local-alpha-request-loop.sh
```

该 script 是 focused E2E proof 的 redacted wrapper。它不是 public onboarding client、product SDK、live PostgreSQL process client 或 release artifact。

该 test 对 login、binding、protected inventory、protected presence 和 logout 使用与 runtime 相同的 protocol frame handler 形态。它使用 test-local repositories、deterministic entropy、deterministic ids 和 in-memory registry，因此不需要 live PostgreSQL server，也不需要提交 secrets。

当前 alpha path details：

1. Local onboarding 通过 `OnboardLocalPlayerWithDeviceCredential` 创建 active player account 和 active digest-only device credential record。
2. Local onboarding 只会在 unit-of-work 成功后一次性返回 raw device credential text。它不会返回 access token 或 runtime session。
3. `runtime.authentication.AuthenticateWithDeviceCredential` 接收 device credential proof，并在验证 stored credential verifier 后返回 opaque access token。
4. Login response envelope 会在 `Envelope.Session` 中携带 runtime session metadata。Session id 是 metadata，不是 authentication proof。
5. `runtime.authentication.BindConnection` 使用 access token，把 server-observed connection id 和 epoch 绑定到 validated player identity。当前仍不使用 WebSocket handshake authentication。
6. `inventory.GrantItem`、`inventory.GetInventory` 和 `runtime.presence.GetPlayerPresence` 这类 protected routes 需要 `vibit.authentication.v1.AuthenticatedRequest`。
7. `runtime.authentication.LogoutAccessToken` 只 revoke presented access token。
8. 后续如果继续用同一个 revoked token 访问 protected route，会被拒绝。

当前 public protocol 不包含 local onboarding route。直接运行 PostgreSQL server process 的开发者，仍需要未来 local tool、request-loop script 或受控 seed path 来调用 local onboarding 并获得第一个 device credential。

Packaged local alpha developer journey 是：

```text
docs/alpha-developer-flow.md
```

它把 README、本 runbook、redacted request-loop script、local status endpoints、acceptance checklist、PostgreSQL manual setup posture、verification commands 和下一步 contribution path 连接起来。Packaged flow 现在已经存在，`docs/release-publishing-decision-gate.md` 已定义 release publishing decision boundary，`docs/release-execution-preparation-gate.md` 已定义 release execution preparation boundary，`docs/release-execution-authorization-gate.md` 已定义 release execution authorization gate criteria，`docs/release-execution-maintainer-decision.md` 已记录 maintainer go decision 以继续 planning，`docs/release-identifier-artifact-plan.md` 已记录 proposed `v0.1.0-alpha.1` identifier 和 source-first artifact plan；下一处缺口是 final maintainer authorization，而不是 release publishing 或 artifact creation 本身。

Local alpha acceptance checklist 是：

```text
docs/alpha-acceptance-checklist.md
```

它记录哪些 alpha items 是 ready、manual、deferred 或 blocked。它不是 release declaration，也不授权 release packaging。

## Manual Verification Paths

### Bootstrap Memory Path

1. 启动 server。
2. 连接 WebSocket client 到 `ws://127.0.0.1:8080/v1/ws`。
3. 对 bootstrap in-memory path，发送 `inventory.GrantItem` 或 `inventory.GetInventory` 的 binary Protobuf `Envelope`。
4. 确认 response 是 binary Protobuf `Envelope`，并且带有相同的 `request_id`。

### Local Alpha Request Loop

运行 minimal local request-loop script：

```bash
examples/local-alpha-request-loop.sh
```

该 script 只打印 redacted path summary 和 Go test status。它不需要 live PostgreSQL、committed verifier keys、raw credentials、raw access tokens 或手写 WebSocket client。

### Authenticated Alpha Proof

运行 focused E2E test：

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

这是当前完整 local alpha flow 的最佳证明。它不需要 live PostgreSQL，也不会打印 raw credentials 或 tokens。

### Alpha Acceptance Checklist

阅读当前 checklist：

```bash
sed -n '1,220p' docs/alpha-acceptance-checklist.md
```

在 packaging 或 publishing work 前，确认 ready、manual、deferred 和 blocked 状态。该 checklist 保持 release publishing、release packaging、public local onboarding、production signup、broad operations/admin behavior、broad product modules 和 direct Nakama/Pitaya API compatibility deferred。

### PostgreSQL Runtime Path

对 PostgreSQL process path：

1. 准备 local PostgreSQL database。
2. 显式 apply SQL migrations。普通 server startup 不会自动 apply migrations。
3. 提供 `VIBIT_RUNTIME_STORE=postgres`、`VIBIT_POSTGRES_DSN` 和全部 verifier key environment variables。
4. 启动 `go run ./cmd/vibit-server`。
5. 通过 `ws://127.0.0.1:8080/v1/ws` 发送 binary Protobuf envelopes。

Process path 当前暴露 login、binding、protected inventory、protected presence 和 logout routes。它不把 local onboarding 暴露为 public protocol route。

## Current Runtime Assumptions

- Runtime 默认使用 in-memory inventory repository。
- `VIBIT_RUNTIME_STORE=postgres` 会启用显式 PostgreSQL composition path，覆盖 persistent inventory、player account、authentication token/credential、runtime session、route protection、logout、connection binding 和 presence-lifecycle wiring。
- Inventory bootstrap permissions 允许 grant 和 read operations。
- Local onboarding/device credential issuance 已作为 application service method 存在，并已在 tests 中证明。它不是 public protocol behavior。
- Authenticated gameplay E2E path 已通过 focused Go protocol test 使用 existing capabilities 证明。
- `examples/local-alpha-request-loop.sh` 是该 proof 上的 minimal local alpha request-loop script。
- `/healthz`、`/readyz`、`/version` 和 `/configz` 为 startup troubleshooting 提供最小 local status surface。
- PostgreSQL persistence 只在显式选择时启用。Persistence boundary 已定义在 `docs/postgresql-persistence-boundary.md`。
- 普通 server startup 不会自动 apply PostgreSQL migrations。
- Optional live PostgreSQL verification 定义在 `docs/postgresql-verification-environment.md`；它要求 `VIBIT_POSTGRES_TEST_DSN`，且不属于默认 server startup。
- Generated route registration 尚未实现；route registration 仍是 handwritten startup/bootstrap code。
- v0.1 alpha path 已有 `docs/alpha-acceptance-checklist.md` 作为 alpha acceptance checklist。
- v0.1 alpha path 仍需要把 alpha developer flow 整理成一条 coherent local developer journey。
- Production signup、external identity providers、password login、account recovery、account merge、multi-device linking、direct Nakama/Pitaya API compatibility 和 broad product modules 仍然 deferred。

这些是第一版 request loop 的 bootstrap assumptions，不是长期 production policy。

## Redaction Rules

不要记录或提交：

- Raw device credential text 或 bytes。
- Raw access tokens。
- Credential 或 token lookup digests。
- Credential 或 token verifier digests。
- HMAC input 或 output bytes。
- Verifier key values。
- Concrete verifier key set ids。
- 带 credentials 的 PostgreSQL DSNs。
- 可能携带 secrets 的 headers、cookies、query strings、WebSocket subprotocol values 或 remote addresses。

Documentation、ADRs、change specs、logs、test output 和 examples 应使用 `redacted-token-text`、`redacted-device-credential` 或 `postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable` 这类 placeholder；只有在明确是非生产 sample text 时才使用。

## Verification Commands

除特别说明外，从仓库根目录运行：

```bash
cd runtime && go test ./...
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
cd runtime && go vet ./...
node tools/vibit check runtime
node tools/vibit check postgres-env
node tools/vibit check all
```

`node tools/vibit check postgres-env` 是静态标准检查。它不会连接 PostgreSQL。Live PostgreSQL verification 继续通过 `VIBIT_POSTGRES_TEST_DSN` 选择性启用。

使用下面的 command 针对 disposable PostgreSQL database 运行当前 live durable inventory verification：

```bash
cd runtime && VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

该 test 会显式 apply inventory migration，并通过 PostgreSQL-backed runtime composition 验证 WebSocket Protobuf `GrantItem` 后接 `GetInventory` 的 request loop。如果未设置 `VIBIT_POSTGRES_TEST_DSN`，test 会 skip，并记录 live PostgreSQL verification 不可用。

该 test 默认使用 `drop_schema` cleanup semantics。其他 cleanup modes 会被该 test 有意 skip，因为 migration apply 必须从 clean schema 验证。
