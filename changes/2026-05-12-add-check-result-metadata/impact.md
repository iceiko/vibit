# Impact Analysis

## Affected Modules

- `inventory`: used by `check module inventory --json`.

## Module Ownership Impact

No module ownership changes.

## Public Contract Impact

Changes the machine-readable output contract for `vibit check ... --json` by adding `rule_id` and `artifact` fields to result items.

Default human-readable text output should remain unchanged.

## Data And Migration Impact

No runtime data or migration impact.

## Test Impact

Verification must confirm that JSON output parses and every result item has the new metadata fields.

## Documentation Impact

Update:

- README
- AGENTS
- Tool JSON output schema
- Change spec
- Conversation log

## Compatibility Risks

Low to medium. The JSON output contract is still version `0.1`, but downstream agents may already rely on the initial `results` shape.

This change is additive for specific checks. For `check all --json`, the aggregate output is intentionally compact and includes subcheck summaries instead of full nested result details; agents should run the specific failing check with `--json` for detailed result items.
