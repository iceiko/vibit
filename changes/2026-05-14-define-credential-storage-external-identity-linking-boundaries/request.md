# Request

## Original Request

The maintainer asked to continue the next project step, while preserving ask-first boundaries and timely commit/push behavior.

## Clarified Requirement

Advance `W-0060`: define credential storage and external identity linking boundaries without adding schema, dependencies, credential records, provider subjects, password hashing, OAuth, OIDC, runtime lookup code, account linking handlers, recovery flows, merge behavior, or player account lifecycle table changes.

## User-Visible Outcome

Agents gain a durable standard that separates player account lifecycle storage from future credential storage and future external identity links. Nakama authentication coverage is mapped into deferred vibit login-method families, and Pitaya session binding vocabulary remains a reference without governing vibit's API.

## Non-Goals

- Do not choose a supported login method.
- Do not add credential tables or external identity tables.
- Do not add password hashing, cryptography, OAuth, OIDC, provider SDKs, or provider dependencies.
- Do not add runtime credential lookup, account linking handlers, recovery flows, or merge behavior.
- Do not add token/session storage or runtime authentication implementation.
- Do not change player account lifecycle table shape.
- Do not change the Protobuf envelope or WebSocket handshake behavior.
- Do not copy Nakama or Pitaya public APIs.

## Acceptance Criteria

- Add an English standard and Simplified Chinese translation for credential storage and external identity linking boundaries.
- Preserve player account lifecycle tables as lifecycle storage only.
- Map Nakama authentication method coverage into deferred vibit login-method families.
- Map Pitaya session binding vocabulary without changing vibit's API surface.
- Update manifests, guides, and module metadata so future agents can discover the boundary.
- Mark `W-0060` completed and `W-0061` next-ready.
- Run repository verification and record results.
