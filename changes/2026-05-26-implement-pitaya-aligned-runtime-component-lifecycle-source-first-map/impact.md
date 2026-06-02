# Impact: Pitaya-Aligned Runtime Component Lifecycle Source-First Map

## Scope

This change adds a source-first inspection surface and repository rule for Pitaya-aligned runtime component lifecycle vocabulary.

Affected areas:

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/*`
- `README*.md`
- `AGENTS*.md`
- `runtime/AGENTS*.md`
- `modules/*/AGENTS*.md`
- `modules/*/module.yaml`
- continuation and roadmap docs
- ADR/change/conversation records

## Behavior

Runtime behavior is unchanged.

The new command is:

- `node tools/vibit inspect pitaya-component-lifecycle --json`

The new check rule is:

- `runtime.pitaya_aligned_runtime_component_lifecycle_source_first_map`

## Compatibility

No protocol shape, generated output, repository interface, PostgreSQL adapter, migration, dependency, hosted surface, SDK, release artifact, distributed runtime behavior, or direct Nakama/Pitaya API compatibility is added.
