# Prototype-Ready Local Development Path Gate

状态：Accepted v0.1
最后更新：2026-05-22
范围：在 prototype-ready implementation packaging 之前，为可重复本地开发路径定义 gate
依赖：`docs/prototype-ready-foundation-execution-plan.md`、`docs/alpha-developer-flow.md`、`docs/runtime-runbook.md`
权威决策：`ADR-0107`

本文是 `docs/prototype-ready-local-development-path-gate.md` 的简体中文译本。英文文件是权威版本。

本文定义让 vibit 的本地开发路径足够可重复、足够适合 prototype 作者使用的 gate。它是一个 gate artifact。它不实现 runtime behavior、不添加 protocol routes、不添加 Protobuf source 或 generated output、不添加 migrations、不添加 dependencies、不扩大 operations/admin behavior、不添加 hosted deployments、不创建 release artifacts、不执行 public announcements、不运行 paid promotion、不改变 authentication/session behavior、不添加 broad product modules，也不添加 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Prototype-ready local development path gate 的记录是：

```yaml
prototype_ready_local_development_path_gate: defined
completed_work_item: W-0199
decision: ADR-0107
check_rule: runtime.prototype_ready_local_development_path_gate
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
source_execution_plan: docs/prototype-ready-foundation-execution-plan.md
gate_standard: docs/prototype-ready-local-development-path-gate.md
gate_standard_translation: docs/prototype-ready-local-development-path-gate.zh-CN.md
future_implementation_work_item: W-0200
future_implementation_direction: prototype_ready_local_development_path_package
supported_prerequisites_recorded: true
startup_expectations_recorded: true
migration_expectations_recorded: true
configuration_secret_posture_recorded: true
example_flow_shape_recorded: true
allowed_future_write_areas_recorded: true
verification_expectations_recorded: true
stop_conditions_recorded: true
planning_only: true
gate_only: true
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

下一个产品问题不只是 vibit 是否能通过检查。Prototype 作者需要一条本地路径，它必须可理解、可重复，并且诚实说明哪些部分仍然是手动步骤。

本地开发路径应该让一个有技术能力的开发者能够：

- 安装需要的本地工具；
- 配置安全的本地 secrets，且不把它们提交进仓库；
- 显式准备或验证 PostgreSQL schema；
- 用 memory mode 启动 runtime 做快速 smoke path；
- 在评估当前 alpha composition 时启动 PostgreSQL runtime path；
- 运行一个有意义的 authenticated flow；
- 检查 health、readiness、version 和 redacted configuration posture；
- 知道哪些缺口是刻意 deferred，而不是靠猜。

这个 gate 保持 source-first 和 local。它不会把 vibit 变成 production deployment、hosted demo、binary distribution、SDK package，也不会把 vibit 变成 Nakama 或 Pitaya 的 compatibility clone。

## 3. Supported Local Prerequisites

第一版本地路径可以假设：

- Go，用于 runtime tests 和 `go run ./cmd/vibit-server`；
- Node.js，用于 `tools/vibit` checks；
- POSIX-compatible shell，用于仓库拥有的 local scripts；
- PostgreSQL，用于 persistent prototype path；
- Buf 和 Protobuf tooling，仅在重新生成 Protobuf output 时需要。

第一版 package 不应要求 Docker、Docker Compose、Kubernetes、cloud infrastructure、external secret managers、package registries、hosted databases 或 paid services。这些以后可能有用，但它们会改变 operations 和 dependency posture，因此需要单独授权。

## 4. Startup Expectations

本地 package 应该记录两种 startup posture：

- `VIBIT_RUNTIME_STORE=memory`，默认 quick smoke posture；
- `VIBIT_RUNTIME_STORE=postgres`，当前 alpha composition，用于 persistent inventory、player accounts、device credential login、token issuance、runtime sessions、request-level route protection、first-message binding、logout 和 presence query。

PostgreSQL posture 需要：

- `VIBIT_POSTGRES_DSN`；
- `VIBIT_AUTH_VERIFIER_KEY_SET_ID`；
- `VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY`；
- `VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY`；
- `VIBIT_AUTH_TOKEN_LOOKUP_KEY`；
- `VIBIT_AUTH_TOKEN_VERIFIER_KEY`；
- `docs/runtime-runbook.md` 中已经记录的可选 token lifetime 和 audience variables。

Package 可以改善围绕 `cd runtime && go run ./cmd/vibit-server`、local ports、status endpoints 和 redacted diagnostics 的文档和脚本。它不得静默启动 hosted services，不得在普通 server startup 中自动 apply migrations，也不得隐藏缺失配置失败。

## 5. Migration Expectations

第一版本地开发路径应该让 migrations 保持显式：

- SQL migration sources 继续放在 `runtime/migrations/postgres/`；
- 普通 `vibit-server` startup 不 apply migrations；
- local setup instructions 应告诉开发者 fresh database 什么时候需要 migration apply 或 status verification；
- 任何未来脚本都必须显式调用已有 repository-owned migration behavior 或文档化命令；
- migration output 不得打印 DSN 或 secrets。

这个 gate 不授权新的 migration files、repository/storage adapter changes、schema changes、automatic startup migration apply behavior 或新的 migration dependencies。

## 6. Configuration And Secret Posture

第一版 local package 可以添加 redacted example environment file 或 documented local environment template，但绝不能提交真实 local secrets。

本地路径必须把以下内容视为 not log-safe：

- raw device credential text 或 bytes；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- HMAC input 或 output bytes；
- verifier key values；
- concrete verifier key set ids；
- 带 credentials 的 PostgreSQL DSNs；
- 可能携带 secrets 的 headers、cookies、query strings、WebSocket subprotocol values 或 remote addresses；
- `.vibit.local.env` 这类 private local environment files。

提交进仓库的 examples 必须使用 placeholders。若之后文档化本地生成 secrets，它们必须由开发者在本地生成，并保留在 version control 之外。

## 7. Example Flow Shape

第一版 prototype-ready local package 应展示一个多步骤 flow，而不是单个孤立 request：

```text
prerequisite check
-> redacted local configuration preparation
-> explicit migration status or apply step
-> runtime startup guidance
-> local onboarding through the existing application-owned path
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

