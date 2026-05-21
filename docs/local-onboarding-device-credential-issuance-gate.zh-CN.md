# Local Onboarding Device Credential Issuance Gate 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：为第一版本地开发者 onboarding flow 定义 gate-only boundary：创建 player account，并签发 server-generated device credential
依赖：`docs/v0.1-alpha-goal.md`、`docs/device-credential-login-service-behavior-gate.md`、`docs/token-credential-material-generation-implementation-gate.md`、`docs/verifier-digest-helper-implementation-gate.md`、`docs/verifier-digest-comparison-helper-gate.md`、`docs/credential-record-schema-boundary.md`、`docs/postgresql-persistence-boundary.md`
Canonical decision：`ADR-0089`

说明：本文件是 `docs/local-onboarding-device-credential-issuance-gate.md` 的简体中文译本。英文版本是权威版本。

## 1. 目的

Alpha path 现在已经有 login、access-token validation、runtime session metadata、first-message connection binding、protected inventory routes、presence lifecycle primitive、protected presence query 和 logout。新的本地开发者仍然缺少进入这条路径所需的第一份 credential。

本 gate 在实现任何 onboarding 行为之前，先定义未来 local onboarding/device credential issuance boundary。

这是 gate-only standard。它不实现 onboarding，不生成或展示 credentials，不通过新流程创建 player accounts，不通过新流程写 credential records，不暴露 public protocol route，不改变 Protobuf sources，不改变 generated output，不改变 migrations，不添加 dependencies，不发布 release，不添加 production signup，不添加 external identity providers，不添加 password login，不添加 account recovery，不添加 multi-device linking，也不采用 direct Nakama/Pitaya API compatibility。

## 2. 核心规则

local onboarding device credential issuance gate 是：

```yaml
local_onboarding_device_credential_issuance_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0181
future_implementation_work_item: W-0182
decision: ADR-0089
check_rule: runtime.local_onboarding_device_credential_issuance_gate
first_surface_candidate: local_developer_onboarding_application_service
surface_visibility: local_only_not_public_signup
future_service_owner: runtime/internal/app/authentication
future_source: runtime/internal/app/authentication/service.go
future_tests: runtime/internal/app/authentication/service_test.go
future_service_method_candidate: OnboardLocalPlayerWithDeviceCredential
player_account_repository_method: CreatePlayerAccount
authentication_repository_method: StoreCredential
credential_material_helper: GenerateDeviceCredentialMaterial
credential_lookup_digest_helper: ComputeCredentialLookupDigest
credential_verifier_digest_helper: ComputeCredentialVerifierDigest
login_route_account_creation_behavior_changed: false
public_protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
production_signup_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

未来 implementation 只有在后续 work item 明确授权 code 之后，才能添加本地开发者 onboarding application service。它不得把现有 public `runtime.authentication.AuthenticateWithDeviceCredential` login route 变成 account creation 行为。

## 3. 未来 Surface

第一版未来 surface candidate 是 application-local：

```yaml
local_onboarding_surface:
  owner: runtime/internal/app/authentication
  visibility: local_developer_only
  public_protocol_route: forbidden_by_this_gate
  websocket_route: forbidden_by_this_gate
  http_route: forbidden_by_this_gate
  startup_auto_onboarding: forbidden_by_this_gate
  production_signup: forbidden_by_this_gate
```

规则：

- 未来方法可以作为 application service method 存在，供 tests 和后续 local tooling 调用。
- 它不是 public game client signup route。
- 除非后续 protocol gate 明确添加 route，否则它不得通过 WebSocket Protobuf routing 访问。
- 它不得在 process startup 时自动创建 player accounts。
- 它不得把现有 login route 上的 `AccountCreationIntent` 解释为创建账号的授权。
- 第一版 posture 不接受 client-generated device credential material。

local-only posture 是有意选择。它让 alpha 开发者能用受控方式获得第一份 credential，同时不把 vibit 提前绑定到 production signup、abuse controls、identity provider linking、recovery、account merge 或 multi-device behavior。

## 4. 未来依赖形状

未来 implementation 只应按 local onboarding 需要扩展 application service dependencies。

```yaml
future_service_dependencies:
  unit_of_work_runner: already_present
  verifier_key_set: already_present
  device_credential_entropy_reader: required
  clock: already_present
  player_id_generator: required
  player_account_event_id_generator: required
  credential_record_id_generator: required
