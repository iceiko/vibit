# Prototype-Ready Local Development Path Package

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Source-first local development package for prototype authors
Depends on: `docs/prototype-ready-local-development-path-gate.md`, `docs/alpha-developer-flow.md`, `docs/runtime-runbook.md`
Canonical decision: `ADR-0108`

The paired Simplified Chinese translation is `docs/prototype-ready-local-development-path-package.zh-CN.md`. The English file is authoritative.

This document packages the current local development path for vibit after the gate recorded in `ADR-0107`. It is a docs, examples, and check-rule package. It does not change production runtime behavior, add protocol routes, add Protobuf source or generated output, add migrations, add dependencies, broaden operations/admin behavior, add hosted deployments, create release artifacts, run public announcements, run paid promotion, change authentication/session behavior, add broad product modules, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The prototype-ready local development path package record is:

```yaml
prototype_ready_local_development_path_package: implemented
completed_work_item: W-0200
decision: ADR-0108
check_rule: runtime.prototype_ready_local_development_path_package
gate_decision: ADR-0107
gate_standard: docs/prototype-ready-local-development-path-gate.md
package_standard: docs/prototype-ready-local-development-path-package.md
package_standard_translation: docs/prototype-ready-local-development-path-package.zh-CN.md
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
package_scope: docs_scripts_examples_static_checks
quick_source_path_recorded: true
prerequisite_check_recorded: true
redacted_local_configuration_template_added: true
gitignore_local_secret_guard_added: true
migration_expectations_packaged: true
runtime_startup_guidance_packaged: true
example_flow_script_recorded: true
request_loop_proof_recorded: true
verification_commands_recorded: true
stop_conditions_recorded: true
next_work_item: W-0201
next_direction: storage_objects_behavior_gate
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

The source alpha already has a real authenticated gameplay loop. This package makes that loop easier to try by collecting the practical path in one place:

```text
clone source
-> check local tools
-> prepare redacted local configuration
-> apply or verify PostgreSQL migrations explicitly when using PostgreSQL
-> start the runtime in memory or PostgreSQL mode
-> run the redacted authenticated request-loop proof
-> inspect health, readiness, version, and config posture
-> continue through the next bounded work item
```

The package is intentionally honest about what is still manual. It does not pretend that source checkout is a packaged product distribution, that PostgreSQL setup is automatic, or that local onboarding is already a public protocol surface.

## 3. Fastest Path

From a source checkout:

```bash
node tools/vibit inspect next
node tools/vibit check all --json
cd runtime && go test ./...
cd .. && examples/local-alpha-request-loop.sh
```

This path proves repository checks, Go tests, and the redacted authenticated request-loop proof. It does not require a live PostgreSQL server, committed verifier keys, raw credentials, raw access tokens, a generated SDK, Docker, hosted infrastructure, or package installation.

## 4. Supported Prerequisites

The first supported local path assumes:

- Go for runtime tests and `go run ./cmd/vibit-server`;
- Node.js for `tools/vibit`;
- a POSIX-compatible shell for repository-owned local scripts;
- PostgreSQL only when evaluating the persistent runtime path;
- Buf and Protobuf tooling only when regenerating Protobuf output.

The package does not require Docker, Docker Compose, Kubernetes, cloud services, external secret managers, package registries, hosted databases, or paid services.

## 5. Local Configuration

The committed placeholder template is:

```text
examples/local.prototype.env.example
```

Use it as a field checklist only. It contains placeholders, not usable secret values.

Local private environment files are intentionally ignored by `.gitignore`:

```text
.vibit.local.env
.env.local
.env.*.local
```

Do not commit, paste, log, or include in ADRs, change records, examples, screenshots, issue reports, or shell transcripts:

- raw device credential text or bytes;
- raw access tokens;
- credential or token lookup digests;
- credential or token verifier digests;
- HMAC input or output bytes;
- verifier key values;
- concrete verifier key set ids;
- PostgreSQL DSNs with credentials;
- headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or transport metadata that may carry secrets.

The repository should be checked for obvious GitHub token leaks before commit or push:

```bash
rg -n -o "ghp_[A-Za-z0-9_]{20,}|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .
```

That command must not read or print private local env contents.

## 6. Migration And Startup Path

The two local startup postures remain:

```text
VIBIT_RUNTIME_STORE=memory
VIBIT_RUNTIME_STORE=postgres
```

The memory store is the default quick smoke posture. It is useful for checking the basic runtime process and the original in-memory inventory request loop.

The PostgreSQL store is the current alpha composition for persistent inventory, player accounts, device credential login, token issuance, runtime sessions, request-level protected routes, first-message connection binding, logout, and protected presence query.

Before using `VIBIT_RUNTIME_STORE=postgres` against a fresh local database:

1. Prepare a local PostgreSQL database.
2. Apply or verify SQL migration sources under `runtime/migrations/postgres/` explicitly.
3. Keep the DSN in private local configuration.
4. Set the required verifier key variables with local-only values.
5. Start `cd runtime && go run ./cmd/vibit-server`.

Normal server startup does not apply migrations automatically. This package does not add automatic startup migration apply behavior, schema changes, repository changes, or storage adapter changes.

## 7. Example Flow

The package records the first executable proof path as:

```bash
examples/local-alpha-request-loop.sh
```

The script wraps:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

The proof covers:

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

It uses existing runtime protocol handlers and test-owned setup. It is not a public onboarding client, product SDK, live PostgreSQL process client, hosted demo, release artifact, or compatibility promise.

## 8. Local Status Surfaces

When the runtime process is running, inspect:

```text
/healthz
/readyz
/version
/configz
```

These are local troubleshooting endpoints. `/configz` reports redacted posture and `secrets_redacted: true`. It must not print verifier keys, raw credentials, raw tokens, DSNs, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.

## 9. Verification

Default verification remains source-first and local:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change implement-prototype-ready-local-development-path-package --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./...
cd .. && examples/local-alpha-request-loop.sh
git diff --check
```

Optional live PostgreSQL verification remains opt-in through `VIBIT_POSTGRES_TEST_DSN` and a disposable database. Default checks must not require live PostgreSQL or private secrets.

## 10. Stop Conditions

Stop and ask for maintainer authorization if improving this package would require:

- production runtime behavior changes;
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
W-0201 Define storage objects behavior gate
```

That work should define the first general storage-object behavior beyond the inventory proof slice. It should remain a gate unless separately authorized to implement runtime behavior, protocol, generated output, migrations, dependencies, operations/admin breadth, authentication/session changes, broad product modules, or direct Nakama/Pitaya compatibility.
