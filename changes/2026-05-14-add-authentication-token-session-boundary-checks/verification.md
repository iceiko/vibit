# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit inspect rule runtime.authentication_token_session_boundary --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check architecture --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-authentication-token-session-boundary-checks --json`
- `node tools/vibit check all --json`
- `node tools/vibit inspect next --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- New Go tests are not required because this change adds repository checks, manifests, and documentation only. Existing Go tests are run by `node tools/vibit check runtime --json`.
