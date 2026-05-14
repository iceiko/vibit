# Authentication Generated Contract Shape Timing Standard 中文版

状态：Draft v0.1
最后更新：2026-05-14
范围：Runtime authentication semantic contracts 与 generated Go contract shape timing
说明：本文件是 `docs/authentication-generated-contract-shape-timing.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准应与 `docs/generated-output.md`、`docs/authentication-contract-error-permission-surfaces.md`、`docs/runtime-authentication-implementation-boundary.md`、`ADR-0017`、`ADR-0028`、`ADR-0036`、`ADR-0037` 和 `ADR-0038` 一起使用。

## 1. 目的

Runtime authentication 现在已有 semantic contracts、storage schema boundaries、storage-neutral repository interface、PostgreSQL adapter，以及 application-owned implementation boundary。

下一个风险是 agent 可能直接从 prose 和 repository code 开始写 service interfaces 或 runtime behavior，绕过已经帮助 inventory 和 player 工作保持可预测的 machine-readable contract shapes。

本标准决定 generated Go authentication contract shapes 的时机和边界，但本次 change 不生成这些文件。

## 2. 时机决策

Generated Go authentication contract shapes 应在 application authentication service interfaces 之前、runtime authentication behavior 之前引入。

推荐顺序是：

1. `contracts/runtime/authentication/` 下的 semantic authentication contracts。
2. Runtime authentication implementation boundary。
3. Authentication generated contract shape timing decision。
4. Runtime authentication family shapes 的 generator 与 check support。
5. Generated authentication contract shape files。
6. Application authentication service interface boundary。
7. Token generation、verifier comparison、login execution、token validation、logout execution、cleanup、protocol carriers 和 runtime behavior 通过后续 gated work 推进。

`W-0088` 只完成第 3 步。它不授权 generated files、service interfaces、handlers、token behavior、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository changes 或 migration schema changes。

## 3. Source And Output

允许的 source set 是已注册的 runtime authentication contract family：

```text
contracts/runtime/authentication/commands/*.yaml
contracts/runtime/authentication/events/*.yaml
contracts/runtime/authentication/errors/*.yaml
contracts/runtime/authentication/permissions/*.yaml
```

Registry source 是 `.arch/contracts.yaml`，位于 runtime 的 `authentication` family 下。

计划的 output root 是：

```text
runtime/internal/generated/contracts/runtime/authentication/
```

计划的文件形态是：

```text
runtime/internal/generated/contracts/runtime/authentication/<contract-type>/<ContractID>.go
```

示例：

```text
runtime/internal/generated/contracts/runtime/authentication/commands/AuthenticateWithDeviceCredential.go
runtime/internal/generated/contracts/runtime/authentication/events/TokenIssued.go
runtime/internal/generated/contracts/runtime/authentication/errors/authentication_errors.go
runtime/internal/generated/contracts/runtime/authentication/permissions/authentication_permissions.go
```

Runtime contract families 需要 family segment，因为 `runtime` 可能拥有多个 semantic families，例如 `session` 和 `authentication`。

第一批 authentication shape files 计划使用的 Go package name 是：

```text
runtimeauthenticationcontracts
```

## 4. Immutability

Generated authentication contract shape files 对 non-system agents 不可变。

如果 generated authentication contract shape output 错了，agents 必须修改以下 source 之一，而不是 patch generated files：

- `contracts/runtime/authentication/` 下的 semantic contract source。
- `.arch/contracts.yaml` 中的 contract registry。
- Generator。
- Generated-output standard。
- 相关 change spec 或 ADR，如果它明确授予 `generated_file_override`。

本标准不授予任何 `generated_file_override`。

## 5. Check Requirements

在提交 generated authentication contract shapes 之前，repository tooling 必须支持这些 checks：

- `node tools/vibit generate contract-shapes all` 能从 semantic contracts 生成 runtime authentication family。
- `node tools/vibit check generated --json` 能发现 runtime authentication family shapes 的 missing、stale 或 drift。
- `node tools/vibit inspect generated --json` 能用 machine-readable form 报告 runtime authentication generated shape status。
- `node tools/vibit check contracts --json` 仍然在 generation 前验证 semantic contract source 和 registry。
- `node tools/vibit check runtime --json` 仍然证明 generated authentication shapes 没有添加 runtime authentication behavior。
- `runtime.selected_login_token_boundary` 和 `runtime.authentication_implementation_boundary` 能区分 metadata-only generated shapes 与 runtime authentication implementation。
- `node tools/vibit check all --json` 包含 generated-shape checks。

第一项 generation work item 必须在生成 output 之前或同时更新 checks。它不得削弱 selected login/token、authentication implementation、generated-output、protocol、WebSocket、dependency、repository 或 migration guards。

## 6. 与 Runtime Behavior 的关系

Generated authentication contract shapes 只是 metadata-only。

它们可以指导 naming、service interface planning、tests 和未来 protocol mapping，但不得实现：

- Login execution。
- Token generation。
- Token verifier comparison。
- Access-token validation。
- Logout execution。
- Refresh behavior。
- Cleanup jobs。
- WebSocket routing。
- Protobuf envelope behavior。
- Persistence。
- Domain dispatch。

Application service interfaces 可以在 generated shapes 存在后设计，但 service code 仍然必须单独 gated。

## 7. Reference Alignment

Nakama 和 Pitaya 仍然是 capability 与 vocabulary references。它们说明 production game backends 需要稳定的 authentication、session、request 和 route concepts。

vibit 应把这个经验转化为 agent-native structure：

- Semantic contracts 定义 names 和 intent。
- Generated contract shapes 让 agents 可以检查 contract surface。
- Handwritten application logic 保持分离，并且必须显式 gated。
- Wire protocol 和 WebSocket carriers 仍然是单独 decisions。

本标准不复制 Nakama 或 Pitaya public APIs。

## 8. Verification

对本 timing decision，verification 应包括：

```bash
node tools/vibit check contracts --json
node tools/vibit inspect generated --json
node tools/vibit check generated --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change decide-authentication-generated-contract-shape-timing --json
node tools/vibit check all --json
git diff --check
```

本 timing decision 不添加或修改 Go runtime behavior，因此不要求 runtime Go tests。

后续 generation work item 在生成文件后，必须运行 generator、generated-output、runtime 和完整 repository checks。

## 9. Migration Path

Migration path 是：

1. 记录本 timing decision 和 output boundary。
2. 添加一个 bounded work item，授权 generator/check support 和 generated authentication shape output。
3. 为 runtime family-aware paths 扩展 generated-output tooling。
4. 通过 tooling 生成文件，不要手写。
5. 验证 source trace、drift、stale files 和 runtime behavior boundaries。
6. 之后再设计 application authentication service interfaces。

