# Dependency Adoption Standard 中文版

状态：Draft v0.1
最后更新：2026-05-12
范围：Foundational runtime and tooling dependencies
说明：本文件是 `docs/dependency-adoption.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文定义 vibit 在把 foundational dependencies 纳入架构前如何评估它们。

本标准应与 `ADR-0010` 一起使用。

## 1. 目的

当成熟开源项目能够降低风险并提升长期质量时，vibit 应使用它们。但不能仅因为某个 dependency 流行或方便，就随意采用它。

Dependency adoption records 让未来 agents 可以检查 dependency decisions。它们解释为什么使用某个 package、它允许出现在哪里、由什么边界封装，以及如何替换。

## 2. 何时需要

在新增或要求任何会影响以下方面的 dependency 之前，必须先有 dependency adoption record：

- Transport
- Protocol generation
- Persistence
- Migration tooling
- Object storage
- Dispatch
- Module loading
- Lifecycle management
- Observability
- Testing strategy
- Generated code conventions

如果小型 developer-only utilities 不影响架构、generated output、public behavior 或 runtime shape，可以在 change spec 中用更轻量的说明。

## 3. 位置

机器可读 dependency status 位于：

```text
.arch/dependencies.yaml
```

可复用 adoption template：

```text
docs/_templates/dependency-adoption.md
```

最终 adoption record 可以是：

- `decisions/` 下的专门 ADR。
- Change spec 中的 dependency section，当该 dependency 很小且范围很窄时适用。

Foundational dependencies 通常应使用专门 ADR。

## 4. 必需评估项

Adoption record 必须评估：

- Problem solved
- Package or tool identity
- Ecosystem role
- Maintenance activity
- License compatibility
- API stability
- Production adoption signals
- Security and supply-chain risk
- Operational fit
- Agent readability
- Testability
- Generated-code compatibility
- Abstraction boundary
- Allowed import or invocation locations
- Forbidden import or invocation locations
- Replacement path
- Verification path

Stars、使用频率和声誉是有用信号，但它们本身不充分。

## 5. 边界规则

Domain modules 不得直接 import 或 invoke transport、protocol、persistence、object storage 或 framework behavior 相关的 foundational third-party dependencies。

推荐 ownership：

- Transport libraries 属于 platform transport adapters。
- Protobuf tooling 属于 generation tooling 和 generated protocol packages。
- PostgreSQL drivers 属于 platform persistence adapters。
- Migration tools 属于 platform migration tooling。
- S3 SDKs 和 MinIO clients 属于 platform object-storage adapters。
- Test frameworks 属于 test infrastructure，而不是 domain logic。

Domain modules 应依赖 vibit-owned interfaces、generated clients、contract types、repositories、policies 和 service abstractions。

## 6. Dependency Status Values

`.arch/dependencies.yaml` 使用以下 status values：

```text
proposed
accepted
deferred
rejected
superseded
```

当 dependency slot 已知但尚未选择实现时，使用 `proposed`。

只有 adoption record 完整且已链接后，才使用 `accepted`。

当 dependency category 真实存在但下一个实现步骤不需要它时，使用 `deferred`。

当某个合理候选不应被选择时，使用 `rejected`。

当之前的 dependency choice 被替换时，使用 `superseded`。

## 7. Agent 规则

Agents 必须：

- 添加 foundational dependencies 前阅读 `.arch/dependencies.yaml`。
- 在 dependency status 变为 `accepted` 前，创建或更新 dependency adoption record。
- 将 dependencies 保持在声明过的 abstraction boundary 后面。
- 当 dependency boundaries 影响 module work 时，更新 AGENTS guides。
- 修改 dependency records 后运行相关 verification commands。

Agents 不得：

- 把 foundational dependency 直接添加到 domain module code。
- 把流行度当成完整评估。
- 把 dependency decisions 藏在 implementation code 里。
- 在没有 replacement path 或 verification path 的情况下接受 dependency。

## 8. 验证方向

当前验证基于文档和 manifest：

```bash
node tools/vibit check architecture
node tools/vibit check schemas
node tools/vibit check memory
node tools/vibit check all
```

未来 checks 应验证：

- 每个 `accepted` foundational dependency 都有 adoption record。
- Domain modules 不 import forbidden dependency packages。
- 只有 platform adapters 可以直接拥有相关 dependencies。
- Generated files 声明它们的 generator dependencies。

## 9. Open Questions

- `.arch/dependencies.yaml` 未来是否应完整 schema-validated？
- Dependency records 是否应包含 machine-readable package coordinates？
- Replacement path 和 license review 是否应成为必需结构化字段？
- Runtime code 存在后，dependency checks 是否应检查 Go imports？
