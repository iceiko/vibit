# Request

## Original Request

The maintainer asked the agent to focus on current project development:

> 你不需要总是考虑GitHub的token的问题，专注于当前项目的开发。

## Clarified Requirement

Continue agent-native tooling by adding a CLI command that gives agents a machine-readable index of project memory.

Target command:

```bash
node tools/vibit inspect memory
```

## User-Visible Outcome

Agents can list change specs, conversation logs, and Agent Decision Records as JSON before deciding which artifacts need full reading.

## Non-Goals

- Do not add external dependencies.
- Do not parse every Markdown section deeply.
- Do not implement search, ranking, or semantic retrieval.
- Do not replace `inspect change`.
- Do not select the server runtime language or instance architecture in this change.

## Unknowns

- Whether memory inspection should later support filters such as `--type conversation`.
- Whether conversation logs should eventually have a machine-readable front matter block.
- Whether decisions and conversations should receive formal schemas beyond required Markdown sections.

## Acceptance Criteria

- [x] `node tools/vibit inspect memory` returns valid JSON.
- [x] The output includes counts for changes, conversations, and decisions.
- [x] The output includes change spec summaries.
- [x] The output includes conversation log IDs, titles, dates, paths, related changes, and related artifacts.
- [x] The output includes decision IDs, titles, statuses, dates, paths, related changes, related conversations, and related artifacts.
- [x] `schema/inspect-output.schema.json` includes `memory_inspection`.
- [x] README and AGENTS mention the command in English and Simplified Chinese.
- [x] Verification is recorded.
