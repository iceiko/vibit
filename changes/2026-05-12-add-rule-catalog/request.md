# Request

## Original Request

The maintainer asked to continue:

> 继续

## Clarified Requirement

Create an explicit rule catalog for machine-readable `rule_id` values emitted by vibit checks.

The first catalog should document existing CLI check rule IDs and make them discoverable to agents.

## User-Visible Outcome

Agents can map a check result `rule_id` to:

- A human-readable title
- A category
- A default severity
- A short description
- Guidance for what to inspect next

## Non-Goals

- Do not stabilize rule IDs as a final public API yet.
- Do not add external validation dependencies.
- Do not implement full JSON Schema validation.
- Do not create a web documentation site.

## Unknowns

- Whether future rule catalogs should be split by tool, module, or standard.
- Whether rule IDs should be generated from code instead of maintained as data.
- Whether rule IDs should include explicit version suffixes after `0.1`.

## Acceptance Criteria

- [x] `rules/check-rules.json` exists.
- [x] `schema/rule-catalog.schema.json` exists.
- [x] The catalog includes all current rule IDs emitted by `tools/vibit`.
- [x] `node tools/vibit check schemas` validates the catalog shape and coverage.
- [x] README and AGENTS mention the rule catalog.
- [x] English and Simplified Chinese schema validation docs mention the rule catalog.
- [x] Verification is recorded.
