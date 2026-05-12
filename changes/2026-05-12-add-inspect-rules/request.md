# Request

## Original Request

The maintainer asked the agent to push the existing work and continue development, while asking the agent to raise questions when confirmation is needed.

## Clarified Requirement

Add an agent-friendly command for listing rule catalog entries, with optional category filtering.

Target commands:

```bash
node tools/vibit inspect rules
node tools/vibit inspect rules --category check
```

## User-Visible Outcome

Agents can discover available check rules without parsing `rules/check-rules.json` directly.

## Non-Goals

- Do not implement fuzzy search.
- Do not add external dependencies.
- Do not change `rules/check-rules.json` shape.
- Do not add pagination yet.

## Unknowns

- Whether future rule discovery should support text search.
- Whether rule categories should become their own catalog.

## Acceptance Criteria

- [x] `node tools/vibit inspect rules` returns valid JSON.
- [x] `node tools/vibit inspect rules --category check` returns valid JSON.
- [x] Category filtering returns only matching rules.
- [x] The inspect output schema covers rules inspection.
- [x] README and AGENTS mention the command.
- [x] Verification is recorded.
