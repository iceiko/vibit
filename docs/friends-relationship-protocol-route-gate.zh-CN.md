# Friends Relationship Protocol Route Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-26
范围：在 application runtime behavior 之后，为未来 client-facing friends relationship protocol routes 定义 gate-only boundary
依赖：`docs/friends-relationship-runtime-behavior-gate.md`、`decisions/ADR-0147-friends-relationship-runtime-behavior-implementation.md`、`docs/runtime-protocol-adapter.md`、`docs/game-protocol.md`、`docs/generated-output.md`、`docs/bound-identity-route-policy-gate.md`、`docs/runtime-session-validation-gate.md`、`docs/reference-game-server-alignment.md`、`docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0148`

英文版 `docs/friends-relationship-protocol-route-gate.md` 是权威版本。本文是配套简体中文翻译。

本文定义 friends relationship protocol route gate。它是一个 gate artifact。本 slice 不添加 protocol route implementation、Protobuf source、generated output、startup wiring、runtime handlers、repository interface changes、PostgreSQL adapter changes、migration changes、dependencies、authentication/session behavior changes、delivery guarantees、stream subscriptions、chat rooms、groups、parties、broadcast fanout、matchmaking、match runtime、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、public announcements、paid promotion、event/audit tables、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Friends relationship protocol route gate 记录是：

```yaml
friends_relationship_protocol_route_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0240
decision: ADR-0148
check_rule: runtime.friends_relationship_protocol_route_gate
source_runtime_behavior_implementation_decision: ADR-0147
source_runtime_behavior_implementation: runtime/internal/app/friends/service.go
source_runtime_behavior_tests: runtime/internal/app/friends/service_test.go
source_runtime_behavior_gate_decision: ADR-0146
source_repository_interface_decision: ADR-0143
repository_interface: runtime/internal/modules/friends.Repository
future_protocol_source_candidate: proto/vibit/friends/v1/friends.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/friends_bridge.go
future_protocol_bridge_test_candidate: runtime/internal/platform/protocol/protobuf/friends_bridge_test.go
future_application_handler_candidate: runtime/internal/app/bootstrap/friends.go
future_application_handler_test_candidate: runtime/internal/app/bootstrap/friends_test.go
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_actor_id_allowed_as_proof: false
client_supplied_target_player_id_allowed: true
first_actor_kind: player
first_payload_package: vibit.friends.v1
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
protocol_route_gate_only: true
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
runtime_handler_added: false
startup_wiring_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
delivery_guarantees_added: false
stream_subscription_added: false
chat_added: false
groups_added: false
parties_added: false
matchmaking_added: false
match_runtime_added: false
event_audit_table_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_protocol_route_implementation_work_item: W-0241
future_protocol_route_implementation_direction: implement_friends_relationship_protocol_route
```

## 2. Purpose

`W-0239` 已在 `runtime/internal/app/friends` 下实现 application-owned friends relationship behavior。下一步有价值的边界不是直接写 route code 或 `.proto`，而是先记录 future WebSocket/Protobuf exposure 如何调用该 service，同时不把 friends behavior 移到 transport、generated files 或 persistence adapters 中。

Nakama 提供 product surface 参考：friends、friend requests、blocks 和 actor-relative relationship status 是常见 game backend social graph capability。vibit 应覆盖这一能力类别。

Pitaya 提供 architecture posture 参考：acceptors、sessions、route handlers、serializers 和 backend services 应保持分离。vibit 通过 credential-neutral WebSocket transport、显式 Protobuf payload bridge、以及 application-owned route handlers 调用 application-owned friends service 来适配这一点。

本 gate 在 implementation 前记录：

- candidate route names；
- candidate request/response message shapes；
- route protection 和 identity handoff posture；
- protocol adapter、application handler 和 startup ownership；
- generated-output expectations；
- public error mapping 和 redaction expectations；
- Nakama/Pitaya reference mapping；
- stop conditions，确保本 slice 不加入 implementation 或 generated artifacts。

## 3. Future Route Surface

第一组 route family 应只为 validated actor 暴露 player-to-player friends relationship operations：

