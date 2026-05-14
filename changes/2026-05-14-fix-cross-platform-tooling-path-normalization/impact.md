# Impact

## Affected Modules

No domain module is affected.

This is a repository tooling bugfix.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or game protocol contract changes.

## Data And Migration Impact

No data or migration impact.

## Test Impact

No Go runtime tests are added. Existing runtime tests still run through `node tools/vibit check runtime --json`.

## Documentation Impact

Updated:

- `docs/agent-tooling.md`
- `docs/agent-tooling.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/conventions.yaml`

## Compatibility Risks

The JSON path normalization is intended to reduce platform-specific differences. Existing consumers expecting backslash paths from Windows tooling should switch to forward-slash repository paths.

This is not a runtime API compatibility change.
