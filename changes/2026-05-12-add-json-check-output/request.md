# Request

## Original Request

The maintainer asked to continue after the recommendation to add JSON output for check commands.

Original maintainer statement:

> 继续

## Clarified Requirement

Add machine-readable JSON output for check commands while preserving existing human-readable text output.

Initial target commands:

```bash
node tools/vibit check all --json
node tools/vibit check schemas --json
node tools/vibit check architecture --json
node tools/vibit check change <change-id> --json
node tools/vibit check module <module> --json
```

## User-Visible Outcome

Humans can keep using the current text output.

Agents can use JSON output to identify:

- Which check ran
- Which rules passed or failed
- Severity
- Message
- Artifact path where relevant
- Summary counts

## Non-Goals

- Do not change existing default text output.
- Do not introduce external dependencies.
- Do not implement full JSON Schema validation yet.
- Do not add CI yet.

## Unknowns

- Final stable JSON schema for check output.
- Whether future checks should support JSON Lines for streaming.
- Whether warnings should be configurable as errors.

## Acceptance Criteria

- [x] `check architecture --json` returns valid JSON.
- [x] `check schemas --json` returns valid JSON.
- [x] `check change <change-id> --json` returns valid JSON.
- [x] `check module <module> --json` returns valid JSON.
- [x] `check all --json` returns valid JSON with nested subchecks.
- [x] Existing text output still works.
- [x] README and AGENTS mention JSON output.
- [x] `schema/inspect-output.schema.json` or a new schema covers check output.
- [x] Verification is recorded.
