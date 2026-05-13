# Verification

Verified:

- `node tools/vibit check work --json`
- `node tools/vibit inspect work --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check change add-work-item-continuation-system --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for unignored files:

```bash
git ls-files --cached --others --exclude-standard -z | node -e '...'
```

Not verified:

- None.

Not applicable:

- Runtime Go tests are not applicable because this change adds workflow metadata and CLI checks only.
- Protobuf generation is not applicable because this change does not add or change `.proto` sources.
- Database migration verification is not applicable because no persistence adapter or migration is added.
