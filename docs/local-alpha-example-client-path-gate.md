# Local Alpha Example Client Path Gate

Status: Accepted v0.1
Last updated: 2026-05-24
Scope: Gate-only boundary for the first clearer source-first local alpha example client or example app path
Depends on: `decisions/ADR-0132-select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof.md`, `docs/agent-native-feature-request-test-workflow.md`, `docs/prototype-ready-local-development-path-package.md`, `docs/alpha-developer-flow.md`, `examples/README.md`
Canonical decision: `ADR-0133`

The paired Simplified Chinese translation is `docs/local-alpha-example-client-path-gate.zh-CN.md`. The English file is authoritative.

This document defines the local alpha example client path gate. It is a gate artifact. It does not implement an example client or example app, publish an SDK, generate client libraries, add dependencies, add runtime behavior, add protocol routes, add Protobuf source, change generated output, add migrations, add persistence, change startup wiring, change authentication/session behavior, add delivery guarantees, add stream subscriptions, add chat rooms, add groups, add broadcast fanout, add matchmaking, add match runtime, add operations/admin behavior, add hosted deployments, create release artifacts, run public announcements, run paid promotion, add Pitaya-style distributed architecture, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The local alpha example client path gate record is:

```yaml
local_alpha_example_client_path_gate: defined
completed_work_item: W-0225
decision: ADR-0133
check_rule: runtime.local_alpha_example_client_path_gate
source_selection_decision: ADR-0132
source_workflow_decision: ADR-0128
standard: docs/local-alpha-example-client-path-gate.md
translation: docs/local-alpha-example-client-path-gate.zh-CN.md
selected_nakama_capability_family: client_sdks_examples_and_developer_experience
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
path_posture: source_first_local_alpha_example_client_path
public_sdk_posture: not_a_public_sdk
live_network_client_posture: deferred_until_public_onboarding_and_client_package_boundary
future_example_docs_candidate: examples/local-alpha-client/README.md
future_example_script_candidate: examples/local-alpha-example-client.sh
future_runtime_proof_candidate: runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
future_implementation_work_item: W-0226
future_implementation_direction: implement_local_alpha_example_client_path
accepted_existing_routes:
  - runtime.authentication.AuthenticateWithDeviceCredential
  - runtime.connection.BindConnection
  - inventory.GrantItem
  - inventory.GetInventory
  - presence.GetPlayerPresence
  - storage.GetOwnStorageObject
  - storage.ListOwnStorageObjects
  - storage.PutOwnStorageObject
  - storage.DeleteOwnStorageObject
  - runtime.authentication.LogoutAccessToken
implementation_added_by_this_gate: false
example_client_implementation_added: false
sdk_added: false
generated_client_library_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
persistence_added: false
startup_wiring_added: false
authentication_session_behavior_changed: false
delivery_guarantees_added: false
stream_subscription_added: false
chat_added: false
group_messaging_added: false
broadcast_fanout_added: false
matchmaking_added: false
match_runtime_added: false
operations_admin_behavior_added: false
hosted_deployment_added: false
release_artifact_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0132` selected `client_sdks_examples_and_developer_experience` as the next Nakama-first prototype-ready capability family. The local alpha already proves authentication, connection binding, protected inventory, presence, storage objects, logout, realtime outbound foundations, and failure-path behavior, but the visible path is still mostly an internal Go E2E proof wrapped by a small shell script.

The next implementation should make that path easier for developers and AI agents to read as a client-like flow before vibit adds more broad product capability. The first example path must stay honest about current alpha constraints:

- local onboarding is an application service behavior, not a public client route;
- generated Protobuf Go output is under `runtime/internal/`, so it is not yet a public client package;
- the current proof is source-first and repository-local;
- secrets and transport metadata must remain redacted;
- the example is not an SDK, hosted demo, package, or compatibility promise.

## 3. Selected Path Shape

The first allowed path is a source-first local alpha example path:

```text
examples/local-alpha-client/README.md
examples/local-alpha-example-client.sh
runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
```

The future implementation may add a clearer example-client README and wrapper script that point to a focused, named runtime proof. If the existing E2E test is too dense for that purpose, the future implementation may add a focused local alpha example-flow test next to the existing Protobuf E2E tests, but it must reuse existing runtime and protocol surfaces.

