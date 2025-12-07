# GitHub Issue 协作规则

## 核心约束（强制）

1. **所有跨项目协作必须使用标准化 Issue 格式**
2. **Issue 内容必须通过 Gist 存储，便于 AI 解析**
3. **必须使用标签管理 Issue 状态**
4. **处理完成的 Issue 必须关闭并标记结果**

## 使用场景

### 何时使用 github-issue

- 向其他项目提交功能请求
- 向其他项目报告 Bug
- 向 CursorToolset 注册新包
- 任何需要项目间自动化协作的场景

### 何时不使用

- 项目内部的 Issue（使用普通 GitHub Issue）
- 简单的问题咨询（直接在 Discussions 提问）

## 命令使用规范

### 创建 Issue（发送方）

```bash
# 必须指定 --repo、--type、--title
github-issue create \
  --repo <owner/repo> \
  --type <feature-request|bug-report|pack-register|pack-sync> \
  --title "<描述性标题>" \
  --payload <file.json>  # 详细内容
```

**AI 行为要求**：
1. 创建前先确认目标仓库存在
2. 使用 `--dry-run` 预览内容
3. payload 必须符合对应 type 的格式规范

### 处理 Issue（接收方）

```bash
# 1. 列出待处理 Issue
github-issue list --status pending

# 2. 获取并解析 Issue
github-issue get <issue-number>

# 3. 标记处理中
github-issue update <issue-number> --status processing

# 4. 处理完成后关闭
github-issue close <issue-number> --result <success|rejected> --comment "<说明>"
```

**AI 行为要求**：
1. 处理前必须先 `github-issue get` 获取完整内容
2. 开始处理时必须更新状态为 `processing`
3. 处理完成必须关闭 Issue 并说明结果

## 包格式规范

### Schema 版本

当前版本：`cursortoolset-issue-v1`

### 基本结构

```json
{
  "$schema": "cursortoolset-issue-v1",
  "meta": {
    "created_at": "ISO 8601 时间",
    "source_project": "owner/repo",
    "github_issue_version": "0.1.0"
  },
  "type": "issue 类型",
  "target": {
    "repo": "目标仓库",
    "pack": "相关包名（可选）"
  },
  "payload": {
    // 具体内容
  }
}
```

### 类型对应的 payload 结构

#### feature-request

```json
{
  "payload": {
    "title": "功能标题",
    "description": "详细描述",
    "use_case": "使用场景",
    "expected_behavior": "期望行为"
  }
}
```

#### bug-report

```json
{
  "payload": {
    "title": "Bug 标题",
    "description": "问题描述",
    "steps_to_reproduce": ["步骤1", "步骤2"],
    "expected_behavior": "期望行为",
    "actual_behavior": "实际行为",
    "environment": {}
  }
}
```

#### pack-register

```json
{
  "payload": {
    "repository": "https://github.com/owner/repo",
    "name": "包名",
    "version": "版本号",
    "description": "包描述"
  }
}
```

## 标签规范

| 标签 | 含义 | 使用时机 |
|------|------|----------|
| `cursortoolset` | 标识来源 | 创建时自动添加 |
| `pending` | 待处理 | 创建时自动添加 |
| `processing` | 处理中 | 开始处理时更新 |
| `processed` | 已完成 | 关闭时添加 |
| `rejected` | 已拒绝 | 关闭时添加 |

## 禁止行为

- ❌ 手动创建不符合格式的 Issue
- ❌ 不使用 Gist 存储详细内容
- ❌ 处理完成后不关闭 Issue
- ❌ 不更新状态标签
- ❌ 重复处理已关闭的 Issue

## 必须遵守

- ✅ 使用 `github-issue` 命令创建标准化 Issue
- ✅ payload 必须符合 schema 规范
- ✅ 处理时更新状态标签
- ✅ 关闭时标记处理结果
- ✅ 添加有意义的处理说明

## 错误处理

| 错误 | 处理方式 |
|------|----------|
| 目标仓库不存在 | 检查仓库地址是否正确 |
| 无权限创建 Issue | 检查 GITHUB_TOKEN 权限 |
| Gist 创建失败 | 检查 token 是否有 gist 权限 |
| 包格式无效 | 使用 `--dry-run` 验证格式 |
