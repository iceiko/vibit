# Release Execution Maintainer Decision

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Maintainer go/no-go decision record for the v0.1 alpha release execution path
Depends on: `docs/release-execution-authorization-gate.md`
Canonical decision: `ADR-0101`

The paired Simplified Chinese translation is `docs/release-execution-maintainer-decision.zh-CN.md`. The English file is authoritative.

This document records the maintainer decision for `W-0193`. It does not publish `v0.1 alpha`, choose or create a final release identifier or tag, create binaries, archives, containers, packages, checksums, provenance files, hosted deployments, GitHub release records, runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, broad product modules, or direct Nakama/Pitaya API compatibility.

## 1. Purpose

`W-0192` intentionally stopped at a blocked maintainer decision gate. The maintainer has now clarified that the decision outcome is approval to continue the release execution path.

The approval is intentionally narrow. It authorizes the repository to move from the blocked go/no-go gate into a bounded release identifier and artifact planning step. It does not authorize immediate release execution commands.

## 2. Decision Record

The maintainer decision is:

```yaml
release_execution_maintainer_decision: recorded
completed_work_item: W-0193
decision: ADR-0101
check_rule: runtime.release_execution_maintainer_decision
maintainer_decision: go_to_release_identifier_artifact_plan
maintainer_approval_recorded: true
release_execution_path_may_continue: true
release_declared: false
release_publishing_authorized_by_this_decision: false
release_execution_commands_authorized_by_this_decision: false
release_identifier_approved_by_this_decision: false
release_tag_authorized_by_this_decision: false
release_artifacts_authorized_by_this_decision: false
release_packaging_authorized_by_this_decision: false
hosted_deployment_authorized_by_this_decision: false
github_release_authorized_by_this_decision: false
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
next_direction: release_identifier_artifact_plan
```

The decision is a go for the next planning step only. It is not a go for any release-producing command.

## 3. Maintainer Input

The maintainer clarified the decision in conversation:

```text
我没看懂，你说需要我决策什么？我决策的结果为同意。
```

This record interprets `同意` as:

- go to continue the release execution path;
- no concrete release identifier approved yet;
- no tag approved yet;
- no artifact family approved yet;
- no publication surface approved yet;
- no release execution command approved yet.

This conservative interpretation preserves the release authorization boundary while allowing the queue to move forward.

## 4. Required Next Step

The next bounded work item is:

```text
W-0194 Define release identifier and artifact plan
```

That work item may propose or define the release identifier and artifact plan, but it still must not create a tag, GitHub release, binary, archive, container, package, checksum, provenance file, hosted deployment, or public release announcement unless a later work item explicitly authorizes execution.

`W-0194` has now recorded that plan in `docs/release-identifier-artifact-plan.md` with proposed identifier `v0.1.0-alpha.1`. The next step is `W-0195 Confirm release execution final authorization`, not release execution.

## 5. Verification State

The decision record should be checked with:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change confirm-release-execution-maintainer-decision --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Go runtime tests are not required by this documentation-only decision record, but they may be run if a later release planning or execution step changes runtime behavior. This record changes no Go runtime code.

## 6. Known Warning Handling

The known repository warning remains:

```text
runtime.identity_boundary on runtime/internal/platform/persistence/postgres/authentication_repository.go
```

This decision accepts that known warning for the purpose of continuing release planning. Any new warning must be triaged before release execution can proceed.

## 7. Stop Conditions

Stop before release execution if any of these occur:

- required checks fail;
- a new warning appears and is not explicitly triaged;
- a final release identifier is selected without a repository record;
- a tag, GitHub release, artifact, checksum, provenance file, hosted deployment, or publication command is attempted without a later execution authorization;
- any tracked artifact contains raw credentials, raw access tokens, verifier keys, lookup digests, verifier digests, HMAC input/output bytes, DSNs with credentials, transport proof carriers, GitHub tokens, or private environment file content;
- runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior, broad product modules, or direct Nakama/Pitaya API compatibility changes inside this decision slice.

## 8. Reference Alignment

Nakama and Pitaya both separate release decisions from artifact creation and deployment mechanics. This decision follows that discipline while preserving vibit's agent-native workflow and without adopting direct Nakama/Pitaya APIs, packaging models, deployment models, or compatibility promises.
