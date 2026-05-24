# Request

## Original Request

```text
继续推进，目标nakama 推进10步以上，推进10个小时以上，不要停止，我会离开10小时。
```

## Clarified Requirement

Complete `M-152/W-0224 Select next Nakama prototype-ready capability after authenticated failure-path proof`.

The authenticated gameplay failure-path proof is now stronger. This slice must choose the next bounded Nakama-first prototype-ready capability or foundation gap without implementing runtime behavior, protocol routes, generated output, migrations, dependencies, SDKs, hosted surfaces, or direct compatibility.

## Selected Direction

```text
define_local_alpha_example_client_path_gate
```

## Selected Nakama Capability Family

```text
client_sdks_examples_and_developer_experience
```

## Selected User Requirement

```text
As a developer evaluating vibit's prototype-ready foundation, I want a clearer source-first example client or example app path that demonstrates the existing login, connection binding, protected gameplay, storage, presence, realtime-observable, logout, and failure behavior without requiring me to reverse-engineer the internal E2E tests.
```

## Rationale

The prototype-ready foundation already has several Nakama-style backend ingredients:

- local onboarding and device credential login;
- opaque access-token validation;
- first-message connection binding;
- protected inventory and presence;
- own-player storage objects;
- logout and revoked-token rejection;
- first server-push/realtime runtime and outbound protocol delivery foundation;
- authenticated failure-path proof.

The current developer-facing proof remains mostly a redacted shell wrapper over Go E2E tests. That is useful for repository confidence, but a prototype developer needs a clearer example path that reads like a client or small app flow, while still preserving vibit's source-first alpha honesty and redaction rules.

This selection deliberately chooses developer experience before broadening into chat, friends, groups, matchmaking, match runtime, operations, SDK publication, or distributed runtime. It makes the existing foundation easier to try before adding more product breadth.

## Follow-Up Work

Open:

```text
M-153/W-0225 Define local alpha example client path gate
```

The follow-up gate should define the allowed example-client shape, test expectations, redaction requirements, and stop conditions before any example implementation is added.

## Non-Goals

- Implement an example client or example app in this selection slice.
- Publish an SDK, package, binary, container, hosted demo, install script, or release artifact.
- Add new protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, authentication/session behavior, token refresh, cleanup jobs, delivery guarantees, stream subscriptions, chat rooms, groups, broadcast fanout, matchmaking, match runtime, or distributed runtime.
- Add broad operations/admin behavior.
- Add direct Nakama/Pitaya API compatibility.
- Reactivate Pitaya as a current architecture driver.

## Acceptance Criteria

- The selected next capability family is recorded as `client_sdks_examples_and_developer_experience`.
- The selected next direction is `define_local_alpha_example_client_path_gate`.
- `ADR-0132` records the selection decision and rationale.
- `M-152/W-0224` is marked completed.
- `M-153/W-0225 Define local alpha example client path gate` is opened as the next-ready work item.
- Repository checks include `runtime.next_nakama_prototype_ready_capability_selection`.
- Pitaya remains deferred as a future distributed architecture reference.
- The selection slice does not add runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted deployment, release artifacts, broad product modules, or direct compatibility.
