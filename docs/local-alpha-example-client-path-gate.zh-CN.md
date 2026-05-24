# Local Alpha Example Client Path Gate

状态：Accepted v0.1
最后更新：2026-05-24
范围：第一条更清晰的 source-first local alpha example client 或 example app path 的 gate-only 边界
依赖：`decisions/ADR-0132-select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof.md`、`docs/agent-native-feature-request-test-workflow.md`、`docs/prototype-ready-local-development-path-package.md`、`docs/alpha-developer-flow.md`、`examples/README.md`
Canonical decision: `ADR-0133`

配套英文源文件是 `docs/local-alpha-example-client-path-gate.md`。英文文件是权威版本。

本文定义 local alpha example client path gate。它是 gate artifact。本 slice 不实现 example client 或 example app，不发布 SDK，不生成 client libraries，不添加依赖，不添加 runtime behavior，不添加 protocol routes，不添加 Protobuf source，不改变 generated output，不添加 migrations，不添加 persistence，不改变 startup wiring，不改变 authentication/session behavior，不添加 delivery guarantees、stream subscriptions、chat rooms、groups、broadcast fanout、matchmaking、match runtime、operations/admin behavior、hosted deployments、release artifacts、public announcements、paid promotion、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

local alpha example client path gate record 是：

```yaml
local_alpha_example_client_path_gate: defined
completed_work_item: W-0225
decision: ADR-0133
check_rule: runtime.local_alpha_example_client_path_gate
source_selection_decision: ADR-0132
source_workflow_decision: ADR-0128
standard: docs/local-alpha-example-client-path-gate.md
translation: docs/local-alpha-example-client-path-gate.zh-CN.md
selected_nakama_capability_family: client_sdks_examples_and_developer_experience
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
path_posture: source_first_local_alpha_example_client_path
public_sdk_posture: not_a_public_sdk
live_network_client_posture: deferred_until_public_onboarding_and_client_package_boundary
future_example_docs_candidate: examples/local-alpha-client/README.md
future_example_script_candidate: examples/local-alpha-example-client.sh
future_runtime_proof_candidate: runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
future_implementation_work_item: W-0226
future_implementation_direction: implement_local_alpha_example_client_path
accepted_existing_routes:
  - runtime.authentication.AuthenticateWithDeviceCredential
  - runtime.connection.BindConnection
  - inventory.GrantItem
  - inventory.GetInventory
  - presence.GetPlayerPresence
  - storage.GetOwnStorageObject
  - storage.ListOwnStorageObjects
  - storage.PutOwnStorageObject
  - storage.DeleteOwnStorageObject
  - runtime.authentication.LogoutAccessToken
implementation_added_by_this_gate: false
example_client_implementation_added: false
sdk_added: false
generated_client_library_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
persistence_added: false
startup_wiring_added: false
authentication_session_behavior_changed: false
delivery_guarantees_added: false
stream_subscription_added: false
chat_added: false
group_messaging_added: false
broadcast_fanout_added: false
matchmaking_added: false
match_runtime_added: false
operations_admin_behavior_added: false
hosted_deployment_added: false
release_artifact_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 目的

`ADR-0132` 选择 `client_sdks_examples_and_developer_experience` 作为下一项 Nakama-first prototype-ready capability family。local alpha 已经证明 authentication、connection binding、protected inventory、presence、storage objects、logout、realtime outbound foundation 和 failure-path behavior，但可见路径仍然主要是一个 shell script 包装内部 Go E2E proof。

下一步 implementation 应让开发者和 AI agents 更容易把当前能力读成一个 client-like flow，然后再扩大产品能力。第一条 example path 必须诚实保留当前 alpha 约束：

- local onboarding 是 application service behavior，不是 public client route；
- generated Protobuf Go output 位于 `runtime/internal/`，还不是 public client package；
- 当前 proof 是 source-first、repo-local；
- secrets 和 transport metadata 必须继续 redacted；
- example 不是 SDK、hosted demo、package 或 compatibility promise。

## 3. 选定 Path Shape

第一条允许的路径是 source-first local alpha example path：

```text
examples/local-alpha-client/README.md
examples/local-alpha-example-client.sh
runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
```

后续 implementation 可以添加更清晰的 example-client README 和 wrapper script，指向 focused、命名清楚的 runtime proof。如果现有 E2E test 对这个目的仍然过密，后续 implementation 可以在现有 Protobuf E2E tests 旁边添加 focused local alpha example-flow test，但必须复用已有 runtime 和 protocol surfaces。

规则：

- example path 必须是 source-first，并且 local to repository。
- top-level example entrypoint 应在 `examples/` 下。
- 任何导入 `runtime/internal/...` 的 Go proof 必须位于 `runtime` tree 或现有 internal-test boundary 下，满足 Go internal package rules。
- path 必须展示当前 local alpha loop，而不是发明新的 product API。
- path 不能创建 public client SDK，也不能暗示稳定的 client package compatibility。

## 4. 展示 Flow

第一条后续 implementation 应展示现有 local alpha flow：

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> protected own-player storage object put/get/list/delete
-> logout
-> rejected post-logout protected request
-> selected failure-path redaction checks
```

