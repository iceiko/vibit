# Verification

Verified:

- `node tools/vibit check all`
- `node tools/vibit check all --json` parsed with `JSON.parse`
- `node tools/vibit check schemas --json` parsed with `JSON.parse`
- `node tools/vibit check architecture --json` parsed with `JSON.parse`
- `node tools/vibit check change add-json-check-output --json` parsed with `JSON.parse`
- `node tools/vibit check module inventory --json` parsed with `JSON.parse`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`

Not verified:

- Full JSON Schema validation is not implemented yet; current schema checks validate critical fields only.

Not applicable:

- Runtime server tests do not apply because this change targets CLI tooling.
