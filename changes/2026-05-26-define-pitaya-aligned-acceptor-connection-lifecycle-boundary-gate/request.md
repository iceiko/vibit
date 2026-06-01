# Request

Define the `W-0265` Pitaya-aligned acceptor and connection lifecycle boundary gate selected by `ADR-0172`.

The change must remain gate-only. It must not implement acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, session binding behavior, kick/disconnect behavior, protocol routes or messages, generated output, persistence, dependencies, metrics/tracing behavior, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility.
