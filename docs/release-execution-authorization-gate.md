# Release Execution Authorization Gate

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Gate-only release execution authorization criteria for the v0.1 alpha path
Depends on: `docs/release-execution-preparation-gate.md`, `docs/release-publishing-decision-gate.md`, `docs/alpha-developer-flow.md`, `docs/alpha-acceptance-checklist.md`, `docs/runtime-runbook.md`
Canonical decision: `ADR-0100`

The paired Simplified Chinese translation is `docs/release-execution-authorization-gate.zh-CN.md`. The English file is authoritative.

This document defines the release execution authorization gate. It does not publish `v0.1 alpha`, choose or create release identifiers or tags, create binaries, archives, containers, packages, checksums, provenance files, hosted deployments, runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, broad product modules, or direct Nakama/Pitaya API compatibility.

## 1. Purpose

The release execution preparation gate defined the inputs and planning boundaries for a future release execution step. The remaining risk is treating readiness criteria as permission to execute the release.

This gate defines the authorization criteria that must be satisfied before maintainers can make a later go/no-go decision. It is an authorization criteria document, not a release execution record.

The gate records final go/no-go criteria, required verification state, release identifier review requirements, artifact authorization boundaries, maintainer approval requirements, and stop conditions.

## 2. Core Rule

The release execution authorization gate is:

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

The gate lets maintainers review whether release execution may be authorized later. It does not itself authorize release execution, select a release identifier, create artifacts, or publish a release.

## 3. Final Go/No-Go Criteria

A later release execution maintainer decision may only be considered if all of these are true:

- `docs/v0.1-alpha-goal.md` still describes the intended `v0.1 alpha` scope accurately.
- `docs/alpha-acceptance-checklist.md` records no unresolved `Blocked` item for the local alpha flow.
- `docs/alpha-developer-flow.md` still points contributors through a coherent local path.
- `docs/runtime-runbook.md` still matches current runtime startup, PostgreSQL setup, verifier key, and redaction posture.
- `docs/release-publishing-decision-gate.md` and `docs/release-execution-preparation-gate.md` remain satisfied.
- `examples/local-alpha-request-loop.sh` still runs the redacted authenticated gameplay proof.
- Repository checks and Go tests pass with no new untriaged warning.
- The known warning, if still present, remains the documented `runtime.identity_boundary` warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- English canonical documents and Simplified Chinese translations remain paired for public-facing release-path documents.
- No tracked artifact contains raw credentials, raw access tokens, verifier keys, lookup digests, verifier digests, HMAC input/output bytes, DSNs with credentials, transport proof carriers, GitHub tokens, or private environment file content.
- The maintainer explicitly chooses to make a release execution go/no-go decision.

If any criterion is not satisfied, the release execution decision must be no-go until the issue is resolved or explicitly deferred by a later repository record.

## 4. Required Verification State

The required default verification state for a later go/no-go decision is:

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

The later release decision record should capture command results, known warnings, and any intentionally skipped optional verification.

Optional live PostgreSQL checks remain opt-in through `VIBIT_POSTGRES_TEST_DSN` and a disposable database. This gate does not make live PostgreSQL verification mandatory for default repository checks.

## 5. Release Identifier Review

This gate does not choose a release identifier.

A later release execution decision must review:

- whether the target identifier is `v0.1 alpha` or another explicitly chosen pre-release identifier;
- whether the identifier already exists as a Git tag, release record, package, image, or external artifact;
- whether the identifier matches project naming and public communication expectations;
- whether the identifier implies a stability level beyond the current pre-alpha/alpha posture;
- whether the identifier can be safely represented in release notes, tags, archives, checksums, packages, or hosted deployments if those are later authorized.

Selecting the final identifier remains an ask-first boundary and must be recorded in a later bounded work item.

## 6. Artifact Authorization Boundary

This gate authorizes no artifacts.

A later maintainer decision must explicitly say which artifact families, if any, are authorized:

- Git version tag.
- GitHub release or equivalent release record.
- Release notes or changelog entry.
- Source archive.
- Checksum file.
- Provenance file.
- Binary build.
- Package.
- Container image.
- Hosted deployment.

Artifact authorization should prefer the smallest surface that lets developers inspect source and run the local alpha path. Optional binaries, packages, containers, provenance files, and hosted deployments should remain deferred unless explicitly approved by a later work item.

## 7. Maintainer Approval Requirements

A later release execution decision requires a durable maintainer approval record that answers:

- Is release execution go or no-go?
- What release identifier, if any, is approved?
- Which artifact families, if any, are approved?
- Which verification results were reviewed?
- Which known warnings were accepted or rejected?
- Which deferrals remain in force?
- Which commands are allowed to create approved artifacts?
- Which stop conditions still apply during execution?

Chat-only approval is not enough. Approval must be represented by a repository artifact and a bounded work item before any release execution command runs.

## 8. Authorization Outcome

The outcome of this gate is:

```yaml
authorization_criteria_defined: true
may_make_maintainer_go_no_go_decision_later: true
may_publish_release_now: false
may_execute_release_now: false
may_select_release_identifier_now: false
may_create_release_artifacts_now: false
```

The repository remains pre-alpha. The next step is a maintainer decision gate, not release execution.

## 9. Stop Conditions

Stop before any later release execution if any of these occur:

- required verification fails;
- a new warning appears and is not explicitly triaged;
- the known warning state changes without a documented decision;
- public-facing English and Simplified Chinese documents diverge materially;
- any tracked artifact contains secrets or unredacted proof material;
- generated output changes without source trace and generation notes;
- runtime behavior, protocol routes, Protobuf sources, migrations, dependencies, operations/admin behavior, authentication/session behavior, broad product modules, or direct Nakama/Pitaya compatibility change inside an authorization-only slice;
- a release identifier is selected without maintainer approval;
- a release tag, artifact, package, checksum, provenance file, hosted deployment, or release record is created;
- the release scope implies production readiness, SLA, hosted service availability, or compatibility promises beyond the alpha goal;
- maintainer approval is not recorded in a repository artifact.

## 10. Redaction Expectations

Authorization records must remain safe to publish in the repository.

Do not include:

- raw device credential text or bytes,
- raw access tokens,
- credential or token lookup digests,
- credential or token verifier digests,
- HMAC input or output bytes,
- verifier key values,
- concrete verifier key set ids,
- PostgreSQL DSNs with credentials,
- headers, cookies, query strings, WebSocket subprotocol values, or remote addresses that may carry secrets,
- local files that contain GitHub tokens or other access credentials.

## 11. Reference Alignment

Nakama and Pitaya demonstrate that serious server frameworks separate release readiness, release decisions, artifacts, and deployment posture. This gate adopts that discipline without adopting their APIs, data models, route names, release packaging, deployment model, cluster model, SDK surfaces, or operations surfaces.

## 12. Next Work

`W-0193 Confirm release execution maintainer decision` has now recorded the maintainer go decision in `docs/release-execution-maintainer-decision.md`. The next bounded work is:

```text
W-0194 Define release identifier and artifact plan
```

That future planning step still must not create release tags, artifacts, hosted deployments, or published release records unless a later execution scope explicitly authorizes those actions and the maintainer approves the ask-first boundary.
