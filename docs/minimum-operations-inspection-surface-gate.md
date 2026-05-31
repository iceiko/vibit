# Minimum Operations Inspection Surface Gate

Status: Accepted v0.1
Last updated: 2026-05-31
Scope: Gate-only boundary for the first minimum source-first operations inspection surface after friends relationship route proof
Depends on: `decisions/ADR-0151-select-next-nakama-prototype-ready-capability-after-friends-route-proof.md`, `docs/runtime-runbook.md`, `docs/alpha-developer-flow.md`, `docs/alpha-acceptance-checklist.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0152`

The paired Simplified Chinese translation is `docs/minimum-operations-inspection-surface-gate.zh-CN.md`. The English file is authoritative.

This document defines the minimum operations inspection surface gate. It is a gate artifact. It does not implement operations/admin endpoints, metrics endpoints, observability pipelines, dashboards, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, event/audit tables, groups, parties, chat, stream subscriptions, matchmaking, match runtime, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The minimum operations inspection surface gate record is:

```yaml
minimum_operations_inspection_surface_gate: defined
completed_work_item: W-0244
decision: ADR-0152
check_rule: runtime.minimum_operations_inspection_surface_gate
source_selection_decision: ADR-0151
source_friends_route_proof_decision: ADR-0150
source_workflow_decision: ADR-0128
standard: docs/minimum-operations-inspection-surface-gate.md
translation: docs/minimum-operations-inspection-surface-gate.zh-CN.md
selected_nakama_capability_family: admin_console_metrics_observability_and_operations
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
surface_posture: source_first_local_operations_inspection
admin_console_posture: deferred
metrics_endpoint_posture: deferred
observability_pipeline_posture: deferred
future_implementation_work_item: W-0245
future_implementation_direction: implement_minimum_operations_inspection_source_first_surface
future_cli_inspection_candidate: tools/vibit inspect operations
future_docs_candidate: docs/runtime-runbook.md
future_acceptance_checklist_candidate: docs/alpha-acceptance-checklist.md
accepted_existing_runtime_surfaces:
  - /healthz
  - /readyz
  - /version
  - /configz
accepted_existing_source_surfaces:
  - .arch/work-items.yaml
  - .arch/runtime.yaml
  - docs/runtime-runbook.md
  - docs/alpha-developer-flow.md
  - docs/alpha-acceptance-checklist.md
  - examples/local-alpha-example-client.sh
  - examples/local-alpha-request-loop.sh
minimum_inspectable_state_categories:
  - process_liveness_and_readiness
  - runtime_version_and_release_posture
  - runtime_store_and_configuration_posture
  - protocol_route_family_inventory
  - local_alpha_flow_steps
  - persistence_and_migration_posture
  - authentication_session_connection_posture
  - generated_output_and_proto_posture
  - repository_check_and_verification_posture
  - deferred_operations_surfaces
redaction_required: true
runtime_behavior_added_by_this_gate: false
operations_admin_endpoint_added: false
admin_console_added: false
metrics_endpoint_added: false
observability_pipeline_added: false
dashboard_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
event_audit_table_added: false
hosted_deployment_added: false
sdk_added: false
distributed_runtime_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0151` selected `admin_console_metrics_observability_and_operations` as the next Nakama-first prototype-ready capability family after the protected friends relationship route family was proven locally.

The current local alpha already exposes small HTTP troubleshooting endpoints:

```text
/healthz
/readyz
/version
/configz
```

It also has source-first inspection through architecture manifests, runbooks, example scripts, repository checks, and focused Go tests. What it does not yet have is a clear minimum operations inspection posture that tells prototype authors which state is inspectable, which state must stay redacted, and which operational surfaces remain future work.

The first operations inspection step should therefore be source-first and local. It should help a developer or agent understand current server state and verification posture without creating an admin console, metrics backend, telemetry pipeline, sensitive state dump, or compatibility promise.

## 3. Minimum Inspectable Categories

The future implementation should expose or document these categories in a source-first way:

```yaml
minimum_inspectable_state_categories:
  process_liveness_and_readiness:
    existing_runtime_surface:
      - /healthz
      - /readyz
    allowed_future_source_inspection:
      - summarize endpoint meaning and expected local-alpha output
  runtime_version_and_release_posture:
    existing_runtime_surface:
      - /version
    allowed_future_source_inspection:
      - show runtime version, pre-alpha status, release identifier posture, and source-first release boundary
  runtime_store_and_configuration_posture:
    existing_runtime_surface:
      - /configz
    allowed_future_source_inspection:
      - report memory versus PostgreSQL posture and whether required configuration is present, never secret values
  protocol_route_family_inventory:
    existing_source_surface:
      - runtime route constants
      - Protobuf bridge registrations
      - docs/runtime-runbook.md
    allowed_future_source_inspection:
      - list known local alpha route families and whether they require authenticated wrappers
  local_alpha_flow_steps:
    existing_source_surface:
      - examples/local-alpha-example-client.sh
      - examples/local-alpha-request-loop.sh
      - runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
    allowed_future_source_inspection:
      - summarize onboarding, login, binding, protected gameplay, friends, presence, storage, logout, and failure proof steps
  persistence_and_migration_posture:
    existing_source_surface:
      - runtime/migrations/postgres/
      - docs/runtime-runbook.md
    allowed_future_source_inspection:
      - list migration source presence and manual-apply expectation, not database contents
  authentication_session_connection_posture:
    existing_source_surface:
      - docs/runtime-runbook.md
      - .arch/runtime.yaml
    allowed_future_source_inspection:
      - summarize proof carrier, route protection, session metadata, connection binding, logout, and redaction posture
  generated_output_and_proto_posture:
    existing_source_surface:
      - proto/
      - runtime/internal/generated/proto/
      - docs/generated-output.md
    allowed_future_source_inspection:
      - show generated-output traceability and no-hand-edit posture
  repository_check_and_verification_posture:
    existing_source_surface:
      - tools/vibit
      - rules/check-rules.json
    allowed_future_source_inspection:
      - summarize required checks and known warnings
  deferred_operations_surfaces:
    allowed_future_source_inspection:
      - state that admin console, metrics endpoint, observability pipeline, live dashboards, player/session/token inspectors, hosted operations, SDKs, and distributed runtime remain deferred
```

