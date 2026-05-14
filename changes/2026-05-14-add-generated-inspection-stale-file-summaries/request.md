# Request

## Original Request

Continue advancing the generator and contract tooling hardening milestone.

## Clarified Requirement

Expose stale generated contract shape file summaries in generated output inspection, using the same expected-output model as `check generated`.

## User-Visible Outcome

`node tools/vibit inspect generated --json` includes `stale_contract_shape_files` and a stale count in `summary`.

## Non-Goals

- Do not delete stale generated files automatically.
- Do not change generated output roots.
- Do not replace the generated inspection output schema.
- Do not change generated contract shape semantics.

## Acceptance Criteria

- Generated inspection includes stale generated contract shape file paths.
- Generated inspection includes stale generated contract shape counts.
- Existing module and type filters continue to work.
- Current generated output inspection reports zero stale files.
