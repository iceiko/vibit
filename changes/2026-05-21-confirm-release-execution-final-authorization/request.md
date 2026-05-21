# Request

Record final release execution authorization for `W-0195`, execute the source-first `v0.1.0-alpha.1` release if verification passes, and make the README more attractive for developers to try.

## Maintainer Input

```text
授权执行 W-0195：go，release=v0.1.0-alpha.1，允许创建并推送 Git tag，允许创建 GitHub Release，仅发布 GitHub source archive，不发布二进制/包/容器/checksum/签名/部署/公告；若发现版本冲突或验证失败则立即停止 同时把我们的README做成特别吸引人去尝试的状态，我们要找到用户。
```

## Interpretation

The maintainer has provided the required final `go` decision and exact release boundary:

- use release identifier `v0.1.0-alpha.1`;
- create and push Git tag `v0.1.0-alpha.1`;
- create GitHub Release `v0.1.0-alpha.1`;
- publish only the GitHub source archive;
- stop on version conflict or verification failure;
- refresh README as the external alpha user entry point.

This request does not authorize:

- binaries,
- packages,
- containers,
- checksum files,
- signing or provenance artifacts,
- hosted deployments,
- install scripts,
- registry publication,
- public announcements beyond the GitHub release record,
- runtime behavior changes,
- protocol route changes,
- Protobuf source or generated output changes,
- migrations,
- dependencies,
- broad operations/admin behavior,
- authentication/session behavior changes,
- broad product module expansion,
- direct Nakama/Pitaya API compatibility.
