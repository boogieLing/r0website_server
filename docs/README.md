# 文档总览

项目的详细文档统一放在 `docs/` 目录下，根目录仅保留项目级说明（`README.md`）和贡献规范（`AGENTS.md`）。

## API 文档

- `API文档.md`：主版 API 文档（带目录与锚点）。
- `API文档_兼容版.md`：兼容旧系统或工具的 API 文档版本。
- `API文档_纯MD.md`：仅使用 Markdown 锚点的纯文本版本，便于脚本处理。

API 文档相关的校验和重构脚本位于 `scripts/`：

- `scripts/test_jump_links.py`
- `scripts/test_pure_md_links.py`
- `scripts/test_real_anchors.py`
- `scripts/check_api_doc.py`

## AI / Agent 使用说明

- `CLAUDE.md`：面向 Claude Code 的项目说明与开发约定。
- `AGENTS.md`（根目录）：面向所有代码代理的仓库级指南。

## 其他说明

- 结构或文档有较大调整时，请同步更新本文件和相关脚本中的路径。
- 新增文档时优先放在 `docs/` 下，并在此处补充索引条目。
