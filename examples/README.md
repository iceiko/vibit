# vibit Examples

Status: Draft v0.1
Last updated: 2026-05-22

The paired Simplified Chinese translation is `examples/README.zh-CN.md`. The English file is authoritative.

This directory contains source-first local examples and templates. These files are not product SDKs, release artifacts, hosted demos, install scripts, package registry publications, or direct Nakama/Pitaya compatibility surfaces.

## Local Alpha Request Loop

Run from the repository root:

```bash
examples/local-alpha-request-loop.sh
```

The script wraps the focused authenticated gameplay E2E proof:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

It proves:

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

The script is intentionally redacted. It must not print raw credentials, raw access tokens, verifier keys, DSNs, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.

## Local Environment Template

`examples/local.prototype.env.example` is a placeholder checklist for local configuration fields. It contains no real secrets and is not meant to be sourced directly without replacing every placeholder with local-only values.

Private local env files such as `.vibit.local.env`, `.env.local`, and `.env.*.local` are ignored by the repository `.gitignore`. Do not commit private local configuration.
