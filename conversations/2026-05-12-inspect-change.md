# Conversation: Inspect Change

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-inspect-change/`

Related artifacts:

- `tools/vibit`
- `schema/inspect-output.schema.json`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`

## Context

After several early CLI improvements, the maintainer redirected the agent to focus on project development rather than repeatedly discussing GitHub token handling.

## Maintainer Narrative

The maintainer said:

> 你不需要总是考虑GitHub的token的问题，专注于当前项目的开发。

## Agent Response Summary

The agent continued with an agent-native tooling slice: adding `node tools/vibit inspect change <change-id>` so agents can read change spec metadata and file presence as JSON during intake or handoff.

## Decisions

- Add direct change inspection instead of requiring agents to manually open every file in `changes/<date>-<change-id>/`.
- Missing change IDs should return structured JSON with `exists: false`, following the existing module inspection pattern.
- Keep YAML parsing lightweight for now and avoid introducing dependencies.

## Artifacts

- Updated `tools/vibit`.
- Updated `schema/inspect-output.schema.json`.
- Updated README and AGENTS bilingual documentation.
- Created change spec `changes/2026-05-12-add-inspect-change/`.

## Open Questions

- Should missing inspect targets eventually use a shared structured error envelope?
- Should `inspect change` include short markdown excerpts in future versions?

## Follow-Up

- Continue adding atomic inspection commands that reduce how much context agents must read manually.

## Redaction Notes

No secrets, tokens, account details, or private data were recorded in this log.
