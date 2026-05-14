# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect contracts --type command --json`
- `node tools/vibit inspect contracts --status draft --json`
- `node tools/vibit inspect contracts --module inventory --type query --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check change add-contract-inspection-filters --json`
- `node tools/vibit check work --json`

Not verified:

- None.

Not applicable:

- Runtime tests are not required for this inspection-only change.
