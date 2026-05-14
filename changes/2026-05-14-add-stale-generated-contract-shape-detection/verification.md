# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check generated --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit check change add-stale-generated-contract-shape-detection --json`

Not verified:

- `node tools/vibit check work --json`

Not applicable:

- Runtime tests are not required for this generated-output check change.
