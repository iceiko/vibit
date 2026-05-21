# Release Execution Preparation Gate

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Gate-only release execution preparation boundary for the v0.1 alpha path
Depends on: `docs/release-publishing-decision-gate.md`, `docs/alpha-developer-flow.md`, `docs/alpha-acceptance-checklist.md`, `docs/runtime-runbook.md`
Canonical decision: `ADR-0099`

The paired Simplified Chinese translation is `docs/release-execution-preparation-gate.zh-CN.md`. The English file is authoritative.

This document defines the release execution preparation gate. It does not publish `v0.1 alpha`, create release tags, binaries, archives, containers, packages, checksums, provenance files, hosted deployments, runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, broad product modules, or direct Nakama/Pitaya API compatibility.

## 1. Purpose

The release publishing decision gate established that local alpha readiness can proceed toward release preparation, but it did not authorize release execution.

This gate defines what a future release execution plan must contain before any artifact creation or publication can be considered. It is a preparation boundary, not a release run.

The gate records planning inputs, release-note boundaries, artifact plan boundaries, maintainer approval points, verification requirements, rollback notes, and stop conditions.

## 2. Core Rule

The release execution preparation gate is:

```yaml
release_execution_preparation_gate: defined
completed_work_item: W-0191
decision: ADR-0099
check_rule: runtime.release_execution_preparation_gate
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
release_execution_plan_defined: true
release_note_inputs_defined: true
artifact_plan_boundaries_defined: true
maintainer_approval_points_defined: true
rollback_stop_conditions_defined: true
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
next_direction: release_execution_authorization_gate
```

The gate allows the project to proceed to a later release execution authorization gate. It does not authorize release execution or artifact creation.

## 3. Preparation Inputs

A later release execution authorization step may only be considered after a preparation record can point to these inputs:

- `docs/v0.1-alpha-goal.md` for the release target and non-goals.
- `docs/alpha-acceptance-checklist.md` for local readiness, manual setup, and deferrals.
- `docs/alpha-developer-flow.md` for contributor-facing local flow.
- `docs/release-publishing-decision-gate.md` for the publication decision boundary.
- `docs/runtime-runbook.md` for runtime startup, PostgreSQL, verifier key, and redaction posture.
- `examples/local-alpha-request-loop.sh` for the redacted local proof command.
- `node tools/vibit inspect next` and repository checks for machine-readable state.
- The current Git branch and commit range that a future release note may summarize.
- The known warning state, currently `runtime.identity_boundary` on `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- Maintainer confirmation that the next step is authorization review, not execution.

This gate does not select a final release identifier, create a tag, sign anything, or create an artifact manifest.

## 4. Release-Note Inputs

Future release notes should be prepared from repository-owned facts only:

- current project status: pre-alpha moving toward `v0.1 alpha`;
- local developer flow summary;
- runtime surfaces: `/v1/ws`, `/healthz`, `/readyz`, `/version`, and `/configz`;
- authenticated gameplay proof path: onboarding -> login -> bind connection -> protected inventory -> presence query -> logout -> revoked-token rejection;
- PostgreSQL manual setup posture;
- redaction expectations and secret handling;
- verification command set and results;
- known warnings and their triage status;
- explicit non-goals and deferrals;
- contribution entry point through `.arch/work-items.yaml`.

Future release notes must not include raw credentials, raw access tokens, verifier keys, digests, DSNs with credentials, GitHub tokens, private environment files, or unredacted transport proof carriers.

This gate does not write final release notes or publish changelog content as a release artifact.

## 5. Artifact Plan Boundary

A later release execution plan may discuss these artifact families:

- Git version tag.
- Release notes or changelog entry.
- Source archive.
- Checksum or provenance file.
- Optional binary build.
- Optional package.
- Optional container image.
- Optional hosted deployment.

This gate creates none of them.

Any future artifact plan must answer:

- which artifact families are in scope;
- which artifact families remain deferred;
- which command would create each artifact;
- where each artifact would be stored;
- which verification must pass before creation;
- which maintainer approval point must occur before creation;
- how to stop safely if creation fails.

The first prepared alpha path should prefer the smallest public release surface that preserves source inspectability and contributor onboarding. Optional binaries, packages, containers, and hosted deployments remain deferred unless a later work item explicitly authorizes them.

## 6. Maintainer Approval Points

The following actions require later explicit maintainer authorization:

- choosing a final release identifier;
- creating a Git tag;
- pushing a Git tag;
- creating a source archive;
- creating checksums or provenance files;
- building binaries, packages, or containers;
- creating hosted deployments;
- creating a GitHub release or equivalent release record;
- announcing public release availability.

Authorization should be represented by a bounded work item and a repository record. Chat-only approval is not enough for durable project state.

## 7. Verification Requirements

Use this command set before treating the preparation gate as complete:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-release-execution-preparation-gate --json
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

## 8. Rollback And Stop Conditions

Because this gate creates no release artifacts, rollback is a documentation and tooling revert.

Stop before any later release execution authorization if any of these occur:

- a required repository check or Go test fails;
- a new warning appears and is not explicitly triaged;
- the known warning state changes without a documented decision;
- a public-facing English document lacks its Simplified Chinese translation;
- any tracked artifact contains raw credential material, raw access tokens, verifier key values, lookup digests, verifier digests, HMAC input/output bytes, DSNs with credentials, transport proof carriers, or GitHub tokens;
- generated output changes without a generation step and source trace;
- Protobuf sources, generated output, migrations, dependencies, or runtime behavior change inside a preparation-only slice;
- any release artifact is created;
- any release tag is created, signed, or pushed;
- hosted deployment work starts;
- broad operations/admin behavior is added;
- authentication/session behavior changes without a separate authorized work item;
- direct Nakama/Pitaya API compatibility is selected;
- the maintainer has not explicitly authorized the next release execution authorization gate.

## 9. Redaction Expectations

Preparation records must remain safe to publish in the repository.

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

## 10. Reference Alignment

Nakama and Pitaya demonstrate that serious server frameworks separate release readiness, artifacts, support posture, and deployment choices. This gate adopts that discipline without adopting their APIs, data models, route names, release packaging, deployment model, cluster model, SDK surfaces, or operations surfaces.

## 11. Next Work

The next bounded contribution is:

```text
W-0192 Define release execution authorization gate
```

That future step may define authorization criteria for release execution, but it still must not publish a release or create release artifacts unless its own scope explicitly says so and the maintainer authorizes the ask-first boundary.
