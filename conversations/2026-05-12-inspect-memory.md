# Conversation: Inspect Memory

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-inspect-memory/`

Related artifacts:

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Context

The maintainer asked the agent to focus on project development. After adding direct change inspection, the next agent-native tooling step was to make project memory discoverable as a structured index.

## Maintainer Narrative

The maintainer said:

> 你不需要总是考虑GitHub的token的问题，专注于当前项目的开发。

Earlier, the maintainer also established that conversations with agents should be recorded so future maintainers and agents can understand how the project arrived at its current direction.

## Agent Response Summary

The agent added `node tools/vibit inspect memory`, which lists change specs, conversation logs, and Agent Decision Records as JSON.

## Decisions

- Add a machine-readable project memory index instead of requiring agents to scan directories manually.
- Keep Markdown parsing shallow and dependency-free.
- Use the existing change inspection shape for change entries.

## Artifacts

- Updated `tools/vibit`.
- Updated `schema/inspect-output.schema.json`.
- Updated README and AGENTS bilingual documentation.
- Created change spec `changes/2026-05-12-add-inspect-memory/`.

## Open Questions

- Should `inspect memory` later support filters?
- Should conversation logs eventually use machine-readable front matter?
- Should project memory move into a formal artifact registry?

## Follow-Up

- Continue adding atomic inspect commands that reduce agent intake cost.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
