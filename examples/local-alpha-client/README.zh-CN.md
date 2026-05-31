# Local Alpha Example Client Path 中文版

状态：Draft v0.1
最后更新：2026-05-31

本文件是 `README.md` 的简体中文译本。英文版本是权威版本。

本目录记录 vibit 的第一条 source-first local alpha example client path。它是 repository-local example path，不是 public SDK、generated client library、package publication、hosted demo、release artifact、install script、live external client guarantee，也不是 direct Nakama/Pitaya API compatibility surface。

## 运行

从 repository root 运行：

```bash
examples/local-alpha-example-client.sh
```

该脚本会在 Go runtime 内运行 focused local alpha Protobuf E2E proof：

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow|TestFriendsRelationshipProtocolRouteLocalAlphaFlow|TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation|TestAuthenticatedGameplayFailurePathsLocalAlphaFlow' -v
```

## 证明流程

该 example path 展示当前 source-first local alpha loop：

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> protected own-player storage object put/get/list/delete
-> protected friends send/status/accept/list/remove/block/unblock/reject
-> presence online/offline proof after close and invalidation
-> logout
-> rejected post-logout protected request
-> protected request failure-path and redaction checks
```

它使用现有 runtime 和 protocol behavior：

```text
runtime.authentication.AuthenticateWithDeviceCredential
runtime.connection.BindConnection
inventory.GrantItem
inventory.GetInventory
presence.GetPlayerPresence
storage.PutOwnStorageObject
storage.GetOwnStorageObject
storage.ListOwnStorageObjects
storage.DeleteOwnStorageObject
friends.SendFriendRequest
friends.GetFriendRelationshipStatus
friends.AcceptFriendRequest
friends.ListFriendRelationships
friends.RemoveFriend
friends.BlockPlayer
friends.UnblockPlayer
friends.RejectFriendRequest
runtime.authentication.LogoutAccessToken
```

## 它是什么

这条路径面向需要理解当前 alpha capability loop 的开发者和 AI agents，避免他们先去反向理解内部 E2E test。它保持 source-first，因为当前 generated Protobuf Go output 位于 `runtime/internal/` 下，而且 local onboarding 仍是 application-owned setup behavior，不是 public client route。

这条路径对应 Nakama 的 developer-experience 压力：backend framework 应该让开发者容易理解 client 如何执行 authentication、protected gameplay requests、player storage、friends relationship social graph behavior、presence、logout 和 failure behavior。vibit 采纳这种产品压力，但不复制 Nakama 的 public routes、payloads、SDK shapes、token semantics、runtime APIs 或 compatibility promises。

## 它不是什么

这条路径不会：

- publish client SDK；
- generate client libraries；
- 添加 public onboarding protocol route；
- 添加 protocol routes 或 Protobuf messages；
- 修改 generated output；
- 修改 runtime、authentication、session、persistence、migration、startup 或 transport behavior；
- 添加 dependencies；
- 添加 hosted deployment 或 release artifacts；
- 添加 stream subscriptions、chat rooms、groups、parties、broadcast fanout、delivery guarantees、matchmaking、match runtime、operations/admin behavior 或 distributed runtime；
- 添加 direct Nakama/Pitaya API compatibility。

Pitaya 仍作为 future distributed architecture reference 延后。本路径不会引入 frontend/backend server roles、RPC、service discovery、groups、cluster routing 或 distributed sessions。

## Redaction

脚本和 example docs 不得打印、持久化、提交或记录：

- raw device credential text or bytes；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- verifier key values；
- concrete verifier key set ids；
- HMAC inputs 或 outputs；
- 带 credentials 的 PostgreSQL DSNs；
- headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 concrete transport metadata。

允许输出的内容仅限 route names、step names、test names、redacted status classes 和 high-level success/failure descriptions。

## Verification

example path 本身通过以下命令验证：

```bash
examples/local-alpha-example-client.sh
```

Architecture 和 continuation verification 还应包括：

```bash
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.local_alpha_example_client_path_implementation
node tools/vibit inspect rule runtime.friends_relationship_protocol_route_local_proof
node tools/vibit check change implement-local-alpha-example-client-path --json
node tools/vibit check change prove-friends-relationship-protocol-route-local-flow --json
node tools/vibit check runtime --json
node tools/vibit check all --json
git diff --check
```
