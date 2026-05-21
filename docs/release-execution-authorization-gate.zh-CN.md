# Release Execution Authorization Gate 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：v0.1 alpha path 的 gate-only release execution authorization criteria
依赖：`docs/release-execution-preparation-gate.md`、`docs/release-publishing-decision-gate.md`、`docs/alpha-developer-flow.md`、`docs/alpha-acceptance-checklist.md`、`docs/runtime-runbook.md`
权威决策：`ADR-0100`
说明：本文件是 `docs/release-execution-authorization-gate.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档定义 release execution authorization gate。它不发布 `v0.1 alpha`，不选择或创建 release identifiers 或 tags，不创建 binaries、archives、containers、packages、checksums、provenance files、hosted deployments、runtime behavior、protocol routes、Protobuf sources、generated output、migrations、dependencies、broad operations/admin behavior、authentication/session behavior changes、broad product modules 或 direct Nakama/Pitaya API compatibility。

## 1. 目的

Release execution preparation gate 已定义 future release execution step 的 inputs 和 planning boundaries。剩余风险是把 readiness criteria 误认为执行 release 的许可。

本 gate 定义 maintainers 在后续做 go/no-go decision 前必须满足的 authorization criteria。它是 authorization criteria document，不是 release execution record。

该 gate 记录 final go/no-go criteria、required verification state、release identifier review requirements、artifact authorization boundaries、maintainer approval requirements 和 stop conditions。

## 2. Core Rule

Release execution authorization gate 是：

```yaml
release_execution_authorization_gate: defined
completed_work_item: W-0192
decision: ADR-0100
check_rule: runtime.release_execution_authorization_gate
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
final_go_no_go_criteria_defined: true
required_verification_state_defined: true
release_identifier_review_defined: true
artifact_authorization_boundaries_defined: true
maintainer_approval_requirements_defined: true
stop_conditions_defined: true
release_identifier_selected: false
release_tag_created: false
release_binary_created: false
release_archive_created: false
release_container_created: false
release_package_created: false
release_checksum_created: false
release_provenance_created: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
direct_nakama_pitaya_api_compatibility_added: false
next_direction: release_execution_maintainer_decision
```

该 gate 让 maintainers 后续可以 review 是否授权 release execution。它本身不授权 release execution，不选择 release identifier，不创建 artifacts，也不发布 release。

## 3. Final Go/No-Go Criteria

只有在以下条件全部满足时，后续 release execution maintainer decision 才可以被考虑：

- `docs/v0.1-alpha-goal.md` 仍准确描述预期 `v0.1 alpha` scope。
- `docs/alpha-acceptance-checklist.md` 对 local alpha flow 不记录 unresolved `Blocked` item。
- `docs/alpha-developer-flow.md` 仍能把 contributors 引导到 coherent local path。
- `docs/runtime-runbook.md` 仍匹配当前 runtime startup、PostgreSQL setup、verifier key 和 redaction posture。
- `docs/release-publishing-decision-gate.md` 和 `docs/release-execution-preparation-gate.md` 仍满足。
- `examples/local-alpha-request-loop.sh` 仍运行 redacted authenticated gameplay proof。
- Repository checks 和 Go tests 通过，且没有新的未 triage warning。
- 如果 known warning 仍存在，它仍是已记录的 `runtime.identity_boundary` warning，对应 `runtime/internal/platform/persistence/postgres/authentication_repository.go`。
- 面向公众的 release-path documents 仍保持 English canonical documents 与 Simplified Chinese translations 成对。
- 任何 tracked artifact 都不包含 raw credentials、raw access tokens、verifier keys、lookup digests、verifier digests、HMAC input/output bytes、带 credentials 的 DSNs、transport proof carriers、GitHub tokens 或 private environment file content。
- Maintainer 明确选择做 release execution go/no-go decision。

如果任一条件不满足，release execution decision 必须是 no-go，直到问题被解决或由后续 repository record 明确 deferred。

## 4. Required Verification State

后续 go/no-go decision 的默认 required verification state 是：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```

后续 release decision record 应记录 command results、known warnings，以及任何 intentionally skipped optional verification。

Optional live PostgreSQL checks 仍通过 `VIBIT_POSTGRES_TEST_DSN` 和 disposable database opt-in。本 gate 不要求默认 repository checks 必须执行 live PostgreSQL verification。

## 5. Release Identifier Review

本 gate 不选择 release identifier。

后续 release execution decision 必须 review：

