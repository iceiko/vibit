# Conversation: Check Contracts

Date: 2026-05-12
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-12-add-check-contracts/`

Related artifacts:

- `tools/vibit`
- `.arch/contracts.yaml`
- `contracts/inventory/`
- `rules/check-rules.json`

## Context

The first inventory contract source files had been added and manually verified. The next useful step was to make that verification repeatable for future agents.

## Maintainer Narrative

The maintainer asked:

```text
继续下一步工作。
```

## Agent Response Summary

The agent chose to add a first-class contract check instead of proceeding directly to runtime implementation.

The check validates that `.arch/contracts.yaml` exists, declares YAML as the contract source format, references ADR-0007, and points to existing source files for registered commands, queries, events, errors, and permissions. It also validates minimal fields inside each contract source file.

## Decisions

- Add `node tools/vibit check contracts`.
- Include the new check in `node tools/vibit check all`.
- Keep the implementation dependency-free for now.

## Artifacts

- Updated `tools/vibit`.
- Updated `rules/check-rules.json`.
- Updated README and AGENTS command references.

## Open Questions

- Whether to add a formal YAML parser later.
- Whether to add a dedicated contract schema after the source format stabilizes.

## Follow-Up

- Use `check contracts` before generating TypeScript runtime artifacts from contract source files.
- Consider adding `inspect contract` after the first generator needs structured contract intake.

## Redaction Notes

No secrets, tokens, credentials, private account details, or unrelated private data are recorded in this log.
