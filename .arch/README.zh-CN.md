# Architecture Manifests 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：vibit 的机器可读架构入口  
说明：本文件是 `.arch/README.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

这个目录包含面向 agents、humans、generators 和未来 verification commands 的 architecture manifests。

这些 manifests 不是装饰性文档。它们的目标是成为可执行的架构上下文。

## 目的

`.arch/` 目录应该回答 agent 在修改代码前必须解决的问题：

- 存在哪些 modules？
- 每个 module 拥有什么？
- 哪些 dependencies 被允许？
- 哪些 contracts 定义 public behavior？
- 哪些 events、commands、queries、errors 和 permissions 已登记？
- 哪些文件是 generated？
- 哪些 tests 或 checks 证明 architecture rules？

## 当前文件

```text
.arch/
  README.md
  README.zh-CN.md
  modules.yaml
  conventions.yaml
  runtime.yaml
  contracts.yaml
```

这是第一版草案。这些文件在实现代码存在前先描述预期形态。

`runtime.yaml` 记录第一版 Go server runtime 方向的 runtime readiness decisions。它指向约束第一语言、服务器实例模型、contract boundary、client protocol、wire format、persistence direction、dependency adoption 和 proof slice 的 Agent Decision Records。

`contracts.yaml` 登记 public command、query、event、error 和 permission contract source files。Contract files 位于 `contracts/` 下，是 semantic source artifacts，不是 generated output。Protobuf wire schemas 计划位于 `proto/` 下，并且必须与这些 semantic contracts 对齐。

## 未来预期文件

```text
.arch/dependencies.yaml
.arch/test-matrix.yaml
.arch/generation.yaml
```

当第一版 prototype 需要它们时，项目应逐步添加这些文件。

## Agent 规则

修改实现代码前，agents 应：

1. 阅读 `CONSTITUTION.md`。
2. 阅读 `AGENTS.md`。
3. 阅读 `.arch/modules.yaml`。
4. 阅读 `.arch/conventions.yaml`。
5. 在修改或创建 runtime implementation code 前，阅读 `.arch/runtime.yaml`。
6. 在新增或修改 public contracts 前，阅读 `.arch/contracts.yaml`。
7. 在相关 module 存在时，阅读其 `module.yaml`。
8. 当 public architecture 变化时，先更新 manifests，再实现。

如果 manifest 缺少安全变更所需的信息，应更新 manifest 或记录这个缺口。

## 验证方向

这些 manifests 最终应支持类似命令：

```bash
vibit check architecture
vibit check module <module>
vibit check contracts
vibit check change <change-id>
```

在这些命令存在前，agents 必须记录 architecture verification 当前不可用。
