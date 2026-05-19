# Impact

This change records the maintainer-authorized direction after `W-0156`.

It selects a narrow implementation slice for `LogoutAccessToken` because the behavior gate already constrained the execution sequence and boundaries. The implementation remains application-owned under `runtime/internal/app/authentication` and remains limited to the presented access-token record.

No runtime session revocation, active WebSocket close behavior, connection registry, reconnect or epoch behavior, Protobuf logout route, protocol session carrier, refresh, logout-all, admin revocation, cleanup job, dependency, memory durable session behavior, broader game backend module, or direct Nakama/Pitaya API compatibility is selected by this confirmation.

Nakama informs the token/session lifecycle pressure. Pitaya informs keeping session/connection infrastructure separate from handlers and credential parsing.
