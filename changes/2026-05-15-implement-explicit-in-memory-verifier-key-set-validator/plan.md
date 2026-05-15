# Plan

1. Add the application-owned authentication package.
2. Implement `VerifierKeySetConfig` and validated `VerifierKeySet`.
3. Add redacted validation errors.
4. Add focused unit tests for validation, copying, immutability, and redaction.
5. Adjust runtime boundary checks only as needed to allow the ADR-0045-authorized validator while preserving deferrals.
6. Update manifests, agent guides, and change records.
7. Run Go and repository verification.
8. Record verification results.
