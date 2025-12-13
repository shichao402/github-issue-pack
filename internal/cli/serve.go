package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/shichao402/github-issue-pack/internal/models"
	"github.com/shichao402/github-issue-pack/internal/service"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 MCP Server 模式",
	Long: `启动 GitHub Issue Pack 的 MCP Server 模式，供 Cursor 等 IDE 调用。

此命令通过 stdin/stdout 与 IDE 通信，使用 JSON-RPC 2.0 协议。

示例：
  github-issue serve`,
	RunE: runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// MCP JSON-RPC 消息类型
type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP 协议消息
type initializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    capabilities `json:"capabilities"`
	ServerInfo      serverInfo   `json:"serverInfo"`
}

type capabilities struct {
	Tools *toolsCapability `json:"tools,omitempty"`
}

type toolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema inputSchema `json:"inputSchema"`
}

type inputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type toolsListResult struct {
	Tools []tool `json:"tools"`
}

type callToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type callToolResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func runServe(cmd *cobra.Command, args []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var request jsonRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		response := handleMCPRequest(&request)
		if response != nil {
			sendMCPResponse(response)
		}
	}

	return scanner.Err()
}

func handleMCPRequest(request *jsonRPCRequest) *jsonRPCResponse {
	switch request.Method {
	case "initialize":
		return handleMCPInitialize(request)
	case "initialized":
		return nil
	case "tools/list":
		return handleMCPToolsList(request)
	case "tools/call":
		return handleMCPToolsCall(request)
	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &rpcError{
				Code:    -32601,
				Message: "Method not found",
				Data:    request.Method,
			},
		}
	}
}

func handleMCPInitialize(request *jsonRPCRequest) *jsonRPCResponse {
	result := initializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: capabilities{
			Tools: &toolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: serverInfo{
			Name:    "github-issue",
			Version: Version,
		},
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func handleMCPToolsList(request *jsonRPCRequest) *jsonRPCResponse {
	tools := []tool{
		{
			Name:        "github_issue_create",
			Description: "创建标准化的 GitHub Issue，自动打包内容到 Gist",
			InputSchema: inputSchema{
				Type: "object",
				Properties: map[string]property{
					"repo": {
						Type:        "string",
						Description: "目标仓库 (格式: owner/repo)",
					},
					"type": {
						Type:        "string",
						Description: "Issue 类型",
						Enum:        []string{"feature-request", "bug-report", "pack-register", "pack-sync", "question", "custom"},
					},
					"title": {
						Type:        "string",
						Description: "Issue 标题",
					},
					"payload": {
						Type:        "string",
						Description: "详细内容 (JSON 字符串)",
					},
				},
				Required: []string{"repo", "type", "title"},
			},
		},
		{
			Name:        "github_issue_list",
			Description: "列出仓库中的标准化 Issue",
			InputSchema: inputSchema{
				Type: "object",
				Properties: map[string]property{
					"repo": {
						Type:        "string",
						Description: "目标仓库 (格式: owner/repo)",
					},
					"status": {
						Type:        "string",
						Description: "状态过滤",
						Enum:        []string{"pending", "processing", "processed", "all"},
					},
					"type": {
						Type:        "string",
						Description: "类型过滤",
						Enum:        []string{"feature-request", "bug-report", "pack-register", "pack-sync"},
					},
					"limit": {
						Type:        "string",
						Description: "返回数量限制 (默认 20)",
					},
				},
				Required: []string{"repo"},
			},
		},
		{
			Name:        "github_issue_get",
			Description: "获取 Issue 详情，包括解析的 payload",
			InputSchema: inputSchema{
				Type: "object",
				Properties: map[string]property{
					"repo": {
						Type:        "string",
						Description: "目标仓库 (格式: owner/repo)",
					},
					"number": {
						Type:        "string",
						Description: "Issue 编号",
					},
				},
				Required: []string{"repo", "number"},
			},
		},
		{
			Name:        "github_issue_update",
			Description: "更新 Issue 状态",
			InputSchema: inputSchema{
				Type: "object",
				Properties: map[string]property{
					"repo": {
						Type:        "string",
						Description: "目标仓库 (格式: owner/repo)",
					},
					"number": {
						Type:        "string",
						Description: "Issue 编号",
					},
					"status": {
						Type:        "string",
						Description: "新状态",
						Enum:        []string{"pending", "processing"},
					},
					"comment": {
						Type:        "string",
						Description: "添加评论 (可选)",
					},
				},
				Required: []string{"repo", "number", "status"},
			},
		},
		{
			Name:        "github_issue_close",
			Description: "关闭 Issue",
			InputSchema: inputSchema{
				Type: "object",
				Properties: map[string]property{
					"repo": {
						Type:        "string",
						Description: "目标仓库 (格式: owner/repo)",
					},
					"number": {
						Type:        "string",
						Description: "Issue 编号",
					},
					"result": {
						Type:        "string",
						Description: "处理结果",
						Enum:        []string{"success", "rejected"},
					},
					"comment": {
						Type:        "string",
						Description: "关闭说明 (可选)",
					},
				},
				Required: []string{"repo", "number", "result"},
			},
		},
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  toolsListResult{Tools: tools},
	}
}

func handleMCPToolsCall(request *jsonRPCRequest) *jsonRPCResponse {
	var params callToolParams
	if err := json.Unmarshal(request.Params, &params); err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &rpcError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    err.Error(),
			},
		}
	}

	var result callToolResult

	switch params.Name {
	case "github_issue_create":
		result = executeCreate(params.Arguments)
	case "github_issue_list":
		result = executeList(params.Arguments)
	case "github_issue_get":
		result = executeGet(params.Arguments)
	case "github_issue_update":
		result = executeUpdate(params.Arguments)
	case "github_issue_close":
		result = executeClose(params.Arguments)
	default:
		result = callToolResult{
			Content: []contentItem{{Type: "text", Text: fmt.Sprintf("未知的工具: %s", params.Name)}},
			IsError: true,
		}
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func getMCPToken() string {
	// 1. 环境变量
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}

	// 2. 从 gh CLI 获取
	token = getGhToken()
	if token != "" {
		return token
	}

	return ""
}

