# Authentication Proof 与 Token Session Contract Dimensions

状态：Draft v0.1
最后更新：2026-05-14
范围：authentication proof、token/session validation、request identity handoff、validation statuses、failure classes、retryability、errors、permissions、commands 和 events 的语义合同维度
依赖：`docs/authentication-token-session-validation.md`

本文件是 `docs/authentication-proof-token-session-contract-dimensions.md` 的简体中文译本。英文版本是权威版本。

## 1. 目的

本标准 ratify 未来 authentication proof 与 token/session validation 工作必须使用的 semantic contract dimensions。

它不选择 login method、token format、token carrier、refresh behavior、signing behavior、credential store、external identity provider、session store、Protobuf envelope change、WebSocket handshake behavior、runtime player handler 或 WebSocket route。

本标准的目标更窄，也更基础：在开始实现前，让未来 Agent 用来描述 authentication 与 session contracts 的词汇和字段稳定下来。

未来工作可以 ratify login、token refresh、logout、session invalidation、connection binding 或 service authentication 等具体 contracts。这些 contracts 必须映射回本文档中的 dimensions，除非后续 ADR 取代它们。

## 2. 必读内容

本标准应与以下文件一起阅读：

- `docs/authentication-token-session-validation.md`
- `docs/player-account-session-contracts.md`
- `docs/player-identity-session-boundary.md`
- `docs/game-protocol.md`
- `docs/runtime-protocol-adapter.md`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `contracts/runtime/session/commands/ValidateSession.yaml`
- `contracts/runtime/session/events/SessionValidated.yaml`
- `contracts/runtime/session/errors/session_errors.yaml`
- `contracts/runtime/session/permissions/session_permissions.yaml`
- `ADR-0019`
- `ADR-0021`
- `ADR-0023`

Nakama 仍然是 account、authentication、session token、refresh token、expiration、revocation 和 realtime socket 能力覆盖的参考。

Pitaya 仍然是 session binding、session context handoff、route handler vocabulary 以及 realtime connection/session 分离的参考。

这两个参考只影响词汇和规划，不支配 vibit 的 API shape。

## 3. 已 Ratify 的维度

每个未来 authentication proof 或 token/session validation contract 都必须说明它拥有、消费或不触碰下列哪些维度。

| Dimension | 含义 | 当前状态 |
| --- | --- | --- |
| `actor_kind` | 已验证或候选 actor 的类型。 | 已 ratify 的词汇。 |
| `actor_id` | validation 后 actor 的稳定标识。 | 已 ratify 的 handoff 字段。 |
| `player_id` | player module 拥有的稳定 player identity。 | validation 前只是 metadata-only。 |
| `session_id` | logical session identifier 或候选 session metadata。 | validation 前只是 metadata-only。 |
| `connection_id` | transport-local connection metadata。 | 不是 proof。 |
| `connection_epoch` | reconnect 或 connection generation metadata。 | 不是 proof。 |
| `validation_status` | domain dispatch 前产生的语义结果分类。 | 已 ratify 的词汇。 |
| `proof_status` | authentication proof 是 absent、present but unverified、proven、rejected、expired、revoked、unsupported 还是 unavailable。 | 供未来 contracts 使用的已 ratify 词汇。 |
| `failure_class` | validation errors 的机器可读失败族。 | 已 ratify 的词汇。 |
| `retryability` | caller 或 runtime 是否可以 retry 同一个语义操作。 | 已 ratify 的期望。 |
| `request_identity_handoff` | 传给 application/domain handlers 的 normalized identity context。 | 由 `runtime/internal/app` 拥有。 |
| `permission_basis` | module permission policies 可依赖的 trust state。 | metadata-only 不足以作为 production permissions 的依据。 |

这些维度是语义层面的，不暗示任何 token shape、storage schema、cryptographic primitive、transport header、envelope field 或 database table。

## 4. Actor Kinds

已 ratify 的 actor kinds：

| Actor kind | 含义 | Production authority |
| --- | --- | --- |
| `unknown` | actor kind 缺失或不可被信任。 | 无 player-owned 或 privileged authority。 |
| `anonymous` | request 被明确作为 unauthenticated 接受。 | 只能访问显式 anonymous capabilities。 |
| `player` | player actor，validation 后由稳定 `player_id` 支撑。 | 只有 validation 成功后才有 player-owned permissions。 |
| `service` | 可信 internal 或 external service actor。 | 只能访问为该 actor 显式建模的 service permissions。 |
| `admin` | administrative actor。 | Deferred；需要未来 admin permission model。 |

当前 Go runtime 在 `runtime/internal/app` 中已有 `unknown`、`player`、`service` 和 `admin`。`anonymous` 被 ratify 为未来 semantic contracts 的 contract vocabulary。把它加入 runtime code 是单独的 implementation change，本设计步骤不要求这样做。

