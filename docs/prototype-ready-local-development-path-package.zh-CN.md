# Prototype-Ready Local Development Path Package 中文版

状态：Accepted v0.1
最后更新：2026-05-22
范围：面向 prototype authors 的 source-first local development package
依赖：`docs/prototype-ready-local-development-path-gate.md`、`docs/alpha-developer-flow.md`、`docs/runtime-runbook.md`
权威决策：`ADR-0108`

本文件是 `docs/prototype-ready-local-development-path-package.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档在 `ADR-0107` 记录的 gate 之后，把 vibit 当前 local development path 打包起来。它是 docs、examples 和 check-rule package。它不改变 production runtime behavior，不添加 protocol routes，不添加 Protobuf source 或 generated output，不添加 migrations，不添加 dependencies，不扩大 operations/admin behavior，不添加 hosted deployments，不创建 release artifacts，不执行 public announcements，不运行 paid promotion，不改变 authentication/session behavior，不添加 broad product modules，也不添加 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Prototype-ready local development path package 记录是：

```yaml
prototype_ready_local_development_path_package: implemented
completed_work_item: W-0200
decision: ADR-0108
check_rule: runtime.prototype_ready_local_development_path_package
gate_decision: ADR-0107
gate_standard: docs/prototype-ready-local-development-path-gate.md
package_standard: docs/prototype-ready-local-development-path-package.md
package_standard_translation: docs/prototype-ready-local-development-path-package.zh-CN.md
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
package_scope: docs_scripts_examples_static_checks
quick_source_path_recorded: true
prerequisite_check_recorded: true
redacted_local_configuration_template_added: true
gitignore_local_secret_guard_added: true
migration_expectations_packaged: true
runtime_startup_guidance_packaged: true
example_flow_script_recorded: true
request_loop_proof_recorded: true
verification_commands_recorded: true
stop_conditions_recorded: true
next_work_item: W-0201
next_direction: storage_objects_behavior_gate
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

Source alpha 已经有真实的 authenticated gameplay loop。本 package 把实践路径集中到一个地方，让它更容易被尝试：

```text
clone source
-> check local tools
-> prepare redacted local configuration
-> 使用 PostgreSQL 时显式 apply 或 verify migrations
-> 以 memory 或 PostgreSQL mode 启动 runtime
-> 运行 redacted authenticated request-loop proof
-> 检查 health、readiness、version 和 config posture
-> 继续下一个 bounded work item
```

这个 package 有意保持诚实：它不会假装 source checkout 已经是 packaged product distribution，不会假装 PostgreSQL setup 会自动完成，也不会假装 local onboarding 已经是 public protocol surface。

## 3. Fastest Path

在 source checkout 中运行：

```bash
node tools/vibit inspect next
node tools/vibit check all --json
cd runtime && go test ./...
cd .. && examples/local-alpha-request-loop.sh
```

这条路径证明 repository checks、Go tests 和 redacted authenticated request-loop proof。它不需要 live PostgreSQL server、committed verifier keys、raw credentials、raw access tokens、generated SDK、Docker、hosted infrastructure 或 package installation。

## 4. Supported Prerequisites

第一条 supported local path 假设：

- Go，用于 runtime tests 和 `go run ./cmd/vibit-server`；
- Node.js，用于 `tools/vibit`；
- POSIX-compatible shell，用于 repository-owned local scripts；
- PostgreSQL，仅在评估 persistent runtime path 时需要；
- Buf 和 Protobuf tooling，仅在重新生成 Protobuf output 时需要。

这个 package 不要求 Docker、Docker Compose、Kubernetes、cloud services、external secret managers、package registries、hosted databases 或 paid services。

## 5. Local Configuration

已提交的 placeholder template 是：

```text
examples/local.prototype.env.example
```

只把它当字段 checklist 使用。它包含 placeholders，不包含可用 secret values。

本地 private environment files 会被 `.gitignore` 忽略：

```text
.vibit.local.env
.env.local
.env.*.local
```