Rules:

- The example path must be source-first and local to the repository.
- The top-level example entrypoint should live under `examples/`.
- Any Go proof that imports `runtime/internal/...` must live inside the `runtime` tree or an existing internal-test boundary that is allowed by Go internal package rules.
- The path must demonstrate the current local alpha loop rather than invent a new product API.
- The path must not create a public client SDK or imply stable client package compatibility.

## 4. Demonstrated Flow

The first future implementation should demonstrate this existing local alpha flow:

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> protected own-player storage object put/get/list/delete
-> logout
-> rejected post-logout protected request
-> selected failure-path redaction checks
```

The implementation may mention the existing realtime outbound foundation as current capability context, but it must not add stream subscriptions, chat rooms, broadcast fanout, delivery guarantees, offline inboxes, acknowledgements, retries, or new realtime behavior.

## 5. Ownership

Future implementation ownership:

```yaml
example_docs_owner: examples/local-alpha-client
example_script_owner: examples
runtime_proof_owner: runtime/internal/platform/protocol/protobuf
application_behavior_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
transport_owner: runtime/internal/platform/transport/ws
```

Rules:

- `examples/` may provide human-facing docs and shell entrypoints only.
- `examples/` must not become a product SDK or package publication source.
- Runtime behavior remains under existing runtime owners.
- Protocol payload and envelope behavior remains under existing Protobuf adapter owners.
- WebSocket transport remains credential-neutral and policy-neutral.
- No domain module should import example code.

## 6. Redaction

The example path must not print, persist, commit, or record:

- raw device credential text or bytes;
- raw access tokens;
- credential or token lookup digests;
- credential or token verifier digests;
- verifier key values;
- concrete verifier key set ids;
- HMAC inputs or outputs;
- PostgreSQL DSNs with credentials;
- headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.

Allowed public output includes route names, step names, redacted status classes, and high-level success/failure descriptions.

## 7. Nakama Mapping

Nakama reference mapping:

- The gate covers the `client_sdks_examples_and_developer_experience` capability family.
- It adopts the product pressure that a backend framework needs a clear client-facing evaluation path.
- It does not copy Nakama public REST paths, client package names, runtime helper names, wire payloads, storage APIs, session token shapes, or compatibility promises.

Pitaya reference mapping:

- Pitaya remains deferred as a future distributed architecture reference.
- This gate does not introduce frontend/backend server roles, RPC, service discovery, groups, cluster routing, or distributed session behavior.

## 8. Future Implementation Work

Open:

```text
M-154/W-0226 Implement local alpha example client path
```

The future work item may:

- add `examples/local-alpha-client/README.md`;
- add `examples/local-alpha-example-client.sh`;
- update `examples/README.md` and `examples/README.zh-CN.md`;
- optionally add or refactor a focused local alpha example-flow test under `runtime/internal/platform/protocol/protobuf`;
- update repository checks and durable memory.

The future work item must not:

- add runtime behavior;
- add new protocol routes;
- add Protobuf source or generated output;
- add migrations, persistence, repository interfaces, adapters, dependencies, startup wiring, SDK publication, generated client libraries, hosted demos, release artifacts, direct compatibility, or Pitaya-style distributed architecture.

## 9. Verification Expectations

The future implementation should verify:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.local_alpha_example_client_path_gate
node tools/vibit check change define-local-alpha-example-client-path-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

If the future implementation changes or adds a Go proof, it must also run focused Go tests and `cd runtime && go test ./...`.

## 10. Stop Conditions

Stop and create a separate gate if the example path requires:

- a public onboarding protocol route;
- public client packages or generated client libraries;
- stable SDK API guarantees;
- new Protobuf source or generated output;
- new runtime behavior;
- new authentication/session behavior;
- startup wiring or live server auto-setup;
- new dependencies;
- live PostgreSQL requirements for default verification;
- hosted deployments or release artifacts;
- chat, groups, stream subscriptions, broadcast fanout, matchmaking, match runtime, operations/admin behavior, distributed runtime, or direct Nakama/Pitaya API compatibility.
