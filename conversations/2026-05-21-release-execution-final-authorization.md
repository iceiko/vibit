# Conversation: Release Execution Final Authorization

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0103-release-execution-final-authorization.md`
Related change: `changes/2026-05-21-confirm-release-execution-final-authorization/`

## Context

The release identifier and artifact plan had completed with proposed identifier `v0.1.0-alpha.1`. The queue was blocked at `W-0195 Confirm release execution final authorization` because release execution would create public tag and GitHub release state.

## Maintainer Narrative

The maintainer authorized execution:

```text
授权执行 W-0195：go，release=v0.1.0-alpha.1，允许创建并推送 Git tag，允许创建 GitHub Release，仅发布 GitHub source archive，不发布二进制/包/容器/checksum/签名/部署/公告；若发现版本冲突或验证失败则立即停止
```

The maintainer also requested that README be made more attractive for developers to try, because the project needs to find users.

## Agent Response Summary

The agent interpreted the authorization as sufficient to record `W-0195` as go and to run release-producing commands after conflict checks and verification pass.

The agent also treated the README refresh as authorized developer-experience work inside the release execution slice, while preserving the runtime, protocol, generated-output, migration, dependency, hosted deployment, broad artifact, and direct compatibility deferrals.

## Decisions

- Accept `ADR-0103`.
- Complete `M-123/W-0195`.
- Authorize release identifier `v0.1.0-alpha.1`.
- Authorize Git tag creation and push for `v0.1.0-alpha.1`.
- Authorize GitHub release record creation for `v0.1.0-alpha.1`.
- Authorize GitHub source archive only.
- Authorize README user-acquisition refresh.
- Defer binaries, packages, containers, checksums, signing/provenance artifacts, hosted deployment, install scripts, registry publication, and public announcements beyond the GitHub release record.
- Preserve runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, operations/admin behavior, authentication/session behavior, broad product modules, and direct Nakama/Pitaya API compatibility.
- Open `M-124/W-0196 Define first alpha user discovery loop` as the next ready direction.

## Artifacts

- `docs/release-execution-final-authorization.md`
- `docs/release-execution-final-authorization.zh-CN.md`
- `decisions/ADR-0103-release-execution-final-authorization.md`
- `changes/2026-05-21-confirm-release-execution-final-authorization/`
- `docs/releases/v0.1.0-alpha.1.md`
- `docs/releases/v0.1.0-alpha.1.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The actual tag and GitHub release creation remain pending until verification and conflict checks pass.
- Public announcements beyond the GitHub release record remain unauthorized.
- The first alpha user discovery loop remains to be defined in `W-0196`.

## Follow-Up

- Run the required conflict checks and verification commands.
- Commit and push the release authorization and README update.
- Create and push tag `v0.1.0-alpha.1` if verification passes.
- Create GitHub Release `v0.1.0-alpha.1` if the tag push succeeds.
- Continue with `W-0196 Define first alpha user discovery loop`.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, private local environment file content, or GitHub tokens were recorded.
