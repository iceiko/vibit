# Impact Analysis

## Affected Modules

- `inventory`: its `module.yaml` is validated by schema checks.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds one CLI command:

```text
vibit check schemas
```

Changes:

```text
vibit check all
```

to include schema checks.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification should include:

- `node tools/vibit check schemas`
- `node tools/vibit check all`
- Inspect commands
- Secret scan

## Documentation Impact

Adds schema validation documentation and initial JSON Schema files.

Updates root docs and agent guide to reference schema validation.

## Compatibility Risks

Medium. Early schema checks may become too strict if the standards are still evolving. Keep first checks focused on critical fields and enum values.
