# Request

## Original Request

Implement source-first feature request scaffolding so AI agents can turn user backend requirements into bounded specs, acceptance criteria, test plans, implementation boundaries, verification records, and durable memory before coding.

## Clarified Requirement

Complete `M-157/W-0229 Implement agent-native feature request scaffolding` by adding a repository-local feature request template, a narrow `tools/vibit scaffold feature` command, documentation, focused repository checks, and durable records. The implementation must stay inside docs, templates, tooling, checks, manifests, ADRs, and change memory.

## User-Visible Outcome

An agent or contributor can run:

```bash
node tools/vibit scaffold feature <change-id> --request "<original request>" --summary "<summary>"
```

and receive a `changes/YYYY-MM-DD-<change-id>/` directory containing `request.md`, `spec.yaml`, `impact.md`, `plan.md`, `checklist.md`, and `verification.md` with Nakama mapping, Pitaya deferral, acceptance criteria, test planning, boundaries, redaction, verification, and durable memory prompts.

## Nakama Capability Mapping

- Capability family: `agent_native_requirement_test_implementation_workflow`.
- Product intent: This slice supports Nakama-first product growth by making every future Nakama-style capability start from a bounded requirement, acceptance criteria, and test plan before implementation.
- API compatibility: This slice does not authorize direct Nakama API compatibility.

## Pitaya Status

Pitaya remains a deferred future distributed architecture reference. This slice does not introduce frontend/backend roles, RPC, service discovery, groups, cluster routing, or distributed session behavior.

## Non-Goals

- Add runtime behavior.
- Add protocol routes or Protobuf source.
- Add generated output.
- Add migrations, persistence, repositories, or adapters.
- Add dependencies.
- Add startup wiring.
- Add SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya compatibility.

## Unknowns

- None blocking. Future user feedback may require richer scaffold arguments, but W-0229 only needs the first deterministic source-first shape.

## Acceptance Criteria

- [x] A feature-request template directory exists at `changes/_template/feature-request/`.
- [x] The template covers `request.md`, `spec.yaml`, `impact.md`, `plan.md`, `checklist.md`, and `verification.md`.
- [x] Template content includes Nakama capability mapping, Pitaya deferral, acceptance criteria, test planning, implementation boundaries, redaction reminders, verification commands, and durable memory expectations.
- [x] `tools/vibit scaffold feature` creates the six required artifacts from the template and refuses to overwrite existing change directories.
- [x] Repository checks validate the implementation shape.
- [x] Runtime, protocol, generated output, migration, dependency, SDK, hosted, distributed runtime, and direct compatibility scope remain deferred.

## Test Expectations

- [x] `node -c tools/vibit` verifies tool syntax.
- [x] A scaffold dry run verifies the non-writing command path.
- [x] A real scaffold invocation generated this change directory.
- [x] `node tools/vibit check runtime --json` validates the rule and template shape.
- [x] Full repository checks are run before completion.

## Redaction Notes

The scaffold and this change do not record raw device credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private user data beyond explicit request text.
