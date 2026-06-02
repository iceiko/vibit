# Change Request: Implement Pitaya-Aligned Runtime Component Lifecycle Source-First Map

Date: 2026-06-02
Status: Accepted

## Request

Implement the next-ready work item, `W-0281 Implement Pitaya-aligned runtime component lifecycle source-first map`.

## Required Outcome

- Add `node tools/vibit inspect pitaya-component-lifecycle --json`.
- Register `runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map`.
- Report runtime component lifecycle boundary vocabulary, current source-first mappings, source surfaces, redaction posture, and implementation deferrals.
- Open `W-0282 Select next Pitaya-aligned direction after runtime component lifecycle map` as next-ready.

## Non-Goals

Do not add runtime component lifecycle behavior, handler registration behavior, component discovery or loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.
