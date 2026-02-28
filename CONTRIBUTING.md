# 贡献指南

感谢您有兴趣为图书管理系统做出贡献！

## 📋 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)

## 行为准则

本项目采用贡献者公约作为行为准则。参与此项目即表示您同意遵守其条款。请阅读 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) 了解详情。

## 如何贡献

### 报告 Bug

如果您发现了 bug，请通过 [GitHub Issues](https://github.com/idleRain/LibraryManagement/issues) 提交报告。

提交 Bug 报告时，请包含：

1. **问题描述**：清晰简洁地描述问题
2. **复现步骤**：详细的复现步骤
3. **预期行为**：您期望发生什么
4. **实际行为**：实际发生了什么
5. **环境信息**：操作系统、浏览器、版本等
6. **截图**：如果适用，添加截图帮助解释问题

### 建议新功能

我们欢迎新功能建议！请通过 [GitHub Issues](https://github.com/idleRain/LibraryManagement/issues) 提交。

提交功能建议时，请包含：

1. **功能描述**：清晰描述您希望添加的功能
2. **使用场景**：描述这个功能的使用场景
3. **预期效果**：描述这个功能应该如何工作

### 改进文档

文档改进包括：

- 修正拼写或语法错误
- 添加缺失的文档
- 改进现有文档的清晰度
- 翻译文档

## 开发流程

### 1. Fork 项目

点击项目页面右上角的 "Fork" 按钮。

### 2. 克隆仓库

```bash
git clone https://github.com/<your-username>/LibraryManagement.git
cd LibraryManagement
```

### 3. 添加上游仓库

```bash
git remote add upstream https://github.com/idleRain/LibraryManagement.git
```

### 4. 创建分支

```bash
git checkout -b feature/your-feature-name
```

分支命名规范：

- `feature/` - 新功能
- `fix/` - Bug 修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

### 5. 安装依赖

```bash
# 安装前端依赖
cd apps/frontend && pnpm install

# 安装后端依赖
cd ../backend && go mod download
```

### 6. 进行开发

编写代码，确保遵循代码规范。

### 7. 运行测试

```bash
# 运行所有测试
make test

# 运行前端测试
make test-frontend

# 运行后端测试
make test-backend
```

### 8. 提交更改

```bash
git add .
git commit -m "feat: 添加新功能"
```

### 9. 推送到 Fork

```bash
git push origin feature/your-feature-name
```

### 10. 创建 Pull Request

在 GitHub 上创建 Pull Request。

## 代码规范

### Go 代码规范

- 遵循 [Effective Go](https://golang.org/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查

```bash
# 格式化
make fmt-backend

# 代码检查
make lint-backend
```

### TypeScript/Svelte 代码规范

- 使用 ESLint 进行代码检查
- 使用 Prettier 格式化代码

```bash
# 格式化
make fmt-frontend

# 代码检查
make lint-frontend
```

### 通用规范

- 使用有意义的变量名和函数名
- 添加必要的注释
- 保持函数简洁，单一职责
- 编写可测试的代码

## 提交规范

本项目使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

### 提交格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

| 类型 | 说明 |
|------|------|
| feat | 新功能 |
| fix | Bug 修复 |
| docs | 文档更新 |
| style | 代码格式（不影响功能） |
| refactor | 代码重构 |
| perf | 性能优化 |
| test | 测试相关 |
| chore | 构建/工具相关 |
| revert | 回退提交 |

### 示例

```
feat: 添加图书推荐功能

- 基于借阅历史推荐
- 添加热门图书推荐
- 添加个性化推荐算法

Closes #123
```

```
fix: 修复用户登录后 Token 未保存的问题

修复了在某些情况下用户登录后 Token 没有正确保存到 localStorage 的问题。

Fixes #456
```

## Pull Request 流程

### PR 检查清单

在提交 PR 之前，请确保：

- [ ] 代码通过所有测试
- [ ] 代码通过 lint 检查
- [ ] 代码已格式化
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
- [ ] 遵循提交规范

### PR 标题

PR 标题应遵循提交规范：

```
feat: 添加图书推荐功能
fix: 修复登录问题
docs: 更新 API 文档
```

### PR 描述

PR 描述应包含：

1. **更改内容**：描述您的更改
2. **相关 Issue**：链接相关 Issue
3. **测试方法**：描述如何测试这些更改
4. **截图**：如果适用，添加截图

### 代码审查

所有 PR 都需要至少一位维护者的审查。审查者可能会提出修改建议，请及时响应。

### 合并

PR 审查通过后，维护者会将其合并到主分支。

## 获取帮助

如果您有任何问题，可以：

- 查看 [文档](docs/README.md)
- 提交 [Issue](https://github.com/idleRain/LibraryManagement/issues)
- 发送邮件至 gold.experience@foxmail.com

## 许可证

通过贡献代码，您同意您的贡献将根据 MIT 许可证授权。

---

再次感谢您的贡献！🙏
