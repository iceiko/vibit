# Protobuf Source Root 中文版

状态：Draft v0.1
最后更新：2026-05-13
范围：`proto/`
说明：本文件是 `proto/README.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本目录是 vibit Protobuf wire schemas 的计划 source root。

Business behavior 的 semantic source of truth 仍然位于 `contracts/` 和 module manifests。Protobuf files 定义 client/server wire message shape，并且必须与 semantic contract sources 保持一致。

## Layout

Module Protobuf sources 应使用：

```text
proto/vibit/<module>/v1/
```

第一版计划 module path 是：

```text
proto/vibit/inventory/v1/
```

生成的 Go Protobuf output 应放在：

```text
runtime/internal/generated/proto/
```

## Rules

- 在创建或修改 `.proto` files 前，运行 `node tools/vibit check protocol`。
- 不要手工编辑生成的 Go Protobuf output。
- Protobuf package names、message names、service names 和 field names 使用英文。
- Public wire schemas 必须显式 versioned。
- Transport adapters 应把 Protobuf wire messages 转换为 vibit commands 和 queries；domain modules 不拥有 Protobuf framing。

## Future Tooling

当 Protobuf generation 开始时，根目录配置应包含：

```text
buf.yaml
buf.gen.yaml
```

Buf linting、formatting、breaking checks 和 generation orchestration 已由 ADR-0013 接受。
