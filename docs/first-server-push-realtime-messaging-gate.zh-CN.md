# First Server Push And Realtime Messaging Gate

状态：Accepted v0.1
最后更新：2026-05-23
范围：storage objects local proof 之后第一版 server push 和 realtime messaging vocabulary 的 gate-only boundary
依赖：`decisions/ADR-0121-next-alpha-direction-first-server-push-realtime-messaging-gate.md`、`docs/game-protocol.md`、`docs/runtime-protocol-adapter.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/prototype-ready-foundation-execution-plan.md`
Canonical decision：`ADR-0122`

英文文件 `docs/first-server-push-realtime-messaging-gate.md` 是权威版本。本文件是配套简体中文翻译。

本文定义第一版 server push 和 realtime messaging gate。它是 gate artifact。它不添加 runtime behavior、transport delivery、protocol routes、Protobuf source、generated output、startup wiring、persistence、migrations、dependencies、authentication/session behavior changes、hosted deployments、release artifacts、public announcements、paid promotion、matchmaking、match runtime、distributed runtime、broad chat/social modules、large object/blob storage、S3-compatible object storage 或 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

第一版 server push 和 realtime messaging gate 记录为：

```yaml
first_server_push_realtime_messaging_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0214
decision: ADR-0122
check_rule: runtime.first_server_push_realtime_messaging_gate
source_next_direction_decision: ADR-0121
source_storage_objects_local_proof_decision: ADR-0120
future_runtime_owner_candidate: runtime/internal/app/realtime
future_runtime_service_source_candidate: runtime/internal/app/realtime/service.go
future_runtime_service_test_candidate: runtime/internal/app/realtime/service_test.go
future_protocol_source_candidate: proto/vibit/realtime/v1/realtime.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/realtime_bridge.go
future_application_handler_candidate: runtime/internal/app/bootstrap/realtime.go
future_transport_delivery_candidate: runtime/internal/platform/transport/ws/outbound.go
future_runtime_slice_work_item: W-0215
future_runtime_slice_direction: first_server_push_realtime_messaging_runtime_slice
first_delivery_model_candidate: single_process_bound_connection_delivery
first_message_intent_vocabulary_recorded: true
first_target_scope_candidates:
  - connection
  - player
  - stream
first_envelope_kind_candidates:
  - event
  - system
websocket_transport_credential_neutral: true
protocol_adapter_payload_mapping_only: true
application_policy_owner_required: true
backend_intent_owner_required: true
realtime_gate_only: true
runtime_behavior_added: false
transport_delivery_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
startup_wiring_added: false
persistence_added: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
matchmaking_added: false
match_runtime_added: false
distributed_runtime_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 目的

`W-0212` 已在 local alpha WebSocket/Protobuf request flow 中证明 own-player storage object routes。`W-0213` 选择第一版 server push 和 realtime messaging gate 作为下一项 prototype-ready 方向。

下一项有用边界不是 implementation，而是 outbound realtime behavior 的 vocabulary 和 ownership gate。没有这个 gate，后续 agent 可能把 server push 藏进 WebSocket transport，把 message policy 和 serializer 混在一起，或复制 Nakama/Pitaya 的 public route/API shape。

Nakama 提供 product pressure：durable storage 之后，常见 game backend 需要 notifications、streams、chat、presence-adjacent signals 和 realtime socket messages。

Pitaya 提供 architecture pressure：acceptors、sessions、handlers、push、groups、broadcast、backend services，以及后续 cluster/RPC topology 必须分离。

vibit 用 agent-native boundary 改造这些参考：

- transport 只移动 bytes 并负责 connection mechanics；
- protocol adapters 负责 encode、decode 和 payload mapping；
- application policy 决定谁可以接收 outbound messages；
- backend/domain services 拥有 message intent 和 invariants；
- persistence、delivery guarantees、retry behavior 和 distributed fanout 保留在独立 gate 后面。

## 3. 所有权

未来第一版 runtime behavior 应由 application 层拥有：

```yaml
future_runtime_owner_candidate: runtime/internal/app/realtime
future_runtime_service_source_candidate: runtime/internal/app/realtime/service.go
future_runtime_service_test_candidate: runtime/internal/app/realtime/service_test.go
connection_registry_owner: runtime/internal/app/connection
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
future_application_handler_candidate: runtime/internal/app/bootstrap/realtime.go
```

规则：

- WebSocket transport 必须保持 credential-neutral 和 payload-policy-neutral。
- Protocol adapters 可以映射 realtime payload bytes 和 envelope metadata，但不得决定 recipients 或 delivery authorization。
- Application-owned realtime behavior 必须从 validated identity、connection registry state、explicit subscriptions 或未来 module policy 中决定 recipient targets。
- Backend/domain services 可以发出 intent，不得直接写 socket。
- Domain modules 不得 import WebSocket transport packages、delivery 相关 generated Protobuf packages 或 Pitaya/Nakama SDKs。
- 未来实现不得使用 direct Nakama/Pitaya API names、route names 或 public compatibility shims。

## 4. 第一版词汇

本 gate 预留一组窄的 vibit-native outbound realtime vocabulary：

```yaml
message_intent_kinds:
  - server_notice
  - domain_event_push
  - stream_message
  - presence_signal
