# Request

Implement `W-0107`, the authentication service behavior skeleton.

The change must add only a skeleton service shape under `runtime/internal/app/authentication`, with typed dependencies, request/result vocabulary, redacted failure classes, and fail-closed behavior.

The change must not implement real device credential login, access-token validation, token issuance, repository lookup or mutation behavior, logout execution, refresh execution, cleanup jobs, protocol carriers, startup wiring, repository changes, migrations, dependencies, or production authentication behavior.