规则：

- `player` 不得从原始 client-supplied `player_id` 推断为 production permissions。
- `service` 和 `admin` 使用前需要单独 ratify proof 与 permission catalogs。
- `anonymous` 必须是显式状态；它不等同于 metadata 缺失。

## 5. Validation Statuses

已 ratify 的 validation statuses：

| Status | 含义 | 当前 production authority |
| --- | --- | --- |
| `unknown` | 没有 validation result。 | 无。 |
| `anonymous` | request 被作为 intentional unauthenticated 接受。 | anonymous-only permissions。 |
| `metadata_only` | identity metadata 在无 authentication proof 的情况下被 normalized。 | 不足以授予 production permissions。 |
| `authentication_proven` | authentication proof 已被验证，但 logical session binding 可能仍是独立步骤。 | Deferred until implementation。 |
| `session_validated` | logical session 已被验证并绑定到 request identity。 | 未来 ratification 后可作为 permission basis。 |
| `service_validated` | service authority 已被验证。 | 未来 ratification 后可作为 service permission basis。 |
| `rejected` | validation 失败，请求不得 dispatch 到 production-sensitive handlers。 | 无。 |

当前 `ValidateSession` semantic contract 使用 `metadata_only` 和 `validated`。本标准为未来词汇做细化，但不要求 runtime code 修改。未来实现时，`validated` 应该在 production permission policies 依赖它之前，被替换或映射到更具体的状态，例如 `authentication_proven`、`session_validated` 或 `service_validated`。

规则：

- `metadata_only` 必须保持 non-authenticated。
- 当 validator 能区分 authentication proof、session binding 或 service authority 时，`validated` 不应作为模糊的长期 production state。
- Domain modules 应消费 normalized request identity，而不是 token/session internals。

## 6. Proof Statuses

供未来 contracts 使用的已 ratify proof statuses：

| Proof status | 含义 |
| --- | --- |
| `not_present` | 没有提供 proof material。 |
| `present_unverified` | proof material 存在，但尚未校验。 |
| `proven` | proof 已由已 ratify 的 authenticator 或 validator 校验通过。 |
| `rejected` | proof 被检查后拒绝。 |
| `expired` | proof 超出允许时间窗口。 |
| `revoked` | proof 被明确 invalidated。 |
| `unsupported` | 当前 runtime 不支持该 proof kind、issuer、method 或 carrier。 |
| `unavailable` | 因 validator 或 dependency 不可用，validation 无法运行。 |

这套词汇有意不选择：

- JWT 或 opaque tokens。
- Access token 或 refresh token 结构。
- Signing、key management、issuer、audience、expiration duration、revocation store 或 replay handling。
- Header、envelope、first-message 或 payload carrier behavior。

## 7. Failure Classes

已 ratify 的 failure classes：

| Failure class | 含义 | 默认 retryability |
| --- | --- | --- |
| `missing_proof` | 必需 proof material 缺失。 | 没有新 proof 时不可 retry。 |
| `malformed_proof` | 已 ratify 的 validator 无法解析 proof material。 | 不改变 input 时不可 retry。 |
| `unsupported_proof` | proof kind、issuer、method 或 carrier 不受支持。 | configuration 或 implementation 改变前不可 retry。 |
| `invalid_proof` | proof 可理解，但被拒绝。 | 同一 proof 不可 retry。 |
| `expired_proof` | proof 形状有效但已过期。 | 只能通过未来已 ratify 的 refresh 或 re-login flow retry。 |
| `revoked_proof` | proof 被明确 invalidated。 | 同一 proof 不可 retry。 |
| `session_not_found` | 引用的 logical session 不存在。 | 通常同一 session 不可 retry。 |
| `session_expired` | logical session 已过期。 | 只能通过未来已 ratify 的 refresh 或 re-login flow retry。 |
| `session_revoked` | logical session 已 invalidated。 | 同一 session 不可 retry。 |
| `actor_disabled` | actor 或 player account 已 disabled。 | 状态改变前不可 retry。 |
| `permission_denied` | validation 成功，但 actor 没有权限。 | authority 改变前不可 retry。 |
| `validator_unavailable` | validator 或其 dependency 不可用。 | 可 retry。 |
| `not_implemented` | 请求的 authentication 或 validation path 尚未实现。 | implementation 改变前不可 retry。 |

未来 error catalogs 必须把每个 validation error 映射到这些 classes 之一；如果需要引入长期存在的新 class，必须通过 change spec 和 ADR 显式说明。

## 8. Retryability Rules

Retryability 是 contract 的一部分，不得由 clients、agents 或 transport adapters 猜测。

规则：

