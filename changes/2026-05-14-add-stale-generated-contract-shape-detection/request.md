# Request

## Original Request

Continue advancing the generator and contract tooling hardening milestone.

## Clarified Requirement

Strengthen generated output checks so stale generated Go contract shape files fail verification when they no longer correspond to registered non-runtime semantic contracts.

## User-Visible Outcome

`node tools/vibit check generated --json` rejects generated contract shape files that exist under `runtime/internal/generated/contracts/` but are not expected from `.arch/contracts.yaml`.

## Non-Goals

- Do not change generated output roots.
- Do not change generated file naming conventions.
- Do not delete generated files unless a failing stale file is actually present.
- Do not turn generated contract shapes into runtime handlers.

## Acceptance Criteria

- Expected generated contract shape paths are derived from registered non-runtime semantic contracts.
- Committed generated contract shape files outside that expected set are reported as stale drift.
- Existing source trace and reproducibility checks remain in place.
- `node tools/vibit check generated --json` passes for the current generated output tree.
