# Request

## Original Request

The maintainer asked:

> 继续

## Clarified Requirement

Add the smallest TypeScript package baseline now that the repository has both command-side and query-side executable runtime paths.

This change should make runtime verification more regular without adding a web framework, persistence layer, or broad package structure.

## User-Visible Outcome

Maintainers and agents can run standard npm scripts for the current TypeScript runtime slice:

```bash
npm run typecheck
npm test
npm run check
```

The existing `node tools/vibit check runtime` remains the canonical vibit check and should include typechecking when the TypeScript baseline exists.

## Non-Goals

- Do not add HTTP, WebSocket, transport, or server process code.
- Do not add persistence.
- Do not add a major runtime framework dependency.
- Do not move existing modules into a new package layout.
- Do not compile emitted JavaScript as part of the default flow.
- Do not make TypeScript the only possible future implementation language.

## Unknowns

- Whether future workspace layout should introduce `runtime/` or `packages/`.
- Whether the long-term test runner should remain Node.js built-in test runner.
- Whether dependency installation should remain minimal or later include dedicated schema tooling.

## Acceptance Criteria

- [x] `package.json` declares the project as private and module-based.
- [x] `package.json` provides `typecheck`, `test`, and `check` scripts.
- [x] `tsconfig.json` typechecks current TypeScript files without emitting JavaScript.
- [x] TypeScript is added as a minimal development dependency.
- [x] `node tools/vibit check runtime` runs typecheck before runtime tests when `package.json` exists.
- [x] `node tools/vibit check all --json` passes.
- [x] Runtime readiness docs/manifests reflect npm as the current package manager and Node built-in tests as the current test runner.
- [x] English and Simplified Chinese docs are updated together.
- [x] Verification is recorded.
