# Pitaya-Aligned Serializer And Message Forwarding Boundary Gate（中文版）

状态：Accepted v0.1
最后更新：2026-06-01
范围：route handler pipeline source-first map 之后，使用 Pitaya-aligned serializer 与 message forwarding 词汇的 gate-only boundary
依赖：`decisions/ADR-0169-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map.md`、`decisions/ADR-0168-pitaya-aligned-route-handler-pipeline-source-first-map.md`、`docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/game-protocol.md`、`docs/reference-game-server-alignment.md`、`.arch/reference.yaml`
规范决策：`ADR-0170`

本文件是 `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md` 的简体中文译本。英文文件是权威版本。

本文只定义 serializer 与 message forwarding vocabulary gate。它不实现 serializer behavior、message forwarding behavior、route handler implementation、handler routing behavior、handler pipeline behavior、pipeline middleware behavior、backend route targeting、cluster-safe session routing behavior、distributed session routing、distributed runtime behavior、distributed groups、room broadcast fanout、delivery guarantees、stream subscriptions、service discovery implementation、service registries、service selectors、node identity、server-to-server RPC、remote calls、frontend/backend server roles、runtime endpoint behavior、metrics endpoints、observability pipelines、dashboards、protocol messages or routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、SDK publication、hosted deployments、release artifacts 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Pitaya-aligned serializer and message forwarding boundary gate 记录为：

```yaml
pitaya_aligned_serializer_message_forwarding_boundary_gate: defined
completed_work_item: W-0262
decision: ADR-0170
check_rule: runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate
previous_direction_decision: ADR-0169
route_handler_pipeline_source_first_map_decision: ADR-0168
route_handler_pipeline_source_first_map_check_rule: runtime.pitaya_aligned_route_handler_pipeline_source_first_map
standard: docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md
translation: docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: serializer_message_forwarding_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_serializer_message_forwarding_vocabulary
future_implementation_work_item: W-0263
future_implementation_direction: implement_pitaya_aligned_serializer_message_forwarding_source_first_map
allowed_serializer_message_forwarding_vocabulary:
  - serializer_boundary
  - serializer_format
  - encode_boundary
  - decode_boundary
  - message_forwarding
  - forwarding_target
  - forwarding_envelope
  - delivery_handoff
```

所有 implementation flags 在本 gate 中保持 `false`，包括 serializer behavior、message forwarding behavior、route handler implementation、handler routing behavior、handler pipeline behavior、pipeline middleware behavior、backend route targeting、service discovery、RPC、distributed runtime、protocol、generated output、persistence、dependencies、hosted surfaces、SDKs 和 direct compatibility。

## 2. Purpose

`ADR-0169` 在 route handler pipeline source-first map 之后选择 serializer 与 message forwarding vocabulary 作为下一项 Pitaya-aligned direction。

风险在于，agent 可能把 serializer 和 forwarding 词汇误解为可以添加 pluggable codecs、protocol shape changes、forwarding workers、backend route targeting、service discovery、RPC、remote delivery 或 distributed runtime behavior。本文只记录词汇和映射；vibit 现有 single-process WebSocket Protobuf flow 保持不变。

## 3. Vocabulary

允许的 serializer 与 message forwarding vocabulary：

- `serializer_boundary`：未来 encode/decode ownership 的规划词汇；当前 concrete boundary 仍是 Protobuf bridge functions。
- `serializer_format`：未来 wire format selection 的规划词汇；当前 envelope format 仍是既有 Protobuf envelope。
- `encode_boundary`：未来 outbound payload encoding 的规划词汇；当前仍使用 explicit generated bridge functions。
- `decode_boundary`：未来 inbound payload decoding 的规划词汇；当前仍使用 explicit generated bridge functions。
- `message_forwarding`：未来把 message forward 到另一个 owner 或 node 的规划词汇；当前 runtime 没有 cross-node forwarding。
- `forwarding_target`：未来 forwarding owner selection 的规划词汇；当前 target scope metadata 不是 backend route targeting。
- `forwarding_envelope`：未来 internal forwarding wrapper 的规划词汇；本 slice 不存在 internal forwarding envelope。
- `delivery_handoff`：未来把 delivery 交给另一个 runtime owner 的规划词汇；当前仍是 single-process WebSocket delivery。

禁止用法：

