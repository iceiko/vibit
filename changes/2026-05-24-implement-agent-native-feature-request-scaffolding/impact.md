# Impact Analysis

## Nakama Product Capability Impact

This implements the `agent_native_requirement_test_implementation_workflow` capability family. It does not add a new game backend feature. It makes future Nakama-style capability work start from scaffolded requirements, acceptance criteria, test planning, boundaries, verification, and durable memory.

## Pitaya Impact

Pitaya remains deferred as a future distributed architecture reference. This change does not add distributed topology, frontend/backend split, RPC, service discovery, groups, cluster routing, distributed sessions, or Pitaya API compatibility.

## Affected Modules

- `repository_workflow`: adds a feature-request template and command.
- `developer_experience`: makes the user requirement intake path executable.
- `runtime` and `reference`: update checks and manifests only.

## Module Ownership Impact

No domain module ownership changes. Storage remains outside this repository-wide workflow slice.

## Public Contract Impact

Adds a repository tooling command:

```text
tools/vibit scaffold feature <change-id> --request <text> [--summary <text>] [--date YYYY-MM-DD] [--dry-run]
```

No server command, query, event, error, permission, route, protocol payload, or schema contract is added.

## Data And Migration Impact

No data, persistence, migration, index, repository, adapter, or transaction behavior changes.

## Test Impact

Verification covers:

- tool syntax;
- scaffold dry-run behavior;
- real scaffold generation through this change directory;
- rule inspection;
- change/work/runtime/memory/schema/all checks;
- diff whitespace check.

Go runtime tests are not required for this slice because it changes no Go runtime behavior. `node tools/vibit check runtime --json` still runs the repository's Go tests as part of the runtime check.

## Documentation And Memory Impact

Adds:

- `docs/agent-native-feature-request-scaffolding.md`
- `docs/agent-native-feature-request-scaffolding.zh-CN.md`
- `ADR-0137`
- W-0229 change artifacts
- W-0229 conversation memory

Updates manifests, agent guides, roadmap docs, `tools/vibit`, and `rules/check-rules.json`.

## Compatibility Risks

The only new user-visible surface is a source-first repository tooling command. The command refuses to overwrite an existing change directory and supports `--dry-run`. There is no API, event, data, protocol, SDK, hosted, distributed runtime, or direct Nakama/Pitaya compatibility impact.