recipient_targets:
  - connection_id_and_epoch
  - player_current_connections
  - stream_subscribers
delivery_outcomes:
  - accepted
  - no_active_recipient
  - recipient_not_authorized
  - payload_invalid
  - delivery_unavailable
```

规则：

- `server_notice` 用于 server-authored operational 或 lifecycle notices。
- `domain_event_push` 用于未来 module 发出的 server facts，不用于 client-authored facts。
- `stream_message` 是 future stream-targeted delivery vocabulary。它不在本 gate 中实现 chat、channels、rooms、groups 或 subscriptions。
- `presence_signal` 保留给未来 presence-adjacent outbound facts。它不改变现有 presence lifecycle behavior。
- `connection_id_and_epoch` 只指向一个 server-observed connection。
- `player_current_connections` 指向一个 validated player 当前 active bound connections。
- `stream_subscribers` 只是未来词汇，直到 subscription ownership 被定义。
- Delivery outcome 不得泄露 raw access tokens、verifier material、credentials、DSNs、storage object values、SQL details、private account data 或隐藏 recipient existence。

## 5. 未来协议形状

第一版未来 protocol source candidate 是：

```text
proto/vibit/realtime/v1/realtime.proto
```

第一版未来 generated output candidate 是：

```text
runtime/internal/generated/proto/vibit/realtime/v1/realtime.pb.go
```

候选未来 message names：

```yaml
future_messages:
  ServerNotice:
    purpose: server-authored notice payload
  RealtimeEnvelope:
    purpose: stable public payload wrapper for outbound realtime messages
  StreamMessage:
    purpose: future stream-targeted message payload
```

规则：

- 本 gate 不添加 Protobuf source 或 generated output。
- 现有 `vibit.protocol.v1.Envelope` 在本 gate 中保持不变。
- 未来 server-to-client messages 应使用现有 envelope `kind` values，例如 `event` 或 `system`，除非后续 protocol ADR 明确改变 envelope semantics。
- 未来 target scopes 可以使用现有 envelope target vocabulary，但 connection-specific targeting 必须是 server-observed，不得信任 client-supplied connection ids 作为 authority。
- 未来 payloads 必须是 vibit-native，不得复制 Nakama notification/channel/stream payloads 或 Pitaya route payload conventions。

## 6. 未来运行时流程

第一版未来 runtime slice 应保持以下顺序：

```yaml
future_outbound_realtime_flow:
  - backend_or_application_service_creates_server_authored_message_intent
  - application_realtime_service_validates_intent_and_policy
  - application_realtime_service_resolves_allowed_recipients
  - protocol_adapter_maps_intent_to_existing_envelope_and_payload_bytes
  - transport_delivery_adapter_writes_binary_frames_to_server_observed_connections
  - delivery_result_reports_redacted_outcomes
