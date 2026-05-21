# Impact

## User Impact

Developers get small local HTTP endpoints for startup troubleshooting:

- `/healthz`
- `/readyz`
- `/version`
- `/configz`

These endpoints help confirm that the runtime process is alive, ready to accept local traffic, exposing `/v1/ws`, using the expected store posture, and running the pre-alpha runtime version.

## Runtime Impact

The runtime HTTP mux now mounts read-only status endpoints in addition to `/v1/ws`.

No gameplay behavior, authentication behavior, session behavior, startup configuration semantics, Protobuf contracts, generated output, migrations, or dependencies are changed.

## Security And Redaction Impact

`/configz` reports only redacted configuration posture:

- runtime store name,
- WebSocket path,
- local alpha request-loop script path,
- whether PostgreSQL configuration is present,
- whether authentication configuration is required,
- and that secrets are redacted.

It does not expose raw credential values, raw access tokens, verifier key values, DSNs, digests, headers, cookies, query strings, subprotocols, remote addresses, or concrete transport metadata.

## Out Of Scope

- Production observability.
- Metrics backend integration.
- Admin console behavior.
- Broad operations surface.
- Release publishing.
- Direct Nakama/Pitaya API compatibility.
