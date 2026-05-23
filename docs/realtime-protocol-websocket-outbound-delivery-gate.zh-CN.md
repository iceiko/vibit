# Realtime Protocol And WebSocket Outbound Delivery Gate

状态：Accepted v0.1
最后更新：2026-05-23
范围：application-owned realtime runtime slice 之后 realtime protocol payloads 和 WebSocket outbound delivery 的 gate-only boundary
依赖：`decisions/ADR-0124-next-alpha-direction-realtime-protocol-websocket-outbound-delivery-gate.md`、`decisions/ADR-0123-first-server-push-realtime-messaging-runtime-slice.md`、`docs/first-server-push-realtime-messaging-gate.md`、`docs/game-protocol.md`、`docs/runtime-protocol-adapter.md`、`docs/generated-output.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision：`ADR-0125`

英文文件 `docs/realtime-protocol-websocket-outbound-delivery-gate.md` 是权威版本。本文件是配套简体中文翻译。

本文定义 realtime protocol 和 WebSocket outbound delivery gate。它是 gate artifact。它不添加 runtime behavior、WebSocket outbound delivery、concrete socket writes、Protobuf source、generated output、protocol bridge behavior、protocol routes、application bootstrap handlers、startup wiring、persistence、migrations、dependencies、authentication/session behavior changes、route-protection changes、hosted deployments、release artifacts、public announcements、paid promotion、stream subscription persistence、offline inboxes、acknowledgements、retries、ordering guarantees、durable offsets、backpressure、distributed fanout、matchmaking、match runtime、broad social modules、blob/S3 storage 或 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

Realtime protocol 和 WebSocket outbound delivery gate 记录为：

```yaml
realtime_protocol_websocket_outbound_delivery_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0217
decision: ADR-0125
check_rule: runtime.realtime_protocol_websocket_outbound_delivery_gate
source_next_direction_decision: ADR-0124
source_realtime_runtime_slice_decision: ADR-0123
source_realtime_gate_decision: ADR-0122
runtime_intent_owner: runtime/internal/app/realtime
connection_registry_owner: runtime/internal/app/connection
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
future_protocol_source_candidate: proto/vibit/realtime/v1/realtime.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/realtime_bridge.go
future_application_handler_candidate: runtime/internal/app/bootstrap/realtime.go
future_transport_delivery_candidate: runtime/internal/platform/transport/ws/outbound.go
future_implementation_slice_work_item: W-0218
future_implementation_slice_direction: realtime_protocol_websocket_outbound_delivery_implementation_slice
first_protocol_delivery_model_candidate: single_process_server_observed_connection_delivery
first_envelope_kind_candidates:
  - event
  - system
first_payload_family_candidates:
  - server_notice
  - domain_event_push
  - stream_message
  - presence_signal
