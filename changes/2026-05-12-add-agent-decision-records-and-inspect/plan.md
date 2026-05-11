# Plan

## Files To Create

- `docs/agent-decision-record.md`
- `docs/agent-decision-record.zh-CN.md`
- `decisions/README.md`
- `decisions/README.zh-CN.md`
- `decisions/_template/adr-agent.md`
- `decisions/ADR-0001-agent-native-maintainability.md`
- `conversations/2026-05-12-external-ai-feedback-on-traceability-tooling-and-immutability.md`

## Files To Edit

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/conventions.yaml`
- `tools/vibit`
- `changes/2026-05-12-add-agent-decision-records-and-inspect/*`

## Generated Artifacts

- None.

## Handwritten Logic

- Add inspect command parsing.
- Read module manifests with lightweight text extraction.
- Return deterministic JSON for module and boundary inspection.
- Include generated-file immutability metadata in docs and inspect output where available.

## Tests

- CLI aggregate check
- Inspect module command
- Inspect boundary command
- Secret scan

## Verification Commands

- `node tools/vibit check all`
- `node tools/vibit inspect module inventory`
- `node tools/vibit inspect boundary --from inventory --to player`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`

## Rollback Or Migration Notes

Remove the ADR-Agent docs, decisions directory, and inspect commands if this standard proves premature.
