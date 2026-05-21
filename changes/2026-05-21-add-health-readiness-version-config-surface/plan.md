# Plan

1. Add small status endpoint constants and redacted runtime info in `runtime/cmd/vibit-server`.
2. Register `/healthz`, `/readyz`, `/version`, and `/configz` on the existing HTTP mux.
3. Add tests for endpoint responses and redaction of DSN, verifier key, token audience, and transport-like values.
4. Update runbook, README, work queue, architecture manifests, ADR, and conversation log.
5. Add repository check coverage for `runtime.health_readiness_version_config_surface`.
6. Run focused Go tests and repository verification.