The first implementation should prefer `tools/vibit inspect operations` or equivalent source-first repository inspection over new runtime endpoints.

## 4. Redaction

Any future inspection surface must not print, persist, log, serialize, or record:

- raw device credential material;
- raw access tokens;
- credential or token lookup digests;
- credential or token verifier digests;
- verifier key values;
- concrete verifier key set ids;
- HMAC inputs or outputs;
- PostgreSQL DSNs with credentials;
- passwords, connection strings, or local secret file contents;
- HTTP headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata;
- full player, session, token, connection, relationship, storage, or presence identifiers unless a later redaction policy explicitly classifies them as safe in that surface;
- database row payloads or arbitrary JSON object values.

Allowed public output includes route names, endpoint names, file paths, high-level capability status, redacted configuration posture, verification command names, and broad state categories.

## 5. Ownership

Future implementation ownership:

```yaml
source_first_inspection_owner: tools/vibit
runtime_status_endpoint_owner: runtime/cmd/vibit-server
runbook_owner: docs/runtime-runbook.md
alpha_flow_owner: docs/alpha-developer-flow.md
acceptance_checklist_owner: docs/alpha-acceptance-checklist.md
architecture_memory_owner:
  - .arch/work-items.yaml
  - .arch/runtime.yaml
  - .arch/reference.yaml
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
```

Rules:

- `tools/vibit` may summarize already committed source, manifests, docs, and check status.
- Runtime endpoints remain limited to the existing local-alpha troubleshooting surface until a later bounded implementation explicitly changes them.
- Runtime behavior, protocol payloads, persistence, migrations, repository interfaces, adapters, authentication/session semantics, and generated output remain owned by their existing boundaries.
- No domain module should become the owner of broad operations inspection.

## 6. Nakama And Pitaya Mapping

Nakama reference mapping:

- This gate covers the `admin_console_metrics_observability_and_operations` capability family.
- It adopts the product pressure that backend developers need to inspect health, version, configuration posture, feature availability, and operational state while evaluating a framework.
- It does not copy Nakama console routes, REST paths, metrics names, dashboard behavior, SDK shapes, account/session/player inspectors, or compatibility promises.

Pitaya reference mapping:

- Pitaya remains deferred as a future distributed architecture reference.
- This gate does not introduce frontend/backend server roles, RPC, service discovery, cluster metrics, distributed session inspection, or distributed operations behavior.

## 7. Future Implementation Work

Open:

```text
M-173/W-0245 Implement minimum operations inspection source-first surface
```

The future work item may:

- add a source-first `tools/vibit inspect operations` command or equivalent inspection subcommand;
- update `docs/runtime-runbook.md` and its translation to describe the minimum inspection workflow;
- update `docs/alpha-developer-flow.md` and checklist references;
- add repository checks that verify redaction and accepted categories;
- summarize existing local alpha routes, endpoints, manifests, and verification posture.

The future work item must not:

- add an admin console, metrics endpoint, observability pipeline, dashboard, hosted operations surface, player/session/token inspector endpoint, or database-state dump;
- add runtime behavior, new HTTP endpoint behavior, protocol messages or routes, Protobuf source, generated output, migrations, dependencies, repository interfaces, PostgreSQL adapters, authentication/session behavior changes, event/audit tables, SDK publication, hosted deployments, release artifacts, distributed runtime, or direct Nakama/Pitaya API compatibility.

## 8. Verification Expectations

The future implementation should verify:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.minimum_operations_inspection_surface_gate
node tools/vibit check change define-minimum-operations-inspection-surface-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

If the future implementation changes Go runtime behavior, protocol code, or tests, it must also run focused Go tests and `cd runtime && go test ./...`. The gate itself does not require Go tests because it adds no Go behavior.

## 9. Stop Conditions

Stop and create a separate gate if the implementation requires:

- a new HTTP operations/admin endpoint;
- metrics, tracing, logging, or observability backend dependencies;
- live player, session, token, connection, relationship, or storage inspectors;
- database row inspection or arbitrary state dumps;
- public or hosted admin UI behavior;
- WebSocket or Protobuf operations routes;
- changes to authentication/session/route protection behavior;
- startup wiring or process lifecycle changes;
- event/audit tables;
- generated client libraries, SDK publication, hosted deployment, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya API compatibility.
