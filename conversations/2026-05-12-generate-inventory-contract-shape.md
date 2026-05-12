# Conversation: Generate Inventory Contract Shape

Date: 2026-05-12
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-12-generate-inventory-contract-shape/`

Related artifacts:

- `tools/vibit`
- `contracts/inventory/commands/GrantItem.yaml`
- `modules/inventory/module.yaml`
- `modules/inventory/generated/contracts/GrantItem.generated.ts`

## Context

The project had added contract source files, contract consistency checks, and a contract inspection command. The next step was to begin the first generated shape without rushing into runtime business logic.

## Maintainer Narrative

The maintainer asked:

```text
继续
```

## Agent Response Summary

The agent chose to generate a single TypeScript contract shape for the inventory `GrantItem` command. This advances the proof chain from contract source to generated shape while keeping runtime behavior, transport, persistence, and package manager choices deferred.

## Decisions

- Add `node tools/vibit generate contract --module <module> --type <type> --id <id>`.
- Generate only `GrantItem.generated.ts` in this change.
- Declare the generated file in `modules/inventory/module.yaml`.
- Add `node tools/vibit check generated` before generated output grows.

## Artifacts

- Updated `tools/vibit`.
- Updated `rules/check-rules.json`.
- Updated `modules/inventory/module.yaml`.
- Added `modules/inventory/generated/contracts/GrantItem.generated.ts`.
- Updated README and AGENTS command references.

## Open Questions

- Whether future generated TypeScript should include runtime validators.
- Whether generated files should remain module-local or move under repository-level `generated/`.
- Whether a real YAML parser should be introduced before generating all contract types.

## Follow-Up

- Use the generated shape as the input boundary before adding handwritten `GrantItem` command handler logic.

## Redaction Notes

No secrets, tokens, credentials, private account details, or unrelated private data are recorded in this log.
