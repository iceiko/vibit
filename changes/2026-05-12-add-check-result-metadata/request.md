# Request

## Original Request

The maintainer asked to continue:

> 继续

## Clarified Requirement

Make machine-readable check output more useful to agents by adding stable result metadata to every check result item.

Initial metadata fields:

- `rule_id`: stable identifier for the checked rule
- `artifact`: repository artifact or logical target related to the result, or `null` when no single artifact applies

## User-Visible Outcome

Agents can consume `node tools/vibit check ... --json` output without parsing human-readable messages to identify which rule passed or failed and which artifact was involved.

Humans can keep using the existing text output.

## Non-Goals

- Do not introduce external dependencies.
- Do not implement a full JSON Schema validator.
- Do not redesign the complete check output contract.
- Do not change server runtime behavior.

## Unknowns

- Final naming convention for long-term rule IDs.
- Whether future rule IDs should be centrally registered in a separate catalog.
- Whether result metadata should include line numbers once checks become more precise.

## Acceptance Criteria

- [x] Every `check ... --json` result item includes `rule_id`.
- [x] Every `check ... --json` result item includes `artifact`.
- [x] `check all --json` includes compact result items for nested subchecks.
- [x] Existing text check output remains usable.
- [x] The tool JSON output schema documents the new fields.
- [x] README and AGENTS mention `rule_id` and `artifact`.
- [x] Verification is recorded.