现有 `examples/local-alpha-request-loop.sh` 已经通过 focused Go E2E test 证明了大部分 flow。未来 package 可以 wrap、document 或扩展 local example ergonomics，但这个 gate 不授权新的 public onboarding route、新的 WebSocket routes、新的 Protobuf messages、SDK publication 或 production client surface。

## 8. Allowed Future Write Areas

`W-0200 Implement prototype-ready local development path package` 可以使用以下 write areas，前提是仍然留在本 gate 内：

- `README.md` 和 `README.zh-CN.md`；
- `docs/runtime-runbook.md`，以及如果之后存在 paired translation，则同步更新；
- `docs/alpha-developer-flow.md` 和 `docs/alpha-developer-flow.zh-CN.md`；
- `docs/alpha-acceptance-checklist.md` 和 `docs/alpha-acceptance-checklist.zh-CN.md`；
- 新的 `docs/` local development path package 文档及 paired Simplified Chinese translations；
- `examples/` scripts、README files 或 placeholder environment templates；
- 需要时添加 `.gitignore` local-only environment files 规则；
- `tools/vibit` 和 `rules/check-rules.json`，用于 static repository checks；
- `.arch/` manifests、change specs、ADRs、AGENTS guides 和 conversation memory；
- focused tests 或 script checks，只能验证既有行为，不能改变 production runtime behavior。

未来 package 在修改 production Go runtime behavior、protocol source、generated output、SQL migrations、repository interfaces、dependencies、release artifacts、hosted deployment surfaces、broad operations/admin behavior、authentication/session semantics 或 broad product modules 前必须先询问。

## 9. Verification Expectations

这个 gate 期望未来 package 保持默认 verification local 且 source-first：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```

Optional live PostgreSQL verification 可以继续通过 `VIBIT_POSTGRES_TEST_DSN` 和 disposable database 选择性启用。默认 repository checks 不得要求运行中的 database 或 private secrets。

## 10. Stop Conditions

如果执行 local path 需要以下任何事项，必须停止并请求 maintainer authorization：

- runtime behavior changes；
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
- GitHub release record 之外的 public announcements；
- paid promotion；
- handling 或 disclosure of secrets。

## 11. Next Work

下一项 bounded direction 是：

```text
W-0200 Implement prototype-ready local development path package
```

该 work 应按照以上 gate 打包本地开发路径，然后把 storage objects、realtime messaging 或 server push、failure/concurrency verification、operations inspection 等后续 product-capability work 留给独立 work items。

## 12. Verification

仓库应验证：

- 本 gate 和它的译文存在；
- `ADR-0107` 记录该决策；
- `.arch` manifests 标记 `W-0199` completed，并打开 `W-0200`；
- README、alpha goal、developer flow、acceptance checklist、AGENTS guides 和 product roadmap 指向新的 next work；
- runtime、protocol、generated output、migration、dependency、operations/admin、authentication/session、product module、hosted deployment、release artifact、public announcement、paid promotion 和 direct compatibility deferrals 保持不变。
