package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Issue GitHub Issue 数据结构
type Issue struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	State     string   `json:"state"`
	HTMLURL   string   `json:"html_url"`
	Labels    []Label  `json:"labels"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	User      User     `json:"user"`
}

// Label 标签
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// User 用户
type User struct {
	Login string `json:"login"`
}

// CreateIssueRequest 创建 Issue 请求
type CreateIssueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels,omitempty"`
}

// UpdateIssueRequest 更新 Issue 请求
type UpdateIssueRequest struct {
	State  string   `json:"state,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

// CreateIssue 创建 Issue
func (c *Client) CreateIssue(owner, repo, title, body string, labels []string) (*Issue, error) {
	req := CreateIssueRequest{
		Title:  title,
		Body:   body,
		Labels: labels,
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues", baseURL, owner, repo)
	respBody, err := c.Post(url, req)
	if err != nil {
		return nil, fmt.Errorf("创建 Issue 失败: %w", err)
	}

	var issue Issue
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return nil, fmt.Errorf("解析 Issue 响应失败: %w", err)
	}

	return &issue, nil
}

// GetIssue 获取 Issue
func (c *Client) GetIssue(owner, repo string, number int) (*Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", baseURL, owner, repo, number)
	respBody, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取 Issue 失败: %w", err)
	}

	var issue Issue
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return nil, fmt.Errorf("解析 Issue 响应失败: %w", err)
	}

	return &issue, nil
}

// ListIssues 列出 Issue
func (c *Client) ListIssues(owner, repo string, labels []string, state string, limit int) ([]Issue, error) {
	params := url.Values{}
	if len(labels) > 0 {
		params.Set("labels", strings.Join(labels, ","))
	}
	if state != "" {
		params.Set("state", state)
	}
	if limit > 0 {
		params.Set("per_page", fmt.Sprintf("%d", limit))
	}

	apiURL := fmt.Sprintf("%s/repos/%s/%s/issues?%s", baseURL, owner, repo, params.Encode())
	respBody, err := c.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("列出 Issue 失败: %w", err)
	}

	var issues []Issue
	if err := json.Unmarshal(respBody, &issues); err != nil {
		return nil, fmt.Errorf("解析 Issue 列表失败: %w", err)
	}

	return issues, nil
}

// UpdateIssue 更新 Issue
func (c *Client) UpdateIssue(owner, repo string, number int, state string, labels []string) (*Issue, error) {
	req := UpdateIssueRequest{
		State:  state,
		Labels: labels,
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", baseURL, owner, repo, number)
	respBody, err := c.Patch(url, req)
	if err != nil {
		return nil, fmt.Errorf("更新 Issue 失败: %w", err)
	}

	var issue Issue
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return nil, fmt.Errorf("解析 Issue 响应失败: %w", err)
	}

	return &issue, nil
}

// AddComment 添加评论
func (c *Client) AddComment(owner, repo string, number int, body string) error {
	req := map[string]string{"body": body}
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", baseURL, owner, repo, number)
	_, err := c.Post(url, req)
	if err != nil {
		return fmt.Errorf("添加评论失败: %w", err)
	}
	return nil
}

// CloseIssue 关闭 Issue
func (c *Client) CloseIssue(owner, repo string, number int) (*Issue, error) {
	return c.UpdateIssue(owner, repo, number, "closed", nil)
}
