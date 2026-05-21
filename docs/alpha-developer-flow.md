# Alpha Developer Flow

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Packaged local developer journey for vibit's v0.1 alpha path

The paired Simplified Chinese translation is `docs/alpha-developer-flow.zh-CN.md`. The English file is authoritative.

This document packages the existing local alpha entry points into one developer journey. It is not a release declaration and does not authorize release publishing, release packaging, hosted deployment, runtime behavior changes, protocol changes, generated output changes, migrations, dependencies, broad operations/admin behavior, product module expansion, or direct Nakama/Pitaya API compatibility.

## 1. Purpose

The local alpha path now has the pieces a technically capable contributor needs to inspect vibit:

- project positioning in `README.md`,
- runtime startup and verification notes in `docs/runtime-runbook.md`,
- a redacted request-loop script at `examples/local-alpha-request-loop.sh`,
- local status endpoints at `/healthz`, `/readyz`, `/version`, and `/configz`,
- acceptance criteria in `docs/alpha-acceptance-checklist.md`,
- and continuation state in `.arch/work-items.yaml`.

This document connects those pieces so a contributor can follow the same sequence without hunting through project memory.

## 2. Current Package State

```text
local_alpha_developer_flow_packaged: true
release_declared: false
release_publishing_authorized_by_this_document: false
release_packaging_authorized_by_this_document: false
release_publishing_decision_gate_defined: true
release_execution_preparation_gate_defined: true
release_execution_authorization_gate_defined: true
release_execution_maintainer_decision_recorded: true
next_direction: release_identifier_artifact_plan
next_work_status: next_ready
```

The repository remains pre-alpha. The packaged flow is ready for local review, and the maintainer decision now allows the next release-path planning step. Publishing `v0.1 alpha`, selecting the final release identifier, creating tags or artifacts, and executing release commands remain deferred to later explicit work items.

## 3. Recommended Journey

1. Read `README.md` to understand vibit, its pre-alpha state, and the agent-native server framework goal.
2. Read `docs/v0.1-alpha-goal.md` for the short-term `v0.1 alpha` target and long-term Nakama/Pitaya-class direction.
3. Read `docs/alpha-acceptance-checklist.md` to see which alpha items are ready, manual, deferred, or blocked.
4. Install local prerequisites: Go, Node.js, PostgreSQL when testing the persistent runtime path, and Buf/Protobuf tooling only when regenerating Protobuf output.
5. Run static repository checks:

```bash
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check all --json
```

6. Run Go tests:

```bash
cd runtime
go test ./...
```

7. Run the redacted local alpha request loop from the repository root:

```bash
examples/local-alpha-request-loop.sh
```

8. Use `docs/runtime-runbook.md` when starting the server manually or evaluating the PostgreSQL runtime path.
9. Use `.arch/work-items.yaml` and `node tools/vibit inspect next` to find the next bounded contribution.

## 4. Runtime Entry Points

The current process exposes:

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

`/v1/ws` is the gameplay WebSocket endpoint. It expects binary `vibit.protocol.v1.Envelope` Protobuf bytes. JSON is not accepted on this endpoint.

`/healthz`, `/readyz`, `/version`, and `/configz` are local alpha troubleshooting endpoints. `/configz` reports redacted runtime posture and includes `secrets_redacted: true`. These endpoints are not a production operations API, admin console, metrics backend, gameplay protocol route, release artifact, or hosted deployment surface.

## 5. Local Proof Flow

The packaged local proof flow is:

```text
local onboarding
-> device credential login
-> connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

The executable entry point is:

```bash
examples/local-alpha-request-loop.sh
```

The script wraps the focused Go E2E proof:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

The proof uses existing runtime protocol handlers. It does not require live PostgreSQL, committed verifier keys, raw credentials, raw access tokens, DSNs, digests, or a hand-built WebSocket client.

## 6. PostgreSQL Path

The PostgreSQL runtime path is the current alpha runtime composition, but it has manual setup requirements:

- prepare a local PostgreSQL database,
- apply or verify SQL migrations explicitly,
- set `VIBIT_RUNTIME_STORE=postgres`,
- set `VIBIT_POSTGRES_DSN`,
- set all authentication verifier key environment variables,
- avoid committing local verifier keys or DSNs.

Normal server startup does not apply migrations automatically. Optional live PostgreSQL verification remains opt-in through `VIBIT_POSTGRES_TEST_DSN` and a disposable database.

## 7. Redaction Contract

Do not record or commit:

- raw device credential text or bytes,
- raw access tokens,
- credential or token lookup digests,
- credential or token verifier digests,
- HMAC input or output bytes,
- verifier key values,
- concrete verifier key set ids,
- PostgreSQL DSNs with credentials,
- headers, cookies, query strings, WebSocket subprotocol values, or remote addresses that may carry secrets.

The request-loop script and `/configz` surface are part of this redaction posture.

## 8. Contribution Entry Point

The next contribution path is always machine-readable:

```bash
node tools/vibit inspect next
```

The release execution maintainer decision is now recorded in `docs/release-execution-maintainer-decision.md`. The next work is `W-0194 Define release identifier and artifact plan`; it may plan the identifier and artifact/publication surface but must not create tags, artifacts, hosted deployments, GitHub release records, or release publication commands.

## 9. Deferred Work

The following remain deferred until later explicit work items:

- publishing `v0.1 alpha`,
- selecting release identifiers,
- creating release tags, binaries, archives, containers, packages, checksums, provenance files, or hosted deployments,
- adding a public local onboarding protocol route,
- adding production signup, external identity providers, password login, account recovery, account merge, or multi-device linking,
- adding broad operations/admin behavior, metrics backend integration, or production observability,
- adding chat, friends, groups, parties, matchmaking, match runtime, SDKs, distributed runtime, or direct Nakama/Pitaya API compatibility.

## 10. Verification

Use this command set when checking the packaged flow:

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
