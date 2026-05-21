# Verification

Status: Verified

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change prepare-first-alpha-feedback-intake-surfaces --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Known accepted warning:

```text
runtime.identity_boundary
```

Notes:

- No Go runtime tests were required for this documentation/check-rule slice because no Go runtime files were changed.
- No live PostgreSQL verification was required or run.
- No release, hosted deployment, package, binary, container, checksum, signing/provenance artifact, install script, registry publication, public announcement, or paid promotion command was run.
