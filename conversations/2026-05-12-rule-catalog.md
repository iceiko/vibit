# Conversation: Rule Catalog

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-rule-catalog/`

Related artifacts:

- `rules/README.md`
- `rules/README.zh-CN.md`
- `rules/check-rules.json`
- `schema/rule-catalog.schema.json`
- `tools/vibit`
- `.arch/conventions.yaml`

## Context

The previous change added `rule_id` and `artifact` metadata to JSON check results. The next step was to make those rule IDs discoverable without reading the CLI source.

## Maintainer Narrative

The maintainer asked to continue:

> 继续

This continued the direction of turning implicit agent context into explicit, machine-readable project artifacts.

## Agent Response Summary

The agent added a first rule catalog at `rules/check-rules.json`, with schema support at `schema/rule-catalog.schema.json`. The CLI schema check now verifies the catalog shape, duplicate rule IDs, required metadata fields, and coverage for all rule IDs currently known to `tools/vibit`.

## Decisions

- Check result `rule_id` values should be documented in a machine-readable catalog.
- The initial catalog is `rules/check-rules.json`.
- Rule catalogs are versioned as `0.1` and remain draft while standards stabilize.
- Coverage for known CLI rule IDs is enforced by `node tools/vibit check schemas`.

## Artifacts

- Added `rules/README.md` and `rules/README.zh-CN.md`.
- Added `rules/check-rules.json`.
- Added `schema/rule-catalog.schema.json`.
- Updated `.arch/conventions.yaml`.
- Updated README, AGENTS, and schema validation docs.
- Updated `tools/vibit` schema checks.

## Open Questions

- Whether rule catalogs should eventually be generated from source.
- Whether rule IDs should have explicit version suffixes after schema `0.1`.
- Whether catalogs should split by tool, standard, or module.

## Follow-Up

- Add rule catalogs for future tools as those tools appear.
- Consider central generation once rule IDs become stable public API.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
