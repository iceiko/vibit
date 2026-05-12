# Schema Validation Standard

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: `schema/`

This document defines the first schema validation layer for vibit.

Schema validation turns standards from readable documents into machine-checkable contracts.

## 1. Purpose

vibit uses manifests, change specs, decision records, and tool JSON output as agent-readable architecture context.

Agents should not rely only on prose or visual inspection. Important structures should be validated by tools.

## 2. Location

Schema files live under:

```text
schema/
```

Initial files:

```text
schema/module-manifest.schema.json
schema/change-spec.schema.json
schema/agent-decision-record.schema.json
schema/inspect-output.schema.json
schema/rule-catalog.schema.json
```

## 3. Validation Strategy

The first implementation is intentionally lightweight.

Current rules:

- Keep schemas as JSON Schema draft 2020-12 documents.
- Validate that schema files exist and parse as JSON.
- Validate critical fields in manifests and specs with no external dependencies.
- Keep strict validation small until standards stabilize.

Future rules:

- Use a full JSON Schema validator when dependency strategy is decided.
- Generate examples from schemas.
- Use schemas as generator inputs where practical.

## 4. CLI Commands

Initial command:

```bash
node tools/vibit check schemas
```

Aggregate command:

```bash
node tools/vibit check all
```

`check all` should include schema checks.

## 5. What Is Checked Now

Initial checks should cover:

- Required schema files exist.
- Schema files parse as JSON.
- `modules/<module>/module.yaml` declares critical fields.
- Change specs declare critical fields and allowed verification status.
- Agent Decision Records contain required sections.
- Tool JSON output schemas cover inspect output and check result output.
- Rule catalogs declare required metadata and cover known check result `rule_id` values.
- Architecture conventions declare schema artifacts.

This is not yet full YAML schema validation.

## 6. Agent Rules

Agents must update schemas when changing the shape of:

- Module manifests
- Change specs
- Agent Decision Records
- Tool JSON output, including inspect output and check result output
- Rule catalogs

If the CLI cannot fully validate a structure yet, agents must record the gap in the relevant change spec.

## 7. Open Questions

- Which full JSON Schema validator should be used later?
- Should YAML manifests be converted to JSON before validation?
- Should schemas drive code generation directly?
- Should tool JSON output have stable versioned schemas from the beginning?
