# Impact

## User Impact

Developers and AI agents get a clearer source-first path for the current local alpha capability loop:

```text
examples/local-alpha-client/README.md
examples/local-alpha-example-client.sh
```

The existing `examples/local-alpha-request-loop.sh` remains as a compatibility wrapper.

## Architecture Impact

This slice is developer-experience work over existing runtime and protocol proofs. It does not change runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, persistence, startup wiring, authentication/session behavior, or transport behavior.

## Reference Impact

Nakama remains the primary product reference. This slice covers the `client_sdks_examples_and_developer_experience` capability family by making the existing alpha loop easier to evaluate without publishing an SDK.

Pitaya remains deferred as a future distributed architecture reference. This slice does not introduce frontend/backend roles, RPC, service discovery, groups, cluster routing, or distributed sessions.

## Risk

Primary risk is over-claiming. The example path must stay explicit that it is source-first, repository-local, and not a public client SDK or compatibility surface.

