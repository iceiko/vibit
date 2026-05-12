# Rule Catalogs

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: Machine-readable rule metadata

This directory stores rule catalogs for machine-readable vibit tool output.

Initial catalog:

```text
rules/check-rules.json
```

The check rule catalog maps `rule_id` values from `node tools/vibit check ... --json` to human-readable metadata.

Each rule should declare:

- Stable `rule_id`
- Category
- Default severity
- Title
- Description
- Agent guidance

Rules are not final public API yet. They are versioned as `0.1` while the standards stabilize.