- 不要引入 Pitaya 或 Nakama 的 public API、package、route、method、wire、handler、serializer、forwarding、registry、selector 或 configuration compatibility names。
- 不要把 serializer 或 forwarding vocabulary 当作添加 codecs、serializer registries、forwarding workers、backend route targeting、service discovery、RPC、remote calls、protocol messages、generated output、persistence、dependencies、topology 或 distributed runtime behavior 的授权。
- 不要把 domain behavior 移进 transport、Protobuf adapters、serializer boundaries、forwarding layers 或 process startup。

## 4. Current Mapping

```yaml
current_single_process_serializer_forwarding_mapping:
  protocol_bridge:
    current: explicit generated Protobuf payload bridge functions
    future_vocabulary: serializer_boundary
    status: no_pluggable_serializer_behavior
  envelope_encoding:
    current: Protobuf envelope owned by the protocol adapter
    future_vocabulary: serializer_format
    status: no_serializer_registry
  payload_encoding:
    current: generated payload bridge functions
    future_vocabulary: encode_boundary
    status: no_custom_encode_pipeline
  payload_decoding:
    current: generated payload bridge functions
    future_vocabulary: decode_boundary
    status: no_custom_decode_pipeline
  outbound_message:
    current: server-push intent converted to protocol envelope
    future_vocabulary: message_forwarding
    status: no_cross_node_forwarding
  target_scope:
    current: metadata-only target scope
    future_vocabulary: forwarding_target
    status: no_backend_route_targeting
  forwarding_envelope:
    current: no internal forwarding envelope
    future_vocabulary: forwarding_envelope
    status: not_implemented
  delivery_handoff:
    current: single-process WebSocket delivery
    future_vocabulary: delivery_handoff
    status: no_remote_delivery_handoff
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: unchanged
transport_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
module_owner: unchanged
dependency_owner: unchanged
```

文档和 manifests 可以定义 serializer/message forwarding vocabulary 与 current mapping。`tools/vibit` 可在后续 implementation work item 授权后输出 source-first map。Runtime、transport、protocol、repository、persistence、generated output、startup wiring、dependencies、service discovery、RPC、remote calls、frontend/backend role behavior、cluster-safe session routing、distributed group behavior 和 room broadcast behavior 在本 gate 中保持不变。

## 6. Nakama And Pitaya Mapping

Nakama 仍是近期 capability breadth 的 primary product reference。Pitaya 仍是 acceptors、sessions、route handlers、frontend/backend roles、RPC/remotes、service discovery、groups、broadcast、cluster routing、handler pipelines、serializers 和 forwarding 的 architecture vocabulary reference。

采纳为词汇：

- serializer boundary 与 serializer format；
- encode boundary 与 decode boundary；
- message forwarding、forwarding target、forwarding envelope 和 delivery handoff。

适配到 vibit：

- 当前 serialization 仍由 Protobuf adapter 和 explicit bridge functions 拥有。
- 当前 server push 仍是 single-process，不意味着 cross-node forwarding。
- 当前 target scope metadata 不意味着 backend route targeting。
- 任何后续 serializer 或 forwarding implementation 都必须另行 gate 和 verification。

当前拒绝：

- direct Pitaya 或 Nakama API compatibility；
- Pitaya 或 Nakama package、method、route、handler、serializer、forwarding、registry、selector 或 configuration naming compatibility；
- serializer behavior、message forwarding behavior、forwarding workers、backend route targeting、service discovery、RPC、remote calls、protocol changes、generated output、persistence、migrations、dependencies、hosted deployment、SDK publication 或 release artifacts。

## 7. Future Implementation Work

下一项：

```text
M-191/W-0263 Implement Pitaya-aligned serializer and message forwarding source-first map
```

该 work item 可以添加 serializer/message forwarding vocabulary 的 source-first repository inspection map，并总结 current protocol bridge、generated payload bridge、outbound message、target-scope metadata 和 delivery handoff mapping。

该 work item 不得添加 serializer behavior、message forwarding behavior、route handler implementation、handler routing behavior、handler pipeline behavior、backend route targeting、service discovery implementation、RPC、remote calls、distributed runtime behavior、protocol messages or routes、Protobuf source、generated output、persistence、migrations、dependencies、hosted deployment、SDK publication 或 direct Nakama/Pitaya API compatibility。

## 8. Verification

Required checks:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate
node tools/vibit check change define-pitaya-aligned-serializer-message-forwarding-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```