```

规则：

- application unit-of-work 仍然是唯一 transaction entry point。
- service 必须从 unit-of-work capability 获取 `NewPlayerAccountRepository()` 和 `NewAuthenticationRepository()`。
- 不得仅为这个 slice 扩展全局 `tx.UnitOfWork` interface。
- Identifier generation 必须注入。本 gate 不选择 UUID、ULID、KSUID、database-generated ids 或外部 id dependency。
- Device credential entropy reader 必须显式传入，避免 tests 依赖不可控 process state。
- 本 gate 不 ratify production default id generation、display names、local operator identity 或 credential lifetime。

## 5. 未来请求和结果形状

候选请求：

```yaml
LocalOnboardingDeviceCredentialIssuanceRequest:
  display_name: required_non_secret_text
  requested_by: optional_local_operator_label
```

候选结果：

```yaml
LocalOnboardingDeviceCredentialIssuanceResult:
  status: created
  player_id: generated_player_id
  credential_record_id: generated_credential_record_id
  device_credential: one_time_raw_credential_text
  created_at: server_time
```

规则：

- 第一版 posture 应由服务器生成 `player_id`、player account event id 和 credential record id。
- 第一版 posture 不应允许 caller-supplied `player_id` 作为 identity proof。
- Display name 不是 proof，且不得嵌入 credential material 或 digests。
- Raw device credential text 只能在 unit of work commit 成功之后出现在成功结果中。
- Raw device credential bytes 和 text 永远不得存储。
- Result 不得包含 access token。Login 仍然通过现有 device credential login route 单独执行。

## 6. 未来必需顺序

当 `W-0182` 或后续 work item 授权 behavior 时，未来 onboarding method 必须按以下顺序执行：

```yaml
local_onboarding_sequence:
  - reject_invalid_local_request_before_unit_of_work
  - generate_device_credential_material_with_explicit_entropy_reader
  - compute_credential_lookup_digest_with_active_VerifierKeySet
  - compute_credential_verifier_digest_with_active_VerifierKeySet
  - enter_application_unit_of_work
  - obtain_player_account_repository_from_unit_of_work_capability
  - obtain_authentication_repository_from_unit_of_work_capability
  - generate_player_account_identifiers_with_injected_generators
  - create_active_player_account
  - create_active_device_credential_record_with_digest_only_storage
  - exit_unit_of_work_successfully
  - return_raw_device_credential_text_once_after_commit
```

规则：

- 基础 request validation 之前不得发生 repository call。
- Credential material generation 必须使用 `GenerateDeviceCredentialMaterial`。
- Credential digest computation 必须使用 `ComputeCredentialLookupDigest` 和 `ComputeCredentialVerifierDigest`。
- Credential record 必须使用 `credential_kind=device_credential_login`。
- Credential record 必须使用 `verifier_algorithm=vibit_hmac_sha256_v1` 和 `verifier_version=1`。
- `verifier_key_id` 必须来自 active `VerifierKeySet.KeySetID()`。
- Player account creation 和 credential record storage 必须一起 commit 或 rollback。
- 如果 player account creation、credential record storage、dependency lookup、id generation、digest computation 或 unit-of-work commit 失败，该方法不得把 raw credential text 作为成功结果返回。

## 7. Repository Handoff

未来 implementation 必须使用现有 repository interfaces。

```yaml
repository_handoff:
  transaction_boundary: UnitOfWorkRunner.WithinUnitOfWork
  player_repository_source: unit_of_work.NewPlayerAccountRepository
  authentication_repository_source: unit_of_work.NewAuthenticationRepository
  player_mutation_method: CreatePlayerAccount
  credential_store_method: StoreCredential
  direct_postgres_import: forbidden
  repository_interface_change: forbidden_by_this_gate
```

Authentication repository 只能接收 digest 和 metadata fields：

```yaml
credential_store_allowed_fields:
  - credential_record_id
  - player_id
  - credential_kind
  - credential_lookup_digest
  - credential_verifier_digest
  - verifier_algorithm
  - verifier_version
  - verifier_key_id
  - occurred_at
  - requested_by
