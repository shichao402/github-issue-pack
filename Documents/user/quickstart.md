# 快速入门

## 安装

```bash
cursortoolset install github-issue
```

## 配置 Token

设置 GitHub Personal Access Token：

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
```

Token 需要以下权限：
- `repo` 或 `public_repo`
- `gist`

## 场景一：提交功能请求

### 1. 准备请求内容

创建 `request.json`：

```json
{
  "title": "添加本地调试支持",
  "description": "希望 github-actions 包能支持使用 act 工具本地调试 workflow",
  "use_case": "开发 workflow 时需要频繁测试，每次推送到 GitHub 太慢",
  "expected_behavior": "提供 act 工具的安装和使用指南"
}
```

### 2. 创建 Issue

```bash
github-issue create \
  --repo shichao402/CursorColdStart \
  --type feature-request \
  --title "github-actions: 添加本地调试支持" \
  --payload request.json
```

### 3. 查看结果

命令会输出创建的 Issue 链接。

## 场景二：处理收到的 Issue

### 1. 列出待处理 Issue

```bash
github-issue list --status pending
```

输出：
```
#  | Type            | Title                              | Created
---|-----------------|------------------------------------|---------
123| feature-request | github-actions: 添加本地调试支持    | 2024-12-07
124| bug-report      | 安装失败                           | 2024-12-07
```

### 2. 获取 Issue 详情

```bash
github-issue get 123
```

输出结构化 JSON，AI 可以直接解析处理。

### 3. 标记处理中

```bash
github-issue update 123 --status processing
```

### 4. 处理完成后关闭

```bash
github-issue close 123 --result success --comment "已在 v1.1.0 中添加"
```

## 场景三：注册新包到 CursorToolset

```bash
github-issue create \
  --repo shichao402/CursorToolset \
  --type pack-register \
  --title "[auto-register] my-awesome-pack" \
  --payload register.json
```

`register.json`：
```json
{
  "repository": "https://github.com/myname/my-awesome-pack",
  "name": "my-awesome-pack",
  "version": "1.0.0",
  "description": "我的工具包"
}
```

## 下一步

- 查看 [命令参考](commands.md) 了解所有命令
- 查看 [数据格式](../design/data-format.md) 了解包格式规范
