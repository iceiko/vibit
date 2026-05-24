# Impact

## Affected Modules

This pilot is repository workflow and roadmap scope. It selects a future runtime proof-hardening slice but does not change runtime behavior.

Relevant future runtime areas for `W-0222`:

- `runtime/internal/app/presence`
- `runtime/internal/app/connection`
- `runtime/internal/platform/protocol/protobuf`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-request-loop.sh`

## Public Behavior Impact

No public behavior changes in this pilot.

The selected future behavior to harden is existing self-presence/status:

- online after authenticated connection binding;
- offline after connection close or invalidation;
- protected self-only query over server-owned presence snapshot.

## Contract And Protocol Impact

No contract or protocol changes in this pilot.

`W-0222` should use existing presence route and payloads unless its own change spec proves that a broader boundary is required. It should not add new Protobuf source or generated output unless a later bounded work item explicitly authorizes it.

## Data And Migration Impact

No data or migration changes.

Presence/status remains single-process and non-durable for this path. Durable/distributed presence remains deferred.

## Test Impact

This pilot records the test plan. It does not add Go tests because it is a planning and workflow-pilot slice.

`W-0222` should add or strengthen tests around:

- active bound connection produces online self-presence;
- closed connection is excluded from active presence;
- invalidated connection is excluded from active presence;
- authenticated local alpha flow can observe the bounded behavior through existing protocol surfaces where feasible;
- raw credentials, raw tokens, verifier data, transport internals, and direct compatibility shapes are not exposed.

## Documentation And Memory Impact

Updates required:

- change spec for the pilot;
- ADR-0129;
- conversation memory;
- `.arch/` work/runtime/reference/contracts/conventions/modules manifests;
- public continuation docs and guides;
- rule catalog and `tools/vibit` check logic.

## Compatibility Risks

Direct Nakama/Pitaya API compatibility remains out of scope.

The pilot uses Nakama only as a product capability reference. It does not copy Nakama routes, payloads, data models, runtime API names, or compatibility shims.

Pitaya remains deferred and does not authorize cluster/RPC/frontend-backend/service-discovery work.