- Error catalogs 必须声明 `retryable: true` 或 `retryable: false`。
- Retryability 描述的是使用同一 proof material 的同一语义请求。
- Refresh、re-login、proof replacement、account recovery 或 support intervention 是独立 flow，不是同一 validation request 的 retry。
- `validator_unavailable` 可以 retry。
- `invalid_proof`、`revoked_proof`、`session_revoked`、`actor_disabled` 和 `not_implemented` 用同一 proof 不可 retry。
- `expired_proof` 和 `session_expired` 只能通过未来已 ratify 的 refresh 或 re-login flow retry。

## 9. Command Dimensions

Authentication proof 与 token/session validation commands 必须声明：

- Actor input semantics。
- 当 proof 在范围内时的 proof input semantics。
- 当 request-level validation 在范围内时的 route 和 target context。
- 作为 metadata only 消费的 transport metadata。
- 该 command 是否可以读取 player accounts、credentials、token state、session state 或 external identity links。
- Output validation status。
- Output request identity handoff fields。
- Failure classes。
- Retryability。
- Required permissions。
- 防止 transport、protocol、player repository 或 domain handlers 拥有 validation behavior 的 invariants。

当前已 ratify command：

```text
ValidateSession
```

`ValidateSession` 由 application 拥有，并且只是 semantic-only。它描述 domain dispatch 前的 request identity handoff。它不实现 token parsing、credential lookup、player account lookup、session persistence、Protobuf envelope changes、WebSocket handshake authentication、runtime player handlers 或 WebSocket routes。

## 10. Query Dimensions

本步骤不 ratify 任何 runtime authentication/session validation query。

未来可以为 read-only inspection、session status、public key metadata 或 service validation metadata ratify queries，但必须先说明 ownership、exposure、permission model 和 information-leakage behavior。

Query 规则：

- Queries 不得泄漏 credential material、token contents、secret keys、provider secrets 或 raw proof material。
- Queries 必须区分 operator/admin inspection 与 gameplay client behavior。
- Queries 必须声明它们对 anonymous、player、service 或 admin actors 是否安全。
- Queries 不得成为 domain modules 内部 token validation 的捷径。

## 11. Event Dimensions

Authentication proof 与 token/session events 必须声明：

- 该 event 是 domain fact、security fact、audit fact 还是 runtime-observation fact。
- 它是否可以发布给 clients。
- 哪些 identifiers 可以安全暴露。
- raw proof material、token strings、credentials、provider subjects 或 secrets 是否被禁止。
- 哪个 command 或 validation path 产生它。
- Compatibility 和 versioning rules。

当前已 ratify event：

```text
SessionValidated
```

`SessionValidated` 只是 semantic fact。它可以描述 decoded request 的 session metadata 被评估成 request identity。它不添加 event bus、durable audit store、token/session store 或 public client event stream。

规则：

- 未来 `AuthenticationSucceeded`、`AuthenticationFailed`、`TokenRefreshed`、`SessionInvalidated` 或类似 event 需要单独 ratification。
- Security-sensitive failure events 在 public exposure 前应先设计 audit/operations 语义。
- Event payloads 不得包含 raw credentials、token strings、password hashes、provider secrets 或完整 third-party identity payloads。

## 12. Error Dimensions

Authentication 和 session errors 必须声明：

- Stable error code。
- Failure class。
- Category。
- Retryability。
- Public-safe message。
- 该 error 是否可以返回给 clients。
- 是否存在 internal-only detail。
- 使用它的 commands 或 queries。

当前 runtime session errors：

- `SESSION_INVALID`
- `SESSION_VALIDATOR_UNAVAILABLE`
- `SESSION_VALIDATION_NOT_IMPLEMENTED`

要求的映射：

| Error code | Failure class | Retryable |
| --- | --- | --- |
| `SESSION_INVALID` | `invalid_proof` 或 `session_not_found`，取决于未来 validator path | `false` |
| `SESSION_VALIDATOR_UNAVAILABLE` | `validator_unavailable` | `true` |
| `SESSION_VALIDATION_NOT_IMPLEMENTED` | `not_implemented` | `false` |

当前 metadata-only runtime 不应仅因为缺少 production proof 就 emit `SESSION_INVALID`，除非未来 validator 已经被 ratify 为要求 proof。

## 13. Permission Dimensions

Authentication 和 session permissions 必须区分：

- 运行 validation infrastructure 的权限。
- 由 validated actor 赋予某 module action 的权限。
- Inspect authentication/session state 的权限。
- Administer 或 revoke sessions 的权限。

当前 runtime session permission：

```text
runtime_session_validate
```

`runtime_session_validate` 允许 application dispatch 在 domain handlers 运行前评估 decoded request metadata。它不是 gameplay permission，也不授予 player-owned domain authority。

规则：

- Domain module permissions 不得把 metadata-only identity 当作 production proof。
- Service 和 admin permissions 需要未来显式 catalogs。
- Token possession 本身不得在没有已 ratify validator 和 request identity handoff 的情况下成为 permission shortcut。