application_policy_owner_required: true
protocol_adapter_payload_mapping_only: true
websocket_transport_payload_policy_neutral: true
server_observed_connection_id_required: true
server_observed_connection_epoch_required: true
client_connection_id_authority_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
direct_nakama_pitaya_api_compatibility_added: false
runtime_behavior_added: false
websocket_outbound_delivery_added: false
socket_write_added: false
protocol_bridge_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
startup_wiring_added: false
persistence_added: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
route_protection_changed: false
delivery_guarantees_added: false
distributed_runtime_added: false
matchmaking_added: false
match_runtime_added: false
broad_social_module_added: false
```

## 2. 目的

`W-0215` 已在 `runtime/internal/app/realtime` 下添加 application-owned realtime service。该 service 校验 server-authored outbound intents 并解析允许的 active recipients，但它有意只返回 delivery intents，不写 WebSocket frames。

`W-0216` 选择本 gate 作为下一项 bounded direction。当前缺失的边界，是从 application-owned delivery intent 到未来 client-visible protocol payloads 与 WebSocket outbound frame delivery 的计划性交接。

没有这个 gate，后续工作可能混合几类职责：

- message authorization 和 recipient resolution；
- Protobuf payload selection 和 generated-output ownership；
- protocol envelope mapping；
- WebSocket connection mechanics 和 socket writes；
- fanout、delivery guarantees、offline storage 和 distributed runtime concerns。

本 gate 在任何 wire 或 socket behavior 实现前把这些职责分开。

## 3. 所有权

未来 outbound path 必须保持以下 owner 分离：

```yaml
runtime_intent_owner: runtime/internal/app/realtime
connection_registry_owner: runtime/internal/app/connection
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
process_wiring_owner: runtime/cmd/vibit-server
application_bootstrap_owner: runtime/internal/app/bootstrap
```

规则：

- `runtime/internal/app/realtime` 拥有 outbound intent validation、recipient target validation 和 policy-facing delivery outcomes。
- `runtime/internal/app/connection` 拥有 server-observed connection id 和 epoch state。
- `runtime/internal/platform/protocol/protobuf` 可以把已授权的 delivery intents 映射成 Protobuf payload bytes 和现有 envelope metadata。
- `runtime/internal/platform/transport/ws` 可以把已编码的 binary frames 写到 server-observed connections。
- WebSocket transport 不得决定 recipient authorization、解析 domain payloads 或构造 domain-specific payloads。
- Protocol adapters 不得决定谁可以接收消息。
- Domain modules 不得 import WebSocket transport package，也不得直接写 socket frames。
- 未来 startup wiring 必须显式，不得把 business behavior 隐藏在 process assembly 中。

## 4. 未来协议形状

第一版未来 protocol source candidate 仍是：

```text
proto/vibit/realtime/v1/realtime.proto
```

第一版未来 generated output candidate 仍是：

```text
runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
```

第一版未来 protocol bridge candidate 是：

```text
runtime/internal/platform/protocol/protobuf/realtime_bridge.go
```

候选未来 payload family：

```yaml
future_payloads:
  ServerNotice:
    intent_kind: server_notice
    purpose: server-authored lifecycle or operational notice visible to an authorized client
  DomainEventPush:
    intent_kind: domain_event_push
    purpose: server-authored module fact prepared for outbound delivery
  StreamMessage:
    intent_kind: stream_message
    purpose: future stream-targeted payload after subscription ownership is defined
  PresenceSignal:
    intent_kind: presence_signal
    purpose: future presence-adjacent outbound signal without changing current presence lifecycle semantics
```

规则：

- 本 gate 不添加 Protobuf source 或 generated output。
- 现有 `vibit.protocol.v1.Envelope` 保持不变。
- 未来 realtime server-to-client payloads 应使用现有 envelope `kind` values，例如 `event` 或 `system`，除非后续 ADR 改变 envelope semantics。
- 未来 payloads 必须是 vibit-native，不得复制 Nakama notification、channel、stream、chat 或 Pitaya route payload conventions。
- Generated Go Protobuf output 仍必须遵循 `docs/generated-output.md`；普通 agent 不得手工编辑 generated files。
- Connection-specific targeting 必须依赖 server-observed connection id 和 epoch，而不是 client-supplied authority。

## 5. 未来 Outbound Flow

第一版未来 implementation slice 应保持以下计划流程：

```yaml
future_outbound_delivery_flow:
  - trusted_runtime_code_creates_server_authored_realtime_intent
  - application_realtime_service_validates_intent_and_recipient_policy
  - application_realtime_service_resolves_delivery_intents
  - protocol_adapter_maps_delivery_intents_to_envelope_and_payload_bytes
  - transport_outbound_adapter_writes_binary_frames_to_server_observed_connections
  - delivery_result_reports_redacted_outcomes