```

规则：

- Server-authored message intent 必须由 trusted runtime code 创建，不能直接接受 client facts。
- 未来 client-originated chat、stream publish 或 room/match messages 需要单独 gate。
- 未来实现可以从 single-process 开始。
- Offline inboxes、persistence、acknowledgements、ordering guarantees、retries、backpressure policy、durable stream offsets 和 distributed fanout 继续 deferred。
- Route handlers 不得直接写 sockets，除非后续 implementation ADR 明确定义 transport delivery adapter boundary。

## 7. 身份与授权

第一版 posture 保守：

```yaml
sender_authority: server_only
client_published_facts_allowed: false
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
connection_id_client_authority_allowed: false
request_token_required_for_client_initiated_realtime_requests: true
```

规则：

- Server push 必须是 server-authored。
- Client-supplied `player_id`、`session_id`、stream id、room id、match id 或 connection id 不得授予 delivery authority。
- 未来 client request 如果请求 publish 或 subscribe，必须先通过现有 route protection 和 application policy。
- 现有 request-token protected route behavior 在本 gate 中保持不变。
- 本 gate 不改变 authentication、access-token validation、runtime session persistence、WebSocket handshake authentication、first-message binding 或 bound identity route policy。

## 8. Nakama/Pitaya 参考映射

Nakama reference mapping：

- 本 gate 让 vibit 朝常见 game-backend realtime surface 前进：notifications、streams、chat 和 presence-adjacent outbound messages。
- 它只采纳 capability family pressure，不采纳 Nakama route paths、REST APIs、runtime API names、payload names 或 compatibility promise。

Pitaya reference mapping：

- 本 gate 采纳 Pitaya 在 acceptors、sessions、handlers、push、groups、broadcast、backend services 和 cluster/RPC vocabulary 上的 separation pressure。
- 它把 group/broadcast/remote/RPC behavior 保留到 single-process delivery 和 application policy 明确之后。
- 它不复制 Pitaya handler names、route conventions、package APIs 或 cluster topology。

## 9. 验证期望

未来 implementation slices 应包含聚焦测试：

- recipient target validation；
- metadata-only identity refusal；
- single-process bound connection resolution；
- redacted delivery errors；
- 若添加 realtime payload，则验证 protocol adapter mapping；
- 若添加 socket writes，则验证 transport delivery behavior；
- no direct Nakama/Pitaya compatibility markers；
- 不泄露 storage object value、token、credential、digest、DSN 或 transport metadata。

本 gate 自身通过以下方式验证：

- document 和 translation 存在；
- ADR 和 change spec 存在；
- `.arch` manifest status markers；
- `tools/vibit` check coverage；
- 不添加 Go runtime behavior、Protobuf source、generated output、migration、dependency、startup wiring 或 direct compatibility。

## 10. 停止条件

添加以下任一内容前必须停止并请求 maintainer authorization：

- runtime service code；
- WebSocket outbound delivery code；
- Protobuf source 或 generated output；
- protocol routes 或 startup registration；
- chat、channels、groups、parties、rooms、matches、matchmaking 或 match runtime；
- stream subscription persistence；
- offline inboxes、acknowledgements、retries、ordering guarantees、durable offsets 或 backpressure mechanisms；
- distributed fanout、frontend/backend split、service discovery、RPC 或 cluster groups；
- credential、token、authentication、session、handshake 或 route-protection semantic changes；
- repository interfaces、PostgreSQL adapters、migrations、dependencies、hosted deployments、release artifacts、public announcements、paid promotion、blob/S3 storage 或 direct Nakama/Pitaya API compatibility。

## 11. 下一步

下一项 bounded work item 是：

```text
W-0215 Implement first server push and realtime messaging runtime slice
```

下一步只应实现本 gate 及其 implementation ADR 授权的最小 runtime slice。除非明确授权，broad chat/social behavior、protocol expansion、generated output、persistence、delivery guarantees、distributed fanout、matchmaking、match runtime 和 direct compatibility 仍必须保留到后续 bounded work items。