## 14. Request Identity Handoff

Request identity handoff 是 validation 与 domain behavior 之间的边界。

Owner：

```text
runtime/internal/app
```

必需语义字段：

- `validation_status`
- `actor_kind`
- `actor_id`
- `player_id`
- `session_id`
- `connection_id`
- `connection_epoch`
- `session_validated`
- `player_id_validated`
- `reason`

规则：

- Domain modules 在 application validation 后接收 request identity。
- Domain modules 不得直接 parse tokens、credentials、WebSocket headers、Protobuf envelope internals 或 session stores。
- `player_id_validated: true` 要求证明 request actor 被允许以该 `player_id` 行动。
- `session_validated: true` 要求证明 logical session 已按照未来已 ratify 的 session model 被验证并绑定。
- `metadata_only` identity 可用于 local proof slices 和 development behavior，但不是 production permission basis。

## 15. Reference Pattern Map

### Nakama Patterns

| Pattern | Vibit position | Contract dimension impact |
| --- | --- | --- |
| Session token | Deferred | Token/session validation contracts 在实现前必须有 proof status、failure class、retryability 和 request identity handoff dimensions。 |
| Refresh token | Deferred | Refresh 是独立 flow；expired proof 不是同一 validation request 的 retry。 |
| Token expiration | Adopted as dimension | 实现前 expiration 必须映射到 `expired_proof` 或 `session_expired`。 |
| Token revocation/logout | Adopted as dimension | 实现前 revocation 必须映射到 `revoked_proof`、`session_revoked` 和显式 events。 |
| Realtime socket bound to authenticated session | Adapted | vibit 保持 request identity handoff 归 application 所有；handshake binding 仍是后续决策。 |

### Pitaya Patterns

| Pattern | Vibit position | Contract dimension impact |
| --- | --- | --- |
| Session object separate from transport | Adopted | `connection_id` 是 metadata；request identity 归 application 所有。 |
| Handler receives session context | Adapted | Domain handlers 接收 `RequestIdentity`，不是 transport-owned session object。 |
| Session binding | Adopted as vocabulary | Binding 需要显式 `session_validated` 和 `player_id_validated` 语义。 |
| Frontend/backend split | Deferred | Distributed session routing 不属于这个 semantic contract step。 |

## 16. Non-Goals

本标准不做以下事情：

- 选择 JWT、opaque tokens、refresh tokens、signing、expiration duration、revocation store、rotation、key management 或 token storage。
- 选择 guest、device、email/password、custom ID、social login、OAuth、OIDC 或 external identity provider login。
- 添加 credential storage、password hashing、provider secrets、external identity tables、token tables 或 session tables。
- 改变 Protobuf envelope。
- 改变 WebSocket handshake behavior。
- 添加 runtime authentication code。
- 添加 token parsing。
- 添加 runtime player handlers 或 WebSocket routes。
- 添加 event bus publication 或 audit persistence。

## 17. Required Future Artifacts

未来实现任何具体 authentication 或 token/session 行为时，必须按需添加或更新：

- `changes/` 下的 change spec。
- 对长期 architecture-sensitive choices 的 ADR。
- `contracts/` 下的 contract sources。
- `.arch/contracts.yaml`。
- `.arch/runtime.yaml`。
- Module manifests 和 guides。
- Error 和 permission catalogs。
- 当存储状态时的 persistence boundary 和 migrations。
- 当引入 security 或 provider dependencies 时的 dependency adoption records。
- 当 wire shape 改变时的 protocol sources 和 generated output。
- Runtime tests 和 architecture checks。
- 英文文档和简体中文翻译。

## 18. Verification

本标准的默认 verification：

```bash
node tools/vibit check contracts --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check work --json
node tools/vibit check change ratify-authentication-proof-token-session-contract-dimensions --json
node tools/vibit check all --json
```

除非 runtime Go code 发生变化，否则本 design-only standard 不要求 Go tests。

## 19. Agent Rules

Agents 必须：

- 在设计 authentication proof 或 token/session validation contracts 时使用已 ratify dimensions。
- 保持 metadata-only identity 为 non-authenticated。
- 把未来 validation failures 映射到显式 failure classes 和 retryability。
- 通过 `runtime/internal/app` 保持 request identity handoff。
- 记录 Nakama 或 Pitaya 词汇是 adopted、adapted、deferred 还是 rejected。

Agents 不得：

- 把 token strings、player account rows、envelope metadata 或 WebSocket connection IDs 本身当作 authentication proof。
- 把 token validation 隐藏在 transport、protocol adapters、repositories 或 domain handlers 中。
- 仅凭本标准添加 concrete login methods、token behavior、credential storage、session persistence、Protobuf envelope changes 或 WebSocket handshake behavior。
