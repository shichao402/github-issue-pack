# GitHub Issue Pack 功能设计

## 概述

本包提供标准化的 GitHub Issue 创建与处理机制，实现项目间的自动化协作。

## 核心功能

### 1. Issue 创建（发送方）

**命令**：`github-issue create`

**流程**：
1. 读取请求内容（从文件或参数）
2. 构建标准化包格式
3. 创建 Gist 存储完整数据
4. 创建 Issue，body 包含摘要 + Gist 链接
5. 添加标准标签

**参数**：
```bash
github-issue create \
  --repo <owner/repo>           # 目标仓库
  --type <issue-type>           # issue 类型
  --title <title>               # 标题
  --payload <file.json>         # 详细内容（可选）
  --attach <file1> <file2>      # 附件（可选）
```

### 2. Issue 列表（接收方）

**命令**：`github-issue list`

**流程**：
1. 查询带有 `cursortoolset` 标签的 issue
2. 按状态过滤（pending/processing/processed）
3. 输出结构化列表

**参数**：
```bash
github-issue list \
  --status <pending|processing|processed|all>  # 状态过滤
  --type <issue-type>                           # 类型过滤
  --limit <n>                                   # 数量限制
```

### 3. Issue 获取（接收方）

**命令**：`github-issue get`

**流程**：
1. 获取指定 issue
2. 解析 body 中的 Gist 链接
3. 下载 Gist 内容
4. 解包并输出结构化数据

**参数**：
```bash
github-issue get <issue-number> \
  --format <json|yaml|text>     # 输出格式
  --output <file>               # 输出到文件（可选）
```

### 4. Issue 关闭（接收方）

**命令**：`github-issue close`

**流程**：
1. 添加处理结果评论
2. 更新标签（processed/rejected）
3. 关闭 issue

**参数**：
```bash
github-issue close <issue-number> \
  --result <success|rejected>   # 处理结果
  --comment <message>           # 处理说明
```

### 5. Issue 状态更新

**命令**：`github-issue update`

**参数**：
```bash
github-issue update <issue-number> \
  --status <processing|pending> # 更新状态标签
  --comment <message>           # 添加评论（可选）
```

## 标签规范

| 标签 | 含义 | 颜色建议 |
|------|------|----------|
| `cursortoolset` | 由本工具创建的 issue | #7057ff |
| `pending` | 待处理 | #fbca04 |
| `processing` | 处理中 | #0e8a16 |
| `processed` | 已处理完成 | #6f42c1 |
| `rejected` | 已拒绝 | #d73a4a |
| `feature-request` | 功能请求 | #a2eeef |
| `bug-report` | Bug 报告 | #d73a4a |
| `pack-register` | 包注册请求 | #0075ca |
| `pack-sync` | 包同步请求 | #0075ca |

## Issue Body 模板

```markdown
## {type}: {title}

**Type:** {type}
**Created by:** cursortoolset v{version}
**Source:** {source_project}

### Summary

{summary}

### Details

📦 [View full payload]({gist_url})

---
<sub>This issue was automatically created by [github-issue-pack](https://github.com/shichao402/github-issue-pack)</sub>
```

## 状态流转

```
[创建 Issue]
     │
     ▼
┌─────────┐
│ pending │ ← 初始状态
└────┬────┘
     │ github-issue update --status processing
     ▼
┌────────────┐
│ processing │ ← 处理中
└─────┬──────┘
      │ github-issue close
      ▼
┌─────────────────────────┐
│ processed OR rejected   │ ← 终态
└─────────────────────────┘
```

## 权限要求

| 操作 | 所需权限 |
|------|----------|
| 创建 Issue | `repo` 或 `public_repo` |
| 创建 Gist | `gist` |
| 读取 Issue | 公开仓库无需权限 |
| 关闭 Issue | 仓库写权限 |

## 错误处理

| 错误码 | 含义 | 处理方式 |
|--------|------|----------|
| `E001` | 目标仓库不存在 | 检查仓库地址 |
| `E002` | 无权限创建 Issue | 检查 token 权限 |
| `E003` | Gist 创建失败 | 检查 token 权限 |
| `E004` | Issue 不存在 | 检查 issue 编号 |
| `E005` | 无效的包格式 | 检查 payload 格式 |
