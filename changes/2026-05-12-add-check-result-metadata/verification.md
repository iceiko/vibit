# Verification

Verified:

- `node tools/vibit check all`
- `node tools/vibit check all --json` parsed with `JSON.parse`
- `node tools/vibit check all --json` metadata assertion for `rule_id` and `artifact`
- `node tools/vibit check schemas --json` metadata assertion for `rule_id` and `artifact`
- `node tools/vibit check architecture --json` metadata assertion for `rule_id` and `artifact`
- `node tools/vibit check change add-check-result-metadata --json` metadata assertion for `rule_id` and `artifact`
- `node tools/vibit check module inventory --json` metadata assertion for `rule_id` and `artifact`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`
- `git diff --check`

Not verified:

- Full JSON Schema validation is not implemented yet; current schema checks validate critical fields only.

Not applicable:

- Runtime server tests do not apply because this change targets CLI tooling.
