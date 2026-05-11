# Verification

Verified:

- `rg -n "ghp_|github_pat_|TOKEN|Token:|api[_-]?key|password|secret" .`
  - Result: matched only redaction rules, template language, and the documented scan command itself. No raw GitHub token value was found in repository files.
- `git status --short`
  - Result: expected documentation and standard files are modified or untracked before commit.
- `git diff --stat`
  - Result: root documentation and architecture convention updates are present.

Not verified:

- No runtime behavior was verified because no implementation code exists for this change.

Not applicable:

- Runtime tests are not applicable because this is a documentation and standards change.
