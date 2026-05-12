# Verification

Verified:

- `node tools/vibit check schemas`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`
- `git diff --check`

Not verified:

- Full JSON Schema validation is not implemented yet; current schema checks validate critical fields and the rule catalog shape with dependency-free logic.

Not applicable:

- Runtime server tests do not apply because this change targets standards and CLI tooling.
