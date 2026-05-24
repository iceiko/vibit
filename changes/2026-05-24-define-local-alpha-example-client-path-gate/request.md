# Request

## Original Request

```text
继续推进，目标nakama 推进10步以上，推进10个小时以上，不要停止，我会离开10小时。
```

## Clarified Requirement

Complete `M-153/W-0225 Define local alpha example client path gate`.

The previous selection chose `client_sdks_examples_and_developer_experience` as the next Nakama-first prototype-ready capability family. This slice must define the allowed source-first example client or example app path before implementation, including ownership, redaction, accepted existing runtime/protocol surfaces, verification expectations, and stop conditions.

## User-Visible Outcome

A future agent can implement a clearer local alpha example client path without guessing whether it is allowed to publish SDKs, generate client libraries, add protocol routes, change runtime behavior, or introduce hosted/demo/release surfaces.

## Selected Direction

```text
implement_local_alpha_example_client_path
```

as the follow-up direction after this gate.

## Selected Nakama Capability Family

```text
client_sdks_examples_and_developer_experience
```

## Non-Goals

- Implement an example client or example app in this gate slice.
- Publish an SDK, generated client library, package, binary, container, hosted demo, install script, or release artifact.
- Add new protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, authentication/session behavior, delivery guarantees, stream subscriptions, chat rooms, groups, broadcast fanout, matchmaking, match runtime, operations/admin behavior, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.
- Reactivate Pitaya as a current architecture driver.

## Acceptance Criteria

- `docs/local-alpha-example-client-path-gate.md` and its Simplified Chinese translation define the gate.
- `ADR-0133` records the decision.
- The gate selects a source-first local alpha example path, not a public SDK.
- The gate records future file candidates under `examples/` and existing runtime proof ownership.
- Redaction rules cover raw credentials, raw access tokens, digests, verifier keys, DSNs, and transport metadata.
- Existing accepted routes and behavior surfaces are listed.
- `M-153/W-0225` is marked completed.
- `M-154/W-0226 Implement local alpha example client path` is opened as next-ready.
- Repository checks include `runtime.local_alpha_example_client_path_gate`.
- Runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted deployment, release artifacts, broad product modules, distributed runtime, and direct compatibility remain deferred.
