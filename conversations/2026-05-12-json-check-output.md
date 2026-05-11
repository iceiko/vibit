# Conversation: JSON Check Output

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-json-check-output/`

Related artifacts:

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Context

The project had already introduced agent-readable inspect commands. The next step was to make verification commands easier for agents to consume without parsing human text output.

## Maintainer Narrative

The maintainer asked to continue:

> 继续

This followed earlier discussion about atomic tooling and machine-readable probes for agents.

## Agent Response Summary

The agent added `--json` support to check commands while preserving default text output for humans. The implementation kept the CLI dependency-free and used the existing `schema/inspect-output.schema.json` file for the first check result output shape.

## Decisions

- Check commands should support `--json` for machine-readable intake, verification, and handoff.
- Text output remains the default behavior.
- `check all --json` should include nested subcheck results instead of forcing agents to run every subcheck manually.
- The first schema is intentionally small and versioned as `0.1`.

## Artifacts

- Added JSON output support to `tools/vibit` check commands.
- Extended `schema/inspect-output.schema.json` with `check_result`.
- Updated README and AGENTS bilingual documentation.
- Created change spec `changes/2026-05-12-add-json-check-output/`.

## Open Questions

- Whether future check output should support JSON Lines for streaming.
- Whether warnings should become configurable failures.
- Whether check result schemas should split from inspect output once the contract stabilizes.

## Follow-Up

- Add stronger schema validation when dependency strategy is decided.
- Consider structured `artifact`, `path`, and `rule_id` fields in check results.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
