# Issue 包格式规范

## 概述

本文档定义了 GitHub Issue 包的标准数据格式，用于项目间的自动化协作。

## Schema 版本

当前版本：`v1`

## 包格式（Gist 内容）

```json
{
  "$schema": "cursortoolset-issue-v1",
  "meta": {
    "created_at": "2024-12-07T10:00:00Z",
    "source_project": "owner/repo",
    "cursortoolset_version": "1.7.0",
    "github_issue_version": "0.1.0"
  },
  "type": "feature-request",
  "target": {
    "repo": "owner/target-repo",
    "pack": "pack-name",
    "version": "1.0.0"
  },
  "payload": {
    // 具体内容，根据 type 不同而不同
  },
  "attachments": [
    {
      "name": "context.json",
      "content": "..."
    }
  ]
}
```

## 字段说明

### meta（元数据）

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| created_at | string | ✅ | ISO 8601 格式的创建时间 |
| source_project | string | ✅ | 来源项目（owner/repo） |
| cursortoolset_version | string | ❌ | CursorToolset 版本 |
| github_issue_version | string | ✅ | 本包版本 |

### type（Issue 类型）

| 类型 | 说明 |
|------|------|
| `feature-request` | 功能请求 |
| `bug-report` | Bug 报告 |
| `pack-register` | 包注册请求 |
| `pack-sync` | 包同步请求 |
| `question` | 问题咨询 |
| `custom` | 自定义类型 |

### target（目标信息）

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| repo | string | ✅ | 目标仓库 |
| pack | string | ❌ | 相关的包名 |
| version | string | ❌ | 相关的版本号 |

### payload（具体内容）

根据 `type` 不同，payload 结构不同。

#### feature-request

```json
{
  "payload": {
    "title": "功能标题",
    "description": "详细描述",
    "use_case": "使用场景",
    "expected_behavior": "期望的行为",
    "alternatives": "替代方案（可选）"
  }
}
```

#### bug-report

```json
{
  "payload": {
    "title": "Bug 标题",
    "description": "问题描述",
    "steps_to_reproduce": [
      "步骤1",
      "步骤2"
    ],
    "expected_behavior": "期望行为",
    "actual_behavior": "实际行为",
    "environment": {
      "os": "macOS 14.0",
      "cursortoolset_version": "1.7.0",
      "pack_version": "1.0.0"
    }
  }
}
```

#### pack-register

```json
{
  "payload": {
    "repository": "https://github.com/owner/pack-repo",
    "name": "pack-name",
    "version": "1.0.0",
    "description": "包描述"
  }
}
```

#### pack-sync

```json
{
  "payload": {
    "repository": "https://github.com/owner/pack-repo",
    "version": "1.1.0",
    "changes": "更新内容摘要"
  }
}
```

### attachments（附件）

附件直接嵌入 Gist 中，每个附件是一个独立文件。

```json
{
  "attachments": [
    {
      "name": "logs.txt",
      "content": "日志内容..."
    },
    {
      "name": "config.json",
      "content": "{\"key\": \"value\"}"
    }
  ]
}
```

## Gist 结构

创建的 Gist 包含以下文件：

```
gist/
├── issue-payload.json    # 主包内容
├── logs.txt              # 附件1（如有）
└── config.json           # 附件2（如有）
```

## 验证规则

1. `$schema` 必须是 `cursortoolset-issue-v1`
2. `meta.created_at` 必须是有效的 ISO 8601 时间
3. `type` 必须是预定义类型之一
4. `payload` 必须符合对应 type 的结构

## 版本兼容

- 解析时应忽略未知字段
- 缺失的可选字段应使用默认值
- schema 版本不兼容时应报错
