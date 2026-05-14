# Request

## Original Request

Continue advancing the generator and contract tooling hardening milestone.

## Clarified Requirement

Add a compact generated contract shape status to generated output inspection so agents can quickly tell whether the requested generated-output slice is complete, missing output, or has stale files.

## User-Visible Outcome

`node tools/vibit inspect generated --json` includes `contract_shape_status`, `missing_contract_shapes`, and a matching `summary.contract_shape_status`.

## Non-Goals

- Do not weaken `check generated`.
- Do not replace detailed generated inspection fields.
- Do not change generated output roots or generated file conventions.

## Acceptance Criteria

- Generated inspection reports `contract_shape_status`.
- Generated inspection reports missing contract shape objects.
- Existing stale file summaries remain available.
- Current generated output reports `contract_shape_status: complete`.
