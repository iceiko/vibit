# Plan

1. Record the final maintainer `go` authorization and exact release boundary in a durable standard, ADR, conversation log, change spec, and manifests.
2. Refresh README and README.zh-CN.md as external alpha user entry points while preserving alpha honesty and deferrals.
3. Add release notes for `v0.1.0-alpha.1` and its Simplified Chinese translation.
4. Add repository check coverage for `runtime.release_execution_final_authorization`.
5. Run conflict checks for local tags, remote tags, and GitHub release records.
6. Run repository and runtime verification.
7. If verification passes, commit and push the authorization/README update, create and push tag `v0.1.0-alpha.1`, and create the GitHub release record.
8. Stop immediately if a version conflict, failed verification, new untriaged warning, unredacted secret risk, or target-commit change occurs.
