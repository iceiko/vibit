# Request

## Original Request

The maintainer asked the agent to continue ten steps according to professional judgment, stopping only when a truly necessary decision required human input.

## Clarified Requirement

Advance `W-0065` by comparing first login-method candidates and recommending the first login-method set for W-0066.

## User-Visible Outcome

The repository gains a candidate comparison document with a clear recommendation: ratify `device_credential_login` first and defer guest/anonymous, custom ID, email/password, external provider, and service-auth families.

## Non-Goals

- Implementing any login method.
- Adding credential tables, external identity tables, token tables, session tables, or migrations.
- Adding password hashing, JWT, OAuth, OIDC, provider SDK, cryptography, key-management, Redis-like, or other major dependencies.
- Changing Protobuf envelope behavior.
- Changing WebSocket handshake authentication behavior.
- Adding runtime player account handlers or WebSocket routes.

## Unknowns

- Whether device credentials are client-generated, server-issued, or both.
- Whether account creation is implicit or explicit.
- Whether credential rotation and recovery are in the first implementation.
- Which token model will be selected later.

## Acceptance Criteria

- Compare all candidate login families listed in W-0065.
- Record benefits, risks, artifact gates, contract impact, storage impact, dependency impact, and implementation complexity.
- Recommend the first login-method set for W-0066.
- Preserve all implementation boundaries.
