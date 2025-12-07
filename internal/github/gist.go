package github

import (
	"encoding/json"
	"fmt"
)

// GistFile Gist 文件
type GistFile struct {
	Content string `json:"content"`
}

// Gist Gist 数据结构
type Gist struct {
	ID          string              `json:"id,omitempty"`
	Description string              `json:"description,omitempty"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
	HTMLURL     string              `json:"html_url,omitempty"`
	CreatedAt   string              `json:"created_at,omitempty"`
}

// CreateGistRequest 创建 Gist 请求
type CreateGistRequest struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

// CreateGist 创建 Gist
func (c *Client) CreateGist(description string, public bool, files map[string]string) (*Gist, error) {
	gistFiles := make(map[string]GistFile)
	for name, content := range files {
		gistFiles[name] = GistFile{Content: content}
	}

	req := CreateGistRequest{
		Description: description,
		Public:      public,
		Files:       gistFiles,
	}

	respBody, err := c.Post(baseURL+"/gists", req)
	if err != nil {
		return nil, fmt.Errorf("创建 Gist 失败: %w", err)
	}

	var gist Gist
	if err := json.Unmarshal(respBody, &gist); err != nil {
		return nil, fmt.Errorf("解析 Gist 响应失败: %w", err)
	}

	return &gist, nil
}

// GetGist 获取 Gist
func (c *Client) GetGist(gistID string) (*Gist, error) {
	respBody, err := c.Get(baseURL + "/gists/" + gistID)
	if err != nil {
		return nil, fmt.Errorf("获取 Gist 失败: %w", err)
	}

	var gist Gist
	if err := json.Unmarshal(respBody, &gist); err != nil {
		return nil, fmt.Errorf("解析 Gist 响应失败: %w", err)
	}

	return &gist, nil
}
