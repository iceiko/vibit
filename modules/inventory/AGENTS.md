# inventory Module Agent Guide

Status: Draft v0.1

## When To Use This Module

Describe the requirements that belong in this module.

## When Not To Use This Module

Describe requirements that should not be implemented in this module.

## Extension Points

- Commands
- Queries
- Events
- Policies
- Tests

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not directly modify data owned by another module.
- Do not add unregistered public commands, queries, events, or permissions.

## Required Tests

See `tests.required` in `module.yaml`.