不要提交、粘贴、记录日志，或放入 ADR、change records、examples、screenshots、issue reports、shell transcripts：

- raw device credential text 或 bytes；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- HMAC input 或 output bytes；
- verifier key values；
- concrete verifier key set ids；
- 带 credentials 的 PostgreSQL DSNs；
- 可能携带 secrets 的 headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 transport metadata。

在 commit 或 push 前，应检查明显的 GitHub token 泄漏：

```bash
rg -n -o "ghp_[A-Za-z0-9_]{20,}|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .
```

该命令不得读取或打印 private local env contents。

## 6. Migration And Startup Path

两种 local startup posture 仍然是：

```text
VIBIT_RUNTIME_STORE=memory
VIBIT_RUNTIME_STORE=postgres
```

Memory store 是默认 quick smoke posture。它适合检查 basic runtime process 和最初的 in-memory inventory request loop。

PostgreSQL store 是当前 alpha composition，覆盖 persistent inventory、player accounts、device credential login、token issuance、runtime sessions、request-level protected routes、first-message connection binding、logout 和 protected presence query。

在 fresh local database 上使用 `VIBIT_RUNTIME_STORE=postgres` 前：

1. 准备 local PostgreSQL database。
2. 显式 apply 或 verify `runtime/migrations/postgres/` 下的 SQL migration sources。
3. 把 DSN 保存在 private local configuration 中。
4. 用 local-only values 设置必需的 verifier key variables。
5. 启动 `cd runtime && go run ./cmd/vibit-server`。

普通 server startup 不会自动 apply migrations。这个 package 不添加 automatic startup migration apply behavior、schema changes、repository changes 或 storage adapter changes。

## 7. Example Flow

Package 记录的第一条 executable proof path 是：

```bash
examples/local-alpha-request-loop.sh
```

该 script 包装：

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

证明覆盖：

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

它使用现有 runtime protocol handlers 和 test-owned setup。它不是 public onboarding client、product SDK、live PostgreSQL process client、hosted demo、release artifact 或 compatibility promise。

## 8. Local Status Surfaces

Runtime process 运行时，检查：

```text
/healthz
/readyz
/version
/configz
```

这些是 local troubleshooting endpoints。`/configz` 报告 redacted posture 和 `secrets_redacted: true`。它不得打印 verifier keys、raw credentials、raw tokens、DSNs、digests、headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 concrete transport metadata。

## 9. Verification

默认 verification 保持 source-first 和 local：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change implement-prototype-ready-local-development-path-package --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
cd .. && examples/local-alpha-request-loop.sh
git diff --check
```

Optional live PostgreSQL verification 仍然通过 `VIBIT_POSTGRES_TEST_DSN` 和 disposable database opt-in。默认 checks 不得要求 live PostgreSQL 或 private secrets。

## 10. Stop Conditions

如果继续改进这个 package 需要以下内容，应停止并请求 maintainer authorization：

- production runtime behavior changes；
- protocol route changes；
- Protobuf source 或 generated output changes；
- SQL migrations、repository interface changes 或 storage adapter changes；
- new dependencies；
- automatic startup migration apply behavior；
- broad operations/admin behavior；
- authentication/session semantic changes；
- public local onboarding routes 或 new proof carriers；
- broad product module expansion；
- direct Nakama/Pitaya API compatibility；
- hosted deployments 或 demos；
- release binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry publications、SDK packages 或 additional release artifacts；
- public announcements beyond the GitHub release record；
- paid promotion；
- handling 或 disclosure of secrets。

## 11. Next Work

下一条 bounded direction 是：

```text
W-0201 Define storage objects behavior gate
```

该工作应定义 inventory proof slice 之外的第一版 general storage-object behavior。除非另行授权，它应保持为 gate，不实现 runtime behavior、protocol、generated output、migrations、dependencies、operations/admin breadth、authentication/session changes、broad product modules 或 direct Nakama/Pitaya compatibility。
