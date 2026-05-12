# Request

Date: 2026-05-12
Change ID: `ratify-go-websocket-protobuf-runtime`
Type: standard

## Maintainer Request

The maintainer challenged the earlier assumption that TypeScript/Node.js and HTTP-first were decided:

```text
你什么时候和我讨论了服务器的开发语言？我什么时候同意你第一参考实验语言用TS了？我什么时候让你协议是HTTP first了？你为什么都不跟我确认和讨论？
```

The maintainer then stated the preferred direction:

```text
我觉得服务器用Go，Go lang更合适。我觉得我们第一个协议应该用WebSocket。另外，我们觉得我们应该用ProtoBuf
```

The maintainer also clarified that architecture decisions should be discussed and analyzed before being recorded:

```text
你应该先和我讨论，并且分析这个决策的正确性，是否合适。
```

The maintainer emphasized that mature open-source projects may be integrated for protocol and networking, and that vibit must not be treated as a demo:

```text
我们可以集成其他的高星高频的GitHub的项目、开源项目。比如协议，比如网络。另外，我们绝不应该有认为这是一个demo的想法，这必然是一个长期维护和研发的系统。
```

The agent asked the maintainer to choose how Protobuf relates to vibit contracts and how to treat the existing TypeScript runtime slice. The maintainer chose option B for both:

```text
两个都选B，继续推进。
```

The selected meanings are:

- vibit manifests and contract files remain the business semantics source of truth; Protobuf owns wire schema.
- The TypeScript runtime slice should not remain the ratified server runtime direction.

## Goal

Correct the architecture direction before more runtime work accumulates.

The repository should clearly state:

- Go is the first server runtime implementation language.
- WebSocket is the first gameplay/client protocol.
- Protobuf is the first client/server wire message format.
- vibit semantic manifests remain separate from Protobuf wire schemas.
- Foundational dependencies require explicit adoption records.
- The previous TypeScript runtime slice and npm package baseline are removed from the mainline direction.

## Non-Goals

- Do not implement the Go runtime yet.
- Do not choose a specific WebSocket library yet.
- Do not choose a specific Protobuf toolchain yet.
- Do not hand-roll a Protobuf generator before the dependency and layout standards are decided.
- Do not rewrite historical change specs to hide earlier work.

## Acceptance Criteria

- [x] `ADR-0003` is marked superseded.
- [x] New ADRs record Go, WebSocket/Protobuf, and dependency adoption decisions.
- [x] `.arch/runtime.yaml` reflects Go/WebSocket/Protobuf direction.
- [x] `.arch/contracts.yaml` distinguishes semantic contracts from Protobuf wire schema.
- [x] Public docs and Chinese translations are updated.
- [x] Misleading TypeScript server runtime files and package baseline are removed from the mainline.
- [x] Runtime check no longer requires TypeScript tests before Go runtime exists.
