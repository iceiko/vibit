# Plan

1. Read the W-0184 work item, ADR-0091, runtime guides, and relevant runtime/protocol/application code.
2. Define the W-0184 boundary in the change spec before treating implementation as authorized.
3. Add a focused Go E2E test over existing runtime capabilities:
   - direct local onboarding service call,
   - existing device credential login protocol route,
   - existing first-message connection binding route,
   - existing protected inventory route through `AuthenticatedRequest`,
   - existing protected presence query through `AuthenticatedRequest`,
   - existing logout route,
   - post-logout protected-route rejection.
4. Keep all helpers test-local and avoid new production protocol, storage, migration, dependency, or generated-output surface.
5. Record ADR-0092 and the conversation log.
6. Update architecture manifests, guides, README/alpha pointers, and repository check rules.
7. Run focused Go tests, full Go tests, repository checks, and whitespace checks.
8. Mark the W-0184 change verified and advance the queue to the runtime runbook refresh work item.
