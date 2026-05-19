# Plan

1. Read the protocol, runtime adapter, authentication validation, and route-protection gate standards.
2. Add the authenticated request Protobuf source.
3. Regenerate Go Protobuf output through Buf.
4. Add application route policy and route access-token validation handoff.
5. Extend the Protobuf frame handler to unwrap authenticated request payloads before protected dispatch.
6. Add tests for missing wrapper, missing token, malformed wrapper, invalid token, store-unavailable failure, metadata-only identity rejection, valid identity dispatch, public auth route behavior, and unchanged envelope route fields.
7. Add a WebSocket transport regression test proving handshake metadata is not treated as a credential carrier.
8. Update repository checks, manifests, module metadata, agent guides, work queue, change spec, and conversation log.
9. Run verification.
