# Request

## Original Request

The maintainer asked the agent to stop repeatedly centering GitHub token concerns and focus on current project development:

> 你不需要总是考虑GitHub的token的问题，专注于当前项目的开发。

## Clarified Requirement

Continue development with a small agent-native tooling slice: add a CLI command for inspecting a change spec directory as structured JSON.

Target command:

```bash
node tools/vibit inspect change <change-id>
```

## User-Visible Outcome

Agents and maintainers can query a change spec by ID and receive a stable JSON summary of the change directory, declared metadata, affected modules, verification state, and required file presence.

## Non-Goals

- Do not introduce external dependencies.
- Do not implement full YAML parsing.
- Do not redesign the change spec standard.
- Do not add structured JSON errors for missing or malformed arguments.
- Do not select the server runtime language or instance architecture in this change.

## Unknowns

- Whether future inspect commands should share a single structured error format.
- Whether `inspect change` should later include excerpts from markdown files.
- Whether change specs should eventually be validated with full JSON Schema-compatible YAML parsing.

## Acceptance Criteria

- [x] `node tools/vibit inspect change add-inspect-change` returns valid JSON.
- [x] The output includes required file presence for the change spec directory.
- [x] The output includes core `spec.yaml` metadata when the spec exists.
- [x] The output returns a structured `exists: false` result for a missing change ID.
- [x] `schema/inspect-output.schema.json` includes `change_inspection`.
- [x] README and AGENTS mention the command in English and Simplified Chinese.
- [x] Verification is recorded.
