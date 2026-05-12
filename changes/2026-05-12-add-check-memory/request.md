# Request

## Original Request

The maintainer asked the agent to focus on current project development:

> 你不需要总是考虑GitHub的token的问题，专注于当前项目的开发。

## Clarified Requirement

Continue improving agent-native project memory by adding a verification command for conversation logs and Agent Decision Records.

Target command:

```bash
node tools/vibit check memory
```

## User-Visible Outcome

Agents can verify that project memory artifacts have the required headings, metadata, and section structure. `check all` includes this verification by default.

## Non-Goals

- Do not perform semantic validation of conversation content.
- Do not parse hidden reasoning or private chain-of-thought.
- Do not introduce external Markdown parsing dependencies.
- Do not redesign conversation or ADR standards.

## Unknowns

- Whether memory checks should later validate links between conversations, changes, and decisions.
- Whether conversation logs should eventually have structured front matter.
- Whether missing related changes should be warnings or errors.

## Acceptance Criteria

- [x] `node tools/vibit check memory` passes.
- [x] `node tools/vibit check memory --json` returns valid check JSON.
- [x] `node tools/vibit check all` includes `check memory`.
- [x] New memory-related `rule_id` values are registered in `rules/check-rules.json`.
- [x] README and AGENTS mention the command in English and Simplified Chinese.
- [x] Verification is recorded.