```yaml
candidate_routes:
  - kind: command
    module: friends
    name: SendFriendRequest
    route_id: friends.SendFriendRequest
    service_method: SendFriendRequest
  - kind: command
    module: friends
    name: AcceptFriendRequest
    route_id: friends.AcceptFriendRequest
    service_method: AcceptFriendRequest
  - kind: command
    module: friends
    name: RejectFriendRequest
    route_id: friends.RejectFriendRequest
    service_method: RejectFriendRequest
  - kind: command
    module: friends
    name: RemoveFriend
    route_id: friends.RemoveFriend
    service_method: RemoveFriend
  - kind: command
    module: friends
    name: BlockPlayer
    route_id: friends.BlockPlayer
    service_method: BlockPlayer
  - kind: command
    module: friends
    name: UnblockPlayer
    route_id: friends.UnblockPlayer
    service_method: UnblockPlayer
  - kind: query
    module: friends
    name: ListFriendRelationships
    route_id: friends.ListFriendRelationships
    service_method: ListFriendRelationships
  - kind: query
    module: friends
    name: GetFriendRelationshipStatus
    route_id: friends.GetFriendRelationshipStatus
    service_method: GetFriendRelationshipStatus
```

规则：

- Route names 必须保持 vibit-native，不复制 Nakama route paths 或 Pitaya route naming。
- Send、accept、reject、remove、block、unblock routes 是 commands。
- List 和 status routes 是 queries。
- 第一组 route family 只面向 validated player actor。不得暴露 arbitrary actor ids 或其他玩家的 private social graph。
- Client payload 可在 service 已有 vocabulary 范围内提供 target player id、expected relationship version、status filter、list limit 和 pagination cursor。
- Groups、parties、chat rooms、broadcast fanout、presence subscriptions、matchmaking filters、match social context、admin social graph inspection、account merge behavior、event/audit streams 和 script hooks 继续 deferred。
- Future route implementation 必须显式注册 routes。不允许 catch-all friends route 或 reflective handler。

## 4. Protocol Shape

第一份 friends relationship protocol source candidate 是：

```text
proto/vibit/friends/v1/friends.proto
```

第一份 generated output candidate 是：

```text
runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go
```

第一组 Protobuf package candidate 是：

```text
vibit.friends.v1
```

Candidate messages：

```yaml
messages:
  FriendRelationship:
    fields:
      relationship_id: string
      player_low_id: string
      player_high_id: string
      requested_by_player_id: string
      lifecycle_state: string
      public_status: string
      version: int64
      created_at: string
      updated_at: string
  FriendRelationshipPage:
    fields:
      relationships: repeated FriendRelationship
      next_pair_token: string
  SendFriendRequestRequest:
    fields:
      target_player_id: string
  SendFriendRequestResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  AcceptFriendRequestRequest:
    fields:
      target_player_id: string
      expected_version: int64
  AcceptFriendRequestResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  RejectFriendRequestRequest:
    fields:
      target_player_id: string
      expected_version: int64
  RejectFriendRequestResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  RemoveFriendRequest:
    fields:
      target_player_id: string
      expected_version: int64
  RemoveFriendResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  BlockPlayerRequest:
    fields:
      target_player_id: string
      expected_version: int64
  BlockPlayerResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  UnblockPlayerRequest:
    fields:
      target_player_id: string
      expected_version: int64
  UnblockPlayerResponse:
    fields:
      relationship: FriendRelationship
      status: string
      version: int64
  ListFriendRelationshipsRequest:
    fields:
      status: string
      limit: int32
      after_pair_token: string
  ListFriendRelationshipsResponse:
    fields:
      page: FriendRelationshipPage
  GetFriendRelationshipStatusRequest:
    fields:
      target_player_id: string
  GetFriendRelationshipStatusResponse:
    fields:
      public_status: string
      relationship: FriendRelationship
      version: int64
```

规则：

