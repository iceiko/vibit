# Prototype-Ready Local Development Path Gate

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Gate for a repeatable local development path before prototype-ready implementation packaging
Depends on: `docs/prototype-ready-foundation-execution-plan.md`, `docs/alpha-developer-flow.md`, `docs/runtime-runbook.md`
Canonical decision: `ADR-0107`

The paired Simplified Chinese translation is `docs/prototype-ready-local-development-path-gate.zh-CN.md`. The English file is authoritative.

This document defines the gate for making vibit's local development path repeatable enough for prototype authors. It is a gate artifact. It does not implement runtime behavior, add protocol routes, add Protobuf source or generated output, add migrations, add dependencies, broaden operations/admin behavior, add hosted deployments, create release artifacts, run public announcements, run paid promotion, change authentication/session behavior, add broad product modules, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The prototype-ready local development path gate record is:

```yaml
prototype_ready_local_development_path_gate: defined
completed_work_item: W-0199
decision: ADR-0107
check_rule: runtime.prototype_ready_local_development_path_gate
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
source_execution_plan: docs/prototype-ready-foundation-execution-plan.md
gate_standard: docs/prototype-ready-local-development-path-gate.md
gate_standard_translation: docs/prototype-ready-local-development-path-gate.zh-CN.md
future_implementation_work_item: W-0200
future_implementation_direction: prototype_ready_local_development_path_package
supported_prerequisites_recorded: true
startup_expectations_recorded: true
migration_expectations_recorded: true
configuration_secret_posture_recorded: true
example_flow_shape_recorded: true
allowed_future_write_areas_recorded: true
verification_expectations_recorded: true
stop_conditions_recorded: true
planning_only: true
gate_only: true
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

The next product problem is not only whether vibit can pass checks. A prototype author needs a local path that is understandable, repeatable, and honest about what is still manual.

The local development path should let a technically capable developer:

- install the required local tools;
- configure safe local secrets without committing them;
- prepare or verify the PostgreSQL schema explicitly;
- start the runtime in memory mode for a quick smoke path;
- start the PostgreSQL runtime path when evaluating the current alpha composition;
- run a meaningful authenticated flow;
- inspect health, readiness, version, and redacted configuration posture;
- know which missing pieces are deferred instead of guessing.

This gate keeps the work source-first and local. It does not turn vibit into a production deployment, a hosted demo, a binary distribution, an SDK package, or a compatibility clone of Nakama or Pitaya.

## 3. Supported Local Prerequisites

The first supported local path may assume:

- Go for runtime tests and `go run ./cmd/vibit-server`;
- Node.js for `tools/vibit` checks;
- a POSIX-compatible shell for repository-owned local scripts;
- PostgreSQL for the persistent prototype path;
- Buf and Protobuf tooling only when regenerating Protobuf output.

The first package should not require Docker, Docker Compose, Kubernetes, cloud infrastructure, external secret managers, package registries, hosted databases, or paid services. Those may be useful later, but they require separate authorization because they change the operations and dependency posture.

## 4. Startup Expectations

The local package should document two startup postures:

- `VIBIT_RUNTIME_STORE=memory`, the default quick smoke posture;
- `VIBIT_RUNTIME_STORE=postgres`, the current alpha composition for persistent inventory, player accounts, device credential login, token issuance, runtime sessions, request-level route protection, first-message binding, logout, and presence query.

The PostgreSQL posture requires:

- `VIBIT_POSTGRES_DSN`;
- `VIBIT_AUTH_VERIFIER_KEY_SET_ID`;
- `VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY`;
- `VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY`;
- `VIBIT_AUTH_TOKEN_LOOKUP_KEY`;
- `VIBIT_AUTH_TOKEN_VERIFIER_KEY`;
- optional token lifetime and audience variables already documented in `docs/runtime-runbook.md`.

The package may improve documentation and scripts around `cd runtime && go run ./cmd/vibit-server`, local ports, status endpoints, and redacted diagnostics. It must not silently start hosted services, apply migrations during normal server startup, or hide missing configuration failures.

## 5. Migration Expectations

The first local development path should make migrations explicit:

- SQL migration sources remain under `runtime/migrations/postgres/`;
- normal `vibit-server` startup does not apply migrations;
- local setup instructions should tell the developer when a fresh database needs migration apply or status verification;
- any future script must call existing repository-owned migration behavior or documented commands explicitly;
- migration output must not print DSNs or secrets.

This gate does not authorize new migration files, repository/storage adapter changes, schema changes, automatic startup migration apply behavior, or new migration dependencies.

## 6. Configuration And Secret Posture

The first local package may add a redacted example environment file or documented local environment template, but it must never commit real local secrets.

The local path must treat the following as not log-safe:

- raw device credential text or bytes;
- raw access tokens;
- credential or token lookup digests;
- credential or token verifier digests;
- HMAC input or output bytes;
- verifier key values;
- concrete verifier key set ids;
- PostgreSQL DSNs with credentials;
- headers, cookies, query strings, WebSocket subprotocol values, or remote addresses that may carry secrets;
- private local environment files such as `.vibit.local.env`.

Committed examples must use placeholders. Generated local secrets, if documented later, must be created by the developer locally and kept outside version control.

## 7. Example Flow Shape

The first prototype-ready local package should demonstrate a multi-step flow, not a single isolated request:

```text
prerequisite check
-> redacted local configuration preparation
-> explicit migration status or apply step
-> runtime startup guidance
-> local onboarding through the existing application-owned path
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

