# Release Publishing Decision Gate

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Gate-only release publishing decision boundary for the v0.1 alpha path
Depends on: `docs/v0.1-alpha-goal.md`, `docs/alpha-acceptance-checklist.md`, `docs/alpha-developer-flow.md`
Canonical decision: `ADR-0098`

The paired Simplified Chinese translation is `docs/release-publishing-decision-gate.zh-CN.md`. The English file is authoritative.

This document defines the release publishing decision gate. It does not publish `v0.1 alpha`, create release tags, binaries, archives, containers, packages, checksums, hosted deployments, runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, broad product modules, or direct Nakama/Pitaya API compatibility.

## 1. Purpose

The local alpha path now has a packaged developer journey, acceptance checklist, runtime runbook, request-loop script, status endpoints, and focused authenticated gameplay proof. The next risk is confusing "ready to prepare a release" with the act of release publishing.

This gate separates those two decisions.

The gate records what must be true before a later release execution step may be prepared. It also records what release artifacts remain forbidden until a later explicit work item authorizes them.

## 2. Core Rule

The release publishing decision gate is:

```yaml
release_publishing_decision_gate: defined
completed_work_item: W-0190
decision: ADR-0098
check_rule: runtime.release_publishing_decision_gate
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
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
next_direction: release_execution_preparation_gate
```

The gate allowed the project to proceed to the release execution preparation gate now recorded in `docs/release-execution-preparation-gate.md`. It does not authorize release execution.

## 3. Publishing Prerequisites

A later release execution preparation step may only be considered when all of these are true:

- `W-0189` has packaged the local alpha developer flow.
- `docs/alpha-acceptance-checklist.md` records the local alpha readiness, manual setup requirements, deferred work, and release deferrals.
- `docs/alpha-developer-flow.md` gives a coherent local developer journey.
- `docs/runtime-runbook.md` describes memory and PostgreSQL startup posture.
- `examples/local-alpha-request-loop.sh` remains a redacted local proof over the authenticated gameplay E2E test.
- `/healthz`, `/readyz`, `/version`, and `/configz` remain local troubleshooting surfaces and `/configz` keeps `secrets_redacted: true`.
- Required repository checks and Go tests pass.
- Any known warning is explicitly documented and does not indicate a new release blocker. The current known warning remains `runtime.identity_boundary` on `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- English canonical docs and Simplified Chinese translations are present for public-facing release-path documents.
- No raw credential, raw access token, verifier key, digest, DSN with credentials, header, cookie, query string, WebSocket subprotocol value, remote address, or GitHub token is recorded in tracked artifacts.
- The maintainer explicitly authorizes a later release execution preparation work item.

These prerequisites are necessary, not sufficient, for release execution.

## 4. Decision Outcome

The decision outcome of this gate is:

```yaml
release_decision_outcome:
  local_alpha_flow_packaged: true
  release_publishing_decision_gate_defined: true
  may_prepare_release_execution_gate_later: true
  may_publish_release_now: false
  may_create_release_artifacts_now: false
```

The repository remains pre-alpha. Future release preparation must be represented by a bounded work item and must keep release execution separate from release publication until explicitly authorized.

## 5. Release Artifact Boundary

Future release execution work may discuss these artifact families only after a later work item authorizes them:

- Git version tags.
- Release notes or changelog entries.
- Source archives.
- Checksums or provenance files.
- Optional binaries.
- Optional packages.
- Optional container images.
- Optional hosted deployments.

This gate creates none of those artifacts.

Forbidden in this gate:

- creating, signing, or pushing tags;
- creating binaries, archives, packages, containers, checksums, or provenance files;
- publishing GitHub releases or equivalent release records;
- deploying a hosted service;
- changing runtime startup behavior for release packaging;
- changing protocol routes or Protobuf sources;
- regenerating output for release reasons;
- adding migrations or dependencies;
- broadening operations/admin behavior;
- changing authentication/session behavior;
- adding broad product modules;
- selecting direct Nakama/Pitaya API compatibility.

## 6. Verification Requirements

Use this command set before treating the gate as complete:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-release-publishing-decision-gate --json
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

Optional live PostgreSQL checks remain opt-in through `VIBIT_POSTGRES_TEST_DSN` and a disposable database. This gate does not make live PostgreSQL verification mandatory for default repository checks.

## 7. Stop Conditions

Stop before any later release execution preparation if any of these occur:

- a required repository check or Go test fails;
- a new warning appears and is not explicitly triaged;
- a public-facing English document lacks its Simplified Chinese translation;
- any tracked artifact contains raw credential material, raw access tokens, verifier key values, lookup digests, verifier digests, HMAC input/output bytes, DSNs with credentials, transport proof carriers, or GitHub tokens;
- generated output changes without a generation step and source trace;
- Protobuf sources, generated output, migrations, dependencies, or runtime behavior change inside a release decision-only slice;
- the release artifact boundary is crossed;
- a hosted deployment is created or implied;
- broad operations/admin behavior is added;
- authentication/session behavior changes without a separate authorized work item;
- direct Nakama/Pitaya API compatibility is selected;
- the maintainer has not explicitly authorized the next release preparation step.

## 8. Redaction Expectations

Release decision records must remain safe to publish in the repository.

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

## 9. Reference Alignment

Nakama and Pitaya both set an expectation that a serious server framework has a coherent local developer path and clear release support posture. This gate uses that expectation only to shape release readiness discipline.

It does not copy Nakama or Pitaya APIs, release packaging, deployment model, cluster model, route names, data models, SDK surfaces, or operations surfaces.

## 10. Next Work

The next bounded contribution is:

```text
W-0192 Define release execution authorization gate
```

`W-0191` completed the release execution preparation gate. The next step may define authorization criteria for release execution, but it still must not publish a release or create release artifacts unless its own scope explicitly says so and the maintainer authorizes the ask-first boundary.
