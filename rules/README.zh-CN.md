# Rule Catalogs 中文版

状态：Draft v0.1  
最后更新：2026-05-12  
范围：机器可读 rule metadata  
说明：本文件是 `rules/README.md` 的简体中文译本。英文版本是权威版本。

本目录存放 vibit 工具输出使用的机器可读 rule catalogs。

初始 catalog：

```text
rules/check-rules.json
```

Check rule catalog 将 `node tools/vibit check ... --json` 输出中的 `rule_id` 映射到人类可读 metadata。

每条 rule 应声明：

- 稳定的 `rule_id`
- Category
- Default severity
- Title
- Description
- Agent guidance

Rules 目前还不是最终 public API。它们在 standards 稳定前使用 `0.1` 版本。
