# Request

## Original Request

```text
继续推进，注意提交和推送。提交和推送的key在Git忽略的文件里有，你找一下。
```

## Clarified Requirement

Complete `M-171/W-0243 Select next Nakama prototype-ready capability after friends relationship route proof`.

The friends relationship protocol route family is now locally proven through the authenticated alpha request flow. This slice must choose the next bounded Nakama-first prototype-ready capability or foundation gap without implementing runtime behavior, protocol routes, generated output, migrations, dependencies, SDKs, hosted surfaces, operations/admin endpoints, or direct compatibility.

## Selected Direction

```text
define_minimum_operations_inspection_surface_gate
```

## Selected Nakama Capability Family

```text
admin_console_metrics_observability_and_operations
```

## Selected User Requirement

```text
As a developer evaluating or extending vibit's prototype-ready foundation, I want a minimum source-first operations inspection posture that tells me what local alpha server state can be inspected, how redaction is handled, and which operational surfaces remain future work before I add broader social, realtime, competitive, or match runtime behavior.
```

## Rationale

The prototype-ready foundation already has several Nakama-style backend ingredients:

- local onboarding and device credential login;
- opaque access-token validation, logout, and revoked-token rejection;
- first-message connection binding and protected route policy;
- protected inventory, presence, storage objects, realtime outbound delivery, and friends routes;
- a source-first example client path;
- feature-request scaffolding and durable repository checks;
- local proof for the friends relationship route family.

The current alpha has `/healthz`, `/readyz`, `/version`, and `/configz`, but it does not yet define what a minimum operations inspection surface should cover for prototype authors who need to understand local server state. Adding more groups, parties, chat, matchmaking, or match runtime breadth before this gate would make the system harder to inspect and debug.

## Follow-Up Work

Open:

```text
M-172/W-0244 Define minimum operations inspection surface gate
```

The follow-up gate should define the accepted source-first operations inspection posture, redaction requirements, inspectable state categories, test expectations, and stop conditions before any operations/admin implementation is added.

## Non-Goals

- Implement an operations inspection surface, admin console, metrics endpoint, observability pipeline, or dashboard in this selection slice.
- Publish an SDK, package, binary, container, hosted demo, install script, or release artifact.
- Add new protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, repository interfaces, PostgreSQL adapters, authentication/session behavior, token refresh, cleanup jobs, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, or distributed runtime.
- Add direct Nakama/Pitaya API compatibility.
- Reactivate Pitaya as a current architecture driver.

## Acceptance Criteria

- The selected next capability family is recorded as `admin_console_metrics_observability_and_operations`.
- The selected next direction is `define_minimum_operations_inspection_surface_gate`.
- `ADR-0151` records the selection decision and rationale.
- `M-171/W-0243` is marked completed.
- `M-172/W-0244 Define minimum operations inspection surface gate` is opened as the next-ready work item.
- Repository checks include `runtime.next_nakama_prototype_ready_capability_after_friends_route_proof`.
- Pitaya remains deferred as a future distributed architecture reference.
- The selection slice does not add runtime behavior, operations/admin endpoints, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted deployment, release artifacts, broad product modules, or direct compatibility.
