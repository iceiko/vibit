# Plan

1. Keep `W-0178` as a Tier 2 functional slice under `ADR-0082`.
2. Define the embedded presence protocol query boundary in `spec.yaml`.
3. Add the protocol source for `vibit.presence.v1.GetPlayerPresence`.
4. Regenerate Go Protobuf output with Buf.
5. Add an application-owned presence query handler over `connection.InMemoryRegistry`.
6. Add the Protobuf bridge and payload registry import.
7. Register the query in PostgreSQL runtime composition.
8. Add focused app, Protobuf bridge/frame, and startup composition tests.
9. Update ADR, conversation memory, manifests, guides, check catalog, and work queue.
10. Verify focused tests, full runtime tests, generated/protocol/runtime/work/change checks, full repository checks, and diff whitespace.
