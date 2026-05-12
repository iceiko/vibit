# Conversation: Go WebSocket Protobuf Direction

Date: 2026-05-12
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-12-ratify-go-websocket-protobuf-runtime/`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/contracts.yaml`
- `decisions/ADR-0003-first-reference-runtime-language.md`
- `decisions/ADR-0008-go-server-runtime-language.md`
- `decisions/ADR-0009-websocket-protobuf-client-protocol.md`
- `decisions/ADR-0010-foundational-dependency-adoption.md`

## Context

The repository had previously recorded TypeScript/Node.js as the first reference runtime and HTTP as a possible first protocol direction. The maintainer challenged that assumption because it had not been explicitly discussed or approved.

## Maintainer Narrative

The maintainer objected to the unconfirmed runtime and protocol assumptions:

```text
你什么时候和我讨论了服务器的开发语言？我什么时候同意你第一参考实验语言用TS了？我什么时候让你协议是HTTP first了？你为什么都不跟我确认和讨论？
```

The maintainer then proposed the preferred direction:

```text
我觉得服务器用Go，Go lang更合适。我觉得我们第一个协议应该用WebSocket。另外，我们觉得我们应该用ProtoBuf
```

The maintainer emphasized that decisions should be analyzed before being written into architecture:

```text
你应该先和我讨论，并且分析这个决策的正确性，是否合适。
```

The maintainer clarified that mature open-source dependencies may be integrated, but vibit must not be treated as a demo:

```text
我们可以集成其他的高星高频的GitHub的项目、开源项目。比如协议，比如网络。另外，我们绝不应该有认为这是一个demo的想法，这必然是一个长期维护和研发的系统。
```

The maintainer confirmed option B for two unresolved choices:

```text
两个都选B，继续推进。
```

The selected choices were:

- vibit manifests remain the business semantics source; Protobuf owns wire schema.
- The TypeScript runtime slice should not remain the ratified server runtime direction.

## Agent Response Summary

The agent acknowledged that silently ratifying TypeScript/Node.js and HTTP-first was wrong for a major architecture decision.

The agent analyzed the corrected direction:

- Go is a better first server runtime for a long-lived server framework with concurrency, networking, deployment, and operational clarity needs.
- WebSocket fits the first game-server client protocol better than HTTP because bidirectional sessions and server push are core expectations.
- Protobuf fits Go and cross-language client/server protocols better than ad hoc JSON once the framework is serious.
- Protobuf should not replace vibit manifests as the source of business semantics.
- High-star and frequently used open-source projects are useful inputs but require explicit adoption records.

## Decisions

- `ADR-0003` is superseded.
- `ADR-0008` ratifies Go as the first server runtime implementation language.
- `ADR-0009` ratifies WebSocket and Protobuf for the first client protocol and wire format.
- `ADR-0010` requires adoption records for foundational dependencies.
- The TypeScript runtime slice and npm package baseline are removed from the mainline.

## Artifacts

- Updated runtime and contract manifests.
- Updated repository and module agent guides.
- Updated README files in English and Simplified Chinese.
- Updated inventory contract implementation metadata.
- Removed TypeScript runtime and npm baseline files.

## Open Questions

- First Go workspace layout.
- First WebSocket library.
- First Protobuf toolchain.
- First `.proto` package and directory convention.
- Manifest-to-proto consistency check design.

## Follow-Up

- Create a dependency adoption record for WebSocket candidates.
- Create a dependency adoption record for Protobuf tooling.
- Define Go runtime workspace layout before implementation code starts.

## Redaction Notes

No secrets, tokens, credentials, private account details, or unrelated private data are recorded in this log.