- 除非后续 protocol ADR 明确改变 envelope semantics，否则现有 `proto/vibit/protocol/v1/envelope.proto` 必须保持不变。
- 暴露时间值时应使用 RFC3339Nano UTC text。
- Optional `expected_version` mapping 必须保留 service 的 optional expected-version vocabulary。未来 implementation 必须在 tests 中明确 absence 与 `0` semantics。
- Public status values 必须映射 service 的 actor-relative statuses：`none`、`outgoing_request_pending`、`incoming_request_pending`、`friends`、`blocked_by_actor`、`blocked_actor`、`mutual_blocked`、`removed`、`rejected`。
- Response `status` values 必须映射 service operation outcomes，例如 `sent`、`accepted`、`request_rejected`、`removed`、`blocked`、`unblocked`、`listed`、`found` 或 `rejected`。
- Protocol shape 不得包含 client-supplied actor id、raw access tokens、credential material、lookup digests、verifier digests、SQL details、private transport metadata、chat payloads、group ids、party ids、matchmaking fields、match runtime fields 或 direct external API compatibility markers。
- Relationship ids、player ids、relationship versions、block state 和 lifecycle state 默认不是 log-safe。

## 5. Route Protection And Identity Handoff

第一组 route-policy posture 是：

```yaml
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_actor_id_allowed_as_proof: false
client_supplied_target_player_id_allowed: true
```

规则：

- Future friends routes 必须是 protected gameplay routes。
- Future handlers 必须从现有 protected-route flow 接收 validated `app.RequestIdentity`。
- Envelope/session metadata 中 metadata-only `player_id` 或 `session_id` 永远不得成为 friends actor proof。
- Client payload 不得选择 actor ids。
- Client payload 可为 player-to-player operations 选择 target player id。
- Service 仍负责在 id generation 或 repository access 前拒绝 invalid identity。
- 本 gate 不改变 authentication、token validation、session persistence、first-message binding、WebSocket handshake behavior、bound-identity policy 或 route-protection semantics。

## 6. Future Route Flow

Future implementation 必须保留以下顺序：

```yaml
future_route_flow:
  - receive WebSocket/Protobuf envelope through existing request path
  - apply protected-route authenticated wrapper policy
  - obtain validated authenticated request identity
  - decode vibit.friends.v1 request payload
  - reject payload actor proof and derive actor only from request identity
  - map payload fields to runtime/internal/app/friends service request
  - call application-owned friends service
  - map service result to vibit.friends.v1 response payload
  - map service public errors to protocol error responses
  - keep transport, Protobuf bridge, application handler, service, repository, and PostgreSQL adapter ownership separated
```

规则：

- WebSocket transport 保持 credential-neutral。
- Protobuf bridge 只应映射 payloads 和 response shapes；它不得拥有 friends relationship behavior。
- Application bootstrap handlers 应拥有 route registration 和 service invocation。
- Application service 继续拥有 identity checks、validation handoff、repository conflict mapping 和 actor-relative public status。
- PostgreSQL adapters 继续只负责 persistence。

## 7. Public Error Mapping

Future route implementation 应映射 service public errors，且不泄漏内部细节：

```yaml
public_error_mapping:
  FRIENDSHIP_INVALID_REQUEST: invalid_request
  FRIENDSHIP_UNAUTHENTICATED: unauthenticated
  FRIENDSHIP_FORBIDDEN: forbidden
  FRIENDSHIP_TARGET_NOT_FOUND: target_not_found
  FRIENDSHIP_RELATIONSHIP_NOT_FOUND: relationship_not_found
  FRIENDSHIP_DUPLICATE_REQUEST: duplicate_request
  FRIENDSHIP_ALREADY_FRIENDS: already_friends
  FRIENDSHIP_BLOCKED_RELATIONSHIP: blocked_relationship
  FRIENDSHIP_INVALID_TRANSITION: invalid_transition
  FRIENDSHIP_VERSION_MISMATCH: version_mismatch
  FRIENDSHIP_UNAVAILABLE: unavailable
```

规则：

