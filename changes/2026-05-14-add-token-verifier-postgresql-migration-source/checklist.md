# Checklist

- [x] Add SQL-first token verifier migration source.
- [x] Include goose Up and Down markers.
- [x] Include `-- Module: runtime.authentication`.
- [x] Create only `authentication_access_tokens`.
- [x] Store non-plaintext lookup and verifier digests only.
- [x] Include player and credential linkage.
- [x] Include expiration, revocation, rotation lineage, retention, and cleanup eligibility columns.
- [x] Preserve player account lifecycle table separation.
- [x] Avoid raw token, refresh token, session, WebSocket, external identity, provider, generic metadata, JWT, OAuth, OIDC, Redis-like, password-hashing, and major authentication dependency scope.
- [x] Update manifests, standards, guides, checks, and conversation memory.
- [x] Run verification.
