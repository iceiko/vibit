# Impact

## Contracts

Adds a new runtime `authentication` contract family with:

- Four command surfaces.
- Six audit-oriented event surfaces.
- One error catalog.
- One permission catalog.

The contracts are semantic source files only.

## Runtime

No runtime authentication behavior is added.

No Go runtime handler, interface, repository, adapter, route, or service implementation is added.

## Protocol

No Protobuf source or generated output changes are made.

The selected request proof carrier remains an explicit semantic payload field until a later protocol decision ratifies a wire shape.

## Transport

No WebSocket handshake or transport behavior changes are made.

The WebSocket transport remains credential-neutral.

## Persistence

No migration, table, repository, or adapter is added.

Credential and token verifier storage remain deferred to W-0071 schema gates.

## Tooling

`tools/vibit` now recognizes runtime authentication contracts for contract inspection and contract checks.

## Documentation

Adds an English standard and Simplified Chinese translation.

Adds ADR-0028 and a conversation log.
