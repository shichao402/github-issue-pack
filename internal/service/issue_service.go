package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shichao402/github-issue-pack/internal/github"
	"github.com/shichao402/github-issue-pack/internal/models"
)

// 标签常量
const (
	LabelCursorToolset  = "cursortoolset"
	LabelPending        = "pending"
	LabelProcessing     = "processing"
	LabelProcessed      = "processed"
	LabelRejected       = "rejected"
	LabelFeatureRequest = "feature-request"
	LabelBugReport      = "bug-report"
	LabelPackRegister   = "pack-register"
	LabelPackSync       = "pack-sync"
)

// IssueService Issue 服务
type IssueService struct {
	client *github.Client
}

// NewIssueService 创建 Issue 服务
func NewIssueService(token string) *IssueService {
	return &IssueService{
		client: github.NewClient(token),
	}
}

// CreateIssueOptions 创建 Issue 的选项
type CreateIssueOptions struct {
	Repo        string
	Type        models.IssueType
	Title       string
	Payload     interface{}
	Attachments []models.Attachment
	DryRun      bool
}

// CreateIssueResult 创建 Issue 的结果
type CreateIssueResult struct {
	IssueURL string
	GistURL  string
	IssueNum int
}

// Create 创建 Issue
func (s *IssueService) Create(opts CreateIssueOptions) (*CreateIssueResult, error) {
	owner, repo, err := parseRepo(opts.Repo)
	if err != nil {
		return nil, err
	}

	// 构建 Issue 包
	pkg, err := models.NewIssuePackage(opts.Type, opts.Repo, opts.Payload)
	if err != nil {
		return nil, fmt.Errorf("构建 Issue 包失败: %w", err)
	}
	pkg.Attachments = opts.Attachments

	// 序列化为 JSON
	pkgJSON, err := pkg.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("序列化 Issue 包失败: %w", err)
	}

	if opts.DryRun {
		fmt.Println("=== Dry Run 模式 ===")
		fmt.Printf("目标仓库: %s/%s\n", owner, repo)
		fmt.Printf("Issue 类型: %s\n", opts.Type)
		fmt.Printf("标题: %s\n", opts.Title)
		fmt.Println("\n=== Gist 内容 ===")
		fmt.Println(pkgJSON)
		return &CreateIssueResult{}, nil
	}

	// 创建 Gist
	gistFiles := map[string]string{
		"issue-payload.json": pkgJSON,
	}
	for _, att := range opts.Attachments {
		gistFiles[att.Name] = att.Content
	}

	gist, err := s.client.CreateGist(
		fmt.Sprintf("[%s] %s", opts.Type, opts.Title),
		false,
		gistFiles,
	)
	if err != nil {
		return nil, fmt.Errorf("创建 Gist 失败: %w", err)
	}

	// 构建 Issue Body
	body := buildIssueBody(opts.Type, opts.Title, gist.HTMLURL)

	// 创建 Issue
	labels := []string{LabelCursorToolset, LabelPending, string(opts.Type)}
	issue, err := s.client.CreateIssue(owner, repo, opts.Title, body, labels)
	if err != nil {
		return nil, fmt.Errorf("创建 Issue 失败: %w", err)
	}

	return &CreateIssueResult{
		IssueURL: issue.HTMLURL,
		GistURL:  gist.HTMLURL,
		IssueNum: issue.Number,
	}, nil
}

// ListOptions 列出 Issue 的选项
type ListOptions struct {
	Repo   string
	Status string // pending, processing, processed, all
	Type   string
	Limit  int
}

// IssueInfo Issue 信息
type IssueInfo struct {
	Number    int
	Title     string
	Type      string
	Status    string
	CreatedAt string
	URL       string
}

// List 列出 Issue
func (s *IssueService) List(opts ListOptions) ([]IssueInfo, error) {
	owner, repo, err := parseRepo(opts.Repo)
	if err != nil {
		return nil, err
	}

	// 构建标签过滤
	labels := []string{LabelCursorToolset}
	if opts.Status != "" && opts.Status != "all" {
		labels = append(labels, opts.Status)
	}
	if opts.Type != "" {
		labels = append(labels, opts.Type)
	}

	state := "open"
	if opts.Status == "processed" || opts.Status == "rejected" {
		state = "closed"
	} else if opts.Status == "all" {
		state = "all"
	}

	issues, err := s.client.ListIssues(owner, repo, labels, state, opts.Limit)
	if err != nil {
		return nil, err
	}

	var result []IssueInfo
	for _, issue := range issues {
		info := IssueInfo{
			Number:    issue.Number,
			Title:     issue.Title,
			CreatedAt: issue.CreatedAt[:10],
			URL:       issue.HTMLURL,
		}

		// 提取类型和状态
		for _, label := range issue.Labels {
			switch label.Name {
			case LabelFeatureRequest, LabelBugReport, LabelPackRegister, LabelPackSync:
				info.Type = label.Name
			case LabelPending, LabelProcessing, LabelProcessed, LabelRejected:
				info.Status = label.Name
			}
		}

		result = append(result, info)
	}

	return result, nil
}

