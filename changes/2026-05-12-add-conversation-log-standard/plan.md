# Plan

## Files To Create

- `docs/conversation-log.md`
- `docs/conversation-log.zh-CN.md`
- `conversations/README.md`
- `conversations/README.zh-CN.md`
- `conversations/_template/session.md`
- `conversations/2026-05-12-founding-session.md`

## Files To Edit

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `.arch/conventions.yaml`

## Generated Artifacts

- None.

## Handwritten Logic

- None. Documentation only.

## Tests

- No implementation tests.

## Verification Commands

- `rg -n "ghp_|github_pat_|TOKEN|Token:" .`
- `git status --short`
- `git diff --stat`

## Rollback Or Migration Notes

This change can be reverted as a documentation-only change if the project chooses a different memory system later.
