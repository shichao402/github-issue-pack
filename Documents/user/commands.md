# 命令参考

## github-issue create

创建标准化的 GitHub Issue。

### 语法

```bash
github-issue create --repo <owner/repo> --type <type> --title <title> [options]
```

### 参数

| 参数 | 必需 | 说明 |
|------|------|------|
| `--repo` | ✅ | 目标仓库（格式：owner/repo） |
| `--type` | ✅ | Issue 类型（feature-request/bug-report/pack-register/pack-sync） |
| `--title` | ✅ | Issue 标题 |
| `--payload` | ❌ | 详细内容文件路径（JSON 格式） |
| `--attach` | ❌ | 附件文件路径（可多次使用） |
| `--dry-run` | ❌ | 预览模式，不实际创建 |

### 示例

```bash
# 创建功能请求
github-issue create \
  --repo shichao402/CursorColdStart \
  --type feature-request \
  --title "添加本地调试支持" \
  --payload request.json

# 预览模式
github-issue create \
  --repo shichao402/CursorColdStart \
  --type feature-request \
  --title "测试" \
  --dry-run
```

---

## github-issue list

列出待处理的 Issue。

### 语法

```bash
github-issue list [options]
```

### 参数

| 参数 | 必需 | 说明 |
|------|------|------|
| `--status` | ❌ | 状态过滤（pending/processing/processed/all），默认 pending |
| `--type` | ❌ | 类型过滤 |
| `--limit` | ❌ | 数量限制，默认 20 |
| `--format` | ❌ | 输出格式（table/json），默认 table |

### 示例

```bash
# 列出所有待处理的 issue
github-issue list

# 列出所有功能请求
github-issue list --type feature-request --status all

# JSON 格式输出
github-issue list --format json
```

---

## github-issue get

获取并解析指定 Issue。

### 语法

```bash
github-issue get <issue-number> [options]
```

### 参数

| 参数 | 必需 | 说明 |
|------|------|------|
| `<issue-number>` | ✅ | Issue 编号 |
| `--format` | ❌ | 输出格式（json/yaml/text），默认 json |
| `--output` | ❌ | 输出到文件 |

### 示例

```bash
# 获取 issue #123 的结构化数据
github-issue get 123

# 输出为 YAML 格式
github-issue get 123 --format yaml

# 保存到文件
github-issue get 123 --output issue-123.json
```

---

## github-issue close

关闭并标记 Issue 处理结果。

### 语法

```bash
github-issue close <issue-number> --result <result> [options]
```

### 参数

| 参数 | 必需 | 说明 |
|------|------|------|
| `<issue-number>` | ✅ | Issue 编号 |
| `--result` | ✅ | 处理结果（success/rejected） |
| `--comment` | ❌ | 处理说明 |

### 示例

```bash
# 标记处理成功
github-issue close 123 --result success --comment "已添加到 v1.1.0"

# 标记拒绝
github-issue close 123 --result rejected --comment "不符合项目规范"
```

---

## github-issue update

更新 Issue 状态。

### 语法

```bash
github-issue update <issue-number> --status <status> [options]
```

### 参数

| 参数 | 必需 | 说明 |
|------|------|------|
| `<issue-number>` | ✅ | Issue 编号 |
| `--status` | ✅ | 新状态（processing/pending） |
| `--comment` | ❌ | 添加评论 |

### 示例

```bash
# 标记为处理中
github-issue update 123 --status processing --comment "开始处理"

# 退回待处理
github-issue update 123 --status pending --comment "需要更多信息"
```

---

## 环境变量

| 变量 | 说明 |
|------|------|
| `GITHUB_TOKEN` | GitHub Personal Access Token（可选，默认使用 gh CLI 认证） |

> **注意**：`--repo` 参数是必需的，不支持默认仓库配置。这是有意为之的设计，遵循「显式优于隐式」原则，避免误操作将 Issue 提交到错误的仓库。

## Token 权限要求

- `repo` 或 `public_repo`：创建/关闭 Issue
- `gist`：创建 Gist 存储 payload