- Public protocol errors 只能暴露 stable public codes/classes；retryability posture 只有后续 ADR 授权时才能暴露。
- Internal repository errors、SQL details、超出 service public error 的 target existence probes、relationship ids、player ids、access-token material、credential material、lookup digests、verifier digests 和 transport metadata 必须保持在默认 logs/errors 之外。
- Authentication 和 route protection failures 必须使用现有 protected-route semantics。本 gate 不发明新的 proof carrier。

## 8. Generated Output Posture

Future generated output 必须遵循 `docs/generated-output.md`。

规则：

- `proto/vibit/friends/v1/friends.proto` 只能由后续 implementation slice 添加。
- `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go` 只能作为 Buf/protoc 生成的 output 添加。
- Generated Go output 必须包含 `protoc-gen-go` generated-code marker，并可追溯到 source `.proto`。
- Agents 不得手工编辑 generated Go Protobuf files。
- 本 gate 不改变 `buf.yaml`、`buf.gen.yaml` 或 generated output。

## 9. Nakama Reference Mapping

Nakama reference mapping：

```yaml
nakama_reference_mapping:
  capability_family: friends_groups_and_parties
  mapped_capabilities:
    - friend_request_send
    - friend_request_accept
    - friend_request_reject
    - friend_remove
    - player_block
    - player_unblock
    - list_relationships
    - get_relationship_status
  direct_api_compatibility: false
```

Nakama 只提供 capability class 参考。vibit 不复制 Nakama route paths、field names、permission semantics、runtime script APIs 或 public API compatibility。

## 10. Pitaya Reference Mapping

Pitaya reference mapping：

```yaml
pitaya_reference_mapping:
  architecture_pressure:
    - acceptor_session_handler_separation
    - serializer_adapter_separation
    - backend_service_boundary
  distributed_architecture_status: deferred
  direct_api_compatibility: false
```

Pitaya 只提供 layering pressure。本 gate 不添加 Pitaya-style distributed topology、frontend/backend split、RPC、groups、service discovery、distributed social graph routing 或 direct Pitaya API compatibility。

## 11. Required Future Tests

Future implementation tests 应覆盖：

- 八个 route ids 的 route registration；
- command/query kind mapping；
- Protobuf request/response bridge mapping；
- optional expected-version mapping；
- 从 validated request identity 派生 actor；
- 通过现有 protected-route wrapper 拒绝 metadata-only identity；
- 拒绝 client-supplied actor proof；
- target player id validation 和 self-targeting behavior；
- none、pending incoming/outgoing、friends、blocked、removed、rejected states 的 public status mapping；
- service public error 到 protocol error 的映射；
- private relationship、player、token、credential、SQL 和 transport details 的 redaction；
- WebSocket transport 或 PostgreSQL adapter packages 中不得拥有 route behavior；
- 如果添加 Protobuf source，必须验证 generated-output traceability。

## 12. Stop Conditions

在添加以下内容前必须停止并创建单独 work item：

- protocol route implementation；
- `proto/vibit/friends/v1/friends.proto`；
- generated Go Protobuf output；
- protocol bridge implementation；
- application bootstrap handlers；
- startup route registration；
- new dependencies；
- migration changes；
- repository interface changes；
- PostgreSQL adapter changes；
- authentication/session behavior changes；
- event/audit tables；
- delivery guarantees；
- stream subscriptions；
- chat rooms；
- groups；
- parties；
- broadcast fanout；
- matchmaking；
- match runtime；
- operations/admin behavior；
- SDK publication；
- generated client libraries；
- hosted deployments；
- release artifacts；
- public announcements；
- paid promotion；
- Pitaya-style distributed architecture；
- direct Nakama/Pitaya API compatibility。

## 13. Verification

本 gate 的 repository check rule 是：

```text
runtime.friends_relationship_protocol_route_gate
```

推荐在本 gate 后运行：

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.friends_relationship_protocol_route_gate
node tools/vibit check change define-friends-relationship-protocol-route-gate --json
node tools/vibit check module friends --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

本 gate 不添加 Go runtime behavior，因此不强制要求 Go tests；但在关闭开发轮次前完整 runtime test run 仍然有用。