The existing `examples/local-alpha-request-loop.sh` already proves much of the flow through a focused Go E2E test. The future package may wrap, document, or extend the local example ergonomics, but this gate does not authorize a new public onboarding route, new WebSocket routes, new Protobuf messages, SDK publication, or a production client surface.

## 8. Allowed Future Write Areas

`W-0200 Implement prototype-ready local development path package` may use these write areas, provided it stays inside this gate:

- `README.md` and `README.zh-CN.md`;
- `docs/runtime-runbook.md` and its translation if a paired translation exists later;
- `docs/alpha-developer-flow.md` and `docs/alpha-developer-flow.zh-CN.md`;
- `docs/alpha-acceptance-checklist.md` and `docs/alpha-acceptance-checklist.zh-CN.md`;
- new `docs/` files for the local development path package and paired Simplified Chinese translations;
- `examples/` scripts, README files, or placeholder environment templates;
- `.gitignore` entries for local-only environment files when needed;
- `tools/vibit` and `rules/check-rules.json` for static repository checks;
- `.arch/` manifests, change specs, ADRs, AGENTS guides, and conversation memory;
- focused tests or script checks that verify existing behavior without changing production runtime behavior.

The future package must ask first before modifying production Go runtime behavior, protocol source, generated output, SQL migrations, repository interfaces, dependencies, release artifacts, hosted deployment surfaces, broad operations/admin behavior, authentication/session semantics, or broad product modules.

## 9. Verification Expectations

The gate expects the future package to keep default verification local and source-first:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```

Optional live PostgreSQL verification may remain opt-in through `VIBIT_POSTGRES_TEST_DSN` and a disposable database. Default repository checks must not require a running database or private secrets.

## 10. Stop Conditions

Stop and ask for maintainer authorization if executing the local path would require:

- runtime behavior changes;
- protocol route changes;
- Protobuf source or generated output changes;
- SQL migrations, repository interface changes, or storage adapter changes;
- new dependencies;
- automatic startup migration apply behavior;
- broad operations/admin behavior;
- authentication/session semantic changes;
- public local onboarding routes or new proof carriers;
- broad product module expansion;
- direct Nakama/Pitaya API compatibility;
- hosted deployments or demos;
- release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, SDK packages, or additional release artifacts;
- public announcements beyond the GitHub release record;
- paid promotion;
- handling or disclosure of secrets.

## 11. Next Work

The next bounded direction is:

```text
W-0200 Implement prototype-ready local development path package
```

That work should package the local development path using the gate above, then leave later product-capability work such as storage objects, realtime messaging or server push, failure/concurrency verification, and operations inspection behind separate work items.

## 12. Verification

The repository should verify:

- this gate and its translation exist;
- `ADR-0107` records the decision;
- `.arch` manifests mark `W-0199` completed and open `W-0200`;
- README, alpha goal, developer flow, acceptance checklist, AGENTS guides, and product roadmap point at the new next work;
- runtime, protocol, generated output, migration, dependency, operations/admin, authentication/session, product module, hosted deployment, release artifact, public announcement, paid promotion, and direct compatibility deferrals remain preserved.
