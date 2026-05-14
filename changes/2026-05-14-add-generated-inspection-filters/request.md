# Request

## Original Request

Continue advancing the generator and contract tooling hardening milestone.

## Clarified Requirement

Extend generated-output inspection so agents can filter generated contract shape output by module or contract type without reading the whole generated tree.

## User-Visible Outcome

Agents can run:

```bash
node tools/vibit inspect generated --module inventory --json
node tools/vibit inspect generated --type command --json
```

## Non-Goals

- Do not change generated output roots.
- Do not change generated file conventions.
- Do not turn generated contract shapes into runtime handlers.

## Acceptance Criteria

- Filtered generated inspection supports `--module`.
- Filtered generated inspection supports `--type`.
- Unfiltered generated inspection remains backward-compatible.
- Agent tooling docs and translation mention the filters.
