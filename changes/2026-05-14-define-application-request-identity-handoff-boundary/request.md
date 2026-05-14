# Request

## Original Request

Continue advancing the project work queue. The active next-ready work item is `W-0024 Define application request identity handoff boundary`.

## Clarified Requirement

Define the application-owned request identity handoff shape in the Go runtime without implementing authentication, changing the Protobuf envelope, changing the WebSocket handshake, adding account persistence, or making inventory depend on the player module.

The handoff must make the current trust boundary explicit:

- Existing envelope `player_id`, `session_id`, `connection_id`, and `connection_epoch` values are metadata only.
- Domain handlers receive an application `RequestIdentity`.
- Validation flags remain false until a future session validation step exists.
- A future validator can replace metadata-only identity with validated application identity before domain handlers run.

## User-Visible Outcome

Agents can now see where request identity belongs in runtime code. They no longer need to infer whether `app.Session`, Protobuf `Session`, or inventory payload `player_id` should be treated as authenticated identity.

## Non-Goals

- No authentication provider.
- No token format.
- No password, guest, device, social, OAuth, OIDC, JWT, or credential model.
- No player account database schema or migrations.
- No Protobuf envelope changes.
- No WebSocket handshake changes.
- No inventory dependency on the player module.

## Unknowns

- The future authentication scheme remains intentionally undecided.
- The future runtime session storage model remains intentionally undecided.
- The future session expiration, refresh, and reconnect model remains intentionally undecided.

## Acceptance Criteria

- `runtime/internal/app` defines request identity and session validation handoff types.
- Protobuf envelope decoding produces metadata-only application identity from existing session metadata.
- Application dispatch fills metadata-only identity when a caller did not provide one.
- Application results preserve request identity metadata.
- Tests verify normalization and metadata-only validation flags.
- Change verification is recorded.
