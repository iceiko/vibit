# Conversation: Check Memory

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-check-memory/`

Related artifacts:

- `tools/vibit`
- `rules/check-rules.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Context

After adding `inspect memory`, the next step was to make project memory verifiable through the standard CLI checks.

## Maintainer Narrative

The maintainer said:

> 你不需要总是考虑GitHub的token的问题，专注于当前项目的开发。

The maintainer's earlier instruction also established conversation history as durable project memory that should help future humans and agents understand how vibit reached its current direction.

## Agent Response Summary

The agent added `node tools/vibit check memory`, registered the new rule IDs it may emit, and added the memory check to `check all`.

## Decisions

- Verify project memory structure separately from schema checks.
- Enforce conversation and ADR structural requirements through `check all`.
- Keep the check structural and avoid judging the semantic content of maintainer narrative.

## Artifacts

- Updated `tools/vibit`.
- Updated `rules/check-rules.json`.
- Updated README and AGENTS bilingual documentation.
- Created change spec `changes/2026-05-12-add-check-memory/`.

## Open Questions

- Should memory checks later validate cross-links between conversations, decisions, and changes?
- Should related artifacts be required for every conversation?

## Follow-Up

- Consider a link-integrity check after the project has more memory artifacts.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
