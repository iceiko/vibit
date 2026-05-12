# Conversation: Add TypeScript Package Baseline

Date: 2026-05-12
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-12-add-typescript-package-baseline/`

Related artifacts:

- `package.json`
- `package-lock.json`
- `tsconfig.json`
- `tools/vibit`
- `.arch/runtime.yaml`

## Context

The project now has two executable TypeScript runtime paths in the inventory module: `GrantItem` and `GetInventory`. Until now, runtime verification used Node.js built-in TypeScript stripping and did not run a TypeScript compiler.

## Maintainer Narrative

The maintainer asked:

```text
继续
```

## Agent Response Summary

The agent chose to add the smallest npm and TypeScript baseline so agents can typecheck runtime code before more executable paths accumulate. The change should keep the runtime shape small and avoid introducing a web framework, persistence layer, or broad package layout.

## Decisions

- Use npm as the current package manager because npm is available with Node.js in the current environment and introduces the least extra tooling.
- Use TypeScript as a development dependency for no-emit typechecking.
- Keep Node.js built-in test runner for current runtime tests.
- Keep `node tools/vibit check runtime` as the vibit-owned runtime verification entry point.

## Artifacts

- Added `package.json` with `typecheck`, `test`, and `check` scripts.
- Generated `package-lock.json`.
- Added `tsconfig.json` for no-emit TypeScript checks.
- Added `tools/package.json` so the existing CommonJS CLI remains unambiguous under a root ESM package.
- Added TypeScript and Node type declarations as development dependencies.
- Updated `node tools/vibit check runtime` to run typecheck before runtime tests.
- Added runtime typecheck rule metadata.
- Updated `.arch/runtime.yaml`, README, AGENTS, and Simplified Chinese translations.

## Open Questions

- Whether a future workspace layout should introduce `runtime/` or `packages/`.
- Whether Node.js built-in tests remain enough once transport and persistence adapters exist.
- Whether schema tooling should become a dependency after the current lightweight YAML parsing reaches its limit.

## Follow-Up

- Consider adding a `check typecheck` command if typechecking grows beyond runtime checks.

## Redaction Notes

No secrets, tokens, credentials, private account details, or unrelated private data are recorded in this log.
