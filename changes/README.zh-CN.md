# Change Specs 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：非平凡变更的持久上下文  
说明：本文件是 `changes/README.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

这个目录保存非平凡工作的 change specs。

权威标准见 `docs/change-spec.md`。

模板文件位于：

```text
changes/_template/
```

开始一个新变更：

1. 复制 `changes/_template/` 到 `changes/YYYY-MM-DD-short-change-id/`。
2. 填写 `request.md`。
3. 填写 `spec.yaml`。
4. 在 implementation 前完成 impact analysis 和 plan。
5. 随工作推进保持 `checklist.md` 和 `verification.md` 最新。

小 typo 修复和窄范围文档编辑可以不需要完整 change spec，但 agents 仍必须在最终交接中记录 verification status。
