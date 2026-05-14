# Verification

Verified:

- `node tools/vibit inspect work --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change close-m007-tooling-hardening-milestone --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Runtime Go tests were already run in W-0047 and are not required for this workflow-only closeout.
- Live PostgreSQL verification is not required because this change does not alter persistence behavior.
