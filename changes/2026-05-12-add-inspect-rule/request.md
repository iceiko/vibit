# Request

## Original Request

The maintainer provided a GitHub token, asked the agent to submit and push the current work, then continue development:

> 提交并推送，然后继续进行开发。有什么需要我确认的，请及时询问和讨论。

## Clarified Requirement

After pushing existing commits, add an agent-friendly command for inspecting rule catalog entries by `rule_id`.

Target command:

```bash
node tools/vibit inspect rule <rule-id>
```

## User-Visible Outcome

Agents can query a single rule by ID and receive JSON containing the rule metadata from `rules/check-rules.json`.

## Non-Goals

- Do not add external dependencies.
- Do not redesign rule catalog shape.
- Do not implement fuzzy search.
- Do not add a web UI.

## Unknowns

- Whether future inspect commands should support category listing.
- Whether missing rules should return structured JSON errors.

## Acceptance Criteria

- [x] `node tools/vibit inspect rule check.subcheck` returns valid JSON.
- [x] The output includes rule metadata from `rules/check-rules.json`.
- [x] Missing rule IDs fail clearly.
- [x] The inspect output schema covers rule inspection.
- [x] README and AGENTS mention the command.
- [x] Verification is recorded.
