# Impact

## User Impact

Maintainers and contributors get a single alpha acceptance checklist that explains what is ready, what remains manual, and what is deferred before `v0.1 alpha` can be packaged and published.

## Runtime Impact

No runtime behavior changes are made. The checklist references existing runtime surfaces and proofs:

- `/v1/ws`
- `/healthz`
- `/readyz`
- `/version`
- `/configz`
- `examples/local-alpha-request-loop.sh`
- the focused authenticated gameplay E2E test

## Documentation Impact

Adds:

- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `ADR-0096`
- W-0188 change and conversation records

Updates the continuation queue and architecture manifests so the next direction becomes packaging the alpha developer flow without release publishing.

## Security And Redaction Impact

The checklist reinforces that verifier keys, raw device credentials, raw access tokens, DSNs, digests, headers, cookies, query strings, subprotocol values, remote addresses, and concrete transport metadata are not log-safe and should not be recorded or committed.

No new secret carrier, credential path, token path, or authentication/session behavior is introduced.

## Out Of Scope

- Publishing `v0.1 alpha`.
- Release tags, binaries, archives, containers, packages, or hosted deployments.
- Runtime behavior changes.
- Protocol route changes.
- Protobuf sources or generated output.
- Migrations or dependencies.
- Broad operations/admin behavior.
- Direct Nakama/Pitaya API compatibility.
