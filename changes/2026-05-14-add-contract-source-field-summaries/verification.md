# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect contracts --type command --json`
- `node tools/vibit inspect contracts --type event --json`
- `node tools/vibit inspect contracts --module inventory --type query --json`
- `node tools/vibit check change add-contract-source-field-summaries --json`
- `node tools/vibit check work --json`
- `node tools/vibit check agent-tooling --json`

Not verified:

- None.

Not applicable:

- Runtime tests are not required for this inspection-only change.
