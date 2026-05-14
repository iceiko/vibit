# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit --help`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect contracts --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect reference --json`
- `node tools/vibit check agent-tooling --json`

Not verified:

- None.

Not applicable:

- Runtime behavior verification is not required because this change only adds agent tooling and standards.
