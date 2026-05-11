# Request

## Original Request

The maintainer asked to continue after the recommendation to add schema validation.

Original maintainer statement:

> 继续

## Clarified Requirement

Add the first schema validation layer for vibit standards.

This should include:

- JSON Schema files for key artifacts
- Bilingual schema validation documentation
- Initial CLI checks that validate critical fields without adding external dependencies
- `check all` integration

## User-Visible Outcome

Maintainers and agents can run:

```bash
node tools/vibit check all
```

and receive validation that key manifests and standards artifacts are structurally coherent.

## Non-Goals

- Do not introduce npm dependencies yet.
- Do not implement full JSON Schema validation yet.
- Do not replace future dedicated schema validation tooling.
- Do not choose server implementation language or server instance architecture.

## Unknowns

- Whether full validation should use AJV, another JSON Schema validator, or a custom small validator later.
- How much YAML parsing should remain custom versus dependency-backed.
- Whether schemas should become the source for generators.

## Acceptance Criteria

- [x] Add `docs/schema-validation.md` and Chinese translation.
- [x] Add schema directory and JSON Schema files.
- [x] Register schema artifacts in `.arch/conventions.yaml`.
- [x] Add `node tools/vibit check schemas`.
- [x] Include schema checks in `node tools/vibit check all`.
- [x] Validate module manifest critical fields.
- [x] Validate change spec critical fields.
- [x] Validate ADR-Agent critical sections.
- [x] Update README and AGENTS.
- [x] Add conversation log.
- [x] Run verification.
