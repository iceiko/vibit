# Conversation: Inspect Contract

Date: 2026-05-12
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-12-add-inspect-contract/`

Related artifacts:

- `tools/vibit`
- `.arch/contracts.yaml`
- `contracts/inventory/`
- `modules/inventory/module.yaml`
- `schema/inspect-output.schema.json`

## Context

The project had just added `check contracts`, which verifies the whole contract registry and all registered contract source files.

The next useful step was to add a more granular intake command so an agent can inspect one contract without manually reading the registry, source file, and module manifest.

## Maintainer Narrative

The maintainer asked:

```text
继续下一步工作。
```

## Agent Response Summary

The agent chose to add `inspect contract` before moving to runtime generation. This keeps the project focused on agent-native intake and verification instead of rushing into implementation code.

## Decisions

- Add `node tools/vibit inspect contract --module <module> --type <type> --id <id>`.
- Keep the command dependency-free for now.
- Treat this command as a read-only intake probe, not a replacement for `check contracts`.

## Artifacts

- Updated `tools/vibit`.
- Updated `schema/inspect-output.schema.json`.
- Updated README and AGENTS command references.

## Open Questions

- Whether future versions should add `inspect contracts --module <module>` for listing all module contracts.
- Whether full contract payload parsing should wait for a real YAML parser.

## Follow-Up

- Use `inspect contract` during generator intake before emitting TypeScript runtime artifacts.

## Redaction Notes

No secrets, tokens, credentials, private account details, or unrelated private data are recorded in this log.
