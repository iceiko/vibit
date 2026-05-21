# Alpha Developer Flow 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：vibit v0.1 alpha path 的 packaged local developer journey
说明：本文件是 `docs/alpha-developer-flow.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档把现有 local alpha entry points 整理成一条 developer journey。它不是 release declaration，也不授权 release publishing、release packaging、hosted deployment、runtime behavior changes、protocol changes、generated output changes、migrations、dependencies、broad operations/admin behavior、product module expansion 或 direct Nakama/Pitaya API compatibility。

## 1. 目的

Local alpha path 现在已经有一个具备技术能力的 contributor 检查 vibit 所需的组件：

- `README.md` 中的 project positioning；
- `docs/runtime-runbook.md` 中的 runtime startup 和 verification notes；
- `examples/local-alpha-request-loop.sh` 中的 redacted request-loop script；
- `/healthz`、`/readyz`、`/version` 和 `/configz` local status endpoints；
- `docs/alpha-acceptance-checklist.md` 中的 acceptance criteria；
- `.arch/work-items.yaml` 中的 continuation state。

本文档把这些组件连接起来，让 contributor 不需要在 project memory 中来回查找也能沿着同一条 sequence 操作。

## 2. 当前 Package State

```text
local_alpha_developer_flow_packaged: true
release_declared: false
release_publishing_authorized_by_this_document: false
release_packaging_authorized_by_this_document: false
release_publishing_decision_gate_defined: true
release_execution_preparation_gate_defined: true
release_execution_authorization_gate_defined: true
next_direction: release_execution_maintainer_decision
next_work_status: blocked
```

Repository 仍是 pre-alpha。Packaged flow 已准备好进行 local review，但发布 `v0.1 alpha` 和执行任何 release step 仍 blocked，等待 maintainer 明确 go/no-go decision。

## 3. 推荐 Journey

1. 阅读 `README.md`，理解 vibit、pre-alpha 状态，以及 agent-native server framework 目标。
2. 阅读 `docs/v0.1-alpha-goal.md`，了解短期 `v0.1 alpha` 目标和长期 Nakama/Pitaya-class 方向。
3. 阅读 `docs/alpha-acceptance-checklist.md`，确认哪些 alpha items 是 ready、manual、deferred 或 blocked。
4. 安装 local prerequisites：Go、Node.js；测试 persistent runtime path 时需要 PostgreSQL；只有重新生成 Protobuf output 时才需要 Buf/Protobuf tooling。
5. 运行 static repository checks：

```bash
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check all --json
```

6. 运行 Go tests：

```bash
cd runtime
go test ./...
```

7. 从 repository root 运行 redacted local alpha request loop：

```bash
examples/local-alpha-request-loop.sh
```

8. 手动启动 server 或评估 PostgreSQL runtime path 时，使用 `docs/runtime-runbook.md`。
9. 使用 `.arch/work-items.yaml` 和 `node tools/vibit inspect next` 找到下一步 bounded contribution。

## 4. Runtime Entry Points

当前 process 暴露：

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

`/v1/ws` 是 gameplay WebSocket endpoint。它期望 binary `vibit.protocol.v1.Envelope` Protobuf bytes。该 endpoint 不接受 JSON。

`/healthz`、`/readyz`、`/version` 和 `/configz` 是 local alpha troubleshooting endpoints。`/configz` 报告 redacted runtime posture，并包含 `secrets_redacted: true`。这些 endpoints 不是 production operations API、admin console、metrics backend、gameplay protocol route、release artifact 或 hosted deployment surface。

## 5. Local Proof Flow

Packaged local proof flow 是：

```text
local onboarding
-> device credential login
-> connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

可执行入口是：

```bash
examples/local-alpha-request-loop.sh
```

该 script 包装 focused Go E2E proof：

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

该 proof 使用现有 runtime protocol handlers。它不需要 live PostgreSQL、committed verifier keys、raw credentials、raw access tokens、DSNs、digests 或手写 WebSocket client。

## 6. PostgreSQL Path

PostgreSQL runtime path 是当前 alpha runtime composition，但它有 manual setup requirements：

- 准备 local PostgreSQL database；
- 显式 apply 或 verify SQL migrations；
- 设置 `VIBIT_RUNTIME_STORE=postgres`；
- 设置 `VIBIT_POSTGRES_DSN`；
- 设置全部 authentication verifier key environment variables；
- 不提交 local verifier keys 或 DSNs。

普通 server startup 不会自动 apply migrations。Optional live PostgreSQL verification 仍通过 `VIBIT_POSTGRES_TEST_DSN` 和 disposable database opt-in。

## 7. Redaction Contract

不要记录或提交：

- raw device credential text 或 bytes；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- HMAC input 或 output bytes；
- verifier key values；
- concrete verifier key set ids；
- 带 credentials 的 PostgreSQL DSNs；
- 可能携带 secrets 的 headers、cookies、query strings、WebSocket subprotocol values 或 remote addresses。

Request-loop script 和 `/configz` surface 都属于该 redaction posture。

## 8. Contribution Entry Point

下一步 contribution path 始终是 machine-readable：

```bash
node tools/vibit inspect next
```

Release execution authorization gate 现在已经定义在 `docs/release-execution-authorization-gate.md`。下一步 work blocked 在 `W-0193 Confirm release execution maintainer decision`；maintainer 必须明确选择 go/no-go，并授权任何 release identifier、tag、artifact、hosted deployment、publication surface 或 release execution command 后，release execution 才能继续。

## 9. Deferred Work

以下内容继续 deferred，直到后续明确 work item：

- 发布 `v0.1 alpha`；
- 选择 release identifiers；
- 创建 release tags、binaries、archives、containers、packages、checksums、provenance files 或 hosted deployments；
- 添加 public local onboarding protocol route；
- 添加 production signup、external identity providers、password login、account recovery、account merge 或 multi-device linking；
- 添加 broad operations/admin behavior、metrics backend integration 或 production observability；
- 添加 chat、friends、groups、parties、matchmaking、match runtime、SDKs、distributed runtime 或 direct Nakama/Pitaya API compatibility。

## 10. Verification

检查 packaged flow 时使用这组命令：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```
