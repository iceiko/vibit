# Plan

1. Read the W-0181 gate and affected runtime/module guides.
2. Extend the authentication service with local onboarding request/result vocabulary and injected dependencies.
3. Implement the gate-defined sequence using existing material generation, digest helpers, player repository, authentication repository, and unit-of-work boundaries.
4. Add startup dependency composition for the service without exposing any route or auto-onboarding behavior.
5. Add focused tests for request validation, ordering, digest-only storage, one-time result behavior after commit, failure handling, no token/session issuance, and unchanged login account-creation behavior.
6. Add ADR-0090, conversation memory, work queue, manifests, guides, check catalog, and repository check coverage.
7. Run focused Go tests, full runtime tests, repository checks, and whitespace checks.
