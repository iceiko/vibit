# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect generated --module inventory --json`
- `node tools/vibit inspect generated --type command --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check change add-generated-inspection-stale-file-summaries --json`
- `node tools/vibit check work --json`

Not verified:

- None.

Not applicable:

- Runtime tests are not required for this inspection-only change.