```

禁止 repository inputs：

- Raw device credential text。
- Raw device credential bytes。
- Encoded credential material。
- Verifier key bytes。
- Encoded verifier key values。
- Access-token material。
- Provider subjects、passwords、OAuth claims、OIDC claims 或 account recovery data。

## 8. 与现有 Login 的关系

本 gate 不改变 device credential login。

```yaml
existing_login_route: runtime.authentication.AuthenticateWithDeviceCredential
existing_login_route_account_creation_behavior_changed: false
existing_login_route_proof_required: true
login_route_returns_access_token: already_implemented
local_onboarding_returns_access_token: false
```

规则：

- `AuthenticateWithDeviceCredential` 仍然要求 presented device credential proof。
- 现有 login request 上的 `AccountCreationIntent` 在后续 route behavior decision 明确改变之前仍不创建账号。
- Local onboarding 创建 credential。Login 使用该 credential 认证。
- Local onboarding 不得通过直接签发 access token 来绕过未来 login proof validation。

## 9. Redaction 要求

未来 implementation 不得把以下值放入 errors、logs、traces、metrics labels、docs examples、test snapshots、ADRs、change specs、conversation logs 或 public responses。唯一例外是成功 local onboarding result 的 one-time carrier：

- Raw device credential text。
- Raw device credential bytes。
- Encoded generated credential material。
- Credential lookup digest bytes。
- Credential verifier digest bytes。
- HMAC input 或 output bytes。
- Verifier key bytes。
- Encoded verifier key values。
- Full concrete `verifier_key_id` values。
- Database connection strings 或 credentials。
- Credential lookup hit/miss details。
- 会泄露 private state 的 player account conflict details。

允许：

- 成功结果中的非 secret `player_id` 和 `credential_record_id`。
- Registered 或 local-only redacted error classes。
- Redacted placeholders，例如 `<device-credential>`、`<credential-lookup-digest>` 和 `<verifier-key-id>`。

生成出来的 credential material 不会因为是 local 就安全。One-time presentation 表示一个成功 local result carrier，不是“可以安全记录一次”。

## 10. 未来 Tests

未来 implementation 必须在以下文件添加 focused tests：

```text
runtime/internal/app/authentication/service_test.go
```

最低 test classes：

```yaml
required_tests:
  onboarding_rejects_invalid_request_without_unit_of_work
  onboarding_generates_device_credential_with_explicit_reader
  onboarding_computes_credential_lookup_and_verifier_digests_before_storage
  onboarding_uses_player_repository_from_unit_of_work_only
  onboarding_uses_authentication_repository_from_unit_of_work_only
  onboarding_creates_active_player_account_before_credential_record
  onboarding_stores_credential_digests_only
  onboarding_returns_raw_device_credential_only_after_commit
  onboarding_does_not_return_credential_when_player_creation_fails
  onboarding_does_not_return_credential_when_credential_storage_fails
  onboarding_does_not_return_credential_when_commit_fails
  onboarding_does_not_issue_access_token_or_runtime_session
  onboarding_errors_do_not_leak_raw_credential_digest_or_key_material
  existing_login_route_still_does_not_create_accounts
```

默认不应要求 live PostgreSQL。除非后续 work item 明确要求 live verification，否则 repository behavior 可以用 fakes 或现有 adapter tests 覆盖。

## 11. Deferrals

本 gate 保留以下 deferrals：

- Runtime onboarding implementation。
- Public signup 或 production account creation。
- Public WebSocket、HTTP 或 CLI surface selection。
- Protobuf request/response messages。
- Generated Go output。
- Migration schema changes。
- Repository interface changes。
- New dependencies。
- External identity providers、OAuth、OIDC、provider SDKs、password login、password hashing、account recovery、account merge 和 multi-device linking。
- Credential rotation 或 replacement behavior。
- Onboarding 直接签发 access-token。
- Onboarding 直接创建 runtime session。
- Memory-store durable authentication behavior。
- Release publishing。
- Direct Nakama/Pitaya public API compatibility。

## 12. Nakama 和 Pitaya 映射

Nakama 表明，可用的 game backend 需要一条创建或认证 player、并获得 secret session/token material 的路径，然后才能发 gameplay requests。vibit 采纳这个 capability need，但把 first local credential issuance 与 production signup 和 direct Nakama API compatibility 分开。

Pitaya 强调 transport acceptors、route handlers 和 identity context 应分离。vibit 通过把 local onboarding 放在 application service orchestration 中，而不是放进 WebSocket transport 或 Protobuf envelope metadata 中，来吸收这个经验。

## 13. 验证

本 gate 的 repository check rule 是：

```text
runtime.local_onboarding_device_credential_issuance_gate
```

推荐验证命令：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-local-onboarding-device-credential-issuance-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Runtime Go tests and Buf generation are not required for this gate-only change because it does not add Go runtime behavior, Protobuf sources, generated output, migrations, dependencies, or release artifacts.
