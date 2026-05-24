# Checklist

- [x] Capture the original user request.
- [x] Map the slice to the Nakama identity/auth/session capability family.
- [x] Record acceptance criteria and non-goals.
- [x] Add local alpha E2E failure-path tests using existing runtime/protocol surfaces.
- [x] Cover missing wrapper, malformed wrapper, malformed token, invalid token, expired token, revoked token, protected presence failure, and redaction.
- [x] Preserve production runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, authentication/session behavior, token refresh, cleanup jobs, and direct compatibility deferrals.
- [x] Add `ADR-0131`.
- [x] Add conversation memory.
- [x] Register `runtime.authenticated_gameplay_failure_path_verification`.
- [x] Update work queue from W-0223 to W-0224.
- [x] Update relevant manifests, guides, README, alpha docs, maturity docs, and roadmap docs.
- [x] Run focused Go tests.
- [x] Run full Go tests.
- [x] Run repository checks.
- [x] Record verification.
