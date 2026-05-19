# Impact Analysis

## Affected Modules

- `runtime`: Work queue direction and runtime manifests record the session creation composition gate.
- `authentication`: Module manifests and guides record that authentication may later compose login-time session creation through application unit-of-work capabilities, but does not own runtime session persistence.

## Module Ownership Impact

No ownership changes are implemented by this direction-selection change. Future session creation composition is planned under `runtime/internal/app`. The session repository remains owned by `runtime/internal/app/session`, and authentication remains token proof validation owner.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or generated contract shape changes.

## Data And Migration Impact

No migration or data behavior changes.

## Test Impact

No runtime tests are required for this direction-selection step.

## Documentation Impact

The work queue, ADR, conversation memory, manifests, guides, rule catalog, and checks are updated to record the selection.

## Compatibility Risks

No runtime compatibility risk. The change records the selected gate direction while preserving explicit deferrals for session creation implementation, route policy, logout, reconnect, transport, protocol, and direct Nakama/Pitaya compatibility.