// GetResult 获取 Issue 的结果
type GetResult struct {
	Issue   *github.Issue
	Package *models.IssuePackage
}

// Get 获取并解析 Issue
func (s *IssueService) Get(repoStr string, number int) (*GetResult, error) {
	owner, repo, err := parseRepo(repoStr)
	if err != nil {
		return nil, err
	}

	issue, err := s.client.GetIssue(owner, repo, number)
	if err != nil {
		return nil, err
	}

	// 从 body 中提取 Gist URL
	gistURL := extractGistURL(issue.Body)
	if gistURL == "" {
		return &GetResult{Issue: issue}, nil
	}

	// 提取 Gist ID
	gistID := extractGistID(gistURL)
	if gistID == "" {
		return &GetResult{Issue: issue}, nil
	}

	// 获取 Gist 内容
	gist, err := s.client.GetGist(gistID)
	if err != nil {
		return &GetResult{Issue: issue}, nil
	}

	// 解析 payload
	if file, ok := gist.Files["issue-payload.json"]; ok {
		pkg, err := models.ParseIssuePackage(file.Content)
		if err == nil {
			return &GetResult{Issue: issue, Package: pkg}, nil
		}
	}

	return &GetResult{Issue: issue}, nil
}

// UpdateStatus 更新 Issue 状态
func (s *IssueService) UpdateStatus(repoStr string, number int, status string, comment string) error {
	owner, repo, err := parseRepo(repoStr)
	if err != nil {
		return err
	}

	// 获取当前 Issue
	issue, err := s.client.GetIssue(owner, repo, number)
	if err != nil {
		return err
	}

	// 更新标签：移除旧状态，添加新状态
	var newLabels []string
	for _, label := range issue.Labels {
		if label.Name != LabelPending && label.Name != LabelProcessing &&
			label.Name != LabelProcessed && label.Name != LabelRejected {
			newLabels = append(newLabels, label.Name)
		}
	}
	newLabels = append(newLabels, status)

	_, err = s.client.UpdateIssue(owner, repo, number, "", newLabels)
	if err != nil {
		return err
	}

	// 添加评论
	if comment != "" {
		err = s.client.AddComment(owner, repo, number, comment)
		if err != nil {
			return fmt.Errorf("添加评论失败: %w", err)
		}
	}

	return nil
}

// Close 关闭 Issue
func (s *IssueService) Close(repoStr string, number int, result string, comment string) error {
	owner, repo, err := parseRepo(repoStr)
	if err != nil {
		return err
	}

	// 获取当前 Issue
	issue, err := s.client.GetIssue(owner, repo, number)
	if err != nil {
		return err
	}

	// 确定最终状态标签
	statusLabel := LabelProcessed
	if result == "rejected" {
		statusLabel = LabelRejected
	}

	// 更新标签
	var newLabels []string
	for _, label := range issue.Labels {
		if label.Name != LabelPending && label.Name != LabelProcessing &&
			label.Name != LabelProcessed && label.Name != LabelRejected {
			newLabels = append(newLabels, label.Name)
		}
	}
	newLabels = append(newLabels, statusLabel)

	// 关闭 Issue 并更新标签
	_, err = s.client.UpdateIssue(owner, repo, number, "closed", newLabels)
	if err != nil {
		return err
	}

	// 添加评论
	if comment != "" {
		err = s.client.AddComment(owner, repo, number, comment)
		if err != nil {
			return fmt.Errorf("添加评论失败: %w", err)
		}
	}

	return nil
}

// parseRepo 解析仓库字符串 "owner/repo"
func parseRepo(repo string) (string, string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("无效的仓库格式，应为 owner/repo")
	}
	return parts[0], parts[1], nil
}

// buildIssueBody 构建 Issue Body
func buildIssueBody(issueType models.IssueType, title string, gistURL string) string {
	return fmt.Sprintf(`## %s: %s

**Type:** %s
**Created by:** github-issue-pack v0.1.0

### Details

📦 [View full payload](%s)

---
<sub>This issue was automatically created by [github-issue-pack](https://github.com/shichao402/github-issue-pack)</sub>
`, issueType, title, issueType, gistURL)
}

// extractGistURL 从 Issue body 中提取 Gist URL
func extractGistURL(body string) string {
	re := regexp.MustCompile(`https://gist\.github\.com/[a-zA-Z0-9_-]+/[a-f0-9]+`)
	match := re.FindString(body)
	return match
}

// extractGistID 从 Gist URL 中提取 ID
func extractGistID(url string) string {
	re := regexp.MustCompile(`/([a-f0-9]+)$`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
