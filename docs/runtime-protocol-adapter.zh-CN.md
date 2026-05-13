# Runtime Protocol Adapter Boundary Standard 中文版

状态：Draft v0.1
最后更新：2026-05-13
范围：WebSocket transport、Protobuf protocol adaptation、application dispatch 和 domain modules 之间的 runtime boundary
说明：本文件是 `docs/runtime-protocol-adapter.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准定义 vibit WebSocket Protobuf server path 的第一版 runtime handoff shape。

本标准应与 `.arch/runtime.yaml`、`.arch/protocol.yaml`、`docs/game-protocol.md`、`docs/generated-output.md`、`ADR-0014`、`ADR-0015`、`ADR-0016`、`ADR-0017` 和 `ADR-0018` 一起使用。

## 1. 目的

vibit 的第一版 runtime path 必须让 agents 容易理解。

主要风险不是写不出 WebSocket server。主要风险是 agents 把 validation、routing、permission checks、session shortcuts、payload decoding 或 domain behavior 放到当时最方便的层里。

本标准在 runtime implementation 开始前定义层与层之间的窄交接，防止这种漂移。

## 2. First Runtime Flow

第一版 gameplay request flow 是：

```text
websocket frame
-> transport frame
-> protocol envelope
-> route request
-> application dispatch
-> domain command or query handler
-> application result
-> protocol envelope
-> websocket frame
```

第一版 server-push flow 是：

```text
domain event
-> application publication decision
-> protocol envelope
-> transport send
```

Server-push persistence 和 outbox behavior 继续 deferred，直到单独的 event delivery decision 存在。

## 3. Layer Ownership

### WebSocket Transport Adapter

Owner：

```text
runtime/internal/platform/transport/ws/
```

职责：

- 拥有 `github.com/coder/websocket`。
- 接受和关闭 WebSocket connections。
- 读写 binary frames。
- 在定义后执行 transport-level size、close 和 lifecycle behavior。
- 把 opaque frame bytes 传给 protocol adapter。
- 发送 protocol adapter 返回的 encoded frame bytes。

不得：

- 解析 domain payloads。
- 直接 dispatch 到 domain modules。
- 解释 command、query、event、permission 或 invariant semantics。
- 构造 module-specific Protobuf payloads。
- 修改 domain state。

第一条 active adapter 是：

```text
runtime/internal/platform/transport/ws/server.go
```

当前行为：

- 暴露兼容 `http.Handler` 的 `Server`。
- 通过 `github.com/coder/websocket` 接受 WebSocket connections。
- 只把 client messages 当作 binary frames 读取。
- 在传给注入的 `FrameHandler` 前复制 frame payload bytes。
- 通过 transport-owned `Frame` type 提供 connection metadata。
- 把每一个 handler response 都写成 binary WebSocket frame。
- 对 text 或其他 non-binary client messages 执行 transport close。
- 由 `runtime/cmd/vibit-server` 挂载到 `/v1/ws`。

该 adapter 必须保持 opaque。它不得 import generated Protobuf packages、application dispatch 或 domain modules。Protocol decoding 和 application routing 属于后续在该 transport package 之外完成的 composition。

### Protobuf Protocol Adapter

Owner：

```text
runtime/internal/platform/protocol/protobuf/
```

职责：

- 拥有 Protobuf envelope encode 和 decode。
- 验证 envelope-level shape。
- 把 `kind`、`module`、`name`、`request_id`、`target`、`session`、`payload_type` 和 `payload` 转换为 application route request。
- 只通过 generated Protobuf packages 解码 generated Protobuf payloads。
- 只通过显式 protocol bridge functions，把 generated wire payloads 转换成 handwritten domain runtime payloads。
- 只通过显式 protocol bridge functions，把 application results 和 events 转换回 generated wire payloads。
- 把 protocol errors 映射为 error envelopes。
- 保留 request correlation。

不得：

- 打开 WebSocket connections。
- 在 envelope metadata interpretation 之外拥有长期 player sessions。
- 执行 domain permission decisions。
- 执行 domain invariants。
- 直接调用 domain repositories。
- 把 business behavior 隐藏在 envelope conversion 中。

当前行为：

- `runtime/internal/platform/protocol/protobuf/frame_handler.go` 提供第一版 frame composition adapter。
- 它通过 Protobuf-owned `FrameRequest` type 接收 frame payload bytes 和 transport metadata。
- 它把 frame payload 解码为 `vibit.protocol.v1.Envelope`。
- 它通过显式 Protobuf/domain bridge，把 command 和 query envelopes 转换为 `app.RouteRequest`。
- 它通过注入的 application dispatcher interface 执行 dispatch。
- 它把成功的 application results 编码为 Protobuf envelopes。
- 它把 `app.ApplicationError` results 编码为 `MESSAGE_KIND_ERROR` envelopes。
- 它返回 encoded envelope bytes，供 WebSocket transport 写回。

该 adapter 有意不 import WebSocket transport package。未来 process wiring layer 可以把 `ws.Frame` 适配为这个 Protobuf-owned `FrameRequest`，同时不把 Protobuf 或 application knowledge 移入 transport package。

### Application Dispatch

Owner：

```text
runtime/internal/app/
```

职责：

- 拥有 command 和 query dispatch。
- 把 route requests 匹配到已登记的 application handlers。
- 为 state-changing commands 创建或调用 unit-of-work boundaries。
- 通过 vibit-owned interfaces 调用 domain module handlers。
- 把 application results 映射为 protocol response metadata。
- 保持 route registration 显式；当 generation 存在时应由 generation 产生。

不得：

- 解析 WebSocket frames。
- 拥有 Protobuf wire framing。
- Import `github.com/coder/websocket`。
- Import generated Protobuf packages，除非后续 adapter decision 明确允许一个窄 bridge。
- 隐藏 module-specific business rules。

第一版 process bootstrap helper 是：

```text
runtime/internal/app/bootstrap/inventory.go
```

它创建第一版 runtime process 使用的 in-memory inventory dispatcher。这个 helper 属于 application composition，不属于 domain behavior。它可以注册 module handlers 并选择临时 bootstrap dependencies，但不得变成长期 persistence、authentication 或 permission model。

### Runtime Process Wiring

Owner：

```text
runtime/cmd/vibit-server/
```

职责：

- 读取 listen address 等 process configuration。
- 组装 bootstrap application dependencies。
- 把 Protobuf frame handler 与 WebSocket transport adapter 组合起来。
- 挂载 `/v1/ws`。
- 启动并拥有 HTTP server lifecycle。

不得：

- 把 business behavior 隐藏进 process startup。
- 直接解码 Protobuf payloads。
- 执行 domain permissions 或 invariants。
- 拥有 authentication/session semantics。
- 在调用 bootstrap assembly 之外拥有 persistence repositories。

第一条 active process entrypoint 是：

```text
runtime/cmd/vibit-server/main.go
```

Manual startup 和 endpoint verification 记录在：

```text
docs/runtime-runbook.md
```

### Protocol-To-Domain Payload Bridges

Owner：

```text
runtime/internal/platform/protocol/protobuf/
```

第一条 active bridge 是：

```text
runtime/internal/platform/protocol/protobuf/inventory_bridge.go
```

职责：

- 在 application dispatch 前，把 decoded generated Protobuf payloads 映射为 handwritten module runtime request structs。
- 在 application dispatch 后，把 application result payloads 映射为 generated Protobuf response payloads。
- 当 events 准备进入 protocol output 时，把 application events 映射为 generated Protobuf event payloads。
- 保持 field mapping 显式并由 tests 覆盖。
- 保留原始 envelope metadata 和 request correlation。

不得：

- 执行 domain permissions 或 invariants。
- 直接调用 repositories。
- 添加 authentication shortcuts。
- 修改 generated Protobuf output。
- 成为隐藏 business behavior 的地方。

规则：

- Domain modules 不得 import generated Protobuf packages。
- `runtime/internal/app/` 不得 import generated Protobuf packages。
- 在后续 generated bridge standard 替代 handwritten bridge 之前，bridge 属于 protocol adapter layer。
- 未知或尚未 bridge 的 routes 只有在 payload 已经是 protocol payload 时才可以原样通过。Inventory routes 遇到 payload type mismatch 时必须 fail fast。

### Domain Module Runtime Logic

Owner：

```text
runtime/internal/modules/<module>/
```

职责：

- 实现手写 command 和 query behavior。
- 执行 module invariants。
- 使用 vibit-owned repository 和 policy interfaces。
- 作为 server facts 发出 domain events。
- 与 module manifests 和 semantic contract sources 保持一致。

不得：

- Import WebSocket libraries。
- 拥有 Protobuf framing。
- 解析 envelopes。
- 直接依赖 `google.golang.org/protobuf`。
- 访问其他 modules 的内部实现。

### Generated Code

Owners：

```text
runtime/internal/generated/contracts/
runtime/internal/generated/proto/
```

职责：

- 提供 generated contract 和 wire shapes。
- 能追溯到 source contracts 或 `.proto` files。
- 对 non-system agents 保持不可变。

不得：

- 包含 handwritten runtime logic。
- 拥有 transport、application dispatch 或 domain behavior。

## 4. First Handoff Types

第一版 runtime implementation 应先引入 vibit-owned handoff types，再绑定 transport 或 generated Protobuf packages。

概念类型名是：

```text
TransportFrame
ProtocolEnvelope
RouteRequest
ApplicationResult
OutboundMessage
```

第一批 Go runtime slices 已经在 `runtime/internal/app/` 下实现 application-owned `RouteRequest` 和 `ApplicationResult` concepts、面向 command 和 query routes 的显式 application dispatcher，在 `runtime/internal/platform/protocol/protobuf/` 下实现 Protobuf-to-application conversion 和 Protobuf-owned `FrameRequest`、`FrameHandler` composition adapter，并在 `runtime/internal/platform/transport/ws/` 下实现 transport-owned `Frame` handoff。其余 concepts 实现时可以使用符合 Go 习惯的名称，但必须保留这些职责。

必需概念：

- `TransportFrame` 携带 frame bytes 和 connection metadata，但不携带 domain semantics。
- `ProtocolEnvelope` 表示 decoded envelope metadata 和 payload bytes。
- `RouteRequest` 携带 `kind`、`module`、`name`、`request_id`、target、session、payload identity 和 decoded command/query payload。
- `ApplicationResult` 携带 response payloads、emitted events 和 application-level errors。
- `OutboundMessage` 携带可由 transport encoding 发送的 protocol-level output。

## 5. Routing Rules

Runtime routing 必须使用 structured route fields：

```text
kind
module
name
```

渲染后的 route key 可以是：

```text
<module>.<name>
```

Route registration 必须显式。当 generators 存在时，route registration 应从 contracts 和 manifests 生成。在 generators 存在前，手写 route registration 必须小、局限于 application dispatch，并由 change specs 覆盖。

当前手写 application dispatcher 只 dispatch `command` 和 `query` route requests。`event`、`error`、`system`、`ack`、`heartbeat`、`input` 和 `state` messages 在后续 standard 定义其 lifecycle 前，不属于 application-dispatchable messages。

Transport handlers 不得从 WebSocket paths 或 message text 构造 ad hoc route strings。

## 6. Error Boundaries

Errors 应在拥有该失败的层映射：

- Transport failures 属于 transport adapter。
- Malformed envelope 和 unknown payload type failures 属于 protocol adapter。
- Unknown route failures 属于 application dispatch。
- Permission、invariant 和 domain validation failures 属于 domain modules 和 application policy boundaries。
- Internal failures 不应把 implementation details 泄漏到 public error payloads。

Public module errors 必须映射到已登记的 error catalogs。

Application errors 由 Protobuf protocol adapter 编码为 `MESSAGE_KIND_ERROR` envelopes。第一条 active mapper 是：

```text
runtime/internal/platform/protocol/protobuf/error_envelope.go
```

规则：

- 保留 application result 中的 `request_id`、route metadata、target metadata 和 session metadata。
- 把 stable application error code 和 public message 复制到 `protocolv1.Error`。
- 把 `Error.request_id` 设置为关联的 request id。
- Error envelopes 中保持 `payload_type` 和 `payload` 为空。
- 在 retryability 从已登记 error catalogs 生成之前，application errors 默认 non-retryable。
- 没有单独的 protocol error handling decision 时，不要通过这个 mapper 暴露 non-application internal errors。

## 7. Session And Target Boundaries

Protocol adapter 可以解析 session 和 target metadata，但不得发明 authentication shortcuts。

在 player 或 auth modules 存在前：

- `player_id` 可以作为 inventory protocol shape 的 planned context。
- Authentication/session validation 仍是 deferred explicit module 或 platform decision。
- Transport connection identity 不得被当作 durable player identity。

`player` 之外的 target scopes 继续 reserved，直到相关 module 和 lifecycle standard 存在。

## 8. Agent Rules

Agents 必须：

- 在添加 WebSocket transport code、Protobuf runtime adapter code、application dispatch code 或 domain runtime handlers 前阅读本标准。
- 保持 transport、protocol、application、domain 和 generated output 职责分离。
- 当新的 boundary rule 可机器验证时，新增或更新 checks。
- 记录不可用的 Go、Buf 或 Protobuf tooling，而不是声称已经运行 generation 或 runtime tests。

Agents 不得：

- 把 business behavior 放进 WebSocket handlers。
- 把 domain permission 或 invariant logic 放进 Protobuf envelope decoding。
- 让 domain modules 直接 import WebSocket 或 Protobuf runtime libraries。
- 让 generated output 变成手写 adapter layer。
- 在 inventory protocol work 中添加 authentication shortcuts。

## 9. Verification

当前 verification：

```bash
node tools/vibit check runtime
node tools/vibit check protocol
node tools/vibit check generated
node tools/vibit check all
```

`check runtime` 应在 Go runtime implementation 开始前，验证 boundary standard、runtime manifest、protocol manifest、runtime agent guide 和 repository guide 都指向 runtime protocol adapter boundary。

当 Go source files 存在时，`check runtime` 也会检查第一批 Go import 和 layer boundaries。Verification 应继续确保：

- `github.com/coder/websocket` imports 只位于 `runtime/internal/platform/transport/ws/`。
- Protobuf runtime imports 只位于 generated Protobuf packages 和 protocol adapters。
- Domain modules 不直接 import transport 或 Protobuf libraries。
- Application dispatch 不解析 WebSocket frames。
- Application 和 domain packages 不得直接 import platform adapters 或 generated Protobuf packages。
- Generated output 不包含 handwritten adapter code。
- Inventory Protobuf/domain bridge code 保持位于 `runtime/internal/platform/protocol/protobuf/`。

## 10. Migration Path

Runtime implementation 前：

1. 保持 package directories 为 skeletons。
2. 定义 runtime protocol adapter boundary。
3. 保持 `.proto` sources 和 generated-output rules 稳定。
4. 只有在 Go、Buf 和 Protobuf tooling 可用时，才添加 toolchain-dependent generated output。

Implementation 开始时：

1. 先添加窄 Go handoff types。
2. 在 wiring WebSocket transport 前添加 protocol adapter tests。
3. 在 domain runtime behavior 扩大前添加 application dispatch tests。
4. 保持第一版 inventory slice 小而 player-scoped。

当前 implementation progress：

1. Narrow Go handoff types 已经覆盖 application requests/results、transport frames 和 Protobuf frame composition。
2. Protocol adapter tests 已覆盖 envelope conversion、inventory payload bridging、error envelope mapping 和 frame composition。
3. Application dispatch tests 已覆盖 command/query routing 和 application errors。
4. `/v1/ws` endpoint mounting 已经存在于 `runtime/cmd/vibit-server`。
