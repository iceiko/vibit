# Storage Objects Behavior Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-22
范围：定义 inventory proof slice 之外第一版 general storage-object behavior 的 gate
依赖：`docs/prototype-ready-local-development-path-package.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/reference-game-server-alignment.md`
权威决策：`ADR-0109`

本文件是 `docs/storage-objects-behavior-gate.md` 的简体中文译本。英文版本是权威版本。

本文定义第一版 storage objects behavior gate。它是 gate artifact。它不实现 storage objects runtime behavior，不添加 protocol routes，不添加 Protobuf source 或 generated output，不添加 migrations，不添加 dependencies，不添加 repository interfaces，不添加 storage adapters，不扩展 operations/admin behavior，不添加 hosted deployments，不创建 release artifacts，不执行 public announcements，不运行 paid promotion，不改变 authentication/session behavior，不实现 broad product module，也不添加 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

storage objects behavior gate 记录如下：

```yaml
storage_objects_behavior_gate: defined
completed_work_item: W-0201
decision: ADR-0109
check_rule: runtime.storage_objects_behavior_gate
source_package_decision: ADR-0108
source_package_standard: docs/prototype-ready-local-development-path-package.md
gate_standard: docs/storage-objects-behavior-gate.md
gate_standard_translation: docs/storage-objects-behavior-gate.zh-CN.md
target_stage: prototype_ready_foundation
reference_capability_family: storage_objects_and_durable_game_state
first_scope_posture: player_owned_small_json_objects
object_identity_tuple: owner_kind_owner_id_collection_key
ownership_posture_recorded: true
scope_key_posture_recorded: true
read_write_semantics_recorded: true
permission_posture_recorded: true
conflict_semantics_recorded: true
protocol_expectations_recorded: true
data_expectations_recorded: true
verification_expectations_recorded: true
stop_conditions_recorded: true
future_schema_gate_work_item: W-0202
future_schema_gate_direction: storage_objects_persistence_schema_gate
gate_only: true
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
repository_interface_changed: false
storage_adapter_changed: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

Inventory proof slice 有价值，但它是 module-specific 的。Prototype authors 还需要一个 general durable state surface，用来保存 loadouts、preferences、tutorial state、quest flags、profile fragments、saved selections，或还不值得变成一等模块的 prototype-specific state。

第一版 storage objects behavior 应让 vibit 对 prototype 更有用，同时保留 agent-native 模型：

- ownership 是显式的；
- object identity 是稳定的；
- permissions 是 route-scoped 且 fail closed 的；
- writes 是 server-authoritative 且 version-aware 的；
- data shape 有边界并且 redaction-aware；
- persistence 仍在 repository 和 adapter gates 之后；
- protocol shape 跟随 behavior，而不是先发明 wire messages。

本 gate 把 Nakama 的 broad durable storage-object capability 作为产品参考，把 Pitaya 的 handler/persistence separation 作为架构参考。它不复制 Nakama public routes、data models、permission values、version strings 或 client API names，也不复制 Pitaya routing 或 cluster internals。

## 3. Storage Object Boundary

在本 gate 中，storage object 指由 server framework 拥有、通过稳定逻辑身份寻址的小型 durable game-state record。

第一版 posture 是：

```text
owner_kind: player
owner_id: authenticated player id
collection: server-defined logical namespace
key: server-defined or validated object key
value: bounded JSON object payload candidate
version: server-issued optimistic concurrency token
```

这不是 large binary object storage、asset storage、replay storage、CDN-backed content、S3-compatible blob storage、file uploads 或任意 document database adoption。`.arch/runtime.yaml` 中已有的 `object_storage` planning entries 保持独立且 deferred。

## 4. Ownership And Scope

第一版 behavior gate 只选择 player-owned objects。

允许的第一版 scope：

- `player` owner kind；
- 每个第一版 object 由一个 authenticated player 拥有；
- owner id 必须来自 validated request identity，不能只来自 client-supplied metadata；
- collection 和 key 在 player scope 内识别 object。

Deferred scopes：

- global objects；
- group、guild、party、room、match 或 server-shard objects；
- cross-player shared writes；
- public catalog objects；
- admin-managed objects；
- object ACL lists；
- object search across arbitrary owners。

这些 scope 会改变 permissions、indexing、operations 和 abuse/failure behavior，因此需要后续 gate。

## 5. Key And Value Posture

第一版 object identity 是这个 tuple：

```text
owner_kind + owner_id + collection + key
```

未来 implementation 应把 `collection` 和 `key` 当作 protocol-visible identifiers，并进行明确 validation。推荐第一版 posture 是 ASCII-safe、length-bounded、case-sensitive strings，且没有 path semantics。

第一版 value posture 是 bounded JSON object payload candidate：

- payload 必须是 JSON object，不是 arbitrary text 或 binary blob；
- maximum size 必须在 implementation 前 ratify；
- nested depth 和 field count 应在 production hardening 前 bounded；
- value contents 默认 not log-safe；
- object ids、collection names 和 keys 只有在 redaction decision 后才能被视为 log-safe。

如果 schema gate 选择更严格的格式，可以修订第一版 value posture，但它必须保持 small-object game state，不能变成 general file 或 blob storage。

## 6. Read Semantics

未来第一版 behavior 应支持 authenticated player 读取自己的 objects，并且是 route-scoped。

候选 read operations：

- 按 collection 和 key 获取一个 object；
- 在一个 collection 中列出 authenticated player 的 objects，并有 bounded pagination。

Read behavior 应：

- 要求 validated player identity；
- 只通过 protected routes 返回 object metadata 和 value；
- 使用稳定 public not-found behavior；
- 避免泄露另一个 player's object 是否存在；
- 第一版 posture 不得用 client-supplied owner id 覆盖 request identity。

Cross-player reads、public reads、admin reads 和 indexed search 均 deferred。

## 7. Write Semantics

未来第一版 behavior 应支持 authenticated player 对自己的 objects 执行 server-authoritative writes。

候选 write operations：

- 按 collection 和 key put 或 replace object；
- 按 collection 和 key delete object。

Write behavior 应：

- 要求 validated player identity；
- 在 persistence 前验证 collection、key、value shape、value size 和 permission；
- 每次 mutation 成功后发放新 version；
- 在添加 protocol messages 前先用 domain terms 记录 create、update 和 delete behavior；
- 避免隐藏 partial writes；
- 把 mutation logic 保持在未来 storage objects module 或 application boundary 内，而不是放到 transport 或 persistence adapters。

Multi-object transactions、batch writes、partial JSON patch、merge semantics、TTL、public ACL changes 和 server-side script hooks 均 deferred。

## 8. Permission Posture

第一版 permission posture 是 fail-closed 且 route-scoped。

第一版 implementation gate 应先定义 permissions 再写代码。候选 permission families：

- read own storage object；
- list own storage objects in a collection；
- write own storage object；
- delete own storage object。

第一版 posture 不得把 metadata-only `player_id` 当作 proof。除非后续 ADR 明确改变 authentication/session behavior，否则它必须组合现有 authenticated request path 和当前 route-protection model。

Client-controlled read/write permission bits、public object permissions、object ACLs 和 admin bypass 均 deferred。

## 9. Conflict Semantics

第一版 write posture 应使用 optimistic concurrency。

候选 version behavior：

- create 或 update 成功后返回 server-issued version；
- update 和 delete 可以接受 expected version；
- missing expected version behavior 必须在 implementation 前选择；
- stale expected version 返回稳定 conflict class；
- malformed version 返回 validation class；
- not-found、permission failure 和 owner mismatch 必须避免 cross-player existence leaks。

Schema gate 应决定第一版 stored version 是 monotonic integer、database revision、opaque server token、hash-derived token 或其他明确 representation。Behavior gate 只要求 version 由 server-issued 且不是 client-authoritative。

## 10. Protocol Expectations

未来 protocol work 应跟随 semantic behavior 和 contract boundaries。

下面的候选 route families 不由本 gate 实现：

- `runtime.storage.GetObject`
- `runtime.storage.ListObjects`
- `runtime.storage.PutObject`
- `runtime.storage.DeleteObject`

以上 route names 只是 planning candidates。未来 protocol gate 必须在 implementation 前定义确切 module/name routing、Protobuf source files、generated output、request/response shapes、error mapping 和 compatibility posture。

第一版 protocol surface 应保留：

- WebSocket-framed Protobuf envelope use；
- protected-route authentication；
- application-owned dispatch and policy；
- transport credential-neutrality；
- generated-output traceability。

## 11. Data Expectations

下一项 bounded direction 是 storage objects persistence schema gate。

该 gate 应决定：

- table name candidate，可能是 `storage_objects`；
- owner kind 和 owner id representation；
- collection/key constraints；
- value representation；
- version representation；
- create/update/delete timestamps；
- uniqueness 和 indexes；
- logs 和 diagnostics 的 redaction posture；
- migration source candidate；
- future repository 和 PostgreSQL adapter boundaries。

本 behavior gate 不添加 SQL migration source、repository interfaces、PostgreSQL adapters、runtime composition、startup wiring 或 migration apply behavior。

## 12. Verification Expectations

未来 implementation gates 应要求 focused tests 覆盖：

- create then read；
- replace with version change；
- stale version conflict；
- delete then not found；
- own-object permission success；
- cross-owner access rejection without existence leak；
- collection/key validation；
- value size and shape validation；
- redacted errors and logs；
- 在 schema 和 adapter 获得授权后的 PostgreSQL persistence behavior。

默认 repository checks 不应要求 live PostgreSQL，除非该 check 明确通过 disposable database environment opt-in。

## 13. Stop Conditions

做以下任何事情前必须停止并请求 maintainer authorization：

- 实现 storage objects runtime behavior；
- 添加 protocol routes；
- 添加 Protobuf source files 或 generated output；
- 添加 SQL migrations；
- 添加 repository interfaces；
- 添加 PostgreSQL 或其他 storage adapters；
- 添加 dependencies；
- 改变 authentication/session semantics；
- 改变 route protection semantics；
- 添加 cross-player、global、group、party、room、match、public、admin 或 ACL object scopes；
- 添加 large object/blob storage 或 S3-compatible object storage；
- 添加 server-side custom logic hooks；
- 添加 broad operations/admin behavior；
- 添加 hosted deployments 或 demos；
- 创建 release binaries、packages、containers、checksums、signing/provenance artifacts、install scripts、registry publications、SDK packages 或 additional release artifacts；
- 执行 GitHub release record 之外的 public announcements；
- 运行 paid promotion；
- 添加 direct Nakama/Pitaya API compatibility。

## 14. Next Work

下一项 bounded direction 是：

```text
W-0202 Define storage objects persistence schema gate
```

该工作应先决定 storage objects 的第一版 persistence schema posture，然后才能添加 migration source、repository interface、adapter、protocol 或 runtime implementation。
