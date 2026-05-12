# Conversation: Inspect Rules

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-inspect-rules/`

Related artifacts:

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Context

After adding single-rule inspection, the next step was to make the rule catalog discoverable through the CLI without requiring agents to parse the catalog file directly.

## Maintainer Narrative

The maintainer asked the agent to continue development and to ask for confirmation when needed.

## Agent Response Summary

The agent added `node tools/vibit inspect rules` and `node tools/vibit inspect rules --category <category>`. The command reads `rules/check-rules.json`, returns the matching rules as JSON, and includes filter metadata.

## Decisions

- Add a lightweight discovery command before adding search or pagination.
- Support category filtering because categories already exist in the rule catalog.
- Keep output JSON-only, consistent with other inspect commands.

## Artifacts

- Updated `tools/vibit`.
- Updated `schema/inspect-output.schema.json`.
- Updated README and AGENTS bilingual documentation.
- Created change spec `changes/2026-05-12-add-inspect-rules/`.

## Open Questions

- Should rule discovery later support text search?
- Should categories become a separate catalog if they gain semantics beyond grouping?

## Follow-Up

- Consider `inspect rules --search <term>` only after the catalog grows enough to need it.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
