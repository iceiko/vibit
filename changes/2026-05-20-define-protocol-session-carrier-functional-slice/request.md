# Request

## Original Request

The maintainer asked to continue advancing toward Nakama/Pitaya-class server capability while keeping confirmation gates lighter so the project can ship.

## Clarified Requirement

Advance `M-104/W-0176` as a Tier 2 functional slice under `ADR-0082`. The slice must embed the protocol session carrier boundary in the change spec and implement only the smallest protocol-visible session metadata behavior needed for lifecycle closure.

## User-Visible Outcome

Successful device-credential login responses now carry the server-created runtime session id and validated player id in the existing `Envelope.Session` metadata. The response keeps the existing envelope shape and does not add a new Protobuf message or field.

For application results that already contain a validated request identity, the Protobuf response builder can also derive response session metadata from that validated identity. Metadata-only identity remains metadata-only and is not upgraded into proof.

## Non-Goals

- No new Protobuf source field.
- No generated Protobuf output.
- No reconnect token.
- No resume token.
- No WebSocket handshake authentication.
- No durable or distributed session routing.
- No close code mapping.
- No player-visible close reason text.
- No logout-triggered socket close.
- No runtime session revocation.
- No presence lifecycle behavior.
- No operations/admin disconnect.
- No dependencies.
- No direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] The change spec embeds the protocol session carrier boundary and records Tier 2 gate density.
- [x] Successful login responses reuse existing `Envelope.Session` to carry the created runtime session id and validated player id.
- [x] Response session metadata can derive from already validated application identity.
- [x] Metadata-only identity is not treated as proof and is not upgraded.
- [x] No Protobuf source, generated output, reconnect/resume token, handshake auth, close behavior, session revocation, dependency, or direct compatibility behavior is added.
- [x] Focused tests cover login response session carrier and metadata-only preservation.