func executeCreate(args map[string]interface{}) callToolResult {
	repo, _ := args["repo"].(string)
	issueType, _ := args["type"].(string)
	title, _ := args["title"].(string)
	payloadStr, _ := args["payload"].(string)

	if repo == "" || issueType == "" || title == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "缺少必需参数: repo, type, title"}},
			IsError: true,
		}
	}

	token := getMCPToken()
	if token == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无法获取 GitHub Token，请设置 GITHUB_TOKEN 环境变量或运行 gh auth login"}},
			IsError: true,
		}
	}

	// 解析 payload
	var payload interface{}
	if payloadStr != "" {
		if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
			return callToolResult{
				Content: []contentItem{{Type: "text", Text: fmt.Sprintf("解析 payload 失败: %v", err)}},
				IsError: true,
			}
		}
	} else {
		payload = map[string]string{
			"title":       title,
			"description": "",
		}
	}

	svc := service.NewIssueService(token)
	result, err := svc.Create(service.CreateIssueOptions{
		Repo:    repo,
		Type:    models.IssueType(issueType),
		Title:   title,
		Payload: payload,
		DryRun:  false,
	})
	if err != nil {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: fmt.Sprintf("创建 Issue 失败: %v", err)}},
			IsError: true,
		}
	}

	return callToolResult{
		Content: []contentItem{{
			Type: "text",
			Text: fmt.Sprintf("✅ Issue 创建成功!\n\nIssue: %s\nGist: %s", result.IssueURL, result.GistURL),
		}},
	}
}

func executeList(args map[string]interface{}) callToolResult {
	repo, _ := args["repo"].(string)
	status, _ := args["status"].(string)
	issueType, _ := args["type"].(string)
	limitStr, _ := args["limit"].(string)

	if repo == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "缺少必需参数: repo"}},
			IsError: true,
		}
	}

	token := getMCPToken()
	if token == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无法获取 GitHub Token"}},
			IsError: true,
		}
	}

	limit := 20
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	svc := service.NewIssueService(token)
	issues, err := svc.List(service.ListOptions{
		Repo:   repo,
		Status: status,
		Type:   issueType,
		Limit:  limit,
	})
	if err != nil {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: fmt.Sprintf("列出 Issue 失败: %v", err)}},
			IsError: true,
		}
	}

	if len(issues) == 0 {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "没有找到匹配的 Issue"}},
		}
	}

	var text string
	text = fmt.Sprintf("找到 %d 个 Issue:\n\n", len(issues))
	for _, issue := range issues {
		text += fmt.Sprintf("- #%d [%s] %s (%s)\n  %s\n", issue.Number, issue.Status, issue.Title, issue.Type, issue.URL)
	}

	return callToolResult{
		Content: []contentItem{{Type: "text", Text: text}},
	}
}

