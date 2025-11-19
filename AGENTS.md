# Repository Guidelines

本文件适用于整个仓库，面向所有贡献者和自动化代理，请在阅读、修改或扩展代码时遵循以下约定。

## 项目结构与模块组织

- `main/main.go`：服务入口。
- `core/`、`global/`、`initialize/`：启动流程、配置加载、日志、MongoDB 与 COS 初始化。
- `api/`、`router/`：Gin 路由和 HTTP 控制器（`base`、`admin`）。
- `service/`、`dao/`：业务逻辑与数据访问。
- `models/po`、`models/vo`、`models/bo`：持久化对象、视图对象和业务对象。
- `middleware/`、`utils/`、`r0Ioc/`：中间件、通用工具和 IoC 容器。
- `config/`：`config.yml` 及相关配置代码；`scripts/`：API 文档与跳转测试脚本。

## 构建、测试与本地开发

- 本地运行：`./run-server.sh`（执行 `go mod tidy`、构建 `r0Website-server` 并使用开发配置启动）。
- 手动构建：`go build -o ./r0Website-server ./main/main.go`。
- Docker：`./docker-build.sh` 构建镜像并在 `8202` 端口启动容器。
- 基本检查：`go vet ./...` 与 `go test ./...`（为新增逻辑补充测试）。
- API 文档校验：如 `python3 scripts/test_jump_links.py` 检查目录锚点是否正确。

## 编码风格与命名约定

- Go 版本 1.17+；提交前务必运行 `gofmt` / `goimports`。
- 包与文件名使用小写加下划线；导出标识符使用 `PascalCase`，非导出使用 `camelCase`。
- 控制器类型统一命名为 `*Controller`，按目录拆分路由、服务与 DAO，保持现有的中文注释风格。

## 测试指南

- 为 `service/`、`utils/` 中的非简单逻辑编写 `_test.go` 单元测试。
- 在提交或发起 PR 前确保 `go test ./...` 全部通过。
- 修改 API 文档时，应运行 `scripts/` 下相关脚本（如 `test_jump_links.py`、`test_pure_md_links.py`）并修复所有报错。

## 提交与 Pull Request 规范

- 提交信息简洁明确，建议使用前缀风格：如 `feat: 添加图片上传缩略图生成功能`、`fix: 修复API文档跳转系统`。
- 单个 PR 聚焦一个主题，描述主要变更、影响范围以及验证步骤（附常用命令），必要时提供日志或截图。
- 如有关联 issue，请在描述中标出，例如 `Fixes #123`。

