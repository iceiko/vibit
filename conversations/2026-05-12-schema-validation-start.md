# Conversation: Schema Validation Start

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-schema-validation/`

Related artifacts:

- `docs/schema-validation.md`
- `schema/`
- `tools/vibit`

## Context

After adding Agent Decision Records and JSON inspect commands, the agent recommended making the next step schema validation. The maintainer asked to continue.

## Maintainer Narrative

> 继续

## English Summary Of Maintainer Intent

The maintainer approved continuing with the proposed next step. No new architectural objection or redirection was introduced.

## Agent Response Summary

The agent added the first schema validation layer while avoiding external dependencies.

The implementation adds:

- Schema validation documentation
- JSON Schema files for module manifests, change specs, Agent Decision Records, and inspect output
- `node tools/vibit check schemas`
- Integration with `node tools/vibit check all`

The CLI currently validates critical fields and schema file parseability. It does not yet run full JSON Schema validation over YAML manifests.

## Decisions

- Add schema files now.
- Keep validation dependency-free for the first version.
- Validate critical fields and required ADR sections.
- Defer full JSON Schema validation until dependency strategy is decided.

## Artifacts

- `docs/schema-validation.md`
- `docs/schema-validation.zh-CN.md`
- `schema/module-manifest.schema.json`
- `schema/change-spec.schema.json`
- `schema/agent-decision-record.schema.json`
- `schema/inspect-output.schema.json`
- Updated `tools/vibit`

## Open Questions

- Which JSON Schema validator should be used later?
- Should YAML manifests be converted to JSON before validation?
- Should schemas become generator inputs?

## Follow-Up

- Add real JSON Schema validation.
- Add schema-driven examples.
- Add machine-readable output for check commands.

## Redaction Notes

No secret values are included in this log.
