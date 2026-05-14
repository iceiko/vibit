# Token Format And Carrier Option Comparison 中文版

状态：Draft v0.1
最后更新：2026-05-14
范围：M-013 的 token format 和 proof carrier comparison
依赖：`docs/first-login-method-set.md`

本文件是 `docs/token-format-carrier-options.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

## 1. 目的

本文在 `device_credential_login` 被 ratify 为第一版 login-method set 后，比较 token format 和 proof carrier options。

它为 W-0068 推荐第一版 token format 和 proof carrier posture。

它不实现 token parsing、signing、validation、refresh、revocation、rotation、replay handling、token storage、session persistence，不修改 Protobuf envelope，不修改 WebSocket handshake authentication，不添加 runtime player handlers，也不添加 WebSocket routes。

## 2. Evaluation Criteria

每个 option 按以下维度评估：

- Production safety。
- Revocation and logout ergonomics。
- Redaction safety。
- Agent-readable implementation shape。
- Dependency load。
- Storage and migration implications。
- Fit with `device_credential_login`。
- Fit with WebSocket-first gameplay。
- Nakama session capability alignment。
- Pitaya session/context vocabulary alignment。
- Reversibility。

## 3. Token Format Summary

| Format | Recommendation | Reason |
| --- | --- | --- |
| `opaque_high_entropy_token` | 推荐作为第一版 token format。 | 它让 validation 保持 explicit 且 storage-backed，容易 redaction，清晰支持 revocation/logout，避免 key-management 和 JWT dependencies，并给 agents 一个窄的 lookup-based implementation path。 |
| `signed_structured_token` | Deferred。 | 后续有用，但 signing、key rotation、issuer/audience/claim drift、revocation 和 replay controls 对第一版 slice 过宽。 |
| `external_provider_token` | Deferred。 | 它属于未来 provider login 和 external identity linking，不属于第一版 device credential path。 |
| `plain_session_id_as_secret` | 第一版 token format 中 rejected。 | 它混淆 identifier 和 proof vocabulary，并有把现有 metadata-only `session_id` 变成 authority 的风险。 |

## 4. Proof Carrier Summary

| Carrier | Recommendation | Reason |
| --- | --- | --- |
| `login_command_response_token` | 推荐用于 successful login 后签发第一版 token。 | 它把 credential exchange 和 token issuance 保持在 application/protocol command flow 中，而不是放进 WebSocket handshake 或 envelope metadata。 |
| `explicit_request_proof_payload` | 推荐作为 authenticated routes 的第一版 request proof posture。 | 它让 proof 在 semantic contracts 中可见，同时不修改当前 Protobuf envelope。它比较 verbose，但 agent-readable 且 reversible。 |
| `first_system_message_binding` | Deferred。 | 未来对 connection binding 有用，但需要 system-message contracts、timeout behavior、reconnect rules 和 connection state。 |
| `protobuf_envelope_extension` | Deferred。 | 后续可能更简洁，但它是 compatibility-sensitive protocol change。 |
| 当前 `Session.session_id` metadata | 作为 proof rejected。 | 它今天是 metadata-only，不能通过 reinterpretation 变成 proof。 |
| WebSocket handshake carrier | Deferred。 | 它可以早拒绝 unauthenticated connections，但有把 authentication 放进 transport 的风险，并需要单独分析 browser/non-browser carrier。 |
| WebSocket subprotocol、cookie 或 query parameter | Deferred。 | 每一种都需要显式 transport risk analysis，不能因为方便就选择。 |

## 5. Recommended First Posture

推荐的第一版 token format：

```text
opaque_high_entropy_token
```

推荐的第一版 proof carrier posture：

```yaml
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
protobuf_envelope_change: false
websocket_handshake_authentication_change: false
current_session_metadata_as_proof: false
first_system_message_binding: deferred
```

这个建议意味着：

- 未来 successful `device_credential_login` command 可以签发 opaque high-entropy access token。
- Raw token value 只能由 clients present，并且其他地方必须 redacted。
- Server-side storage 如被使用，必须存 lookup-safe hash 或等价的 non-plaintext verifier，不存 raw token strings。
- Authenticated gameplay requests 应通过未来 explicit semantic contract fields 携带 proof，直到后续 protocol decision ratify 更清晰的 carrier。
- 现有 Protobuf `Session.session_id`、`Session.player_id`、`Session.connection_id` 和 `Session.connection_epoch` 保持 metadata-only。
- WebSocket transport 保持 credential-neutral。

W-0068 必须用 issuer、verifier、subject、audience、expiration、refresh、revocation、rotation、replay、redaction 和 storage implications 来 ratify 或调整本建议。

## 6. Token Format Details

### Opaque High-Entropy Token

Position:

```text
recommend_first
```

Benefits:

- 第一版 slice 不需要 signing 或 key-management dependency。
- Server-side validation 可以 explicit 且 contract-checkable。
- 如果存在 token records，revocation、logout 和 forced invalidation 很直接。
- 因为 token 没有 client-inspectable claims，所以 redaction 简单。
- 适合 `device_credential_login`：login 证明 credential possession；token 证明后续 session 或 access grant。
- 把 claim evolution 留在 client-visible token contents 之外。

Risks:

- 需要 token lookup storage 或等价 verifier。
- 需要谨慎的 token hashing 和 indexing rules。
- 需要 cleanup 和 expiration logic。
- Implementation 前会增加 database 或 session-store work。
- 每次 validation 可能命中 storage，除非后续 ratify caching 或 session binding。

实现前 required artifacts：

- Token issuance contract。
- Token validation contract。
- Token storage 或 verifier schema gate。
- Hash and lookup rule。
- Expiration rule。
- Revocation/logout rule。
- Redaction rule。
- Replay and collision tests。
- 防止 plaintext token storage 和 metadata-only proof shortcuts 的 repository checks。

推荐状态：

```yaml
format: opaque_high_entropy_token
recommended_for_first_posture: true
requires_signing_dependency: false
requires_key_management: false
requires_token_storage_or_verifier: true
revocation_fit: high
redaction_fit: high
agent_clarity: high
confidence: high
```

### Signed Structured Token

Position:

```text
defer
```

Benefits:

- 可不通过 storage lookup 验证。
- 可携带 issuer、audience、subject、expiration 和 claims。
- 许多 backend teams 熟悉。

Risks:

- 需要 signing dependency 和 key-management posture。
- 没有 server-side denylist 或 session store 时，revocation 更难。
- Claims 可能与 server truth 漂移。
- Agents 可能倾向把 authorization facts 放进 token claims，而不是 module-owned permission logic。
- Key rotation、algorithm agility、clock skew、audience validation 和 replay posture 必须显式。

Recommendation:

```text
Defer until key management, revocation, and claim ownership are ratified.
```

### External Provider Token

Position:

```text
defer
```

Benefits:

- 对未来 platform、social、OAuth 或 OIDC login 有用。
- 当 provider validation 被 ratify 后，可复用 provider-issued proof。

Risks:

- 它不是 vibit 自己的 session token。
- Provider issuer、audience、key、metadata、outage 和 refresh semantics 不同。
- 需要 external identity linking 和 provider dependency decisions。

Recommendation:

```text
Defer until external provider login and external identity linking are ratified.
```

### Plain Session ID As Secret

Position:

```text
reject_for_first
```

Benefits:

- Vocabulary 简单。
- 看起来类似常见 session-cookie systems。

Risks:

- 容易混淆 identifier 和 proof。
- 与当前 metadata-only `Session.session_id` 冲突。
- 鼓励后续 agents 把现有 envelope fields 视为 authenticated。
- 不如显式命名 opaque access token 和后续 logical session 清楚。

Recommendation:

```text
Do not use plain session ID as the first token format.
```

未来 runtime session 可以拥有 high-entropy session identifier，但这必须作为 session persistence 被 ratify，而不是 reinterpret 当前 metadata。

## 7. Carrier Details

### Login Command Response Token

Position:

```text
recommend_for_issuance
```

第一版 token 应由未来 semantic login command response 在 credential validation 成功后签发。

Benefits:

- 保持 transport credential-neutral。
- 避免 Protobuf envelope changes。
- 让 authentication result 绑定到显式 command contracts。
- 便于 agents 作为 request/response behavior 测试。

Risks:

- 需要未来 login command 和 response contracts。
- 需要围绕 response payloads 的 redaction 和 logging rules。
- 它本身不绑定 WebSocket connection。

### Explicit Request Proof Payload

Position:

```text
recommend_for_first_authenticated_routes
```

Authenticated commands 或 queries 应在 explicit contract-owned payload fields 中携带 token proof，直到后续 carrier 被 ratify。

Benefits:

- 不需要 envelope 或 handshake change。
- 在 semantic contracts 中清晰。
- 可跨 transports 工作。
- 避免把 metadata-only fields 当成 authority。
- 如果后续 ratify envelope 或 session-binding carrier，可逆。

Risks:

- 比 envelope-level 或 connection-bound proof 更 verbose。
- 需要每个 authenticated contract 声明 proof semantics，或后续共享 generated proof wrapper conventions。
- 如果 redaction 未强制，可能增加 token exposure。

### First System Message Binding

Position:

```text
defer
```

后续它可能成为强 WebSocket-first model，尤其适合 realtime gameplay。它需要 system-message contracts、binding state、timeout behavior、reconnect rules 和 connection lifecycle tests。

### Protobuf Envelope Extension

Position:

```text
defer
```

Envelope-level proof field 后续可能减少 per-contract repetition，但它会改变 public wire schema 和 generated output。它需要 protocol change spec 和 compatibility review。

### Current Session Metadata

Position:

```text
reject_as_proof
```

当前 `Session.session_id`、`Session.player_id`、`Session.connection_id` 和 `Session.connection_epoch` fields 保持 metadata-only。它们只可在允许的地方作为 metadata copy、normalize 或 log。它们不得授权 player-owned behavior。

### WebSocket Handshake Carrier

Position:

```text
defer
```

Handshake proof 后续对 early rejection 可能有用，但它触及 transport/process boundary，必须通过专门 handshake decision ratify。不能在比较 token formats 时添加。

## 8. Comparative Matrix

Scores 是 qualitative：

```text
high = favorable
medium = manageable with gates
low = unfavorable for first slice
```

| Token format | Production safety | Revocation | Redaction | Dependency load | Storage complexity | Agent clarity | First-slice fit |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `opaque_high_entropy_token` | High | High | High | High | Medium | High | High |
| `signed_structured_token` | Medium | Low | Medium | Low | Low to medium | Medium | Medium |
| `external_provider_token` | Medium | Medium | Low | Low | Medium | Low | Low |
| `plain_session_id_as_secret` | Low | Medium | Medium | High | Medium | Low | Low |

| Carrier | Boundary clarity | Protocol compatibility | WebSocket fit | Agent clarity | First-slice fit |
| --- | --- | --- | --- | --- | --- |
| `login_command_response_token` | High | High | High | High | High |
| `explicit_request_proof_payload` | High | High | Medium | High | High |
| `first_system_message_binding` | Medium | Medium | High | Medium | Medium |
| `protobuf_envelope_extension` | Medium | Low | High | Medium | Medium |
| Current `Session.session_id` metadata | Low | High | Medium | Low | Low |
| WebSocket handshake carrier | Medium | Medium | High | Medium | Medium |

## 9. Reference Mapping

### Nakama

Nakama 通常在 authentication 后返回 session tokens，并支持基于 refresh-token 的 continuation。

vibit 将其改造为：

- Adopted concept：authentication 产生 server-recognized token 或 session proof。
- Recommended first token format：opaque high-entropy token，不是 Nakama-compatible token format。
- Deferred concept：refresh token lifecycle。
- Deferred concept：realtime socket bound to authenticated session。
- 当前 rejected：direct Nakama API compatibility。

### Pitaya

Pitaya 的相关输入是 session context 和 handler binding。

vibit 将其改造为：

- 未来 handlers 在 validation 后接收 normalized request identity。
- Proof carrier 必须在 domain dispatch 前产生 `RequestIdentity`。
- Connection binding 和 session object behavior 继续 deferred。
- Authentication 不得隐藏在 transport acceptors 或 route handlers 中。

## 10. Recommended W-0068 Ratification Packet

W-0068 应 ratify 或调整：

```yaml
first_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
access_token: selected
refresh_token: deferred
session_token_vocabulary: deferred_until_session_persistence
protobuf_envelope_change: false
websocket_handshake_authentication_change: false
current_session_metadata_as_proof: false
requires_before_implementation:
  - token_issuance_contract
  - token_validation_contract
  - token_hash_lookup_rule
  - token_storage_or_verifier_schema_gate
  - expiration_rule
  - revocation_logout_rule
  - redaction_rules
  - repository_checks
  - focused_tests
does_not_authorize:
  - token_parser_code
  - token_tables
  - session_tables
  - credential_tables
  - jwt_or_signing_dependency
  - protobuf_envelope_change
  - websocket_handshake_authentication
  - runtime_player_handlers
  - websocket_routes
```

## 11. Open Questions For Later Work

- 第一版 opaque token 在 public contracts 中叫 access token、session token，还是另一个 vibit term。
- Token storage 是立即使用 PostgreSQL，还是后续 session store。
- 第一版 production implementation 是否包含 refresh tokens。
- Expiration 是否足够短以避免第一版 implementation 中的 refresh。
- Revocation 和 logout 是否必须先于第一版 runtime login implementation。
- Request proof payload fields 是在每个 authenticated command 中重复，还是后续通过 shared contract wrapper 生成。

这些问题不是 W-0067 的 blockers。它们应由 W-0068 到 W-0071 回答。