func executeGet(args map[string]interface{}) callToolResult {
	repo, _ := args["repo"].(string)
	numberStr, _ := args["number"].(string)

	if repo == "" || numberStr == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "缺少必需参数: repo, number"}},
			IsError: true,
		}
	}

	var number int
	fmt.Sscanf(numberStr, "%d", &number)
	if number <= 0 {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无效的 Issue 编号"}},
			IsError: true,
		}
	}

	token := getMCPToken()
	if token == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无法获取 GitHub Token"}},
			IsError: true,
		}
	}

	svc := service.NewIssueService(token)
	result, err := svc.Get(repo, number)
	if err != nil {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: fmt.Sprintf("获取 Issue 失败: %v", err)}},
			IsError: true,
		}
	}

	var text string
	text = fmt.Sprintf("Issue #%d: %s\n\n", result.Issue.Number, result.Issue.Title)
	text += fmt.Sprintf("状态: %s\n", result.Issue.State)
	text += fmt.Sprintf("创建时间: %s\n", result.Issue.CreatedAt)
	text += fmt.Sprintf("URL: %s\n", result.Issue.HTMLURL)

	if result.Package != nil {
		text += fmt.Sprintf("\n类型: %s\n", result.Package.Type)
		if result.Package.Payload != nil {
			payloadJSON, _ := json.MarshalIndent(result.Package.Payload, "", "  ")
			text += fmt.Sprintf("\nPayload:\n%s\n", string(payloadJSON))
		}
	}

	return callToolResult{
		Content: []contentItem{{Type: "text", Text: text}},
	}
}

func executeUpdate(args map[string]interface{}) callToolResult {
	repo, _ := args["repo"].(string)
	numberStr, _ := args["number"].(string)
	status, _ := args["status"].(string)
	comment, _ := args["comment"].(string)

	if repo == "" || numberStr == "" || status == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "缺少必需参数: repo, number, status"}},
			IsError: true,
		}
	}

	var number int
	fmt.Sscanf(numberStr, "%d", &number)
	if number <= 0 {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无效的 Issue 编号"}},
			IsError: true,
		}
	}

	token := getMCPToken()
	if token == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无法获取 GitHub Token"}},
			IsError: true,
		}
	}

	svc := service.NewIssueService(token)
	err := svc.UpdateStatus(repo, number, status, comment)
	if err != nil {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: fmt.Sprintf("更新 Issue 失败: %v", err)}},
			IsError: true,
		}
	}

	return callToolResult{
		Content: []contentItem{{
			Type: "text",
			Text: fmt.Sprintf("✅ Issue #%d 状态已更新为 %s", number, status),
		}},
	}
}

func executeClose(args map[string]interface{}) callToolResult {
	repo, _ := args["repo"].(string)
	numberStr, _ := args["number"].(string)
	result, _ := args["result"].(string)
	comment, _ := args["comment"].(string)

	if repo == "" || numberStr == "" || result == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "缺少必需参数: repo, number, result"}},
			IsError: true,
		}
	}

	var number int
	fmt.Sscanf(numberStr, "%d", &number)
	if number <= 0 {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无效的 Issue 编号"}},
			IsError: true,
		}
	}

	token := getMCPToken()
	if token == "" {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: "无法获取 GitHub Token"}},
			IsError: true,
		}
	}

	svc := service.NewIssueService(token)
	err := svc.Close(repo, number, result, comment)
	if err != nil {
		return callToolResult{
			Content: []contentItem{{Type: "text", Text: fmt.Sprintf("关闭 Issue 失败: %v", err)}},
			IsError: true,
		}
	}

	statusText := "成功"
	if result == "rejected" {
		statusText = "已拒绝"
	}

	return callToolResult{
		Content: []contentItem{{
			Type: "text",
			Text: fmt.Sprintf("✅ Issue #%d 已关闭 (%s)", number, statusText),
		}},
	}
}

func sendMCPResponse(response *jsonRPCResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		return
	}
	fmt.Println(string(data))
}

func sendError(id interface{}, code int, message string, data interface{}) {
	response := &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &rpcError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	sendMCPResponse(response)
}
