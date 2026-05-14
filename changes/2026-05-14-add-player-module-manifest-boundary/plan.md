# Plan

1. Inspect the current work queue and confirm `W-0023` is the only `next_ready` item.
2. Read the module manifest standard, inventory module boundary, player identity/session boundary standard, and ADR-0021.
3. Generate the `player` module skeleton.
4. Replace generated module placeholders with a boundary-only player manifest and bilingual module guides.
5. Register `player` in `.arch/modules.yaml`.
6. Record `player` in `.arch/contracts.yaml` as a boundary-only module with no public contracts yet.
7. Adjust checks so registered boundary-only modules do not force premature contracts or Protobuf files.
8. Complete `W-0023`, add the next bounded work item, and run verification.

## Ask-First Boundaries

Stop and ask before:

- Choosing a concrete authentication scheme.
- Choosing a token format.
- Adding JWT, OAuth, OIDC, password hashing, credential storage, or identity-provider dependencies.
- Adding player account migrations.
- Changing the Protobuf envelope.
- Changing the WebSocket handshake contract.
- Moving inventory ownership into the player module.
