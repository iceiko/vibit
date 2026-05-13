# Protobuf Source Root 中文版

状态：Draft v0.1
最后更新：2026-05-13
范围：`proto/`
说明：本文件是 `proto/README.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本目录是 vibit Protobuf wire schemas 的 source root。

Business behavior 的 semantic source of truth 仍然位于 `contracts/` 和 module manifests。Protobuf files 定义 client/server wire message shape，并且必须与 semantic contract sources 保持一致。

Game protocol envelope standard 由 `docs/game-protocol.md`、`.arch/protocol.yaml` 和 `ADR-0015` 定义。创建 `.proto` files 前必须阅读这些 artifacts。第一版 protocol model 是 WebSocket-framed Protobuf envelope，使用显式 `kind`、`module` 和 `name` routing fields，并包含 session metadata、target metadata、server-authoritative message rules 和 error mapping。

## Layout

Module Protobuf sources 应使用：

```text
proto/vibit/<module>/v1/
```

第一版计划 module path 是：

```text
proto/vibit/inventory/v1/
```

第一版 protocol envelope source 是：

```text
proto/vibit/protocol/v1/envelope.proto
```

生成的 Go Protobuf output 应放在：

```text
runtime/internal/generated/proto/
```

Generated output rules 定义在：

```text
docs/generated-output.md
```

根目录 Protobuf generation 配置是：

```text
buf.yaml
buf.gen.yaml
```

`buf.yaml` 声明 `proto/` 为 source root，并使用 standard linting 和 file-level breaking checks。`buf.gen.yaml` 声明官方 Go Protobuf generation plugin 和计划中的 generated output path。

## Rules

- 在创建或修改 `.proto` files 前，运行 `node tools/vibit check protocol`。
- 保持 envelope 和 module payload schemas 与 `.arch/protocol.yaml` 对齐。
- 不要手工编辑生成的 Go Protobuf output。
- 不要在 `runtime/internal/generated/proto/` 下创建手写 runtime code。
- `runtime/internal/generated/proto/` 下的 generated Go Protobuf files 必须是 `*.pb.go`，包含 `protoc-gen-go` generated-code marker，并包含能解析到现有 `.proto` file 的 source trace。
- Protobuf package names、message names、service names 和 field names 使用英文。
- Public wire schemas 必须显式 versioned。
- 保持 `option go_package` 与 `runtime/internal/generated/proto/` 对齐。
- Transport adapters 应把 Protobuf wire messages 转换为 vibit commands 和 queries；domain modules 不拥有 Protobuf framing。
- 在相关 modules 和 standards 存在前，不要在 Protobuf files 中实现 room state sync、matchmaking、allocation、reconnect replay、presence、streams、realtime input 或 state patches。

## Future Tooling

Buf linting、formatting、breaking checks 和 generation orchestration 已由 ADR-0013 接受。

只有在本地 Buf CLI 与 Go Protobuf generator 可用，并且 generated-output rules 对本次变更足够稳定时，才应运行本地 generation。如果本地 toolchain 不可用，应把 generation 记录为 not verified，而不是手工创建 generated files。

计划中的 generation command：

```text
buf generate
```

Generation 后运行：

```text
node tools/vibit check generated
```