implementation 可以把现有 realtime outbound foundation 作为当前能力上下文提到，但不能添加 stream subscriptions、chat rooms、broadcast fanout、delivery guarantees、offline inboxes、acknowledgements、retries 或新的 realtime behavior。

## 5. Ownership

后续 implementation ownership：

```yaml
example_docs_owner: examples/local-alpha-client
example_script_owner: examples
runtime_proof_owner: runtime/internal/platform/protocol/protobuf
application_behavior_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
transport_owner: runtime/internal/platform/transport/ws
```

规则：

- `examples/` 只能提供 human-facing docs 和 shell entrypoints。
- `examples/` 不能变成 product SDK 或 package publication source。
- Runtime behavior 继续由现有 runtime owners 负责。
- Protocol payload 和 envelope behavior 继续由现有 Protobuf adapter owners 负责。
- WebSocket transport 继续保持 credential-neutral 和 policy-neutral。
- domain module 不应 import example code。

## 6. Redaction

example path 不能打印、持久化、提交或记录：

- raw device credential text 或 bytes；
- raw access tokens；
- credential/token lookup digests；
- credential/token verifier digests；
- verifier key values；
- concrete verifier key set ids；
- HMAC inputs 或 outputs；
- 带 credentials 的 PostgreSQL DSNs；
- headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 concrete transport metadata。

允许输出 route names、step names、redacted status classes 和 high-level success/failure descriptions。

## 7. Nakama Mapping

Nakama reference mapping：

- 本 gate 覆盖 `client_sdks_examples_and_developer_experience` capability family。
- 它采纳 backend framework 需要清晰 client-facing evaluation path 的产品压力。
- 它不复制 Nakama public REST paths、client package names、runtime helper names、wire payloads、storage APIs、session token shapes 或 compatibility promises。

Pitaya reference mapping：

- Pitaya 继续 deferred as future distributed architecture reference。
- 本 gate 不引入 frontend/backend server roles、RPC、service discovery、groups、cluster routing 或 distributed session behavior。

## 8. 后续 Implementation Work

打开：

```text
M-154/W-0226 Implement local alpha example client path
```

后续 work item 可以：

- 添加 `examples/local-alpha-client/README.md`；
- 添加 `examples/local-alpha-example-client.sh`；
- 更新 `examples/README.md` 和 `examples/README.zh-CN.md`；
- 可选地在 `runtime/internal/platform/protocol/protobuf` 下添加或整理 focused local alpha example-flow test；
- 更新 repository checks 和 durable memory。

后续 work item 不能：

- 添加 runtime behavior；
- 添加新 protocol routes；
- 添加 Protobuf source 或 generated output；
- 添加 migrations、persistence、repository interfaces、adapters、dependencies、startup wiring、SDK publication、generated client libraries、hosted demos、release artifacts、direct compatibility 或 Pitaya-style distributed architecture。

## 9. Verification Expectations

后续 implementation 应验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.local_alpha_example_client_path_gate
node tools/vibit check change define-local-alpha-example-client-path-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

如果后续 implementation 修改或添加 Go proof，也必须运行 focused Go tests 和 `cd runtime && go test ./...`。

## 10. Stop Conditions

如果 example path 需要以下内容，应停止并创建单独 gate：

- public onboarding protocol route；
- public client packages 或 generated client libraries；
- stable SDK API guarantees；
- new Protobuf source 或 generated output；
- new runtime behavior；
- new authentication/session behavior；
- startup wiring 或 live server auto-setup；
- new dependencies；
- 默认验证依赖 live PostgreSQL；
- hosted deployments 或 release artifacts；
- chat、groups、stream subscriptions、broadcast fanout、matchmaking、match runtime、operations/admin behavior、distributed runtime 或 direct Nakama/Pitaya API compatibility。
