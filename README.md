# GitHub Issue Pack

标准化的 GitHub Issue 创建与处理工具，支持项目间自动化协作。

## 功能概述

本包提供了一套标准化的 Issue 协作机制，让项目之间可以通过 Issue 进行自动化通信：

- **标准化打包**：将请求内容打包成规范格式，上传到 Gist
- **标准化解包**：从 Issue 中提取 Gist 链接，解析结构化数据
- **状态管理**：通过标签管理 Issue 的处理状态
- **AI 友好**：结构化数据便于 AI 理解和处理

## 使用场景

1. **功能请求**：向其他项目提交 feature request
2. **Bug 报告**：向其他项目报告 bug
3. **包注册**：向 CursorToolset 注册新包
4. **包同步**：通知 CursorToolset 同步包版本

## 安装

```bash
cursortoolset install github-issue
```

## 命令

详见 [Documents/user/commands.md](Documents/user/commands.md)

## 文档

- [用户指南](Documents/user/)
- [开发文档](Documents/development/)
- [设计文档](Documents/design/)
