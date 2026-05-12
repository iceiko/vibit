# Conversation: Check Result Metadata

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-check-result-metadata/`

Related artifacts:

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Context

After adding JSON output for check commands, the next step was to make the output easier for agents to consume reliably.

## Maintainer Narrative

The maintainer asked to continue:

> 继续

This continued the tooling direction around atomic, machine-readable probes for agent workflows.

## Agent Response Summary

The agent added stable metadata to check result items: `rule_id` identifies the checked rule, and `artifact` identifies the related file, directory, command, or logical target. During verification, `check all --json` proved too large when it embedded all nested result details, so it was changed to return a compact overview. Agents should run the specific failing check with `--json` to retrieve detailed result items.

## Decisions

- JSON check result items should include `rule_id` and `artifact`.
- `check all --json` is a compact overview, not a full nested detail dump.
- Specific check commands keep detailed `results` arrays for agent diagnosis.
- The first metadata shape remains versioned as schema `0.1`.

## Artifacts

- Updated `tools/vibit` check JSON output.
- Updated `schema/inspect-output.schema.json`.
- Updated README and AGENTS bilingual documentation.
- Created change spec `changes/2026-05-12-add-check-result-metadata/`.

## Open Questions

- Whether `rule_id` values should be centrally registered in a catalog.
- Whether future results should include line numbers or structured locations.
- Whether `check all --json` should offer an explicit verbose mode later.

## Follow-Up

- Add a dedicated rule catalog if rule IDs become stable public API.
- Add structured location metadata when checks become more precise.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
