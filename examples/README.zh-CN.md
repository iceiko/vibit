# vibit Examples 中文版

状态：Draft v0.1
最后更新：2026-05-22

本文件是 `examples/README.md` 的简体中文译本。英文版本是权威版本。

本目录包含 source-first local examples 和 templates。这些文件不是 product SDKs、release artifacts、hosted demos、install scripts、package registry publications 或 direct Nakama/Pitaya compatibility surfaces。

## Local Alpha Request Loop

从 repository root 运行：

```bash
examples/local-alpha-request-loop.sh
```

该 script 包装 focused authenticated gameplay E2E proof：

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

它证明：

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

该 script 有意保持 redacted。它不得打印 raw credentials、raw access tokens、verifier keys、DSNs、digests、headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 concrete transport metadata。

## Local Environment Template

`examples/local.prototype.env.example` 是 local configuration fields 的 placeholder checklist。它不包含真实 secrets，也不应在未替换所有 placeholders 为 local-only values 前直接 source。

`.vibit.local.env`、`.env.local` 和 `.env.*.local` 等 private local env files 会被 repository `.gitignore` 忽略。不要提交 private local configuration。
