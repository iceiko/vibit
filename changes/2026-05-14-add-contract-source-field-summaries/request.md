# Request

## Original Request

Continue advancing the generator and contract tooling hardening milestone.

## Clarified Requirement

Expose machine-readable contract source field summaries in contract index inspection so agents can understand input, output, and event payload shapes without opening every YAML source file.

## User-Visible Outcome

`node tools/vibit inspect contracts --type command --json` includes field summaries such as `input_required`, `input_properties`, `output_required`, and `output_properties`.

## Non-Goals

- Do not change contract source format.
- Do not make summaries authoritative over contract source files.
- Do not change Protobuf schemas or generated output.
- Do not implement runtime behavior.

## Acceptance Criteria

- Contract inspection includes input field summaries for command and query sources.
- Contract inspection includes output field summaries for command and query sources.
- Event contract inspection includes payload field summaries.
- Existing contract inspection filters remain compatible.