- target identifier 是 `v0.1 alpha`，还是另一个被明确选择的 pre-release identifier；
- 该 identifier 是否已作为 Git tag、release record、package、image 或 external artifact 存在；
- 该 identifier 是否匹配 project naming 和 public communication expectations；
- 该 identifier 是否暗示超出当前 pre-alpha/alpha posture 的稳定性级别；
- 如果后续授权 tags、archives、checksums、packages 或 hosted deployments，该 identifier 是否能安全表达在 release notes、tags、archives、checksums、packages 或 hosted deployments 中。

Final identifier selection 仍是 ask-first boundary，必须在后续 bounded work item 中记录。

## 6. Artifact Authorization Boundary

本 gate 不授权任何 artifacts。

后续 maintainer decision 必须明确说明哪些 artifact families 被授权：

- Git version tag。
- GitHub release 或等价 release record。
- Release notes 或 changelog entry。
- Source archive。
- Checksum file。
- Provenance file。
- Binary build。
- Package。
- Container image。
- Hosted deployment。

Artifact authorization 应优先选择最小 surface，让 developers 能 inspect source 并运行 local alpha path。Optional binaries、packages、containers、provenance files 和 hosted deployments 应继续 deferred，除非后续 work item 明确批准。

## 7. Maintainer Approval Requirements

后续 release execution decision 需要 durable maintainer approval record，并回答：

- Release execution 是 go 还是 no-go？
- 哪个 release identifier 被批准（如有）？
- 哪些 artifact families 被批准（如有）？
- Review 了哪些 verification results？
- 哪些 known warnings 被接受或拒绝？
- 哪些 deferrals 仍然有效？
- 哪些 commands 被允许创建已批准 artifacts？
- Execution 过程中哪些 stop conditions 仍适用？

仅聊天中的 approval 不足够。任何 release execution command 运行前，approval 必须由 repository artifact 和 bounded work item 表达。

## 8. Authorization Outcome

本 gate 的 outcome 是：

```yaml
authorization_criteria_defined: true
may_make_maintainer_go_no_go_decision_later: true
may_publish_release_now: false
may_execute_release_now: false
may_select_release_identifier_now: false
may_create_release_artifacts_now: false
```

Repository 仍是 pre-alpha。下一步是 maintainer decision gate，而不是 release execution。

## 9. Stop Conditions

如果发生以下任何情况，在后续 release execution 前停止：

- required verification 失败；
- 出现新的 warning 且没有明确 triage；
- known warning state 变化但没有 documented decision；
- 面向公众的 English 和 Simplified Chinese documents 出现实质 divergence；
- 任何 tracked artifact 包含 secrets 或 unredacted proof material；
- generated output 改变但没有 source trace 和 generation notes；
- runtime behavior、protocol routes、Protobuf sources、migrations、dependencies、operations/admin behavior、authentication/session behavior、broad product modules 或 direct Nakama/Pitaya compatibility 在 authorization-only slice 中变化；
- release identifier 未经 maintainer approval 被选择；
- 创建 release tag、artifact、package、checksum、provenance file、hosted deployment 或 release record；
- release scope 暗示超出 alpha goal 的 production readiness、SLA、hosted service availability 或 compatibility promises；
- maintainer approval 没有记录在 repository artifact 中。

## 10. Redaction Expectations

Authorization records 必须保持可以安全提交到 repository。

不要包含：

- raw device credential text 或 bytes；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- HMAC input 或 output bytes；
- verifier key values；
- concrete verifier key set ids；
- 带 credentials 的 PostgreSQL DSNs；
- 可能携带 secrets 的 headers、cookies、query strings、WebSocket subprotocol values 或 remote addresses；
- 包含 GitHub tokens 或其他 access credentials 的 local files。

## 11. Reference Alignment

Nakama 和 Pitaya 都体现了严肃 server framework 应分离 release readiness、release decisions、artifacts 和 deployment posture。本 gate 采用这种纪律，但不采用它们的 APIs、data models、route names、release packaging、deployment model、cluster model、SDK surfaces 或 operations surfaces。

## 12. Next Work

`W-0193 Confirm release execution maintainer decision` 现在已经在 `docs/release-execution-maintainer-decision.md` 中记录 maintainer go decision。下一步 bounded work 是：

```text
W-0194 Define release identifier and artifact plan
```

该 future planning step 仍不得创建 release tags、artifacts、hosted deployments 或 published release records，除非后续 execution scope 明确授权这些 actions，并且 maintainer 批准 ask-first boundary。