```

规则：

- Application realtime service 仍是 policy owner。
- Protocol mapping 只在 application service 返回 accepted delivery intent 后开始。
- Transport delivery 只在 protocol adapter 返回 encoded binary frame bytes 后开始。
- Delivery results 必须 redacted，不得泄露 raw credentials、raw tokens、verifier digests、DSNs、SQL details、private account data 或授权 public class 之外的隐藏 recipient existence。
- Write 失败不得改变 domain state，除非后续 delivery guarantee gate 明确定义 durable delivery semantics。

## 6. WebSocket Transport Delivery Boundary

未来 transport delivery candidate 是：

```text
runtime/internal/platform/transport/ws/outbound.go
```

未来 transport delivery 可以：

- 接受 encoded binary frame；
- 指向 server-observed connection id 和 epoch；
- 通过现有 WebSocket connection owner 写出；
- 向调用方报告 redacted write outcomes。

未来 transport delivery 不得：

- 选择 recipients；
- 检查 domain payload semantics；
- 校验 authentication credentials；
- 创建 Protobuf payloads；
- 实现 stream subscriptions、chat、groups、broadcast fanout、offline inboxes、delivery guarantees、retries、ordering、backpressure、durable offsets、cluster routing、RPC、service discovery 或 direct compatibility shims。

## 7. 身份与授权

第一版 posture 仍是 server-authored 且保守：

```yaml
sender_authority: server_or_admin_validated_identity_only
client_published_facts_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
connection_id_client_authority_allowed: false
route_protection_changed_by_this_gate: false
websocket_handshake_authentication_changed_by_this_gate: false
```

规则：

- Client-supplied player id、session id、stream id、room id、match id 或 connection id 不得授予 outbound delivery authority。
- 未来 client-originated publish、subscribe、chat、stream、room 或 match messages 需要单独 gate。
- 现有 route-protection behavior、first-message connection binding behavior、runtime session behavior、logout behavior 和 access-token validation behavior 保持不变。
- Metadata-only identity 仍然不是充分 proof。

## 8. Nakama/Pitaya 参考映射

Nakama reference mapping：

- 本 gate 让 vibit 向 client-visible outbound realtime capability 前进，这些能力支撑 notifications、streams、chat 和 presence-adjacent features。
- 它只采纳 capability pressure，不复制 Nakama public APIs、route paths、runtime helper names、payload names 或 compatibility promises。

Pitaya reference mapping：

- 本 gate 保留 Pitaya-style separation：acceptor/transport mechanics、sessions and connection state、handlers、protocol serialization、backend service intent、push/group/broadcast vocabulary 和后续 cluster concerns 彼此分离。
- Group membership、broadcast fanout、remote calls、RPC、service discovery 和 frontend/backend role separation 继续 deferred。
- 它不复制 Pitaya route conventions、package APIs、handler naming 或 cluster topology。

## 9. 验证预期

未来 implementation slices 应包含聚焦测试：

- 从 accepted realtime delivery intents 到 protocol payload 的 mapping；
- 拒绝错误 intent kinds、target kinds 和 malformed payloads；
- binary frame write handoff 到 server-observed connection id 和 epoch；
- stale epoch 或 missing connection handling；
- transport write error redaction；
- application policy 不进入 protocol 和 transport adapters；
- domain modules 不 import WebSocket transport packages；
- 不手工编辑 generated output；
- 不添加 direct Nakama/Pitaya compatibility markers。

本 gate 自身通过以下内容验证：

- 英文和简体中文标准存在；
- ADR 和 change spec 存在；
- `.arch` manifest status markers；
- `runtime.realtime_protocol_websocket_outbound_delivery_gate` check coverage；
- 不添加 Go runtime behavior、Protobuf source、generated output、protocol bridge、WebSocket outbound writer、startup wiring、migration、dependency、delivery guarantee、distributed runtime 或 direct compatibility。

## 10. 停止条件

在后续明确授权的 implementation work item 之外，添加以下任一内容前必须停止并询问：

- `proto/vibit/realtime/v1/realtime.proto`；
- `runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go`；
- `runtime/internal/platform/protocol/protobuf/realtime_bridge.go`；
- `runtime/internal/platform/transport/ws/outbound.go`；
- `runtime/internal/app/bootstrap/realtime.go`；
- outbound delivery 的 startup wiring；
- concrete socket writes；
- protocol routes 或 client publish routes；
- stream subscription persistence；
- chat、channels、groups、parties、rooms、matches、matchmaking 或 match runtime；
- offline inboxes、acknowledgements、retries、ordering guarantees、durable offsets 或 backpressure mechanisms；
- distributed fanout、frontend/backend split、service discovery、RPC 或 cluster groups；
- credential、token、authentication、session、handshake 或 route-protection semantic changes；
- repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployments、release artifacts、public announcements、paid promotion、blob/S3 storage 或 direct Nakama/Pitaya API compatibility。

## 11. 下一步工作

下一项 bounded work item 是：

```text
W-0218 Implement realtime protocol and WebSocket outbound delivery slice
```

下一步应只实现本 gate 和配套 implementation ADR 授权的最小 slice。Stream subscription ownership、chat semantics、groups、broadcast fanout、offline inboxes、acknowledgements、retries、ordering guarantees、durable offsets、backpressure、distributed fanout、matchmaking、match runtime 和 direct compatibility 继续保留在后续 bounded work items 后面，除非被明确授权。
